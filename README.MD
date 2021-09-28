# Go Web Autogen

Extremely small and lightweight code generator for API routes and swagger. Inspired by ASP.NET

# Example

Learning by examples is probably the easiest way to get started. You can take a look at the projects in the `examples/` folder, but here's a quick summary:

Let's assume you have a project structured like this:
```
controllers/
    users_controller.go
models/
    users_model.go
main.go
```

In the `users_controller.go` file, you would have code like this:
```go
package controllers

// @HttpGet("/users/v1/get-user")
func GetUserById(id uint64) models.UserEntry {
    // Do whatever here
    return models.UserEntry{
        Id = 1,
        Username = "User123"
    }
}
```

You don't have to manually add swagger docs, or routes, or anything else. Just run `go run main.go` to generate your "autogen" folder. The generated code will look like this:
```go
package autogen_web

type AutoGenRegister struct {

}

// @Produce json
// @Router //users/v1/get-user [get]
// @Success 200 {object} models.UserEntry
func (a *AutoGenRegister) users_v1_getuser(app *Fiber.app) {
    app.Get("/users/v1/get-user", func(c *fiber.Ctx) error {
        response := controllers.GetUserById()
        return c.JSON(response)
    })
}

func (a* AutoGenRegister) Run(app *Fiber.app) {
    a.users_v1_getuser(app)
}

```

In your `main.go`, call `AutoGenRegister.Run(app)` with your fiber app before calling `app.Listen()`, like this:

```go
package main

import (
	autogen_web "project/autogen"
)

func main() {
	app := fiber.New()
	var registry = autogen_web.AutoGenRegister{}
	registry.Run(app)
	app.Listen(":3000")
}
```
