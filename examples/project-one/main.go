package main

// @title Example Project 1
// @version 1.0
// @description Example description 1

import (
	autogen_web "project-one/autogen"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Static("/", "./public")
	var registry = autogen_web.AutoGenRegister{}

	registry.Run(app)
	app.Listen(":4012")
}
