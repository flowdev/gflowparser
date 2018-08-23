package parser

import (
	"bytes"
	"errors"

	"github.com/flowdev/gparselib"
)

// CheckFeedback converts parser errors into a single error and
// additional feedback.
func CheckFeedback(pr *gparselib.ParseResult) (string, error) {
	if pr.HasError() {
		return "", errors.New("Found errors while parsing flow:\n" +
			feedbackToString(pr))
	}
	return feedbackToString(pr), nil
}
func feedbackToString(pr *gparselib.ParseResult) string {
	buf := bytes.Buffer{}
	for _, fb := range pr.Feedback {
		buf.WriteString(fb.String())
		buf.WriteString("\n")
	}
	return buf.String()
}
