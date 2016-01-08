package parser

import (
	"github.com/flowdev/gflowparser/data"
	"testing"
)

func TestParseChainMiddle(t *testing.T) {
	p := NewParseChainMiddle()

	runTest(t, p, "no match 1", "->(B)", nil, 4)
	runTest(t, p, "no match 2", "(Bla)", nil, 3)
	runTest(t, p, "simple 1", "->(Bla)",
		[]interface{}{"", &data.Operation{Name: "bla", Type: "Bla",
			InPorts:  []*data.PortData{&data.PortData{Name: "in", CapName: "In", SrcPos: 2}},
			OutPorts: []*data.PortData{&data.PortData{Name: "out", CapName: "Out", SrcPos: 7}}}}, 0)
	runTest(t, p, "simple 2", "-> in.2 \t bla() \t error ",
		[]interface{}{"", &data.Operation{Name: "bla",
			InPorts:  []*data.PortData{&data.PortData{Name: "in", CapName: "In", HasIndex: true, Index: 2, SrcPos: 3}},
			OutPorts: []*data.PortData{&data.PortData{Name: "error", CapName: "Error", SrcPos: 18}}}}, 0)
	runTest(t, p, "simple 3", "[DataType]-> /* comm */\n \t xIn.1   bla(Blu)outY.123",
		[]interface{}{"DataType", &data.Operation{Name: "bla", Type: "Blu",
			InPorts:  []*data.PortData{&data.PortData{Name: "xIn", CapName: "XIn", HasIndex: true, Index: 1, SrcPos: 27}},
			OutPorts: []*data.PortData{&data.PortData{Name: "outY", CapName: "OutY", HasIndex: true, Index: 123, SrcPos: 43}}}}, 0)
}

func TestParseChainEnd(t *testing.T) {
	p := NewParseChainEnd()

	runTest(t, p, "empty", "", nil, 3)
	runTest(t, p, "no match 1", "-", nil, 3)
	runTest(t, p, "no match 3", " /* \n */ \t [Bla]>", nil, 3)
	runTest(t, p, "no ports, no type", "->",
		&data.Connection{FromPort: data.NewPort("", 0), ToPort: data.NewPort("", 2)}, 0)
	runTest(t, p, "no ports but a type", " \t [Bla]-> ",
		&data.Connection{DataType: "Bla", FromPort: data.NewPort("", 0), ToPort: data.NewPort("", 11)}, 0)
	runTest(t, p, "out-port and no type", " \r\n // blu \n \t -> \r\n outX.3",
		&data.Connection{FromPort: data.NewPort("", 0), ToPort: data.NewIdxPort("outX", 3, 21)}, 0)
	runTest(t, p, "out-port and type", "\n \t /* Bla */ [ \t Blu \t ]->  \t outX.7",
		&data.Connection{DataType: "Blu", FromPort: data.NewPort("", 0), ToPort: data.NewIdxPort("outX", 7, 31)}, 0)
}
