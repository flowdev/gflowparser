package parser

import (
	"github.com/flowdev/gflowparser/data"
	"github.com/flowdev/gparselib"
)

// ------------ Utility functions:
func parseData(dat interface{}) *gparselib.ParseData {
	md := dat.(*data.MainData)
	return md.ParseData
}
func setParseData(dat interface{}, subData *gparselib.ParseData) interface{} {
	md := dat.(*data.MainData)
	md.ParseData = subData
	return md
}

// ------------ TextSemantic:
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
type ParseOptSpc struct {
	optSpc     *gparselib.ParseOptional
	semantic   *TextSemantic
	parseSpace *gparselib.ParseSpace
	InPort     func(interface{})
	SetOutPort func(func(interface{}))
}

func NewParseOptSpc() *ParseOptSpc {
	f := &ParseOptSpc{}
	f.optSpc = gparselib.NewParseOptional(parseData, setParseData)
	f.semantic = NewTextSemantic()
	f.parseSpace = gparselib.NewParseSpace(parseData, setParseData, false)

	f.optSpc.SetSemOutPort(f.semantic.InPort)
	f.semantic.SetOutPort(f.optSpc.SemInPort)
	f.optSpc.SetSubOutPort(f.parseSpace.InPort)
	f.parseSpace.SetOutPort(f.optSpc.SubInPort)

	f.InPort = f.optSpc.InPort
	f.SetOutPort = f.optSpc.SetOutPort

	return f
}

// ------------ ParseSpaceComment:
type ParseSpaceComment struct {
	spcComs           *gparselib.ParseMulti0
	semantic          *TextSemantic
	spcOrCom          *gparselib.ParseAny
	parseSpace        *gparselib.ParseSpace
	parseLineComment  *gparselib.ParseLineComment
	parseBlockComment *gparselib.ParseBlockComment
	InPort            func(interface{})
	SetOutPort        func(func(interface{}))
}

func NewParseSpaceComment() *ParseSpaceComment {
	f := &ParseSpaceComment{}
	f.spcComs = gparselib.NewParseMulti0(parseData, setParseData)
	f.semantic = NewTextSemantic()
	f.spcOrCom = gparselib.NewParseAny(parseData, setParseData)
	f.parseSpace = gparselib.NewParseSpace(parseData, setParseData, true)
	f.parseLineComment = gparselib.NewParseLineComment(parseData, setParseData, "//")
	f.parseBlockComment = gparselib.NewParseBlockComment(parseData, setParseData, "/*", "*/")

	f.spcComs.SetSemOutPort(f.semantic.InPort)
	f.semantic.SetOutPort(f.spcComs.SemInPort)
	f.spcComs.SetSubOutPort(f.spcOrCom.InPort)
	f.spcOrCom.SetOutPort(f.spcComs.SubInPort)
	f.spcOrCom.AppendSubOutPort(f.parseSpace.InPort)
	f.parseSpace.SetOutPort(f.spcOrCom.SubInPort)
	f.spcOrCom.AppendSubOutPort(f.parseLineComment.InPort)
	f.parseLineComment.SetOutPort(f.spcOrCom.SubInPort)
	f.spcOrCom.AppendSubOutPort(f.parseBlockComment.InPort)
	f.parseBlockComment.SetOutPort(f.spcOrCom.SubInPort)

	f.InPort = f.spcComs.InPort
	f.SetOutPort = f.spcComs.SetOutPort

	return f
}

// ------------ ParseStatementEnd:
type ParseStatementEnd struct {
	stmtEnd      *gparselib.ParseAll
	semantic     *TextSemantic
	optSpc1      *ParseSpaceComment
	parseLiteral *gparselib.ParseLiteral
	optSpc2      *ParseSpaceComment
	InPort       func(interface{})
	SetOutPort   func(func(interface{}))
}

func NewParseStatementEnd() *ParseStatementEnd {
	f := &ParseStatementEnd{}
	f.stmtEnd = gparselib.NewParseAll(parseData, setParseData)
	f.semantic = NewTextSemantic()
	f.optSpc1 = NewParseSpaceComment()
	f.parseLiteral = gparselib.NewParseLiteral(parseData, setParseData, ";")
	f.optSpc2 = NewParseSpaceComment()

	f.stmtEnd.SetSemOutPort(f.semantic.InPort)
	f.semantic.SetOutPort(f.stmtEnd.SemInPort)
	f.stmtEnd.AppendSubOutPort(f.optSpc1.InPort)
	f.optSpc1.SetOutPort(f.stmtEnd.SubInPort)
	f.stmtEnd.AppendSubOutPort(f.parseLiteral.InPort)
	f.parseLiteral.SetOutPort(f.stmtEnd.SubInPort)
	f.stmtEnd.AppendSubOutPort(f.optSpc2.InPort)
	f.optSpc2.SetOutPort(f.stmtEnd.SubInPort)

	f.InPort = f.stmtEnd.InPort
	f.SetOutPort = f.stmtEnd.SetOutPort

	return f
}
