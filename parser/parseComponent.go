package parser

import (
	"strings"

	"github.com/flowdev/gparselib"
)

// ParseType parses a type declaration including optional package.
// The semantic result is the optional package name and the local type name.
//
// flow:
//     in (ParseData)-> [pOpt gparselib.ParseOptional [subparser = ParsePackageIdent]] -> out
//     in (ParseData)-> [gparselib.ParseAll [subparser = pOpt, ParseLocalTypeIdent]] -> out
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

// NewParseType creates a new parser for a type declaration.
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
	pOpt := func(pd2 *gparselib.ParseData, ctx2 interface{},
	) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseOptional(pd2, ctx2, p.pPack.In, nil)
	}
	return gparselib.ParseAll(
		pd, ctx,
		[]gparselib.SubparserOp{pOpt, p.pLocalType.In},
		parseTypeSemantic,
	)
}
func parseTypeSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	val0 := pd.SubResults[0].Value
	pack := ""
	if val0 != nil {
		pack = (val0).(string)
	}
	pd.Result.Value = &TypeSemValue{
		Package:   pack,
		LocalType: (pd.SubResults[1].Value).(string),
	}
	return pd, ctx
}

// ParseOpDecl parses an operation declaration.
// The semantic result is the name and the type.
//
// flow:
//     in (ParseData)-> [pAll gparselib.ParseAll [ParseNameIdent, ParseASpc]] -> out
//     in (ParseData)-> [pOpt gparselib.ParseOptional [pAll]] -> out
//     in (ParseData)-> [gparselib.ParseAll [pOpt, ParseType]] -> out
//
// Details:
type ParseOpDecl struct {
	pName *ParseNameIdent
	pType *ParseType
}

// OpDeclSemValue is the semantic representation of a type declaration.
type OpDeclSemValue struct {
	Name string
	Type *TypeSemValue
}

// NewParseOpDecl creates a new parser for an operation declaration.
// If any regular expression used by the subparsers is invalid an error is
// returned.
func NewParseOpDecl() (*ParseOpDecl, error) {
	pName, err := NewParseNameIdent()
	if err != nil {
		return nil, err
	}
	pType, err := NewParseType()
	if err != nil {
		return nil, err
	}
	return &ParseOpDecl{pName: pName, pType: pType}, nil
}

// In is the input port of the ParseOpDecl operation.
func (p *ParseOpDecl) In(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	pAll := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseAll(
			pd, ctx,
			[]gparselib.SubparserOp{p.pName.In, ParseASpc},
			func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
				pd2.Result.Value = pd2.SubResults[0].Value
				return pd2, ctx2
			},
		)
	}
	pOpt := func(pd2 *gparselib.ParseData, ctx2 interface{},
	) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseOptional(pd2, ctx2, pAll, nil)
	}
	return gparselib.ParseAll(
		pd, ctx,
		[]gparselib.SubparserOp{pOpt, p.pType.In},
		parseOpDeclSemantic,
	)
}
func parseOpDeclSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	val0 := pd.SubResults[0].Value
	typeVal := (pd.SubResults[1].Value).(*TypeSemValue)
	name := ""
	if val0 != nil {
		name = (val0).(string)
	} else {
		name = strings.ToLower(typeVal.LocalType[:1]) + typeVal.LocalType[1:]
	}
	pd.Result.Value = &OpDeclSemValue{
		Name: name,
		Type: typeVal,
	}
	return pd, ctx
}

// ParseTypeList parses types separated by commas.
// The semantic result is a slice of *TypeSemValue.
//
// flow:
//     in (ParseData)-> [gparselib.ParseAll []] -> out
//
// Details:
type ParseTypeList struct {
	pt *ParseType
}

// NewParseTypeList creates a new parser for a type list.
// If any regular expression used by the subparsers is invalid an error is
// returned.
func NewParseTypeList() (*ParseTypeList, error) {
	p, err := NewParseType()
	if err != nil {
		return nil, err
	}
	return &ParseTypeList{pt: p}, nil
}

// In is the input port of the ParseTypeList operation.
func (p *ParseTypeList) In(pd *gparselib.ParseData, ctx interface{},
) (*gparselib.ParseData, interface{}) {
	pComma := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseLiteral(pd, ctx, nil, `,`)
	}
	pAdditionalType := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseAll(
			pd, ctx,
			[]gparselib.SubparserOp{ParseSpaceComment, pComma, ParseSpaceComment, p.pt.In},
			func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
				pd2.Result.Value = pd2.SubResults[3].Value
				return pd2, ctx2
			},
		)
	}
	pAdditionalTypes := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseMulti0(pd, ctx, pAdditionalType, nil)
	}
	return gparselib.ParseAll(
		pd, ctx,
		[]gparselib.SubparserOp{p.pt.In, pAdditionalTypes},
		parseTypeListSemantic,
	)
}
func parseTypeListSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	firstType := pd.SubResults[0].Value
	additionalTypes := (pd.SubResults[1].Value).([]interface{})
	alltypes := make([](*TypeSemValue), len(additionalTypes)+1)
	alltypes[0] = firstType.(*TypeSemValue)

	for i, typ := range additionalTypes {
		alltypes[i+1] = typ.(*TypeSemValue)
	}
	pd.Result.Value = alltypes
	return pd, ctx
}
