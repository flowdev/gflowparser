package parser

import (
	"strings"

	"github.com/flowdev/gflowparser/data"
	"github.com/flowdev/gparselib"
)

// ParseType parses a type declaration including optional package.
// Semantic result: The optional package name and the local type name.
//
// flow:
//     in (ParseData)-> [pOpt gparselib.ParseOptional [ParsePackageIdent]] -> out
//     in (ParseData)-> [gparselib.ParseAll [pOpt, ParseLocalTypeIdent]] -> out
//
// Details:
type ParseType struct {
	pLocalType *LocalTypeIdentParser
	pPack      *PackageIdentParser
}

// NewParseType creates a new parser for a type declaration.
// If any regular expression used by the subparsers is invalid an error is
// returned.
func NewParseType() (*ParseType, error) {
	pPack, err := NewPackageIdentParser()
	if err != nil {
		return nil, err
	}
	pLType, err := NewLocalTypeIdentParser()
	if err != nil {
		return nil, err
	}
	return &ParseType{pPack: pPack, pLocalType: pLType}, nil
}

// In is the input port of the ParseType operation.
func (p *ParseType) In(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	pOpt := func(pd2 *gparselib.ParseData, ctx2 interface{},
	) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseOptional(pd2, ctx2, p.pPack.ParsePackageIdent, nil)
	}
	return gparselib.ParseAll(
		pd, ctx,
		[]gparselib.SubparserOp{pOpt, p.pLocalType.ParseLocalTypeIdent},
		parseTypeSemantic,
	)
}
func parseTypeSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	val0 := pd.SubResults[0].Value
	pack := ""
	if val0 != nil {
		pack = (val0).(string)
	}
	pd.Result.Value = data.Type{
		Package:   pack,
		LocalType: (pd.SubResults[1].Value).(string),
		SrcPos:    pd.Result.Pos,
	}
	return pd, ctx
}

// ParseCompDecl parses an operation declaration.
// Semantic result: The name and the type.
//
// flow:
//     in (ParseData)-> [pAll gparselib.ParseAll [ParseNameIdent, ParseASpc]] -> out
//     in (ParseData)-> [pOpt gparselib.ParseOptional [pAll]] -> out
//     in (ParseData)-> [gparselib.ParseAll [pOpt, ParseType]] -> out
//
// Details:
type ParseCompDecl struct {
	pName *NameIdentParser
	pType *ParseType
}

// NewParseCompDecl creates a new parser for an operation declaration.
// If any regular expression used by the subparsers is invalid an error is
// returned.
func NewParseCompDecl() (*ParseCompDecl, error) {
	pName, err := NewNameIdentParser()
	if err != nil {
		return nil, err
	}
	pType, err := NewParseType()
	if err != nil {
		return nil, err
	}
	return &ParseCompDecl{pName: pName, pType: pType}, nil
}

// In is the input port of the ParseCompDecl operation.
func (p *ParseCompDecl) In(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	pLong := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseAll(
			pd, ctx,
			[]gparselib.SubparserOp{p.pName.ParseNameIdent, ParseASpc, p.pType.In},
			parseCompDeclSemantic,
		)
	}
	return gparselib.ParseAny(
		pd, ctx,
		[]gparselib.SubparserOp{pLong, p.pType.In},
		func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
			if typ, ok := pd2.Result.Value.(data.Type); ok {
				name := nameFromType(typ.LocalType)
				pd2.Result.Value = data.CompDecl{
					Name:      name,
					Type:      typ,
					VagueType: name == typ.LocalType && typ.Package == "",
					SrcPos:    pd.Result.Pos,
				}
			}
			return pd2, ctx2
		},
	)
}
func parseCompDeclSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	name := (pd.SubResults[0].Value).(string)
	typeVal := (pd.SubResults[2].Value).(data.Type)
	pd.Result.Value = data.CompDecl{
		Name:      name,
		Type:      typeVal,
		VagueType: false,
		SrcPos:    pd.Result.Pos,
	}
	return pd, ctx
}
func nameFromType(localType string) string {
	return strings.ToLower(localType[:1]) + localType[1:]
}

// ParseTypeList parses types separated by commas.
// Semantic result: []data.Type
//
// flow:
//     in (ParseData)-> [pAdditionalType gparselib.ParseAll
//                          [ParseSpaceComment, ParseLiteral, ParseSpaceComment, ParseType]
//                      ] -> out
//     in (ParseData)-> [pAdditionalTypes gparselib.ParseMulti0 [pAdditionalType]] -> out
//     in (ParseData)-> [gparselib.ParseAll
//                          [ParseType, pAdditionalTypes]
//                      ] -> out
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
	alltypes := make([](data.Type), len(additionalTypes)+1)
	alltypes[0] = firstType.(data.Type)

	for i, typ := range additionalTypes {
		alltypes[i+1] = typ.(data.Type)
	}
	pd.Result.Value = alltypes
	return pd, ctx
}

