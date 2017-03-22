package parser

import (
	"github.com/flowdev/gflowparser/data"
	"github.com/flowdev/gflowparser/semantic"
	"github.com/flowdev/gparselib"
)

// ------------ ParseConnections:
// semantic result: (flow data.Flow{})
type ParseConnections struct {
	//connections *gparselib.ParseMulti1
	//semantic    *SemanticConnections
	//chain       *gparselib.ParseAll
	//chainBeg    *ParseChainBegin
	//chainMids   *gparselib.ParseMulti0
	//optChainEnd *gparselib.ParseOptional
	//stmtEnd     *ParseStatementEnd
	//chainMid    *ParseChainMiddle
	//chainEnd    *ParseChainEnd
	InPort     func(interface{})
	SetOutPort func(func(interface{}))
}

func NewParseConnections() *ParseConnections {
	f := &ParseConnections{}
	connections := gparselib.NewParseMulti1(parseData, setParseData)
	semantic := semantic.NewSemanticConnections()
	chain := gparselib.NewParseAll(parseData, setParseData)
	chainBeg := NewParseChainBegin()
	chainMids := gparselib.NewParseMulti0(parseData, setParseData)
	optChainEnd := gparselib.NewParseOptional(parseData, setParseData)
	stmtEnd := NewParseStatementEnd()
	chainMid := NewParseChainMiddle()
	chainEnd := NewParseChainEnd()

	connections.SetSemOutPort(semantic.InPort)
	semantic.SetOutPort(connections.SemInPort)
	connections.SetSubOutPort(chain.InPort)
	chain.SetOutPort(connections.SubInPort)
	chain.AppendSubOutPort(chainBeg.InPort)
	chainBeg.SetOutPort(chain.SubInPort)
	chain.AppendSubOutPort(chainMids.InPort)
	chainMids.SetOutPort(chain.SubInPort)
	chain.AppendSubOutPort(optChainEnd.InPort)
	optChainEnd.SetOutPort(chain.SubInPort)
	chain.AppendSubOutPort(stmtEnd.InPort)
	stmtEnd.SetOutPort(chain.SubInPort)
	chainMids.SetSubOutPort(chainMid.InPort)
	chainMid.SetOutPort(chainMids.SubInPort)
	optChainEnd.SetSubOutPort(chainEnd.InPort)
	chainEnd.SetOutPort(optChainEnd.SubInPort)

	f.InPort = connections.InPort
	f.SetOutPort = connections.SetOutPort

	return f
}

// ------------ ParseChainBegin:
// semantic result: { (conn data.Connection{FromPort, DataType, ShowDataType, ToPort, ToOp}?), (oper data.Operation{Name, Type, SrcPos, OutPorts}) }
type SemanticChainBeginMin struct {
	outPort func(interface{})
}

func NewSemanticChainBeginMin() *SemanticChainBeginMin {
	return &SemanticChainBeginMin{}
}
func (op *SemanticChainBeginMin) InPort(dat interface{}) {
	var port *data.PortData
	var oper *data.Operation
	md := dat.(*data.MainData)
	res := md.ParseData.Result
	subVals := res.Value.([]interface{})
	subRes := md.ParseData.SubResults

	oper = subVals[0].(*data.Operation)
	if subVals[1] != nil {
		port = subVals[1].(*data.PortData)
	} else {
		port = data.DefaultOutPort(subRes[1].Pos)
	}
	oper.OutPorts = append(oper.OutPorts, port)
	md.ParseData.Result.Value = []interface{}{nil, oper}
	op.outPort(md)
}
func (op *SemanticChainBeginMin) SetOutPort(port func(interface{})) {
	op.outPort = port
}

type SemanticChainBeginMax struct {
	outPort func(interface{})
}

func NewSemanticChainBeginMax() *SemanticChainBeginMax {
	return &SemanticChainBeginMax{}
}
func (op *SemanticChainBeginMax) InPort(dat interface{}) {
	var port *data.PortData
	var oper *data.Operation
	dataType := ""
	md := dat.(*data.MainData)
	res := md.ParseData.Result
	subVals := res.Value.([]interface{})
	subRes := md.ParseData.SubResults
	conn := &data.Connection{}

	if subVals[0] != nil {
		port = subVals[0].(*data.PortData)
	}
	subSubVals := subVals[1].([]interface{})
	dataType = subSubVals[0].(string)
	oper = subSubVals[1].(*data.Operation)
	if port == nil {
		port = data.CopyPort(oper.InPorts[0], subRes[0].Pos)
	}
	conn.FromPort = port
	conn.DataType = dataType
	conn.ShowDataType = (dataType != "")
	conn.ToPort = oper.InPorts[0]
	conn.ToOp = oper
	md.ParseData.Result.Value = []interface{}{conn, oper}
	op.outPort(md)
}
func (op *SemanticChainBeginMax) SetOutPort(port func(interface{})) {
	op.outPort = port
}

type ParseChainBegin struct {
	//chainBeg     *gparselib.ParseAny
	//chainBegMax  *gparselib.ParseAll
	//maxSemantic  *SemanticChainBeginMax
	//optPortMax   *ParseOptPort
	//chainMid     *ParseChainMiddle
	//chainBegMin  *gparselib.ParseAll
	//minSemantic  *SemanticChainBeginMin
	//opNameParens *ParseOperationNameParens
	//optPortMin   *ParseOptPort
	InPort     func(interface{})
	SetOutPort func(func(interface{}))
}

