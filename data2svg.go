package gflowparser

import (
	"bytes"
	"errors"
	"fmt"
	"sort"

	"github.com/flowdev/gflowparser/data"
	"github.com/flowdev/gflowparser/parser"
	"github.com/flowdev/gflowparser/svg"
	"github.com/flowdev/gparselib"
)

// Error messages.
const (
	errMsgDeclAndRef = "in one flow line the component '%s' is declared at " +
		"index [%d, %d] and referenced at index [%d, %d]"
	errMsg2Decls = "a component with the name '%s' is declared in the " +
		"two index positions [%d, %d] and [%d, %d]"
	errMsgPartType = "Found illegal flow part type '%T' at index [%d, %d]"
)

type declOrRef struct {
	i, j                 int
	isDecl               bool
	splitRefs, mergeRefs int
}

// TODO: Prevent circles (for now)!
// TODO: Handle lines with multiple splits/merges (clusterStart, clusterEnd, ...)
// TODO: Restructure easily
type decl struct {
	name         string
	i, j         int
	svgOp        *svg.Op
	svgMerge     *svg.Merge
	svgSplit     *svg.Split
	clusterStart int
	clusterEnd   int
}

// checkParserFeedback converts parser errors into a single error.
func checkParserFeedback(pd *gparselib.ParseData) (string, error) {
	if pd.Result.HasError() {
		return "", errors.New("Found errors while parsing flow:\n" +
			feedbackToString(pd))
	}
	return feedbackToString(pd), nil
}
func feedbackToString(pd *gparselib.ParseData) string {
	buf := bytes.Buffer{}
	for _, fb := range pd.Result.Feedback {
		buf.WriteString(fb.String())
		buf.WriteString("\n")
	}
	return buf.String()
}

// parserPartsToSVGData converts the parts of a data.Flow one to one into SVG
// diagram shapes.
// So this operation does only a simple translation but doesn't change the form
// of the part table.
// Components are special since they can be translated in 3 ways:
// 1. Into a decl struct if it is a declaration (the first occurence).
// 2. Into a merge if there are more parts before it.
// 3. Into a string (the name of the decl) if it is only used for a split.
func parserPartsToSVGData(flowDat data.Flow,
) (flow [][]interface{}, decls map[string]*decl, err error) {
	svgDat := make([][]interface{}, len(flowDat.Parts))
	decls = make(map[string]*decl)

	for i, partLine := range flowDat.Parts {
		m := len(partLine) - 1
		svgLine := make([]interface{}, m+1)
		for j, part := range partLine {
			switch p := part.(type) {
			case data.Arrow:
				svgLine[j] = arrowToSVGData(p, j > 0, j < m)
			case data.Component:
				if dcl, ok := decls[p.Decl.Name]; ok {
					if p.Decl.VagueType {
						if dcl.clusterStart >= i {
							return nil, nil, fmt.Errorf(errMsgDeclAndRef,
								dcl.name, dcl.i, dcl.j, i, j)
						}
					} else {
						return nil, nil, fmt.Errorf(errMsg2Decls,
							dcl.name, dcl.i, dcl.j, i, j)
					}
					if j > 0 {
						dcl.svgMerge.Size++
						dcl.clusterEnd = max(dcl.clusterEnd, i)
						svgLine[j] = dcl.svgMerge
					}
					if j < m {
						if svgLine[j] == nil {
							svgLine[j] = dcl.name
						}
					}
				} else {
					dcl := &decl{
						name:         p.Decl.Name,
						i:            i,
						j:            j,
						svgOp:        compToSVGData(p),
						svgMerge:     &svg.Merge{ID: p.Decl.Name},
						clusterStart: i,
						clusterEnd:   i,
					}
					decls[dcl.name] = dcl
					svgLine[j] = dcl
				}
			default:
				panic(fmt.Sprintf(errMsgPartType, part, i, j))
			}
		}
		svgDat[i] = svgLine
	}
	return svgDat, decls, nil
}

func arrowToSVGData(arr data.Arrow, hasSrcOp, hasDstOp bool) *svg.Arrow {
	return &svg.Arrow{
		DataType: arrDataToSVGData(arr.Data),
		HasSrcOp: hasSrcOp, SrcPort: portToSVGData(arr.FromPort),
		HasDstOp: hasDstOp, DstPort: portToSVGData(arr.ToPort),
	}
}

func compToSVGData(comp data.Component) *svg.Op {
	plugs := make([]*svg.Plugin, len(comp.Plugins))
	for i, plug := range comp.Plugins {
		plugs[i] = pluginToSVGData(plug)
	}
	return &svg.Op{
		Main:    &svg.Rect{Text: compDeclToSVGData(comp.Decl)},
		Plugins: plugs,
	}
}

func compDeclToSVGData(decl data.CompDecl) []string {
	typ := typeToSVGData(decl.Type)
	if typ == decl.Name {
		return []string{decl.Name}
	}
	return []string{decl.Name, typ}
}

