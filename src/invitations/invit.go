package invitations

type Invitation struct {
	Id             uint64 `json:"id"`
	InvitingUserId uint64 `json:"inviting_user_id"`
	InvitedUserId  uint64 `json:"invited_user_id"`
	ListId         uint64 `json:"list_id"`
	EditingRights  bool   `json:"editing_rights"`
}

func (invit Invitation) Export(listName string) interface{} {
	return struct {
		Invitation
		ListName string `json:"list_name"`
	}{
		Invitation: invit,
		ListName:   listName,
	}
}
