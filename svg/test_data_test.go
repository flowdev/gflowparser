package svg

var BigTestFlowData = &Flow{
	Shapes: [][]interface{}{
		{
			&Arrow{
				DataType: "Data",
				HasSrcOp: false, SrcPort: "in",
				HasDstOp: true, DstPort: "",
			},
			&Op{
				Main: &Rect{
					Text: []string{"ra", "(MiSo)"},
				},
			},
			&Split{
				Shapes: [][]interface{}{
					{
						&Arrow{
							DataType: "Data",
							HasSrcOp: true, SrcPort: "special",
							HasDstOp: true, DstPort: "in",
						},
						&Op{
							Main: &Rect{
								Text: []string{"do"},
							},
							Fills: []*Fill{
								{
									Title: "semantics:",
									Rects: []*Rect{
										{Text: []string{"TextSemantics"}},
									},
								},
								{
									Title: "subParser:",
									Rects: []*Rect{
										{Text: []string{"LitralParser"}},
										{Text: []string{"NaturalParser"}},
									},
								},
							},
						},
						&Arrow{
							DataType: "BigDataType",
							HasSrcOp: true, SrcPort: "out",
							HasDstOp: true, DstPort: "in1",
						},
						&Merge{
							ID:   "BigMerge",
							Size: 3,
						},
					}, {
						&Arrow{
							DataType: "Data",
							HasSrcOp: true, SrcPort: "out",
							HasDstOp: true, DstPort: "in",
						},
						&Op{
							Main: &Rect{
								Text: []string{"bla", "(Blue)"},
							},
						},
						&Arrow{
							DataType: "Data2",
							HasSrcOp: true, SrcPort: "",
							HasDstOp: false, DstPort: "...",
						},
					},
				},
			},
		}, {
			&Split{
				Shapes: [][]interface{}{
					{
						&Arrow{
							DataType: "",
							HasSrcOp: false, SrcPort: "...",
							HasDstOp: true, DstPort: "in",
						},
						&Op{
							Main: &Rect{
								Text: []string{"bla", "(Blue)"},
							},
						},
						&Arrow{
							DataType: "Data",
							HasSrcOp: true, SrcPort: "out",
							HasDstOp: true, DstPort: "in2",
						},
						&Merge{
							ID:   "BigMerge",
							Size: 3,
						},
					}, {
						&Arrow{
							DataType: "Data3",
							HasSrcOp: false, SrcPort: "in2",
							HasDstOp: true, DstPort: "in",
						},
						&Op{
							Main: &Rect{
								Text: []string{"megaParser", "(MegaParser)"},
							},
							Fills: []*Fill{
								{
									Title: "semantics:",
									Rects: []*Rect{
										{Text: []string{"TextSemantics"}},
									},
								},
								{
									Title: "subParser:",
									Rects: []*Rect{
										{Text: []string{"LitralParser"}},
										{Text: []string{"NaturalParser"}},
									},
								},
							},
						},
						&Arrow{
							DataType: "Data",
							HasSrcOp: true, SrcPort: "out",
							HasDstOp: true, DstPort: "in3",
						},
						&Merge{
							ID:   "BigMerge",
							Size: 3,
						},
					},
				},
			},
		}, {
			&Op{
				Main: &Rect{
					Text: []string{"BigMerge"},
				},
			},
		},
	},
}
