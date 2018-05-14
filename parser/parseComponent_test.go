package parser

import "testing"

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
			expectedValue:    &TypeSemValue{Package: "", LocalType: "Ab"},
			expectedErrCount: 0,
		}, {
			givenName:        "simple 2",
			givenContent:     `a0`,
			expectedValue:    &TypeSemValue{Package: "", LocalType: "a0"},
			expectedErrCount: 0,
		}, {
			givenName:        "simple 3",
			givenContent:     `Ab_cd`,
			expectedValue:    &TypeSemValue{Package: "", LocalType: "Ab"},
			expectedErrCount: 0,
		}, {
			givenName:        "simple 4",
			givenContent:     `abcDef`,
			expectedValue:    &TypeSemValue{Package: "", LocalType: "abcDef"},
			expectedErrCount: 0,
		}, {
			givenName:        "complex 1",
			givenContent:     `p.Ab1Cd`,
			expectedValue:    &TypeSemValue{Package: "p", LocalType: "Ab1Cd"},
			expectedErrCount: 0,
		}, {
			givenName:        "complex 2",
			givenContent:     `pack.a1Bc_d`,
			expectedValue:    &TypeSemValue{Package: "pack", LocalType: "a1Bc"},
			expectedErrCount: 0,
		},
	})
}

