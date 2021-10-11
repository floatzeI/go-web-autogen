package decoratorparser_test

import (
	"testing"

	"github.com/floatzeI/go-web-autogen/decoratorparser"
)

func TestParseResponse(t *testing.T) {
	comments := `@HttpGet("/api/example")
@Response(400, "Invalid Request", "Models.Error")`
	var parser = decoratorparser.New(comments, "Test", "Test")
	parser.GetFunctions()
	result := parser.ParseResponse()
	if len(result) != 1 {
		t.Error("Expected len=1, got len=", len(result))
	}
	entry := result[0]
	if entry.Description != "Invalid Request" {
		t.Error("Expected specified description, got", entry.Description)
	}
	if entry.Model != "Models.Error" {
		t.Error("Expected Models.Error model, got", entry.Model)
	}
	if entry.StatusCode != 400 {
		t.Error("Expected 400 status code, got", entry.StatusCode)
	}
}

func TestStringifyResponse(t *testing.T) {
	// 200
	entry := decoratorparser.ResponseEntry{
		StatusCode:  200,
		Description: "Desc Here",
		Model:       "models.Success",
	}
	expectedResult := `@Success 200 {object} models.Success "Desc Here"`
	var result = decoratorparser.StringifyResponse(entry)
	if result != expectedResult {
		t.Error("Unexpected result: ", result, "!=", expectedResult)
	}
	// 400 (failure)
	entry.StatusCode = 400
	entry.Model = "models.Error"
	expectedResult = `@Failure 400 {object} models.Error "Desc Here"`
	result = decoratorparser.StringifyResponse(entry)
	if result != expectedResult {
		t.Error("Unexpected result: ", result, "!=", expectedResult)
	}
}
