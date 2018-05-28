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

// parserToSVGData converts a data.Flow into SVG diagram data.
func parserToSVGData(flowDat data.Flow) *svg.Flow {
	svgDat := make([][]interface{}, 0, len(flowDat.Parts))
	for i, partLine := range flowDat.Parts {
		m := len(partLine) - 1
		svgLine := make([]interface{}, 0, m+1)
		for j, part := range partLine {
			switch p := part.(type) {
			case data.Arrow:
				svgLine = append(svgLine, arrowToSVGData(p, j > 0, j < m))
			case data.Component:
				svgLine = append(svgLine, compToSVGData(p))
			default:
				panic(fmt.Sprintf("Found illegal flow part type '%T' at position: %d, %d", part, i, j))
			}
		}
		svgDat = append(svgDat, svgLine)
	}
	return &svg.Flow{
		Shapes: svgDat,
	}
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
	return []string{decl.Name, typeToSVGData(decl.Type)}
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

// checkParserFeedback converts parser errors into a single error.
func checkParserFeedback(pd *gparselib.ParseData) (string, error) {
	if pd.Result.HasError() {
		return "", errors.New(feedbackToString(pd))
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

// FlowToSVG converts a flow DSL string into a SVG diagram string.
//
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
func (fts *FlowToSVG) ConvertFlowToSVG(flowContent, flowName string) ([]byte, string, error) {
	pd := gparselib.NewParseData(flowName, flowContent)
	pd, _ = fts.pFlow.In(pd, nil)

	fb, err := checkParserFeedback(pd)
	if err != nil {
		return nil, "", err
	}

	flow := parserToSVGData(pd.Result.Value.(data.Flow))

	buf, err := svg.FromFlowData(flow)
	if err != nil {
		return nil, "", err
	}

	return buf, fb, nil
}
