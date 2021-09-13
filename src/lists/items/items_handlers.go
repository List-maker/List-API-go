package items

import (
	"listes_back/src/lists"
	"listes_back/src/users"
	"listes_back/src/utils"
	"net/http"
)

func CreateItem(w http.ResponseWriter, r *http.Request, user users.User) {
	listId, found, valid := utils.ExtractUintFromRequest("list_id", r)
	if !found {
		utils.Prettier(w, "no id provided", nil, http.StatusBadRequest)
		return
	}
	if !valid {
		utils.Prettier(w, "invalid list id", nil, http.StatusBadRequest)
		return
	}

	list, found := lists.LoadListById(listId)
	if !found {
		utils.Prettier(w, "list not found", nil, http.StatusBadRequest)
		return
	}

	if !list.CanEdit(user.Id) {
		utils.Prettier(w, "you do not have the permission !", nil, http.StatusUnauthorized)
		return
	}

	_ = r.ParseForm()
	name := r.FormValue("name")
	if err := utils.IsNameValid(name); err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	newItemId, err := createItem(listId, name)
	if err != nil {
		utils.Prettier(w, "failed to create item", nil, http.StatusInternalServerError)
		utils.PrintError(err)
		return
	}

	err = list.SetItems(append(list.Items, newItemId))
	if err != nil {
		utils.Prettier(w, "failed add item to the list", nil, http.StatusInternalServerError)
		utils.PrintError(err)
		return
	}

	utils.Prettier(w, "item created !", newItemId, http.StatusOK)
}

func GetItem(w http.ResponseWriter, r *http.Request, user users.User) {
	item, err := LoadItemFromRequest(r)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	if !item.CanView(user.Id) {
		utils.Prettier(w, "you do not have the permission", nil, http.StatusUnauthorized)
		return
	}

	utils.Prettier(w, "item informations", item.Export(), http.StatusOK)
}

func UpdateItem(w http.ResponseWriter, r *http.Request, user users.User) {
	item, err := LoadItemFromRequest(r)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	if !item.CanEdit(user.Id) {
		utils.Prettier(w, "you do not have the permission", nil, http.StatusUnauthorized)
		return
	}

	_ = r.ParseForm()
	name := r.FormValue("name")
	if name == "" {
		utils.Prettier(w, "no name provided", nil, http.StatusBadRequest)
		return
	}
	if name == item.Name {
		utils.Prettier(w, "the new name cannot be the same as the old !", nil, http.StatusBadRequest)
		return
	}
	if err = utils.IsNameValid(name); err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	err = updateItem(item.Id, "name", name)
	if err != nil {
		utils.Prettier(w, "failed to update item", nil, http.StatusInternalServerError)
		return
	}

	item.Name = name
	utils.Prettier(w, "item updated successfully !", item.Export(), http.StatusOK)
}

func CheckItem(w http.ResponseWriter, r *http.Request, user users.User) {
	item, err := LoadItemFromRequest(r)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	if !item.CanEdit(user.Id) {
		utils.Prettier(w, "you do not have the permission", nil, http.StatusUnauthorized)
		return
	}

	newCheckState := !item.Checked

	err = updateItem(item.Id, "checked", newCheckState)
	if err != nil {
		utils.Prettier(w, "failed to update item", nil, http.StatusInternalServerError)
		return
	}

	item.Checked = newCheckState
	utils.Prettier(w, "item updated successfully !", item.Export(), http.StatusOK)
}

func DeleteItem(w http.ResponseWriter, r *http.Request, user users.User) {
	item, err := LoadItemFromRequest(r)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	list, found := lists.LoadListById(item.ParentId)
	if !found {
		utils.Prettier(w, "list not found", nil, http.StatusBadRequest)
		return
	}

	if !list.CanEdit(user.Id) {
		utils.Prettier(w, "you do not have the permission !", nil, http.StatusUnauthorized)
		return
	}

	err = deleteItem(item.Id)
	if err != nil {
		utils.Prettier(w, "failed to delete item", nil, http.StatusInternalServerError)
		return
	}

	err = list.SetItems(utils.RemoveFromUint64Slice(list.Items, item.Id))
	if err != nil {
		utils.Prettier(w, "failed to remove item from list", nil, http.StatusInternalServerError)
		return
	}

	utils.Prettier(w, "item deleted successfully !", nil, http.StatusOK)
}
