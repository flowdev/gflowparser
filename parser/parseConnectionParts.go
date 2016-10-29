package parser

import (
	"github.com/flowdev/gflowparser/data"
	"github.com/flowdev/gparselib"
	"math"
	"strconv"
	"strings"
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
	connPart     *gparselib.ParseAll
	semantic     *SemanticConnectionPart
	optInPort    *ParseOptPortSpc
	opNameParens *ParseOperationNameParens
	optOutPort   *ParseOptPort
	InPort       func(interface{})
	SetOutPort   func(func(interface{}))
}

func NewParseConnectionPart() *ParseConnectionPart {
	f := &ParseConnectionPart{}
	f.connPart = gparselib.NewParseAll(parseData, setParseData)
	f.semantic = NewSemanticConnectionPart()
	f.optInPort = NewParseOptPortSpc()
	f.opNameParens = NewParseOperationNameParens()
	f.optOutPort = NewParseOptPort()

	f.connPart.SetSemOutPort(f.semantic.InPort)
	f.semantic.SetOutPort(f.connPart.SemInPort)
	f.connPart.AppendSubOutPort(f.optInPort.InPort)
	f.optInPort.SetOutPort(f.connPart.SubInPort)
	f.connPart.AppendSubOutPort(f.opNameParens.InPort)
	f.opNameParens.SetOutPort(f.connPart.SubInPort)
	f.connPart.AppendSubOutPort(f.optOutPort.InPort)
	f.optOutPort.SetOutPort(f.connPart.SubInPort)

	f.InPort = f.connPart.InPort
	f.SetOutPort = f.connPart.SetOutPort

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
		errPos := md.ParseData.SubResults[0].Pos
		gparselib.AddError(errPos, "At least an operation name or an operation type have to be provided",
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
	opNameParens *gparselib.ParseAll
	semantic     *SemanticOperationNameParens
	optOpName    *gparselib.ParseOptional
	openType     *gparselib.ParseLiteral
	spc1         *ParseOptSpc
	optOpType    *ParseOptOperationType
	closeType    *gparselib.ParseLiteral
	spc2         *ParseOptSpc
	opName       *gparselib.ParseAll
	smallIdent   *ParseSmallIdent
	spc3         *ParseOptSpc
	InPort       func(interface{})
	SetOutPort   func(func(interface{}))
}

func NewParseOperationNameParens() *ParseOperationNameParens {

	f := &ParseOperationNameParens{}
	f.opNameParens = gparselib.NewParseAll(parseData, setParseData)
	f.semantic = NewSemanticOperationNameParens()
	f.optOpName = gparselib.NewParseOptional(parseData, setParseData)
	f.openType = gparselib.NewParseLiteral(parseData, setParseData, "(")
	f.spc1 = NewParseOptSpc()
	f.optOpType = NewParseOptOperationType()
	f.closeType = gparselib.NewParseLiteral(parseData, setParseData, ")")
	f.spc2 = NewParseOptSpc()
	f.opName = gparselib.NewParseAll(parseData, setParseData)
	f.smallIdent = NewParseSmallIdent()
	f.spc3 = NewParseOptSpc()

	f.opNameParens.SetSemOutPort(f.semantic.InPort)
	f.semantic.SetOutPort(f.opNameParens.SemInPort)
	f.opNameParens.AppendSubOutPort(f.optOpName.InPort)
	f.optOpName.SetOutPort(f.opNameParens.SubInPort)
	f.opNameParens.AppendSubOutPort(f.openType.InPort)
	f.openType.SetOutPort(f.opNameParens.SubInPort)
	f.opNameParens.AppendSubOutPort(f.spc1.InPort)
	f.spc1.SetOutPort(f.opNameParens.SubInPort)
	f.opNameParens.AppendSubOutPort(f.optOpType.InPort)
	f.optOpType.SetOutPort(f.opNameParens.SubInPort)
	f.opNameParens.AppendSubOutPort(f.closeType.InPort)
	f.closeType.SetOutPort(f.opNameParens.SubInPort)
	f.opNameParens.AppendSubOutPort(f.spc2.InPort)
	f.spc2.SetOutPort(f.opNameParens.SubInPort)
	f.optOpName.SetSubOutPort(f.opName.InPort)
	f.opName.SetOutPort(f.optOpName.SubInPort)
	f.opName.AppendSubOutPort(f.smallIdent.InPort)
	f.smallIdent.SetOutPort(f.opName.SubInPort)
	f.opName.AppendSubOutPort(f.spc3.InPort)
	f.spc3.SetOutPort(f.opName.SubInPort)

	f.InPort = f.opNameParens.InPort
	f.SetOutPort = f.opNameParens.SetOutPort

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
	optOpType     *gparselib.ParseOptional
	opType        *gparselib.ParseAll
	semantic      *SemanticOptOperationType
	parseBigIdent *ParseBigIdent
	parseOptSpc   *ParseOptSpc
	InPort        func(interface{})
	SetOutPort    func(func(interface{}))
}

func NewParseOptOperationType() *ParseOptOperationType {
	f := &ParseOptOperationType{}
	f.optOpType = gparselib.NewParseOptional(parseData, setParseData)
	f.opType = gparselib.NewParseAll(parseData, setParseData)
	f.semantic = NewSemanticOptOperationType()
	f.parseBigIdent = NewParseBigIdent()
	f.parseOptSpc = NewParseOptSpc()

	f.optOpType.SetSubOutPort(f.opType.InPort)
	f.opType.SetOutPort(f.optOpType.SubInPort)
	f.opType.SetSemOutPort(f.semantic.InPort)
	f.semantic.SetOutPort(f.opType.SemInPort)
	f.opType.AppendSubOutPort(f.parseBigIdent.InPort)
	f.parseBigIdent.SetOutPort(f.opType.SubInPort)
	f.opType.AppendSubOutPort(f.parseOptSpc.InPort)
	f.parseOptSpc.SetOutPort(f.opType.SubInPort)

	f.InPort = f.optOpType.InPort
	f.SetOutPort = f.optOpType.SetOutPort

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
	arrow      *gparselib.ParseAll
	semantic   *SemanticArrow
	spcCom1    *ParseSpaceComment
	optType    *gparselib.ParseOptional
	optCall    *gparselib.ParseOptional
	litArr     *gparselib.ParseLiteral
	spcCom2    *ParseSpaceComment
	typ        *gparselib.ParseAll
	call       *gparselib.ParseRegexp
	openType   *gparselib.ParseLiteral
	spc1       *ParseOptSpc
	typeName   *ParseBigIdent
	spc2       *ParseOptSpc
	closeType  *gparselib.ParseLiteral
	InPort     func(interface{})
	SetOutPort func(func(interface{}))
}

func NewParseArrow() *ParseArrow {
	f := &ParseArrow{}
	f.arrow = gparselib.NewParseAll(parseData, setParseData)
	f.semantic = NewSemanticArrow()
	f.spcCom1 = NewParseSpaceComment()
	f.optType = gparselib.NewParseOptional(parseData, setParseData)
	f.optCall = gparselib.NewParseOptional(parseData, setParseData)
	f.litArr = gparselib.NewParseLiteral(parseData, setParseData, "->")
	f.spcCom2 = NewParseSpaceComment()
	f.typ = gparselib.NewParseAll(parseData, setParseData)
	f.call = gparselib.NewParseRegexp(parseData, setParseData, "[saip]")
	f.openType = gparselib.NewParseLiteral(parseData, setParseData, "[")
	f.spc1 = NewParseOptSpc()
	f.typeName = NewParseBigIdent()
	f.spc2 = NewParseOptSpc()
	f.closeType = gparselib.NewParseLiteral(parseData, setParseData, "]")

	f.arrow.SetSemOutPort(f.semantic.InPort)
	f.semantic.SetOutPort(f.arrow.SemInPort)
	f.arrow.AppendSubOutPort(f.spcCom1.InPort)
	f.spcCom1.SetOutPort(f.arrow.SubInPort)
	f.arrow.AppendSubOutPort(f.optType.InPort)
	f.optType.SetOutPort(f.arrow.SubInPort)
	f.arrow.AppendSubOutPort(f.optCall.InPort)
	f.optCall.SetOutPort(f.arrow.SubInPort)
	f.arrow.AppendSubOutPort(f.litArr.InPort)
	f.litArr.SetOutPort(f.arrow.SubInPort)
	f.arrow.AppendSubOutPort(f.spcCom2.InPort)
	f.spcCom2.SetOutPort(f.arrow.SubInPort)
	f.optType.SetSubOutPort(f.typ.InPort)
	f.typ.SetOutPort(f.optType.SubInPort)
	f.optCall.SetSubOutPort(f.call.InPort)
	f.call.SetOutPort(f.optCall.SubInPort)
	f.typ.AppendSubOutPort(f.openType.InPort)
	f.openType.SetOutPort(f.typ.SubInPort)
	f.typ.AppendSubOutPort(f.spc1.InPort)
	f.spc1.SetOutPort(f.typ.SubInPort)
	f.typ.AppendSubOutPort(f.typeName.InPort)
	f.typeName.SetOutPort(f.typ.SubInPort)
	f.typ.AppendSubOutPort(f.spc2.InPort)
	f.spc2.SetOutPort(f.typ.SubInPort)
	f.typ.AppendSubOutPort(f.closeType.InPort)
	f.closeType.SetOutPort(f.typ.SubInPort)

	f.InPort = f.arrow.InPort
	f.SetOutPort = f.arrow.SetOutPort

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
	optPortSpc *gparselib.ParseOptional
	portSpc    *gparselib.ParseAll
	semantic   *SemanticOptPortSpc
	pport      *ParsePort
	space      *gparselib.ParseSpace
	InPort     func(interface{})
	SetOutPort func(func(interface{}))
}

func NewParseOptPortSpc() *ParseOptPortSpc {
	f := &ParseOptPortSpc{}
	f.optPortSpc = gparselib.NewParseOptional(parseData, setParseData)
	f.portSpc = gparselib.NewParseAll(parseData, setParseData)
	f.semantic = NewSemanticOptPortSpc()
	f.pport = NewParsePort()
	f.space = gparselib.NewParseSpace(parseData, setParseData, false)

	f.optPortSpc.SetSubOutPort(f.portSpc.InPort)
	f.portSpc.SetOutPort(f.optPortSpc.SubInPort)
	f.portSpc.SetSemOutPort(f.semantic.InPort)
	f.semantic.SetOutPort(f.portSpc.SemInPort)
	f.portSpc.AppendSubOutPort(f.pport.InPort)
	f.pport.SetOutPort(f.portSpc.SubInPort)
	f.portSpc.AppendSubOutPort(f.space.InPort)
	f.space.SetOutPort(f.portSpc.SubInPort)

	f.InPort = f.optPortSpc.InPort
	f.SetOutPort = f.optPortSpc.SetOutPort

	return f
}

// ------------ ParseOptPort:
// semantic result: port data.Port
type ParseOptPort struct {
	optPort    *gparselib.ParseOptional
	pport      *ParsePort
	InPort     func(interface{})
	SetOutPort func(func(interface{}))
}

func NewParseOptPort() *ParseOptPort {
	f := &ParseOptPort{}
	f.optPort = gparselib.NewParseOptional(parseData, setParseData)
	f.pport = NewParsePort()

	f.optPort.SetSubOutPort(f.pport.InPort)
	f.pport.SetOutPort(f.optPort.SubInPort)

	f.InPort = f.optPort.InPort
	f.SetOutPort = f.optPort.SetOutPort

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
	port       *gparselib.ParseAll
	semantic   *SemanticPort
	portName   *ParseSmallIdent
	optPortNum *gparselib.ParseOptional
	portNum    *gparselib.ParseAll
	dot        *gparselib.ParseLiteral
	num        *gparselib.ParseNatural
	InPort     func(interface{})
	SetOutPort func(func(interface{}))
}

func NewParsePort() *ParsePort {
	f := &ParsePort{}
	f.port = gparselib.NewParseAll(parseData, setParseData)
	f.semantic = NewSemanticPort()
	f.portName = NewParseSmallIdent()
	f.optPortNum = gparselib.NewParseOptional(parseData, setParseData)
	f.portNum = gparselib.NewParseAll(parseData, setParseData)
	f.dot = gparselib.NewParseLiteral(parseData, setParseData, ".")
	f.num = gparselib.NewParseNatural(parseData, setParseData, 10)

	f.port.SetSemOutPort(f.semantic.InPort)
	f.semantic.SetOutPort(f.port.SemInPort)
	f.port.AppendSubOutPort(f.portName.InPort)
	f.portName.SetOutPort(f.port.SubInPort)
	f.port.AppendSubOutPort(f.optPortNum.InPort)
	f.optPortNum.SetOutPort(f.port.SubInPort)
	f.optPortNum.SetSubOutPort(f.portNum.InPort)
	f.portNum.SetOutPort(f.optPortNum.SubInPort)
	f.portNum.AppendSubOutPort(f.dot.InPort)
	f.dot.SetOutPort(f.portNum.SubInPort)
	f.portNum.AppendSubOutPort(f.num.InPort)
	f.num.SetOutPort(f.portNum.SubInPort)

	f.InPort = f.port.InPort
	f.SetOutPort = f.port.SetOutPort

	return f
}
