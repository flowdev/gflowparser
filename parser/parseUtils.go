// Utilities for the Flow Parser
//
// This file contains some utilities that help building the flow parser.
// Most of them are themself simple parsers.

package parser

import (
	"github.com/flowdev/gparselib"
)

// ParseNameIdent parses an identifier that starts with a lower case character
// (a - z). Potentially followed by more valid identifier characters
// (A - Z, a - z or 0 - 9).  The semantic result is the parsed text.
//
// flow:
//     in (ParseData)-> [gparselib.ParseRegexp[semantics=TextSemantic]] -> out
//
// Details:
//  - [ParseData](https://github.com/flowdev/gparselib/blob/master/base.go#L74-L79)
//  - [ParseRegexp](https://github.com/flowdev/gparselib/blob/master/simpleParser.go#L163)
//  - [TextSemantic](./parseUtils.md#textsemantic)
type ParseNameIdent gparselib.ParseRegexp

// NewParseNameIdent creates a new parser for the given regular expression.
// If the regular expression is invalid an error is returned.
func NewParseNameIdent() (*ParseNameIdent, error) {
	p, err := gparselib.NewParseRegexp(`^[a-z][a-zA-Z0-9]*`)
	return (*ParseNameIdent)(p), err
}

// In is the input port of the ParseNameIdent operation.
func (p *ParseNameIdent) In(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	return ((*gparselib.ParseRegexp)(p)).In(pd, ctx, TextSemantic)
}

// ParsePackageIdent parses an identifier that starts with a lower case character
// (a - z). Potentially followed by more valid lower case identifier characters
// (a - z or 0 - 9).  The semantic result is the parsed text.
//
// flow:
//     in (ParseData)-> [gparselib.ParseRegexp[semantics=TextSemantic]] -> out
//
// Details:
//  - [ParseData](https://github.com/flowdev/gparselib/blob/master/base.go#L74-L79)
//  - [ParseRegexp](https://github.com/flowdev/gparselib/blob/master/simpleParser.go#L163)
//  - [TextSemantic](./parseUtils.md#textsemantic)
type ParsePackageIdent gparselib.ParseRegexp

// NewParsePackageIdent creates a new parser for the given regular expression.
// If the regular expression is invalid an error is returned.
func NewParsePackageIdent() (*ParsePackageIdent, error) {
	p, err := gparselib.NewParseRegexp(`^[a-z][a-z0-9]*`)
	return (*ParsePackageIdent)(p), err
}

// In is the input port of the ParsePackageIdent operation.
func (p *ParsePackageIdent) In(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	return ((*gparselib.ParseRegexp)(p)).In(pd, ctx, TextSemantic)
}

// ParseLocalTypeIdent parses an identifier that starts with an upper case character
// (A - Z). Potentially followed by more valid identifier characters
// (A - Z, a - z or 0 - 9).  The semantic result is the parsed text.
//
// flow:
//     in (ParseData)-> [gparselib.ParseRegexp[semantics=TextSemantic]] -> out
//
// Details:
//  - [ParseData](https://github.com/flowdev/gparselib/blob/master/base.go#L74-L79)
//  - [ParseRegexp](https://github.com/flowdev/gparselib/blob/master/simpleParser.go#L163)
//  - [TextSemantic](./parseUtils.md#textsemantic)
type ParseLocalTypeIdent gparselib.ParseRegexp

// NewParseLocalTypeIdent creates a new parser for the given regular expression.
// If the regular expression is invalid an error is returned.
func NewParseLocalTypeIdent() (*ParseLocalTypeIdent, error) {
	p, err := gparselib.NewParseRegexp(`^[A-Z][a-zA-Z0-9]*`)
	return (*ParseLocalTypeIdent)(p), err
}

// In is the input port of the ParseLocalTypeIdent operation.
func (p *ParseLocalTypeIdent) In(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	return ((*gparselib.ParseRegexp)(p)).In(pd, ctx, TextSemantic)
}

// ParseType parses a type declaration including optional package.
// The semantic result is the optional package name and the local type name.
//
// flow:
//     in (ParseData)-> [gparselib.ParseAll []] -> out
//
// Details:
type ParseType struct {
	pLocalType *ParseLocalTypeIdent
	pPack      *ParsePackageIdent
}

