package parser

import (
	"testing"

	"github.com/flowdev/gflowparser/data"
)

func TestParseConnections(t *testing.T) {
	p := NewParseConnections()
	doIt := &data.Operation{Name: "doIt", Type: "DoIt", SrcPos: 2, InPorts: []*data.PortData{data.NewPort("in", 2)}, OutPorts: []*data.PortData{}}
	minFlow := &data.Flow{Ops: []*data.Operation{&data.Operation{Name: "bla", Type: "Bla", OutPorts: []*data.PortData{}}}, Conns: []*data.Connection{}}
	simpleFlow1 := &data.Flow{Ops: []*data.Operation{doIt}, Conns: []*data.Connection{
		&data.Connection{FromPort: data.NewPort("in", 0), ToPort: doIt.InPorts[0], ToOp: doIt},
	}}

	blue := &data.Operation{Name: "blue", Type: "Blue", SrcPos: 3, InPorts: []*data.PortData{data.NewPort("in", 3)},
		OutPorts: []*data.PortData{data.NewPort("out", 10)}}
	simpleFlow2 := &data.Flow{Ops: []*data.Operation{blue}, Conns: []*data.Connection{
		&data.Connection{FromPort: data.NewPort("in", 0), ToPort: blue.InPorts[0], ToOp: blue},
		&data.Connection{FromOp: blue, FromPort: blue.OutPorts[0], ToPort: data.NewPort("out", 13)},
		&data.Connection{FromPort: data.NewPort("in2", 15), ToPort: blue.InPorts[0], ToOp: blue},
		&data.Connection{FromPort: data.NewPort("in3", 30), ToPort: blue.InPorts[0], ToOp: blue},
	}}

	bla := &data.Operation{Name: "bla", Type: "Bla", SrcPos: 3, InPorts: []*data.PortData{data.NewPort("in", 3)},
		OutPorts: []*data.PortData{data.NewPort("out", 9)}}
	blue2 := &data.Operation{Name: "blue", Type: "Blue", SrcPos: 12, InPorts: []*data.PortData{data.NewPort("in", 12)},
		OutPorts: []*data.PortData{data.NewPort("out", 19), data.NewPort("out2", 38)}}
	simpleFlow3 := &data.Flow{Ops: []*data.Operation{bla, blue2}, Conns: []*data.Connection{
		&data.Connection{FromPort: data.NewPort("in", 0), ToPort: bla.InPorts[0], ToOp: bla},
		&data.Connection{FromOp: bla, FromPort: bla.OutPorts[0], ToPort: blue2.InPorts[0], ToOp: blue2},
		&data.Connection{FromOp: blue2, FromPort: blue2.OutPorts[0], ToPort: data.NewPort("out", 22)},
		&data.Connection{FromPort: data.NewPort("in2", 24), ToPort: blue2.InPorts[0], ToOp: blue2},
		&data.Connection{FromOp: blue2, FromPort: blue2.OutPorts[1], ToPort: data.NewPort("out2", 45)},
	}}

	blaa := &data.Operation{Name: "blaa", Type: "Bla", SrcPos: 11, InPorts: []*data.PortData{data.NewIdxPort("i", 0, 7)},
		OutPorts: []*data.PortData{data.NewIdxPort("o", 0, 21)}}
	bluu := &data.Operation{Name: "bluu", Type: "Blue", SrcPos: 32, InPorts: []*data.PortData{data.NewIdxPort("i", 1, 28), data.NewPort("in", 62)},
		OutPorts: []*data.PortData{data.NewIdxPort("o", 3, 43), data.NewIdxPort("o", 2, 69)}}
	ab := &data.Operation{Name: "ab", Type: "Ab", SrcPos: 84, OutPorts: []*data.PortData{data.NewPort("out", 89)}}
	complexFlow := &data.Flow{Ops: []*data.Operation{blaa, bluu, ab}, Conns: []*data.Connection{
		&data.Connection{FromPort: data.NewIdxPort("i", 1, 0), ToPort: blaa.InPorts[0], ToOp: blaa},
		&data.Connection{FromOp: blaa, FromPort: blaa.OutPorts[0], ToPort: bluu.InPorts[0], ToOp: bluu},
		&data.Connection{FromOp: bluu, FromPort: bluu.OutPorts[0], ToPort: data.NewIdxPort("o", 3, 50)},
		&data.Connection{FromPort: data.NewIdxPort("in", 2, 54), ToPort: bluu.InPorts[1], ToOp: bluu},
		&data.Connection{FromOp: bluu, FromPort: bluu.OutPorts[1], ToPort: data.NewPort("out2", 76)},
		&data.Connection{FromOp: ab, FromPort: ab.OutPorts[0], ToPort: data.NewPort("out1", 92)},
	}}

	runTest(t, p, "no match 1", "-> (Bla) -> ", nil, 2)
	runTest(t, p, "min flow", "(Bla) \r\n\t ;", minFlow, 0)
	runTest(t, p, "(un)indexed port error", "-> (Blue) -> ;\nin2 -> in.2 (Blue) out.2 -> out2;", nil, 3)
	runTest(t, p, "simple flow 1", "->doIt(DoIt);", simpleFlow1, 0)
	runTest(t, p, "multiple input ports flow", "-> (Blue) -> ; in2 -> (Blue); in3 -> blue();", simpleFlow2, 0)
	runTest(t, p, "split out port error", "-> (Blue) -> ; in2 -> (Blue) ->out2;", nil, 1)
	runTest(t, p, "simple flow 3", "-> (Bla) -> (Blue) -> ; in2 -> (Blue) out2 ->out2;", simpleFlow3, 0)
	runTest(t, p, "complex flow", "i.1 -> i.0 blaa(Bla) o.0 -> i.1 bluu(Blue) o.3 -> ;\n  in.2 -> bluu() o.2 -> out2;\n  (Ab) -> out1;", complexFlow, 0)
}

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
