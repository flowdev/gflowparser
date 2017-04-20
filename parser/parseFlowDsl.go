package parser

import (
	"github.com/flowdev/gflowparser/data"
	"github.com/flowdev/gparselib"
)

// ------------ ParseFlowFile:
// semantic result: flowFile &data.FlowFile{FileName, Version, Flows}
type SemanticFlowFile struct {
	outPort func(interface{})
}

func NewSemanticFlowFile() *SemanticFlowFile {
	return &SemanticFlowFile{}
}
func (op *SemanticFlowFile) InPort(dat interface{}) {
	md := dat.(*data.MainData)
	res := md.ParseData.Result
	subVals := res.Value.([]interface{})

	version := subVals[0].(data.Version)
	flows := subVals[1].([]*data.Flow)

	res.Value = &data.FlowFile{FileName: md.ParseData.Source.Name, Version: version, Flows: flows}
	op.outPort(md)
}
func (op *SemanticFlowFile) SetOutPort(port func(interface{})) {
	op.outPort = port
}

type ParseFlowFile struct {
	flowFile *gparselib.ParseAll
	//	semantic *SemanticFlowFile
	//version *ParseVersion
	//flows   *gparselib.ParseMulti1
	//eof     *gparselib.ParseEof
	//flow    *ParseFlow
	InPort func(interface{})
}

func NewParseFlowFile() *ParseFlowFile {
	flowFile := gparselib.NewParseAll(parseData, setParseData)
	semantic := NewSemanticFlowFile()
	version := NewParseVersion()
	flows := gparselib.NewParseMulti1(parseData, setParseData)
	eof := gparselib.NewParseEof(parseData, setParseData)
	flow := NewParseFlow()

	flowFile.SetSemOutPort(semantic.InPort)
	semantic.SetOutPort(flowFile.SemInPort)
	flowFile.AppendSubOutPort(version.InPort)
	version.SetOutPort(flowFile.SubInPort)
	flowFile.AppendSubOutPort(flows.InPort)
	flows.SetOutPort(flowFile.SubInPort)
	flowFile.AppendSubOutPort(eof.InPort)
	eof.SetOutPort(flowFile.SubInPort)
	flows.SetSubOutPort(flow.InPort)
	flow.SetOutPort(flows.SubInPort)

	return &ParseFlowFile{flowFile: flowFile, InPort: flowFile.InPort}
}
func (f *ParseFlowFile) SetOutPort(port func(interface{})) {
	f.flowFile.SetOutPort(port)
}

// ------------ ParseVersion:
// semantic result: vers &data.Version{Political, Major}
type SemanticVersion struct {
	outPort func(interface{})
}

func NewSemanticVersion() *SemanticVersion {
	return &SemanticVersion{}
}
func (op *SemanticVersion) InPort(dat interface{}) {
	md := dat.(*data.MainData)
	res := md.ParseData.Result
	subVals := res.Value.([]interface{})

	political := subVals[3].(uint64)
	major := subVals[5].(uint64)

	res.Value = &data.Version{Political: int(political), Major: int(major)}
	op.outPort(md)
}
func (op *SemanticVersion) SetOutPort(port func(interface{})) {
	op.outPort = port
}

type ParseVersion struct {
	version *gparselib.ParseAll
	//semantic   *SemanticVersion
	//spcCommBeg *ParseSpaceComment
	//vers       *gparselib.ParseLiteral
	//aspc       *gparselib.ParseSpace
	//political  *gparselib.ParseNatural
	//dot        *gparselib.ParseLiteral
	//major      *gparselib.ParseNatural
	//spcCommEnd *ParseSpaceComment
	InPort func(interface{})
}

