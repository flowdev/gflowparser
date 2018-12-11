package data2svg_test

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/flowdev/gflowparser/data"
	"github.com/flowdev/gflowparser/data2svg"
	"github.com/flowdev/gflowparser/svg"
	"github.com/flowdev/gparselib"
	"github.com/sanity-io/litter"
)

func TestConvert(t *testing.T) {
	specs := []struct {
		name     string
		given    data.Flow
		expected svg.Flow
		hasError bool
	}{
		{
			name: "simple",
			given: data.Flow{
				Parts: [][]interface{}{
					{
						data.Arrow{
							FromPort: &data.Port{Name: "a"},
							Data:     []data.Type{data.Type{LocalType: "b"}},
						},
						data.Component{
							Decl: data.CompDecl{
								Name:      "c",
								Type:      data.Type{LocalType: "c"},
								VagueType: true,
							},
						},
					},
				},
			},
			expected: svg.Flow{
				Shapes: [][]interface{}{
					{
						&svg.Arrow{
							DataType: "(b)",
							HasSrcOp: false, SrcPort: "a",
							HasDstOp: true,
						},
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{"c"},
							},
							Plugins: []*svg.Plugin{},
						},
					},
				},
			},
			hasError: false,
		}, {
			name: "full_arrows",
			given: data.Flow{
				Parts: [][]interface{}{
					{
						data.Arrow{
							FromPort: &data.Port{Name: "a"},
							Data: []data.Type{
								data.Type{Package: "pack", LocalType: "b"},
							},
							ToPort: &data.Port{Name: "in"},
						},
						data.Component{
							Decl: data.CompDecl{
								Name:      "c",
								Type:      data.Type{LocalType: "c"},
								VagueType: true,
							},
						},
						data.Arrow{
							FromPort: &data.Port{
								Name:     "out",
								HasIndex: true,
								Index:    3,
							},
							Data: []data.Type{
								data.Type{LocalType: "b"},
								data.Type{Package: "pack", LocalType: "Btype"},
								data.Type{Package: "pack", LocalType: "btype"},
							},
							ToPort: &data.Port{
								Name:     "in",
								HasIndex: true,
								Index:    2,
							},
						},
						data.Component{
							Decl: data.CompDecl{
								Name:      "d",
								Type:      data.Type{LocalType: "D"},
								VagueType: false,
							},
						},
						data.Arrow{
							FromPort: &data.Port{Name: "out"},
							Data:     []data.Type{data.Type{LocalType: "b"}},
							ToPort:   &data.Port{Name: "out"},
						},
					},
				},
			},
			expected: svg.Flow{
				Shapes: [][]interface{}{
					{
						&svg.Arrow{
							DataType: "(pack.b)",
							HasSrcOp: false, SrcPort: "a",
							HasDstOp: true, DstPort: "in",
						},
						&svg.Op{
							Main:    &svg.Rect{Text: []string{"c"}},
							Plugins: []*svg.Plugin{},
						},
						&svg.Arrow{
							DataType: "(b, pack.Btype, pack.btype)",
							HasSrcOp: true, SrcPort: "out[3]",
							HasDstOp: true, DstPort: "in[2]",
						},
						&svg.Op{
							Main:    &svg.Rect{Text: []string{"d", "D"}},
							Plugins: []*svg.Plugin{},
						},
						&svg.Arrow{
							DataType: "(b)",
							HasSrcOp: true, SrcPort: "out",
							HasDstOp: false, DstPort: "out",
						},
					},
				},
			},
			hasError: false,
		}, {
			name: "continuations",
			given: data.Flow{
				Parts: [][]interface{}{
					{
						data.Arrow{
							FromPort: &data.Port{Name: "a"},
							Data: []data.Type{
								data.Type{Package: "pack", LocalType: "b"},
							},
							ToPort: &data.Port{Name: "in"},
						},
						data.Component{
							Decl: data.CompDecl{
								Name:      "c",
								Type:      data.Type{LocalType: "c"},
								VagueType: true,
							},
						},
						data.Arrow{
							FromPort: &data.Port{
								Name:     "out",
								HasIndex: true,
								Index:    3,
							},
							ToPort: &data.Port{
								Name:  "...",
								Index: 2,
							},
						},
					}, {
						data.Arrow{
							FromPort: &data.Port{
								Name:  "...",
								Index: 2,
							},
							Data: []data.Type{
								data.Type{LocalType: "b"},
								data.Type{Package: "pack", LocalType: "Btype"},
								data.Type{Package: "pack", LocalType: "btype"},
							},
							ToPort: &data.Port{
								Name:     "in",
								HasIndex: true,
								Index:    2,
							},
						},
						data.Component{
							Decl: data.CompDecl{
								Name:      "d",
								Type:      data.Type{LocalType: "D"},
								VagueType: false,
							},
						},
					},
				},
			},
			expected: svg.Flow{
				Shapes: [][]interface{}{
					{
						&svg.Arrow{
							DataType: "(pack.b)",
							HasSrcOp: false, SrcPort: "a",
							HasDstOp: true, DstPort: "in",
						},
						&svg.Op{
							Main:    &svg.Rect{Text: []string{"c"}},
							Plugins: []*svg.Plugin{},
						},
						&svg.Arrow{
							HasSrcOp: true, SrcPort: "out[3]",
							HasDstOp: false, DstPort: "...2",
						},
					}, {}, {
						&svg.Arrow{
							DataType: "(b, pack.Btype, pack.btype)",
							HasSrcOp: false, SrcPort: "...2",
							HasDstOp: true, DstPort: "in[2]",
						},
						&svg.Op{
							Main:    &svg.Rect{Text: []string{"d", "D"}},
							Plugins: []*svg.Plugin{},
						},
					},
				},
			},
			hasError: false,
		}, {
			name: "full_components",
			given: data.Flow{
				Parts: [][]interface{}{
					{
						data.Component{
							Decl: data.CompDecl{
								Name:      "a",
								Type:      data.Type{Package: "p", LocalType: "A"},
								VagueType: false,
							},
							Plugins: []data.Plugin{
								{
									Types: []data.Type{data.Type{Package: "q", LocalType: "Z"}},
								},
							},
						},
						data.Arrow{Data: []data.Type{data.Type{LocalType: "data"}}},
						data.Component{
							Decl: data.CompDecl{
								Name:      "b",
								Type:      data.Type{Package: "p", LocalType: "B"},
								VagueType: false,
							},
							Plugins: []data.Plugin{
								{
									Types: []data.Type{
										data.Type{Package: "q", LocalType: "Z"},
										data.Type{Package: "q", LocalType: "Y"},
										data.Type{Package: "q", LocalType: "X"},
										data.Type{Package: "q", LocalType: "W"},
									},
								},
							},
						},
						data.Arrow{},
						data.Component{
							Decl: data.CompDecl{
								Name:      "c",
								Type:      data.Type{LocalType: "c"},
								VagueType: true,
							},
							Plugins: []data.Plugin{
								{
									Name: "plugin1",
									Types: []data.Type{
										data.Type{Package: "q", LocalType: "V"},
										data.Type{Package: "q", LocalType: "U"},
									},
								}, {
									Name: "plugin2",
									Types: []data.Type{
										data.Type{Package: "q", LocalType: "T"},
									},
								}, {
									Types: []data.Type{
										data.Type{Package: "q", LocalType: "Plugin3"},
									},
								},
							},
						},
						data.Arrow{},
						data.Component{
							Decl: data.CompDecl{
								Name:      "d",
								Type:      data.Type{LocalType: "D"},
								VagueType: false,
							},
							Plugins: []data.Plugin{
								{
									Types: []data.Type{
										data.Type{Package: "q", LocalType: "T"},
									},
								}, {
									Types: []data.Type{
										data.Type{Package: "q", LocalType: "U"},
									},
								}, {
									Types: []data.Type{
										data.Type{Package: "q", LocalType: "V"},
									},
								}, {
									Types: []data.Type{
										data.Type{Package: "q", LocalType: "W"},
									},
								}, {
									Types: []data.Type{
										data.Type{Package: "q", LocalType: "X"},
									},
								}, {
									Types: []data.Type{
										data.Type{Package: "q", LocalType: "Y"},
									},
								}, {
									Types: []data.Type{
										data.Type{Package: "q", LocalType: "Z"},
									},
								},
							},
						},
					},
				},
			},
			expected: svg.Flow{
				Shapes: [][]interface{}{
					{
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{"a", "p.A"},
							},
							Plugins: []*svg.Plugin{
								&svg.Plugin{
									Rects: []*svg.Rect{
										&svg.Rect{Text: []string{"q.Z"}},
									},
								},
							},
						},
						&svg.Arrow{
							DataType: "(data)",
							HasSrcOp: true, HasDstOp: true,
						},
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{"b", "p.B"},
							},
							Plugins: []*svg.Plugin{
								&svg.Plugin{
									Rects: []*svg.Rect{
										&svg.Rect{Text: []string{"q.Z"}},
										&svg.Rect{Text: []string{"q.Y"}},
										&svg.Rect{Text: []string{"q.X"}},
										&svg.Rect{Text: []string{"q.W"}},
									},
								},
							},
						},
						&svg.Arrow{HasSrcOp: true, HasDstOp: true},
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{"c"},
							},
							Plugins: []*svg.Plugin{
								&svg.Plugin{
									Title: "plugin1",
									Rects: []*svg.Rect{
										&svg.Rect{Text: []string{"q.V"}},
										&svg.Rect{Text: []string{"q.U"}},
									},
								},
								&svg.Plugin{
									Title: "plugin2",
									Rects: []*svg.Rect{
										&svg.Rect{Text: []string{"q.T"}},
									},
								},
								&svg.Plugin{
									Rects: []*svg.Rect{
										&svg.Rect{Text: []string{"q.Plugin3"}},
									},
								},
							},
						},
						&svg.Arrow{HasSrcOp: true, HasDstOp: true},
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{"d", "D"},
							},
							Plugins: []*svg.Plugin{
								&svg.Plugin{
									Rects: []*svg.Rect{
										&svg.Rect{Text: []string{"q.T"}},
									},
								},
								&svg.Plugin{
									Rects: []*svg.Rect{
										&svg.Rect{Text: []string{"q.U"}},
									},
								},
								&svg.Plugin{
									Rects: []*svg.Rect{
										&svg.Rect{Text: []string{"q.V"}},
									},
								},
								&svg.Plugin{
									Rects: []*svg.Rect{
										&svg.Rect{Text: []string{"q.W"}},
									},
								},
								&svg.Plugin{
									Rects: []*svg.Rect{
										&svg.Rect{Text: []string{"q.X"}},
									},
								},
								&svg.Plugin{
									Rects: []*svg.Rect{
										&svg.Rect{Text: []string{"q.Y"}},
									},
								},
								&svg.Plugin{
									Rects: []*svg.Rect{
										&svg.Rect{Text: []string{"q.Z"}},
									},
								},
							},
						},
					},
				},
			},
			hasError: false,
		}, {
			name: "splits_n_merges",
			given: data.Flow{
				Parts: [][]interface{}{
					{ // [a A] (data)-> [b B] -> [c C]
						data.Component{
							Decl: data.CompDecl{
								Name: "a", VagueType: false,
								Type: data.Type{LocalType: "A"},
							},
						},
						data.Arrow{Data: []data.Type{data.Type{LocalType: "data"}}},
						data.Component{
							Decl: data.CompDecl{
								Name: "b", VagueType: false,
								Type: data.Type{LocalType: "B"},
							},
						},
						data.Arrow{},
						data.Component{
							Decl: data.CompDecl{
								Name: "c", VagueType: false,
								Type: data.Type{LocalType: "C"},
							},
						},
					}, { // [a] (data)-> [h] -> [b] -> [i]
						data.Component{
							Decl: data.CompDecl{
								Name: "a", VagueType: true,
								Type: data.Type{LocalType: "a"},
							},
						},
						data.Arrow{Data: []data.Type{data.Type{LocalType: "data"}}},
						data.Component{
							Decl: data.CompDecl{
								Name: "h", VagueType: true,
								Type: data.Type{LocalType: "h"},
							},
						},
						data.Arrow{},
						data.Component{
							Decl: data.CompDecl{
								Name: "b", VagueType: true,
								Type: data.Type{LocalType: "b"},
							},
						},
						data.Arrow{},
						data.Component{
							Decl: data.CompDecl{
								Name: "i", VagueType: true,
								Type: data.Type{LocalType: "i"},
							},
						},
					}, { // [a] (data)-> [d D] -> [c]
						data.Component{
							Decl: data.CompDecl{
								Name: "a", VagueType: true,
								Type: data.Type{LocalType: "a"},
							},
						},
						data.Arrow{Data: []data.Type{data.Type{LocalType: "data"}}},
						data.Component{
							Decl: data.CompDecl{
								Name: "d", VagueType: false,
								Type: data.Type{LocalType: "D"},
							},
						},
						data.Arrow{},
						data.Component{
							Decl: data.CompDecl{
								Name: "c", VagueType: true,
								Type: data.Type{LocalType: "c"},
							},
						},
					}, { // [a] (data)-> [e E] -> [c]
						data.Component{
							Decl: data.CompDecl{
								Name: "a", VagueType: true,
								Type: data.Type{LocalType: "a"},
							},
						},
						data.Arrow{Data: []data.Type{data.Type{LocalType: "data"}}},
						data.Component{
							Decl: data.CompDecl{
								Name: "e", VagueType: false,
								Type: data.Type{LocalType: "E"},
							},
						},
						data.Arrow{},
						data.Component{
							Decl: data.CompDecl{
								Name: "c", VagueType: true,
								Type: data.Type{LocalType: "c"},
							},
						},
					}, { // [b] (data)-> [f F] -> [g] -> [c]
						data.Component{
							Decl: data.CompDecl{
								Name: "b", VagueType: true,
								Type: data.Type{LocalType: "b"},
							},
						},
						data.Arrow{Data: []data.Type{data.Type{LocalType: "data"}}},
						data.Component{
							Decl: data.CompDecl{
								Name: "f", VagueType: false,
								Type: data.Type{LocalType: "F"},
							},
						},
						data.Arrow{},
						data.Component{
							Decl: data.CompDecl{
								Name: "g", VagueType: true,
								Type: data.Type{LocalType: "g"},
							},
						},
						data.Arrow{},
						data.Component{
							Decl: data.CompDecl{
								Name: "c", VagueType: true,
								Type: data.Type{LocalType: "c"},
							},
						},
					}, { // [j J] (data)-> [a]
						data.Component{
							Decl: data.CompDecl{
								Name: "j", VagueType: false,
								Type: data.Type{LocalType: "J"},
							},
						},
						data.Arrow{Data: []data.Type{data.Type{LocalType: "data"}}},
						data.Component{
							Decl: data.CompDecl{
								Name: "a", VagueType: true,
								Type: data.Type{LocalType: "a"},
							},
						},
					}, { // [k K] (data)-> [a]
						data.Component{
							Decl: data.CompDecl{
								Name: "k", VagueType: false,
								Type: data.Type{LocalType: "K"},
							},
						},
						data.Arrow{Data: []data.Type{data.Type{LocalType: "data"}}},
						data.Component{
							Decl: data.CompDecl{
								Name: "a", VagueType: true,
								Type: data.Type{LocalType: "a"},
							},
						},
					}, { // [l] (data)-> [k]
						data.Component{
							Decl: data.CompDecl{
								Name: "l", VagueType: true,
								Type: data.Type{LocalType: "l"},
							},
						},
						data.Arrow{Data: []data.Type{data.Type{LocalType: "data"}}},
						data.Component{
							Decl: data.CompDecl{
								Name: "k", VagueType: true,
								Type: data.Type{LocalType: "k"},
							},
						},
					},
				},
			},
			expected: svg.Flow{
				Shapes: [][]interface{}{
					{
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{"j", "J"},
							},
							Plugins: []*svg.Plugin{},
						},
						&svg.Arrow{
							DataType: "(data)",
							HasSrcOp: true, HasDstOp: true,
						},
						&svg.Merge{ID: "a", Size: 2},
					}, {}, { // TODO: The unnecessary empty slice isn't nice.
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{"l"},
							},
							Plugins: []*svg.Plugin{},
						},
						&svg.Arrow{
							DataType: "(data)",
							HasSrcOp: true, HasDstOp: true,
						},
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{"k", "K"},
							},
							Plugins: []*svg.Plugin{},
						},
						&svg.Arrow{
							DataType: "(data)",
							HasSrcOp: true, HasDstOp: true,
						},
						&svg.Merge{ID: "a", Size: 2},
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{"a", "A"},
							},
							Plugins: []*svg.Plugin{},
						},
						&svg.Split{
							Shapes: [][]interface{}{
								{
									&svg.Arrow{
										DataType: "(data)",
										HasSrcOp: true, HasDstOp: true,
									},
									&svg.Merge{ID: "b", Size: 2},
								}, {
									&svg.Arrow{
										DataType: "(data)",
										HasSrcOp: true, HasDstOp: true,
									},
									&svg.Op{
										Main: &svg.Rect{
											Text: []string{"h"},
										},
										Plugins: []*svg.Plugin{},
									},
									&svg.Arrow{HasSrcOp: true, HasDstOp: true},
									&svg.Merge{ID: "b", Size: 2},
									&svg.Op{
										Main: &svg.Rect{
											Text: []string{"b", "B"},
										},
										Plugins: []*svg.Plugin{},
									},
									&svg.Split{
										Shapes: [][]interface{}{
											{
												&svg.Arrow{HasSrcOp: true, HasDstOp: true},
												&svg.Merge{ID: "c", Size: 4},
											}, {
												&svg.Arrow{HasSrcOp: true, HasDstOp: true},
												&svg.Op{
													Main: &svg.Rect{
														Text: []string{"i"},
													},
													Plugins: []*svg.Plugin{},
												},
											}, {
												&svg.Arrow{
													DataType: "(data)",
													HasSrcOp: true, HasDstOp: true,
												},
												&svg.Op{
													Main: &svg.Rect{
														Text: []string{"f", "F"},
													},
													Plugins: []*svg.Plugin{},
												},
												&svg.Arrow{HasSrcOp: true, HasDstOp: true},
												&svg.Op{
													Main: &svg.Rect{
														Text: []string{"g"},
													},
													Plugins: []*svg.Plugin{},
												},
												&svg.Arrow{HasSrcOp: true, HasDstOp: true},
												&svg.Merge{ID: "c", Size: 4},
											},
										},
									},
								}, {
									&svg.Arrow{
										DataType: "(data)",
										HasSrcOp: true, HasDstOp: true,
									},
									&svg.Op{
										Main: &svg.Rect{
											Text: []string{"d", "D"},
										},
										Plugins: []*svg.Plugin{},
									},
									&svg.Arrow{HasSrcOp: true, HasDstOp: true},
									&svg.Merge{ID: "c", Size: 4},
								}, {
									&svg.Arrow{
										DataType: "(data)",
										HasSrcOp: true, HasDstOp: true,
									},
									&svg.Op{
										Main: &svg.Rect{
											Text: []string{"e", "E"},
										},
										Plugins: []*svg.Plugin{},
									},
									&svg.Arrow{HasSrcOp: true, HasDstOp: true},
									&svg.Merge{ID: "c", Size: 4},
									&svg.Op{
										Main: &svg.Rect{
											Text: []string{
												"c",
												"C",
											},
										},
										Plugins: []*svg.Plugin{},
									},
								},
							},
						},
					},
				},
			},
			hasError: false,
		}, {
			name: "circles",
			given: data.Flow{
				Parts: [][]interface{}{
					{
						data.Arrow{
							FromPort: &data.Port{Name: "i"},
							Data:     []data.Type{data.Type{LocalType: "data"}},
						},
						data.Component{
							Decl: data.CompDecl{
								Name:      "a",
								Type:      data.Type{LocalType: "A"},
								VagueType: false,
							},
						},
						data.Arrow{},
						data.Component{
							Decl: data.CompDecl{
								Name:      "b",
								Type:      data.Type{LocalType: "B"},
								VagueType: false,
							},
						},
						data.Arrow{},
						data.Component{
							Decl: data.CompDecl{
								Name:      "a",
								Type:      data.Type{LocalType: "a"},
								VagueType: true,
							},
						},
					}, {
						data.Component{
							Decl: data.CompDecl{
								Name:      "b",
								Type:      data.Type{LocalType: "b"},
								VagueType: true,
							},
						},
						data.Arrow{
							Data: []data.Type{data.Type{LocalType: "data"}},
						},
						data.Component{
							Decl: data.CompDecl{
								Name:      "c",
								Type:      data.Type{LocalType: "C"},
								VagueType: false,
							},
						},
						data.Arrow{},
						data.Component{
							Decl: data.CompDecl{
								Name:      "a",
								Type:      data.Type{LocalType: "a"},
								VagueType: true,
							},
						},
					},
				},
			},
			expected: svg.Flow{
				Shapes: [][]interface{}{
					{
						&svg.Arrow{
							DataType: "(data)", SrcPort: "i",
							HasSrcOp: false, HasDstOp: true,
						},
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{"a", "A"},
							},
							Plugins: []*svg.Plugin{},
						},
						&svg.Arrow{HasSrcOp: true, HasDstOp: true},
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{"b", "B"},
							},
							Plugins: []*svg.Plugin{},
						},
						&svg.Split{
							Shapes: [][]interface{}{
								{
									&svg.Arrow{HasSrcOp: true, HasDstOp: true},
									&svg.Rect{
										Text: []string{"a"},
									},
								}, {
									&svg.Arrow{
										DataType: "(data)",
										HasSrcOp: true, HasDstOp: true,
									},
									&svg.Op{
										Main: &svg.Rect{
											Text: []string{"c", "C"},
										},
										Plugins: []*svg.Plugin{},
									},
									&svg.Arrow{HasSrcOp: true, HasDstOp: true},
									&svg.Rect{
										Text: []string{"a"},
									},
								},
							},
						},
					},
				},
			},
			hasError: false,
		}, {
			name: "empty_rows",
			given: data.Flow{
				Parts: [][]interface{}{
					{
						data.Component{
							Decl: data.CompDecl{
								Name:      "a",
								Type:      data.Type{LocalType: "A"},
								VagueType: false,
							},
						},
						data.Arrow{
							Data: []data.Type{data.Type{LocalType: "data"}},
						},
						data.Component{
							Decl: data.CompDecl{
								Name:      "b",
								Type:      data.Type{LocalType: "B"},
								VagueType: false,
							},
						},
					}, {
						data.Component{
							Decl: data.CompDecl{
								Name:      "c",
								Type:      data.Type{LocalType: "C"},
								VagueType: false,
							},
						},
						data.Arrow{
							Data: []data.Type{data.Type{LocalType: "data"}},
						},
						data.Component{
							Decl: data.CompDecl{
								Name:      "d",
								Type:      data.Type{LocalType: "D"},
								VagueType: false,
							},
						},
					}, {
						data.Component{
							Decl: data.CompDecl{
								Name:      "e",
								Type:      data.Type{LocalType: "E"},
								VagueType: false,
							},
						},
						data.Arrow{
							Data: []data.Type{data.Type{LocalType: "data"}},
						},
						data.Component{
							Decl: data.CompDecl{
								Name:      "f",
								Type:      data.Type{LocalType: "F"},
								VagueType: false,
							},
						},
					},
				},
			},
			expected: svg.Flow{
				Shapes: [][]interface{}{
					[]interface{}{
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{"a", "A"},
							},
							Plugins: []*svg.Plugin{},
						},
						&svg.Arrow{
							DataType: "(data)",
							HasSrcOp: true, HasDstOp: true,
						},
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{"b", "B"},
							},
							Plugins: []*svg.Plugin{},
						},
					}, {}, {
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{"c", "C"},
							},
							Plugins: []*svg.Plugin{},
						},
						&svg.Arrow{
							DataType: "(data)",
							HasSrcOp: true, HasDstOp: true,
						},
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{"d", "D"},
							},
							Plugins: []*svg.Plugin{},
						},
					}, {}, {
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{
									"e",
									"E",
								},
							},
							Plugins: []*svg.Plugin{},
						},
						&svg.Arrow{
							DataType: "(data)",
							HasSrcOp: true,
							SrcPort:  "",
							HasDstOp: true,
							DstPort:  "",
						},
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{
									"f",
									"F",
								},
							},
							Plugins: []*svg.Plugin{},
						},
					},
				},
			},
			hasError: false,
		}, {
			name: "error",
			given: data.Flow{
				Parts: [][]interface{}{
					{
						data.Arrow{
							FromPort: &data.Port{Name: "a"},
							Data:     []data.Type{data.Type{LocalType: "b"}},
						},
						data.Component{
							Decl: data.CompDecl{
								Name:      "c",
								Type:      data.Type{LocalType: "C"},
								VagueType: false,
							},
						},
						data.Arrow{},
						data.Component{
							Decl: data.CompDecl{
								Name:      "c",
								Type:      data.Type{LocalType: "c"},
								VagueType: false,
							},
						},
					},
				},
			},
			expected: svg.Flow{},
			hasError: true,
		},
	}

	for _, spec := range specs {
		t.Logf("Testing spec: %s\n", spec.name)
		got, err := data2svg.Convert(
			spec.given,
			gparselib.NewSourceData("test data", "sad but true: <undefined>"),
		)
		if spec.hasError && err != nil {
			continue
		} else if spec.hasError && err == nil {
			t.Error("Expected an error but didn't get one.")
			continue
		} else if !spec.hasError && err != nil {
			t.Errorf("Expected no error but got: %s", err)
			continue
		}

		checkValue(spec.expected, got, spec.name, t)
	}
}

func checkValue(expected, got interface{}, name string, t *testing.T) {
	if expected != nil && got == nil {
		t.Errorf("Expected a value.")
	} else if !reflect.DeepEqual(got, expected) {
		goprint := litter.Options{
			HidePrivateFields: false,
			HomePackage:       "data2svg",
			StripPackageNames: false,
		}
		t.Errorf("The acutal value isn't equal to the expected one. Please diff '%s.exp' and '%s.got'.", name, name)
		if err := ioutil.WriteFile(name+".exp", []byte(goprint.Sdump(expected)), 0644); err != nil {
			panic("Unable to write file: " + name + ".exp")
		}
		if err := ioutil.WriteFile(name+".got", []byte(goprint.Sdump(got)), 0644); err != nil {
			panic("Unable to write file: " + name + ".got")
		}
	}
}
