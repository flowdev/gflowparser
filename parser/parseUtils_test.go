package parser

import (
	"reflect"
	"testing"

	"github.com/flowdev/gflowparser/data"
	"github.com/flowdev/gparselib"
)

// testParseOp is the interface of all parsers to be tested.
type testParseOp func(outPort func(interface{})) (inPort func(interface{}))

type parseTestData struct {
	givenName        string
	givenContent     string
	expectedValue    interface{}
	expectedErrCount int
}

func TestParseSmallIdent(t *testing.T) {
	runTests(t, ParseSmallIdent, []parseTestData{
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

func TestParseBigIdent(t *testing.T) {
	runTests(t, ParseBigIdent, []parseTestData{
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
			givenContent:     `A`,
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:        "simple 1",
			givenContent:     `Ab`,
			expectedValue:    "Ab",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 2",
			givenContent:     `A0`,
			expectedValue:    "A0",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 3",
			givenContent:     `Ab_cd`,
			expectedValue:    "Ab",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 4",
			givenContent:     `AbcDef`,
			expectedValue:    "AbcDef",
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
			expectedValue:    "",
			expectedErrCount: 0,
		}, {
			givenName:        "no match",
			givenContent:     `baaa`,
			expectedValue:    "",
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

func TestParseSpaceComment(t *testing.T) {
	runTests(t, ParseSpaceComment, []parseTestData{
		{
			givenName:        "empty",
			givenContent:     ``,
			expectedValue:    "",
			expectedErrCount: 0,
		}, {
			givenName:        "no match",
			givenContent:     `baaa`,
			expectedValue:    "",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 1",
			givenContent:     " i",
			expectedValue:    " ",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 2",
			givenContent:     "\t0",
			expectedValue:    "\t",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 3",
			givenContent:     " /* bla */ _t",
			expectedValue:    " /* bla */ ",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 4",
			givenContent:     " // comment! \n lilalo",
			expectedValue:    " // comment! \n ",
			expectedErrCount: 0,
		}, {
			givenName:        "complex",
			givenContent:     " /* bla\n */ \t // com!\n \t \r\n/** blu */ _t",
			expectedValue:    " /* bla\n */ \t // com!\n \t \r\n/** blu */ ",
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
			givenContent:     " /* bla\n */ \t // com!\n \t \r\n/** blu ; */ ",
			expectedValue:    nil,
			expectedErrCount: 1,
		}, {
			givenName:        "simple 1",
			givenContent:     ";",
			expectedValue:    ";",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 2",
			givenContent:     "\t;0",
			expectedValue:    "\t;",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 3",
			givenContent:     " /* bla */; _t",
			expectedValue:    " /* bla */; ",
			expectedErrCount: 0,
		}, {
			givenName:        "simple 4",
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
	var mainData2 *data.MainData
	portIn := p(func(dat interface{}) { mainData2 = dat.(*data.MainData) })
	for _, spec := range specs {
		t.Logf("Parsing source '%s'.", spec.givenName)
		mainData := &data.MainData{}
		mainData.ParseData = gparselib.NewParseData(spec.givenName, spec.givenContent)
		portIn(mainData)
		pd2 := mainData2.ParseData

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
			}
			if pd2.Result.Feedback[spec.expectedErrCount-1].Msg.String() == "" {
				t.Logf("Actual errors are: %s", printErrors(pd2.Result.Feedback))
				t.Errorf("Expected an error message.")
			}
		}

	}
}
func runTest(t *testing.T, fp interface{}, name string, content string, ev interface{}, errCount int) {
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
