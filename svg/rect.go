package svg

func rectDataToSVG(r *Rect, sf *svgFlow, x int, y int) (nsf *svgFlow, nx, ny int) {
	txt := "... back to: " + r.Text[0]
	width := len(txt) * 12

	y += 12 + 24 - 6
	sf.Texts = append(sf.Texts, &svgText{
		X: x, Y: y,
		Width: width,
		Text:  txt,
	})

	x += width

	return sf, x + width + 12, y + 12
}