// TypeSemValue is the semantic representation of a type declaration.
type TypeSemValue struct {
	Package   string
	LocalType string
}

// NewParseType creates a new parser for e type declaration.
// If any regular expression used by the subparsers is invalid an error is
// returned.
func NewParseType() (*ParseType, error) {
	pPack, err := NewParsePackageIdent()
	if err != nil {
		return nil, err
	}
	pLType, err := NewParseLocalTypeIdent()
	if err != nil {
		return nil, err
	}
	return &ParseType{pPack: pPack, pLocalType: pLType}, nil
}

// In is the input port of the ParseType operation.
func (p *ParseType) In(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	pDot := func(pd2 *gparselib.ParseData, ctx2 interface{},
	) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseLiteral(pd2, ctx2, TextSemantic, `.`)
	}
	pAll1 := func(pd2 *gparselib.ParseData, ctx2 interface{},
	) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseAll(
			pd2, ctx2,
			[]gparselib.SubparserOp{p.pPack.In, pDot},
			TextSemantic,
		)
	}
	pOpt := func(pd2 *gparselib.ParseData, ctx2 interface{},
	) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseOptional(pd2, ctx2, pAll1, nil)
	}
	return gparselib.ParseAll(
		pd, ctx,
		[]gparselib.SubparserOp{pOpt, p.pLocalType.In},
		parseTypeSemantic,
	)
}
func parseTypeSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	pack := ""
	if pd.SubResults[0].Value != nil {
		optVals := (pd.SubResults[0].Value).([]interface{})
		pack = (optVals[0]).(string)
	}
	pd.Result.Value = &TypeSemValue{
		Package:   pack,
		LocalType: pd.SubResults[1].Text,
	}
	pd.SubResults = nil
	return pd, ctx
}

// ParseOptSpc parses optional space but no newline.
// The semantic result is the parsed text.
func ParseOptSpc(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	pSpc := func(pd2 *gparselib.ParseData, ctx2 interface{},
	) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseSpace(pd2, ctx2, TextSemantic, false)
	}
	return gparselib.ParseOptional(pd, ctx, pSpc, TextSemantic)
}

// ParseASpc parses space but no newline.
// The semantic result is the parsed text.
func ParseASpc(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	return gparselib.ParseSpace(pd, ctx, TextSemantic, false)
}

// ParseSpaceComment parses any amount of space (including newline) and line
// (`//` ... <NL>) and block (`/*` ... `*/`) comments.  The semantic result is
// the parsed text.
func ParseSpaceComment(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	pSpc := func(pd2 *gparselib.ParseData, ctx2 interface{},
	) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseSpace(pd2, ctx2, TextSemantic, true)
	}
	pLnCmnt := func(pd2 *gparselib.ParseData, ctx2 interface{},
	) (*gparselib.ParseData, interface{}) {
		var err error
		pd2, ctx2, err = gparselib.ParseLineComment(pd2, ctx2, TextSemantic, `//`)
		if err != nil {
			panic(err) // can only be a programming error!
		}
		return pd2, ctx2
	}
	pBlkCmnt := func(pd2 *gparselib.ParseData, ctx2 interface{},
	) (*gparselib.ParseData, interface{}) {
		var err error
		pd2, ctx2, err = gparselib.ParseBlockComment(pd2, ctx2, TextSemantic, `/*`, `*/`)
		if err != nil {
			panic(err) // can only be a programming error!
		}
		return pd2, ctx2
	}
	pAny := func(pd2 *gparselib.ParseData, ctx2 interface{},
	) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseAny(
			pd2, ctx2,
			[]gparselib.SubparserOp{pSpc, pLnCmnt, pBlkCmnt},
			TextSemantic,
		)
	}
	return gparselib.ParseMulti0(pd, ctx, pAny, TextSemantic)
}

// ParseStatementEnd parses optional space and comments as defined by
// `ParseSpaceComment` followed by a semicolon (`;`) and more optional space
// and comments.  The semantic result is the parsed text.
/*
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
*/

// TextSemantic returns the successfully parsed text as semantic value.
func TextSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	pd.Result.Value = pd.Result.Text
	pd.SubResults = nil
	return pd, ctx
}
