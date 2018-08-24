package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/flowdev/gflowparser"
)

func main() {
	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"ERROR: Unable to read flow DSL from standard input: %s.\n", err)
		os.Exit(2)
	}

	buf, fb, err := gflowparser.ConvertFlowDSLToSVG(string(buf), "standard input")
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"ERROR: Unable to convert flow to SVG:\n%s", err)
		os.Exit(3)
	}
	os.Stderr.WriteString(fb)

	_, err = os.Stdout.Write(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"ERROR: Unable to write SVG to standard output: %s.\n", err)
		os.Exit(7)
	}
}
