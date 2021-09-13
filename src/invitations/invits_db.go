package invitations

import (
	"database/sql"
	"listes_back/src/database"
	"listes_back/src/utils"
)

func createInvit(invitingUserId, invitedUserId, listId uint64, editingRights bool) (uint64, error) {
	conn, err := database.GetDb().GetConnection()
	if err != nil {
		return 0, err
	}
	defer database.CloseConnection(conn)

	var newInvitId uint64
	err = conn.QueryRow("INSERT INTO list_invitations(inviting_user_id, invited_user_id, list_id, editing_rights) VALUES($1, $2, $3, $4) returning id;", invitingUserId, invitedUserId, listId, editingRights).Scan(&newInvitId)
	if err != nil {
		return 0, err
	}
	return newInvitId, nil
}

func isUserInvitedToList(invitedUserId, listId uint64) bool {
	conn, err := database.GetDb().GetConnection()
	if err != nil {
		utils.PrintError(err)
		return false
	}
	defer database.CloseConnection(conn)

	var invited bool
	err = conn.QueryRow("SELECT EXISTS(SELECT * FROM list_invitations WHERE invited_user_id = $1 AND list_id = $2);", invitedUserId, listId).Scan(&invited)
	if err != nil {
		utils.PrintError(err)
		return false
	}
	return invited
}

func getUserInvitations(userId uint64) (fromUser, toUser []uint64, err error) {
	conn, err := database.GetDb().GetConnection()
	if err != nil {
		return nil, nil, err
	}
	defer database.CloseConnection(conn)

	rows, err := conn.Query("SELECT id FROM list_invitations WHERE inviting_user_id = $1 ORDER BY id;", userId)
	if err != nil {
		return nil, nil, err
	}

	fromUser = []uint64{}
	for rows.Next() {
		var invitId uint64
		err = rows.Scan(&invitId)
		if err != nil {
			return nil, nil, err
		}
		fromUser = append(fromUser, invitId)
	}

	rows, err = conn.Query("SELECT id FROM list_invitations WHERE invited_user_id = $1 ORDER BY id;", userId)
	if err != nil {
		return nil, nil, err
	}

	toUser = []uint64{}
	for rows.Next() {
		var invitId uint64
		err = rows.Scan(&invitId)
		if err != nil {
			return nil, nil, err
		}
		toUser = append(toUser, invitId)
	}

	return fromUser, toUser, nil
}

func loadInvitById(invitId uint64, getListName bool) (Invitation, string, bool) {
	conn, err := database.GetDb().GetConnection()
	if err != nil {
		utils.PrintError(err)
		return Invitation{}, "", false
	}
	defer database.CloseConnection(conn)

	row := conn.QueryRow("SELECT id, inviting_user_id, invited_user_id, list_id, editing_rights FROM list_invitations WHERE id = $1;", invitId)
	if row.Err() != nil {
		if err != sql.ErrNoRows {
			utils.PrintError(err)
		}
		return Invitation{}, "", false
	}

	var invit Invitation
	err = row.Scan(&invit.Id, &invit.InvitingUserId, &invit.InvitedUserId, &invit.ListId, &invit.EditingRights)
	if err != nil {
		if err != sql.ErrNoRows {
			utils.PrintError(err)
		}
		return Invitation{}, "", false
	}

	var listName string
	if getListName {
		err = conn.QueryRow("SELECT name FROM user_lists WHERE id = $1;", invit.ListId).Scan(&listName)
		if err != nil {
			utils.PrintError(err)
			listName = "/!\\ FAILED TO RETRIEVE LIST NAME /!\\"
		}
	}

	return invit, listName, true
}

func deleteInvit(invitId uint64) error {
	conn, err := database.GetDb().GetConnection()
	if err != nil {
		return err
	}
	defer database.CloseConnection(conn)

	_, err = conn.Exec("DELETE FROM list_invitations WHERE id = $1;", invitId)
	return err
}
