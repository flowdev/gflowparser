package parser

import (
	"github.com/flowdev/gflowparser/data"
	"testing"
)

func TestParseChainBegin(t *testing.T) {
	p := NewParseChainBegin()
	maxOpBlaNoPorts1 := &data.Operation{Name: "bla", Type: "Bla", SrcPos: 2,
		InPorts:  []*data.PortData{&data.PortData{Name: "in", CapName: "In", SrcPos: 2}},
		OutPorts: []*data.PortData{&data.PortData{Name: "out", CapName: "Out", SrcPos: 7}}}
	maxOpBlaNoPorts2 := &data.Operation{Name: "bla", Type: "Bla", SrcPos: 15,
		InPorts:  []*data.PortData{&data.PortData{Name: "in", CapName: "In", SrcPos: 15}},
		OutPorts: []*data.PortData{&data.PortData{Name: "out", CapName: "Out", SrcPos: 20}}}
	maxConnNoTypeNoPorts := &data.Connection{FromPort: &data.PortData{Name: "in", CapName: "In", SrcPos: 0},
		ToPort: maxOpBlaNoPorts1.InPorts[0], ToOp: maxOpBlaNoPorts1}
	maxConnTypeNoPorts := &data.Connection{FromPort: &data.PortData{Name: "in", CapName: "In", SrcPos: 0},
		ToPort: maxOpBlaNoPorts2.InPorts[0], ToOp: maxOpBlaNoPorts2, DataType: "BlaFlowData", ShowDataType: true}

	maxOpBluPorts := &data.Operation{Name: "bla", Type: "Blu", SrcPos: 30,
		InPorts:  []*data.PortData{&data.PortData{Name: "xIn", CapName: "XIn", SrcPos: 22, HasIndex: true, Index: 1}},
		OutPorts: []*data.PortData{&data.PortData{Name: "outY", CapName: "OutY", SrcPos: 38, HasIndex: true, Index: 123}}}
	maxConnTypePorts := &data.Connection{FromPort: &data.PortData{Name: "ourIn", CapName: "OurIn", SrcPos: 0},
		ToPort: maxOpBluPorts.InPorts[0], ToOp: maxOpBluPorts, DataType: "BlaFlowData", ShowDataType: true}

	minOpBlaNoPorts := &data.Operation{Name: "bla", Type: "Bla", SrcPos: 0,
		OutPorts: []*data.PortData{&data.PortData{Name: "out", CapName: "Out", SrcPos: 5}}}
	minOpBlaPorts := &data.Operation{Name: "bla", SrcPos: 0,
		OutPorts: []*data.PortData{&data.PortData{Name: "error", CapName: "Error", SrcPos: 6, HasIndex: true, Index: 3}}}
	minOpBluePorts := &data.Operation{Name: "bla", Type: "Blue", SrcPos: 0,
		OutPorts: []*data.PortData{&data.PortData{Name: "error", CapName: "Error", SrcPos: 10, HasIndex: true, Index: 3}}}

	runTest(t, p, "no match 1", "->(B)", nil, 3)
	runTest(t, p, "no match 2", "(B)", nil, 3)
	runTest(t, p, "simple max 1", "->(Bla)", []interface{}{maxConnNoTypeNoPorts, maxOpBlaNoPorts1}, 0)
	runTest(t, p, "simple max 2", "[BlaFlowData]->(Bla)", []interface{}{maxConnTypeNoPorts, maxOpBlaNoPorts2}, 0)
	runTest(t, p, "full max", "ourIn [BlaFlowData]-> xIn.1   bla(Blu)outY.123", []interface{}{maxConnTypePorts, maxOpBluPorts}, 0)
	runTest(t, p, "simple min 1", "(Bla)", []interface{}{nil, minOpBlaNoPorts}, 0)
	runTest(t, p, "simple min 2", "bla() error.3", []interface{}{nil, minOpBlaPorts}, 0)
	runTest(t, p, "full min", "bla(Blue) error.3", []interface{}{nil, minOpBluePorts}, 0)
}

func TestParseChainMiddle(t *testing.T) {
	p := NewParseChainMiddle()

	runTest(t, p, "no match 1", "->(B)", nil, 1)
	runTest(t, p, "no match 2", "(Bla)", nil, 1)
	runTest(t, p, "simple 1", "->(Bla)",
		[]interface{}{"", &data.Operation{Name: "bla", Type: "Bla", SrcPos: 2,
			InPorts:  []*data.PortData{&data.PortData{Name: "in", CapName: "In", SrcPos: 2}},
			OutPorts: []*data.PortData{&data.PortData{Name: "out", CapName: "Out", SrcPos: 7}}}}, 0)
	runTest(t, p, "simple 2", "-> in.2 \t bla() \t error ",
		[]interface{}{"", &data.Operation{Name: "bla", SrcPos: 10,
			InPorts:  []*data.PortData{&data.PortData{Name: "in", CapName: "In", HasIndex: true, Index: 2, SrcPos: 3}},
			OutPorts: []*data.PortData{&data.PortData{Name: "error", CapName: "Error", SrcPos: 18}}}}, 0)
	runTest(t, p, "simple 3", "[DataType]-> /* comm */\n \t xIn.1   bla(Blu)outY.123",
		[]interface{}{"DataType", &data.Operation{Name: "bla", Type: "Blu", SrcPos: 35,
			InPorts:  []*data.PortData{&data.PortData{Name: "xIn", CapName: "XIn", HasIndex: true, Index: 1, SrcPos: 27}},
			OutPorts: []*data.PortData{&data.PortData{Name: "outY", CapName: "OutY", HasIndex: true, Index: 123, SrcPos: 43}}}}, 0)
}

func TestParseChainEnd(t *testing.T) {
	p := NewParseChainEnd()

	runTest(t, p, "empty", "", nil, 1)
	runTest(t, p, "no match 1", "-", nil, 1)
	runTest(t, p, "no match 3", " /* \n */ \t [Bla]>", nil, 1)
	runTest(t, p, "no ports, no type", "->",
		&data.Connection{FromPort: data.NewPort("", 0), ToPort: data.NewPort("", 2)}, 0)
	runTest(t, p, "no ports but a type", " \t [Bla]-> ",
		&data.Connection{DataType: "Bla", FromPort: data.NewPort("", 0), ToPort: data.NewPort("", 11)}, 0)
	runTest(t, p, "out-port and no type", " \r\n // blu \n \t -> \r\n outX.3",
		&data.Connection{FromPort: data.NewPort("", 0), ToPort: data.NewIdxPort("outX", 3, 21)}, 0)
	runTest(t, p, "out-port and type", "\n \t /* Bla */ [ \t Blu \t ]->  \t outX.7",
		&data.Connection{DataType: "Blu", FromPort: data.NewPort("", 0), ToPort: data.NewIdxPort("outX", 7, 31)}, 0)
}
