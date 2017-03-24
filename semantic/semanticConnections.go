package semantic

import (
	"strconv"

	"github.com/flowdev/gflowparser/data"
	"github.com/flowdev/gparselib"
)

type SemanticConnectionsData struct {
	mainData *data.MainData
	// stuff needed for flow:
	ops   []*data.Operation
	conns []*data.Connection
	opMap map[string]*data.Operation
	// intermediate stuff:
	chainBegOp   *data.Operation
	chainBegConn *data.Connection
	chainMids    []interface{}
	chainEnd     *data.Connection
	addOpResult  *AddOpResult
	newOp        *data.Operation
}

type AddOpResult struct {
	op           *data.Operation
	outPort      *data.PortData
	outPortAdded bool
}

// text input:
// ( optInPort  [OptDataType]-> optInPort )? opName(OpType) optOutPort
// ( [OptDataType]-> optInPort opName(OpType) optOutPort )*
// ( [OptDataType]-> optOutPort )?
//
// semantic input:
// Multi1[
//   All[
//     [ (conn data.Connection{FromPort, DataType, ShowDataType, ToPort, ToOp}?), (oper data.Operation{Name, Type, SrcPos, OutPorts}) ],
//     [ (bigIdentDataType string), (op data.Operation{Name, Type, SrcPos, InPorts, OutPorts}) ]*,
//     (connection data.Connection{FromPort{}, DataType, ToPort})?
//   ]
// ]
//
// semantic result: (flow data.Flow{})
type SemanticConnections struct {
	//createConns     *CreateConnections
	//verifyOutPorts  *VerifyOutPortsUsedOnlyOnce
	//handleChainBeg  *HandleChainBeg
	//handleChainMids *HandleChainMids
	//handleChainEnd  *HandleChainEnd
	//begAddLastOp    *AddLastOp
	//midAddLastOp    *AddLastOp
	InPort     func(interface{})
	SetOutPort func(func(interface{}))
}

func NewSemanticConnections() *SemanticConnections {
	f := &SemanticConnections{}
	createConns := NewCreateConnections()
	verifyOutPorts := NewVerifyOutPortsUsedOnlyOnce()
	handleChainBeg := NewHandleChainBeg()
	handleChainMids := NewHandleChainMids()
	handleChainEnd := NewHandleChainEnd()
	begAddLastOp := NewAddLastOp()
	midAddLastOp := NewAddLastOp()

	createConns.SetOutPort(verifyOutPorts.InPort)
	createConns.SetChainOutPort(handleChainBeg.InPort)
	handleChainBeg.SetOutPort(handleChainMids.InPort)
	handleChainMids.SetOutPort(handleChainEnd.InPort)
	handleChainEnd.SetOutPort(createConns.ChainInPort)
	handleChainBeg.SetAddOpOutPort(begAddLastOp.InPort)
	begAddLastOp.SetOutPort(handleChainBeg.AddOpInPort)
	handleChainMids.SetAddOpOutPort(midAddLastOp.InPort)
	midAddLastOp.SetOutPort(handleChainMids.AddOpInPort)

	f.InPort = createConns.InPort
	f.SetOutPort = verifyOutPorts.SetOutPort

	return f
}

// ------------ CreateConnections:
// semantic input:
// Multi1[
//   All[
//     [ (conn data.Connection{FromPort, DataType, ShowDataType, ToPort, ToOp}?), (oper data.Operation{Name, Type, SrcPos, OutPorts}) ],
//     [ (bigIdentDataType string), (op data.Operation{Name, Type, SrcPos, InPorts, OutPorts}) ]*,
//     (connection data.Connection{FromPort{}, DataType, ToPort})?
//   ]
// ]
//
// semantic result: (flow data.Flow{Conns, Ops})
type CreateConnections struct {
	chainOutPort func(*SemanticConnectionsData)
	outPort      func(interface{})
}

