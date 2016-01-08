package parser

import (
	"github.com/flowdev/gflowparser/data"
	"github.com/flowdev/gparselib"
)

// ------------ ParseConnections:
type ParseConnections struct {
	connections *gparselib.ParseMulti1
	//	semantic    *SemanticCreateChain
	chain       *gparselib.ParseAll
	chainBeg    *ParseChainBegin
	chainMids   *gparselib.ParseMulti0
	optChainEnd *gparselib.ParseOptional
	stmtEnd     *ParseStatementEnd
	chainMid    *ParseChainMiddle
	chainEnd    *ParseChainEnd
	InPort      func(interface{})
}

func NewParseConnections() *ParseConnections {
	f := &ParseConnections{}
	f.connections = gparselib.NewParseMulti1(parseData, setParseData)
	//	f.semantic = NewSemanticCreateChain()
	f.chain = gparselib.NewParseAll(parseData, setParseData)
	f.chainBeg = NewParseChainBegin()
	f.chainMids = gparselib.NewParseMulti0(parseData, setParseData)
	f.optChainEnd = gparselib.NewParseOptional(parseData, setParseData)
	f.stmtEnd = NewParseStatementEnd()
	f.chainMid = NewParseChainMiddle()
	f.chainEnd = NewParseChainEnd()

	//	f.connections.SetSemOutPort(f.semantic.InPort)
	//	f.semantic.SetOutPort(f.connections.SemInPort)
	f.connections.SetSubOutPort(f.chain.InPort)
	f.chain.SetOutPort(f.connections.SubInPort)
	f.chain.AppendSubOutPort(f.chainBeg.InPort)
	f.chainBeg.SetOutPort(f.chain.SubInPort)
	f.chain.AppendSubOutPort(f.chainMids.InPort)
	f.chainMids.SetOutPort(f.chain.SubInPort)
	f.chain.AppendSubOutPort(f.optChainEnd.InPort)
	f.optChainEnd.SetOutPort(f.chain.SubInPort)
	f.chain.AppendSubOutPort(f.stmtEnd.InPort)
	f.stmtEnd.SetOutPort(f.chain.SubInPort)
	f.chainMids.SetSubOutPort(f.chainMid.InPort)
	f.chainMid.SetOutPort(f.chainMids.SubInPort)
	f.optChainEnd.SetSubOutPort(f.chainEnd.InPort)
	f.chainEnd.SetOutPort(f.optChainEnd.SubInPort)

	f.InPort = f.connections.InPort

	return f
}
func (f *ParseConnections) SetOutPort(port func(interface{})) {
	f.connections.SetOutPort(port)
}

// ------------ ParseChainBegin:
type ParseChainBegin struct {
	chainBeg    *gparselib.ParseAny
	chainBegMax *gparselib.ParseAll
	chainBegMin *gparselib.ParseAll
	//	maxSemantic  *SemanticCreateChainBegMax
	optPortMax *ParseOptPort
	chainMid   *ParseChainMiddle
	//	minSemantic  *SemanticCreateChainBegMin
	opNameParens *ParseOperationNameParens
	optPortMin   *ParseOptPort
	InPort       func(interface{})
}

func NewParseChainBegin() *ParseChainBegin {
	f := &ParseChainBegin{}
	f.chainBeg = gparselib.NewParseAny(parseData, setParseData)
	f.chainBegMax = gparselib.NewParseAll(parseData, setParseData)
	f.chainBegMin = gparselib.NewParseAll(parseData, setParseData)
	//	f.maxSemantic = NewSemanticCreateChainBegMax()
	f.optPortMax = NewParseOptPort()
	f.chainMid = NewParseChainMiddle()
	//	f.minSemantic = NewSemanticCreateChainBegMin()
	f.opNameParens = NewParseOperationNameParens()
	f.optPortMin = NewParseOptPort()

	f.chainBeg.AppendSubOutPort(f.chainBegMax.InPort)
	f.chainBegMax.SetOutPort(f.chainBeg.SubInPort)
	f.chainBeg.AppendSubOutPort(f.chainBegMin.InPort)
	f.chainBegMin.SetOutPort(f.chainBeg.SubInPort)
	//	f.chainBegMax.SetSemOutPort(f.maxSemantic.InPort)
	//	f.maxSemantic.SetOutPort(f.chainBegMax.SemInPort)
	f.chainBegMax.AppendSubOutPort(f.optPortMax.InPort)
	f.optPortMax.SetOutPort(f.chainBegMax.SubInPort)
	f.chainBegMax.AppendSubOutPort(f.chainMid.InPort)
	f.chainMid.SetOutPort(f.chainBegMax.SubInPort)
	//	f.chainBegMin.SetSemOutPort(f.minSemantic.InPort)
	//	f.minSemantic.SetOutPort(f.chainBegMin.SemInPort)
	f.chainBegMin.AppendSubOutPort(f.opNameParens.InPort)
	f.opNameParens.SetOutPort(f.chainBegMin.SubInPort)
	f.chainBegMin.AppendSubOutPort(f.optPortMin.InPort)
	f.optPortMin.SetOutPort(f.chainBegMin.SubInPort)

	f.InPort = f.chainBeg.InPort

	return f
}
func (f *ParseChainBegin) SetOutPort(port func(interface{})) {
	f.chainBeg.SetOutPort(port)
}

// ------------ ParseChainMiddle:
type ParseChainMiddle struct {
	chainMid   *gparselib.ParseAll
	arrow      *ParseArrow
	connPart   *ParseConnectionPart
	InPort     func(interface{})
	SetOutPort func(func(interface{}))
}

func NewParseChainMiddle() *ParseChainMiddle {
	f := &ParseChainMiddle{}
	f.chainMid = gparselib.NewParseAll(parseData, setParseData)
	f.arrow = NewParseArrow()
	f.connPart = NewParseConnectionPart()

	f.chainMid.AppendSubOutPort(f.arrow.InPort)
	f.arrow.SetOutPort(f.chainMid.SubInPort)
	f.chainMid.AppendSubOutPort(f.connPart.InPort)
	f.connPart.SetOutPort(f.chainMid.SubInPort)

	f.InPort = f.chainMid.InPort
	f.SetOutPort = f.chainMid.SetOutPort

	return f
}

// ------------ ParseChainEnd:
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
	chainEnd *gparselib.ParseAll
	semantic *SemanticChainEnd
	arrow    *ParseArrow
	optPort  *ParseOptPort
	InPort   func(interface{})
}

func NewParseChainEnd() *ParseChainEnd {
	f := &ParseChainEnd{}
	f.chainEnd = gparselib.NewParseAll(parseData, setParseData)
	f.semantic = NewSemanticChainEnd()
	f.arrow = NewParseArrow()
	f.optPort = NewParseOptPort()

	f.chainEnd.SetSemOutPort(f.semantic.InPort)
	f.semantic.SetOutPort(f.chainEnd.SemInPort)
	f.chainEnd.AppendSubOutPort(f.arrow.InPort)
	f.arrow.SetOutPort(f.chainEnd.SubInPort)
	f.chainEnd.AppendSubOutPort(f.optPort.InPort)
	f.optPort.SetOutPort(f.chainEnd.SubInPort)

	f.InPort = f.chainEnd.InPort

	return f
}
func (f *ParseChainEnd) SetOutPort(port func(interface{})) {
	f.chainEnd.SetOutPort(port)
}
