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
// * Semantic result: data.Port
//
// flow:
//     in (gparselib.ParseData)-> [pIndex gparselib.ParseAll [gparselib.ParseLiteral, gparselib.ParseNatural]] -> out
//     in (gparselib.ParseData)-> [pOptIdx gparselib.ParseOptional [pIndex]] -> out
//     in (gparselib.ParseData)-> [pNormPort gparselib.ParseAll [ParseNameIdent, pOptIdx]] -> out
//     in (gparselib.ParseData)-> [pDots gparselib.ParseLiteral] -> out
//     in (gparselib.ParseData)-> [pContinuation gparselib.ParseAll [pDots, gparselib.ParseNatural]] -> out
//     in (gparselib.ParseData)-> [gparselib.ParseAll [pContinuation, pNormPort]] -> out
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

// MultiTypeListParser is a parser for multiple type lists.
type MultiTypeListParser struct {
	ptl *TypeListParser
}

// NewMultiTypeListParser creates a new parser for type lists (types separated
// by ',') separated by '|'.
// If a subparser can't be created an error is returned.
func NewMultiTypeListParser() (*MultiTypeListParser, error) {
	ptl, err := NewTypeListParser()
	if err != nil {
		return nil, err
	}
	return &MultiTypeListParser{ptl: ptl}, nil
}

// ParseMultiTypeList parses multiple type lists (types separated
// by ',') separated by '|'.
// * Semantic result: []data.Type containing data.SeparatorType
//
// flow:
//     in (gparselib.ParseData)-> [pAdditionalTypeList gparselib.ParseAll
//                          [ParseSpaceComment, gparselib.ParseLiteral, ParseSpaceComment, ParseTypeList]
//                      ] -> out
//     in (gparselib.ParseData)-> [pAdditionalTypeLists gparselib.ParseMulti0 [pAdditionalTypeList]] -> out
//     in (gparselib.ParseData)-> [gparselib.ParseAll
//                          [ParseTypeList, pAdditionalTypeLists]
//                      ] -> out
func (p *MultiTypeListParser) ParseMultiTypeList(pd *gparselib.ParseData, ctx interface{},
) (*gparselib.ParseData, interface{}) {
	pBar := gparselib.NewParseLiteralPlugin(nil, `|`)
	pAdditionalTypeList := gparselib.NewParseAllPlugin(
		[]gparselib.SubparserOp{ParseSpaceComment, pBar, ParseSpaceComment, p.ptl.ParseTypeList},
		func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
			pd2.Result.Value = pd2.SubResults[3].Value
			return pd2, ctx2
		},
	)
	pAdditionalTypeLists := gparselib.NewParseMulti0Plugin(pAdditionalTypeList, nil)
	return gparselib.ParseAll(
		pd, ctx,
		[]gparselib.SubparserOp{p.ptl.ParseTypeList, pAdditionalTypeLists},
		parseMultiTypeListSemantic,
	)
}
func parseMultiTypeListSemantic(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	firstList := pd.SubResults[0].Value
	additionalLists := (pd.SubResults[1].Value).([]interface{})
	alltypes := make([]data.Type, 0, 64)
	alltypes = append(alltypes, firstList.([]data.Type)...)

	for _, list := range additionalLists {
		alltypes = append(alltypes, data.SeparatorType)
		alltypes = append(alltypes, list.([]data.Type)...)
	}
	pd.Result.Value = alltypes
	return pd, ctx
}

// ArrowParser is a parser for a flow arrow including ports and data types.
type ArrowParser struct {
	pPort *PortParser
	pData *MultiTypeListParser
}

// NewArrowParser creates a new parser for a flow arrow.
// If any regular expression used by the subparsers is invalid an error is
// returned.
func NewArrowParser() (*ArrowParser, error) {
	pPort, err := NewPortParser()
	if err != nil {
		return nil, err
	}
	pData, err := NewMultiTypeListParser()
	if err != nil {
		return nil, err
	}
	return &ArrowParser{pPort: pPort, pData: pData}, nil
}