func NewCreateConnections() *CreateConnections {
	return &CreateConnections{}
}
func (op *CreateConnections) InPort(dat interface{}) {
	md := dat.(*data.MainData)
	connsData := &SemanticConnectionsData{
		mainData: md,
		ops:      make([]*data.Operation, 0, 32),
		conns:    make([]*data.Connection, 0, 32),
		opMap:    make(map[string]*data.Operation),
	}

	for _, subResult := range md.ParseData.SubResults {
		chain := subResult.Value.([]interface{})
		chainBeg := chain[0].([]interface{})
		connsData.addOpResult = &AddOpResult{}
		connsData.chainBegConn = chainBeg[0].(*data.Connection)
		connsData.chainBegOp = chainBeg[1].(*data.Operation)
		connsData.chainMids = chain[1].([]interface{})
		connsData.chainEnd = chain[2].(*data.Connection)

		op.chainOutPort(connsData)
	}
	md.ParseData.Result.Value = nil

	if md.ParseData.Result.ErrPos < 0 {
		flow := &data.Flow{
			Conns: connsData.conns[:],
			Ops:   connsData.ops[:],
		}
		md.ParseData.Result.Value = flow
	}
	op.outPort(md)
}
func (op *CreateConnections) ChainInPort(dat *SemanticConnectionsData) {
	// WARNING: We make use of the knowledge that all calls in this subflow (package semantic) are synchronous!
	// On the other hand this means that this method needn't do anything at all and we save quite some stack space. :-)
}
func (op *CreateConnections) SetChainOutPort(port func(*SemanticConnectionsData)) {
	op.chainOutPort = port
}
func (op *CreateConnections) SetOutPort(port func(interface{})) {
	op.outPort = port
}

// ------------ VerifyOutPortsUsedOnlyOnce:
// semantic input: (flow data.Flow{Conns, Ops})
// semantic result: (flow data.Flow{Conns, Ops}) only possible changes are added errors
type VerifyOutPortsUsedOnlyOnce struct {
	outPort func(interface{})
}

func NewVerifyOutPortsUsedOnlyOnce() *VerifyOutPortsUsedOnlyOnce {
	return &VerifyOutPortsUsedOnlyOnce{}
}
func (op *VerifyOutPortsUsedOnlyOnce) InPort(dat interface{}) {
	md := dat.(*data.MainData)
	flow := md.ParseData.Result.Value.(*data.Flow)
	if flow == nil {
		op.outPort(md)
		return
	}
	// check for output ports that are connected to multiple input ports:
	connMap := make(map[string]map[string]int)
	for _, conn := range flow.Conns {
		fromPort := describePort(conn.FromOp, conn.FromPort)
		toPort := describePort(conn.ToOp, conn.ToPort)
		toPorts, ok := connMap[fromPort]
		if !ok {
			toPorts = make(map[string]int)
			connMap[fromPort] = toPorts
		}
		toPorts[toPort] = max(toPorts[toPort], conn.ToPort.SrcPos)
	}

	for fromPort, toPorts := range connMap {
		if len(toPorts) > 1 {
			gparselib.AddError(lastSrcPos(toPorts), "The output port '"+fromPort+"' is connected to multiple input ports: "+enumeratePorts(toPorts), nil, md.ParseData)
		}
	}
	op.outPort(md)
}
func (op *VerifyOutPortsUsedOnlyOnce) SetOutPort(port func(interface{})) {
	op.outPort = port
}

// ------------ HandleChainBeg:
// semantic input: dat.chainBegOp, dat.chainBegConn plus dat.opMap has to be up to date
// semantic result: dat.addOpResult, dat.conns is updated (if dat.chainBegConn != nil)
//                  plus dat.ops and dat.opMap are updated as necessary
type HandleChainBeg struct {
	addOpOutPort func(*SemanticConnectionsData)
	outPort      func(*SemanticConnectionsData)
}

func NewHandleChainBeg() *HandleChainBeg {
	return &HandleChainBeg{}
}
func (op *HandleChainBeg) InPort(dat *SemanticConnectionsData) {
	// first add the operation:
	dat.newOp = dat.chainBegOp
	op.addOpOutPort(dat)
}
func (op *HandleChainBeg) AddOpInPort(dat *SemanticConnectionsData) {
	lastOp := dat.addOpResult.op
	connBeg := dat.chainBegConn

	// now add the connection if it exists:
	if connBeg != nil {
		connBeg.ToOp = lastOp
		correctToPort(connBeg, lastOp)
		dat.conns = append(dat.conns, connBeg)
	}

	op.outPort(dat)
}
func (op *HandleChainBeg) SetAddOpOutPort(port func(*SemanticConnectionsData)) {
	op.addOpOutPort = port
}
func (op *HandleChainBeg) SetOutPort(port func(*SemanticConnectionsData)) {
	op.outPort = port
}

