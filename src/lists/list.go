package lists

import (
	"listes_back/src/database"
	"listes_back/src/utils"
)

type List struct {
	Id      uint64   `json:"id"`
	Name    string   `json:"name"`
	Items   []uint64 `json:"items"`
	Editors []uint64 `json:"editors"`
	Viewers []uint64 `json:"viewers"`
}

func (list List) CanEdit(userId uint64) bool {
	return utils.ContainsUint64(list.Editors, userId)
}

func (list List) CanView(userId uint64) bool {
	return list.CanEdit(userId) || utils.ContainsUint64(list.Viewers, userId)
}

func (list List) ExportFor(userId uint64) interface{} {
	if utils.ContainsUint64(list.Editors, userId) {
		return list
	}
	return struct {
		Id      uint64   `json:"id"`
		Name    string   `json:"name"`
		Items   []uint64 `json:"items"`
		Members []uint64 `json:"members"`
	}{
		Id:      list.Id,
		Name:    list.Name,
		Items:   list.Items,
		Members: append(list.Editors, list.Viewers...),
	}
}

func (list *List) SetItems(items []uint64) error {
	err := updateList(list.Id, "items", database.Uint64Slice(items))
	if err != nil {
		return err
	}
	list.Items = items
	return nil
}

func (list *List) SetEditors(editors []uint64) error {
	err := updateList(list.Id, "editors", database.Uint64Slice(editors))
	if err != nil {
		return err
	}
	list.Editors = editors
	return nil
}

func (list *List) SetViewers(viewers []uint64) error {
	err := updateList(list.Id, "viewers", database.Uint64Slice(viewers))
	if err != nil {
		return err
	}
	list.Viewers = viewers
	return nil
}
