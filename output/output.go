package output

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/flowdev/gflowparser/data"
)

var AllowedFormat = map[string]bool{
	"dot": true,
	"go":  true,
}
var specialFormat = map[string]bool{
	"dot": true,
}

// ------------ ProduceFormats:
// input:  *data.MainData{SelectedFormats, FlowFile}
// output: *data.MainData{SelectedFormats, FlowFile, CurrentFormat} to the correct output port.
type ProduceFormats struct {
	outPort        func(*data.MainData)
	specialOutPort func(*data.MainData)
}

func NewOutputFormats() *ProduceFormats {
	return &ProduceFormats{}
}
func (op *ProduceFormats) InPort(md *data.MainData) {
	for _, format := range md.SelectedFormats {
		md.CurrentFormat = format
		if specialFormat[format] {
			op.specialOutPort(md)
		} else {
			op.outPort(md)
		}
	}
}
func (op *ProduceFormats) SetOutPort(port func(*data.MainData)) {
	op.outPort = port
}
func (op *ProduceFormats) SetSpecialOutPort(port func(*data.MainData)) {
	op.specialOutPort = port
}

// ------------ FillPortPairs:
// input:  flows with operations with ports but without port pairs
// output: flows with operations with port pairs filled
type FillPortPairs struct {
	outPort func(*data.MainData)
}

func NewFillPortPairs() *FillPortPairs {
	return &FillPortPairs{}
}
func (op *FillPortPairs) InPort(md *data.MainData) {
	for _, flow := range md.FlowFile.Flows {
		for _, op := range flow.Ops {
			fillPortPairs4Op(op)
		}
	}
	op.outPort(md)
}
func (op *FillPortPairs) SetOutPort(port func(*data.MainData)) {
	op.outPort = port
}

func fillPortPairs4Op(op *data.Operation) {
	l := len(op.InPorts)
	m := len(op.OutPorts)
	n := max(l, m)
	portPairs := make([]*data.PortPair, n)
	for i := 0; i < n; i++ {
		p := &data.PortPair{}
		if i < l {
			p.InPort = op.InPorts[i]
		}
		if i < m {
			p.OutPort = op.OutPorts[i]
		}
		portPairs[i] = p
	}
}
func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

// ------------ FillTemplate:
// input:  flows with operations with ports but without port pairs
// output: flows with operations with port pairs filled
const templateDot string = `digraph {{.flow.Name}} {
{{if .horizontal}}  rankdir=LR;{{end}}
  node [shape=Mrecord,style=filled,fillcolor="#428bca",rank=same];

  {{range .flow.Ops -}}
    {{.Name}} [label="{{.Name}}\n({{.Type}})
    {{- range .PortPairs -}}
      |{ {{with .InPort}}<{{.Name}}{{if .HasIndex}}.{{.Index}}{{end}}> {{.Name}}{{if .HasIndex}}.{{.Index}}{{end}}{{end -}}
	  |{{with .OutPort}}<{{.Name}}{{if .HasIndex}}.{{.Index}}{{end}}> {{.Name}}{{if .HasIndex}}.{{.Index}}{{end}}{{end}} }
	{{- end}}"] ;
  {{end}}

  node [shape=plaintext,style=plain,rank=same];

  {{range .flow.Conns -}}
    {{if .FromOp}}{{.FromOp.Name}}:{{with .FromPort}}{{.Name}}{{if .HasIndex}}.{{.Index}}{{end}}{{end}} {{end -}}
	{{else}}{{with .FromPort}}"{{.Name}}{{if .HasIndex}}.{{.Index}}{{end}}"{{end}} {{end -}}
    ->
	{{- if .ToOp}} {{.ToOp.Name}}:{{with .ToPort}}{{.Name}}{{if .HasIndex}}.{{.Index}}{{end}}{{end}}{{end}}
	{{else}} {{with .ToPort}}"{{.Name}}{{if .HasIndex}}.{{.Index}}{{end}}"{{end}}{{end}}
	{{- if .ShowDataType}} [label="{{.DataType}}"]{{end}} ;
  {{end}}
}`
const templateGo string = `type {{.flow.Name}} struct {
  {{range .flow.Conns -}}
    {{if .FromOp}}{{.FromOp.Name}}:{{with .FromPort}}{{.Name}}{{if .HasIndex}}.{{.Index}}{{end}}{{end}} {{end -}}
	{{else}}{{with .FromPort}}"{{.Name}}{{if .HasIndex}}.{{.Index}}{{end}}"{{end}} {{end -}}
    ->
	{{- if .ToOp}} {{.ToOp.Name}}:{{with .ToPort}}{{.Name}}{{if .HasIndex}}.{{.Index}}{{end}}{{end}}{{end}}
	{{else}} {{with .ToPort}}"{{.Name}}{{if .HasIndex}}.{{.Index}}{{end}}"{{end}}{{end}}
	{{- if .ShowDataType}} [label="{{.DataType}}"]{{end}} ;
  {{end}}
}

package main


{{with .flow -}}
type {{.Name}} struct {
{{#operations}} {{name}} *{{type}}
{{/operations}}
{{#connections}}{{^fromOp}}	{{fromPort.capName}}Port func({{dataType}})
{{/fromOp}}{{/connections}}
}

func New{{name}}() *{{name}} {
    f := &{{name}}{}
{{#operations}} f.{{name}} = New{{type}}()
{{/operations}}

{{#connections}}{{#fromOp}}{{#toOp}}    f.{{fromOp.name}}.Set{{fromPort.capName}}Port({{#fromPort.hasIndex}}{{fromPort.index}}, {{/fromPort.hasIndex}}f.{{toOp.name}}.{{toPort.capName}}Port{{#toPort.hasIndex}}[{{toPort.index}}]{{/toPort.hasIndex}})
{{/toOp}}{{/fromOp}}{{/connections}}

{{#connections}}{{^fromOp}}	f.{{fromPort.capName}}Port = f.{{toOp.name}}.{{toPort.capName}}Port
{{/fromOp}}{{/connections}}

    return f
}
{{range .Connections}}{{if not .ToOp}}func (f *{{$.flow.Name}}) Set{{.ToPort.CapName}}Port(port func({{.DataType}})) {
	f.{{.FromOp.Name}}.Set{{.FromPort.CapName}}Port(port)
}
{{end}}{{end}}

{{end}}
`

type FillTemplate struct {
	outPort   func(*data.MainData)
	templates map[string]*template.Template
}

func NewFillTemplate() *FillTemplate {
	// TODO: Compile all templates:
	tmpls := make(map[string]*template.Template)
	tmpls["dot"] = template.Must(template.New("dot").Parse(templateDot))
	tmpls["go"] = template.Must(template.New("go").Parse(templateGo))
	return &FillTemplate{templates: tmpls}
}
func (op *FillTemplate) InPort(md *data.MainData) {
	tplData := map[string]interface{}{
		"horizontal": md.Horizontal,
	}
	t := op.templates[md.CurrentFormat]
	var b bytes.Buffer // A Buffer needs no initialization.
	fmt.Fprintf(&b, "world!")
	for i, flow := range md.FlowFile.Flows {
		if i > 0 {
			b.Write([]byte("\n"))
		}
		err := t.Execute(os.Stdout, tplData)
		if err != nil {
			// TODO: use error port instead!
			log.Println("executing template:", err)
		}
	}
	op.outPort(md)
}
func (op *FillTemplate) SetOutPort(port func(*data.MainData)) {
	op.outPort = port
}
