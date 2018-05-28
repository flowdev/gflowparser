package parser

import (
	"testing"

	"github.com/flowdev/gflowparser/data"
)

func TestParseType(t *testing.T) {
	p, err := NewParseType()
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
			givenContent:     `1A`,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:        "no match 2",
			givenContent:     `_A`,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:        "simple 1",
			givenContent:     `Ab`,
			expectedValue:    data.Type{Package: "", LocalType: "Ab"},
			expectedErrCount: 0,
		}, {
			givenName:        "simple 2",
			givenContent:     `a0`,
			expectedValue:    data.Type{Package: "", LocalType: "a0"},
			expectedErrCount: 0,
		}, {
			givenName:        "simple 3",
			givenContent:     `Ab_cd`,
			expectedValue:    data.Type{Package: "", LocalType: "Ab"},
			expectedErrCount: 0,
		}, {
			givenName:        "simple 4",
			givenContent:     `abcDef`,
			expectedValue:    data.Type{Package: "", LocalType: "abcDef"},
			expectedErrCount: 0,
		}, {
			givenName:        "complex 1",
			givenContent:     `p.Ab1Cd`,
			expectedValue:    data.Type{Package: "p", LocalType: "Ab1Cd"},
			expectedErrCount: 0,
		}, {
			givenName:        "complex 2",
			givenContent:     `pack.a1Bc_d`,
			expectedValue:    data.Type{Package: "pack", LocalType: "a1Bc"},
			expectedErrCount: 0,
		},
	})
}

