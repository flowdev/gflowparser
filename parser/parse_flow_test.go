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
				ToPort:   &data.Port{Name: "bPort"},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 3",
			givenContent: `(Data)->`,
			expectedValue: data.Arrow{
				Data: []data.Type{data.Type{LocalType: "Data"}},
			},
			expectedErrCount: 0,
		}, {
			givenName:    "complex",
			givenContent: "aPort // comment1\n ( \t Data \t ) \t -> // comment2\n bPort",
			expectedValue: data.Arrow{
				FromPort: &data.Port{Name: "aPort"},
				Data:     []data.Type{data.Type{LocalType: "Data"}},
				ToPort:   &data.Port{Name: "bPort"},
			},
			expectedErrCount: 0,
		},
	})
}
