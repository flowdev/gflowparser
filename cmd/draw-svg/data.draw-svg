package main

import "github.com/flowdev/gflowparser/svg"

var flowData = &svg.Flow{
	Shapes: [][]interface{}{
		{
			&svg.Arrow{
				DataType: "flowData",
				HasSrcOp: false, SrcPort: "in",
				HasDstOp: true, DstPort: "in",
			},
			&svg.Op{
				Main: &svg.Rect{
					Text: []string{"flowDataToSVGOp", "[synchronous]"},
				},
			},
			&svg.Split{
				Shapes: [][]interface{}{
					{
						&svg.Arrow{
							DataType: "(shapes, *svgFlow, *x, *y)",
							HasSrcOp: true, SrcPort: "shapesOut",
							HasDstOp: true, DstPort: "in",
						},
						&svg.Merge{
							ID:   "shapesToSVG",
							Size: 2,
						},
					}, {
						&svg.Arrow{
							DataType: "svgFlow",
							HasSrcOp: true, SrcPort: "out",
							HasDstOp: false, DstPort: "out",
						},
					},
				},
			},
		}, {
			&svg.Arrow{
				HasSrcOp: false, SrcPort: "...backFrom 'splitDataToSVG'",
				HasDstOp: true, DstPort: "in",
			},
			&svg.Merge{
				ID:   "shapesToSVG",
				Size: 2,
			},
		}, {
			&svg.Op{
				Main: &svg.Rect{
					Text: []string{"shapesToSVG", "[synchronous]"},
				},
			},
			&svg.Split{
				Shapes: [][]interface{}{
					{
						&svg.Arrow{
							DataType: "(arrow, *svgFlow, *x, *y, *moveData)",
							HasSrcOp: true, SrcPort: "arrowOut",
							HasDstOp: true, DstPort: "in",
						},
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{"arrowDataToSVG"},
							},
						},
						&svg.Arrow{
							HasSrcOp: true, SrcPort: "",
							HasDstOp: false, DstPort: "RETURN",
						},
					}, {
						&svg.Arrow{
							DataType: "(op, *svgFlow, completedMerge, *svgMainRect, *y0, *x, *y)",
							HasSrcOp: true, SrcPort: "opOut",
							HasDstOp: true, DstPort: "in",
						},
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{"opDataToSVG"},
							},
						},
						&svg.Arrow{
							HasSrcOp: true, SrcPort: "",
							HasDstOp: false, DstPort: "RETURN",
						},
					}, {
						&svg.Arrow{
							DataType: "(merge, moveData, *completedMerge, x, y)",
							HasSrcOp: true, SrcPort: "mergeOut",
							HasDstOp: true, DstPort: "in",
						},
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{"mergeDataToSVG"},
							},
						},
						&svg.Arrow{
							HasSrcOp: true, SrcPort: "",
							HasDstOp: false, DstPort: "RETURN",
						},
					}, {
						&svg.Arrow{
							DataType: "(split, *svgFlow, *svgMainRect, *x, *y)",
							HasSrcOp: true, SrcPort: "splitOut",
							HasDstOp: true, DstPort: "in",
						},
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{"splitDataToSVG", "[synchronous]"},
							},
						},
						&svg.Split{
							Shapes: [][]interface{}{
								{
									&svg.Arrow{
										DataType: "(shapes, *svgFlow, *x, *y)",
										HasSrcOp: true, SrcPort: "out",
										HasDstOp: false, DstPort: "...backTo 'shapesToSVG'",
									},
								}, {
									&svg.Arrow{
										HasSrcOp: true, SrcPort: "",
										HasDstOp: false, DstPort: "RETURN",
									},
								},
							},
						},
					}, {
						&svg.Arrow{
							HasSrcOp: true, SrcPort: "",
							HasDstOp: false, DstPort: "RETURN",
						},
					},
				},
			},
		},
	},
}
