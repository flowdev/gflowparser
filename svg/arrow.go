package svg

func arrowDataToSVG(a *Arrow, sf *svgFlow, lsr *svgRect, x int, y int,
) (nsf *svgFlow, nx, ny int, mod *moveData) {
	var srcPortText, dstPortText *svgText
	dataTexts := make([]*svgText, 0, 8)

	y += 24
	portLen := 0 // length in chars NOT pixels
	if a.HasSrcOp {
		portLen = len(a.SrcPort)
	}
	if a.HasDstOp {
		portLen += len(a.DstPort)
	}

	dataLen := maxLen(a.DataType)
	width := max(
		portLen+2,
		dataLen+2,
	)*12 + 6 + // 6 so the source port text isn't glued to the op
		12 // last 12 is for tip of arrow

	sf.Texts, x = addSrcPort(a, sf.Texts, x, y)
	if a.SrcPort != "" { // remember this text as we might have to move it down
		srcPortText = sf.Texts[len(sf.Texts)-1]
	}

	if len(a.DataType) != 0 {
		dataX := x + ((width-12)-dataLen*12)/2
		for i, text := range a.DataType {
			if i > 0 {
				y += 22
				if srcPortText != nil {
					srcPortText.Y += 22
				}
			}
			st := &svgText{
				X: dataX, Y: y - 8,
				Width: len(text) * 12,
				Text:  text,
			}
			sf.Texts = append(sf.Texts, st)
			dataTexts = append(dataTexts, st)
		}
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

	yn := y + 24
	adjustLastRect(lsr, yn-12)

	return sf, x, yn, &moveData{
		arrow:       sf.Arrows[len(sf.Arrows)-1],
		dstPortText: dstPortText,
		dataTexts:   dataTexts,
		yn:          yn,
	}
}

func addSrcPort(a *Arrow, sts []*svgText, x, y int) ([]*svgText, int) {
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
		if a.SrcPort != "" {
			sts = append(sts, &svgText{
				X: x + 6, Y: y + 20,
				Width: len(a.SrcPort) * 12,
				Text:  a.SrcPort,
			})
		}
	}
	return sts, x
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
