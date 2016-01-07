package parser

import (
	"github.com/flowdev/gflowparser/data"
	"testing"
)

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
