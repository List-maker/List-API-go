package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"listes_back/src/invitations"
	"listes_back/src/lists"
	"listes_back/src/lists/items"
	"listes_back/src/users"
	"net/http"
)

func initRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello there !")
	})

	// Auth
	r.HandleFunc("/api/auth/login", users.Login)
	r.HandleFunc("/api/auth/register", users.Register)
	r.HandleFunc("/api/auth/logout", users.AuthRequired(users.Logout))

	// Users
	r.HandleFunc("/api/user", users.AuthRequired(users.GetCurrentUser)).Methods(http.MethodGet)
	r.HandleFunc("/api/user/{id}", users.GetUser).Methods(http.MethodGet)
	r.HandleFunc("/api/user", users.UserRequired(users.UpdateCurrentUser)).Methods(http.MethodPut)
	// r.HandleFunc("/api/user/password", users.AuthRequired(users.UpdatePassword)).Methods(http.MethodPut)
	r.HandleFunc("/api/user/{id}/avatar", users.GetAvatar).Methods(http.MethodGet)
	r.HandleFunc("/api/user/avatar", users.UserRequired(users.UpdateAvatar)).Methods(http.MethodPut)
	r.HandleFunc("/api/user/avatar", users.UserRequired(users.DeleteAvatar)).Methods(http.MethodDelete)

	// Lists
	r.HandleFunc("/api/list", users.UserRequired(lists.CreateList)).Methods(http.MethodPost)
	r.HandleFunc("/api/list/{id}", users.UserRequired(lists.GetList)).Methods(http.MethodGet)
	r.HandleFunc("/api/list/{id}", users.UserRequired(lists.UpdateList)).Methods(http.MethodPut)
	r.HandleFunc("/api/list/{id}/pin", users.UserRequired(lists.PinList)).Methods(http.MethodPut)
	r.HandleFunc("/api/list/{id}", users.UserRequired(lists.DeleteList)).Methods(http.MethodDelete)
	r.HandleFunc("api/list/user/{id}", users.UserRequired(lists.GetUserLists)).Methods(http.MethodGet)

	// Items
	r.HandleFunc("/api/list/{list_id}/add", users.UserRequired(items.CreateItem)).Methods(http.MethodPost)
	r.HandleFunc("/api/item/{id}", users.UserRequired(items.GetItem)).Methods(http.MethodGet)
	r.HandleFunc("/api/item/{id}", users.UserRequired(items.UpdateItem)).Methods(http.MethodPut)
	r.HandleFunc("/api/item/{id}/check", users.UserRequired(items.CheckItem)).Methods(http.MethodPut)
	r.HandleFunc("/api/item/{id}", users.UserRequired(items.DeleteItem)).Methods(http.MethodDelete)

	// Invitations
	r.HandleFunc("/api/invitation/new", users.UserRequired(invitations.CreateInvit)).Methods(http.MethodPost)
	r.HandleFunc("/api/invitation/list", users.UserRequired(invitations.ListInvits)).Methods(http.MethodGet)
	r.HandleFunc("/api/invitation/{id}", users.UserRequired(invitations.GetInvit)).Methods(http.MethodGet)
	r.HandleFunc("/api/invitation/{id}/accept", users.UserRequired(invitations.AcceptInvit)).Methods(http.MethodPut)
	r.HandleFunc("/api/invitation/{id}/delete", users.UserRequired(invitations.DeleteInvit)).Methods(http.MethodPut)
	return r
}
