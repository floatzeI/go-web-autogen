// Package autogen_web
// This file was automatically generated by go-web-autogen
// It should not be edited manually
// Generated At: 2021-09-29T12:44:28-04:00
package autogen_web

import (
    items "project-one/models/items"
    controllers "project-one/controllers"

	"github.com/gofiber/fiber/v2"
)

type AutoGenRegister struct {
}

// @Produce json
// @Router /api/v1/items/{itemId}/info [get]
// @Success 200 {object} items.ItemEntry
// @Param itemId path int64 false "The itemId"
func get_apiv1itemsitemidinfo(app *fiber.App) {
	app.Get("/api/v1/items/{itemId}/info", func(c *fiber.Ctx) error {
		
		var result = controllers.NewOtherController().GetItemById(NewArgumentParser(c, "path", "itemId", false).GetInt64())
		
		return c.JSON(result)
	})
}

// @Produce json
// @Router /api/v1/items/create [post]
// @Success 200 {object} items.CreateResponse
// @Param request body items.CreateRequest true "The request"
// @Param skipIfAlreadyExists query bool false "The skipIfAlreadyExists"
func post_apiv1itemscreate(app *fiber.App) {
	app.Post("/api/v1/items/create", func(c *fiber.Ctx) error {
		
		var request_decoded items.CreateRequest
		if err := c.BodyParser(&request_decoded); err != nil {
            return err
        }

		var result = controllers.NewOtherController().CreateItem(request_decoded, NewArgumentParser(c, "query", "skipIfAlreadyExists", false).GetBool())
		
		return c.JSON(result)
	})
}

// @Produce json
// @Router /api/v1/users/{userId}/info [get]
// @Failure 400 {object} models.ErrorResponse "InvalidUserId: UserId is invalid"
// @Success 200 {object} users.GetUserById
// @Param userId path int64 true "The userId"
func get_apiv1usersuseridinfo(app *fiber.App) {
	app.Get("/api/v1/users/{userId}/info", func(c *fiber.Ctx) error {
		
		var result = controllers.GetUserById(NewArgumentParser(c, "path", "userId", true).GetInt64())
		
		return c.JSON(result)
	})
}

// @Produce json
// @Router /api/v1/users/username [get]
// @Success 200 {object} users.GetUserById
// @Param username query string false "The username"
func get_apiv1usersusername(app *fiber.App) {
	app.Get("/api/v1/users/username", func(c *fiber.Ctx) error {
		
		var result = controllers.GetUserByName(NewArgumentParser(c, "query", "username", false).GetString())
		
		return c.JSON(result)
	})
}

func (a *AutoGenRegister) Run(app *fiber.App) {
    get_apiv1itemsitemidinfo(app)
    post_apiv1itemscreate(app)
    get_apiv1usersuseridinfo(app)
    get_apiv1usersusername(app)
}