// ParseTitledTypes parses a name followed by the equals sign and types separated by commas.
// Semantic result: The title and a slice of *data.Type.
//
// flow:
//     in (ParseData)-> [gparselib.ParseAll
//                          [ ParseNameIdent, ParseSpaceComment, ParseLiteral,
//                            ParseSpaceComment, ParseTypeList                 ]
//                      ] -> out
//
// Details:
type ParseTitledTypes struct {
	pn  *NameIdentParser
	ptl *ParseTypeList
	pt  *ParseType
}

// NewParseTitledTypes creates a new parser for a titled type list.
// If any regular expression used by the subparsers is invalid an error is
// returned.
func NewParseTitledTypes() (*ParseTitledTypes, error) {
	pn, err := NewNameIdentParser()
	if err != nil {
		return nil, err
	}
	ptl, err := NewParseTypeList()
	if err != nil {
		return nil, err
	}
	pt, err := NewParseType()
	if err != nil {
		return nil, err
	}
	return &ParseTitledTypes{pn: pn, ptl: ptl, pt: pt}, nil
}

// In is the input port of the ParseTitledTypes operation.
func (p *ParseTitledTypes) In(pd *gparselib.ParseData, ctx interface{},
) (*gparselib.ParseData, interface{}) {
	pEqual := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseLiteral(pd, ctx, nil, `=`)
	}
	pBig := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseAll(pd, ctx,
			[]gparselib.SubparserOp{p.pn.ParseNameIdent, ParseSpaceComment, pEqual, ParseSpaceComment, p.ptl.In},
			func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
				val0 := pd2.SubResults[0].Value
				val4 := pd2.SubResults[4].Value
				pd2.Result.Value = data.NameNTypes{
					Name:   val0.(string),
					Types:  val4.([]data.Type),
					SrcPos: pd.Result.Pos,
				}
				return pd2, ctx2
			},
		)
	}
	return gparselib.ParseAny(
		pd, ctx,
		[]gparselib.SubparserOp{pBig, p.pt.In},
		func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
			if typ, ok := pd2.Result.Value.(data.Type); ok {
				pd2.Result.Value = data.NameNTypes{
					Types:  []data.Type{typ},
					SrcPos: pd.Result.Pos,
				}
			}
			return pd2, ctx2
		},
	)
}

// ParseTitledTypesList parses TitledTypes separated by a pipe '|' character.
// Semantic result: A slice of data.NameNTypes.
//
// flow:
//     in (ParseData)-> [pAdditionalList gparselib.ParseAll
//                          [ParseSpaceComment, ParseLiteral, ParseSpaceComment, ParseTitledTypes]
//                      ] -> out
//     in (ParseData)-> [pAdditionalLists gparselib.ParseMulti0 [pAdditionalList]] -> out
//     in (ParseData)-> [gparselib.ParseAll
//                          [ParseTitledTypes, pAdditionalLists]
//                      ] -> out
//
// Details:
type ParseTitledTypesList struct {
	ptt *ParseTitledTypes
}

// NewParseTitledTypesList creates a new parser for multiple titled type lists.
// If any regular expression used by the subparsers is invalid an error is
// returned.
func NewParseTitledTypesList() (*ParseTitledTypesList, error) {
	ptt, err := NewParseTitledTypes()
	if err != nil {
		return nil, err
	}
	return &ParseTitledTypesList{ptt: ptt}, nil
}

// In is the input port of the ParseTitledTypesList operation.
func (p *ParseTitledTypesList) In(pd *gparselib.ParseData, ctx interface{},
) (*gparselib.ParseData, interface{}) {
	pBar := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseLiteral(pd, ctx, nil, `|`)
	}
	pAdditionalList := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseAll(
			pd, ctx,
			[]gparselib.SubparserOp{ParseSpaceComment, pBar, ParseSpaceComment, p.ptt.In},
			func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
				pd2.Result.Value = pd2.SubResults[3].Value
				return pd2, ctx2
			},
		)
	}
	pAdditionalLists := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseMulti0(pd, ctx, pAdditionalList, nil)
	}
	return gparselib.ParseAll(
		pd, ctx,
		[]gparselib.SubparserOp{p.ptt.In, pAdditionalLists},
		parseTitledTypesListSemantic,
	)
}
func parseTitledTypesListSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	firstList := pd.SubResults[0].Value
	additionalLists := (pd.SubResults[1].Value).([]interface{})
	alllists := make([](data.NameNTypes), len(additionalLists)+1)
	alllists[0] = firstList.(data.NameNTypes)

	for i, typ := range additionalLists {
		alllists[i+1] = typ.(data.NameNTypes)
	}
	pd.Result.Value = alllists
	return pd, ctx
}

