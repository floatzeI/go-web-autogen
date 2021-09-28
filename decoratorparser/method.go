package decoratorparser

import (
	"errors"
	"regexp"
	"strings"
	"web-autogen/utils"
)

type HttpMethod struct {
	HttpMethod string
	Url        string
}

var httpMethodRegex = regexp.MustCompile("Http(Get|Post|Patch|Delete|Put|Options|Head)")

func ParseMethod(functions FunctionsResponse) HttpMethod {
	for _, f := range functions.FunctionCalls {
		var result = httpMethodRegex.FindAllStringSubmatch(f.Name, -1)
		if len(result) == 0 {
			continue
		}
		method := result[0][1]
		url := f.Arguments[0]
		url = utils.StringTrimmer(url)
		return HttpMethod{
			HttpMethod: method,
			Url:        url,
		}
	}
	panic(errors.New("no HTTP Method found"))
}

type StringifyMethodResponse struct {
	Comment string
}

func StringifyMethod(method HttpMethod) StringifyMethodResponse {
	var str = "@Router " + method.Url + " [" + strings.ToLower(method.HttpMethod) + "]"
	return StringifyMethodResponse{
		Comment: str,
	}
}
