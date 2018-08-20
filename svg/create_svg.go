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
	}, 2, 1
}
func adjustDimensions(sf *svgFlow, xn, yn int) *svgFlow {
	sf.TotalWidth = xn + 2
	sf.TotalHeight = yn + 3
	return sf
}

// Unfortunately this has to be global as the next op can be in another
// shapesToSVG call than the last merge.
// ATTENTION: The op has to come directly after the last merge!
var completedMerge *myMergeData

func shapesToSVG(
	shapes [][]interface{}, sf *svgFlow, x0 int, y0 int,
	pluginArrowDataToSVG func(*Arrow, *svgFlow, int, int) (*svgFlow, int, int, *moveData),
	pluginOpDataToSVG func(*Op, *svgFlow, *myMergeData, int, int) (*svgFlow, *svgRect, int, int, int),
	pluginRectDataToSVG func(*Rect, *svgFlow, int, int) (*svgFlow, int, int),
	pluginSplitDataToSVG func(*Split, *svgFlow, *svgRect, int, int) (*svgFlow, int, int),
	pluginMergeDataToSVG func(*Merge, *moveData, int, int) *myMergeData,
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
				sf, lsr, y0, x, y = pluginOpDataToSVG(s, sf, completedMerge, x, y0)
				completedMerge = nil
			case *Rect:
				sf, x, y = pluginRectDataToSVG(s, sf, x, y)
			case *Split:
				sf, x, y = pluginSplitDataToSVG(s, sf, lsr, x, y)
				lsr = nil
			case *Merge:
				completedMerge = pluginMergeDataToSVG(s, mod, x, y)
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

// Unfortunately this has to be global as merges can be in different
// shapesToSVG calls.
// ATTENTION: The op has to come directly after the last merge!
var allMerges = make(map[string]*myMergeData)

func mergeDataToSVG(m *Merge, mod *moveData, x0, y0 int,
) (completedMerge *myMergeData) {
	md := allMerges[m.ID]
	if md == nil { // first merge
		md = &myMergeData{
			x0:       x0,
			y0:       y0,
			yn:       mod.yn,
			curSize:  1,
			moveData: []*moveData{mod},
		}
		allMerges[m.ID] = md
	} else { // additional merge
		md.x0 = max(md.x0, x0)
		md.y0 = min(md.y0, y0)
		md.yn = max(md.yn, mod.yn)
		md.curSize++
		md.moveData = append(md.moveData, mod)
	}
	if md.curSize >= m.Size { // merge is comleted!
		moveXTo(md, md.x0)
		return md
	}
	return nil
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

func splitDataToSVG(s *Split, sf *svgFlow, lsr *svgRect, x0, y0 int,
) (nsf *svgFlow, xn, yn int) {
	nsf, xn, yn = shapesToSVG(
		s.Shapes,
		sf, x0, y0,
		arrowDataToSVG,
		opDataToSVG,
		rectDataToSVG,
		splitDataToSVG,
		mergeDataToSVG,
	)
	adjustLastRect(lsr, yn)
	return
}

func adjustLastRect(lsr *svgRect, yn int) {
	if lsr != nil {
		if lsr.Y+lsr.Height < yn {
			lsr.Height = yn - lsr.Y
		}
	}
}

func arrowDataToSVG(a *Arrow, sf *svgFlow, x int, y int,
) (nsf *svgFlow, nx, ny int, mod *moveData) {
	var dstPortText, dataText *svgText
	y += 24
	portLen := 0 // length in chars NOT pixels

	sf.Texts, x, portLen = addSrcPort(a, sf.Texts, x, y)

	if a.HasDstOp {
		portLen += len(a.DstPort)
	}
	width := max(
		portLen+2,
		len(a.DataType)+2,
	)*12 + 6 + // 6 so the source port text isn't glued to the op
		12 // last 12 is for tip of arrow

	if a.DataType != "" {
		dataText = &svgText{
			X: x + ((width-12)-len(a.DataType)*12)/2, Y: y - 8,
			Width: len(a.DataType) * 12,
			Text:  a.DataType,
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
	if a.DstPort != "" {
		dstPortText = sf.Texts[len(sf.Texts)-1]
	}

	return sf, x, y + 36, &moveData{
		arrow:       sf.Arrows[len(sf.Arrows)-1],
		dstPortText: dstPortText,
		dataText:    dataText,
		yn:          y + 24,
	}
}
func addSrcPort(a *Arrow, sts []*svgText, x, y int) ([]*svgText, int, int) {
	portLen := 0
	if !a.HasSrcOp { // text before the arrow
		if a.SrcPort != "" {
			sts = append(sts, &svgText{
				X: x + 1, Y: y + 6,
				Width: len(a.SrcPort)*12 - 2,
				Text:  a.SrcPort,
			})
		}
		x += 12 * len(a.SrcPort)
	} else { // text under the arrow
		portLen += len(a.SrcPort)
		if a.SrcPort != "" {
			sts = append(sts, &svgText{
				X: x + 6, Y: y + 20,
				Width: len(a.SrcPort) * 12,
				Text:  a.SrcPort,
			})
		}
	}
	return sts, x, portLen
}
func addDstPort(a *Arrow, sts []*svgText, x, y int) ([]*svgText, int) {
	if !a.HasDstOp {
		if a.DstPort != "" { // text after the arrow
			sts = append(sts, &svgText{
				X: x + 3, Y: y + 6,
				Width: len(a.DstPort)*12 - 2,
				Text:  a.DstPort,
			})
		}
		x += 3 + 12*len(a.DstPort)
	} else if a.DstPort != "" { // text under the arrow
		sts = append(sts, &svgText{
			X: x - len(a.DstPort)*12 - 12, Y: y + 20,
			Width: len(a.DstPort) * 12,
			Text:  a.DstPort,
		})
	}
	return sts, x
}

func rectDataToSVG(r *Rect, sf *svgFlow, x int, y int) (nsf *svgFlow, nx, ny int) {
	txt := "... back to: " + r.Text[0]
	width := len(txt) * 12

	sf.Texts = append(sf.Texts, &svgText{
		X: x, Y: y + 6,
		Width: width,
		Text:  txt,
	})

	x += width

	return sf, x + width + 12, y + 12
}

func opDataToSVG(op *Op, sf *svgFlow, completedMerge *myMergeData, x0, y0 int,
) (nsf *svgFlow, lsr *svgRect, ny0 int, xn, yn int) {
	var y int

	opW, opH := textDimensions(op.Main)
	opW += 2 * 12
	opH += 6 + 10
	for _, f := range op.Plugins {
		w, l := fillDimensions(f)
		opH += l
		opW = max(opW, w)
	}
	if len(op.Plugins) > 0 {
		opH += 6
	}

	if completedMerge != nil {
		x0 = completedMerge.x0
		y0 = completedMerge.y0
		ny0 = y0
		opH = max(opH, completedMerge.yn-y0)
	}

	lsr, y, xn, yn = outerOpToSVG(op.Main, opW, opH, sf, x0, y0)
	for _, f := range op.Plugins {
		y = fillDataToSVG(f, xn-x0, sf, x0, y)
	}

	return sf, lsr, y0, xn, yn
}
func textDimensions(r *Rect) (width int, height int) {
	width = maxLen(r.Text) * 12
	height += len(r.Text) * 24
	return
}
func outerOpToSVG(r *Rect, w int, h int, sf *svgFlow, x0, y0 int,
) (svgMainRect *svgRect, y02 int, xn int, yn int) {
	x := x0
	y := y0 + 6
	h0 := len(r.Text)*24 + 6*2

	svgMainRect = &svgRect{
		X: x, Y: y,
		Width: w, Height: h,
		IsPlugin: false,
	}
	sf.Rects = append(sf.Rects, svgMainRect)

	y += 6
	for _, t := range r.Text {
		sf.Texts = append(sf.Texts, &svgText{
			X: x + 12, Y: y + 24 - 6,
			Width: len(t) * 12,
			Text:  t,
		})
		y += 24
	}

	return svgMainRect, y0 + 6 + h0, x + w, y0 + h + 2*6
}
func fillDataToSVG(
	f *Plugin,
	width int,
	sf *svgFlow,
	x0, y0 int,
) (yn int) {
	x := x0
	y := y0

	y += 3
	if f.Title != "" {
		sf.Texts = append(sf.Texts, &svgText{
			X: x + 6, Y: y + 24 - 6,
			Width: (len(f.Title) + 1) * 12,
			Text:  f.Title + ":",
		})
		y += 24
	}

	for _, r := range f.Rects {
		sf.Lines = append(sf.Lines, &svgLine{
			X1: x0, Y1: y,
			X2: x0 + width, Y2: y,
		})
		y += 3
		for _, t := range r.Text {
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
		Width:    width,
		Height:   y - y0,
		IsPlugin: true,
	})

	return y
}
func fillDimensions(f *Plugin) (width int, height int) {
	if f.Title != "" {
		height = 24                       // title text
		width = (len(f.Title)+1)*12 + 2*6 // title text and padding
	}
	height += 2 * 3 // padding
	for _, r := range f.Rects {
		w, h := textDimensions(r)
		height += h + 3
		width = max(width, w+2*6)
	}
	return width, height
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
