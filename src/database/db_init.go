package database

import (
	"database/sql"
	"fmt"
)

func (db Database) init() error {
	fmt.Println("Initializing database ...")

	conn, err := db.GetConnection()
	if err != nil {
		return err
	}
	defer CloseConnection(conn)

	// Test the connection to the database
	if err = conn.Ping(); err != nil {
		return err
	}

	// Init the tables of the database
	if err = initUsersTable(conn); err != nil {
		return err
	}
	if err = initListInvitationsTable(conn); err != nil {
		return err
	}
	if err = initUserListsTable(conn); err != nil {
		return err
	}
	if err = initListItems(conn); err != nil {
		return err
	}

	return nil
}

func initUsersTable(conn *sql.DB) error {
	_, err := conn.Exec(`CREATE TABLE IF NOT EXISTS users (
		id bigserial,
		username text NOT NULL,
		password text NOT NULL,
		email text NOT NULL,
		last_password_refresh timestamp without time zone DEFAULT now(),
		pinned_lists json DEFAULT '[]',
		settings json NOT NULL,
		PRIMARY KEY (id)
);`)
	return err
}

func initListInvitationsTable(conn *sql.DB) error {
	_, err := conn.Exec(`CREATE TABLE IF NOT EXISTS list_invitations (
		id bigserial,
		inviting_user_id bigint NOT NULL,
		invited_user_id bigint NOT NULL,
		list_id bigint NOT NULL,
		editing_rights boolean NOT NULL,
		PRIMARY KEY (id)
);`)
	return err
}

func initUserListsTable(conn *sql.DB) error {
	_, err := conn.Exec(`CREATE TABLE IF NOT EXISTS user_lists (
		id bigserial,
		name text NOT NULL,
		items json DEFAULT '[]',
		editors json NOT NULL,
		viewers json DEFAULT '[]',
		PRIMARY KEY (id)
);`)
	return err
}

func initListItems(conn *sql.DB) error {
	_, err := conn.Exec(`CREATE TABLE IF NOT EXISTS list_items (
		id bigserial,
		parent_id bigint NOT NULL,
		name text NOT NULL,
		checked boolean DEFAULT false,
		PRIMARY KEY (id)
);`)
	return err
}
