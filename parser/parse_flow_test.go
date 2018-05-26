package parser

import (
	"testing"

	"github.com/flowdev/gflowparser/data"
)

func TestParsePort(t *testing.T) {
	p, err := NewParsePort()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	runTests(t, p.In, []parseTestData{
		{
			givenName:        "empty",
			givenContent:     ``,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:        "no match 1",
			givenContent:     `1:a`,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:        "no match 2",
			givenContent:     `_a`,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:        "simple 1",
			givenContent:     `aB`,
			expectedValue:    data.Port{Name: "aB"},
			expectedErrCount: 0,
		}, {
			givenName:        "simple 2",
			givenContent:     `a0`,
			expectedValue:    data.Port{Name: "a0"},
			expectedErrCount: 0,
		}, {
			givenName:        "simple 3",
			givenContent:     `aB_cd`,
			expectedValue:    data.Port{Name: "aB"},
			expectedErrCount: 0,
		}, {
			givenName:        "simple 4",
			givenContent:     `abcDef`,
			expectedValue:    data.Port{Name: "abcDef"},
			expectedErrCount: 0,
		}, {
			givenName:        "complex 1",
			givenContent:     `ab1Cd:1`,
			expectedValue:    data.Port{Name: "ab1Cd", HasIndex: true, Index: 1},
			expectedErrCount: 0,
		}, {
			givenName:        "complex 2",
			givenContent:     `a1Bc:003`,
			expectedValue:    data.Port{Name: "a1Bc", HasIndex: true, Index: 3},
			expectedErrCount: 0,
		},
	})
}

func TestParseArrow(t *testing.T) {
	p, err := NewParseArrow()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	runTests(t, p.In, []parseTestData{
		{
			givenName:        "empty",
			givenContent:     ``,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:        "no match 1",
			givenContent:     `>`,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:        "no match 2",
			givenContent:     `aPort (pack.Data)- bPort`,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:        "simple 1",
			givenContent:     `->`,
			expectedValue:    data.Arrow{},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 2",
			givenContent: `aPort->bPort`,
			expectedValue: data.Arrow{
				FromPort: &data.Port{Name: "aPort"},
				ToPort:   &data.Port{Name: "bPort", SrcPos: 7},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 3",
			givenContent: `(Data)->`,
			expectedValue: data.Arrow{
				Data: []data.Type{data.Type{LocalType: "Data", SrcPos: 1}},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "complex",
			givenContent: "aPort // comment1\n ( \t Data \t ) \t -> // comment2\n bPort",
			expectedValue: data.Arrow{
				FromPort: &data.Port{Name: "aPort"},
				Data:     []data.Type{data.Type{LocalType: "Data", SrcPos: 23}},
				ToPort:   &data.Port{Name: "bPort", SrcPos: 50},
			},
			expectedErrCount: 0,
		},
	})
}

func TestParseFlow(t *testing.T) {
	p, err := NewParseFlow()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	runTests(t, p.In, []parseTestData{
		{
			givenName:        "empty",
			givenContent:     ``,
			expectedValue:    nil,
			expectedErrCount: 5,
		}, {
			givenName:        "first input port missing",
			givenContent:     `(dat)->[A]`,
			expectedValue:    nil,
			expectedErrCount: 2,
		}, {
			givenName:        "component missing",
			givenContent:     `aPort (pack.Data)-> bPort`,
			expectedValue:    nil,
			expectedErrCount: 5,
		}, {
			givenName:        "data of first arrow missing",
			givenContent:     `[A]->out`,
			expectedValue:    nil,
			expectedErrCount: 2,
		}, {
			givenName:        "two consecutive arrows",
			givenContent:     `in(Data)->->out`,
			expectedValue:    nil,
			expectedErrCount: 2,
		}, {
			givenName:        "two consecutive components",
			givenContent:     `[A][B]`,
			expectedValue:    nil,
			expectedErrCount: 2,
		}, {
			givenName:        "no end",
			givenContent:     `a(b)->[c]`,
			expectedValue:    nil,
			expectedErrCount: 2,
		}, {
			givenName:    "simple 1",
			givenContent: `a(b)->[c];`,
			expectedValue: data.Flow{
				Parts: [][]interface{}{
					{
						data.Arrow{
							FromPort: &data.Port{Name: "a"},
							Data:     []data.Type{data.Type{LocalType: "b", SrcPos: 2}},
						},
						data.Component{
							Decl: data.CompDecl{
								Name:      "c",
								Type:      data.Type{LocalType: "c", SrcPos: 7},
								VagueType: true,
								SrcPos:    7,
							},
							SrcPos: 6,
						},
					},
				},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 2",
			givenContent: "[A](b)->c // my comment\n",
			expectedValue: data.Flow{
				Parts: [][]interface{}{
					{
						data.Component{Decl: data.CompDecl{
							Name:   "a",
							Type:   data.Type{LocalType: "A", SrcPos: 1},
							SrcPos: 1,
						}},
						data.Arrow{
							Data:   []data.Type{data.Type{LocalType: "b", SrcPos: 4}},
							ToPort: &data.Port{Name: "c", SrcPos: 8},
							SrcPos: 3,
						},
					},
				},
			},
			expectedErrCount: 0,
		},
	})
}
