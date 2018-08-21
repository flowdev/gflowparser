package svg

func mergeDataToSVG(m *Merge, sf *svgFlow, mod *moveData, x0, y0 int,
) (completedMerge *myMergeData) {
	md := sf.allMerges[m.ID]
	if md == nil { // first merge
		md = &myMergeData{
			x0:       x0,
			y0:       y0,
			yn:       mod.yn,
			curSize:  1,
			moveData: []*moveData{mod},
		}
		sf.allMerges[m.ID] = md
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
