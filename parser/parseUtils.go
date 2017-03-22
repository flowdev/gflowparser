package parser

import (
	"github.com/flowdev/gflowparser/data"
	"github.com/flowdev/gparselib"
)

// ------------ Utility functions (needed because the gparselib doesn't know about our MainData of course):
func parseData(dat interface{}) *gparselib.ParseData {
	md := dat.(*data.MainData)
	return md.ParseData
}
func setParseData(dat interface{}, subData *gparselib.ParseData) interface{} {
	md := dat.(*data.MainData)
	md.ParseData = subData
	return md
}

// ------------ TextSemantic (Semantics are called by gparselib and thus have to accept the empty interface):
type TextSemantic struct {
	outPort func(interface{})
}

func NewTextSemantic() *TextSemantic {
	return &TextSemantic{}
}
func (f *TextSemantic) InPort(dat interface{}) {
	md := dat.(*data.MainData)
	md.ParseData.Result.Value = md.ParseData.Result.Text
	f.outPort(md)
}
func (f *TextSemantic) SetOutPort(port func(interface{})) {
	f.outPort = port
}

// ------------ ParseSmallIdent:
// semantic result: text string
type ParseSmallIdent struct {
	parseRegex *gparselib.ParseRegexp
	InPort     func(interface{})
	SetOutPort func(func(interface{}))
}

func NewParseSmallIdent() *ParseSmallIdent {
	f := &ParseSmallIdent{}
	f.parseRegex = gparselib.NewParseRegexp(parseData, setParseData, "[a-z][a-zA-Z0-9]*")

	f.InPort = f.parseRegex.InPort
	f.SetOutPort = f.parseRegex.SetOutPort

	return f
}

// ------------ ParseBigIdent:
// semantic result: text string
type ParseBigIdent struct {
	parseRegex *gparselib.ParseRegexp
	InPort     func(interface{})
	SetOutPort func(func(interface{}))
}

func NewParseBigIdent() *ParseBigIdent {
	f := &ParseBigIdent{}
	f.parseRegex = gparselib.NewParseRegexp(parseData, setParseData, "[A-Z][a-zA-Z0-9]+")

	f.InPort = f.parseRegex.InPort
	f.SetOutPort = f.parseRegex.SetOutPort

	return f
}

// ------------ ParseOptSpc:
// semantic result: text string
type ParseOptSpc struct {
	//optSpc *gparselib.ParseOptional
	//semantic   *TextSemantic
	//parseSpace *gparselib.ParseSpace
	InPort     func(interface{})
	SetOutPort func(func(interface{}))
}

func NewParseOptSpc() *ParseOptSpc {
	f := &ParseOptSpc{}
	optSpc := gparselib.NewParseOptional(parseData, setParseData)
	semantic := NewTextSemantic()
	parseSpace := gparselib.NewParseSpace(parseData, setParseData, false)

	optSpc.SetSemOutPort(semantic.InPort)
	semantic.SetOutPort(optSpc.SemInPort)
	optSpc.SetSubOutPort(parseSpace.InPort)
	parseSpace.SetOutPort(optSpc.SubInPort)

	f.InPort = optSpc.InPort
	f.SetOutPort = optSpc.SetOutPort

	return f
}

// ------------ ParseSpaceComment:
// semantic result: text string
type ParseSpaceComment struct {
	//spcComs *gparselib.ParseMulti0
	//semantic          *TextSemantic
	//spcOrCom          *gparselib.ParseAny
	//parseSpace        *gparselib.ParseSpace
	//parseLineComment  *gparselib.ParseLineComment
	//parseBlockComment *gparselib.ParseBlockComment
	InPort     func(interface{})
	SetOutPort func(func(interface{}))
}

func NewParseSpaceComment() *ParseSpaceComment {
	f := &ParseSpaceComment{}
	spcComs := gparselib.NewParseMulti0(parseData, setParseData)
	semantic := NewTextSemantic()
	spcOrCom := gparselib.NewParseAny(parseData, setParseData)
	parseSpace := gparselib.NewParseSpace(parseData, setParseData, true)
	parseLineComment := gparselib.NewParseLineComment(parseData, setParseData, "//")
	parseBlockComment := gparselib.NewParseBlockComment(parseData, setParseData, "/*", "*/")

	spcComs.SetSemOutPort(semantic.InPort)
	semantic.SetOutPort(spcComs.SemInPort)
	spcComs.SetSubOutPort(spcOrCom.InPort)
	spcOrCom.SetOutPort(spcComs.SubInPort)
	spcOrCom.AppendSubOutPort(parseSpace.InPort)
	parseSpace.SetOutPort(spcOrCom.SubInPort)
	spcOrCom.AppendSubOutPort(parseLineComment.InPort)
	parseLineComment.SetOutPort(spcOrCom.SubInPort)
	spcOrCom.AppendSubOutPort(parseBlockComment.InPort)
	parseBlockComment.SetOutPort(spcOrCom.SubInPort)

	f.InPort = spcComs.InPort
	f.SetOutPort = spcComs.SetOutPort

	return f
}

// ------------ ParseStatementEnd:
// semantic result: text string
type ParseStatementEnd struct {
	//stmtEnd      *gparselib.ParseAll
	//semantic     *TextSemantic
	//optSpc1      *ParseSpaceComment
	//parseLiteral *gparselib.ParseLiteral
	//optSpc2      *ParseSpaceComment
	InPort     func(interface{})
	SetOutPort func(func(interface{}))
}

func NewParseStatementEnd() *ParseStatementEnd {
	f := &ParseStatementEnd{}
	stmtEnd := gparselib.NewParseAll(parseData, setParseData)
	semantic := NewTextSemantic()
	optSpc1 := NewParseSpaceComment()
	parseLiteral := gparselib.NewParseLiteral(parseData, setParseData, ";")
	optSpc2 := NewParseSpaceComment()

	stmtEnd.SetSemOutPort(semantic.InPort)
	semantic.SetOutPort(stmtEnd.SemInPort)
	stmtEnd.AppendSubOutPort(optSpc1.InPort)
	optSpc1.SetOutPort(stmtEnd.SubInPort)
	stmtEnd.AppendSubOutPort(parseLiteral.InPort)
	parseLiteral.SetOutPort(stmtEnd.SubInPort)
	stmtEnd.AppendSubOutPort(optSpc2.InPort)
	optSpc2.SetOutPort(stmtEnd.SubInPort)

	f.InPort = stmtEnd.InPort
	f.SetOutPort = stmtEnd.SetOutPort

	return f
}