// ------------ HandleChainMids:
// semantic input: dat.chainMids, dat.addOpResult plus dat.opMap has to be up to date
// semantic result: dat.addOpResult, dat.conns is updated plus dat.ops and dat.opMap are updated as necessary
type HandleChainMids struct {
	addOpOutPort func(*SemanticConnectionsData)
	outPort      func(*SemanticConnectionsData)
}

func NewHandleChainMids() *HandleChainMids {
	return &HandleChainMids{}
}
func (op *HandleChainMids) InPort(dat *SemanticConnectionsData) {
	if len(dat.chainMids) <= 0 {
		op.outPort(dat)
		return
	}
	fromOp := dat.addOpResult.op

	for _, chainMidIf := range dat.chainMids {
		chainMid := chainMidIf.([]interface{})
		arrowType := chainMid[0].(string)
		fromPort := dat.addOpResult.outPort
		toOp := chainMid[1].(*data.Operation)
		toPort := toOp.InPorts[0] // TODO: Should this line be executed after the op is added???

		// add the operation:
		dat.newOp = toOp
		op.addOpOutPort(dat)
		toOp = dat.addOpResult.op

		// now add the connection:
		connMid := &data.Connection{
			DataType:     arrowType,
			ShowDataType: (arrowType != ""),
			FromOp:       fromOp,
			FromPort:     fromPort,
			ToOp:         toOp,
			ToPort:       toPort,
		}
		correctToPort(connMid, toOp)
		dat.conns = append(dat.conns, connMid)

		fromOp = toOp
	}

	op.outPort(dat)
}
func (op *HandleChainMids) AddOpInPort(dat *SemanticConnectionsData) {
	// WARNING: We make use of the knowledge that all calls in this subflow (package semantic) are synchronous!
	// On the other hand this means that this method needn't do anything at all and we save quite some stack space. :-)
}
func (op *HandleChainMids) SetAddOpOutPort(port func(*SemanticConnectionsData)) {
	op.addOpOutPort = port
}
func (op *HandleChainMids) SetOutPort(port func(*SemanticConnectionsData)) {
	op.outPort = port
}

// ------------ HandleChainEnd:
// semantic input: dat.chainEnd, dat.addOpResult plus dat.opMap has to be up to date
// semantic result: dat.addOpResult, dat.conns is updated plus dat.ops and dat.opMap are updated as necessary
type HandleChainEnd struct {
	outPort func(*SemanticConnectionsData)
}

func NewHandleChainEnd() *HandleChainEnd {
	return &HandleChainEnd{}
}
func (op *HandleChainEnd) InPort(dat *SemanticConnectionsData) {
	chainEnd := dat.chainEnd
	addOpResult := dat.addOpResult
	lastOp := addOpResult.op

	if chainEnd != nil {
		chainEnd.FromOp = lastOp
		if addOpResult.outPort != nil {
			chainEnd.FromPort = addOpResult.outPort
		}
		if chainEnd.FromPort.Name == "" && chainEnd.ToPort.Name == "" {
			chainEnd.FromPort = data.DefaultOutPort(chainEnd.FromPort.SrcPos)
			chainEnd.ToPort = data.CopyPort(chainEnd.FromPort, chainEnd.ToPort.SrcPos)
		} else if chainEnd.ToPort.Name == "" {
			chainEnd.ToPort = data.CopyPort(chainEnd.FromPort, chainEnd.ToPort.SrcPos)
		} else if chainEnd.FromPort.Name == "" {
			chainEnd.FromPort = data.DefaultOutPort(chainEnd.FromPort.SrcPos)
			addPort(lastOp, chainEnd.FromPort, dat.mainData.ParseData, addOpResult)
		}
		correctFromPort(chainEnd, lastOp)
		dat.conns = append(dat.conns, chainEnd)
	} else if addOpResult.outPortAdded {
		lastOp.OutPorts = lastOp.OutPorts[:len(lastOp.OutPorts)-1]
	}

	op.outPort(dat)
}
func (op *HandleChainEnd) SetOutPort(port func(*SemanticConnectionsData)) {
	op.outPort = port
}

// ------------ AddLastOp:
// semantic input: dat.newOp, dat.opMap has to be up to date
// semantic result: dat.addOpResult = &AddOpResult{op *data.Operation, outPort *data.PortData, outPortAdded bool},
//		dat.ops and dat.opMap are updated as necessary
type AddLastOp struct {
	outPort func(*SemanticConnectionsData)
}

