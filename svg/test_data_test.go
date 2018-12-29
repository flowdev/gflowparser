package svg

var BigTestFlowData = Flow{
	Shapes: [][]interface{}{
		{
			&Arrow{
				DataType: []string{"Data"},
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
							DataType: []string{"Data"},
							HasSrcOp: true, SrcPort: "special",
							HasDstOp: true, DstPort: "in",
						},
						&Op{
							Main: &Rect{
								Text: []string{"do"},
							},
							Plugins: []*Plugin{
								{
									Title: "semantics",
									Rects: []*Rect{
										{Text: []string{"TextSemantics"}},
									},
								},
								{
									Title: "subParser",
									Rects: []*Rect{
										{Text: []string{"LitralParser"}},
										{Text: []string{"NaturalParser"}},
									},
								},
							},
						},
						&Arrow{
							DataType: []string{"BigDataType"},
							HasSrcOp: true, SrcPort: "out",
							HasDstOp: true, DstPort: "in1",
						},
						&Merge{
							ID:   "BigMerge",
							Size: 3,
						},
					}, {
						&Arrow{
							DataType: []string{"Data"},
							HasSrcOp: true, SrcPort: "out",
							HasDstOp: true, DstPort: "in",
						},
						&Op{
							Main: &Rect{
								Text: []string{"bla", "(Blue)"},
							},
						},
						&Arrow{
							DataType: []string{"Data2"},
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
							DataType: []string{},
							HasSrcOp: false, SrcPort: "...",
							HasDstOp: true, DstPort: "in",
						},
						&Op{
							Main: &Rect{
								Text: []string{"bla", "(Blue)"},
							},
						},
						&Arrow{
							DataType: []string{"Data"},
							HasSrcOp: true, SrcPort: "out",
							HasDstOp: true, DstPort: "in2",
						},
						&Merge{
							ID:   "BigMerge",
							Size: 3,
						},
					}, {
						&Arrow{
							DataType: []string{"Data3"},
							HasSrcOp: false, SrcPort: "in2",
							HasDstOp: true, DstPort: "in",
						},
						&Op{
							Main: &Rect{
								Text: []string{"megaParser", "(MegaParser)"},
							},
							Plugins: []*Plugin{
								{
									Title: "semantics",
									Rects: []*Rect{
										{Text: []string{"TextSemantics"}},
									},
								},
								{
									Title: "subParser",
									Rects: []*Rect{
										{Text: []string{"LitralParser"}},
										{Text: []string{"NaturalParser"}},
									},
								},
							},
						},
						&Arrow{
							DataType: []string{"(Data,", " data2,", " Data3)"},
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
		}, { // empty to force more space
		}, {
			&Arrow{
				DataType: []string{"(Data,", " data2,", " Data3)"},
				HasSrcOp: false, SrcPort: "in3",
				HasDstOp: true, DstPort: "",
			},
			&Op{
				Main: &Rect{
					Text: []string{"recursive"},
				},
			},
			&Arrow{
				DataType: []string{"(Data)"},
				HasSrcOp: true, SrcPort: "",
				HasDstOp: true, DstPort: "",
			},
			&Op{
				Main: &Rect{
					Text: []string{"secondOp"},
				},
			},
			&Arrow{
				DataType: []string{"(Data,", " data2,", " Data3)"},
				HasSrcOp: true, SrcPort: "out",
				HasDstOp: true, DstPort: "",
			},
			&Rect{
				Text: []string{"recursive"},
			},
		},
	},
}
