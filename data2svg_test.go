package gflowparser

import (
	"reflect"
	"testing"

	"github.com/flowdev/gflowparser/data"
	"github.com/flowdev/gflowparser/svg"
)

// func parserToSVGData(flowDat data.Flow) *svg.Flow {
func TestParserToSVGData(t *testing.T) {
	specs := []struct {
		name     string
		given    data.Flow
		expected *svg.Flow
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
			expected: &svg.Flow{
				Shapes: [][]interface{}{
					{
						&svg.Arrow{
							HasSrcOp: false, SrcPort: "a",
							DataType: "(b)",
							HasDstOp: true, DstPort: "",
						},
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{"c", "c"},
							},
							Plugins: []*svg.Plugin{},
						},
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
			expected: &svg.Flow{
				Shapes: [][]interface{}{
					{
						&svg.Arrow{
							HasSrcOp: false, SrcPort: "a",
							DataType: "(b)",
							HasDstOp: true, DstPort: "in",
						},
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{"c", "c"},
							},
							Plugins: []*svg.Plugin{},
						},
						&svg.Arrow{
							HasSrcOp: true, SrcPort: "out",
							DataType: "(b)",
							HasDstOp: true, DstPort: "in",
						},
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{"d", "D"},
							},
							Plugins: []*svg.Plugin{},
						},
						&svg.Arrow{
							HasSrcOp: true, SrcPort: "out",
							DataType: "(b)",
							HasDstOp: false, DstPort: "out",
						},
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
							Plugins: []data.NameNTypes{
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
							Plugins: []data.NameNTypes{
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
							Plugins: []data.NameNTypes{
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
									Name: "plugin3",
									Types: []data.Type{
										data.Type{Package: "q", LocalType: "S"},
									},
								},
							},
						},
					},
				},
			},
			expected: &svg.Flow{
				Shapes: [][]interface{}{
					{
						&svg.Op{
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
						&svg.Arrow{HasSrcOp: true, HasDstOp: true},
						&svg.Op{
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
						&svg.Arrow{HasSrcOp: true, HasDstOp: true},
						&svg.Op{
							Main: &svg.Rect{
								Text: []string{"c", "c"},
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
									Title: "plugin3",
									Rects: []*svg.Rect{
										{Text: []string{"q.S"}},
									},
								},
							},
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
								VagueType: true,
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
		gotAll, _, err := parserPartsToSVGData(spec.given)
		if spec.hasError && err != nil {
			continue
		} else if spec.hasError && err == nil {
			t.Error("Expected an error but didn't get one.")
			continue
		} else if !spec.hasError && err != nil {
			t.Errorf("Expected no error but got: %s", err)
			continue
		}
		if len(gotAll.Shapes) != len(spec.expected.Shapes) {
			t.Errorf("Expected %d part lines, got: %d",
				len(spec.expected.Shapes), len(gotAll.Shapes))
			continue
		}
		for i, expectedLine := range spec.expected.Shapes {
			t.Logf("Testing part line: %d\n", i+1)
			gotLine := gotAll.Shapes[i]
			if len(gotLine) != len(expectedLine) {
				t.Errorf("Expected %d parts in line %d, got: %d",
					len(expectedLine), i+1, len(gotLine))
				continue
			}
			for j, expectedPart := range expectedLine {
				t.Logf("Testing part: %d\n", j+1)
				gotPart := gotLine[j]
				checkValue(gotPart, expectedPart, t)
			}
		}
	}
}

func checkValue(got, expected interface{}, t *testing.T) {
	if expected != nil && got == nil {
		t.Errorf("Expected a value.")
	} else if !reflect.DeepEqual(got, expected) {
		t.Logf("The acutal value isn't equal to the expected one:")
		t.Logf("Expected value of type '%T', got '%T'.", expected, got)
		t.Errorf("Expected value '%#v', got '%#v'.", expected, got)
	}
}
