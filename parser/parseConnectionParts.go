package parser

import (
	"math"
	"strconv"
	"strings"

	"github.com/flowdev/gflowparser/data"
	"github.com/flowdev/gparselib"
)

// ------------ ParseConnectionPart:
// semantic result: op data.Operation{+InPorts, +OutPorts}
type SemanticConnectionPart struct {
	outPort func(interface{})
}

func NewSemanticConnectionPart() *SemanticConnectionPart {
	return &SemanticConnectionPart{}
}
func (op *SemanticConnectionPart) InPort(dat interface{}) {
	md := dat.(*data.MainData)
	res := md.ParseData.Result
	subVals := res.Value.([]interface{})

	var port1 *data.PortData
	oper := subVals[1].(*data.Operation)
	var port2 *data.PortData

	if subVals[0] == nil {
		port1 = data.DefaultInPort(md.ParseData.SubResults[0].Pos)
	} else {
		port1 = subVals[0].(*data.PortData)
	}
	if subVals[2] == nil {
		port2 = data.DefaultOutPort(md.ParseData.SubResults[2].Pos)
	} else {
		port2 = subVals[2].(*data.PortData)
	}
	oper.InPorts = append(oper.InPorts, port1)
	oper.OutPorts = append(oper.OutPorts, port2)
	res.Value = oper

	op.outPort(md)
}
func (op *SemanticConnectionPart) SetOutPort(port func(interface{})) {
	op.outPort = port
}

type ParseConnectionPart struct {
	//connPart     *gparselib.ParseAll
	//semantic     *SemanticConnectionPart
	//optInPort    *ParseOptPortSpc
	//opNameParens *ParseOperationNameParens
	//optOutPort   *ParseOptPort
	InPort     func(interface{})
	SetOutPort func(func(interface{}))
}

func NewParseConnectionPart() *ParseConnectionPart {
	f := &ParseConnectionPart{}
	connPart := gparselib.NewParseAll(parseData, setParseData)
	semantic := NewSemanticConnectionPart()
	optInPort := NewParseOptPortSpc()
	opNameParens := NewParseOperationNameParens()
	optOutPort := NewParseOptPort()

	connPart.SetSemOutPort(semantic.InPort)
	semantic.SetOutPort(connPart.SemInPort)
	connPart.AppendSubOutPort(optInPort.InPort)
	optInPort.SetOutPort(connPart.SubInPort)
	connPart.AppendSubOutPort(opNameParens.InPort)
	opNameParens.SetOutPort(connPart.SubInPort)
	connPart.AppendSubOutPort(optOutPort.InPort)
	optOutPort.SetOutPort(connPart.SubInPort)

	f.InPort = connPart.InPort
	f.SetOutPort = connPart.SetOutPort

	return f
}

// ------------ ParseOperationNameParens:
// semantic result: op data.Operation{Name, Type, SrcPos}
type SemanticOperationNameParens struct {
	outPort func(interface{})
}

func NewSemanticOperationNameParens() *SemanticOperationNameParens {
	return &SemanticOperationNameParens{}
}
func (op *SemanticOperationNameParens) InPort(dat interface{}) {
	md := dat.(*data.MainData)
	res := md.ParseData.Result
	subRes := md.ParseData.SubResults
	subVals := res.Value.([]interface{})
	oper := &data.Operation{}

	if subVals[0] != nil {
		opNameVal := subVals[0].([]interface{})
		oper.Name = opNameVal[0].(string)
		oper.SrcPos = subRes[0].Pos
	}
	if subVals[3] != nil {
		oper.Type = subVals[3].(string)
		if subVals[0] == nil {
			oper.SrcPos = subRes[1].Pos
		}
	}
	if len(oper.Name) <= 0 && len(oper.Type) <= 0 {
		gparselib.AddError(subRes[0].Pos, "At least an operation name or an operation type have to be provided",
			nil, md.ParseData)
	} else if len(oper.Name) <= 0 {
		oper.Name = strings.ToLower(oper.Type[0:1]) + oper.Type[1:]
	}
	res.Value = oper

	op.outPort(md)
}
func (op *SemanticOperationNameParens) SetOutPort(port func(interface{})) {
	op.outPort = port
}

