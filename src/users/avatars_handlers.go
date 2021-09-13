package users

import (
	"fmt"
	"listes_back/src/utils"
	"net/http"
)

func GetAvatar(w http.ResponseWriter, r *http.Request) {
	user, err := ExtractUserFromRequest(r)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	printAvatar(w, user.Id)
}

func UpdateAvatar(w http.ResponseWriter, r *http.Request, user User) {
	err, status := updateAvatar(user.Id, r)
	if err != nil {
		utils.Prettier(w, fmt.Sprintf("failed to update avatar: %v", err), nil, status)
		return
	}
	utils.Prettier(w, "avatar updated successfully !", nil, http.StatusOK)
}

func DeleteAvatar(w http.ResponseWriter, r *http.Request, user User) {
	err := deleteAvatar(user.Id)
	if err != nil {
		utils.Prettier(w, "failed to delete avatar: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}
	utils.Prettier(w, "avatar deleted successfully !", nil, http.StatusOK)
}
