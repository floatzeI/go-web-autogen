package services

import "project-one/models/items"

type IItemsService interface {
	GetItemById(id int64) items.ItemEntry
}
