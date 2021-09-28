package decoratorparser

import (
	"errors"
	"strconv"
	"web-autogen/utils"
)

type ResponseEntry struct {
	StatusCode  int
	Description string
	Model       string
}

// ParseResponse parses all of the comments for the function, and returns a slice of response entries (if there are any)
func ParseResponse(calls FunctionsResponse) []ResponseEntry {
	resp := make([]ResponseEntry, 0)
	var possibleResponses = calls.GetCallsByName("Response")
	for _, entry := range possibleResponses {
		params := entry.Arguments
		if len(entry.Arguments) < 1 {
			panic(errors.New("Incorrect argument count for Response(). Len=" + strconv.Itoa(len(entry.Arguments))))
		}
		statusCode, err := strconv.Atoi(params[0])
		if err != nil {
			panic(errors.New("Unexpected status code: " + entry.Arguments[0] + ". Error = " + err.Error()))
		}
		var Description = ""
		var Model = ""

		if len(params) > 1 {
			Description = utils.StringTrimmer(params[1])
		}
		if len(params) > 2 {
			Model = utils.StringTrimmer(params[2])
		}

		resp = append(resp, ResponseEntry{
			StatusCode:  statusCode,
			Description: Description,
			Model:       Model,
		})
	}

	return resp
}

func StringifyResponse(resp ResponseEntry) string {
	var str = "@"
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		str += "Success"
	} else {
		str += "Failure"
	}
	str += " " + strconv.Itoa(resp.StatusCode) + " {object} " + resp.Model
	if resp.Description != "" && resp.Model != "" {
		str += " " + "\"" + resp.Description + "\""
	}
	return str
}
