package data2svg

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	data "github.com/flowdev/gflowparser"
	"github.com/flowdev/gflowparser/svg"
	"github.com/flowdev/gparselib"
)

func init() {
	spew.Config.Indent = "    "
	spew.Config.DisablePointerAddresses = true
	spew.Config.DisableCapacities = true
	spew.Config.SortKeys = true
	spew.Config.SpewKeys = true
}

// func parserToSVGData(flowDat data.Flow) *svg.Flow {
func TestParserToSVGData(t *testing.T) {
	specs := []struct {
		name     string
		given    data.Flow
		expected [][]interface{}
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
			expected: [][]interface{}{
				{
					&svg.Arrow{
						HasSrcOp: false, SrcPort: "a",
						DataType: "(b)",
						HasDstOp: true, DstPort: "",
					},
					&decl{
						name: "c",
						i:    0, j: 1,
						svgOp: &svg.Op{
							Main: &svg.Rect{
								Text: []string{"c"},
							},
							Plugins: []*svg.Plugin{},
						},
						svgMerge: &svg.Merge{ID: "c", Size: 1},
					},
				},
			},
			hasError: false,
		}, {
			name: "full arrows",
			given: data.Flow{
				Parts: [][]interface{}{
					{
						data.Arrow{
							FromPort: &data.Port{Name: "a"},
							Data:     []data.Type{data.Type{LocalType: "b"}},
							ToPort:   &data.Port{Name: "in"},
						},
						data.Component{
							Decl: data.CompDecl{
								Name:      "c",
								Type:      data.Type{LocalType: "c"},
								VagueType: true,
							},
						},
						data.Arrow{
							FromPort: &data.Port{Name: "out"},
							Data:     []data.Type{data.Type{LocalType: "b"}},
							ToPort:   &data.Port{Name: "in"},
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
			expected: [][]interface{}{
				{
					&svg.Arrow{
						HasSrcOp: false, SrcPort: "a",
						DataType: "(b)",
						HasDstOp: true, DstPort: "in",
					},
					&decl{
						name: "c",
						i:    0, j: 1,
						svgOp: &svg.Op{
							Main: &svg.Rect{
								Text: []string{"c"},
							},
							Plugins: []*svg.Plugin{},
						},
						svgMerge: &svg.Merge{ID: "c", Size: 1},
					},
					&svg.Arrow{
						HasSrcOp: true, SrcPort: "out",
						DataType: "(b)",
						HasDstOp: true, DstPort: "in",
					},
					&decl{
						name: "d",
						i:    0, j: 3,
						svgOp: &svg.Op{
							Main: &svg.Rect{
								Text: []string{"d", "D"},
							},
							Plugins: []*svg.Plugin{},
						},
						svgMerge: &svg.Merge{ID: "d", Size: 1},
					},
					&svg.Arrow{
						HasSrcOp: true, SrcPort: "out",
						DataType: "(b)",
						HasDstOp: false, DstPort: "out",
					},
				},
			},
			hasError: false,
		}, {
			name: "full components",
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
						data.Arrow{},
						data.Component{
							Decl: data.CompDecl{
								Name:      "b",
								Type:      data.Type{Package: "p", LocalType: "B"},
								VagueType: false,
							},
							Plugins: []data.Plugin{
								{
									Types: []data.Type{
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
					},
				},
			},
			expected: [][]interface{}{
				{
					&decl{
						name: "a",
						i:    0, j: 0,
						svgOp: &svg.Op{
							Main: &svg.Rect{
								Text: []string{"a", "p.A"},
							},
							Plugins: []*svg.Plugin{
								{
									Rects: []*svg.Rect{
										{Text: []string{"q.Z"}},
									},
								},
							},
						},
						svgMerge: &svg.Merge{ID: "a", Size: 0},
					},
					&svg.Arrow{HasSrcOp: true, HasDstOp: true},
					&decl{
						name: "b",
						i:    0, j: 2,
						svgOp: &svg.Op{
							Main: &svg.Rect{
								Text: []string{"b", "p.B"},
							},
							Plugins: []*svg.Plugin{
								{
									Rects: []*svg.Rect{
										{Text: []string{"q.Y"}},
										{Text: []string{"q.X"}},
										{Text: []string{"q.W"}},
									},
								},
							},
						},
						svgMerge: &svg.Merge{ID: "b", Size: 1},
					},
					&svg.Arrow{HasSrcOp: true, HasDstOp: true},
					&decl{
						name: "c",
						i:    0, j: 4,
						svgOp: &svg.Op{
							Main: &svg.Rect{
								Text: []string{"c"},
							},
							Plugins: []*svg.Plugin{
								{
									Title: "plugin1",
									Rects: []*svg.Rect{
										{Text: []string{"q.V"}},
										{Text: []string{"q.U"}},
									},
								}, {
									Title: "plugin2",
									Rects: []*svg.Rect{
										{Text: []string{"q.T"}},
									},
								}, {
									Rects: []*svg.Rect{
										{Text: []string{"q.Plugin3"}},
									},
								},
							},
						},
						svgMerge: &svg.Merge{ID: "c", Size: 1},
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
			expected: nil,
			hasError: true,
		},
	}

	for _, spec := range specs {
		t.Logf("Testing spec: %s\n", spec.name)
		gotAll, _, _, err := parserPartsToSVGData(
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

		if len(gotAll) != len(spec.expected) {
			t.Errorf("Expected %d part lines, got: %d",
				len(spec.expected), len(gotAll))
			continue
		}
		for i, expectedLine := range spec.expected {
			t.Logf("Testing part line: %d\n", i+1)
			gotLine := gotAll[i]
			if len(gotLine) != len(expectedLine) {
				t.Errorf("Expected %d parts in line %d, got: %d",
					len(expectedLine), i+1, len(gotLine))
				continue
			}
			for j, expectedPart := range expectedLine {
				t.Logf("Testing part: %d\n", j+1)
				gotPart := gotLine[j]
				checkValue(expectedPart, gotPart, t)
			}
		}
	}
}

func checkValue(expected, got interface{}, t *testing.T) {
	if expected != nil && got == nil {
		t.Errorf("Expected a value.")
	} else if !reflect.DeepEqual(got, expected) {
		t.Logf("The acutal value isn't equal to the expected one:")
		// Use this for compact report (excellent for quick manual usage):
		t.Error(spew.Sprintf("Expected value:\n '%#v',\nGot value:\n '%#v'", expected, got))
		// Use this for multi-line report (excellent for diff):
		//t.Errorf("Expected value:\n%s\nGot value:\n%s", spew.Sdump(expected), spew.Sdump(got))
	}
}
