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
