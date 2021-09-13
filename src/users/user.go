package users

import (
	"listes_back/src/database"
	"time"
)

type User struct {
	Id                  uint64       `json:"id"`
	Username            string       `json:"username"`
	Password            string       `json:"password"`
	Email               string       `json:"email"`
	LastPasswordRefresh time.Time    `json:"last_password_refresh"`
	PinnedLists         []uint64     `json:"pinned_lists"`
	Settings            UserSettings `json:"settings"`
}

func (user User) withPinnedLists(pinnedLists []uint64) User {
	user.PinnedLists = pinnedLists
	return user
}

func (user User) withSettings(settings UserSettings) User {
	user.Settings = settings
	return user
}

func (user User) ExportPublic() interface{} {
	return struct {
		Id       uint64 `json:"id"`
		Username string `json:"username"`
	}{
		Id:       user.Id,
		Username: user.Username,
	}
}

func (user User) ExportPrivate() interface{} {
	return struct {
		Id          uint64       `json:"id"`
		Username    string       `json:"username"`
		Email       string       `json:"email"`
		PinnedLists []uint64     `json:"pinned_lists"`
		Settings    UserSettings `json:"settings"`
	}{
		Id:          user.Id,
		Username:    user.Username,
		Email:       user.Email,
		PinnedLists: user.PinnedLists,
		Settings:    user.Settings,
	}
}

func (user *User) SetPinnedLists(pinned []uint64) error {
	err := updateUserById(user.Id, "pinned_lists", database.Uint64Slice(pinned))
	if err != nil {
		return err
	}
	user.PinnedLists = pinned
	return nil
}
