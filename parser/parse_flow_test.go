package parser

import (
	"testing"

	"github.com/flowdev/gflowparser/data"
)

func TestParsePort(t *testing.T) {
	p, err := NewPortParser()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	runTests(t, p.ParsePort, []parseTestData{
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
	p, err := NewArrowParser()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	runTests(t, p.ParseArrow, []parseTestData{
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
			givenContent: "aPort \t ( // comment1\n Data // comment2\n ) \t -> \t bPort",
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
	p, err := NewFlowParser()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	runTests(t, p.ParseFlow, []parseTestData{
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
			givenName:        "wrong new line",
			givenContent:     "a(b)->\n[c];",
			expectedValue:    nil,
			expectedErrCount: 5,
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
		}, {
			givenName:    "complex 1",
			givenContent: "[A](b)->c // my comment\nd \t (e)-> \t f \t [G] \t -> \t h;",
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
					}, {
						data.Arrow{
							FromPort: &data.Port{Name: "d", SrcPos: 24},
							Data:     []data.Type{data.Type{LocalType: "e", SrcPos: 29}},
							ToPort:   &data.Port{Name: "f", SrcPos: 36},
							SrcPos:   24,
						},
						data.Component{
							Decl: data.CompDecl{
								Name:   "g",
								Type:   data.Type{LocalType: "G", SrcPos: 41},
								SrcPos: 41,
							},
							SrcPos: 40,
						},
						data.Arrow{
							ToPort: &data.Port{Name: "h", SrcPos: 51},
							SrcPos: 46,
						},
					},
				},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "complex 2",
			givenContent: "[a B [c=D,E|f=G,H]] (i)-> out\n",
			expectedValue: data.Flow{
				Parts: [][]interface{}{
					{
						data.Component{
							Decl: data.CompDecl{
								Name: "a", Type: data.Type{LocalType: "B", SrcPos: 3},
								SrcPos: 1,
							},
							Plugins: []data.NameNTypes{
								data.NameNTypes{
									Name: "c",
									Types: []data.Type{
										data.Type{LocalType: "D", SrcPos: 8},
										data.Type{LocalType: "E", SrcPos: 10},
									},
									SrcPos: 6,
								},
								data.NameNTypes{
									Name: "f",
									Types: []data.Type{
										data.Type{LocalType: "G", SrcPos: 14},
										data.Type{LocalType: "H", SrcPos: 16},
									},
									SrcPos: 12,
								},
							},
						},
						data.Arrow{
							Data:   []data.Type{data.Type{LocalType: "i", SrcPos: 21}},
							ToPort: &data.Port{Name: "out", SrcPos: 26},
							SrcPos: 20,
						},
					},
				},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "complex 3",
			givenContent: "[A[p1|p2|p3]] (b)->c [D[plug=sp1,sp2,sp3]] (e)-> f[G[sp1,sp2]] -> h;",
			expectedValue: data.Flow{
				Parts: [][]interface{}{
					{
						data.Component{
							Decl: data.CompDecl{
								Name:   "a",
								Type:   data.Type{LocalType: "A", SrcPos: 1},
								SrcPos: 1,
							},
							Plugins: []data.NameNTypes{
								data.NameNTypes{
									Types: []data.Type{
										data.Type{LocalType: "p1", SrcPos: 3},
									},
									SrcPos: 3,
								},
								data.NameNTypes{
									Types: []data.Type{
										data.Type{LocalType: "p2", SrcPos: 6},
									},
									SrcPos: 6,
								},
								data.NameNTypes{
									Types: []data.Type{
										data.Type{LocalType: "p3", SrcPos: 9},
									},
									SrcPos: 9,
								},
							},
						},
						data.Arrow{
							Data:   []data.Type{data.Type{LocalType: "b", SrcPos: 15}},
							ToPort: &data.Port{Name: "c", SrcPos: 19},
							SrcPos: 14,
						},
						data.Component{
							Decl: data.CompDecl{
								Name:   "d",
								Type:   data.Type{LocalType: "D", SrcPos: 22},
								SrcPos: 22,
							},
							Plugins: []data.NameNTypes{
								data.NameNTypes{
									Name: "plug",
									Types: []data.Type{
										data.Type{LocalType: "sp1", SrcPos: 29},
										data.Type{LocalType: "sp2", SrcPos: 33},
										data.Type{LocalType: "sp3", SrcPos: 37},
									},
									SrcPos: 24,
								},
							},
							SrcPos: 21,
						},
						data.Arrow{
							Data:   []data.Type{data.Type{LocalType: "e", SrcPos: 44}},
							ToPort: &data.Port{Name: "f", SrcPos: 49},
							SrcPos: 43,
						},
						data.Component{
							Decl: data.CompDecl{
								Name:   "g",
								Type:   data.Type{LocalType: "G", SrcPos: 51},
								SrcPos: 51,
							},
							Plugins: []data.NameNTypes{
								data.NameNTypes{
									Types: []data.Type{
										data.Type{LocalType: "sp1", SrcPos: 53},
										data.Type{LocalType: "sp2", SrcPos: 57},
									},
									SrcPos: 53,
								},
							},
							SrcPos: 50,
						},
						data.Arrow{
							ToPort: &data.Port{Name: "h", SrcPos: 66},
							SrcPos: 63,
						},
					},
				},
			},
			expectedErrCount: 0,
		},
	})
}
