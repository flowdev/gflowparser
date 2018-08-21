package svg

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
