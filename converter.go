package gflowparser

import (
	"github.com/flowdev/gflowparser/data"
	"github.com/flowdev/gflowparser/data2svg"
	"github.com/flowdev/gflowparser/parser"
	"github.com/flowdev/gflowparser/svg"
	"github.com/flowdev/gparselib"
)

// ConvertFlowDSLToSVG transforms a flow given as DSL string into a SVG image
// plus (currently empty) feedback string and potential error.
func ConvertFlowDSLToSVG(flowContent, flowName string) ([]byte, string, error) {
	pd := gparselib.NewParseData(flowName, flowContent)
	pFlow, err := parser.NewFlowParser()
	if err != nil {
		return nil, "", err
	}
	pd, _ = pFlow.ParseFlow(pd, nil)

	fb, err := parser.CheckFeedback(pd.Result)
	if err != nil {
		return nil, "", err
	}

	sf, err := data2svg.Convert(pd.Result.Value.(data.Flow), pd.Source)
	if err != nil {
		return nil, "", err
	}

	buf, err := svg.FromFlowData(sf)
	if err != nil {
		return nil, "", err
	}

	return buf, fb, nil
}
