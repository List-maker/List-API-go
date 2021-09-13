package items

import "listes_back/src/lists"

type Item struct {
	Id       uint64 `json:"id"`
	ParentId uint64 `json:"parent_id"`
	Name     string `json:"name"`
	Checked  bool   `json:"checked"`
}

func (item Item) CanEdit(userId uint64) bool {
	parentList, found := lists.LoadListById(item.ParentId)
	if !found {
		return false
	}
	return parentList.CanEdit(userId)
}

func (item Item) CanView(userId uint64) bool {
	parentList, found := lists.LoadListById(item.ParentId)
	if !found {
		return false
	}
	return parentList.CanView(userId)
}

func (item Item) Export() interface{} {
	return item
}
