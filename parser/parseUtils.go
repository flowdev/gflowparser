// Utilities for the Flow Parser
//
// This file contains some utilities that help building the flow parser.
// Most of them are themself simple parsers.

package parser

import (
	"github.com/flowdev/gflowparser/data"
	"github.com/flowdev/gparselib"
)

// Utility functions needed because gparselib doesn't know about our MainData:

func getParseData(dat interface{}) *gparselib.ParseData {
	md := dat.(*data.MainData)
	return md.ParseData
}
func setParseData(dat interface{}, subData *gparselib.ParseData) interface{} {
	md := dat.(*data.MainData)
	md.ParseData = subData
	return md
}

// ParseSmallIdent parses an identifier that starts with a lower case character
// (a - z). Potentially followed by more valid identifier characters
// (A - Z, a - z or 0 - 9).  The semantic result is the parsed text.
//
// flow:
//     (MainData)-> p:gparselib.ParseRegexp[semantics: TextSemantic] ->
//
// Details:
//  - [MainData](../data/data.md#maindata)
//  - [ParseRegexp](https://github.com/flowdev/gparselib/blob/master/simpleParser.go#L163)
//  - [TextSemantic](./parseUtils.md#textsemantic)
func ParseSmallIdent(portOut func(interface{})) (portIn func(interface{})) {
	ptIn, err := gparselib.ParseRegexp(
		portOut, TextSemantic,
		getParseData, setParseData,
		`[a-z][a-zA-Z0-9]*`,
	)
	if err != nil {
		panic(err)
	}
	return ptIn
}

// ParseBigIdent parses an identifier that starts with an upper case character
// (A - Z) followed by at least one other valid identifier character
// (A - Z, a - z or 0 - 9).  The semantic result is the parsed text.
//
// flow:
//   MainData->p(gparselib.ParseRegexp)->
//   p MainData=> (TextSemantic) => p
func ParseBigIdent(portOut func(interface{})) (portIn func(interface{})) {
	ptIn, err := gparselib.ParseRegexp(
		portOut, TextSemantic,
		getParseData, setParseData,
		`[A-Z][a-zA-Z0-9]+`,
	)
	if err != nil {
		panic(err)
	}
	return ptIn
}

// ParseOptSpc parses optional space but no newline.
// The semantic result is the parsed text.
func ParseOptSpc(portOut func(interface{})) (portIn func(interface{})) {
	pSpc := func(portOut func(interface{})) (portIn func(interface{})) {
		return gparselib.ParseSpace(
			portOut, nil,
			getParseData, setParseData,
			false,
		)
	}
	portIn = gparselib.ParseOptional(
		portOut, pSpc, TextSemantic,
		getParseData, setParseData,
	)
	return
}

// ParseSpaceComment parses any amount of space (including newline) and line
// (`//` ... <NL>) and block (`/*` ... `*/`) comments.  The semantic result is
// the parsed text.
func ParseSpaceComment(portOut func(interface{})) (portIn func(interface{})) {
	pSpc := func(portOut func(interface{})) (portIn func(interface{})) {
		return gparselib.ParseSpace(
			portOut, nil,
			getParseData, setParseData,
			true,
		)
	}
	pLnCmnt := func(portOut func(interface{})) (portIn func(interface{})) {
		ptIn, err := gparselib.ParseLineComment(
			portOut, nil,
			getParseData, setParseData,
			`//`,
		)
		if err != nil {
			panic(err)
		}
		return ptIn
	}
	pBlkCmnt := func(portOut func(interface{})) (portIn func(interface{})) {
		ptIn, err := gparselib.ParseBlockComment(
			portOut, nil,
			getParseData, setParseData,
			`/*`, `*/`,
		)
		if err != nil {
			panic(err)
		}
		return ptIn
	}
	pAny := func(portOut func(interface{})) (portIn func(interface{})) {
		return gparselib.ParseAny(
			portOut, []gparselib.SubparserOp{pSpc, pLnCmnt, pBlkCmnt}, TextSemantic,
			getParseData, setParseData,
		)
	}
	portIn = gparselib.ParseMulti0(
		portOut, pAny, TextSemantic,
		getParseData, setParseData,
	)
	return
}

// ParseStatementEnd parses optional space and comments as defined by
// `ParseSpaceComment` followed by a semicolon (`;`) and more optional space
// and comments.  The semantic result is the parsed text.
func ParseStatementEnd(portOut func(interface{})) (portIn func(interface{})) {
	pSemicolon := func(portOut func(interface{})) (portIn func(interface{})) {
		return gparselib.ParseLiteral(
			portOut, nil,
			getParseData, setParseData,
			`;`,
		)
	}
	portIn = gparselib.ParseAll(
		portOut,
		[]gparselib.SubparserOp{ParseSpaceComment, pSemicolon, ParseSpaceComment},
		TextSemantic,
		getParseData, setParseData,
	)
	return
}

// TextSemantic returns the successfully parsed text as semantic value.
// Semantics are called by gparselib and thus have to accept the empty interface.
func TextSemantic(portOut func(interface{})) (portIn func(interface{})) {
	portIn = func(dat interface{}) {
		md := dat.(*data.MainData)
		md.ParseData.Result.Value = md.ParseData.Result.Text
		md.ParseData.SubResults = nil
		portOut(md)
	}
	return
}
