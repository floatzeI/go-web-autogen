package decoratorparser_test

import (
	"testing"

	"github.com/floatzeI/go-web-autogen/decoratorparser"
)

func TestFunction(t *testing.T) {
	// parse correct
	comments := `@FunctionCallOne("string", "otherString", true, 123)
@FunctionCallNoArgs()`
	result := decoratorparser.New(comments, "test", "test").GetFunctions()
	if len(result.FunctionCalls) != 2 {
		t.Error("Bad functioncalls len=", len(result.FunctionCalls))
	}
	var callOne decoratorparser.FunctionEntry
	var callTwo decoratorparser.FunctionEntry
	for _, item := range result.FunctionCalls {
		if item.Name == "FunctionCallOne" {
			callOne = item
		} else if item.Name == "FunctionCallNoArgs" {
			callTwo = item
		}
	}
	if len(callOne.Arguments) != 4 {
		t.Error("FunctionCallOne should have 4 arguments, but len=", len(callOne.Arguments))
	}

	if len(callTwo.Arguments) != 0 {
		t.Error("FunctionCallNoArgs() should have no arguments, but args len=", len(callTwo.Arguments))
	}

	if callOne.Arguments[0] != "string" {
		t.Error(callOne.Arguments[0], "!=", "string")
	}
	if callOne.Arguments[1] != "otherString" {
		t.Error(callOne.Arguments[0], "!=", "otherString")
	}
	if callOne.Arguments[2] != "true" {
		t.Error(callOne.Arguments[0], "!=", "true")
	}
	if callOne.Arguments[3] != "123" {
		t.Error(callOne.Arguments[0], "!=", "123")
	}
	// try to get a call by name
	var firstCall = result.GetFirstCallByName("FunctionCallNoArgs")
	if firstCall == nil {
		t.Error("Nil firstCall")
		return
	}
	if firstCall.Name != "FunctionCallNoArgs" {
		t.Error(firstCall.Name, "!=", "FunctionCallNoArgs")
	}
}
