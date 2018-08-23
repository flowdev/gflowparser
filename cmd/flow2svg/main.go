package main

import (
	"fmt"
	"io/ioutil"
	"os"

	data "github.com/flowdev/gflowparser"
	"github.com/flowdev/gflowparser/data2svg"
	"github.com/flowdev/gflowparser/parser"
	"github.com/flowdev/gflowparser/svg"
	"github.com/flowdev/gparselib"
)

func main() {
	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"ERROR: Unable to read flow DSL from standard input: %s.\n",
			err)
		os.Exit(2)
	}

	flowName := "standard input"
	flowContent := string(buf)

	pFlow, err := parser.NewFlowParser()
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"ERROR: Unable to create new flow parser: %s.\n", err)
		os.Exit(3)
	}
	pd := gparselib.NewParseData(flowName, flowContent)
	pd, _ = pFlow.ParseFlow(pd, nil)

	fb, err := parser.CheckFeedback(pd.Result)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(4)
	}
	os.Stderr.WriteString(fb)

	sf, err := data2svg.Convert(pd.Result.Value.(data.Flow), pd.Source)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"ERROR: Unable to convert flow data to SVG flow data: %s.\n", err)
		os.Exit(5)
	}

	buf, err = svg.FromFlowData(sf)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"ERROR: Unable to convert SVG flow data to SVG: %s.\n", err)
		os.Exit(6)
	}

	_, err = os.Stdout.Write(buf)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"ERROR: Unable to write SVG to standard output: %s.\n",
			err,
		)
		os.Exit(7)
	}
}
