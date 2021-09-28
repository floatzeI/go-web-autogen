// Package autogen_web
// This file is autogenerated by go-web-autogen. Do not manually edit it.

package autogen_web

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type BadRequestError struct {
	Message   string
	Parameter string
	Type      string
}

func ArgNotProvidedError(param string, pType string, t string) *BadRequestError {
	return &BadRequestError{
		Type:      t,
		Parameter: param,
		Message:   "Required parameter was not provided. Name=" + param,
	}
}

func InvalidArgumentError(param string, pType string, t string) *BadRequestError {
	return &BadRequestError{
		Type:      t,
		Parameter: param,
		Message:   "Parameter was not correct format. Name=" + param,
	}
}

type ArgumentParser struct {
	ctx           *fiber.Ctx
	pType         string
	parameterName string
	required      bool
	rawValue      string
}

func NewArgumentParser(ctx *fiber.Ctx, pType string, parameterName string, required bool) *ArgumentParser {
	var v string
	if pType == "query" {
		v = ctx.Query(parameterName)
	} else if pType == "path" {
		v = ctx.Params(parameterName)
	} else {
		panic(errors.New("Unsupported provider type: " + pType))
	}

	if required && v == "" {
		panic(InvalidArgumentError(parameterName, pType, ""))
	}

	return &ArgumentParser{
		ctx:           ctx,
		pType:         pType,
		parameterName: parameterName,
		required:      required,
		rawValue:      v,
	}
}

func (p *ArgumentParser) GetInt() int {
	if p.rawValue == "" {
		return 0
	}
	parse, parseErr := strconv.Atoi(p.rawValue)

	if parseErr != nil {
		if p.required {
			panic(InvalidArgumentError(p.parameterName, p.pType, "int"))
		}
		return 0
	}
	return parse
}

func (p *ArgumentParser) GetInt64() int64 {
	if p.rawValue == "" {
		return 0
	}
	parse, parseErr := strconv.ParseInt(p.rawValue, 10, 64)

	if parseErr != nil {
		if p.required {
			panic(InvalidArgumentError(p.parameterName, p.pType, "int"))
		}
		return 0
	}
	return parse
}

func (p *ArgumentParser) GetUint() uint {
	if p.rawValue == "" {
		return 0
	}
	parse, parseErr := strconv.ParseUint(p.rawValue, 10, 32)

	if parseErr != nil {
		if p.required {
			panic(InvalidArgumentError(p.parameterName, p.pType, "int"))
		}
		return 0
	}
	return uint(parse)
}

func (p *ArgumentParser) GetUint64() uint64 {
	if p.rawValue == "" {
		return 0
	}
	parse, parseErr := strconv.ParseUint(p.rawValue, 10, 32)

	if parseErr != nil {
		if p.required {
			panic(InvalidArgumentError(p.parameterName, p.pType, "int"))
		}
		return 0
	}
	return parse
}

func (p *ArgumentParser) GetString() string {
	return p.rawValue
}

func (p *ArgumentParser) GetBool() bool {
	str := strings.ToLower(p.rawValue)
	if str == "true" {
		return true
	}
	if str == "false" {
		return false
	}
	if p.required {
		panic(ArgNotProvidedError(p.parameterName, p.pType, "bool"))
	}
	return false
}
