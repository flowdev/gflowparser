package parser

import (
	"fmt"
	"math"

	"github.com/flowdev/gflowparser/data"
	"github.com/flowdev/gparselib"
)

// PortParser is a parser for a port including optional index.
type PortParser struct {
	pName *NameIdentParser
}

// NewPortParser creates a new parser for a port.
// If any regular expression used by the subparsers is invalid an error is
// returned.
func NewPortParser() (*PortParser, error) {
	pName, err := NewNameIdentParser()
	if err != nil {
		return nil, err
	}
	return &PortParser{pName: pName}, nil
}

// ParsePort parses a port including optional index.
// Semantic result: data.Port
//
// flow:
//     in (ParseData)-> [pColon gparselib.ParseLiteral] -> out
//     in (ParseData)-> [pNumber gparselib.ParseNatural] -> out
//     in (ParseData)-> [pIndex gparselib.ParseAll [pColon, pNumber]] -> out
//     in (ParseData)-> [pOptIdx gparselib.ParseOptional [pIndex]] -> out
//     in (ParseData)-> [pNormPort gparselib.ParseAll [ParseNameIdent, pOptIdx]] -> out
//     in (ParseData)-> [pDots gparselib.ParseLiteral] -> out
//     in (ParseData)-> [pContinuation gparselib.ParseAll [pDots, pNumber]] -> out
//     in (ParseData)-> [gparselib.ParseAll [pNormPort, pContinuation]] -> out
func (p *PortParser) ParsePort(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	pColon := gparselib.NewParseLiteralPlugin(nil, `:`)
	pNumber, err := gparselib.NewParseNaturalPlugin(nil, 10)
	if err != nil {
		panic(err)
	}
	pIndex := gparselib.NewParseAllPlugin([]gparselib.SubparserOp{pColon, pNumber},
		func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
			pd2.Result.Value = pd2.SubResults[1].Value
			return pd2, ctx2
		},
	)
	pOptIdx := gparselib.NewParseOptionalPlugin(pIndex, nil)
	pNormPort := gparselib.NewParseAllPlugin([]gparselib.SubparserOp{p.pName.ParseNameIdent, pOptIdx}, parsePortSemantic)

	pDots := gparselib.NewParseLiteralPlugin(nil, `...`)
	pContinuation := gparselib.NewParseAllPlugin([]gparselib.SubparserOp{pDots, pNumber}, parseContinuationPortSemantic)

	return gparselib.ParseAny(pd, ctx, []gparselib.SubparserOp{pContinuation, pNormPort}, nil)
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
func parseContinuationPortSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	val1 := pd.SubResults[1].Value.(uint64)
	pd.Result.Value = data.Port{
		Name:   data.ContinuationSignal,
		Index:  int(val1),
		SrcPos: pd.Result.Pos,
	}
	return pd, ctx
}

// ArrowParser parses a flow arrow including ports and data types.
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
type ArrowParser struct {
	pPort *PortParser
	pData *TypeListParser
}

// NewArrowParser creates a new parser for a flow arrow.
// If any regular expression used by the subparsers is invalid an error is
// returned.
func NewArrowParser() (*ArrowParser, error) {
	pPort, err := NewPortParser()
	if err != nil {
		return nil, err
	}
	pData, err := NewTypeListParser()
	if err != nil {
		return nil, err
	}
	return &ArrowParser{pPort: pPort, pData: pData}, nil
}

