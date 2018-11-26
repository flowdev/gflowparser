// Utilities for the Flow Parser
//
// This file contains some utilities that help building the flow parser.
// Most of them are themself simple parsers.

package parser

import (
	"strings"

	"github.com/flowdev/gparselib"
)

// NameIdentParser is a RegexpParser for parsing a name identifier.
type NameIdentParser gparselib.RegexpParser

// NewNameIdentParser creates a new parser for the given regular expression.
// If the regular expression is invalid an error is returned.
func NewNameIdentParser() (*NameIdentParser, error) {
	p, err := gparselib.NewRegexpParser(`^[a-z][a-zA-Z0-9]*`)
	return (*NameIdentParser)(p), err
}

// ParseNameIdent parses a name identifier.
// * Regexp: [a-z][a-zA-Z0-9]*
// * Semantic result: The parsed text.
//
// flow:
//     in (gparselib.ParseData)-> [gparselib.ParseRegexp[semantics=TextSemantic]] -> out
func (p *NameIdentParser) ParseNameIdent(
	pd *gparselib.ParseData, ctx interface{},
) (*gparselib.ParseData, interface{}) {
	return ((*gparselib.RegexpParser)(p)).ParseRegexp(pd, ctx, TextSemantic)
}

// PackageIdentParser is a RegexpParser for parsing a package identifier.
type PackageIdentParser gparselib.RegexpParser

// NewPackageIdentParser creates a new parser for the given regular expression.
// If the regular expression is invalid an error is returned.
func NewPackageIdentParser() (*PackageIdentParser, error) {
	p, err := gparselib.NewRegexpParser(`^[a-z][a-z0-9]*\.`)
	return (*PackageIdentParser)(p), err
}

// ParsePackageIdent parses a package identifier.
// * Regexp: [a-z][a-z0-9]*\.
// * Semantic result: The parsed text (without the dot).
//
// flow:
//     in (gparselib.ParseData)-> [gparselib.ParseRegexp[semantics=TextSemantic]] -> out
func (p *PackageIdentParser) ParsePackageIdent(
	pd *gparselib.ParseData, ctx interface{},
) (*gparselib.ParseData, interface{}) {
	return ((*gparselib.RegexpParser)(p)).ParseRegexp(pd, ctx,
		func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
			pd.Result.Value = pd.Result.Text[:len(pd.Result.Text)-1]
			return pd, ctx
		})
}

// LocalTypeIdentParser parses a local (without package) type identifier.
// Regexp: [A-Za-z][a-zA-Z0-9]*
// Semantic result: The parsed text.
//
// flow:
//     in (ParseData)-> [gparselib.ParseRegexp[semantics=TextSemantic]] -> out
//
// Details:
//  - [ParseData](https://github.com/flowdev/gparselib/blob/master/base.go#L74-L79)
//  - [ParseRegexp](https://github.com/flowdev/gparselib/blob/master/simpleParser.go#L163)
//  - [TextSemantic](./parseUtils.md#textsemantic)
type LocalTypeIdentParser gparselib.RegexpParser

// NewLocalTypeIdentParser creates a new parser for the given regular expression.
// If the regular expression is invalid an error is returned.
func NewLocalTypeIdentParser() (*LocalTypeIdentParser, error) {
	p, err := gparselib.NewRegexpParser(`^[A-Za-z][a-zA-Z0-9]*`)
	return (*LocalTypeIdentParser)(p), err
}

// ParseLocalTypeIdent is the input port of the LocalTypeIdentParser operation.
func (p *LocalTypeIdentParser) ParseLocalTypeIdent(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	return ((*gparselib.RegexpParser)(p)).ParseRegexp(pd, ctx, TextSemantic)
}

// ParseOptSpc parses optional space but no newline.
// Semantic result: The parsed text.
func ParseOptSpc(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	pSpc := func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseSpace(pd2, ctx2, TextSemantic, false)
	}
	return gparselib.ParseOptional(pd, ctx, pSpc, TextSemantic)
}