func pluginToSVGData(plug data.NameNTypes) *svg.Plugin {
	rects := make([]*svg.Rect, len(plug.Types))
	for i, typ := range plug.Types {
		rects[i] = &svg.Rect{
			Text: []string{typeToSVGData(typ)},
		}
	}
	return &svg.Plugin{
		Title: plug.Name,
		Rects: rects,
	}
}

func arrDataToSVGData(dat []data.Type) string {
	if len(dat) == 0 {
		return ""
	}
	buf := bytes.Buffer{}
	buf.WriteString("(")
	buf.WriteString(typeToSVGData(dat[0]))
	for i := 1; i < len(dat); i++ {
		buf.WriteString(", ")
		buf.WriteString(typeToSVGData(dat[i]))
	}
	buf.WriteString(")")
	return buf.String()
}

func portToSVGData(port *data.Port) string {
	if port == nil {
		return ""
	}
	if port.HasIndex {
		return fmt.Sprintf("%s[%d]", port.Name, port.Index)
	}
	return port.Name
}

func typeToSVGData(typ data.Type) string {
	if typ.Package != "" {
		return typ.Package + "." + typ.LocalType
	}
	return typ.LocalType
}

// sortAndUniqIdxs sorts the indices of declarations and makes them unique.
// WARNING: The second index has to be sorted in decreasing order or splits
// can't be added easily.
func sortAndUniqDeclRefs(declRefs []declOrRef) []declOrRef {
	if len(declRefs) == 0 {
		return declRefs
	}
	sort.Slice(declRefs, func(i, j int) bool {
		return declRefs[i].i < declRefs[j].i ||
			(declRefs[i].i == declRefs[j].i && declRefs[i].j >= declRefs[j].j)
	})

	i := 0
	for j := 1; i < len(declRefs); j++ {
		if declRefs[i].i != declRefs[j].i || declRefs[i].j != declRefs[j].j {
			i++
			declRefs[i] = declRefs[j]
		}
	}
	return declRefs[:i+1]
}

// enhanceDecls adds knowledge to the declarations declOrRef
// whether the decl is target for a split and the number of merges.
func enhanceDecls(declRefs []declOrRef, flow *svg.Flow) []declOrRef {
	svgDat := flow.Shapes
	for i, decl := range declRefs {
		if !decl.isDecl {
			continue
		}
		split, merge := 0, 0
		for _, refRef := range declRefs {
			if refRef.isDecl {
				continue
			}
			ref := svgDat[refRef.i][refRef.j].(declOrRef)
			if ref.i != decl.i || ref.j != decl.j {
				continue
			}
			sl := svgDat[refRef.i]
			if refRef.j < len(sl)-1 {
				split = 1
			}
			if refRef.j > 0 {
				merge++
			}
		}
		decl.splitRefs = split
		decl.mergeRefs = merge
		declRefs[i] = decl
	}
	return declRefs
}

type mergeData struct {
	line                   []interface{}
	svgMerge               *svg.Merge
	found                  int
	lastMergeI, lastMergeJ int
}

