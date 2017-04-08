package output

import (
	"github.com/flowdev/gflowparser/data"
)

var specialFormat = map[string]bool{
	"dot": true,
}

// ------------ ProduceFormats:
// input:  *data.MainData{SelectedFormats, FlowFile}
// output: *data.MainData{SelectedFormats, FlowFile, CurrentFormat} to the correct output port.
type ProduceFormats struct {
	outPort        func(*data.MainData)
	specialOutPort func(*data.MainData)
}

func NewOutputFormats() *ProduceFormats {
	return &ProduceFormats{}
}
func (op *ProduceFormats) InPort(md *data.MainData) {
	for _, format := range md.SelectedFormats {
		md.CurrentFormat = format
		if specialFormat[format] {
			op.specialOutPort(md)
		} else {
			op.outPort(md)
		}
	}
}
func (op *ProduceFormats) SetOutPort(port func(*data.MainData)) {
	op.outPort = port
}
func (op *ProduceFormats) SetSpecialOutPort(port func(*data.MainData)) {
	op.specialOutPort = port
}
