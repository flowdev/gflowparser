package data2svg

import (
	"bytes"
	"fmt"

	"github.com/flowdev/gflowparser/data"
	"github.com/flowdev/gflowparser/svg"
)

// Error messages.
const (
	errMsgDeclAndRef = "Circular flows aren't allowed yet, but the component " +
		"'%s' is declared here:\n%s\n... and referenced again here:\n%s"
	errMsg2Decls = "A component with the name '%s' is declared two times, " +
		"here:\n%s\n... and again here:\n%s"
	errMsgPartType = "Found illegal flow part type '%T' at index [%d, %d]"
	errMsgLoneComp = "Component reference with name '%s' without " +
		"input or output found:\n%s"
)

// Whereer can give a human readable description of a source position.
type Whereer interface {
	Where(pos int) string
}

type decl struct {
	name     string
	srcPos   int
	i, j     int
	svgOp    *svg.Op
	svgMerge *svg.Merge
	svgSplit *svg.Split
}

type merge struct {
	name   string
	srcPos int
	svg    *svg.Merge
}

type split struct {
	name   string
	srcPos int
}

// parserPartsToSVGData converts the parts of a data.Flow one to one into SVG
// diagram shapes.
// So this operation does only a simple translation but doesn't change the form
// of the part table.
// Components are special since they can be translated in 3 ways:
// 1. Into a decl struct if it is a declaration (the first occurence).
// 2. Into a merge if there are more parts before it.
// 3. Into a split if it is only used for a split.
func parserPartsToSVGData(flowDat data.Flow, w Whereer,
) (shapes [][]interface{}, decls map[string]*decl, clsts clusters, err error) {
	svgDat := make([][]interface{}, len(flowDat.Parts))
	decls = make(map[string]*decl)
	clsts = clusters(nil)

	for i, partLine := range flowDat.Parts {
		m := len(partLine) - 1
		svgLine := make([]interface{}, m+1)
		for j, part := range partLine {
			switch p := part.(type) {
			case data.Arrow:
				svgLine[j] = arrowToSVGData(p, j > 0, j < m)
			case data.Component:
				if dcl, ok := decls[p.Decl.Name]; ok {
					if !p.Decl.VagueType { // prevent double declaration
						return nil, nil, nil, fmt.Errorf(errMsg2Decls,
							dcl.name, w.Where(dcl.srcPos), w.Where(p.SrcPos))
					}
					if j > 0 { // we probably need a merge
						dcl.svgMerge.Size++
						clsts.addCluster(dcl.i, i)
						svgLine[j] = &merge{
							name:   dcl.name,
							srcPos: p.SrcPos,
							svg:    dcl.svgMerge,
						}
					} else if j < m { // we only need a split
						svgLine[j] = &split{name: dcl.name, srcPos: p.SrcPos}
					} else { // we don't need anything at all???!!!
						return nil, nil, nil, fmt.Errorf(errMsgLoneComp,
							dcl.name, w.Where(p.SrcPos))
					}
				} else {
					dcl := &decl{
						name:     p.Decl.Name,
						srcPos:   p.SrcPos,
						i:        i,
						j:        j,
						svgOp:    compToSVGData(p),
						svgMerge: &svg.Merge{ID: p.Decl.Name},
					}
					if j > 0 { // we might need a merge
						dcl.svgMerge.Size++
					}
					decls[dcl.name] = dcl
					svgLine[j] = dcl
				}
			default:
				panic(fmt.Sprintf(errMsgPartType, part, i, j))
			}
		}
		svgDat[i] = svgLine
	}
	return svgDat, decls, clsts, nil
}

func arrowToSVGData(arr data.Arrow, hasSrcOp, hasDstOp bool) *svg.Arrow {
	return &svg.Arrow{
		DataType: arrDataToSVGData(arr.Data),
		HasSrcOp: hasSrcOp, SrcPort: portToSVGData(arr.FromPort),
		HasDstOp: hasDstOp, DstPort: portToSVGData(arr.ToPort),
	}
}

func compToSVGData(comp data.Component) *svg.Op {
	plugs := make([]*svg.Plugin, len(comp.Plugins))
	for i, plug := range comp.Plugins {
		plugs[i] = pluginToSVGData(plug)
	}
	return &svg.Op{
		Main:    &svg.Rect{Text: compDeclToSVGData(comp.Decl)},
		Plugins: plugs,
	}
}

func compDeclToSVGData(decl data.CompDecl) []string {
	typ := typeToSVGData(decl.Type)
	if typ == decl.Name {
		return []string{decl.Name}
	}
	return []string{decl.Name, typ}
}

func pluginToSVGData(plug data.Plugin) *svg.Plugin {
	rects := make([]*svg.Rect, len(plug.Types))
	for i, typ := range plug.Types {
		rects[i] = &svg.Rect{
			Text: []string{typeToSVGData(typ)},
		}
	}
	return &svg.Plugin{
		Title: plug.Name,
		Rects: rects,
	}
}