func NewAddLastOp() *AddLastOp {
	return &AddLastOp{}
}
func (op *AddLastOp) InPort(dat *SemanticConnectionsData) {
	newOp := dat.newOp
	result := &AddOpResult{}
	dat.addOpResult = result

	existingOp := dat.opMap[newOp.Name]
	if existingOp != nil {
		updateExistingOp(existingOp, newOp, dat.mainData.ParseData, result)
	} else {
		addNewOp(dat, newOp, result)
	}

	op.outPort(dat)
}
func (op *AddLastOp) SetOutPort(port func(*SemanticConnectionsData)) {
	op.outPort = port
}

func updateExistingOp(existingOp *data.Operation, newOp *data.Operation, pd *gparselib.ParseData, result *AddOpResult) {
	if existingOp.Type == "" {
		existingOp.Type = newOp.Type
	} else if newOp.Type != "" && newOp.Type != existingOp.Type {
		gparselib.AddError(newOp.SrcPos, "The operation '"+newOp.Name+"' has got two different types '"+existingOp.Type+"' and '"+newOp.Type+"'!", nil, pd)
	}
	if len(newOp.InPorts) > 0 {
		addPort(existingOp, newOp.InPorts[0], pd, nil)
	}
	if len(newOp.OutPorts) > 0 {
		addPort(existingOp, newOp.OutPorts[0], pd, result)
	}
	result.op = existingOp
}

func addNewOp(dat *SemanticConnectionsData, newOp *data.Operation, result *AddOpResult) {
	dat.opMap[newOp.Name] = newOp
	dat.ops = append(dat.ops, newOp)
	result.op = newOp
	if len(newOp.OutPorts) > 0 {
		result.outPort = newOp.OutPorts[0]
		result.outPortAdded = true
	}
}

// Utility functions
const (
	PortsEqual = iota
	PortsDiffer
	PortsConflict
)

func addPort(op *data.Operation, newPort *data.PortData, pd *gparselib.ParseData, result *AddOpResult) {
	var ports []*data.PortData
	var portType string
	if result == nil {
		ports = op.InPorts
		portType = "input"
	} else {
		ports = op.OutPorts
		portType = "output"
	}

	for _, oldPort := range ports {
		c := comparePorts(oldPort, newPort)
		if c == PortsConflict {
			gparselib.AddError(max(newPort.SrcPos, oldPort.SrcPos),
				"The "+portType+" port '"+newPort.Name+"' of the operation '"+op.Name+"' is used as indexed and unindexed port in the same flow!", nil, pd)
			return
		}
		if c == PortsEqual {
			if result != nil {
				result.outPort = oldPort
				result.outPortAdded = false
			}
			return
		}
	}

	if result == nil {
		op.InPorts = append(ports, newPort)
	} else {
		op.OutPorts = append(ports, newPort)
		result.outPort = newPort
		result.outPortAdded = true
	}
}
func correctFromPort(conn *data.Connection, op *data.Operation) {
	for _, p := range op.OutPorts {
		c := comparePorts(p, conn.FromPort)
		if c != PortsDiffer {
			conn.FromPort = p
			break
		}
	}
}
func correctToPort(conn *data.Connection, op *data.Operation) {
	for _, p := range op.InPorts {
		c := comparePorts(p, conn.ToPort)
		if c != PortsDiffer {
			conn.ToPort = p
		}
	}
}
func comparePorts(p1, p2 *data.PortData) int {
	if p1 == nil && p2 == nil {
		return PortsEqual
	}
	if p1 == nil || p2 == nil {
		return PortsDiffer
	}
	if p1.Name != p2.Name {
		return PortsDiffer
	}
	if p1.HasIndex != p2.HasIndex {
		return PortsConflict
	}
	if !p1.HasIndex {
		return PortsEqual
	}
	if p1.Index == p2.Index {
		return PortsEqual
	}
	return PortsDiffer
}
func describePort(op *data.Operation, port *data.PortData) string {
	var desc string
	if op == nil {
		desc = "<FLOW>"
	} else {
		desc = op.Name
	}
	desc += ":" + port.Name

	if port.HasIndex {
		desc += "." + strconv.Itoa(port.Index)
	}
	return desc
}
func lastSrcPos(pMap map[string]int) int {
	lastPos := -1
	for _, pos := range pMap {
		if pos > lastPos {
			lastPos = pos
		}
	}
	return lastPos
}
func enumeratePorts(pMap map[string]int) string {
	enum := ""
	for p := range pMap {
		if enum != "" {
			enum += ", "
		}
		enum += p
	}
	return enum
}
func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
