package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/flowdev/gflowparser/data"
	"github.com/flowdev/gparselib"
	. "github.com/smartystreets/goconvey/convey"
)

func TestParseStatementEnd(t *testing.T) {
	p := NewParseStatementEnd()

	runTest(t, p, "empty", "", nil, 1)
	runTest(t, p, "no match 1", "baaa", nil, 1)
	runTest(t, p, "no match 2", " /* bla\n */ \t // com!\n \t \r\n/** blu ; */ ", nil, 1)
	runTest(t, p, "simple 1", ";", ";", 0)
	runTest(t, p, "simple 2", "\t;0", "\t;", 0)
	runTest(t, p, "simple 3", " /* bla */; _t", " /* bla */; ", 0)
	runTest(t, p, "simple 4", " // comment! \n ;lilalo", " // comment! \n ;", 0)
	runTest(t, p, "complex", " /* bla\n */ \t; // com!\n \t \r\n/** blu */ _t",
		" /* bla\n */ \t; // com!\n \t \r\n/** blu */ ", 0)

}

func TestParseSpaceComment(t *testing.T) {
	p := NewParseSpaceComment()

	runTest(t, p, "empty", "", "", 0)
	runTest(t, p, "no match", "baaa", "", 0)
	runTest(t, p, "simple 1", " i", " ", 0)
	runTest(t, p, "simple 2", "\t0", "\t", 0)
	runTest(t, p, "simple 3", " /* bla */ _t", " /* bla */ ", 0)
	runTest(t, p, "simple 4", " // comment! \n lilalo", " // comment! \n ", 0)
	runTest(t, p, "complex", " /* bla\n */ \t // com!\n \t \r\n/** blu */ _t",
		" /* bla\n */ \t // com!\n \t \r\n/** blu */ ", 0)
}

func TestParseOptSpc(t *testing.T) {
	p := NewParseOptSpc()

	runTest(t, p, "empty", "", "", 0)
	runTest(t, p, "no match", "baaa", "", 0)
	runTest(t, p, "simple 1", " i", " ", 0)
	runTest(t, p, "simple 2", "\t0", "\t", 0)
	runTest(t, p, "simple 3", " \t _t", " \t ", 0)
	runTest(t, p, "simple 4", " \n ", " ", 0)
}

func TestParseBigIdent(t *testing.T) {
	p := NewParseBigIdent()

	runTest(t, p, "empty", "", nil, 1)
	runTest(t, p, "no match 1", "baaa", nil, 1)
	runTest(t, p, "no match 2", "A", nil, 1)
	runTest(t, p, "simple 1", "Ab", "Ab", 0)
	runTest(t, p, "simple 2", "A0", "A0", 0)
	runTest(t, p, "simple 3", "Ab_cd", "Ab", 0)
	runTest(t, p, "simple 4", "AbcDef", "AbcDef", 0)
	runTest(t, p, "simple 5", "Abc123Def", "Abc123Def", 0)
}

func TestParseSmallIdent(t *testing.T) {
	p := NewParseSmallIdent()

	runTest(t, p, "empty", "", nil, 1)
	runTest(t, p, "no match 1", "ABCD", nil, 1)
	runTest(t, p, "no match 2", "123", nil, 1)
	runTest(t, p, "simple 1", "aB", "aB", 0)
	runTest(t, p, "simple 2", "a0", "a0", 0)
	runTest(t, p, "simple 3", "aB_CD", "aB", 0)
	runTest(t, p, "simple 4", "aBCdEF", "aBCdEF", 0)
	runTest(t, p, "simple 5", "aBC123dEF", "aBC123dEF", 0)
}

func runTest(t *testing.T, fp interface{}, name string, content string, ev interface{}, errCount int) {
	var mainData2 *data.MainData
	mainData := &data.MainData{}
	mainData.ParseData = gparselib.NewParseData(name, content)
	callWithInterfaceFunc(fp, "SetOutPort", func(dat interface{}) { mainData2 = dat.(*data.MainData) })
	callWithMainData(fp, "InPort", mainData)

	Convey("Parsing '"+name+"', ...", t, func() {
		Convey(`... should create a result.`, func() {
			So(mainData2.ParseData.Result, ShouldNotBeNil)
		})
		valueTest(mainData2.ParseData.Result.Value, ev)
		Convey(`... should give the right number of errors.`, func() {
			r := mainData2.ParseData.Result
			if errCount <= 0 {
				So(r.HasError(), ShouldBeFalse)
			} else {
				So(len(r.Feedback), ShouldEqual, errCount)
			}
		})
		Convey(`... should give the right error text.`, func() {
			errs := mainData2.ParseData.Result.Feedback
			if errCount > 0 {
				if len(errs) != errCount {
					So(printErrors(errs), ShouldBeBlank)
				} else {
					So(errs[errCount-1].Msg.String(), ShouldNotBeBlank)
				}
			}
		})
	})
}
func callWithMainData(p interface{}, method string, value *data.MainData) {
	fVal := getFunc(p, method)
	if fVal.IsValid() {
		fVal.Call([]reflect.Value{reflect.ValueOf(value)})
	}
}
func callWithInterfaceFunc(p interface{}, method string, value func(interface{})) {
	fVal := getFunc(p, method)
	if fVal.IsValid() {
		fVal.Call([]reflect.Value{reflect.ValueOf(value)})
	}
}
func printErrors(errs []*gparselib.FeedbackItem) string {
	result := ""
	for _, err := range errs {
		if err.Kind == gparselib.FeedbackError {
			result += err.String() + "\n"
		}
	}
	if result == "" {
		result = "<EMPTY>"
	}
	return result
}
func getFunc(p interface{}, method string) reflect.Value {
	pVal := reflect.ValueOf(p)
	fVal := pVal.MethodByName(method)
	if !fVal.IsValid() {
		fVal = pVal.Elem().FieldByName(method)
	}
	if !fVal.IsValid() {
		fmt.Printf("ERROR: Unable to find method or field '%s' in fp (type: %T, value: %#v)!!!\n", method, pVal, pVal)
	}
	return fVal
}
func valueTest(actual, expected interface{}) {
	if expected == nil {
		Convey(`... should create no value.`, func() {
			So(actual, ShouldBeNil)
		})
	} else {
		Convey(`... should create the right value.`, func() {
			So(fmt.Sprintf("%T", actual), ShouldEqual, fmt.Sprintf("%T", expected))
			So(spew.Sprintf("%v", actual), ShouldEqual, spew.Sprintf("%v", expected))
		})
	}
}