func arrDataToSVGData(dat []data.Type) string {
	if len(dat) == 0 {
		return ""
	}
	buf := bytes.Buffer{}
	buf.WriteString("(")
	buf.WriteString(typeToSVGData(dat[0]))
	for i := 1; i < len(dat); i++ {
		buf.WriteString(", ")
		buf.WriteString(typeToSVGData(dat[i]))
	}
	buf.WriteString(")")
	return buf.String()
}

func portToSVGData(port *data.Port) string {
	if port == nil {
		return ""
	}
	if port.HasIndex {
		return fmt.Sprintf("%s[%d]", port.Name, port.Index)
	}
	return port.Name
}

func typeToSVGData(typ data.Type) string {
	if typ.Package != "" {
		return typ.Package + "." + typ.LocalType
	}
	return typ.LocalType
}

// handleSplits handles splits (and merges with splits).
// After this function finishes the following conditions hold true:
// - There are no *split shapes anymore.
// - *merge shapes are always at the end of their line.
// - *decl shapes can be anywhere but always before their merges.
func handleSplits(shapes [][]interface{}, decls map[string]*decl, clsts clusters,
) ([][]interface{}, clusters) {
	for i := len(shapes) - 1; i >= 0; i-- {
		sl := shapes[i]
		for j := len(sl) - 1; j >= 0; j-- {
			switch s := sl[j].(type) {
			case *merge:
				if j < len(sl)-1 { // add rest of line to split
					dcl := decls[s.name]
					dcl.svgSplit = prependShapeLine(dcl.svgSplit, sl[j+1:])
					shapes[i] = sl[:j+1]
					sl = shapes[i]
				}
			case *split: // add rest of line to split
				if j > 0 || len(sl) <= 1 {
					panic(fmt.Sprintf("illegal split at index[%d, %d]; row size: %d", i, j, len(sl)))
				}
				dcl := decls[s.name]
				dcl.svgSplit = prependShapeLine(dcl.svgSplit, sl[j+1:])
				shapes, clsts = deleteShapeLine(shapes, clsts, i)
			case *decl:
				if j < len(sl)-1 { // add rest of line to split
					s.svgSplit = prependShapeLine(s.svgSplit, sl[j+1:])
					shapes[i] = sl[:j+1]
					sl = shapes[i]
				}
				if s.svgSplit != nil && len(s.svgSplit.Shapes) == 1 { // no real split
					shapes[i] = append(sl, s.svgSplit.Shapes[0]...)
					sl = shapes[i]
				}
				if s.svgSplit != nil && len(s.svgSplit.Shapes) <= 1 { // no split at all
					s.svgSplit = nil
				}
			}
		}
	}
	return shapes, clsts
}
func prependShapeLine(split *svg.Split, sl []interface{}) *svg.Split {
	if split == nil || len(split.Shapes) == 0 {
		return &svg.Split{Shapes: [][]interface{}{sl}}
	}
	split.Shapes = append([][]interface{}{copyShapeLine(sl)}, split.Shapes...)
	return split
}
func copyShapeLine(sl []interface{}) []interface{} {
	result := make([]interface{}, len(sl))
	copy(result, sl)
	return result
}
func deleteShapeLine(shapes [][]interface{}, clsts clusters, i int,
) ([][]interface{}, clusters) {
	shapes = append(shapes[:i], shapes[i+1:]...)
	clsts = clsts.deleteLine(i)
	return shapes, clsts
}

// breakCircles replaces back pointing merges with simple *svg.Rects.
// After this function finishes the following conditions hold true:
// - There are no *merge shapes for a decl in the splits of the decl.
// - There are no *merge shapes for a decl in the same row of the decl.
func breakCircles(shapes [][]interface{}, seenDecls []string) {
	for _, sl := range shapes {
		breakCirclesInLine(sl, seenDecls)
	}
}
func breakCirclesInLine(sl []interface{}, seenDecls []string) {
	if len(sl) == 0 {
		return
	}
	for j, si := range sl {
		switch s := si.(type) {
		case *merge:
			if hasDecl(s.name, seenDecls) { // it's a backreference -> convert
				s.svg.Size--
				sl[j] = &svg.Rect{Text: []string{s.name}}
			}
		case *decl:
			newSeenDecls := append([]string{s.name}, seenDecls...)
			if s.svgSplit != nil {
				breakCircles(
					s.svgSplit.Shapes,
					newSeenDecls,
				)
			}
			if j < len(sl)-1 {
				breakCirclesInLine(
					sl[j+1:],
					newSeenDecls,
				)
			}
		}
	}
}
func hasDecl(name string, decls []string) bool {
	for _, d := range decls {
		if d == name {
			return true
		}
	}
	return false
}

