package parser

import (
	"reflect"
	"testing"

	"github.com/flowdev/gparselib"
)

// testParseOp is the interface of all parsers to be tested.
type testParseOp func(*gparselib.ParseData, interface{}) (*gparselib.ParseData, interface{})

type parseTestData struct {
	givenName        string
	givenContent     string
	expectedValue    interface{}
	expectedErrCount int
}

func TestParseNameIdent(t *testing.T) {
	p, err := NewNameIdentParser()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	runTests(t, p.ParseNameIdent, []parseTestData{
		{
			givenName:        "empty",
			givenContent:     ``,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:        "no match 1",
			givenContent:     `ABCD`,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:        "no match 2",
			givenContent:     `123`,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:        "simple 1",
			givenContent:     `aB`,
			expectedValue:    "aB",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 2",
			givenContent:     `a0`,
			expectedValue:    "a0",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 3",
			givenContent:     `aB_CD`,
			expectedValue:    "aB",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 4",
			givenContent:     `aBCdEF`,
			expectedValue:    "aBCdEF",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 5",
			givenContent:     `aBC123dEF`,
			expectedValue:    "aBC123dEF",
			expectedErrCount: 0,
		},
	})
}

func TestParsePackageIdent(t *testing.T) {
	p, err := NewPackageIdentParser()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	runTests(t, p.ParsePackageIdent, []parseTestData{
		{
			givenName:        "empty",
			givenContent:     ``,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:        "no match 1",
			givenContent:     `ABCD.`,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:        "no match 2",
			givenContent:     `123.`,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:        "simple 1",
			givenContent:     `a.`,
			expectedValue:    "a",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 2",
			givenContent:     `a0.`,
			expectedValue:    "a0",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 3",
			givenContent:     `ab._CD`,
			expectedValue:    "ab",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 4",
			givenContent:     `abcd.EF`,
			expectedValue:    "abcd",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 5",
			givenContent:     `abc123d.EF`,
			expectedValue:    "abc123d",
			expectedErrCount: 0,
		},
	})
}

func TestParseLocalTypeIdent(t *testing.T) {
	p, err := NewLocalTypeIdentParser()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	runTests(t, p.ParseLocalTypeIdent, []parseTestData{
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
			expectedValue:    "Ab",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 2",
			givenContent:     `a0`,
			expectedValue:    "a0",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 3",
			givenContent:     `Ab_cd`,
			expectedValue:    "Ab",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 4",
			givenContent:     `abcDef`,
			expectedValue:    "abcDef",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 5",
			givenContent:     `Abc123Def`,
			expectedValue:    "Abc123Def",
			expectedErrCount: 0,
		},
	})
}

func TestParseOptSpc(t *testing.T) {
	runTests(t, ParseOptSpc, []parseTestData{
		{
			givenName:        "empty",
			givenContent:     ``,
			expectedValue:    nil,
			expectedErrCount: 0,
		}, {
			givenName:        "no match",
			givenContent:     `baaa`,
			expectedValue:    nil,
			expectedErrCount: 0,
		}, {
			givenName:        "simple 1",
			givenContent:     ` i`,
			expectedValue:    " ",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 2",
			givenContent:     "\t0",
			expectedValue:    "\t",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 3",
			givenContent:     " \t _t",
			expectedValue:    " \t ",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 4",
			givenContent:     " \n ",
			expectedValue:    " ",
			expectedErrCount: 0,
		},
	})
}

func TestParseASpc(t *testing.T) {
	runTests(t, ParseASpc, []parseTestData{
		{
			givenName:        "empty",
			givenContent:     ``,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:        "no match",
			givenContent:     `baaa`,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:        "simple 1",
			givenContent:     ` i`,
			expectedValue:    " ",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 2",
			givenContent:     "\t0",
			expectedValue:    "\t",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 3",
			givenContent:     " \t _t",
			expectedValue:    " \t ",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 4",
			givenContent:     " \n ",
			expectedValue:    " ",
			expectedErrCount: 0,
		},
	})
}

