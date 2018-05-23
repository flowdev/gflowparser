package parser

import (
	"github.com/flowdev/gflowparser/data"
	"github.com/flowdev/gparselib"
)

// ParsePort parses a port including optional index.
// Semantic result: data.Port
//
// flow:
//     in (ParseData)-> [pIndex gparselib.ParseAll
//                          [gparselib.ParseLiteral, gparselib.ParseNatural]
//                      ] -> out
//     in (ParseData)-> [pOpt gparselib.ParseOptional [pIndex]] -> out
//     in (ParseData)-> [gparselib.ParseAll [ParseNameIdent, pOpt]] -> out
//
// Details:
type ParsePort struct {
	pName *ParseNameIdent
}

// NewParsePort creates a new parser for a port.
// If any regular expression used by the subparsers is invalid an error is
// returned.
func NewParsePort() (*ParsePort, error) {
	pName, err := NewParseNameIdent()
	if err != nil {
		return nil, err
	}
	return &ParsePort{pName: pName}, nil
}

// In is the input port of the ParsePort operation.
func (p *ParsePort) In(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	pColon := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseLiteral(pd, ctx, nil, `:`)
	}
	pNumber := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		pd2, ctx2, err := gparselib.ParseNatural(pd, ctx, nil, 10)
		if err != nil {
			panic(err)
		}
		return pd2, ctx2
	}
	pIndex := func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseAll(pd, ctx,
			[]gparselib.SubparserOp{pColon, pNumber},
			func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
				pd2.Result.Value = pd2.SubResults[1].Value
				return pd2, ctx2
			},
		)
	}
	pOpt := func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseOptional(pd2, ctx2, pIndex, nil)
	}
	return gparselib.ParseAll(pd, ctx,
		[]gparselib.SubparserOp{p.pName.In, pOpt},
		parsePortSemantic,
	)
}
func parsePortSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	val1 := pd.SubResults[1].Value
	port := data.Port{
		Name: (pd.SubResults[0].Value).(string),
	}
	if val1 != nil {
		port.HasIndex = true
		port.Index = int((val1).(uint64))
	}
	pd.Result.Value = port
	return pd, ctx
}

// ParseArrow parses a flow arrow including ports and data types.
// Semantic result: data.Arrow
//
// flow:
//     in (ParseData)-> [pOptPort gparselib.ParseOptional [ParsePort]] -> out
//     in (ParseData)-> [pLeftParen gparselib.ParseLiteral] -> out
//     in (ParseData)-> [pRightParen gparselib.ParseLiteral] -> out
//     in (ParseData)-> [pArrow gparselib.ParseLiteral] -> out
//     in (ParseData)-> [pData gparselib.ParseAll
//                          [pLeftParen, ParseOptSpc,
//                           ParseTypeList, ParseOptSpc,
//                           pRightParen, ParseOptSpc
//                          ]
//                      ] -> out
//     in (ParseData)-> [pOptData gparselib.ParseOptional [pData]] -> out
//     in (ParseData)-> [gparselib.ParseAll
//                          [pOptPort, ParseSpaceComment, pOptData,
//                           pArrow, ParseSpaceComment, pOptPort
//                          ]
//                      ] -> out
//
// Details:
type ParseArrow struct {
	pPort *ParsePort
	pData *ParseTypeList
}

// NewParseArrow creates a new parser for a flow arrow.
// If any regular expression used by the subparsers is invalid an error is
// returned.
func NewParseArrow() (*ParseArrow, error) {
	pPort, err := NewParsePort()
	if err != nil {
		return nil, err
	}
	pData, err := NewParseTypeList()
	if err != nil {
		return nil, err
	}
	return &ParseArrow{pPort: pPort, pData: pData}, nil
}

