package parser

import (
	"fmt"
	"io"
	"os"

	"github.com/flowdev/gflowparser/data"
	"github.com/flowdev/gparselib"
)

// ------------ HandleParseErrors:
// output: MainData with FlowFile filled if no errors were found or none otherwise
type HandleParseResult struct {
	outPort func(interface{})
}

func NewHandleParseResult() *HandleParseResult {
	return &HandleParseResult{}
}
func (op *HandleParseResult) InPort(dat interface{}) {
	md := dat.(*data.MainData)
	r := md.ParseData.Result

	if !r.HasError() {
		outputFeedback(os.Stdout, r.Feedback)
		md.FlowFile = r.Value.(*data.FlowFile)
		op.outPort(md)
	} else {
		outputFeedback(os.Stdout, r.Feedback)
	}
}
func (op *HandleParseResult) SetOutPort(port func(interface{})) {
	op.outPort = port
}
func outputFeedback(w io.Writer, msgs []*gparselib.FeedbackItem) {
	for msg := range msgs {
		fmt.Fprintln(w, msg)
	}
}