// ParseASpc parses space but no newline.
// Semantic result: The parsed text.
func ParseASpc(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	return gparselib.ParseSpace(pd, ctx, TextSemantic, false)
}

// SpaceCommentSemValue is the semantic representation of space and comments.
// It specifically informs whether a newline has been parsed.
type SpaceCommentSemValue struct {
	Text    string
	NewLine bool
}

const newLineRune = 10

// spaceCommentSemantic returns the successfully parsed text as semantic value
// plus a signal whether a newline has been parsed.
func spaceCommentSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	semVal := SpaceCommentSemValue{Text: pd.Result.Text}
	semVal.NewLine = strings.ContainsRune(semVal.Text, newLineRune)
	pd.Result.Value = semVal
	return pd, ctx
}

// ParseSpaceComment parses any amount of space (including newline) and line
// (`//` ... <NL>) and block (`/*` ... `*/`) comments.
// Semantic result: The parsed text plus a signal whether a newline was
// parsed.
func ParseSpaceComment(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	pSpc := func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseSpace(pd2, ctx2, TextSemantic, true)
	}
	pLnCmnt := func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
		var err error
		pd2, ctx2, err = gparselib.ParseLineComment(pd2, ctx2, TextSemantic, `//`)
		if err != nil {
			panic(err) // can only be a programming error!
		}
		return pd2, ctx2
	}
	pBlkCmnt := func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
		var err error
		pd2, ctx2, err = gparselib.ParseBlockComment(pd2, ctx2, TextSemantic, `/*`, `*/`)
		if err != nil {
			panic(err) // can only be a programming error!
		}
		return pd2, ctx2
	}
	pAny := func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseAny(
			pd2, ctx2,
			[]gparselib.SubparserOp{pSpc, pLnCmnt, pBlkCmnt},
			TextSemantic,
		)
	}
	return gparselib.ParseMulti0(pd, ctx, pAny, spaceCommentSemantic)
}

// Error messages for semantic errors.
const (
	errMsgNoEnd = "A statement must be ended by a semicolon (';'), a new line or the end of the input"
)

// ParseStatementEnd parses optional space and comments as defined by
// `ParseSpaceComment` followed by a semicolon (`;`) and more optional space
// and comments.
// The semicolon can be omited if the space or comments contain a new line or
// at the end of the input.
// Semantic result: The parsed text.
func ParseStatementEnd(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	pSemicolon := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseLiteral(pd, ctx, TextSemantic, `;`)
	}
	pOptSemi := func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseOptional(pd, ctx, pSemicolon, nil)
	}
	pEOF := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseEOF(pd, ctx,
			func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
				pd2.Result.Value = true
				return pd2, ctx2
			},
		)
	}
	pOptEOF := func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseOptional(pd, ctx, pEOF, nil)
	}

	return gparselib.ParseAll(pd, ctx,
		[]gparselib.SubparserOp{ParseSpaceComment, pOptSemi, ParseSpaceComment, pOptEOF},
		func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
			spcCmnt1 := pd2.SubResults[0].Value.(SpaceCommentSemValue)
			semi := pd2.SubResults[1].Value
			spcCmnt2 := pd2.SubResults[2].Value.(SpaceCommentSemValue)
			eof := pd2.SubResults[3].Value
			if spcCmnt1.NewLine || semi != nil || spcCmnt2.NewLine || eof != nil {
				pd2.Result.Value = pd2.Result.Text
			} else {
				pd2.AddError(pd2.Result.Pos, errMsgNoEnd, nil)
				pd2.Result.Value = nil
			}
			return pd2, ctx2
		},
	)
}

// TextSemantic returns the successfully parsed text as semantic value.
func TextSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	pd.Result.Value = pd.Result.Text
	return pd, ctx
}