func NewParseChainBegin() *ParseChainBegin {
	f := &ParseChainBegin{}
	chainBeg := gparselib.NewParseAny(parseData, setParseData)
	chainBegMax := gparselib.NewParseAll(parseData, setParseData)
	chainBegMin := gparselib.NewParseAll(parseData, setParseData)
	maxSemantic := NewSemanticChainBeginMax()
	optPortMax := NewParseOptPort()
	chainMid := NewParseChainMiddle()
	minSemantic := NewSemanticChainBeginMin()
	opNameParens := NewParseOperationNameParens()
	optPortMin := NewParseOptPort()

	chainBeg.AppendSubOutPort(chainBegMax.InPort)
	chainBegMax.SetOutPort(chainBeg.SubInPort)
	chainBegMax.SetSemOutPort(maxSemantic.InPort)
	maxSemantic.SetOutPort(chainBegMax.SemInPort)
	chainBegMax.AppendSubOutPort(optPortMax.InPort)
	optPortMax.SetOutPort(chainBegMax.SubInPort)
	chainBegMax.AppendSubOutPort(chainMid.InPort)
	chainMid.SetOutPort(chainBegMax.SubInPort)
	chainBeg.AppendSubOutPort(chainBegMin.InPort)
	chainBegMin.SetOutPort(chainBeg.SubInPort)
	chainBegMin.SetSemOutPort(minSemantic.InPort)
	minSemantic.SetOutPort(chainBegMin.SemInPort)
	chainBegMin.AppendSubOutPort(opNameParens.InPort)
	opNameParens.SetOutPort(chainBegMin.SubInPort)
	chainBegMin.AppendSubOutPort(optPortMin.InPort)
	optPortMin.SetOutPort(chainBegMin.SubInPort)

	f.InPort = chainBeg.InPort
	f.SetOutPort = chainBeg.SetOutPort

	return f
}

// ------------ ParseChainMiddle:
// semantic result: { (bigIdentDataType string), (op data.Operation{Name, Type, SrcPos, InPorts, OutPorts}) }
type ParseChainMiddle struct {
	//chainMid   *gparselib.ParseAll
	//arrow      *ParseArrow
	//connPart   *ParseConnectionPart
	InPort     func(interface{})
	SetOutPort func(func(interface{}))
}

func NewParseChainMiddle() *ParseChainMiddle {
	f := &ParseChainMiddle{}
	chainMid := gparselib.NewParseAll(parseData, setParseData)
	arrow := NewParseArrow()
	connPart := NewParseConnectionPart()

	chainMid.AppendSubOutPort(arrow.InPort)
	arrow.SetOutPort(chainMid.SubInPort)
	chainMid.AppendSubOutPort(connPart.InPort)
	connPart.SetOutPort(chainMid.SubInPort)

	f.InPort = chainMid.InPort
	f.SetOutPort = chainMid.SetOutPort

	return f
}

// ------------ ParseChainEnd:
// semantic result: connection data.Connection{FromPort{}, DataType, ToPort}
type SemanticChainEnd struct {
	outPort func(interface{})
}

func NewSemanticChainEnd() *SemanticChainEnd {
	return &SemanticChainEnd{}
}
func (op *SemanticChainEnd) InPort(dat interface{}) {
	md := dat.(*data.MainData)
	res := md.ParseData.Result
	subVals := res.Value.([]interface{})
	subRes := md.ParseData.SubResults
	conn := &data.Connection{FromPort: data.NewPort("", subRes[0].Pos)}
	res.Value = conn

	if subVals[0] != nil {
		conn.DataType = subVals[0].(string)
	}
	if subVals[1] != nil {
		conn.ToPort = subVals[1].(*data.PortData)
	} else {
		conn.ToPort = data.NewPort("", subRes[1].Pos)
	}
	op.outPort(md)
}
func (op *SemanticChainEnd) SetOutPort(port func(interface{})) {
	op.outPort = port
}

type ParseChainEnd struct {
	//chainEnd   *gparselib.ParseAll
	//semantic   *SemanticChainEnd
	//arrow      *ParseArrow
	//optPort    *ParseOptPort
	InPort     func(interface{})
	SetOutPort func(func(interface{}))
}

func NewParseChainEnd() *ParseChainEnd {
	f := &ParseChainEnd{}
	chainEnd := gparselib.NewParseAll(parseData, setParseData)
	semantic := NewSemanticChainEnd()
	arrow := NewParseArrow()
	optPort := NewParseOptPort()

	chainEnd.SetSemOutPort(semantic.InPort)
	semantic.SetOutPort(chainEnd.SemInPort)
	chainEnd.AppendSubOutPort(arrow.InPort)
	arrow.SetOutPort(chainEnd.SubInPort)
	chainEnd.AppendSubOutPort(optPort.InPort)
	optPort.SetOutPort(chainEnd.SubInPort)

	f.InPort = chainEnd.InPort
	f.SetOutPort = chainEnd.SetOutPort

	return f
}
