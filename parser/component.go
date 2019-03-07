package parser

import (
	"strings"

	"github.com/flowdev/gflowparser/data"
	"github.com/flowdev/gparselib"
)

// TypeParser parses a type declaration including optional package.
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

// ParseType parses a type declaration including optional package.
// * Semantic result: The optional package name and the local type name
//   including possible subtypes in case of a map or a list (data.Type).
//
// flow:
//     in (gparselib.ParseData)-> [pOpt gparselib.ParseOptional [ParsePackageIdent]] -> out
//     in (gparselib.ParseData)-> [gparselib.ParseAll [pOpt, ParseLocalTypeIdent]] -> out
func (p *TypeParser) ParseType(pd *gparselib.ParseData, ctx interface{},
) (*gparselib.ParseData, interface{}) {
	pCloseParen := gparselib.NewParseLiteralPlugin(TextSemantic, `)`)
	pList := gparselib.NewParseAllPlugin(
		[]gparselib.SubparserOp{
			gparselib.NewParseLiteralPlugin(nil, `list(`), ParseSpaceComment,
			p.ParseType, ParseSpaceComment,
			pCloseParen,
		},
		func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
			t := (pd2.SubResults[2].Value).(data.Type)
			pd2.Result.Value = data.Type{
				ListType: &t,
				SrcPos:   pd.Result.Pos,
			}
			return pd2, ctx2
		},
	)

	pMap := gparselib.NewParseAllPlugin(
		[]gparselib.SubparserOp{
			gparselib.NewParseLiteralPlugin(TextSemantic, `map(`), ParseSpaceComment,
			p.ParseType, ParseSpaceComment,
			gparselib.NewParseLiteralPlugin(nil, `,`), ParseSpaceComment,
			p.ParseType, ParseSpaceComment,
			pCloseParen,
		},
		func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
			tKey := (pd2.SubResults[2].Value).(data.Type)
			tValue := (pd2.SubResults[6].Value).(data.Type)
			pd2.Result.Value = data.Type{
				MapKeyType:   &tKey,
				MapValueType: &tValue,
				SrcPos:       pd.Result.Pos,
			}
			return pd2, ctx2
		},
	)

	pOptPack := gparselib.NewParseOptionalPlugin(p.pPack.ParsePackageIdent, nil)
	pSimpleType := gparselib.NewParseAllPlugin(
		[]gparselib.SubparserOp{pOptPack, p.pLocalType.ParseLocalTypeIdent},
		parseSimpleTypeSemantic,
	)

	return gparselib.ParseAny(
		pd, ctx,
		[]gparselib.SubparserOp{pList, pMap, pSimpleType},
		nil,
	)
}
func parseSimpleTypeSemantic(pd *gparselib.ParseData, ctx interface{},
) (*gparselib.ParseData, interface{}) {
	val0 := pd.SubResults[0].Value
	pack := ""
	if val0 != nil {
		pack = (val0).(string)
	}
	lType := (pd.SubResults[1].Value).(string)
	if pack == "" && (lType == "list" || lType == "map") {
		pd.AddError(pd.Result.Pos, "keyword '"+lType+"' not allowed as type", nil)
		pd.Result.Value = nil
		return pd, ctx
	}
	pd.Result.Value = data.Type{
		Package:   pack,
		LocalType: lType,
		SrcPos:    pd.Result.Pos,
	}
	return pd, ctx
}

// CompDeclParser is a parser for a component declaration.
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

