package utils_test

import (
	"testing"
	"web-autogen/utils"
)

func TestStringTimmer(t *testing.T) {
	var input = "   \"example string here\"  "
	var expectedOutput = "example string here"
	var result = utils.StringTrimmer(input)
	if result != expectedOutput {
		t.Error("Expected", expectedOutput, "got", result)
	}
}
