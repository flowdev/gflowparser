package data

import (
	"strings"
)

type PortData struct {
	Name     string
	CapName  string
	HasIndex bool
	Index    int
	SrcPos   int
}

func NewPort(name string, srcPos int) *PortData {
	return newPort(name, false, 0, srcPos)
}
func NewIdxPort(name string, idx int, srcPos int) *PortData {
	return newPort(name, true, idx, srcPos)
}
func CopyPort(port *PortData, srcPos int) *PortData {
	return &PortData{port.Name, port.CapName, port.HasIndex, port.Index, srcPos}
}
func DefaultInPort(srcPos int) *PortData {
	return NewPort("in", srcPos)
}
func DefaultOutPort(srcPos int) *PortData {
	return NewPort("out", srcPos)
}
func newPort(name string, hasIdx bool, idx int, srcPos int) *PortData {
	capName := name
	if len(name) > 0 {
		capName = strings.ToUpper(name[0:1]) + name[1:]
	}
	return &PortData{name, capName, hasIdx, idx, srcPos}
}

type PortPair struct {
	InPort  *PortData
	OutPort *PortData
}

type Operation struct {
	Name      string
	Type      string
	InPorts   []*PortData
	OutPorts  []*PortData
	SrcPos    int
	PortPairs []*PortPair
}

type Connection struct {
	FromOp       *Operation
	FromPort     *PortData
	DataType     string
	ShowDataType bool
	ToOp         *Operation
	ToPort       *PortData
}

// Port is the semantic representation of a port.
type Port struct {
	Name     string
	HasIndex bool
	Index    int
}

// Component is the semantic representation of a component.
type Component struct {
	Decl    CompDecl
	Plugins []NameNTypes
}

// NameNTypes is the semantic representation of a Name plus a list of types.
type NameNTypes struct {
	Name  string
	Types []Type
}

// CompDecl is the semantic representation of an operation declaration.
type CompDecl struct {
	Name      string
	Type      Type
	VagueType bool
}

// Type is the semantic representation of a type declaration.
type Type struct {
	Package   string
	LocalType string
}
