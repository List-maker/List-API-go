package invitations

import (
	"listes_back/src/lists"
	"listes_back/src/users"
	"listes_back/src/utils"
	"net/http"
	"strconv"
)

func CreateInvit(w http.ResponseWriter, r *http.Request, user users.User) {
	_ = r.ParseForm()

	invitedUserIdStr, found := r.Form["invited_user_id"]
	if !found {
		utils.Prettier(w, "no invited user id provided", nil, http.StatusBadRequest)
		return
	}
	invitedUserId, err := strconv.ParseUint(invitedUserIdStr[0], 10, 64)
	if err != nil {
		utils.Prettier(w, "invalid invited user id", nil, http.StatusBadRequest)
		return
	}
	if user.Id == invitedUserId {
		utils.Prettier(w, "you can't invite yourself", nil, http.StatusBadRequest)
		return
	}
	invitedUser, found := users.LoadUserById(invitedUserId)
	if !found {
		utils.Prettier(w, "invited user not found", nil, http.StatusBadRequest)
		return
	}

	listIdStr, found := r.Form["list_id"]
	if !found {
		utils.Prettier(w, "no list id provided", nil, http.StatusBadRequest)
		return
	}
	listId, err := strconv.ParseUint(listIdStr[0], 10, 64)
	if err != nil {
		utils.Prettier(w, "invalid list id", nil, http.StatusBadRequest)
		return
	}
	list, found := lists.LoadListById(listId)
	if !found {
		utils.Prettier(w, "list not found", nil, http.StatusBadRequest)
		return
	}

	editingRightsStr, found := r.Form["editing_rights"]
	if !found {
		/*utils.Prettier(w, "no editing rights provided", nil, http.StatusBadRequest)
		return*/
		editingRightsStr = []string{"false"} // no rights provided = no editing rights for the invitation
	}
	editingRights, err := strconv.ParseBool(editingRightsStr[0])
	if err != nil {
		utils.Prettier(w, "invalid editing rights", nil, http.StatusBadRequest)
		return
	}

	if editingRights && !list.CanEdit(user.Id) {
		utils.Prettier(w, "you do not have the permission to share this list with editing rights", nil, http.StatusUnauthorized)
		return
	}

	if list.CanView(invitedUserId) {
		utils.Prettier(w, "this user has already access to this list", nil, http.StatusBadRequest)
		return
	}

	if isUserInvitedToList(invitedUserId, listId) {
		utils.Prettier(w, "this user is already invited to this list", nil, http.StatusBadRequest)
		return
	}

	invitId, err := createInvit(user.Id, invitedUser.Id, list.Id, editingRights)
	if err != nil {
		utils.PrintError(err)
		utils.Prettier(w, "failed to create invitation", nil, http.StatusInternalServerError)
		return
	}

	utils.Prettier(w, "invitation created successfully", invitId, http.StatusOK)
}

func ListInvits(w http.ResponseWriter, r *http.Request, user users.User) {
	fromUser, toUser, err := getUserInvitations(user.Id)
	if err != nil {
		utils.PrintError(err)
		utils.Prettier(w, "failed to retrieve invitations", nil, http.StatusInternalServerError)
		return
	}

	invitations := struct {
		FromYou []uint64 `json:"from_you"`
		ToYou   []uint64 `json:"to_you"`
	}{
		FromYou: fromUser,
		ToYou:   toUser,
	}

	utils.Prettier(w, "user invitations", invitations, http.StatusOK)
}

func GetInvit(w http.ResponseWriter, r *http.Request, user users.User) {
	invitId, found, valid := utils.ExtractUintFromRequest("id", r)
	if !found {
		utils.Prettier(w, "no invitation id provided", nil, http.StatusBadRequest)
		return
	}
	if !valid {
		utils.Prettier(w, "invalid invitation id", nil, http.StatusBadRequest)
		return
	}
	invit, listName, found := loadInvitById(invitId, true)
	if !found {
		utils.Prettier(w, "invitation not found", nil, http.StatusBadRequest)
		return
	}

	if user.Id != invit.InvitingUserId && user.Id != invit.InvitedUserId {
		utils.Prettier(w, "you do not have the permission !", nil, http.StatusUnauthorized)
		return
	}

	utils.Prettier(w, "invitation informations", invit.Export(listName), http.StatusOK)
}

func AcceptInvit(w http.ResponseWriter, r *http.Request, user users.User) {
	invitId, found, valid := utils.ExtractUintFromRequest("id", r)
	if !found {
		utils.Prettier(w, "no invitation id provided", nil, http.StatusBadRequest)
		return
	}
	if !valid {
		utils.Prettier(w, "invalid invitation id", nil, http.StatusBadRequest)
		return
	}
	invit, _, found := loadInvitById(invitId, false)
	if !found {
		utils.Prettier(w, "invitation not found", nil, http.StatusBadRequest)
		return
	}

	if invit.InvitedUserId != user.Id {
		utils.Prettier(w, "you do not have the permission !", nil, http.StatusUnauthorized)
		return
	}

	list, found := lists.LoadListById(invit.ListId)
	if !found {
		utils.Prettier(w, "the list no longer exists", nil, http.StatusBadRequest)
		if err := deleteInvit(invit.Id); err != nil {
			utils.PrintError(err)
		}
		return
	}

	if list.CanView(user.Id) {
		utils.Prettier(w, "you are already on this list", nil, http.StatusBadRequest)
		if err := deleteInvit(invit.Id); err != nil {
			utils.PrintError(err)
		}
		return
	}

	var err error
	if invit.EditingRights {
		err = list.SetEditors(append(list.Editors, user.Id))
	} else {
		err = list.SetViewers(append(list.Viewers, user.Id))
	}
	if err == nil {
		err = deleteInvit(invit.Id)
	}
	if err != nil {
		utils.Prettier(w, "failed to accept invitation", nil, http.StatusInternalServerError)
		utils.PrintError(err)
		return
	}

	utils.Prettier(w, "invitation accepted successfully !", list.ExportFor(user.Id), http.StatusOK)
}

func DeleteInvit(w http.ResponseWriter, r *http.Request, user users.User) {
	invitId, found, valid := utils.ExtractUintFromRequest("id", r)
	if !found {
		utils.Prettier(w, "no invitation id provided", nil, http.StatusBadRequest)
		return
	}
	if !valid {
		utils.Prettier(w, "invalid invitation id", nil, http.StatusBadRequest)
		return
	}
	invit, _, found := loadInvitById(invitId, false)
	if !found {
		utils.Prettier(w, "invitation not found", nil, http.StatusBadRequest)
		return
	}

	if invit.InvitingUserId != user.Id && invit.InvitedUserId != user.Id {
		utils.Prettier(w, "you do not have the permission !", nil, http.StatusUnauthorized)
		return
	}

	err := deleteInvit(invit.Id)
	if err != nil {
		utils.Prettier(w, "failed to delete invitation", nil, http.StatusInternalServerError)
		utils.PrintError(err)
		return
	}

	utils.Prettier(w, "invitation deleted successfully !", nil, http.StatusOK)
}
