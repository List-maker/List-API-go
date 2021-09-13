package users

type Invitation struct {
	Id             uint64
	InvitingUserId uint64
	InvitedUserId  uint64
	ListId         uint64
}