func TestParseSpaceComment(t *testing.T) {
	runTests(t, ParseSpaceComment, []parseTestData{
		{
			givenName:        "empty",
			givenContent:     ``,
			expectedValue:    SpaceCommentSemValue{Text: "", NewLine: false},
			expectedErrCount: 0,
		}, {
			givenName:        "no match",
			givenContent:     `baaa`,
			expectedValue:    SpaceCommentSemValue{Text: "", NewLine: false},
			expectedErrCount: 0,
		}, {
			givenName:        "simple 1",
			givenContent:     " i",
			expectedValue:    SpaceCommentSemValue{Text: " ", NewLine: false},
			expectedErrCount: 0,
		}, {
			givenName:        "simple 2",
			givenContent:     "\t0",
			expectedValue:    SpaceCommentSemValue{Text: "\t", NewLine: false},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 3",
			givenContent: " /* bla */ _t",
			expectedValue: SpaceCommentSemValue{
				Text:    " /* bla */ ",
				NewLine: false,
			},
			expectedErrCount: 0,
		}, {
			givenName:    "simple 4",
			givenContent: " // comment! \n lilalo",
			expectedValue: SpaceCommentSemValue{
				Text:    " // comment! \n ",
				NewLine: true,
			},
			expectedErrCount: 0,
		}, {
			givenName:    "complex",
			givenContent: " /* bla\n */ \t // com!\n \t \r\n/** blu */ _t",
			expectedValue: SpaceCommentSemValue{
				Text:    " /* bla\n */ \t // com!\n \t \r\n/** blu */ ",
				NewLine: true,
			},
			expectedErrCount: 0,
		},
	})
}

func TestParseStatementEnd(t *testing.T) {
	runTests(t, ParseStatementEnd, []parseTestData{
		{
			givenName:        "empty",
			givenContent:     ``,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:        "no match 1",
			givenContent:     `baaa`,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:        "no match 2",
			givenContent:     " /* bla */ \t /** blu ; */ \t ",
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:        "simple 1",
			givenContent:     ";",
			expectedValue:    ";",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 2",
			givenContent:     "\n",
			expectedValue:    "\n",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 3",
			givenContent:     "\t;0",
			expectedValue:    "\t;",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 4",
			givenContent:     " /* bla */\r\n _t",
			expectedValue:    " /* bla */\r\n ",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 5",
			givenContent:     " // comment! \n ;lilalo",
			expectedValue:    " // comment! \n ;",
			expectedErrCount: 0,
		}, {
			givenName:        "complex",
			givenContent:     " /* bla\n */ \t; // com!\n \t \r\n/** blu */ _t",
			expectedValue:    " /* bla\n */ \t; // com!\n \t \r\n/** blu */ ",
			expectedErrCount: 0,
		},
	})
}

func runTests(t *testing.T, p testParseOp, specs []parseTestData) {
	for _, spec := range specs {
		t.Logf("Parsing source '%s'.", spec.givenName)
		pd := gparselib.NewParseData(spec.givenName, spec.givenContent)
		pd2, _ := p(pd, nil)

		if spec.expectedValue != nil && pd2.Result.Value == nil {
			t.Errorf("Expected a semantic result.")
		} else if !reflect.DeepEqual(pd2.Result.Value, spec.expectedValue) {
			t.Logf("The acutal value isn't equal to the expected one:")
			t.Logf(
				"Expected value of type '%T', got '%T'.",
				spec.expectedValue, pd2.Result.Value,
			)
			t.Errorf(
				"Expected value '%#v', got '%#v'.",
				spec.expectedValue, pd2.Result.Value,
			)
		}
		if spec.expectedErrCount <= 0 {
			if pd2.Result.HasError() {
				t.Logf("Actual errors are: %s", printErrors(pd2.Result.Feedback))
				t.Errorf("Expected no error but found at least one.")
			}
		} else {
			if len(pd2.Result.Feedback) != spec.expectedErrCount {
				t.Logf("Actual errors are: %s", printErrors(pd2.Result.Feedback))
				t.Fatalf(
					"Expected %d errors, got %d.",
					spec.expectedErrCount, len(pd2.Result.Feedback),
				)
			} else if pd2.Result.Feedback[spec.expectedErrCount-1].Msg.String() == "" {
				t.Logf("Actual errors are: %s", printErrors(pd2.Result.Feedback))
				t.Errorf("Expected an error message.")
			}
		}

	}
}
func printErrors(errs []*gparselib.FeedbackItem) string {
	result := ""
	for _, err := range errs {
		if err.Kind == gparselib.FeedbackError {
			result += err.String() + "\n"
		}
	}
	if result == "" {
		result = "<EMPTY>"
	}
	return result
}
