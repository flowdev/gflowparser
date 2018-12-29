package main

import "github.com/flowdev/gflowparser/svg"

var flowData = svg.Flow{
	Shapes: [][]interface{}{
		{
			&svg.Arrow{
				DataType: []string{"flowData"},
				HasSrcOp: false, SrcPort: "in",
				HasDstOp: true, DstPort: "",
			},
			&svg.Op{
				Main: &svg.Rect{
					Text: []string{"validateFlowData"},
				},
			},
			&svg.Split{
				Shapes: [][]interface{}{
					{
						&svg.Arrow{
							DataType: []string{"flowData"},
							HasSrcOp: true, SrcPort: "",
							HasDstOp: true, DstPort: "",
						},
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{"flowDataToSVGFlow"},
							},
						},
						&svg.Arrow{
							DataType: []string{"svgFlow"},
							HasSrcOp: true, SrcPort: "",
							HasDstOp: true, DstPort: "",
						},
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{"svgFlowToBytes"},
							},
						},
						&svg.Split{
							Shapes: [][]interface{}{
								{
									&svg.Arrow{
										DataType: []string{"bytes"},
										HasSrcOp: true, SrcPort: "",
										HasDstOp: false, DstPort: "out",
									},
								}, {
									&svg.Arrow{
										DataType: []string{"error"},
										HasSrcOp: true, SrcPort: "err",
										HasDstOp: false, DstPort: "err",
									},
								},
							},
						},
					}, {
						&svg.Arrow{
							DataType: []string{"error"},
							HasSrcOp: true, SrcPort: "err",
							HasDstOp: false, DstPort: "err",
						},
					},
				},
			},
		}, {}, {
			&svg.Arrow{
				DataType: []string{"flowData"},
				HasSrcOp: false, SrcPort: "in",
				HasDstOp: true, DstPort: "",
			},
			&svg.Op{
				Main: &svg.Rect{
					Text: []string{"initSVGData"},
				},
			},
			&svg.Arrow{
				DataType: []string{"(flowShapes, svgFlow, x0, y0)"},
				HasSrcOp: true, SrcPort: "",
				HasDstOp: true, DstPort: "",
			},
			&svg.Op{
				Main: &svg.Rect{
					Text: []string{"shapesToSVG"},
				},
				Plugins: []*svg.Plugin{
					{Title: "arrowDataToSVG"},
					{Title: "opDataToSVG"},
					{Title: "splitDataToSVG"},
					{Title: "mergeDataToSVG"},
				},
			},
			&svg.Arrow{
				DataType: []string{"(svgFlow, xn, yn)"},
				HasSrcOp: true, SrcPort: "",
				HasDstOp: true, DstPort: "",
			},
			&svg.Op{
				Main: &svg.Rect{
					Text: []string{"adjustDimensions"},
				},
			},
			&svg.Arrow{
				DataType: []string{"svgFlow"},
				HasSrcOp: true, SrcPort: "",
				HasDstOp: false, DstPort: "out",
			},
		},
	},
}
