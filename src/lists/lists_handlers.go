package lists

import (
	"listes_back/src/users"
	"listes_back/src/utils"
	"net/http"
)

func CreateList(w http.ResponseWriter, r *http.Request, user users.User) {
	_ = r.ParseForm()
	listName := r.FormValue("name")
	if listName == "" {
		utils.Prettier(w, "no name provided", nil, http.StatusBadRequest)
		return
	}
	if err := utils.IsNameValid(listName); err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	listId, err := createList(user.Id, listName)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	utils.Prettier(w, "list created !", listId, http.StatusOK)
}

func GetList(w http.ResponseWriter, r *http.Request, user users.User) {
	listId, found, valid := utils.ExtractUintFromRequest("id", r)
	if !found {
		utils.Prettier(w, "no id provided", nil, http.StatusBadRequest)
		return
	}
	if !valid {
		utils.Prettier(w, "invalid list id", nil, http.StatusBadRequest)
		return
	}

	list, found := LoadListById(listId)
	if !found {
		utils.Prettier(w, "list not found", nil, http.StatusBadRequest)
		return
	}

	if !list.CanView(user.Id) {
		utils.Prettier(w, "you do not have the permission !", nil, http.StatusUnauthorized)
		return
	}

	utils.Prettier(w, "list informations", list.ExportFor(user.Id), http.StatusOK)
}

func UpdateList(w http.ResponseWriter, r *http.Request, user users.User) {
	listId, found, valid := utils.ExtractUintFromRequest("id", r)
	if !found {
		utils.Prettier(w, "no id provided", nil, http.StatusBadRequest)
		return
	}
	if !valid {
		utils.Prettier(w, "invalid list id", nil, http.StatusBadRequest)
		return
	}

	list, found := LoadListById(listId)
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
	if name == "" {
		utils.Prettier(w, "no name provided", nil, http.StatusBadRequest)
		return
	}
	if name == list.Name {
		utils.Prettier(w, "the new name cannot be the same as the old !", nil, http.StatusBadRequest)
		return
	}
	if err := utils.IsNameValid(name); err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	err := updateList(listId, "name", name)
	if err != nil {
		utils.Prettier(w, "failed to update the list", nil, http.StatusInternalServerError)
		utils.PrintError(err)
		return
	}

	list.Name = name
	utils.Prettier(w, "list updated successfully", list.ExportFor(user.Id), http.StatusOK)
}

func PinList(w http.ResponseWriter, r *http.Request, user users.User) {
	listId, found, valid := utils.ExtractUintFromRequest("id", r)
	if !found {
		utils.Prettier(w, "no id provided", nil, http.StatusBadRequest)
		return
	}
	if !valid {
		utils.Prettier(w, "invalid list id", nil, http.StatusBadRequest)
		return
	}

	list, found := LoadListById(listId)
	if !found {
		utils.Prettier(w, "list not found", nil, http.StatusBadRequest)
		return
	}

	if !list.CanView(user.Id) {
		utils.Prettier(w, "you do not have the permission !", nil, http.StatusUnauthorized)
		return
	}

	if utils.ContainsUint64(user.PinnedLists, list.Id) {
		err := user.SetPinnedLists(utils.RemoveFromUint64Slice(user.PinnedLists, list.Id))
		if err != nil {
			utils.Prettier(w, "failed to unpin the list", nil, http.StatusInternalServerError)
			utils.PrintError(err)
			return
		}
	} else {
		err := user.SetPinnedLists(append(user.PinnedLists, list.Id))
		if err != nil {
			utils.Prettier(w, "failed to pin the list", nil, http.StatusInternalServerError)
			utils.PrintError(err)
			return
		}
	}

	utils.Prettier(w, "pinned lists updated successfully", user.ExportPrivate(), http.StatusOK)
}

func DeleteList(w http.ResponseWriter, r *http.Request, user users.User) {
	listId, found, valid := utils.ExtractUintFromRequest("id", r)
	if !found {
		utils.Prettier(w, "no id provided", nil, http.StatusBadRequest)
		return
	}
	if !valid {
		utils.Prettier(w, "invalid list id", nil, http.StatusBadRequest)
		return
	}

	list, found := LoadListById(listId)
	if !found {
		utils.Prettier(w, "list not found", nil, http.StatusBadRequest)
		return
	}

	if !list.CanEdit(user.Id) {
		utils.Prettier(w, "you do not have the permission !", nil, http.StatusUnauthorized)
		return
	}

	err := deleteList(listId)
	if err != nil {
		utils.Prettier(w, "failed to delete the list", nil, http.StatusInternalServerError)
		utils.PrintError(err)
		return
	}

	utils.Prettier(w, "list deleted successfully", nil, http.StatusOK)
}