// ParseArrow parses a flow arrow including ports and data types.
// * Semantic result: data.Arrow
//
// flow:
//     in (gparselib.ParseData)-> [pOptPort gparselib.ParseOptional [ParsePort]] -> out
//     in (gparselib.ParseData)-> [pLeftParen gparselib.ParseLiteral] -> out
//     in (gparselib.ParseData)-> [pRightParen gparselib.ParseLiteral] -> out
//     in (gparselib.ParseData)-> [pArrow gparselib.ParseLiteral] -> out
//     in (gparselib.ParseData)-> [pData gparselib.ParseAll
//                          [pLeftParen, ParseSpaceComment,
//                           ParseMultiTypeList, ParseSpaceComment,
//                           pRightParen, ParseOptSpc
//                          ]
//                      ] -> out
//     in (gparselib.ParseData)-> [pOptData gparselib.ParseOptional [pData]] -> out
//     in (gparselib.ParseData)-> [gparselib.ParseAll
//                          [pOptPort, ParseOptSpc, pOptData,
//                           pArrow, ParseOptSpc, pOptPort
//                          ]
//                      ] -> out
func (p *ArrowParser) ParseArrow(pd *gparselib.ParseData, ctx interface{}) (*gparselib.ParseData, interface{}) {
	pOptPort := gparselib.NewParseOptionalPlugin(p.pPort.ParsePort, nil)
	pLeftParen := gparselib.NewParseLiteralPlugin(nil, `(`)
	pRightParen := gparselib.NewParseLiteralPlugin(nil, `)`)
	pArrow := gparselib.NewParseLiteralPlugin(nil, `->`)
	pData := gparselib.NewParseAllPlugin(
		[]gparselib.SubparserOp{
			pLeftParen, ParseSpaceComment,
			p.pData.ParseMultiTypeList, ParseSpaceComment,
			pRightParen, ParseOptSpc,
		},
		func(pd2 *gparselib.ParseData, ctx2 interface{}) (*gparselib.ParseData, interface{}) {
			pd2.Result.Value = pd2.SubResults[2].Value
			return pd2, ctx2
		},
	)
	pOptData := gparselib.NewParseOptionalPlugin(pData, nil)

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

// FlowParser is a parser for a complete flow.
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

// ParseFlow parses a complete flow.
// * Semantic result: data.Flow
//
// flow:
//     in (gparselib.ParseData)-> [pAnyPart gparselib.ParseAny [ParseArrow, ParseComponent]] -> out
//     in (gparselib.ParseData)-> [pFullPart gparselib.ParseAll [pAnyPart, ParseOptSpc]] -> out
//     in (gparselib.ParseData)-> [pPartSequence gparselib.ParseMulti [pFullPart]] -> out
//     in (gparselib.ParseData)-> [pPartLine gparselib.ParseAll
//                          [pPartSequence, ParseStatementEnd]
//                      ] -> out
//     in (gparselib.ParseData)-> [pLines gparselib.ParseMulti1 [pPartLine]] -> out
//     in (gparselib.ParseData)-> [gparselib.ParseAll [pLines, gparselib.ParseEOF]] -> out
func (p *FlowParser) ParseFlow(pd *gparselib.ParseData, ctx interface{},
) (*gparselib.ParseData, interface{}) {
	pAnyPart := gparselib.NewParseAnyPlugin(
		[]gparselib.SubparserOp{p.pArrow.ParseArrow, p.pComp.ParseComponent},
		nil,
	)
	pFullPart := gparselib.NewParseAllPlugin(
		[]gparselib.SubparserOp{pAnyPart, ParseOptSpc},
		func(pd3 *gparselib.ParseData, ctx3 interface{}) (*gparselib.ParseData, interface{}) {
			pd3.Result.Value = pd3.SubResults[0].Value
			return pd3, ctx3
		},
	)
	pPartSequence := gparselib.NewParseMultiPlugin(pFullPart, nil, 2, math.MaxInt32)
	pPartLine := gparselib.NewParseAllPlugin(
		[]gparselib.SubparserOp{pPartSequence, ParseStatementEnd},
		parsePartLineSemantic,
	)
	pLines := gparselib.NewParseMulti1Plugin(pPartLine, parseFlowSemantic)
	pEOF := gparselib.NewParseEOFPlugin(nil)
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
	// TODO: keep this lenient data parsing???
	//if len(firstArrow.Data) == 0 {
	//	pd.AddError(pd.Result.Pos, errMsgFirstData, nil)
	//}
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