func NewParseVersion() *ParseVersion {
	version := gparselib.NewParseAll(parseData, setParseData)
	semantic := NewSemanticVersion()
	spcCommBeg := NewParseSpaceComment()
	vers := gparselib.NewParseLiteral(parseData, setParseData, "version")
	aspc := gparselib.NewParseSpace(parseData, setParseData, false)
	political := gparselib.NewParseNatural(parseData, setParseData, 10)
	dot := gparselib.NewParseLiteral(parseData, setParseData, ".")
	major := gparselib.NewParseNatural(parseData, setParseData, 10)
	spcCommEnd := NewParseSpaceComment()

	version.SetSemOutPort(semantic.InPort)
	semantic.SetOutPort(version.SemInPort)
	version.AppendSubOutPort(spcCommBeg.InPort)
	spcCommBeg.SetOutPort(version.SubInPort)
	version.AppendSubOutPort(vers.InPort)
	vers.SetOutPort(version.SubInPort)
	version.AppendSubOutPort(aspc.InPort)
	aspc.SetOutPort(version.SubInPort)
	version.AppendSubOutPort(political.InPort)
	political.SetOutPort(version.SubInPort)
	version.AppendSubOutPort(dot.InPort)
	dot.SetOutPort(version.SubInPort)
	version.AppendSubOutPort(major.InPort)
	major.SetOutPort(version.SubInPort)
	version.AppendSubOutPort(spcCommEnd.InPort)
	spcCommEnd.SetOutPort(version.SubInPort)

	return &ParseVersion{version: version, InPort: version.InPort}
}
func (f *ParseVersion) SetOutPort(port func(interface{})) {
	f.version.SetOutPort(port)
}

// ------------ ParseFlow:
// semantic result: flow *data.Flow including name
type SemanticFlow struct {
	outPort func(interface{})
}

func NewSemanticFlow() *SemanticFlow {
	return &SemanticFlow{}
}
func (op *SemanticFlow) InPort(dat interface{}) {
	md := dat.(*data.MainData)
	res := md.ParseData.Result
	subVals := res.Value.([]interface{})

	name := subVals[2].(string)
	flow := subVals[6].(*data.Flow)
	flow.Name = name

	res.Value = flow
	op.outPort(md)
}
func (op *SemanticFlow) SetOutPort(port func(interface{})) {
	op.outPort = port
}

type ParseFlow struct {
	flow *gparselib.ParseAll
	//semantic    *SemanticFlow
	//flowLiteral *gparselib.ParseLiteral
	//aspc        *gparselib.ParseSpace
	//name        *ParseBigIdent
	//spcComm1    *ParseSpaceComment
	//openFlow    *gparselib.ParseLiteral
	//spcComm2    *ParseSpaceComment
	//connections *ParseConnections
	//closeFlow   *gparselib.ParseLiteral
	//spcComm3    *ParseSpaceComment
	InPort func(interface{})
}

func NewParseFlow() *ParseFlow {
	flow := gparselib.NewParseAll(parseData, setParseData)
	semantic := NewSemanticFlow()
	flowLiteral := gparselib.NewParseLiteral(parseData, setParseData, "flow")
	aspc := gparselib.NewParseSpace(parseData, setParseData, false)
	name := NewParseBigIdent()
	spcComm1 := NewParseSpaceComment()
	openFlow := gparselib.NewParseLiteral(parseData, setParseData, "{")
	spcComm2 := NewParseSpaceComment()
	connections := NewParseConnections()
	closeFlow := gparselib.NewParseLiteral(parseData, setParseData, "}")
	spcComm3 := NewParseSpaceComment()

	flow.SetSemOutPort(semantic.InPort)
	semantic.SetOutPort(flow.SemInPort)
	flow.AppendSubOutPort(flowLiteral.InPort)
	flowLiteral.SetOutPort(flow.SubInPort)
	flow.AppendSubOutPort(aspc.InPort)
	aspc.SetOutPort(flow.SubInPort)
	flow.AppendSubOutPort(name.InPort)
	name.SetOutPort(flow.SubInPort)
	flow.AppendSubOutPort(spcComm1.InPort)
	spcComm1.SetOutPort(flow.SubInPort)
	flow.AppendSubOutPort(openFlow.InPort)
	openFlow.SetOutPort(flow.SubInPort)
	flow.AppendSubOutPort(spcComm2.InPort)
	spcComm2.SetOutPort(flow.SubInPort)
	flow.AppendSubOutPort(connections.InPort)
	connections.SetOutPort(flow.SubInPort)
	flow.AppendSubOutPort(closeFlow.InPort)
	closeFlow.SetOutPort(flow.SubInPort)
	flow.AppendSubOutPort(spcComm3.InPort)
	spcComm3.SetOutPort(flow.SubInPort)

	return &ParseFlow{flow: flow, InPort: flow.InPort}
}
func (f *ParseFlow) SetOutPort(port func(interface{})) {
	f.flow.SetOutPort(port)
}
