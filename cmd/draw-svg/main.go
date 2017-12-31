package main

import (
	"fmt"
	"log"
	"os"
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
{{- if .IsFill}}
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

type svgArrow struct {
	X1, Y1       int
	X2, Y2       int
	XTip1, YTip1 int
	XTip2, YTip2 int
}

type svgRect struct {
	X, Y   int
	Width  int
	Height int
	IsFill bool
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
}

type arrow struct {
	dataType string
	hasSrcOp bool
	srcPort  string
	hasDstOp bool
	dstPort  string
}

type rect struct {
	width  int
	height int
	isFill bool
	text   []string
}
type fill struct {
	title string
	rects []*rect
}

type op struct {
	main  *rect
	fills []*fill
}

type flow struct {
	shapes []interface{}
}

func flowDataToSVG(f *flow) *svgFlow {
	sas := make([]*svgArrow, 0, len(f.shapes))
	srs := make([]*svgRect, 0, len(f.shapes))
	sls := make([]*svgLine, 0, 64)
	sts := make([]*svgText, 0, 64)
	x0 := 2 // don't start directly at the edge
	y0 := 0 // the shapes leave head room themself
	var y, ymax int

	for _, is := range f.shapes {
		switch s := is.(type) {
		case *arrow:
			sas, sts, x0, y = arrowDataToSVG(s, sas, sts, x0, y0)
		case *op:
			srs, sls, sts, x0, y = opDataToSVG(s, srs, sls, sts, x0, y0)
		default:
			panic(fmt.Sprintf("unsupported type: %T", is))
		}
		ymax = max(ymax, y)
	}
	return &svgFlow{
		TotalWidth: x0 + 2, TotalHeight: ymax,
		Arrows: sas, Rects: srs, Lines: sls, Texts: sts,
	}
}
func arrowDataToSVG(
	a *arrow,
	sas []*svgArrow,
	sts []*svgText,
	x0, y0 int,
) ([]*svgArrow, []*svgText, int, int) {
	x := x0
	y := y0 + 24
	portLen := 0 // length in chars NOT pixels

	sts, x, portLen = addSrcPort(a, sts, x, y)

	if a.hasDstOp {
		portLen += len(a.dstPort)
	}
	width := max(
		portLen+1,
		len(a.dataType)+2,
	)*12 + 6 + // 6 so the source port text isn't glued to the op
		12 // last 12 is for tip of arrow

	if a.dataType != "" {
		sts = append(sts, &svgText{
			X: x + ((width-12)-len(a.dataType)*12)/2, Y: y - 8,
			Width: len(a.dataType) * 12,
			Text:  a.dataType,
		})
	}

	sas = append(sas, &svgArrow{
		X1: x, Y1: y,
		X2: x + width, Y2: y,
		XTip1: x + width - 8, YTip1: y - 8,
		XTip2: x + width - 8, YTip2: y + 8,
	})
	x += width

	sts, x = addDstPort(a, sts, x, y)

	return sas, sts, x, y + 36
}
func addSrcPort(a *arrow, sts []*svgText, x, y int) ([]*svgText, int, int) {
	portLen := 0
	if !a.hasSrcOp { // text before the arrow
		if a.srcPort != "" {
			sts = append(sts, &svgText{
				X: x + 1, Y: y + 6,
				Width: len(a.srcPort)*12 - 2,
				Text:  a.srcPort,
			})
		}
		x += 12 * len(a.srcPort)
	} else { // text under the arrow
		portLen += len(a.srcPort)
		if a.srcPort != "" {
			sts = append(sts, &svgText{
				X: x + 6, Y: y + 20,
				Width: len(a.srcPort) * 12,
				Text:  a.srcPort,
			})
		}
	}
	return sts, x, portLen
}
func addDstPort(a *arrow, sts []*svgText, x, y int) ([]*svgText, int) {
	if !a.hasDstOp {
		if a.dstPort != "" { // text after the arrow
			sts = append(sts, &svgText{
				X: x + 1, Y: y + 6,
				Width: len(a.dstPort)*12 - 2,
				Text:  a.dstPort,
			})
		}
		x += 12 * len(a.dstPort)
	} else if a.dstPort != "" { // text under the arrow
		sts = append(sts, &svgText{
			X: x - len(a.dstPort)*12 - 12, Y: y + 20,
			Width: len(a.dstPort) * 12,
			Text:  a.dstPort,
		})
	}
	return sts, x
}

func opDataToSVG(
	op *op,
	srs []*svgRect,
	sls []*svgLine,
	sts []*svgText,
	x0, y0 int,
) ([]*svgRect, []*svgLine, []*svgText, int, int) {
	var xn, yn int
	opW, opH := textDimensions(op.main)
	opW += 2 * 12
	opH += 6 + 10
	for _, f := range op.fills {
		w, l := fillDimensions(f)
		opH += l
		opW = max(opW, w)
	}
	srs, sts, y0, xn, yn = outerOpToSVG(op.main, opW, opH, srs, sts, x0, y0)
	for _, f := range op.fills {
		srs, sls, sts, y0 = fillDataToSVG(f, xn-x0, srs, sls, sts, x0, y0)
	}
	return srs, sls, sts, xn, yn
}
func textDimensions(r *rect) (width int, height int) {
	width = maxLen(r.text) * 12
	height += len(r.text) * 24
	return
}

func outerOpToSVG(
	r *rect,
	w int,
	h int,
	srs []*svgRect,
	sts []*svgText,
	x0, y0 int,
) (srs2 []*svgRect, sts2 []*svgText, y02 int, xn int, yn int) {
	x := x0
	y := y0 + 6
	h0 := len(r.text)*24 + 6*2

	srs = append(srs, &svgRect{
		X: x, Y: y,
		Width:  w,
		Height: h,
		IsFill: false,
	})

	y += 6
	for _, t := range r.text {
		sts = append(sts, &svgText{
			X: x + 12, Y: y + 24 - 6,
			Width: len(t) * 12,
			Text:  t,
		})
		y += 24
	}

	return srs, sts, y0 + 6 + h0, x + w, y0 + h + 2*6
}

func fillDataToSVG(
	f *fill,
	width int,
	srs []*svgRect,
	sls []*svgLine,
	sts []*svgText,
	x0, y0 int,
) (srs2 []*svgRect, sls2 []*svgLine, sts2 []*svgText, yn int) {
	x := x0
	y := y0

	y += 3
	sts = append(sts, &svgText{
		X: x + 6, Y: y + 24 - 6,
		Width: len(f.title) * 12,
		Text:  f.title,
	})
	y += 24

	for _, r := range f.rects {
		sls = append(sls, &svgLine{
			X1: x0, Y1: y,
			X2: x0 + width, Y2: y,
		})
		y += 3
		for _, t := range r.text {
			sts = append(sts, &svgText{
				X: x + 6, Y: y + 24 - 6,
				Width: len(t) * 12,
				Text:  t,
			})
			y += 24
		}
	}

	y += 3
	srs = append(srs, &svgRect{
		X: x0, Y: y0,
		Width:  width,
		Height: y - y0,
		IsFill: true,
	})

	return srs, sls, sts, y
}
func fillDimensions(f *fill) (width int, height int) {
	height = 24 + 2*3             // title text and padding
	width = len(f.title)*12 + 2*6 // title text and padding
	for _, r := range f.rects {
		w, h := textDimensions(r)
		height += h + 2*3
		width = max(width, w+2*6)
	}
	return
}

func main() {
	flow := &flow{
		shapes: []interface{}{
			&arrow{
				dataType: "",
				hasSrcOp: false, srcPort: "in",
				hasDstOp: true, dstPort: "",
			},
			&op{
				main: &rect{
					text: []string{"do"},
				},
				fills: []*fill{
					{
						title: "semantics:",
						rects: []*rect{
							{text: []string{"TextSemantics"}},
						},
					},
					{
						title: "subParser:",
						rects: []*rect{
							{text: []string{"LitralParser"}},
							{text: []string{"NaturalParser"}},
						},
					},
				},
			},
			&arrow{
				dataType: "Data",
				hasSrcOp: true, srcPort: "out",
				hasDstOp: true, dstPort: "in",
			},
			&op{
				main: &rect{
					text: []string{"ra", "(MiSo)"},
				},
			},
		},
	}
	svgflow := flowDataToSVG(flow)

	// compile and execute template
	t := template.Must(template.New("diagram").Parse(svgDiagram))
	err := t.Execute(os.Stdout, svgflow)
	if err != nil {
		log.Println("executing template:", err)
	}
}

func maxLen(ss []string) int {
	m := 0
	for _, s := range ss {
		m = max(m, len(s))
	}
	return m
}
func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
