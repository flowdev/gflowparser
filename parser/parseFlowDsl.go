package parser

import (
	"github.com/flowdev/gparselib"
)

// ------------ ParseFlowFile:
type ParseFlowFile struct {
	flowFile *gparselib.ParseAll
	//	semantic *SemanticCreateFlowFileData
	version *ParseVersion
	flows   *gparselib.ParseMulti1
	eof     *gparselib.ParseEof
	flow    *ParseFlow
	InPort  func(interface{})
}

func NewParseFlowFile() *ParseFlowFile {
	f := &ParseFlowFile{}
	f.flowFile = gparselib.NewParseAll(parseData, setParseData)
	//	f.semantic = NewSemanticCreateFlowFileData()
	f.version = NewParseVersion()
	f.flows = gparselib.NewParseMulti1(parseData, setParseData)
	f.eof = gparselib.NewParseEof(parseData, setParseData)
	f.flow = NewParseFlow()

	//	f.flowFile.SetSemOutPort(f.semantic.InPort)
	//	f.semantic.SetOutPort(f.flowFile.SemInPort)
	f.flowFile.AppendSubOutPort(f.version.InPort)
	f.version.SetOutPort(f.flowFile.SubInPort)
	f.flowFile.AppendSubOutPort(f.flows.InPort)
	f.flows.SetOutPort(f.flowFile.SubInPort)
	f.flowFile.AppendSubOutPort(f.eof.InPort)
	f.eof.SetOutPort(f.flowFile.SubInPort)
	f.flows.SetSubOutPort(f.flow.InPort)
	f.flow.SetOutPort(f.flows.SubInPort)

	f.InPort = f.flowFile.InPort

	return f
}
func (f *ParseFlowFile) SetOutPort(port func(interface{})) { // datatype: FlowFile ?
	f.flowFile.SetOutPort(port)
}

// ------------ ParseVersion:
type ParseVersion struct {
	version *gparselib.ParseAll
	//	semantic  *SemanticCreateVersion
	spcComm   *ParseSpaceComment
	vers      *gparselib.ParseLiteral
	aspc      *gparselib.ParseSpace
	political *gparselib.ParseNatural
	dot       *gparselib.ParseLiteral
	major     *gparselib.ParseNatural
	InPort    func(interface{})
}

func NewParseVersion() *ParseVersion {
	f := &ParseVersion{}
	f.version = gparselib.NewParseAll(parseData, setParseData)
	//	f.semantic = NewSemanticCreateVersion()
	f.spcComm = NewParseSpaceComment()
	f.vers = gparselib.NewParseLiteral(parseData, setParseData, "version")
	f.aspc = gparselib.NewParseSpace(parseData, setParseData, false)
	f.political = gparselib.NewParseNatural(parseData, setParseData, 10)
	f.dot = gparselib.NewParseLiteral(parseData, setParseData, ".")
	f.major = gparselib.NewParseNatural(parseData, setParseData, 10)

	//	f.version.SetSemOutPort(f.semantic.InPort)
	//	f.semantic.SetOutPort(f.version.SemInPort)
	f.version.AppendSubOutPort(f.spcComm.InPort)
	f.spcComm.SetOutPort(f.version.SubInPort)
	f.version.AppendSubOutPort(f.vers.InPort)
	f.vers.SetOutPort(f.version.SubInPort)
	f.version.AppendSubOutPort(f.aspc.InPort)
	f.aspc.SetOutPort(f.version.SubInPort)
	f.version.AppendSubOutPort(f.political.InPort)
	f.political.SetOutPort(f.version.SubInPort)
	f.version.AppendSubOutPort(f.dot.InPort)
	f.dot.SetOutPort(f.version.SubInPort)
	f.version.AppendSubOutPort(f.major.InPort)
	f.major.SetOutPort(f.version.SubInPort)
	f.version.AppendSubOutPort(f.spcComm.InPort)
	f.spcComm.SetOutPort(f.version.SubInPort)

	f.InPort = f.version.InPort

	return f
}
func (f *ParseVersion) SetOutPort(port func(interface{})) {
	f.version.SetOutPort(port)
}

// ------------ ParseFlow:
type ParseFlow struct {
	flow *gparselib.ParseAll
	//	semantic    *SemanticCreateFlow
	flowLiteral *gparselib.ParseLiteral
	aspc        *gparselib.ParseSpace
	name        *ParseBigIdent
	spcComm1    *ParseSpaceComment
	openFlow    *gparselib.ParseLiteral
	spcComm2    *ParseSpaceComment
	connections *ParseConnections
	closeFlow   *gparselib.ParseLiteral
	spcComm3    *ParseSpaceComment
	InPort      func(interface{})
}

func NewParseFlow() *ParseFlow {
	f := &ParseFlow{}
	f.flow = gparselib.NewParseAll(parseData, setParseData)
	//	f.semantic = NewSemanticCreateFlow()
	f.flowLiteral = gparselib.NewParseLiteral(parseData, setParseData, "flow")
	f.aspc = gparselib.NewParseSpace(parseData, setParseData, false)
	f.name = NewParseBigIdent()
	f.spcComm1 = NewParseSpaceComment()
	f.openFlow = gparselib.NewParseLiteral(parseData, setParseData, "{")
	f.spcComm2 = NewParseSpaceComment()
	f.connections = NewParseConnections()
	f.closeFlow = gparselib.NewParseLiteral(parseData, setParseData, "}")
	f.spcComm3 = NewParseSpaceComment()

	//	f.flow.SetSemOutPort(f.semantic.InPort)
	//	f.semantic.SetOutPort(f.flow.SemInPort)
	f.flow.AppendSubOutPort(f.flowLiteral.InPort)
	f.flowLiteral.SetOutPort(f.flow.SubInPort)
	f.flow.AppendSubOutPort(f.aspc.InPort)
	f.aspc.SetOutPort(f.flow.SubInPort)
	f.flow.AppendSubOutPort(f.name.InPort)
	f.name.SetOutPort(f.flow.SubInPort)
	f.flow.AppendSubOutPort(f.spcComm1.InPort)
	f.spcComm1.SetOutPort(f.flow.SubInPort)
	f.flow.AppendSubOutPort(f.openFlow.InPort)
	f.openFlow.SetOutPort(f.flow.SubInPort)
	f.flow.AppendSubOutPort(f.spcComm2.InPort)
	f.spcComm2.SetOutPort(f.flow.SubInPort)
	f.flow.AppendSubOutPort(f.connections.InPort)
	f.connections.SetOutPort(f.flow.SubInPort)
	f.flow.AppendSubOutPort(f.closeFlow.InPort)
	f.closeFlow.SetOutPort(f.flow.SubInPort)
	f.flow.AppendSubOutPort(f.spcComm3.InPort)
	f.spcComm3.SetOutPort(f.flow.SubInPort)

	f.InPort = f.flow.InPort

	return f
}
func (f *ParseFlow) SetOutPort(port func(interface{})) {
	f.flow.SetOutPort(port)
}