// ParseCompDecl parses a component declaration.
// * Semantic result: The name and the type (data.CompDecl).
//
// flow:
//     in (gparselib.ParseData)-> [pAll gparselib.ParseAll [ParseNameIdent, ParseASpc]] -> out
//     in (gparselib.ParseData)-> [pOpt gparselib.ParseOptional [pAll]] -> out
//     in (gparselib.ParseData)-> [gparselib.ParseAll [pOpt, ParseType]] -> out
func (p *CompDeclParser) ParseCompDecl(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	pLong := gparselib.NewParseAllPlugin(
		[]gparselib.SubparserOp{p.pName.ParseNameIdent, ParseASpc, p.pType.ParseType},
		parseCompDeclSemantic,
	)
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

// TypeListParser is a parser for types separated by commas.
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

// ParseTypeList parses types separated by commas.
// * Semantic result: []data.Type
//
// flow:
//     in (gparselib.ParseData)-> [pAdditionalType gparselib.ParseAll
//                          [ParseSpaceComment, gparselib.ParseLiteral, ParseSpaceComment, ParseType]
//                      ] -> out
//     in (gparselib.ParseData)-> [pAdditionalTypes gparselib.ParseMulti0 [pAdditionalType]] -> out
//     in (gparselib.ParseData)-> [gparselib.ParseAll
//                          [ParseType, pAdditionalTypes]
//                      ] -> out
func (p *TypeListParser) ParseTypeList(pd *gparselib.ParseData, ctx interface{},
) (*gparselib.ParseData, interface{}) {
	pComma := gparselib.NewParseLiteralPlugin(nil, `,`)
	pAdditionalType := gparselib.NewParseAllPlugin(
		[]gparselib.SubparserOp{ParseSpaceComment, pComma, ParseSpaceComment, p.pt.ParseType},
		func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
			pd2.Result.Value = pd2.SubResults[3].Value
			return pd2, ctx2
		},
	)
	pAdditionalTypes := gparselib.NewParseMulti0Plugin(pAdditionalType, nil)
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

// PluginParser is a parser for a plugin.
type PluginParser struct {
	pn  *NameIdentParser
	ptl *TypeListParser
	pt  *TypeParser
}

// NewPluginParser creates a new parser for a plugin.
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

// ParsePlugin parses a name followed by an equals sign and types separated by commas.
// Alternatively a single type is parsed.
// * Semantic result: The title and a slice of data.Type (data.Plugin).
//
// flow:
//     in (gparselib.ParseData)-> [pFullPlugin gparselib.ParseAll
//                          [ ParseNameIdent, ParseSpaceComment, gparselib.ParseLiteral,
//                            ParseSpaceComment, ParseTypeList                          ]
//                      ] -> out
//     in (gparselib.ParseData)-> [gparselib.ParseAny [pFullPlugin, ParseType]] -> out
func (p *PluginParser) ParsePlugin(pd *gparselib.ParseData, ctx interface{},
) (*gparselib.ParseData, interface{}) {
	pEqual := gparselib.NewParseLiteralPlugin(nil, `=`)
	pFullPlugin := gparselib.NewParseAllPlugin(
		[]gparselib.SubparserOp{p.pn.ParseNameIdent, ParseSpaceComment, pEqual, ParseSpaceComment, p.ptl.ParseTypeList},
		func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
			val0 := pd2.SubResults[0].Value
			val4 := pd2.SubResults[4].Value
			pd2.Result.Value = data.Plugin{
				Name:   val0.(string),
				Types:  val4.([]data.Type),
				SrcPos: pd.Result.Pos,
			}
			return pd2, ctx2
		},
	)
	return gparselib.ParseAny(
		pd, ctx,
		[]gparselib.SubparserOp{pFullPlugin, p.pt.ParseType},
		func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
			if typ, ok := pd2.Result.Value.(data.Type); ok {
				pd2.Result.Value = data.Plugin{
					Types:  []data.Type{typ},
					SrcPos: pd.Result.Pos,
				}
			}
			return pd2, ctx2
		},
	)
}

// PluginListParser is a parser for Plugins separated by a pipe '|' character.
type PluginListParser struct {
	pp *PluginParser
}

// NewPluginListParser creates a new parser for multiple plugins.
// If any regular expression used by the subparsers is invalid an error is
// returned.
func NewPluginListParser() (*PluginListParser, error) {
	pp, err := NewPluginParser()
	if err != nil {
		return nil, err
	}
	return &PluginListParser{pp: pp}, nil
}

// ParsePluginList parses Plugins separated by a pipe '|' character.
// * Semantic result: A slice of data.Plugin.
//
// flow:
//     in (gparselib.ParseData)-> [pAdditionalList gparselib.ParseAll
//                          [ ParseSpaceComment, gparselib.ParseLiteral,
//                            ParseSpaceComment, ParsePlugin            ]
//                      ] -> out
//     in (gparselib.ParseData)-> [pAdditionalLists gparselib.ParseMulti0 [pAdditionalList]] -> out
//     in (gparselib.ParseData)-> [gparselib.ParseAll
//                          [ParsePlugin, pAdditionalLists]
//                      ] -> out
func (p *PluginListParser) ParsePluginList(
	pd *gparselib.ParseData, ctx interface{},
) (*gparselib.ParseData, interface{}) {
	pBar := gparselib.NewParseLiteralPlugin(nil, `|`)
	pAdditionalPlugin := gparselib.NewParseAllPlugin(
		[]gparselib.SubparserOp{ParseSpaceComment, pBar, ParseSpaceComment, p.pp.ParsePlugin},
		func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
			pd2.Result.Value = pd2.SubResults[3].Value
			return pd2, ctx2
		},
	)
	pAdditionalPlugins := gparselib.NewParseMulti0Plugin(pAdditionalPlugin, nil)
	return gparselib.ParseAll(
		pd, ctx,
		[]gparselib.SubparserOp{p.pp.ParsePlugin, pAdditionalPlugins},
		parsePluginListSemantic,
	)
}
func parsePluginListSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	firstPlugin := pd.SubResults[0].Value
	additionalPlugins := (pd.SubResults[1].Value).([]interface{})
	allPlugins := make([](data.Plugin), len(additionalPlugins)+1)
	allPlugins[0] = firstPlugin.(data.Plugin)

	for i, plug := range additionalPlugins {
		allPlugins[i+1] = plug.(data.Plugin)
	}
	pd.Result.Value = allPlugins
	return pd, ctx
}