// ParsePlugins parses the plugins of an operation starting with a '[' followed
// by a TitledTypesList or a TypeList and a closing ']'.
// Semantic result: A slice of data.NameNTypes.
//
// flow:
//     in (ParseData)-> [pList gparselib.ParseAny
//                          [ParseTitledTypesList, ParseTypeList]
//                      ] -> out
//     in (ParseData)-> [gparselib.ParseAll
//                          [ ParseLiteral, ParseSpaceComment, pList,
//                            ParseSpaceComment, ParseLiteral         ]
//                      ] -> out
//
// Details:
type ParsePlugins struct {
	pttl *ParseTitledTypesList
	ptl  *ParseTypeList
}

// NewParsePlugins creates a new parser for the plugins of an operation.
// If any regular expression used by the subparsers is invalid an error is
// returned.
func NewParsePlugins() (*ParsePlugins, error) {
	pttl, err := NewParseTitledTypesList()
	if err != nil {
		return nil, err
	}
	ptl, err := NewParseTypeList()
	if err != nil {
		return nil, err
	}
	return &ParsePlugins{pttl: pttl, ptl: ptl}, nil
}

// In is the input port of the ParsePlugins operation.
func (p *ParsePlugins) In(pd *gparselib.ParseData, ctx interface{},
) (*gparselib.ParseData, interface{}) {
	pList := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseBest(
			pd, ctx,
			[]gparselib.SubparserOp{p.pttl.In, p.ptl.In},
			nil,
		)
	}
	pOpen := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseLiteral(pd, ctx, nil, `[`)
	}
	pClose := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseLiteral(pd, ctx, nil, `]`)
	}
	return gparselib.ParseAll(
		pd, ctx,
		[]gparselib.SubparserOp{pOpen, ParseSpaceComment, pList, ParseSpaceComment, pClose},
		parsePluginsSemantic,
	)
}
func parsePluginsSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	list := pd.SubResults[2].Value
	if v, ok := list.([](data.Type)); ok {
		pd.Result.Value = [](data.NameNTypes){
			data.NameNTypes{Name: "", Types: v, SrcPos: v[0].SrcPos},
		}
	} else {
		pd.Result.Value = list
	}

	return pd, ctx
}

// ParseComponent parses a component including declaration and its plugins.
// Semantic result: A data.Component.
//
// flow:
//     in (ParseData)-> [pPlugins gparselib.ParseAll
//                          [ParseSpaceComment, ParsePlugins]
//                      ] -> out
//     in (ParseData)-> [pOpt gparselib.ParseOptional [pPlugins]] -> out
//     in (ParseData)-> [gparselib.ParseAll
//                          [ ParseLiteral, ParseSpaceComment, ParseCompDecl,
//                            pOpt, ParseSpaceComment, ParseLiteral          ]
//                      ] -> out
//
// Details:
type ParseComponent struct {
	pod *ParseCompDecl
	pp  *ParsePlugins
}

// NewParseComponent creates a new parser for a complete component.
// If any regular expression used by the subparsers is invalid an error is
// returned.
func NewParseComponent() (*ParseComponent, error) {
	pod, err := NewParseCompDecl()
	if err != nil {
		return nil, err
	}
	pp, err := NewParsePlugins()
	if err != nil {
		return nil, err
	}
	return &ParseComponent{pod: pod, pp: pp}, nil
}

// In is the input port of the ParseComponent operation.
func (p *ParseComponent) In(pd *gparselib.ParseData, ctx interface{},
) (*gparselib.ParseData, interface{}) {
	pPlugins := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseAll(
			pd, ctx,
			[]gparselib.SubparserOp{ParseSpaceComment, p.pp.In},
			func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
				pd2.Result.Value = pd2.SubResults[1].Value
				return pd2, ctx2
			},
		)
	}
	pOpt := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseOptional(pd, ctx, pPlugins, nil)
	}
	pOpen := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseLiteral(pd, ctx, nil, `[`)
	}
	pClose := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseLiteral(pd, ctx, nil, `]`)
	}
	return gparselib.ParseAll(
		pd, ctx,
		[]gparselib.SubparserOp{pOpen, ParseSpaceComment, p.pod.In, pOpt, ParseSpaceComment, pClose},
		parseComponentSemantic,
	)
}
func parseComponentSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	semVal := data.Component{
		Decl:   (pd.SubResults[2].Value).(data.CompDecl),
		SrcPos: pd.Result.Pos,
	}
	if pd.SubResults[3].Value != nil {
		semVal.Plugins = (pd.SubResults[3].Value).([]data.NameNTypes)
	}
	pd.Result.Value = semVal
	return pd, ctx
}
