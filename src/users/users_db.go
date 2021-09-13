package users

import (
	"database/sql"
	"errors"
	"fmt"
	"listes_back/src/database"
	"listes_back/src/utils"
)

func createUser(username, password, email string) (User, error) {
	conn, err := database.GetDb().GetConnection()
	if err != nil {
		return User{}, err
	}
	defer database.CloseConnection(conn)

	// _, err = conn.Exec(fmt.Sprintf("INSERT INTO users(username, password, email, settings) VALUES('%s', '%s', '%s', '%s')", username, password, email, jsonDefaultSettings))
	var newUserId uint64
	err = conn.QueryRow("INSERT INTO users(username, password, email, settings) VALUES($1, $2, $3, $4) returning id;", username, password, email, jsonDefaultSettings).Scan(&newUserId)
	if err != nil {
		return User{}, err
	}

	user, found := LoadUserById(newUserId)
	if !found {
		return User{}, errors.New("failed to create user")
	}
	return user, nil
}

func LoadUserById(id uint64) (User, bool) {
	conn, err := database.GetDb().GetConnection()
	if err != nil {
		utils.PrintError(err)
		return User{}, false
	}
	defer database.CloseConnection(conn)

	row := conn.QueryRow("SELECT id, username, password, email, last_password_refresh, pinned_lists, settings FROM users WHERE id = $1 LIMIT 1;", id)
	if err = row.Err(); err != nil {
		utils.PrintError(err)
		return User{}, false
	}

	var user User
	var pinnedLists database.Uint64Slice
	var rawSettings string
	err = row.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.LastPasswordRefresh, &pinnedLists, &rawSettings)
	if err != nil {
		if err != sql.ErrNoRows {
			utils.PrintError(err)
		}
		return User{}, false
	}

	return user.withPinnedLists(pinnedLists).withSettings(parseUserSettings(rawSettings)), true
}

func loadUsersBy(field string, value interface{}) ([]User, error) {
	conn, err := database.GetDb().GetConnection()
	if err != nil {
		return []User{}, err
	}
	defer database.CloseConnection(conn)

	rows, err := conn.Query(fmt.Sprintf("SELECT id, username, password, email, last_password_refresh, pinned_lists, settings FROM users WHERE %s = $1;", field), value)
	if err != nil {
		return []User{}, err
	}

	foundUsers := []User{}
	for rows.Next() {
		var user User
		var pinnedLists database.Uint64Slice
		var rawSettings string
		err = rows.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.LastPasswordRefresh, &pinnedLists, &rawSettings)
		if err != nil {
			return []User{}, err
		}
		foundUsers = append(foundUsers, user.withPinnedLists(pinnedLists).withSettings(parseUserSettings(rawSettings)))
	}
	return foundUsers, nil
}

func checkUserExistence(username, email string) (string, bool, error) {
	conn, err := database.GetDb().GetConnection()
	if err != nil {
		return "", false, err
	}
	defer database.CloseConnection(conn)

	var exist bool
	err = conn.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);", username).Scan(&exist)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("ici")
			return "", false, err
		}
	}
	if exist { // username is already registered
		return "username", true, nil
	}

	err = conn.QueryRow(fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM users WHERE email = '%s')", email)).Scan(&exist)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("là")
			return "", false, err
		}
	}
	if exist { // email is already registered
		return "email", true, nil
	}

	return "", false, nil
}

func updateUserById(id uint64, field string, value interface{}) error {
	conn, err := database.GetDb().GetConnection()
	if err != nil {
		return err
	}
	defer database.CloseConnection(conn)

	_, err = conn.Exec(fmt.Sprintf("UPDATE users SET %s = $1 WHERE id = $2;", field), value, id)
	return err
}
