package svg

func opDataToSVG(op *Op, sf *svgFlow, x0, y0 int,
) (nsf *svgFlow, lsr *svgRect, ny0 int, xn, yn int) {
	var y int

	opW := maxTextWidth(op.Main) + 2*12 // text + padding
	opH := 0                            // outerOpToSVG should calculate itself
	for _, f := range op.Plugins {
		w := maxPluginWidth(f)
		opW = max(opW, w)
	}

	if sf.completedMerge != nil {
		x0 = sf.completedMerge.x0
		y0 = sf.completedMerge.y0
		ny0 = y0
		opH = sf.completedMerge.yn - y0 // now we have a minimum to enforce
	}

	lsr, y, xn, yn = outerOpToSVG(op.Main, opW, opH, sf, x0, y0)
	for _, f := range op.Plugins {
		y = pluginDataToSVG(f, xn-x0, sf, x0, y)
	}
	if len(op.Plugins) > 0 {
		y += 6
		lsr.Height = max(lsr.Height+6, y-y0)
		yn = max(yn, y0+lsr.Height+2*6)
	}

	return sf, lsr, y0, xn, yn
}

func outerOpToSVG(r *Rect, w int, h int, sf *svgFlow, x0, y0 int,
) (svgMainRect *svgRect, y02 int, xn int, yn int) {
	x := x0
	y := y0 + 6
	h0 := len(r.Text)*24 + 6*2
	h = max(h, h0)

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

func pluginDataToSVG(
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

	for i, r := range f.Rects {
		if i > 0 || f.Title != "" {
			sf.Lines = append(sf.Lines, &svgLine{
				X1: x0, Y1: y,
				X2: x0 + width, Y2: y,
			})
			y += 3
		}
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

func maxPluginWidth(f *Plugin) int {
	width := 0
	if f.Title != "" {
		width = (len(f.Title)+1)*12 + 2*6 // title text and padding
	}
	for _, r := range f.Rects {
		w := maxTextWidth(r)
		width = max(width, w+2*6)
	}
	return width
}

func maxTextWidth(r *Rect) int {
	return maxLen(r.Text) * 12
}

func maxLen(ss []string) int {
	m := 0
	for _, s := range ss {
		m = max(m, len(s))
	}
	return m
}
