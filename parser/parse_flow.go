package parser

import (
	"fmt"
	"math"

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
	pName *NameIdentParser
}

// NewParsePort creates a new parser for a port.
// If any regular expression used by the subparsers is invalid an error is
// returned.
func NewParsePort() (*ParsePort, error) {
	pName, err := NewNameIdentParser()
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
		[]gparselib.SubparserOp{p.pName.ParseNameIdent, pOpt},
		parsePortSemantic,
	)
}
func parsePortSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	val1 := pd.SubResults[1].Value
	port := data.Port{
		Name:   (pd.SubResults[0].Value).(string),
		SrcPos: pd.Result.Pos,
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
//                          [pLeftParen, ParseSpaceComment,
//                           ParseTypeList, ParseSpaceComment,
//                           pRightParen, ParseOptSpc
//                          ]
//                      ] -> out
//     in (ParseData)-> [pOptData gparselib.ParseOptional [pData]] -> out
//     in (ParseData)-> [gparselib.ParseAll
//                          [pOptPort, ParseOptSpc, pOptData,
//                           pArrow, ParseOptSpc, pOptPort
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
				pLeftParen, ParseSpaceComment,
				p.pData.In, ParseSpaceComment,
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
			pOptPort, ParseOptSpc, pOptData,
			pArrow, ParseOptSpc, pOptPort,
		},
		parseArrowSemantic,
	)
}
func parseArrowSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	val0 := pd.SubResults[0].Value
	val2 := pd.SubResults[2].Value
	val5 := pd.SubResults[5].Value
	arrow := data.Arrow{SrcPos: pd.Result.Pos}
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
//     in (ParseData)-> [pAnyPart gparselib.ParseAny [ParseArrow, ParseComponent]] -> out
//     in (ParseData)-> [pFullPart gparselib.ParseAll [pAnyPart, ParseOptSpc]] -> out
//     in (ParseData)-> [pPartString gparselib.ParseMulti [pFullPart]] -> out
//     in (ParseData)-> [pPartLine gparselib.ParseAll
//                          [pPartString, ParseStatementEnd]
//                      ] -> out
//     in (ParseData)-> [gparselib.ParseMulti1 [pPartLine]] -> out
//
// Details:
type ParseFlow struct {
	pArrow *ParseArrow
	pComp  *ParseComponent
}

// Error messages for semantic errors.
const (
	errMsg2Arrows = "A flow line must contain alternating arrows and components " +
		"but this one has got two consecutive arrows at position %d"
	errMsg2Comps = "A flow line must contain alternating arrows and components " +
		"but this one has got two consecutive components at position %d"
	errMsgPartType = "A flow line must only contain arrows and components " +
		"but this one has got a %T at position %d"
	errMsgFirstPort = "The first arrow of this flow line is missing a source port"
	errMsgLastPort  = "The last arrow of this flow line is missing a destination port"
	errMsgFirstData = "The first arrow of this flow line is missing its data declaration"
)

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

// In is the input port of the ParseFlow operation.
func (p *ParseFlow) In(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	pAnyPart := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseAny(
			pd, ctx,
			[]gparselib.SubparserOp{p.pArrow.In, p.pComp.In},
			nil,
		)
	}
	pFullPart := func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseAll(pd, ctx,
			[]gparselib.SubparserOp{pAnyPart, ParseOptSpc},
			func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
				pd2.Result.Value = pd2.SubResults[0].Value
				return pd2, ctx2
			},
		)
	}
	pPartString := func(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseMulti(pd, ctx, pFullPart, nil, 2, math.MaxInt32)
	}
	pPartLine := func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseAll(pd, ctx,
			[]gparselib.SubparserOp{pPartString, ParseStatementEnd},
			parsePartLineSemantic,
		)
	}
	return gparselib.ParseMulti1(pd, ctx, pPartLine, parseFlowSemantic)
}
func parsePartLineSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	partLine := pd.SubResults[0].Value.([]interface{})
	n := len(partLine)

	var lastIsArrow, lastIsComp bool
	for i, part := range partLine {
		switch v := part.(type) {
		case data.Arrow:
			if lastIsArrow {
				pd.AddError(v.SrcPos, fmt.Sprintf(errMsg2Arrows, i+1), nil)
				return pd, ctx
			}
			lastIsArrow = true
			lastIsComp = false
		case data.Component:
			if lastIsComp {
				pd.AddError(v.SrcPos, fmt.Sprintf(errMsg2Comps, i+1), nil)
				return pd, ctx
			}
			lastIsComp = true
			lastIsArrow = false
		default:
			pd.AddError(pd.Result.Pos, fmt.Sprintf(errMsgPartType, part, i+1), nil)
			return pd, ctx
		}
	}
	var firstArrow data.Arrow
	if v, ok := partLine[0].(data.Arrow); ok {
		firstArrow = v
		if firstArrow.FromPort == nil {
			pd.AddError(pd.Result.Pos, errMsgFirstPort, nil)
		}
	} else {
		firstArrow = partLine[1].(data.Arrow)
	}
	if len(firstArrow.Data) == 0 {
		pd.AddError(pd.Result.Pos, errMsgFirstData, nil)
	}
	if lastArrow, ok := partLine[n-1].(data.Arrow); ok {
		if lastArrow.ToPort == nil {
			pd.AddError(pd.Result.Pos, errMsgLastPort, nil)
		}
	}
	if !pd.Result.HasError() {
		pd.Result.Value = partLine
	}
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
