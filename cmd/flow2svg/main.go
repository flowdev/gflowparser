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
			"ERROR: Unable to read flow DSL from standard input: %s.\n",
			err)
		os.Exit(2)
	}

	fts, err := gflowparser.NewFlowToSVG()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s.\n", err)
		os.Exit(3)
	}
	buf, fb, err := fts.ConvertFlowToSVG(string(buf), "standard input")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s.\n", err)
		os.Exit(4)
	}
	os.Stderr.WriteString(fb)

	_, err = os.Stdout.Write(buf)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"ERROR: Unable to write SVG to standard output: %s.\n",
			err,
		)
		os.Exit(5)
	}
}
