package data

// Flow is the semantic representation of a complete flow.
type Flow struct {
	Parts [][]interface{}
}

// Arrow is the semantic representation of a flow arrow including data type and
// ports.
type Arrow struct {
	FromPort *Port
	ToPort   *Port
	Data     []Type
	SrcPos   int
}

// Port is the semantic representation of a port.
type Port struct {
	Name     string
	HasIndex bool
	Index    int
	SrcPos   int
}

// Component is the semantic representation of a component.
type Component struct {
	Decl    CompDecl
	Plugins []Plugin
	SrcPos  int
}

// Plugin is the semantic representation of a component plugin.
type Plugin struct {
	Name   string
	Types  []Type
	SrcPos int
}

// CompDecl is the semantic representation of an component declaration.
type CompDecl struct {
	Name      string
	Type      Type
	VagueType bool
	SrcPos    int
}

// Type is the semantic representation of a type declaration.
type Type struct {
	Package   string
	LocalType string
	SrcPos    int
}
