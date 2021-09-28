package decoratorparser_test

import (
	"testing"
	"web-autogen/decoratorparser"
)

func TestParseMethod(t *testing.T) {
	var txt = "@HttpGet(\"/api/v1/testing\")"
	var expectedOutputMethod = "Get"
	var expectedOutputUrl = "/api/v1/testing"

	var p = decoratorparser.New(txt, "Test", "Test")
	var result = p.ParseMethod()
	if result.HttpMethod != expectedOutputMethod {
		t.Error("Method does not match: ", expectedOutputMethod, "!=", result.HttpMethod)
	}
	if result.Url != expectedOutputUrl {
		t.Error("Url does not match:", expectedOutputUrl, result.Url)
	}
}

func TestStringifyMethod(t *testing.T) {
	var result = decoratorparser.StringifyMethod(decoratorparser.HttpMethod{
		Url:        "/api/v1/test",
		HttpMethod: "Get",
	})
	var expected = "@Router /api/v1/test [get]"
	if result.Comment != expected {
		t.Error("Comment does not match", expected, result)
	}
}
