package main

import (
	"fmt"
	"os"

	"github.com/flowdev/gflowparser/svg"
)

func main() {
	buf, err := svg.FromFlowData(flowData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s.\n", err)
		os.Exit(2)
	}
	_, err = os.Stdout.Write(buf)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"ERROR: Unable to write SVG to standard out: %s.\n",
			err,
		)
		os.Exit(3)
	}
}
