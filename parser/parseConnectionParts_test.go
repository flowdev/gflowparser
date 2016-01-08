package parser

import (
	"github.com/flowdev/gflowparser/data"
	"testing"
)

func TestParseConnectionPart(t *testing.T) {
	p := NewParseConnectionPart()

	runTest(t, p, "empty", "", nil, 3)
	runTest(t, p, "no match 1", "()",
		&data.Operation{
			InPorts:  []*data.PortData{&data.PortData{Name: "in", CapName: "In", SrcPos: 0}},
			OutPorts: []*data.PortData{&data.PortData{Name: "out", CapName: "Out", SrcPos: 2}}}, 1)
	runTest(t, p, "no match 2", "bla", nil, 3)
	runTest(t, p, "no match 3", "Bla", nil, 3)
	runTest(t, p, "simple 1", "(Bla)",
		&data.Operation{Name: "bla", Type: "Bla",
			InPorts:  []*data.PortData{&data.PortData{Name: "in", CapName: "In", HasIndex: false, SrcPos: 0}},
			OutPorts: []*data.PortData{&data.PortData{Name: "out", CapName: "Out", HasIndex: false, SrcPos: 5}}}, 0)
	runTest(t, p, "simple 2", "in.2 \t bla() \t error ",
		&data.Operation{Name: "bla", Type: "",
			InPorts:  []*data.PortData{&data.PortData{Name: "in", CapName: "In", HasIndex: true, Index: 2, SrcPos: 0}},
			OutPorts: []*data.PortData{&data.PortData{Name: "error", CapName: "Error", HasIndex: false, SrcPos: 15}}}, 0)
	runTest(t, p, "simple 3", "xIn.1   bla(Blu)outY.123",
		&data.Operation{Name: "bla", Type: "Blu",
			InPorts:  []*data.PortData{&data.PortData{Name: "xIn", CapName: "XIn", HasIndex: true, Index: 1, SrcPos: 0}},
			OutPorts: []*data.PortData{&data.PortData{Name: "outY", CapName: "OutY", HasIndex: true, Index: 123, SrcPos: 16}}}, 0)
}

func TestParseOperationNameParens(t *testing.T) {
	p := NewParseOperationNameParens()

	runTest(t, p, "empty", "", nil, 2)
	runTest(t, p, "no match 1", "()", &data.Operation{}, 1)
	runTest(t, p, "no match 2", "bla", nil, 2)
	runTest(t, p, "no match 3", "Bla", nil, 2)
	runTest(t, p, "simple 1", "(Bla)", &data.Operation{Name: "bla", Type: "Bla"}, 0)
	runTest(t, p, "simple 2", "bla()", &data.Operation{Name: "bla", Type: ""}, 0)
	runTest(t, p, "simple 3", "bla(Blu)", &data.Operation{Name: "bla", Type: "Blu"}, 0)
	runTest(t, p, "simple 4", "bla \t ( \t Blu \t ) \t ", &data.Operation{Name: "bla", Type: "Blu"}, 0)
}

func TestParseOptOperationType(t *testing.T) {
	p := NewParseOptOperationType()

	runTest(t, p, "empty", "", nil, 0)
	runTest(t, p, "no match 1", "blu", nil, 0)
	runTest(t, p, "no match 2", "B", nil, 0)
	runTest(t, p, "no match 3", " ", nil, 0)
	runTest(t, p, "simple 1", "Bla", "Bla", 0)
	runTest(t, p, "simple 2", "Blu  ", "Blu", 0)
	runTest(t, p, "simple 3", "Blu  \t  \t ", "Blu", 0)
}

func TestParseArrow(t *testing.T) {
	p := NewParseArrow()

	runTest(t, p, "empty", "", nil, 2)
	runTest(t, p, "no match 1", "-", nil, 2)
	runTest(t, p, "no match 3", " /* \n */ \t [Bla]>", nil, 2)
	runTest(t, p, "simple 1", "[Bla]->", "Bla", 0)
	runTest(t, p, "simple 2", "->", "", 0)
	runTest(t, p, "simple 3", "\n \t /* Blu */ [ \t Bla \t ]->  \t   ", "Bla", 0)
	runTest(t, p, "simple 4", " \r\n // blu \n \t -> \r\n \t", "", 0)
}

func TestParseOptPortSpc(t *testing.T) {
	p := NewParseOptPortSpc()

	runTest(t, p, "empty", "", nil, 0)
	runTest(t, p, "no space", "p.1", nil, 0)
	runTest(t, p, "no match 2", "pt. ", nil, 0)
	runTest(t, p, "simple 1", "p ", data.NewPort("p", 0), 0)
	runTest(t, p, "simple 2", "pt.0\t", data.NewIdxPort("pt", 0, 0), 0)
	runTest(t, p, "long port name", "looooongPortName \t ", data.NewPort("looooongPortName", 0), 0)
	runTest(t, p, "name and index", "port.123 \t ", data.NewIdxPort("port", 123, 0), 0)
	runTest(t, p, "too large index", "port.9999999999 ", nil, 1)
}

func TestParseOptPort(t *testing.T) {
	p := NewParseOptPort()

	runTest(t, p, "empty", "", nil, 0)
	runTest(t, p, "no match", ".1", nil, 0)
	runTest(t, p, "half match 1", "pt.", data.NewPort("pt", 0), 0)
	runTest(t, p, "half match 2", "pt_1", data.NewPort("pt", 0), 0)
	runTest(t, p, "simple 1", "p", data.NewPort("p", 0), 0)
	runTest(t, p, "simple 2", "pt.0", data.NewIdxPort("pt", 0, 0), 0)
	runTest(t, p, "long port name", "looooongPortName", data.NewPort("looooongPortName", 0), 0)
	runTest(t, p, "name and index", "port.123", data.NewIdxPort("port", 123, 0), 0)
	runTest(t, p, "too large index", "port.9999999999", nil, 1)
}

func TestParsePort(t *testing.T) {
	p := NewParsePort()

	runTest(t, p, "empty", "", nil, 2)
	runTest(t, p, "no match", ".1", nil, 2)
	runTest(t, p, "half match 1", "pt.", data.NewPort("pt", 0), 0)
	runTest(t, p, "half match 2", "pt_1", data.NewPort("pt", 0), 0)
	runTest(t, p, "simple 1", "p", data.NewPort("p", 0), 0)
	runTest(t, p, "simple 2", "pt.0", data.NewIdxPort("pt", 0, 0), 0)
	runTest(t, p, "long port name", "looooongPortName", data.NewPort("looooongPortName", 0), 0)
	runTest(t, p, "name and index", "port.123", data.NewIdxPort("port", 123, 0), 0)
	runTest(t, p, "too large index", "port.9999999999", nil, 1)
}
