package gflowparser_test

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/flowdev/gflowparser"
	"github.com/flowdev/gflowparser/data"
)

func TestConvertFlowDSLToSVG(t *testing.T) {
	specs := []struct {
		givenFlowName     string
		givenFlowContent  string
		expectedSVG       string
		expectedCompTypes []data.Type
		expectedDataTypes []data.Type
		expectedFeedback  string
		expectedError     string
	}{
		{
			givenFlowName:    "simple success",
			givenFlowContent: "in (data)-> [a] -> out",
			expectedSVG: `<?xml version="1.0" ?>
<svg version="1.1" xmlns="http://www.w3.org/2000/svg" width="259px" height="64px">
<!-- Generated by FlowDev tool. -->
	<rect fill="rgb(255,255,255)" fill-opacity="1" stroke="none" stroke-opacity="1" stroke-width="0.0" width="259" height="64" x="0" y="0"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="26" y1="25" x2="140" y2="25"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="132" y1="17" x2="140" y2="25"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="132" y1="33" x2="140" y2="25"/>

	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="176" y1="25" x2="218" y2="25"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="210" y1="17" x2="218" y2="25"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="210" y1="33" x2="218" y2="25"/>

	<rect fill="rgb(96,196,255)" fill-opacity="1.0" stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" width="36" height="36" x="140" y="7" rx="10" ry="10"/>


	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="3" y="31" textLength="22" lengthAdjust="spacingAndGlyphs">in</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="41" y="17" textLength="72" lengthAdjust="spacingAndGlyphs">(data)</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="152" y="31" textLength="12" lengthAdjust="spacingAndGlyphs">a</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="221" y="31" textLength="34" lengthAdjust="spacingAndGlyphs">out</text>
</svg>
`,
			expectedCompTypes: []data.Type{
				{LocalType: "a", SrcPos: 13},
			},
			expectedDataTypes: []data.Type{
				{LocalType: "data", SrcPos: 4},
			},
			expectedFeedback: "",
			expectedError:    ``,
		}, {
			givenFlowName:    "missing first data",
			givenFlowContent: "in ()-> [a] -> out",
			expectedSVG:      ``,
			expectedFeedback: "",
			expectedError: `Found errors while parsing flow:
ERROR: File 'missing first data', line 1, column 1:
in ()-> [a] -> out
At least 1 matches expected but got only 0.
ERROR: File 'missing first data', line 1, column 1:
in ()-> [a] -> out
At least 2 matches expected but got only 0.
ERROR: File 'missing first data', line 1, column 1:
in ()-> [a] -> out
Any subparser should match. But all 2 subparsers failed.
ERROR: File 'missing first data', line 1, column 4:
in ()-> [a] -> out
Literal '->' expected.
ERROR: File 'missing first data', line 1, column 1:
in ()-> [a] -> out
Literal '[' expected.
`,
		},
	}
	for _, spec := range specs {
		t.Logf("Testing flow: %s\n", spec.givenFlowName)
		gotSVG, gotCompTypes, gotDataTypes, gotFeedback, gotError := gflowparser.ConvertFlowDSLToSVG(
			spec.givenFlowContent, spec.givenFlowName)

		if spec.expectedError != "" && gotError != nil {
			if spec.expectedError != gotError.Error() {
				t.Errorf("Expected error '%s' but got: '%s'",
					spec.expectedError, gotError)
			}
		} else if spec.expectedError != "" && gotError == nil {
			t.Error("Expected an error but didn't get one.")
		} else if spec.expectedError == "" && gotError != nil {
			t.Errorf("Expected no error but got: %s", gotError)
		}
		if spec.expectedError != "" || gotError != nil {
			continue
		}
		if !reflect.DeepEqual(gotCompTypes, spec.expectedCompTypes) {
			t.Errorf("Expected component types '%#v' but got: '%#v'",
				spec.expectedCompTypes, gotCompTypes)
		}
		if !reflect.DeepEqual(gotDataTypes, spec.expectedDataTypes) {
			t.Errorf("Expected component types '%#v' but got: '%#v'",
				spec.expectedDataTypes, gotDataTypes)
		}
		if spec.expectedFeedback != gotFeedback {
			t.Errorf("Expected feedback '%s' but got: '%s'",
				spec.expectedFeedback, gotFeedback)
		}
		if spec.expectedSVG != string(gotSVG) {
			ioutil.WriteFile("fail.svg", gotSVG, os.FileMode(0644))
			t.Fatal("Got unexpected SVG, please look into 'fail.svg'.")
		}
	}
}