package users

import (
	"listes_back/src/utils"
	"net/http"
)

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user, err := ExtractCurrentUserFromRequest(r)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	utils.Prettier(w, "private user informations", user.ExportPrivate(), http.StatusOK)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := ExtractUserFromRequest(r)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	utils.Prettier(w, "public user informations", user.ExportPublic(), http.StatusOK)
}

func UpdateCurrentUser(w http.ResponseWriter, r *http.Request, user User) {
	_ = r.ParseForm()
	username := r.FormValue("username")
	updateUsername := username != ""
	if updateUsername { // provided
		if username == user.Username {
			utils.Prettier(w, "your new username cannot be the same as the old !", user.ExportPrivate(), http.StatusBadRequest)
			return
		}
		if invalidUserNameError := utils.IsUserNameValid(username); invalidUserNameError != nil {
			utils.Prettier(w, invalidUserNameError.Error(), nil, http.StatusBadRequest)
			return
		}
	}

	email := r.FormValue("email")
	updateEmail := email != ""
	if updateEmail { // provided
		if email == user.Email {
			utils.Prettier(w, "your new email cannot be the same as the old !", user.ExportPrivate(), http.StatusBadRequest)
			return
		}
		if invalidEmailError := utils.IsEmailValid(email); invalidEmailError != nil {
			utils.Prettier(w, invalidEmailError.Error(), nil, http.StatusBadRequest)
			return
		}
	}

	if !updateUsername && !updateEmail {
		utils.Prettier(w, "no information provided", nil, http.StatusBadRequest)
		return
	}

	if updateUsername {
		err := updateUserById(user.Id, "username", username)
		if err != nil {
			utils.Prettier(w, err.Error(), nil, http.StatusInternalServerError)
			return
		}
		user.Username = username
	}

	if updateEmail {
		err := updateUserById(user.Id, "email", email)
		if err != nil {
			utils.Prettier(w, err.Error(), nil, http.StatusInternalServerError)
			return
		}
		user.Email = email
	}

	utils.Prettier(w, "user updated successfully", user.ExportPrivate(), http.StatusOK)
}