func TestParseOpDecl(t *testing.T) {
	p, err := NewParseOpDecl()
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
			expectedValue: &OpDeclSemValue{
				Name: "a",
				Type: &TypeSemValue{LocalType: "A"},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 2",
			givenContent: `a0`,
			expectedValue: &OpDeclSemValue{
				Name: "a0",
				Type: &TypeSemValue{LocalType: "a0"},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 3",
			givenContent: `p.Ab_cd`,
			expectedValue: &OpDeclSemValue{
				Name: "ab",
				Type: &TypeSemValue{Package: "p", LocalType: "Ab"},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 4",
			givenContent: `abcDef`,
			expectedValue: &OpDeclSemValue{
				Name: "abcDef",
				Type: &TypeSemValue{LocalType: "abcDef"},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "complex 1",
			givenContent: `n p.Ab1Cd`,
			expectedValue: &OpDeclSemValue{
				Name: "n",
				Type: &TypeSemValue{Package: "p", LocalType: "Ab1Cd"},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "complex 2",
			givenContent: "nam \t pack.a1Bc_d",
			expectedValue: &OpDeclSemValue{
				Name: "nam",
				Type: &TypeSemValue{Package: "pack", LocalType: "a1Bc"},
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
			expectedValue: []*TypeSemValue{
				&TypeSemValue{LocalType: "A"},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 2",
			givenContent: `a,b`,
			expectedValue: []*TypeSemValue{
				&TypeSemValue{LocalType: "a"},
				&TypeSemValue{LocalType: "b"},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 3",
			givenContent: `p.A , q.B`,
			expectedValue: []*TypeSemValue{
				&TypeSemValue{Package: "p", LocalType: "A"},
				&TypeSemValue{Package: "q", LocalType: "B"},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "complex",
			givenContent: "a, B \t \n, /* comment */ p.C, q.D",
			expectedValue: []*TypeSemValue{
				&TypeSemValue{LocalType: "a"},
				&TypeSemValue{LocalType: "B"},
				&TypeSemValue{Package: "p", LocalType: "C"},
				&TypeSemValue{Package: "q", LocalType: "D"},
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
			expectedValue: &TitledTypesSemValue{
				Title: "a",
				Types: []*TypeSemValue{
					&TypeSemValue{LocalType: "A"},
				},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 2",
			givenContent: `a=b, C`,
			expectedValue: &TitledTypesSemValue{
				Title: "a",
				Types: []*TypeSemValue{
					&TypeSemValue{LocalType: "b"},
					&TypeSemValue{LocalType: "C"},
				},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 3",
			givenContent: `tiTle = p.A, q.B`,
			expectedValue: &TitledTypesSemValue{
				Title: "tiTle",
				Types: []*TypeSemValue{
					&TypeSemValue{Package: "p", LocalType: "A"},
					&TypeSemValue{Package: "q", LocalType: "B"},
				},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "complex",
			givenContent: "t \t \n= /* comment */ p.C, q.D",
			expectedValue: &TitledTypesSemValue{
				Title: "t",
				Types: []*TypeSemValue{
					&TypeSemValue{Package: "p", LocalType: "C"},
					&TypeSemValue{Package: "q", LocalType: "D"},
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
			expectedValue: []*TitledTypesSemValue{
				&TitledTypesSemValue{
					Title: "a",
					Types: []*TypeSemValue{&TypeSemValue{LocalType: "A"}},
				},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 2",
			givenContent: `a=b|c=D`,
			expectedValue: []*TitledTypesSemValue{
				&TitledTypesSemValue{
					Title: "a",
					Types: []*TypeSemValue{&TypeSemValue{LocalType: "b"}},
				},
				&TitledTypesSemValue{
					Title: "c",
					Types: []*TypeSemValue{&TypeSemValue{LocalType: "D"}},
				},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 3",
			givenContent: `a=b | c=D`,
			expectedValue: []*TitledTypesSemValue{
				&TitledTypesSemValue{
					Title: "a",
					Types: []*TypeSemValue{&TypeSemValue{LocalType: "b"}},
				},
				&TitledTypesSemValue{
					Title: "c",
					Types: []*TypeSemValue{&TypeSemValue{LocalType: "D"}},
				},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "complex",
			givenContent: "a=b \t \n| /* comment */ c=D",
			expectedValue: []*TitledTypesSemValue{
				&TitledTypesSemValue{
					Title: "a",
					Types: []*TypeSemValue{&TypeSemValue{LocalType: "b"}},
				},
				&TitledTypesSemValue{
					Title: "c",
					Types: []*TypeSemValue{&TypeSemValue{LocalType: "D"}},
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
			expectedValue: []*TitledTypesSemValue{
				&TitledTypesSemValue{
					Title: "a",
					Types: []*TypeSemValue{&TypeSemValue{LocalType: "A"}},
				},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 2",
			givenContent: `[a]`,
			expectedValue: []*TitledTypesSemValue{
				&TitledTypesSemValue{
					Title: "",
					Types: []*TypeSemValue{&TypeSemValue{LocalType: "a"}},
				},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 3",
			givenContent: `[ a=b,D ]`,
			expectedValue: []*TitledTypesSemValue{
				&TitledTypesSemValue{
					Title: "a",
					Types: []*TypeSemValue{
						&TypeSemValue{LocalType: "b"},
						&TypeSemValue{LocalType: "D"},
					},
				},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "complex",
			givenContent: "[ \t \na=b|c=D /* comment */ ]",
			expectedValue: []*TitledTypesSemValue{
				&TitledTypesSemValue{
					Title: "a",
					Types: []*TypeSemValue{&TypeSemValue{LocalType: "b"}},
				},
				&TitledTypesSemValue{
					Title: "c",
					Types: []*TypeSemValue{&TypeSemValue{LocalType: "D"}},
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
			expectedValue: &ComponentSemValue{
				Decl: &OpDeclSemValue{
					Name: "a", Type: &TypeSemValue{LocalType: "A"},
				},
				Plugins: nil,
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 2",
			givenContent: `[a B[c]]`,
			expectedValue: &ComponentSemValue{
				Decl: &OpDeclSemValue{
					Name: "a", Type: &TypeSemValue{LocalType: "B"},
				},
				Plugins: []*TitledTypesSemValue{
					&TitledTypesSemValue{
						Title: "",
						Types: []*TypeSemValue{&TypeSemValue{LocalType: "c"}},
					},
				},
			},
		}, {
			givenName:    "simple 3",
			givenContent: `[ a B [c=D] ]`,
			expectedValue: &ComponentSemValue{
				Decl: &OpDeclSemValue{
					Name: "a", Type: &TypeSemValue{LocalType: "B"},
				},
				Plugins: []*TitledTypesSemValue{
					&TitledTypesSemValue{
						Title: "c",
						Types: []*TypeSemValue{&TypeSemValue{LocalType: "D"}},
					},
				},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "complex",
			givenContent: "[ \t \na B /* comment 1 */ [c=D] // comment 2\n ]",
			expectedValue: &ComponentSemValue{
				Decl: &OpDeclSemValue{
					Name: "a", Type: &TypeSemValue{LocalType: "B"},
				},
				Plugins: []*TitledTypesSemValue{
					&TitledTypesSemValue{
						Title: "c",
						Types: []*TypeSemValue{&TypeSemValue{LocalType: "D"}},
					},
				},
			},
			expectedErrCount: 0,
		},
	})
}
