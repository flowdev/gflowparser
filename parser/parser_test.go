package parser

import (
	"bytes"
	"testing"

	"github.com/flowdev/gparselib"
)

func TestCheckFeedback(t *testing.T) {
	errPrefix := "Found errors while parsing flow:\n"
	specs := []struct {
		name        string
		givenResult *gparselib.ParseResult
		// Feedback: []*FeedbackItem
		expectedFeedback string
		expectedError    string
	}{
		{
			name: "nil feedback",
			givenResult: &gparselib.ParseResult{
				Feedback: nil,
			},
			expectedFeedback: "",
			expectedError:    "",
		}, {
			name: "empty feedback",
			givenResult: &gparselib.ParseResult{
				Feedback: []*gparselib.FeedbackItem{},
			},
			expectedFeedback: "",
			expectedError:    "",
		}, {
			name: "one error",
			givenResult: &gparselib.ParseResult{
				Feedback: []*gparselib.FeedbackItem{
					&gparselib.FeedbackItem{
						Kind: gparselib.FeedbackError,
						Msg:  bytes.NewBufferString("No flow found!"),
					},
				},
			},
			expectedFeedback: "",
			expectedError:    errPrefix + "ERROR: No flow found!\n",
		}, {
			name: "two errors",
			givenResult: &gparselib.ParseResult{
				Feedback: []*gparselib.FeedbackItem{
					&gparselib.FeedbackItem{
						Kind: gparselib.FeedbackError,
						Msg:  bytes.NewBufferString("No flow found!"),
					},
					&gparselib.FeedbackItem{
						Kind: gparselib.FeedbackError,
						Msg:  bytes.NewBufferString("Nothing found!"),
					},
				},
			},
			expectedFeedback: "",
			expectedError: errPrefix + "ERROR: No flow found!\n" +
				"ERROR: Nothing found!\n",
		},
	}
	for _, spec := range specs {
		t.Logf("Testing spec: %s\n", spec.name)
		gotFeedback, gotError := CheckFeedback(spec.givenResult)

		if spec.expectedError != "" && gotError != nil {
			if spec.expectedError != gotError.Error() {
				t.Errorf("Expected error '%s' but got: '%s'",
					spec.expectedError, gotError)
			}
		} else if spec.expectedError != "" && gotError == nil {
			t.Error("Expected an error but didn't get one.")
		} else if spec.expectedError == "" && gotError != nil {
			t.Errorf("Expected no error but got: %s", gotError)
		}
		if spec.expectedFeedback != gotFeedback {
			t.Errorf("Expected feedback '%s' but got: '%s'",
				spec.expectedFeedback, gotFeedback)
		}
	}
}