type ParseOperationNameParens struct {
	//opNameParens *gparselib.ParseAll
	//semantic     *SemanticOperationNameParens
	//optOpName    *gparselib.ParseOptional
	//openType     *gparselib.ParseLiteral
	//spc1         *ParseOptSpc
	//optOpType    *ParseOptOperationType
	//closeType    *gparselib.ParseLiteral
	//spc2         *ParseOptSpc
	//opName       *gparselib.ParseAll
	//smallIdent   *ParseSmallIdent
	//spc3         *ParseOptSpc
	InPort     func(interface{})
	SetOutPort func(func(interface{}))
}

func NewParseOperationNameParens() *ParseOperationNameParens {
	f := &ParseOperationNameParens{}
	opNameParens := gparselib.NewParseAll(parseData, setParseData)
	semantic := NewSemanticOperationNameParens()
	optOpName := gparselib.NewParseOptional(parseData, setParseData)
	openType := gparselib.NewParseLiteral(parseData, setParseData, "(")
	spc1 := NewParseOptSpc()
	optOpType := NewParseOptOperationType()
	closeType := gparselib.NewParseLiteral(parseData, setParseData, ")")
	spc2 := NewParseOptSpc()
	opName := gparselib.NewParseAll(parseData, setParseData)
	smallIdent := NewParseSmallIdent()
	spc3 := NewParseOptSpc()

	opNameParens.SetSemOutPort(semantic.InPort)
	semantic.SetOutPort(opNameParens.SemInPort)
	opNameParens.AppendSubOutPort(optOpName.InPort)
	optOpName.SetOutPort(opNameParens.SubInPort)
	opNameParens.AppendSubOutPort(openType.InPort)
	openType.SetOutPort(opNameParens.SubInPort)
	opNameParens.AppendSubOutPort(spc1.InPort)
	spc1.SetOutPort(opNameParens.SubInPort)
	opNameParens.AppendSubOutPort(optOpType.InPort)
	optOpType.SetOutPort(opNameParens.SubInPort)
	opNameParens.AppendSubOutPort(closeType.InPort)
	closeType.SetOutPort(opNameParens.SubInPort)
	opNameParens.AppendSubOutPort(spc2.InPort)
	spc2.SetOutPort(opNameParens.SubInPort)
	optOpName.SetSubOutPort(opName.InPort)
	opName.SetOutPort(optOpName.SubInPort)
	opName.AppendSubOutPort(smallIdent.InPort)
	smallIdent.SetOutPort(opName.SubInPort)
	opName.AppendSubOutPort(spc3.InPort)
	spc3.SetOutPort(opName.SubInPort)

	f.InPort = opNameParens.InPort
	f.SetOutPort = opNameParens.SetOutPort

	return f
}

// ------------ ParseOptOperationType:
// semantic result: bigIdentOperationType string
type SemanticOptOperationType struct {
	outPort func(interface{})
}

func NewSemanticOptOperationType() *SemanticOptOperationType {
	return &SemanticOptOperationType{}
}
func (op *SemanticOptOperationType) InPort(dat interface{}) {
	md := dat.(*data.MainData)
	res := md.ParseData.Result
	subVals := res.Value.([]interface{})
	res.Value = subVals[0]

	op.outPort(md)
}
func (op *SemanticOptOperationType) SetOutPort(port func(interface{})) {
	op.outPort = port
}

type ParseOptOperationType struct {
	//optOpType     *gparselib.ParseOptional
	//opType        *gparselib.ParseAll
	//semantic      *SemanticOptOperationType
	//parseBigIdent *ParseBigIdent
	//parseOptSpc   *ParseOptSpc
	InPort     func(interface{})
	SetOutPort func(func(interface{}))
}

func NewParseOptOperationType() *ParseOptOperationType {
	f := &ParseOptOperationType{}
	optOpType := gparselib.NewParseOptional(parseData, setParseData)
	opType := gparselib.NewParseAll(parseData, setParseData)
	semantic := NewSemanticOptOperationType()
	parseBigIdent := NewParseBigIdent()
	parseOptSpc := NewParseOptSpc()

	optOpType.SetSubOutPort(opType.InPort)
	opType.SetOutPort(optOpType.SubInPort)
	opType.SetSemOutPort(semantic.InPort)
	semantic.SetOutPort(opType.SemInPort)
	opType.AppendSubOutPort(parseBigIdent.InPort)
	parseBigIdent.SetOutPort(opType.SubInPort)
	opType.AppendSubOutPort(parseOptSpc.InPort)
	parseOptSpc.SetOutPort(opType.SubInPort)

	f.InPort = optOpType.InPort
	f.SetOutPort = optOpType.SetOutPort

	return f
}

// ------------ ParseArrow:
// semantic result: bigIdentDataType string
type SemanticArrow struct {
	outPort func(interface{})
}

