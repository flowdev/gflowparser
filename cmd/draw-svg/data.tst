package main

var flowData = &flow{
	shapes: [][]interface{}{
		{
			&arrow{
				dataType: "Data",
				hasSrcOp: false, srcPort: "in",
				hasDstOp: true, dstPort: "",
			},
			&op{
				main: &rect{
					text: []string{"ra", "(MiSo)"},
				},
			},
			&split{
				shapes: [][]interface{}{
					{
						&arrow{
							dataType: "Data",
							hasSrcOp: true, srcPort: "special",
							hasDstOp: true, dstPort: "in",
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
					}, {
						&arrow{
							dataType: "Data",
							hasSrcOp: true, srcPort: "out",
							hasDstOp: true, dstPort: "in",
						},
						&op{
							main: &rect{
								text: []string{"bla", "(Blue)"},
							},
						},
						&arrow{
							dataType: "Data2",
							hasSrcOp: true, srcPort: "",
							hasDstOp: false, dstPort: "...",
						},
					},
				},
			},
		}, {
			&split{
				shapes: [][]interface{}{
					{
						&arrow{
							dataType: "",
							hasSrcOp: false, srcPort: "...",
							hasDstOp: true, dstPort: "in",
						},
						&op{
							main: &rect{
								text: []string{"bla", "(Blue)"},
							},
						},
					}, {
						&arrow{
							dataType: "Data3",
							hasSrcOp: false, srcPort: "in2",
							hasDstOp: true, dstPort: "in",
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
			},
		},
	},
}
