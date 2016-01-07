package main

import (
	"fmt"
	"github.com/flowdev/gflowparser/data"
	"github.com/flowdev/gflowparser/parser"
	"github.com/flowdev/gparselib"
)

func main() {
	var mainData2 *data.MainData
	mainData := &data.MainData{}
	mainData.ParseData = gparselib.NewParseData("testfile", "smallIdent")

	p := parser.NewParseSmallIdent()
	p.SetOutPort(func(dat interface{}) { mainData2 = dat.(*data.MainData) })
	p.InPort(mainData)

	fmt.Printf("Result.Value: %#v\n", mainData2.ParseData.Result.Value)
}
