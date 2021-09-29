// This is an example using the DI/Struct approach. It's a bit more complicated but should be fairly simple to understand.
// This improves testability at the expensive of greater technical debt (since if you wanted to leave go-web-autogen, you would
// have to make your own DI framework)

package controllers

import (
	"project-one/models/items"
	"project-one/services"
)

// @Controller("/api/v1/items/")
type OtherController struct {
	items services.IItemsService
}

func NewOtherController(itemsService services.IItemsService) *OtherController {
	return &OtherController{
		items: itemsService,
	}
}

// @HttpGet("{itemId}/info")
func (c *OtherController) GetItemById(itemId int64) items.ItemEntry {
	return c.items.GetItemById(itemId)
}

// @HttpPost("create")
func (c *OtherController) CreateItem(request items.CreateRequest, skipIfAlreadyExists bool) items.CreateResponse {
	return items.CreateResponse{
		Id: 123,
	}
}
