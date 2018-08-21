package svg

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
