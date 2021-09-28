package decoratorparser

import (
	"regexp"
	"strings"
	"web-autogen/utils"
)

var parserRegex = regexp.MustCompile(`@([a-zA-Z]+)\((.+?)\)`)

type FunctionEntry struct {
	Name      string
	Arguments []string
}

type FunctionsResponse struct {
	FunctionCalls []FunctionEntry
}

// GetFirstCallByName, or return nil if no calls exist
func (f *FunctionsResponse) GetFirstCallByName(name string) *FunctionEntry {
	for _, f := range f.FunctionCalls {
		if f.Name == name {
			return &f
		}
	}
	return nil
}

func (f *FunctionsResponse) GetCallsByName(name string) []FunctionEntry {
	var s = make([]FunctionEntry, 0)
	for _, f := range f.FunctionCalls {
		if f.Name == name {
			s = append(s, f)
		}
	}
	return s
}

// GetFunctions will return a list of functions from the comment text
func GetFunctions(comments string) FunctionsResponse {
	result := make([]FunctionEntry, 0)
	var groups = parserRegex.FindAllStringSubmatch(comments, -1)
	for _, function := range groups {
		var functionName = function[1]
		var arguments = strings.Split(function[2], ",")
		for idx := range arguments {
			arguments[idx] = utils.StringTrimmer(arguments[idx])
		}
		result = append(result, FunctionEntry{
			Name:      functionName,
			Arguments: arguments,
		})
	}
	return FunctionsResponse{
		FunctionCalls: result,
	}
}