// handleMerges handles merges by appending the decl after the last merge.
func handleMerges(allShapes [][]interface{}, myShapes [][]interface{},
) [][]interface{} {
	for i := len(myShapes) - 1; i >= 0; i-- {
		sl := myShapes[i]
		for j := len(sl) - 1; j >= 0; j-- {
			if s, ok := sl[j].(*decl); ok {
				if s.svgSplit != nil { // handle merges in splits
					s.svgSplit.Shapes = handleMerges(
						allShapes,
						s.svgSplit.Shapes,
					)
				}
				if s.svgMerge == nil || s.svgMerge.Size <= 0 ||
					(s.svgMerge.Size == 1 && j > 0) { // no merge
					s.svgMerge = nil
					continue
				}
				found := addDeclLineAfterLastMerge(
					allShapes,
					copyShapeLine(sl[j:]),
					s.name,
				)
				if !found {
					panic("Unable to find merge with ID: " + s.name)
				}
				if j > 0 {
					myShapes[i] = append(sl[:j], s.svgMerge)
					sl = myShapes[i]
				} else { // remove row (can only happen in allShapes)
					myShapes = append(myShapes[:i], myShapes[i+1:]...)
					allShapes = myShapes
				}
				s.svgMerge = nil // prevent endless loop
			}
		}
	}
	return myShapes
}

// addDeclLineAfterLastMerge adds the declLine directly after the last merge.
// It doesn't add a new line. This saves us a lot of headache:
// - No clusters are modified.
// - shapes itself doesn't change (only one of its rows might be modified).
func addDeclLineAfterLastMerge(shapes [][]interface{}, dl []interface{}, name string,
) (found bool) {
	for i := len(shapes) - 1; i >= 0; i-- {
		sl := shapes[i]
		for j := len(sl) - 1; j >= 0; j-- {
			switch s := sl[j].(type) {
			case *merge:
				if s.name == name { // this is the last merge
					if s.svg.Size <= 1 { // remove merge
						shapes[i] = append(sl[:j], dl...)
					} else {
						shapes[i] = append(sl, dl...)
					}
					return true
				}
			case *decl:
				if s.svgSplit != nil {
					found = addDeclLineAfterLastMerge(
						s.svgSplit.Shapes,
						dl,
						name,
					)
					if found {
						return true
					}
				}
			}
		}
	}
	return false
}

// addEmptyRows adds empty lines on the top level to form visible clusters.
func addEmptyRows(shapes [][]interface{}, clsts clusters) [][]interface{} {
	m := len(shapes) - 1
	for i := m; i >= 0; {
		mn, mx := clsts.getCluster(i)
		if mx < m {
			shapes = insertEmptyShapeRow(shapes, mx+1)
		}
		i = mn - 1
	}
	return shapes
}
func insertEmptyShapeRow(shapes [][]interface{}, i int) [][]interface{} {
	shapes = append(shapes, []interface{}{})
	for j := len(shapes) - 1; j > i; j-- {
		shapes[j] = shapes[j-1]
	}
	shapes[i] = []interface{}{}
	return shapes
}

// cleanSVGData replaces all special data structures with pure SVG ones.
// shapes itself doesn't change (only some of its rows might be appended to).
func cleanSVGData(shapes [][]interface{}) {
	for i, sl := range shapes {
		var svgSplit *svg.Split
		for j, si := range sl {
			switch s := si.(type) {
			case *merge:
				sl[j] = s.svg
			case *decl:
				sl[j] = s.svgOp
				if s.svgSplit != nil && len(s.svgSplit.Shapes) > 0 {
					cleanSVGData(s.svgSplit.Shapes)
					if j < len(sl)-1 {
						panic("decls with splits have to be at the end of their row!")
					}
					svgSplit = s.svgSplit
				}
			case *svg.Merge, *svg.Rect, *svg.Arrow, *svg.Op:
				// nothing to do
			default:
				panic(fmt.Sprintf("found unexpected shape type: %T", si))
			}
			if svgSplit != nil {
				shapes[i] = append(sl, svgSplit)
			}
		}
	}
}

// Convert converts a flow data structure (as generated by the parser) into a
// SVG diagram (as data structure).
// If the flow is invalid an error and no diagram data is returned.
func Convert(flow data.Flow, wh Whereer) (svg.Flow, error) {
	shapes, decls, clsts, err := parserPartsToSVGData(flow, wh)
	if err != nil {
		return svg.Flow{}, err
	}

	shapes, clsts = handleSplits(shapes, decls, clsts)

	breakCircles(shapes, nil)

	shapes = handleMerges(shapes, shapes)

	shapes = addEmptyRows(shapes, clsts)

	cleanSVGData(shapes)

	return svg.Flow{Shapes: shapes}, nil
}