// FullPluginsParser parses the plugins of an operation including '[' and ']'.
type FullPluginsParser struct {
	pttl *PluginListParser
	ptl  *TypeListParser
}

// NewFullPluginsParser creates a new parser for the plugins of a component.
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

// ParseFullPlugins parses the plugins of an operation starting with a '[' followed
// by a plugin list or a type list and a closing ']'.
// * Semantic result: A slice of data.Plugin.
//
// flow:
//     in (gparselib.ParseData)-> [pList gparselib.ParseAny
//                          [ParsePluginList, ParseTypeList]
//                      ] -> out
//     in (gparselib.ParseData)-> [gparselib.ParseAll
//                          [ gparselib.ParseLiteral, ParseSpaceComment, pList,
//                            ParseSpaceComment, gparselib.ParseLiteral         ]
//                      ] -> out
func (p *FullPluginsParser) ParseFullPlugins(pd *gparselib.ParseData, ctx interface{},
) (*gparselib.ParseData, interface{}) {
	pList := gparselib.NewParseBestPlugin(
		[]gparselib.SubparserOp{p.pttl.ParsePluginList, p.ptl.ParseTypeList},
		nil,
	)
	pOpen := gparselib.NewParseLiteralPlugin(nil, `[`)
	pClose := gparselib.NewParseLiteralPlugin(nil, `]`)
	return gparselib.ParseAll(
		pd, ctx,
		[]gparselib.SubparserOp{pOpen, ParseSpaceComment, pList, ParseSpaceComment, pClose},
		parseFullPluginsSemantic,
	)
}
func parseFullPluginsSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	list := pd.SubResults[2].Value
	if v, ok := list.([](data.Type)); ok {
		pd.Result.Value = [](data.Plugin){
			data.Plugin{Name: "", Types: v, SrcPos: v[0].SrcPos},
		}
	} else {
		pd.Result.Value = list
	}

	return pd, ctx
}

// ComponentParser is a parser for a component including declaration and plugins.
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

// ParseComponent parses a component including declaration and plugins.
// * Semantic result: A data.Component.
//
// flow:
//     in (gparselib.ParseData)-> [pPlugins gparselib.ParseAll
//                          [ParseSpaceComment, ParseFullPlugins]
//                      ] -> out
//     in (gparselib.ParseData)-> [pOpt gparselib.ParseOptional [pPlugins]] -> out
//     in (gparselib.ParseData)-> [gparselib.ParseAll
//                          [ gparselib.ParseLiteral, ParseSpaceComment, ParseCompDecl,
//                            pOpt, ParseSpaceComment, gparselib.ParseLiteral          ]
//                      ] -> out
func (p *ComponentParser) ParseComponent(pd *gparselib.ParseData, ctx interface{},
) (*gparselib.ParseData, interface{}) {
	pPlugins := gparselib.NewParseAllPlugin(
		[]gparselib.SubparserOp{ParseSpaceComment, p.pfp.ParseFullPlugins},
		func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
			pd2.Result.Value = pd2.SubResults[1].Value
			return pd2, ctx2
		},
	)
	pOpt := gparselib.NewParseOptionalPlugin(pPlugins, nil)
	pOpen := gparselib.NewParseLiteralPlugin(nil, `[`)
	pClose := gparselib.NewParseLiteralPlugin(nil, `]`)
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
		semVal.Plugins = (pd.SubResults[3].Value).([]data.Plugin)
	}
	pd.Result.Value = semVal
	return pd, ctx
}
