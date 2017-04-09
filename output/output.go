package output

import (
	"github.com/flowdev/gflowparser/data"
)

var AllowedFormat = map[string]bool{
	"dot": true,
	"go":  true,
}
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

// ------------ FillPortPairs:
// input:  *data.MainData{SelectedFormats, FlowFile}
// output: *data.MainData{SelectedFormats, FlowFile, CurrentFormat} to the correct output port.
type FillPortPairs struct {
	outPort func(*data.MainData)
}

func NewFillPortPairs() *FillPortPairs {
	return &FillPortPairs{}
}
func (op *FillPortPairs) InPort(md *data.MainData) {
	for _, flow := range md.FlowFile.Flows {
		for _, op := range flow.Ops {
			fillPortPairs4Op(op)
		}
	}
	op.outPort(md)
}
func (op *FillPortPairs) SetOutPort(port func(*data.MainData)) {
	op.outPort = port
}

func fillPortPairs4Op(op *data.Operation) {
	l := len(op.InPorts)
	m := len(op.OutPorts)
	n := max(l, m)
	portPairs := make([]*data.PortPair, n)
	for i := 0; i < n; i++ {
		p := &data.PortPair{}
		if i < l {
			p.InPort = op.InPorts[i]
		}
		if i < m {
			p.OutPort = op.OutPorts[i]
		}
		portPairs[i] = p
	}
}
func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
