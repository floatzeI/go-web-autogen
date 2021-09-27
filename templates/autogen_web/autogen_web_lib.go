package autogen_web

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type ParameterHelper struct {
	ctx *fiber.Ctx
}

func (p *ParameterHelper) GetIntQuery(parameterName string) int {
	var v = p.ctx.Query(parameterName)
	if v == "" {
		return 0
	}
	parse, parseErr := strconv.Atoi(v)
	if parseErr != nil {
		return 0
	}
	return parse
}