package svg

import (
	"bytes"
	"fmt"
	"text/template"
)

const svgDiagram = `<?xml version="1.0" ?>
<svg version="1.1" xmlns="http://www.w3.org/2000/svg" width="{{.TotalWidth}}px" height="{{.TotalHeight}}px">
<!-- Generated by simple FlowDev draw-svg tool. -->
	<rect fill="rgb(255,255,255)" fill-opacity="1" stroke="none" stroke-opacity="1" stroke-width="0.0" width="{{.TotalWidth}}" height="{{.TotalHeight}}" x="0" y="0"/>
{{- range .Arrows}}
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="{{.X1}}" y1="{{.Y1}}" x2="{{.X2}}" y2="{{.Y2}}"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="{{.XTip1}}" y1="{{.YTip1}}" x2="{{.X2}}" y2="{{.Y2}}"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="{{.XTip2}}" y1="{{.YTip2}}" x2="{{.X2}}" y2="{{.Y2}}"/>
{{end}}
{{- range .Rects}}
{{- if .IsPlugin}}
	<rect fill="rgb(32,224,32)" fill-opacity="1.0" stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" width="{{.Width}}" height="{{.Height}}" x="{{.X}}" y="{{.Y}}"/>
{{- else}}
	<rect fill="rgb(96,196,255)" fill-opacity="1.0" stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" width="{{.Width}}" height="{{.Height}}" x="{{.X}}" y="{{.Y}}" rx="10" ry="10"/>
{{- end}}
{{- end}}
{{range .Lines}}
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="1.0" x1="{{.X1}}" y1="{{.Y1}}" x2="{{.X2}}" y2="{{.Y2}}"/>
{{- end}}
{{range .Texts}}
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="{{.X}}" y="{{.Y}}" textLength="{{.Width}}" lengthAdjust="spacingAndGlyphs">{{.Text}}</text>
{{- end}}
</svg>
`

// Arrow contains all information for displaying an Arrow including data type
// and ports.
type Arrow struct {
	DataType string
	HasSrcOp bool
	SrcPort  string
	HasDstOp bool
	DstPort  string
}

// Rect just contains the text lines to display in a rectangle.
type Rect struct {
	Text []string
}

// Plugin is a helper operation that is used inside a proper operation.
type Plugin struct {
	Title string
	Rects []*Rect
}

// Op holds all data to describe a single operation including possible plugins.
type Op struct {
	Main    *Rect
	Plugins []*Plugin
}

// Split contains data for multiple paths/arrows originating from a single Op.
type Split struct {
	Shapes [][]interface{}
}

// Merge holds data for merging multiple paths/arrows into a single Op.
type Merge struct {
	ID   string
	Size int
}

// Flow contains data for a whole flow.
// The data is organized in rows and individual shapes per row.
// Valid shapes are Arrow, Op, Split and Merge.
type Flow struct {
	Shapes [][]interface{}
}

type svgArrow struct {
	X1, Y1       int
	X2, Y2       int
	XTip1, YTip1 int
	XTip2, YTip2 int
}

type svgRect struct {
	X, Y     int
	Width    int
	Height   int
	IsPlugin bool
}

type svgLine struct {
	X1, Y1 int
	X2, Y2 int
}

type svgText struct {
	X, Y  int
	Width int
	Text  string
}

type svgFlow struct {
	TotalWidth  int
	TotalHeight int
	Arrows      []*svgArrow
	Rects       []*svgRect
	Lines       []*svgLine
	Texts       []*svgText

	completedMerge *myMergeData
	allMerges      map[string]*myMergeData
}

type myMergeData struct {
	moveData []*moveData
	curSize  int
	x0, y0   int
	yn       int
}
type moveData struct {
	arrow       *svgArrow
	dataText    *svgText
	dstPortText *svgText
	yn          int
}

var tmpl = template.Must(template.New("diagram").Parse(svgDiagram))

