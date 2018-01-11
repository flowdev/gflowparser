package main

var flowData = &flow{
	shapes: [][]interface{}{
		{
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
			&arrow{
				dataType: "",
				hasSrcOp: true, srcPort: "",
				hasDstOp: false, dstPort: "...",
			},
		}, {
			&arrow{
				dataType: "",
				hasSrcOp: false, srcPort: "...",
				hasDstOp: true, dstPort: "",
			},
			&op{
				main: &rect{
					text: []string{"megaParser", "(MegaParser)"},
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
		},
	},
}
