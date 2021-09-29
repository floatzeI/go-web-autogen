package decoratorparser

import (
	"fmt"
	"regexp"
	"strings"
	"web-autogen/utils"
)

type HttpMethod struct {
	HttpMethod         string
	Url                string
	IsControllerAction bool
}

var httpMethodRegex = regexp.MustCompile("Http(Get|Post|Patch|Delete|Put|Options|Head)")

func (p *DecoratorParser) ParseMethod() HttpMethod {
	for _, f := range p.functions.FunctionCalls {
		var result = httpMethodRegex.FindAllStringSubmatch(f.Name, -1)
		if len(result) == 0 {
			continue
		}
		method := result[0][1]
		url := f.Arguments[0]
		url = utils.StringTrimmer(url)
		url = strings.ReplaceAll(url, "{", ":")
		url = strings.ReplaceAll(url, "}", "")
		return HttpMethod{
			IsControllerAction: true,
			HttpMethod:         method,
			Url:                url,
		}
	}
	fmt.Println("[warning] no HTTP Method found for " + p.controller + "." + p.controllerAction)
	return HttpMethod{
		IsControllerAction: false,
	}
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