// FromFlowData creates a SVG diagram from flow data.
// If the flow data isn't valid or the SVG diagram can't be created with its
// template, an error is returned.
func FromFlowData(f *Flow) ([]byte, error) {
	var err error
	f, err = validateFlowData(f)
	if err != nil {
		return nil, err
	}

	sf := flowDataToSVGFlow(f)

	return svgFlowToBytes(sf)
}

func validateFlowData(f *Flow) (*Flow, error) {
	if f == nil || len(f.Shapes) <= 0 {
		return nil, fmt.Errorf("flow is empty")
	}
	for i, row := range f.Shapes {
		for j, ishape := range row {
			switch ishape.(type) {
			case *Arrow, *Op, *Split, *Merge:
				break
			default:
				return nil, fmt.Errorf(
					"unsupported shape type %T at row index %d and column index %d",
					ishape, i, j)
			}
		}
	}
	return f, nil
}

func svgFlowToBytes(sf *svgFlow) ([]byte, error) {
	buf := bytes.Buffer{}
	err := tmpl.Execute(&buf, sf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func flowDataToSVGFlow(f *Flow) *svgFlow {
	sf, x, y := initSVGData()
	sf, x, y = shapesToSVG(
		f.Shapes,
		sf, x, y,
		arrowDataToSVG,
		opDataToSVG,
		rectDataToSVG,
		splitDataToSVG,
		mergeDataToSVG,
	)
	return adjustDimensions(sf, x, y)
}

func initSVGData() (sf *svgFlow, x0, y0 int) {
	return &svgFlow{
		Arrows: make([]*svgArrow, 0, 64),
		Rects:  make([]*svgRect, 0, 64),
		Lines:  make([]*svgLine, 0, 64),
		Texts:  make([]*svgText, 0, 64),

		allMerges: make(map[string]*myMergeData),
	}, 2, 1
}
func adjustDimensions(sf *svgFlow, xn, yn int) *svgFlow {
	sf.TotalWidth = xn + 2
	sf.TotalHeight = yn + 3
	return sf
}

func shapesToSVG(
	shapes [][]interface{}, sf *svgFlow, x0 int, y0 int,
	pluginArrowDataToSVG func(*Arrow, *svgFlow, int, int) (*svgFlow, int, int, *moveData),
	pluginOpDataToSVG func(*Op, *svgFlow, int, int) (*svgFlow, *svgRect, int, int, int),
	pluginRectDataToSVG func(*Rect, *svgFlow, int, int) (*svgFlow, int, int),
	pluginSplitDataToSVG func(*Split, *svgFlow, *svgRect, int, int) (*svgFlow, int, int),
	pluginMergeDataToSVG func(*Merge, *svgFlow, *moveData, int, int) *myMergeData,
) (nsf *svgFlow, xn, yn int) {
	var xmax, ymax int
	var mod *moveData
	var lsr *svgRect

	for _, ss := range shapes {
		x := x0
		lsr = nil
		if len(ss) < 1 {
			y0 += 48
			continue
		}
		for _, is := range ss {
			y := y0
			switch s := is.(type) {
			case *Arrow:
				sf, x, y, mod = pluginArrowDataToSVG(s, sf, x, y)
				lsr = nil
			case *Op:
				sf, lsr, y0, x, y = pluginOpDataToSVG(s, sf, x, y0)
				sf.completedMerge = nil
			case *Rect:
				sf, x, y = pluginRectDataToSVG(s, sf, x, y)
			case *Split:
				sf, x, y = pluginSplitDataToSVG(s, sf, lsr, x, y)
				lsr = nil
			case *Merge:
				sf.completedMerge = pluginMergeDataToSVG(s, sf, mod, x, y)
				mod = nil
			default:
				panic(fmt.Sprintf("unsupported type: %T", is))
			}

			ymax = max(ymax, y)
		}
		xmax = max(xmax, x)
		y0 = ymax + 5
	}
	return sf, xmax, ymax
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}