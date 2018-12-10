package gflowparser

import (
	"github.com/flowdev/gflowparser/data"
	"github.com/flowdev/gflowparser/data2svg"
	"github.com/flowdev/gflowparser/parser"
	"github.com/flowdev/gflowparser/svg"
	"github.com/flowdev/gparselib"
)

// ConvertFlowDSLToSVG transforms a flow given as DSL string into a SVG image
// plus component (subflow) types, data types, (currently empty) feedback
// string and potential error.
func ConvertFlowDSLToSVG(flowContent, flowName string,
) (
	svgData []byte,
	compTypes []data.Type,
	dataTypes []data.Type,
	feedback string,
	err error,
) {
	pd := gparselib.NewParseData(flowName, flowContent)
	pFlow, err := parser.NewFlowParser()
	if err != nil {
		return nil, nil, nil, "", err
	}
	pd, _ = pFlow.ParseFlow(pd, nil)

	fb, err := parser.CheckFeedback(pd.Result)
	if err != nil {
		return nil, nil, nil, "", err
	}

	flow := pd.Result.Value.(data.Flow)

	sf, err := data2svg.Convert(flow, pd.Source)
	if err != nil {
		return nil, nil, nil, "", err
	}
	compTypes, dataTypes = extractTypes(flow)

	buf, err := svg.FromFlowData(sf)
	if err != nil {
		return nil, nil, nil, "", err
	}

	return buf, compTypes, dataTypes, fb, nil
}

func extractTypes(flow data.Flow) (compTypes []data.Type, dataTypes []data.Type) {
	dataMap := make(map[string]data.Type)
	compMap := make(map[string]data.Type)

	for _, partLine := range flow.Parts {
		for _, part := range partLine {
			switch p := part.(type) {
			case data.Arrow:
				dataMap = addTypes(dataMap, p.Data)
			case data.Component:
				// check component, plugins, ...
				compMap[typToString(p.Decl.Type)] = p.Decl.Type
				for _, plugin := range p.Plugins {
					compMap = addTypes(compMap, plugin.Types)
				}
			}
		}
	}
	return valuesOf(compMap), valuesOf(dataMap)
}
func valuesOf(typeMap map[string]data.Type) []data.Type {
	types := make([]data.Type, 0, len(typeMap))
	for _, t := range typeMap {
		types = append(types, t)
	}
	return types
}
func addTypes(typeMap map[string]data.Type, types []data.Type) map[string]data.Type {
	for _, t := range types {
		typeMap[typToString(t)] = t
	}
	return typeMap
}
func typToString(t data.Type) string {
	return t.Package + "." + t.LocalType
}