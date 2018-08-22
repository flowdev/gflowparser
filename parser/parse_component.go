package parser

import (
	"strings"

	"github.com/flowdev/gflowparser/data"
	"github.com/flowdev/gparselib"
)

// TypeParser parses a type declaration including optional package.
// Semantic result: The optional package name and the local type name.
//
// flow:
//     in (ParseData)-> [pOpt gparselib.ParseOptional [ParsePackageIdent]] -> out
//     in (ParseData)-> [gparselib.ParseAll [pOpt, ParseLocalTypeIdent]] -> out
//
// Details:
type TypeParser struct {
	pLocalType *LocalTypeIdentParser
	pPack      *PackageIdentParser
}

// NewTypeParser creates a new parser for a type declaration.
// If any regular expression used by the subparsers is invalid an error is
// returned.
func NewTypeParser() (*TypeParser, error) {
	pPack, err := NewPackageIdentParser()
	if err != nil {
		return nil, err
	}
	pLType, err := NewLocalTypeIdentParser()
	if err != nil {
		return nil, err
	}
	return &TypeParser{pPack: pPack, pLocalType: pLType}, nil
}

// ParseType is the input port of the TypeParser operation.
func (p *TypeParser) ParseType(pd *gparselib.ParseData, ctx interface{},
) (*gparselib.ParseData, interface{}) {
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
func parseTypeSemantic(pd *gparselib.ParseData, ctx interface{},
) (*gparselib.ParseData, interface{}) {
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

// CompDeclParser parses a component declaration.
// Semantic result: The name and the type.
//
// flow:
//     in (ParseData)-> [pAll gparselib.ParseAll [ParseNameIdent, ParseASpc]] -> out
//     in (ParseData)-> [pOpt gparselib.ParseOptional [pAll]] -> out
//     in (ParseData)-> [gparselib.ParseAll [pOpt, ParseType]] -> out
//
// Details:
type CompDeclParser struct {
	pName *NameIdentParser
	pType *TypeParser
}

// NewCompDeclParser creates a new parser for an operation declaration.
// If any regular expression used by the subparsers is invalid an error is
// returned.
func NewCompDeclParser() (*CompDeclParser, error) {
	pName, err := NewNameIdentParser()
	if err != nil {
		return nil, err
	}
	pType, err := NewTypeParser()
	if err != nil {
		return nil, err
	}
	return &CompDeclParser{pName: pName, pType: pType}, nil
}

// ParseCompDecl is the input port of the CompDeclParser operation.
func (p *CompDeclParser) ParseCompDecl(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	pLong := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseAll(
			pd, ctx,
			[]gparselib.SubparserOp{p.pName.ParseNameIdent, ParseASpc, p.pType.ParseType},
			parseCompDeclSemantic,
		)
	}
	return gparselib.ParseAny(
		pd, ctx,
		[]gparselib.SubparserOp{pLong, p.pType.ParseType},
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

// TypeListParser parses types separated by commas.
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
type TypeListParser struct {
	pt *TypeParser
}

// NewTypeListParser creates a new parser for a type list.
// If any regular expression used by the subparsers is invalid an error is
// returned.
func NewTypeListParser() (*TypeListParser, error) {
	p, err := NewTypeParser()
	if err != nil {
		return nil, err
	}
	return &TypeListParser{pt: p}, nil
}

// ParseTypeList is the input port of the TypeListParser operation.
func (p *TypeListParser) ParseTypeList(pd *gparselib.ParseData, ctx interface{},
) (*gparselib.ParseData, interface{}) {
	pComma := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseLiteral(pd, ctx, nil, `,`)
	}
	pAdditionalType := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseAll(
			pd, ctx,
			[]gparselib.SubparserOp{ParseSpaceComment, pComma, ParseSpaceComment, p.pt.ParseType},
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
		[]gparselib.SubparserOp{p.pt.ParseType, pAdditionalTypes},
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

// PluginParser parses a name followed by the equals sign and types separated by commas.
// Semantic result: The title and a slice of *data.Type.
//
// flow:
//     in (ParseData)-> [gparselib.ParseAll
//                          [ ParseNameIdent, ParseSpaceComment, ParseLiteral,
//                            ParseSpaceComment, ParseTypeList                 ]
//                      ] -> out
//
// Details:
type PluginParser struct {
	pn  *NameIdentParser
	ptl *TypeListParser
	pt  *TypeParser
}

// NewPluginParser creates a new parser for a titled type list.
// If any regular expression used by the subparsers is invalid an error is
// returned.
func NewPluginParser() (*PluginParser, error) {
	pn, err := NewNameIdentParser()
	if err != nil {
		return nil, err
	}
	ptl, err := NewTypeListParser()
	if err != nil {
		return nil, err
	}
	pt, err := NewTypeParser()
	if err != nil {
		return nil, err
	}
	return &PluginParser{pn: pn, ptl: ptl, pt: pt}, nil
}

// ParsePlugin is the input port of the PluginParser operation.
func (p *PluginParser) ParsePlugin(pd *gparselib.ParseData, ctx interface{},
) (*gparselib.ParseData, interface{}) {
	pEqual := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseLiteral(pd, ctx, nil, `=`)
	}
	pBig := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseAll(pd, ctx,
			[]gparselib.SubparserOp{p.pn.ParseNameIdent, ParseSpaceComment, pEqual, ParseSpaceComment, p.ptl.ParseTypeList},
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
		[]gparselib.SubparserOp{pBig, p.pt.ParseType},
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

// PluginListParser parses Plugins separated by a pipe '|' character.
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
type PluginListParser struct {
	pp *PluginParser
}

// NewPluginListParser creates a new parser for multiple titled type lists.
// If any regular expression used by the subparsers is invalid an error is
// returned.
func NewPluginListParser() (*PluginListParser, error) {
	pp, err := NewPluginParser()
	if err != nil {
		return nil, err
	}
	return &PluginListParser{pp: pp}, nil
}

// ParsePluginList is the input port of the PluginListParser
// operation.
func (p *PluginListParser) ParsePluginList(
	pd *gparselib.ParseData, ctx interface{},
) (*gparselib.ParseData, interface{}) {
	pBar := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseLiteral(pd, ctx, nil, `|`)
	}
	pAdditionalList := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseAll(
			pd, ctx,
			[]gparselib.SubparserOp{ParseSpaceComment, pBar, ParseSpaceComment, p.pp.ParsePlugin},
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
		[]gparselib.SubparserOp{p.pp.ParsePlugin, pAdditionalLists},
		parsePluginListSemantic,
	)
}
func parsePluginListSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
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

// FullPluginsParser parses the plugins of an operation starting with a '[' followed
// by a PluginList or a TypeList and a closing ']'.
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
type FullPluginsParser struct {
	pttl *PluginListParser
	ptl  *TypeListParser
}

// NewFullPluginsParser creates a new parser for the plugins of an operation.
// If any regular expression used by the subparsers is invalid an error is
// returned.
func NewFullPluginsParser() (*FullPluginsParser, error) {
	pttl, err := NewPluginListParser()
	if err != nil {
		return nil, err
	}
	ptl, err := NewTypeListParser()
	if err != nil {
		return nil, err
	}
	return &FullPluginsParser{pttl: pttl, ptl: ptl}, nil
}

// ParseFullPlugins is the input port of the FullPluginsParser operation.
func (p *FullPluginsParser) ParseFullPlugins(pd *gparselib.ParseData, ctx interface{},
) (*gparselib.ParseData, interface{}) {
	pList := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseBest(
			pd, ctx,
			[]gparselib.SubparserOp{p.pttl.ParsePluginList, p.ptl.ParseTypeList},
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
		parseFullPluginsSemantic,
	)
}
func parseFullPluginsSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
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

// ComponentParser parses a component including declaration and its plugins.
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
type ComponentParser struct {
	pcd *CompDeclParser
	pfp *FullPluginsParser
}

// NewParseComponent creates a new parser for a complete component.
// If any regular expression used by the subparsers is invalid an error is
// returned.
func NewParseComponent() (*ComponentParser, error) {
	pcd, err := NewCompDeclParser()
	if err != nil {
		return nil, err
	}
	pfp, err := NewFullPluginsParser()
	if err != nil {
		return nil, err
	}
	return &ComponentParser{pcd: pcd, pfp: pfp}, nil
}

// ParseComponent is the input port of the ComponentParser operation.
func (p *ComponentParser) ParseComponent(pd *gparselib.ParseData, ctx interface{},
) (*gparselib.ParseData, interface{}) {
	pPlugins := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseAll(
			pd, ctx,
			[]gparselib.SubparserOp{ParseSpaceComment, p.pfp.ParseFullPlugins},
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
		[]gparselib.SubparserOp{pOpen, ParseSpaceComment, p.pcd.ParseCompDecl, pOpt, ParseSpaceComment, pClose},
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
