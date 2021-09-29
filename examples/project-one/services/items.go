package services

import "project-one/models/items"

// @ResolveScoped("IItemsService")
type ItemsService struct{}

func NewItemsService() IItemsService {
	return &ItemsService{}
}

func (i *ItemsService) GetItemById(id int64) items.ItemEntry {
	var dest items.ItemEntry
	pretendThisIsADbCall("select id from items where id = @id limit 1", map[string]interface{}{
		"id": id,
	}, &dest)
	return dest
}

// Pretend this is a MustExec()-type function, so the caller doesn't need "if err != nil"
func pretendThisIsADbCall(sql string, sqlParams interface{}, dest interface{}) {
	// ...
}
