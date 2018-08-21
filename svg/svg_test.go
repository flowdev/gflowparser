package svg_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/flowdev/gflowparser/svg"
)

const expSVG = `<?xml version="1.0" ?>
<svg version="1.1" xmlns="http://www.w3.org/2000/svg" width="826px" height="589px">
<!-- Generated by simple FlowDev draw-svg tool. -->
	<rect fill="rgb(255,255,255)" fill-opacity="1" stroke="none" stroke-opacity="1" stroke-width="0.0" width="826" height="589" x="0" y="0"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="26" y1="25" x2="116" y2="25"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="108" y1="17" x2="116" y2="25"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="108" y1="33" x2="116" y2="25"/>

	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="212" y1="25" x2="362" y2="25"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="354" y1="17" x2="362" y2="25"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="354" y1="33" x2="362" y2="25"/>

	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="530" y1="25" x2="704" y2="25"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="696" y1="17" x2="704" y2="25"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="696" y1="33" x2="704" y2="25"/>

	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="212" y1="231" x2="314" y2="231"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="306" y1="223" x2="314" y2="231"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="306" y1="239" x2="314" y2="231"/>

	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="410" y1="231" x2="512" y2="231"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="504" y1="223" x2="512" y2="231"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="504" y1="239" x2="512" y2="231"/>

	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="38" y1="308" x2="104" y2="308"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="96" y1="300" x2="104" y2="308"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="96" y1="316" x2="104" y2="308"/>

	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="200" y1="308" x2="704" y2="308"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="696" y1="300" x2="704" y2="308"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="696" y1="316" x2="704" y2="308"/>

	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="38" y1="385" x2="140" y2="385"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="132" y1="377" x2="140" y2="385"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="132" y1="393" x2="140" y2="385"/>

	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="308" y1="385" x2="704" y2="385"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="696" y1="377" x2="704" y2="385"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" x1="696" y1="393" x2="704" y2="385"/>

	<rect fill="rgb(96,196,255)" fill-opacity="1.0" stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" width="96" height="272" x="116" y="7" rx="10" ry="10"/>
	<rect fill="rgb(96,196,255)" fill-opacity="1.0" stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" width="168" height="189" x="362" y="7" rx="10" ry="10"/>
	<rect fill="rgb(32,224,32)" fill-opacity="1.0" stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" width="168" height="57" x="362" y="43"/>
	<rect fill="rgb(32,224,32)" fill-opacity="1.0" stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" width="168" height="84" x="362" y="100"/>
	<rect fill="rgb(96,196,255)" fill-opacity="1.0" stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" width="96" height="60" x="314" y="213" rx="10" ry="10"/>
	<rect fill="rgb(96,196,255)" fill-opacity="1.0" stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" width="96" height="60" x="104" y="290" rx="10" ry="10"/>
	<rect fill="rgb(96,196,255)" fill-opacity="1.0" stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" width="168" height="213" x="140" y="367" rx="10" ry="10"/>
	<rect fill="rgb(32,224,32)" fill-opacity="1.0" stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" width="168" height="57" x="140" y="427"/>
	<rect fill="rgb(32,224,32)" fill-opacity="1.0" stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" width="168" height="84" x="140" y="484"/>
	<rect fill="rgb(96,196,255)" fill-opacity="1.0" stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="2.5" width="120" height="408" x="704" y="7" rx="10" ry="10"/>

	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="1.0" x1="362" y1="70" x2="530" y2="70"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="1.0" x1="362" y1="127" x2="530" y2="127"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="1.0" x1="362" y1="154" x2="530" y2="154"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="1.0" x1="140" y1="454" x2="308" y2="454"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="1.0" x1="140" y1="511" x2="308" y2="511"/>
	<line stroke="rgb(0,0,0)" stroke-opacity="1.0" stroke-width="1.0" x1="140" y1="538" x2="308" y2="538"/>

	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="3" y="31" textLength="22" lengthAdjust="spacingAndGlyphs">in</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="41" y="17" textLength="48" lengthAdjust="spacingAndGlyphs">Data</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="128" y="31" textLength="24" lengthAdjust="spacingAndGlyphs">ra</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="128" y="55" textLength="72" lengthAdjust="spacingAndGlyphs">(MiSo)</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="218" y="45" textLength="84" lengthAdjust="spacingAndGlyphs">special</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="257" y="17" textLength="48" lengthAdjust="spacingAndGlyphs">Data</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="326" y="45" textLength="24" lengthAdjust="spacingAndGlyphs">in</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="374" y="31" textLength="24" lengthAdjust="spacingAndGlyphs">do</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="368" y="64" textLength="120" lengthAdjust="spacingAndGlyphs">semantics:</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="368" y="91" textLength="156" lengthAdjust="spacingAndGlyphs">TextSemantics</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="368" y="121" textLength="120" lengthAdjust="spacingAndGlyphs">subParser:</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="368" y="148" textLength="144" lengthAdjust="spacingAndGlyphs">LitralParser</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="368" y="175" textLength="156" lengthAdjust="spacingAndGlyphs">NaturalParser</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="536" y="45" textLength="36" lengthAdjust="spacingAndGlyphs">out</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="545" y="17" textLength="132" lengthAdjust="spacingAndGlyphs">BigDataType</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="656" y="45" textLength="36" lengthAdjust="spacingAndGlyphs">in1</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="218" y="251" textLength="36" lengthAdjust="spacingAndGlyphs">out</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="233" y="223" textLength="48" lengthAdjust="spacingAndGlyphs">Data</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="278" y="251" textLength="24" lengthAdjust="spacingAndGlyphs">in</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="326" y="237" textLength="36" lengthAdjust="spacingAndGlyphs">bla</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="326" y="261" textLength="72" lengthAdjust="spacingAndGlyphs">(Blue)</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="425" y="223" textLength="60" lengthAdjust="spacingAndGlyphs">Data2</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="515" y="237" textLength="34" lengthAdjust="spacingAndGlyphs">...</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="3" y="314" textLength="34" lengthAdjust="spacingAndGlyphs">...</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="68" y="328" textLength="24" lengthAdjust="spacingAndGlyphs">in</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="116" y="314" textLength="36" lengthAdjust="spacingAndGlyphs">bla</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="116" y="338" textLength="72" lengthAdjust="spacingAndGlyphs">(Blue)</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="206" y="328" textLength="36" lengthAdjust="spacingAndGlyphs">out</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="422" y="300" textLength="48" lengthAdjust="spacingAndGlyphs">Data</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="656" y="328" textLength="36" lengthAdjust="spacingAndGlyphs">in2</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="3" y="391" textLength="34" lengthAdjust="spacingAndGlyphs">in2</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="53" y="377" textLength="60" lengthAdjust="spacingAndGlyphs">Data3</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="104" y="405" textLength="24" lengthAdjust="spacingAndGlyphs">in</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="152" y="391" textLength="120" lengthAdjust="spacingAndGlyphs">megaParser</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="152" y="415" textLength="144" lengthAdjust="spacingAndGlyphs">(MegaParser)</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="146" y="448" textLength="120" lengthAdjust="spacingAndGlyphs">semantics:</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="146" y="475" textLength="156" lengthAdjust="spacingAndGlyphs">TextSemantics</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="146" y="505" textLength="120" lengthAdjust="spacingAndGlyphs">subParser:</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="146" y="532" textLength="144" lengthAdjust="spacingAndGlyphs">LitralParser</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="146" y="559" textLength="156" lengthAdjust="spacingAndGlyphs">NaturalParser</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="314" y="405" textLength="36" lengthAdjust="spacingAndGlyphs">out</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="476" y="377" textLength="48" lengthAdjust="spacingAndGlyphs">Data</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="656" y="405" textLength="36" lengthAdjust="spacingAndGlyphs">in3</text>
	<text fill="rgb(0,0,0)" fill-opacity="1.0" font-family="monospace" font-size="16" x="716" y="31" textLength="96" lengthAdjust="spacingAndGlyphs">BigMerge</text>
</svg>
`

func TestFromFlowData(t *testing.T) {
	testFlow := svg.BigTestFlowData

	gotBytes, gotErr := svg.FromFlowData(testFlow)

	if gotErr != nil {
		t.Fatalf("Unexpected error: %s", gotErr)
	}
	if expSVG != string(gotBytes) {
		t.Error("Got unexpected SVG, please look into 'fail.svg'.")
		ioutil.WriteFile("fail.svg", gotBytes, os.FileMode(0644))
	}
}
