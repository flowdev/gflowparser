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

type split struct {
	shapes [][]interface{}
}

type merge struct {
	id   string
	size int
}

type flow struct {
	shapes [][]interface{}
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

var (
	allMerges      = make(map[string]*myMergeData)
	completedMerge *myMergeData
	adts           func(*arrow, *svgFlow, *int, *int, *moveData)
)

func flowDataToSVG(f *flow) *svgFlow {
	sf := &svgFlow{
		Arrows: make([]*svgArrow, 0, len(f.shapes)),
		Rects:  make([]*svgRect, 0, len(f.shapes)),
		Lines:  make([]*svgLine, 0, 64),
		Texts:  make([]*svgText, 0, 64),
	}
	x, y := shapesToSVG(f.shapes, sf, 2, 0)

	sf.TotalWidth = x + 2
	sf.TotalHeight = y + 3
	return sf
}
func shapesToSVG(shapes [][]interface{}, sf *svgFlow, x0, y0 int) (int, int) {
	var xmax, ymax int
	var mod *moveData
	var lsr *svgRect

	for _, ss := range shapes {
		x := x0
		lsr = nil
		for _, is := range ss {
			y := y0
			switch s := is.(type) {
			case *arrow:
				mod = &moveData{}
				adts(s, sf, &x, &y, mod)
				lsr = nil
			case *op:
				y0, x, y, lsr = opDataToSVG(s, sf, x, y0)
				mod = nil
			case *split:
				x, y = splitDataToSVG(s, sf, lsr, x, y0)
				mod = nil
				lsr = nil
			case *merge:
				mergeDataToSVG(s, mod, x, y0)
				mod = nil
				lsr = nil
			default:
				panic(fmt.Sprintf("unsupported type: %T", is))
			}
			ymax = max(ymax, y)
		}
		xmax = max(xmax, x)
		y0 = ymax + 5
	}
	return xmax, ymax
}
func mergeDataToSVG(m *merge, mod *moveData, x0, y0 int) {
	md := allMerges[m.id]
	if md == nil { // first merge
		md = &myMergeData{
			x0:       x0,
			y0:       y0,
			yn:       mod.yn,
			curSize:  1,
			moveData: []*moveData{mod},
		}
		allMerges[m.id] = md
	} else { // additional merge
		md.x0 = max(md.x0, x0)
		md.y0 = min(md.y0, y0)
		md.yn = max(md.yn, mod.yn)
		md.curSize++
		md.moveData = append(md.moveData, mod)
	}
	if md.curSize >= m.size { // merge is comleted!
		moveXTo(md, md.x0)
		completedMerge = md
	}
}
func moveXTo(med *myMergeData, newX int) {
	for _, mod := range med.moveData {
		xShift := newX - mod.arrow.X2

		mod.arrow.X2 = newX
		mod.arrow.XTip1 = newX - 8
		mod.arrow.XTip2 = newX - 8

		if mod.dstPortText != nil {
			mod.dstPortText.X += xShift
		}
		if mod.dataText != nil {
			mod.dataText.X += xShift / 2
		}
	}
}

func splitDataToSVG(s *split, sf *svgFlow, lsr *svgRect, x0, y0 int) (int, int) {
	x0, y0 = shapesToSVG(s.shapes, sf, x0, y0)
	if lsr != nil {
		if lsr.Y+lsr.Height < y0 {
			lsr.Height = y0 - lsr.Y
		}
	}
	return x0, y0
}

func arrowDataToSVG() (portIn func(*arrow, *svgFlow, *int, *int, *moveData)) {
	portIn = func(a *arrow, sf *svgFlow, px *int, py *int, mod *moveData) {
		var dstPortText, dataText *svgText
		x := *px
		y := *py + 24
		portLen := 0 // length in chars NOT pixels

		sf.Texts, x, portLen = addSrcPort(a, sf.Texts, x, y)

		if a.hasDstOp {
			portLen += len(a.dstPort)
		}
		width := max(
			portLen+2,
			len(a.dataType)+2,
		)*12 + 6 + // 6 so the source port text isn't glued to the op
			12 // last 12 is for tip of arrow

		if a.dataType != "" {
			dataText = &svgText{
				X: x + ((width-12)-len(a.dataType)*12)/2, Y: y - 8,
				Width: len(a.dataType) * 12,
				Text:  a.dataType,
			}
			sf.Texts = append(sf.Texts, dataText)
		}

		sf.Arrows = append(sf.Arrows, &svgArrow{
			X1: x, Y1: y,
			X2: x + width, Y2: y,
			XTip1: x + width - 8, YTip1: y - 8,
			XTip2: x + width - 8, YTip2: y + 8,
		})
		x += width

		sf.Texts, x = addDstPort(a, sf.Texts, x, y)
		if a.dstPort != "" {
			dstPortText = sf.Texts[len(sf.Texts)-1]
		}

		*px = x
		*py = y + 36
		mod.arrow = sf.Arrows[len(sf.Arrows)-1]
		mod.dstPortText = dstPortText
		mod.dataText = dataText
		mod.yn = y + 24
	}
	return
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

func opDataToSVG(op *op, sf *svgFlow, x0, y0 int) (int, int, int, *svgRect) {
	var y, xn, yn int
	opW, opH := textDimensions(op.main)
	opW += 2 * 12
	opH += 6 + 10
	for _, f := range op.fills {
		w, l := fillDimensions(f)
		opH += l
		opW = max(opW, w)
	}
	if completedMerge != nil {
		x0 = completedMerge.x0
		y0 = completedMerge.y0
		opH = max(opH, completedMerge.yn-y0)
		completedMerge = nil
	}
	y, xn, yn = outerOpToSVG(op.main, opW, opH, sf, x0, y0)
	lsr := sf.Rects[len(sf.Rects)-1]
	for _, f := range op.fills {
		y = fillDataToSVG(f, xn-x0, sf, x0, y)
	}
	return y0, xn, yn, lsr
}
func textDimensions(r *rect) (width int, height int) {
	width = maxLen(r.text) * 12
	height += len(r.text) * 24
	return
}

func outerOpToSVG(r *rect, w int, h int, sf *svgFlow, x0, y0 int,
) (y02 int, xn int, yn int) {
	x := x0
	y := y0 + 6
	h0 := len(r.text)*24 + 6*2

	sf.Rects = append(sf.Rects, &svgRect{
		X: x, Y: y,
		Width:  w,
		Height: h,
		IsFill: false,
	})

	y += 6
	for _, t := range r.text {
		sf.Texts = append(sf.Texts, &svgText{
			X: x + 12, Y: y + 24 - 6,
			Width: len(t) * 12,
			Text:  t,
		})
		y += 24
	}

	return y0 + 6 + h0, x + w, y0 + h + 2*6
}

func fillDataToSVG(
	f *fill,
	width int,
	sf *svgFlow,
	x0, y0 int,
) (yn int) {
	x := x0
	y := y0

	y += 3
	sf.Texts = append(sf.Texts, &svgText{
		X: x + 6, Y: y + 24 - 6,
		Width: len(f.title) * 12,
		Text:  f.title,
	})
	y += 24

	for _, r := range f.rects {
		sf.Lines = append(sf.Lines, &svgLine{
			X1: x0, Y1: y,
			X2: x0 + width, Y2: y,
		})
		y += 3
		for _, t := range r.text {
			sf.Texts = append(sf.Texts, &svgText{
				X: x + 6, Y: y + 24 - 6,
				Width: len(t) * 12,
				Text:  t,
			})
			y += 24
		}
	}

	y += 3
	sf.Rects = append(sf.Rects, &svgRect{
		X: x0, Y: y0,
		Width:  width,
		Height: y - y0,
		IsFill: true,
	})

	return y
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
	adts = arrowDataToSVG()
	svgflow := flowDataToSVG(flowData)

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
