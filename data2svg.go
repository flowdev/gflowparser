package gflowparser

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/flowdev/gflowparser/data"
	"github.com/flowdev/gflowparser/parser"
	"github.com/flowdev/gflowparser/svg"
	"github.com/flowdev/gparselib"
)

// Error messages.
const (
	errMsgDeclAndRef = "Circular flows aren't allowed yet, but the component " +
		"'%s' is declared here:\n%s\n... and referenced again here:\n%s"
	errMsg2Decls = "A component with the name '%s' is declared two times, " +
		"here:\n%s\n... and again here:\n%s"
	errMsgPartType = "Found illegal flow part type '%T' at index [%d, %d]"
	errMsgLoneComp = "Component reference with name '%s' without " +
		"input or output found:\n%s"
)

type whereer interface {
	Where(pos int) string
}

type decl struct {
	name     string
	srcPos   int
	i, j     int
	svgOp    *svg.Op
	svgMerge *svg.Merge
	svgSplit *svg.Split
}

type merge struct {
	name   string
	srcPos int
	svg    *svg.Merge
}

type split struct {
	name   string
	srcPos int
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
// 3. Into a split if it is only used for a split.
func parserPartsToSVGData(flowDat data.Flow, w whereer,
) (shapes [][]interface{}, decls map[string]*decl, clsts clusters, err error) {
	svgDat := make([][]interface{}, len(flowDat.Parts))
	decls = make(map[string]*decl)
	clsts = clusters(nil)

	for i, partLine := range flowDat.Parts {
		m := len(partLine) - 1
		svgLine := make([]interface{}, m+1)
		for j, part := range partLine {
			switch p := part.(type) {
			case data.Arrow:
				svgLine[j] = arrowToSVGData(p, j > 0, j < m)
			case data.Component:
				if dcl, ok := decls[p.Decl.Name]; ok {
					if !p.Decl.VagueType { // prevent double declaration
						return nil, nil, nil, fmt.Errorf(errMsg2Decls,
							dcl.name, w.Where(dcl.srcPos), w.Where(p.SrcPos))
					}
					if j > 0 { // we need a merge
						dcl.svgMerge.Size++
						clsts.addCluster(dcl.i, i)
						svgLine[j] = &merge{
							name:   dcl.name,
							srcPos: p.SrcPos,
							svg:    dcl.svgMerge,
						}
					} else if j < m { // we only need a split
						svgLine[j] = &split{name: dcl.name, srcPos: p.SrcPos}
					} else { // we don't need anything at all???!!!
						return nil, nil, nil, fmt.Errorf(errMsgLoneComp,
							dcl.name, w.Where(p.SrcPos))
					}
				} else {
					dcl := &decl{
						name:     p.Decl.Name,
						srcPos:   p.SrcPos,
						i:        i,
						j:        j,
						svgOp:    compToSVGData(p),
						svgMerge: &svg.Merge{ID: p.Decl.Name},
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
	return svgDat, decls, clsts, nil
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

// reshapeSVGData handles merges and splits
func reshapeSVGData(shapes [][]interface{}, decls map[string]*decl, clsts clusters,
) (nshapes [][]interface{}, ndecls map[string]*decl, nclsts clusters, err error) {
	return shapes, decls, clsts, nil
}

// breakCircles replaces back pointing merges with simple svg.Rects.
func breakCircles(shapes [][]interface{}, decls map[string]*decl, clsts clusters,
) (nshapes [][]interface{}, err error) {
	return shapes, nil
}

// cleanSVGData replaces all special data structures with pure SVG ones.
func cleanSVGData(shapes [][]interface{}) [][]interface{} {
	return shapes
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

	shapes, decls, clsts, err := parserPartsToSVGData(
		pd.Result.Value.(data.Flow),
		pd.Source,
	)
	if err != nil {
		return nil, "", err
	}

	// TODO: Restructure: insert merge comp -> addLine
	// TODO: Restructure: move split comp -> deleteLine
	shapes, decls, clsts, err = reshapeSVGData(shapes, decls, clsts)
	if err != nil {
		return nil, "", err
	}

	shapes, err = breakCircles(shapes, decls, clsts)
	if err != nil {
		return nil, "", err
	}

	shapes = cleanSVGData(shapes)

	buf, err := svg.FromFlowData(&svg.Flow{Shapes: shapes})
	if err != nil {
		return nil, "", err
	}

	return buf, fb, nil
}