// ParseArrow is the input port of the ArrowParser operation.
func (p *ArrowParser) ParseArrow(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	pOptPort := func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseOptional(pd2, ctx2, p.pPort.ParsePort, nil)
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
				p.pData.ParseTypeList, ParseSpaceComment,
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

// FlowParser parses a complete flow.
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
type FlowParser struct {
	pArrow *ArrowParser
	pComp  *ComponentParser
}

// Error messages for semantic errors.
const (
	errMsg2Arrows = "A flow line must contain alternating arrows and components " +
		"but this one has got two consecutive arrows at position %d"
	errMsg2Comps = "A flow line must contain alternating arrows and components " +
		"but this one has got two consecutive components at position %d"
	errMsgPartType = "A flow line must only contain arrows and components " +
		"but this one has got a %T at position %d"
	errMsgFirstPort     = "The first arrow of this flow line is missing a source port"
	errMsgLastPort      = "The last arrow of this flow line is missing a destination port"
	errMsgFirstData     = "The first arrow of this flow line is missing its data declaration"
	errMsgContInMidLine = "A continuation is only allowed at the very start or end of a flow line, " +
		"but not at position %d"
	errMsgContNoMatch = "The continuation at the very end of flow line %d is missing its counter part"
	errMsgContStart   = "The continuation at the very start of flow line %d can't exists before its counter part"
	errMsg2ContEnd    = "The continuation at the very end of flow line %d is doubled at flow line %d"
	errMsgContData    = "The continuation at the very end of flow line %d has got an invalid data annotation"
)

// NewFlowParser creates a new parser for a flow.
// If any regular expression used by the subparsers is invalid an error is
// returned.
func NewFlowParser() (*FlowParser, error) {
	pArrow, err := NewArrowParser()
	if err != nil {
		return nil, err
	}
	pComp, err := NewParseComponent()
	if err != nil {
		return nil, err
	}
	return &FlowParser{pArrow: pArrow, pComp: pComp}, nil
}

// ParseFlow is the input port of the FlowParser operation.
func (p *FlowParser) ParseFlow(pd *gparselib.ParseData, ctx interface{},
) (*gparselib.ParseData, interface{}) {
	pAnyPart := func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseAny(
			pd2, ctx2,
			[]gparselib.SubparserOp{p.pArrow.ParseArrow, p.pComp.ParseComponent},
			nil,
		)
	}
	pFullPart := func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseAll(pd2, ctx2,
			[]gparselib.SubparserOp{pAnyPart, ParseOptSpc},
			func(pd3 *gparselib.ParseData, ctx3 interface{}) (*gparselib.ParseData, interface{}) {
				pd3.Result.Value = pd3.SubResults[0].Value
				return pd3, ctx3
			},
		)
	}
	pPartString := func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseMulti(pd2, ctx2, pFullPart, nil, 2, math.MaxInt32)
	}
	pPartLine := func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseAll(pd2, ctx2,
			[]gparselib.SubparserOp{pPartString, ParseStatementEnd},
			parsePartLineSemantic,
		)
	}
	pLines := func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseMulti1(pd2, ctx2, pPartLine, parseFlowSemantic)
	}
	pEOF := func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
		return gparselib.ParseEOF(pd2, ctx2, nil)
	}
	return gparselib.ParseAll(pd, ctx,
		[]gparselib.SubparserOp{pLines, pEOF},
		func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
			pd2.Result.Value = pd2.SubResults[0].Value
			return pd2, ctx2
		},
	)
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
			if (v.FromPort != nil && v.FromPort.Continuation() && i != 0) ||
				(v.ToPort != nil && v.ToPort.Continuation() && i != n-1) {

				pd.AddError(v.SrcPos, fmt.Sprintf(errMsgContInMidLine, i+1), nil)
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
	} else {
		pd.Result.Value = nil
		pd.ResetSourcePos(-1)
	}
	return pd, ctx
}
func parseFlowSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	lines := make([][]interface{}, len(pd.SubResults))
	for i, subResult := range pd.SubResults {
		line := subResult.Value.([]interface{})
		lines[i] = line
	}
	pd = checkContinuations(lines, pd)
	if !pd.Result.HasError() {
		pd.Result.Value = data.Flow{
			Parts: lines,
		}
	} else {
		pd.Result.Value = nil
		pd.ResetSourcePos(-1)
	}
	return pd, ctx
}
func checkContinuations(lines [][]interface{}, pd *gparselib.ParseData) *gparselib.ParseData {
	endConts := make(map[int]int, 64)
	for i, line := range lines {
		if v, ok := line[0].(data.Arrow); ok {
			if v.FromPort.Continuation() {
				if _, ok := endConts[v.FromPort.Index]; ok {
					delete(endConts, v.FromPort.Index)
				} else {
					pd.AddError(v.FromPort.SrcPos, fmt.Sprintf(errMsgContStart, i+1), nil)
				}
			}
		}
		n := len(line)
		if v, ok := line[n-1].(data.Arrow); ok {
			if v.ToPort.Continuation() {
				if len(v.Data) > 0 {
					pd.AddError(v.Data[0].SrcPos, fmt.Sprintf(errMsgContData, i+1), nil)
				}
				if j, ok := endConts[v.ToPort.Index]; ok {
					pd.AddError(v.ToPort.SrcPos, fmt.Sprintf(errMsg2ContEnd, j+1, i+1), nil)
				} else {
					endConts[v.ToPort.Index] = i
				}
			}
		}
	}
	for _, v := range endConts {
		line := lines[v]
		n := len(line)
		arr := line[n-1].(data.Arrow)
		pd.AddError(arr.ToPort.SrcPos, fmt.Sprintf(errMsgContNoMatch, v+1), nil)
	}

	return pd
}