func NewSemanticArrow() *SemanticArrow {
	return &SemanticArrow{}
}
func (op *SemanticArrow) InPort(dat interface{}) {
	md := dat.(*data.MainData)
	res := md.ParseData.Result
	subVals := res.Value.([]interface{})
	if subVals[1] != nil { // type exists
		subSubVals := subVals[1].([]interface{})
		res.Value = subSubVals[2]
	} else {
		res.Value = ""
	}

	op.outPort(md)
}
func (op *SemanticArrow) SetOutPort(port func(interface{})) {
	op.outPort = port
}

type ParseArrow struct {
	//arrow      *gparselib.ParseAll
	//semantic   *SemanticArrow
	//spcCom1    *ParseSpaceComment
	//optType    *gparselib.ParseOptional
	//litArr     *gparselib.ParseLiteral
	//spcCom2    *ParseSpaceComment
	//typ        *gparselib.ParseAll
	//openType   *gparselib.ParseLiteral
	//spc1       *ParseOptSpc
	//typeName   *ParseBigIdent
	//spc2       *ParseOptSpc
	//closeType  *gparselib.ParseLiteral
	InPort     func(interface{})
	SetOutPort func(func(interface{}))
}

func NewParseArrow() *ParseArrow {
	f := &ParseArrow{}
	arrow := gparselib.NewParseAll(parseData, setParseData)
	semantic := NewSemanticArrow()
	spcCom1 := NewParseSpaceComment()
	optType := gparselib.NewParseOptional(parseData, setParseData)
	litArr := gparselib.NewParseLiteral(parseData, setParseData, "->")
	spcCom2 := NewParseSpaceComment()
	typ := gparselib.NewParseAll(parseData, setParseData)
	openType := gparselib.NewParseLiteral(parseData, setParseData, "[")
	spc1 := NewParseOptSpc()
	typeName := NewParseBigIdent()
	spc2 := NewParseOptSpc()
	closeType := gparselib.NewParseLiteral(parseData, setParseData, "]")

	arrow.SetSemOutPort(semantic.InPort)
	semantic.SetOutPort(arrow.SemInPort)
	arrow.AppendSubOutPort(spcCom1.InPort)
	spcCom1.SetOutPort(arrow.SubInPort)
	arrow.AppendSubOutPort(optType.InPort)
	optType.SetOutPort(arrow.SubInPort)
	arrow.AppendSubOutPort(litArr.InPort)
	litArr.SetOutPort(arrow.SubInPort)
	arrow.AppendSubOutPort(spcCom2.InPort)
	spcCom2.SetOutPort(arrow.SubInPort)
	optType.SetSubOutPort(typ.InPort)
	typ.SetOutPort(optType.SubInPort)
	typ.AppendSubOutPort(openType.InPort)
	openType.SetOutPort(typ.SubInPort)
	typ.AppendSubOutPort(spc1.InPort)
	spc1.SetOutPort(typ.SubInPort)
	typ.AppendSubOutPort(typeName.InPort)
	typeName.SetOutPort(typ.SubInPort)
	typ.AppendSubOutPort(spc2.InPort)
	spc2.SetOutPort(typ.SubInPort)
	typ.AppendSubOutPort(closeType.InPort)
	closeType.SetOutPort(typ.SubInPort)

	f.InPort = arrow.InPort
	f.SetOutPort = arrow.SetOutPort

	return f
}

// ------------ ParseOptPortSpc:
// semantic result: port data.Port
type SemanticOptPortSpc struct {
	outPort func(interface{})
}

func NewSemanticOptPortSpc() *SemanticOptPortSpc {
	return &SemanticOptPortSpc{}
}
func (op *SemanticOptPortSpc) InPort(dat interface{}) {
	md := dat.(*data.MainData)
	res := md.ParseData.Result
	subVals := res.Value.([]interface{})
	res.Value = subVals[0]

	op.outPort(md)
}
func (op *SemanticOptPortSpc) SetOutPort(port func(interface{})) {
	op.outPort = port
}

type ParseOptPortSpc struct {
	//optPortSpc *gparselib.ParseOptional
	//portSpc    *gparselib.ParseAll
	//semantic   *SemanticOptPortSpc
	//pport      *ParsePort
	//space      *gparselib.ParseSpace
	InPort     func(interface{})
	SetOutPort func(func(interface{}))
}

