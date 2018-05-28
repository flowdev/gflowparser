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
		},
	}

	for _, spec := range specs {
		t.Logf("Testing spec: %s\n", spec.name)
		gotAll := parserToSVGData(spec.given)
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