// In is the input port of the ParsePort operation.
func (p *ParseArrow) In(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	pOptPort := func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseOptional(pd2, ctx2, p.pPort.In, nil)
	}
	pLeftParen := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseLiteral(pd, ctx, nil, `(`)
	}
	pRightParen := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseLiteral(pd, ctx, nil, `)`)
	}
	pArrow := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseLiteral(pd, ctx, nil, `->`)
	}
	pData := func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseAll(pd, ctx,
			[]gparselib.SubparserOp{
				pLeftParen, ParseOptSpc,
				p.pData.In, ParseOptSpc,
				pRightParen, ParseOptSpc,
			},
			func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
				pd2.Result.Value = pd2.SubResults[2].Value
				return pd2, ctx2
			},
		)
	}
	pOptData := func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseOptional(pd2, ctx2, pData, nil)
	}

	return gparselib.ParseAll(pd, ctx,
		[]gparselib.SubparserOp{
			pOptPort, ParseSpaceComment, pOptData,
			pArrow, ParseSpaceComment, pOptPort,
		},
		parseArrowSemantic,
	)
}
func parseArrowSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	val0 := pd.SubResults[0].Value
	val2 := pd.SubResults[2].Value
	val5 := pd.SubResults[5].Value
	arrow := data.Arrow{}
	if val0 != nil {
		port := (val0).(data.Port)
		arrow.FromPort = &port
	}
	if val2 != nil {
		arrow.Data = (val2).([]data.Type)
	}
	if val5 != nil {
		port := (val5).(data.Port)
		arrow.ToPort = &port
	}
	pd.Result.Value = arrow
	return pd, ctx
}

// ParseFlow parses a complete flow.
// Semantic result: data.Flow
//
// flow:
//     in (ParseData)-> [pOptArrow gparselib.ParseOptional [ParseArrow]] -> out
//     in (ParseData)-> [pPair gparselib.ParseAll
//                          [ParseComponent, ParseOptSpc,
//                           ParseArrow, ParseOptSpc
//                          ]
//                      ] -> out
//     in (ParseData)-> [pPairs gparselib.ParseMulti1 [pPair]] -> out
//     in (ParseData)-> [pPartLine gparselib.ParseAll
//                          [pOptArrow, ParseOptSpc, pPairs, ParseSpaceComment]
//                      ] -> out
//     in (ParseData)-> [gparselib.ParseMulti1 [pPartLine]] -> out
//
// Details:
type ParseFlow struct {
	pArrow *ParseArrow
	pComp  *ParseComponent
}

// NewParseFlow creates a new parser for a flow.
// If any regular expression used by the subparsers is invalid an error is
// returned.
func NewParseFlow() (*ParseFlow, error) {
	pArrow, err := NewParseArrow()
	if err != nil {
		return nil, err
	}
	pComp, err := NewParseComponent()
	if err != nil {
		return nil, err
	}
	return &ParseFlow{pArrow: pArrow, pComp: pComp}, nil
}

// In is the input port of the ParsePort operation.
func (p *ParseFlow) In(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	pOptArrow := func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseOptional(pd2, ctx2, p.pArrow.In, nil)
	}
	pPair := func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseAll(pd, ctx,
			[]gparselib.SubparserOp{
				p.pComp.In, ParseOptSpc,
				p.pArrow.In, ParseOptSpc,
			},
			func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
				pd2.Result.Value = []interface{}{
					pd2.SubResults[0].Value,
					pd2.SubResults[2].Value,
				}
				return pd2, ctx2
			},
		)
	}
	pPairs := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseMulti1(pd, ctx, pPair, parsePairsSemantic)
	}
	pPartLine := func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseAll(pd, ctx,
			[]gparselib.SubparserOp{
				pOptArrow, ParseOptSpc, pPairs, ParseSpaceComment,
			},
			parsePartLineSemantic,
		)
	}
	return gparselib.ParseMulti1(pd, ctx, pPartLine, parseFlowSemantic)
}
func parsePairsSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	parts := make([]interface{}, len(pd.SubResults)*2)
	j := 0
	for _, subResult := range pd.SubResults {
		pair := subResult.Value.([]interface{})
		parts[j] = pair[0]
		parts[j+1] = pair[1]
		j += 2
	}
	pd.Result.Value = parts
	return pd, ctx
}
func parsePartLineSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	partline := []interface{}{
		pd.SubResults[0].Value,
	}
	pairs := pd.SubResults[2].Value.([]interface{})
	pd.Result.Value = append(partline, pairs...)
	return pd, ctx
}
func parseFlowSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	lines := make([][]interface{}, len(pd.SubResults))
	for i, subResult := range pd.SubResults {
		line := subResult.Value.([]interface{})
		lines[i] = line
	}
	pd.Result.Value = data.Flow{
		Parts: lines,
	}
	return pd, ctx
}
