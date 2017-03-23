package semantic

import (
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
//     { (conn data.Connection{FromPort, DataType, ShowDataType, ToPort, ToOp}?), (oper data.Operation{Name, Type, SrcPos, OutPorts}) },
//     { (bigIdentDataType string), (op data.Operation{Name, Type, SrcPos, InPorts, OutPorts}) }*,
//     (connection data.Connection{FromPort{}, DataType, ToPort})?
//   ]
// ]
//
// semantic result: (flow data.Flow{})

type SemanticConnections struct {
	/*
		createConns     *CreateConnections
		verifyOutPorts  *VerifyOutPortsUsedOnlyOnce
		handleChainBeg  *HandleChainBeg
		handleChainMids *HandleChainMids
		handleChainEnd  *HandleChainEnd
		begAddLastOp    *AddLastOp
		midAddLastOp    *AddLastOp
	*/
	InPort     func(interface{})
	SetOutPort func(func(interface{}))
}

func NewSemanticConnections() *SemanticConnections {
	f := &SemanticConnections{}
	/*
		f.createConns = NewCreateConnections()
		f.verifyOutPorts = NewVerifyOutPortsUsedOnlyOnce()
		f.handleChainBeg = NewHandleChainBeg()
		f.handleChainMids = NewHandleChainMids()
		f.handleChainEnd = NewHandleChainEnd()
		f.begAddLastOp = NewAddLastOp()
		f.midAddLastOp = NewAddLastOp()

		f.createConns.SetOutPort(f.verifyOutPorts.InPort)
		f.createConns.SetChainOutPort(f.handleChainBeg.InPort)
		f.handleChainBeg.SetOutPort(f.handleChainMids.InPort)
		f.handleChainMids.SetOutPort(f.handleChainEnd.InPort)
		f.handleChainEnd.SetOutPort(f.createConns.ChainInPort)
		f.handleChainBeg.SetAddOpOutPort(f.begAddLastOp.InPort)
		f.begAddLastOp.SetOutPort(f.handleChainBeg.AddOpInPort)
		f.handleChainMids.SetAddOpOutPort(f.midAddLastOp.InPort)
		f.midAddLastOp.SetOutPort(f.handleChainMids.AddOpInPort)

		f.InPort = f.createConns.InPort
		f.SetOutPort = f.verifyOutPorts.SetOutPort
	*/

	return f
}

// ------------ HandleChainBeg:
// semantic input: dat.chainBegOp, dat.chainBegConn plus dat.opMap has to be up to date
// semantic result: dat.conns is updated (if dat.chainBegConn != nil) plus dat.ops and dat.opMap are updated as necessary
type HandleChainBeg struct {
	//begAddLastOp *AddLastOp
	addOpOutPort func(*SemanticConnectionsData)
	outPort      func(*SemanticConnectionsData)
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

// ------------ AddLastOp:
// semantic input: dat.newOp, dat.opMap has to be up to date
// semantic result: dat.addOpResult = &AddOpResult{op *data.Operation, outPort *data.PortData, outPortAdded bool},
//		dat.ops and dat.opMap are updated as necessary
type AddLastOp struct {
	outPort func(*SemanticConnectionsData)
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
func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