// addSplitsAndMergesToSVGData rearranges SVG flow shapes to accomodate splits
// and merges in flows.
func addSplitsAndMergesToSVGData(flow *svg.Flow, declRefs []declOrRef,
) (*svg.Flow, []mergeData) {
	splits := make(map[int]map[int]*svg.Split)
	merges := make(map[int]map[int]mergeData)
	svgDat := flow.Shapes
	for _, dor := range declRefs {
		if dor.isDecl {
			svgDat, splits, merges = handleDecl(svgDat, dor, splits, merges)
		} else {
			svgDat = handleRef(svgDat, dor, splits, merges)
		}
	}
	ms := simplifyMerges(merges)
	ms = sortMerges(ms)
	return &svg.Flow{Shapes: svgDat}, ms
}
func handleDecl(
	svgDat [][]interface{},
	decl declOrRef,
	splits map[int]map[int]*svg.Split,
	merges map[int]map[int]mergeData,
) ([][]interface{}, map[int]map[int]*svg.Split, map[int]map[int]mergeData) {
	// Splits:
	// Splits have to be directly after the decl/comp.
	if decl.splitRefs > 0 {
		var split *svg.Split
		svgDat, split = addSplit(svgDat, decl.i, decl.j)
		rememberSplit(splits, split, decl.i, decl.j)
	}

	// Merges:
	// Merges replace the ref/decl and end the line.
	// The decl/comp itself and everything after that has to follow on an
	// own line directly after the last merge.
	if decl.mergeRefs > 0 {
		md := mergeData{
			svgMerge: &svg.Merge{
				ID:   mergeID(svgDat, decl.i, decl.j),
				Size: decl.mergeRefs + 1,
			},
			line:  svgDat[decl.i][decl.j:],
			found: 1,
		}
		rememberMerge(merges, md, decl.i, decl.j)
		svgDat[decl.i][decl.j] = md.svgMerge
		svgDat[decl.i] = svgDat[decl.i][:decl.j+1]
	}
	return svgDat, splits, merges
}
func handleRef(
	svgDat [][]interface{},
	refRef declOrRef,
	splits map[int]map[int]*svg.Split,
	merges map[int]map[int]mergeData,
) [][]interface{} {
	sl := svgDat[refRef.i]
	ref := sl[refRef.j].(declOrRef)

	// Splits:
	// Add rest of line to split
	if refRef.j < len(sl)-1 {
		split := splits[ref.i][ref.j]
		split.Shapes = append(split.Shapes, sl[refRef.j+1:])
		sl = sl[:refRef.j]
		svgDat[refRef.i] = sl
	}

	// Merges:
	// Merges replace the ref/decl and end the line.
	// The decl/comp itself and everything after that has to follow on an
	// own line directly after the last merge.
	if refRef.j > 0 {
		merge := merges[ref.i][ref.j]
		if len(sl) > refRef.j {
			sl[refRef.j] = merge.svgMerge
		} else {
			sl = append(sl, merge.svgMerge)
			svgDat[refRef.i] = sl
		}
		merge.found++
		if merge.found == merge.svgMerge.Size {
			merge.lastMergeI = refRef.i
			merge.lastMergeJ = refRef.j
		}
	}
	return svgDat
}
func rememberSplit(splits map[int]map[int]*svg.Split, split *svg.Split, i, j int) {
	subMap := splits[i]
	if subMap == nil {
		subMap = make(map[int]*svg.Split)
		splits[i] = subMap
	}
	subMap[j] = split
}
func rememberMerge(merges map[int]map[int]mergeData, merge mergeData, i, j int) {
	subMap := merges[i]
	if subMap == nil {
		subMap = make(map[int]mergeData)
		merges[i] = subMap
	}
	subMap[j] = merge
}
func addSplit(svgDat [][]interface{}, i, j int) ([][]interface{}, *svg.Split) {
	sl := svgDat[i]
	rest := sl[j+1:]
	split := &svg.Split{
		Shapes: make([][]interface{}, 0, 16),
	}
	if len(rest) > 0 {
		split.Shapes = append(split.Shapes, rest)
	}
	sl = append(sl[:j+1], split)
	svgDat[i] = sl
	return svgDat, split
}
func mergeID(svgData [][]interface{}, i, j int) string {
	return svgData[i][j].(svg.Op).Main.Text[0]
}
func simplifyMerges(merges map[int]map[int]mergeData) []mergeData {
	ms := make([]mergeData, 0, 64)
	for _, mrgs := range merges {
		for _, m := range mrgs {
			ms = append(ms, m)
		}
	}
	return ms
}
func sortMerges(ms []mergeData) []mergeData {
	sort.Slice(ms, func(i, j int) bool {
		return ms[i].lastMergeI > ms[j].lastMergeI ||
			(ms[i].lastMergeI == ms[j].lastMergeI &&
				ms[i].lastMergeJ >= ms[j].lastMergeJ)
	})
	return ms
}

// addMergesAndSpaceToSVGData adds an empty line after every normal flow line.
func addMergesAndSpaceToSVGData(flow *svg.Flow, ms []mergeData) *svg.Flow {
	src := flow.Shapes
	dst := make([][]interface{}, len(src)*2-1+len(ms))
	di := len(dst) - 1
	si := len(src) - 1
	mi := 0
	for mi < len(ms) || si >= 0 {
		for ms[mi].lastMergeI >= si {
			dst[di] = ms[mi].line
			mi++
			di--
		}
		dst[di] = src[si]
		si--
		di--
		if si >= 0 { // TODO: less empty lines
			dst[di] = []interface{}{}
			di--
		}
	}
	return &svg.Flow{Shapes: dst[di+1:]}
}

// FlowToSVG converts a flow DSL string into a SVG diagram string.
type FlowToSVG struct {
	pFlow *parser.ParseFlow
}

// NewFlowToSVG creates a new converter for a flow.
// If any regular expression used by the parser is invalid an error is
// returned.
func NewFlowToSVG() (*FlowToSVG, error) {
	pFlow, err := parser.NewParseFlow()
	if err != nil {
		return nil, err
	}
	return &FlowToSVG{pFlow: pFlow}, nil
}

// ConvertFlowToSVG converts a flow string into a SVG diagram and parser
// feedback.
// If the flow is invalid or some other error happens an error and no diagram
// or feedback is returned.
func (fts *FlowToSVG) ConvertFlowToSVG(flowContent, flowName string,
) ([]byte, string, error) {
	pd := gparselib.NewParseData(flowName, flowContent)
	pd, _ = fts.pFlow.In(pd, nil)

	fb, err := checkParserFeedback(pd)
	if err != nil {
		return nil, "", err
	}

	//svgData, decls, err := parserPartsToSVGData(pd.Result.Value.(data.Flow))
	_, _, err = parserPartsToSVGData(pd.Result.Value.(data.Flow))
	if err != nil {
		return nil, "", err
	}
	/*
		flow, merges := addSplitsAndMergesToSVGData(flow, declRefs)

		flow = addMergesAndSpaceToSVGData(flow, merges)
	*/
	buf, err := svg.FromFlowData(nil)
	if err != nil {
		return nil, "", err
	}

	return buf, fb, nil
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