func TestParseCompDecl(t *testing.T) {
	p, err := NewParseCompDecl()
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
			givenContent:     `1A`,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:        "no match 2",
			givenContent:     `_A`,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:    "simple 1",
			givenContent: `A`,
			expectedValue: data.CompDecl{
				Name: "a",
				Type: data.Type{LocalType: "A"},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 2",
			givenContent: `a0`,
			expectedValue: data.CompDecl{
				Name:      "a0",
				Type:      data.Type{LocalType: "a0"},
				VagueType: true,
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 3",
			givenContent: `p.Ab_cd`,
			expectedValue: data.CompDecl{
				Name: "ab",
				Type: data.Type{Package: "p", LocalType: "Ab"},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 4",
			givenContent: `abcDef`,
			expectedValue: data.CompDecl{
				Name:      "abcDef",
				Type:      data.Type{LocalType: "abcDef"},
				VagueType: true,
			},
			expectedErrCount: 0,
		}, {
			givenName:    "complex 1",
			givenContent: `n p.Ab1Cd`,
			expectedValue: data.CompDecl{
				Name: "n",
				Type: data.Type{Package: "p", LocalType: "Ab1Cd", SrcPos: 2},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "complex 2",
			givenContent: "nam \t pack.a1Bc_d",
			expectedValue: data.CompDecl{
				Name: "nam",
				Type: data.Type{Package: "pack", LocalType: "a1Bc", SrcPos: 6},
			},
			expectedErrCount: 0,
		},
	})
}

func TestParseTypeList(t *testing.T) {
	p, err := NewParseTypeList()
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
			givenContent:     `1A`,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:        "no match 2",
			givenContent:     `_A`,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:    "simple 1",
			givenContent: `A`,
			expectedValue: []data.Type{
				data.Type{LocalType: "A"},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 2",
			givenContent: `a,b`,
			expectedValue: []data.Type{
				data.Type{LocalType: "a"},
				data.Type{LocalType: "b", SrcPos: 2},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 3",
			givenContent: `p.A , q.B`,
			expectedValue: []data.Type{
				data.Type{Package: "p", LocalType: "A"},
				data.Type{Package: "q", LocalType: "B", SrcPos: 6},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "complex",
			givenContent: "a, B \t \n, /* comment */ p.C, q.D",
			expectedValue: []data.Type{
				data.Type{LocalType: "a"},
				data.Type{LocalType: "B", SrcPos: 3},
				data.Type{Package: "p", LocalType: "C", SrcPos: 24},
				data.Type{Package: "q", LocalType: "D", SrcPos: 29},
			},
			expectedErrCount: 0,
		},
	})
}

func TestParseTitledTypes(t *testing.T) {
	p, err := NewParseTitledTypes()
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
			givenContent:     `1A=bla`,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:        "no match 2",
			givenContent:     `a:b`,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:        "no match 3",
			givenContent:     `a=1a`,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:    "simple 1",
			givenContent: `a=A`,
			expectedValue: data.NameNTypes{
				Name: "a",
				Types: []data.Type{
					data.Type{LocalType: "A", SrcPos: 2},
				},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 2",
			givenContent: `a=b, C`,
			expectedValue: data.NameNTypes{
				Name: "a",
				Types: []data.Type{
					data.Type{LocalType: "b", SrcPos: 2},
					data.Type{LocalType: "C", SrcPos: 5},
				},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 3",
			givenContent: `tiTle = p.A, q.B`,
			expectedValue: data.NameNTypes{
				Name: "tiTle",
				Types: []data.Type{
					data.Type{Package: "p", LocalType: "A", SrcPos: 8},
					data.Type{Package: "q", LocalType: "B", SrcPos: 13},
				},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "complex",
			givenContent: "t \t \n= /* comment */ p.C, q.D",
			expectedValue: data.NameNTypes{
				Name: "t",
				Types: []data.Type{
					data.Type{Package: "p", LocalType: "C", SrcPos: 21},
					data.Type{Package: "q", LocalType: "D", SrcPos: 26},
				},
			},
			expectedErrCount: 0,
		},
	})
}

func TestParseTitledTypesList(t *testing.T) {
	p, err := NewParseTitledTypesList()
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
			givenName:        "no match",
			givenContent:     `|a=b`,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:    "simple 1",
			givenContent: `a=A`,
			expectedValue: []data.NameNTypes{
				data.NameNTypes{
					Name:  "a",
					Types: []data.Type{data.Type{LocalType: "A", SrcPos: 2}},
				},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 2",
			givenContent: `a=b|c=D`,
			expectedValue: []data.NameNTypes{
				data.NameNTypes{
					Name:  "a",
					Types: []data.Type{data.Type{LocalType: "b", SrcPos: 2}},
				},
				data.NameNTypes{
					Name:   "c",
					Types:  []data.Type{data.Type{LocalType: "D", SrcPos: 6}},
					SrcPos: 4,
				},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 3",
			givenContent: `a=b | c=D`,
			expectedValue: []data.NameNTypes{
				data.NameNTypes{
					Name:  "a",
					Types: []data.Type{data.Type{LocalType: "b", SrcPos: 2}},
				},
				data.NameNTypes{
					Name:   "c",
					Types:  []data.Type{data.Type{LocalType: "D", SrcPos: 8}},
					SrcPos: 6,
				},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "complex",
			givenContent: "a=b \t \n| /* comment */ c=D",
			expectedValue: []data.NameNTypes{
				data.NameNTypes{
					Name:  "a",
					Types: []data.Type{data.Type{LocalType: "b", SrcPos: 2}},
				},
				data.NameNTypes{
					Name:   "c",
					Types:  []data.Type{data.Type{LocalType: "D", SrcPos: 25}},
					SrcPos: 23,
				},
			},
			expectedErrCount: 0,
		},
	})
}

func TestParsePlugins(t *testing.T) {
	p, err := NewParsePlugins()
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
			givenName:        "no match",
			givenContent:     `[a=b`,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:    "simple 1",
			givenContent: `[a=A]`,
			expectedValue: []data.NameNTypes{
				data.NameNTypes{
					Name:   "a",
					Types:  []data.Type{data.Type{LocalType: "A", SrcPos: 3}},
					SrcPos: 1,
				},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 2",
			givenContent: `[a]`,
			expectedValue: []data.NameNTypes{
				data.NameNTypes{
					Name:   "",
					Types:  []data.Type{data.Type{LocalType: "a", SrcPos: 1}},
					SrcPos: 1,
				},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 3",
			givenContent: `[ a=b,D ]`,
			expectedValue: []data.NameNTypes{
				data.NameNTypes{
					Name: "a",
					Types: []data.Type{
						data.Type{LocalType: "b", SrcPos: 4},
						data.Type{LocalType: "D", SrcPos: 6},
					},
					SrcPos: 2,
				},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "complex",
			givenContent: "[ \t \na=b|c=D /* comment */ ]",
			expectedValue: []data.NameNTypes{
				data.NameNTypes{
					Name:   "a",
					Types:  []data.Type{data.Type{LocalType: "b", SrcPos: 7}},
					SrcPos: 5,
				},
				data.NameNTypes{
					Name:   "c",
					Types:  []data.Type{data.Type{LocalType: "D", SrcPos: 11}},
					SrcPos: 9,
				},
			},
			expectedErrCount: 0,
		},
	})
}

func TestParseComponent(t *testing.T) {
	p, err := NewParseComponent()
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
			givenName:        "no match",
			givenContent:     `[a [a]]`,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:    "simple 1",
			givenContent: `[a A]`,
			expectedValue: data.Component{
				Decl: data.CompDecl{
					Name: "a", Type: data.Type{LocalType: "A", SrcPos: 3},
					SrcPos: 1,
				},
				Plugins: nil,
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 2",
			givenContent: `[a B[c]]`,
			expectedValue: data.Component{
				Decl: data.CompDecl{
					Name: "a", Type: data.Type{LocalType: "B", SrcPos: 3},
					SrcPos: 1,
				},
				Plugins: []data.NameNTypes{
					data.NameNTypes{
						Name:   "",
						Types:  []data.Type{data.Type{LocalType: "c", SrcPos: 5}},
						SrcPos: 5,
					},
				},
			},
		}, {
			givenName:    "simple 3",
			givenContent: `[ a B [c=D] ]`,
			expectedValue: data.Component{
				Decl: data.CompDecl{
					Name: "a", Type: data.Type{LocalType: "B", SrcPos: 4},
					SrcPos: 2,
				},
				Plugins: []data.NameNTypes{
					data.NameNTypes{
						Name:   "c",
						Types:  []data.Type{data.Type{LocalType: "D", SrcPos: 9}},
						SrcPos: 7,
					},
				},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "complex",
			givenContent: "[ \t \na B /* comment 1 */ [c=D] // comment 2\n ]",
			expectedValue: data.Component{
				Decl: data.CompDecl{
					Name: "a", Type: data.Type{LocalType: "B", SrcPos: 7},
					SrcPos: 5,
				},
				Plugins: []data.NameNTypes{
					data.NameNTypes{
						Name:   "c",
						Types:  []data.Type{data.Type{LocalType: "D", SrcPos: 28}},
						SrcPos: 26,
					},
				},
			},
			expectedErrCount: 0,
		},
	})
}