func NewParseOptPortSpc() *ParseOptPortSpc {
	f := &ParseOptPortSpc{}
	optPortSpc := gparselib.NewParseOptional(parseData, setParseData)
	portSpc := gparselib.NewParseAll(parseData, setParseData)
	semantic := NewSemanticOptPortSpc()
	pport := NewParsePort()
	space := gparselib.NewParseSpace(parseData, setParseData, false)

	optPortSpc.SetSubOutPort(portSpc.InPort)
	portSpc.SetOutPort(optPortSpc.SubInPort)
	portSpc.SetSemOutPort(semantic.InPort)
	semantic.SetOutPort(portSpc.SemInPort)
	portSpc.AppendSubOutPort(pport.InPort)
	pport.SetOutPort(portSpc.SubInPort)
	portSpc.AppendSubOutPort(space.InPort)
	space.SetOutPort(portSpc.SubInPort)

	f.InPort = optPortSpc.InPort
	f.SetOutPort = optPortSpc.SetOutPort

	return f
}

// ------------ ParseOptPort:
// semantic result: port data.Port
type ParseOptPort struct {
	//optPort    *gparselib.ParseOptional
	//pport      *ParsePort
	InPort     func(interface{})
	SetOutPort func(func(interface{}))
}

func NewParseOptPort() *ParseOptPort {
	f := &ParseOptPort{}
	optPort := gparselib.NewParseOptional(parseData, setParseData)
	pport := NewParsePort()

	optPort.SetSubOutPort(pport.InPort)
	pport.SetOutPort(optPort.SubInPort)

	f.InPort = optPort.InPort
	f.SetOutPort = optPort.SetOutPort

	return f
}

// ------------ ParsePort:
// semantic result: port data.Port
type SemanticPort struct {
	outPort func(interface{})
}

func NewSemanticPort() *SemanticPort {
	return &SemanticPort{}
}
func (op *SemanticPort) InPort(dat interface{}) {
	md := dat.(*data.MainData)
	pd := md.ParseData
	nameRes := pd.SubResults[0]

	if pd.SubResults[1].Value == nil {
		md.ParseData.Result.Value = data.NewPort(nameRes.Text, nameRes.Pos)
	} else {
		val := pd.SubResults[1].Value
		idx64 := val.([]interface{})[1].(uint64)
		if idx64 > uint64(math.MaxInt32) {
			errPos := pd.SubResults[1].Pos + 1
			gparselib.AddError(errPos, "Ridiculous large port index "+strconv.FormatUint(idx64, 10), nil, pd)
			md.ParseData.Result.ErrPos = -1 // just a semantic error, no syntax error!
			md.ParseData.Result.Value = nil
		} else {
			md.ParseData.Result.Value = data.NewIdxPort(nameRes.Text, int(idx64), nameRes.Pos)
		}
	}
	op.outPort(md)
}
func (op *SemanticPort) SetOutPort(port func(interface{})) {
	op.outPort = port
}

type ParsePort struct {
	//port       *gparselib.ParseAll
	//semantic   *SemanticPort
	//portName   *ParseSmallIdent
	//optPortNum *gparselib.ParseOptional
	//portNum    *gparselib.ParseAll
	//dot        *gparselib.ParseLiteral
	//num        *gparselib.ParseNatural
	InPort     func(interface{})
	SetOutPort func(func(interface{}))
}

func NewParsePort() *ParsePort {
	f := &ParsePort{}
	port := gparselib.NewParseAll(parseData, setParseData)
	semantic := NewSemanticPort()
	portName := NewParseSmallIdent()
	optPortNum := gparselib.NewParseOptional(parseData, setParseData)
	portNum := gparselib.NewParseAll(parseData, setParseData)
	dot := gparselib.NewParseLiteral(parseData, setParseData, ".")
	num := gparselib.NewParseNatural(parseData, setParseData, 10)

	port.SetSemOutPort(semantic.InPort)
	semantic.SetOutPort(port.SemInPort)
	port.AppendSubOutPort(portName.InPort)
	portName.SetOutPort(port.SubInPort)
	port.AppendSubOutPort(optPortNum.InPort)
	optPortNum.SetOutPort(port.SubInPort)
	optPortNum.SetSubOutPort(portNum.InPort)
	portNum.SetOutPort(optPortNum.SubInPort)
	portNum.AppendSubOutPort(dot.InPort)
	dot.SetOutPort(portNum.SubInPort)
	portNum.AppendSubOutPort(num.InPort)
	num.SetOutPort(portNum.SubInPort)

	f.InPort = port.InPort
	f.SetOutPort = port.SetOutPort

	return f
}
