package lists

import (
	"database/sql"
	"fmt"
	"listes_back/src/database"
	"listes_back/src/utils"
)

func createList(userId uint64, listName string) (uint64, error) {
	conn, err := database.GetDb().GetConnection()
	if err != nil {
		return 0, err
	}
	defer database.CloseConnection(conn)

	editors := fmt.Sprintf("[%d]", userId) // quick marshal
	var newListId uint64
	err = conn.QueryRow("INSERT INTO user_lists(name, editors) VALUES($1, $2) returning id;", listName, editors).Scan(&newListId)
	if err != nil {
		return 0, err
	}

	return newListId, nil
}

func QueryUserLists(userId uint64) []uint64 {
	conn, err := database.GetDb().GetConnection()
	if err != nil {
		return []uint64{}
	}
	defer database.CloseConnection(conn)

	rows, err := conn.Query("SELECT id FROM user_lists WHERE editors::jsonb @> $1 OR viewers::jsonb @> $1;", userId)
	if err != nil {
		utils.PrintError(err)
		return []uint64{}
	}

	var allId []uint64
	for rows.Next() {
		var listId uint64
		err = rows.Scan(&listId)
		if err != nil {
			utils.PrintError(err)
			return []uint64{}
		}
		allId = append(allId, listId)
	}
	return allId
}

func LoadListById(listId uint64) (List, bool) {
	conn, err := database.GetDb().GetConnection()
	if err != nil {
		return List{}, false
	}
	defer database.CloseConnection(conn)

	row := conn.QueryRow("SELECT id, name, items, editors, viewers FROM user_lists WHERE id = $1 LIMIT 1;", listId)
	if err = row.Err(); err != nil {
		utils.PrintError(err)
		return List{}, false
	}

	var list List
	var items, editors, viewers database.Uint64Slice
	items = database.Uint64Slice{}
	err = row.Scan(&list.Id, &list.Name, &items, &editors, &viewers)
	if err != nil {
		if err != sql.ErrNoRows {
			utils.PrintError(err)
		}
		return List{}, false
	}
	list.Items = items
	list.Editors = editors
	list.Viewers = viewers

	return list, true
}

func updateList(listId uint64, field string, value interface{}) error {
	conn, err := database.GetDb().GetConnection()
	if err != nil {
		return err
	}
	defer database.CloseConnection(conn)

	_, err = conn.Exec(fmt.Sprintf("UPDATE user_lists SET %s = $1 WHERE id = $2;", field), value, listId)
	return err
}

func deleteList(listId uint64) error {
	conn, err := database.GetDb().GetConnection()
	if err != nil {
		return err
	}
	defer database.CloseConnection(conn)

	_, err = conn.Exec("DELETE FROM user_lists WHERE id = $1;", listId)
	if err != nil {
		return err
	}

	_, err = conn.Exec("DELETE FROM list_items WHERE parent_id = $1;", listId)
	if err != nil {
		return err
	}

	_, err = conn.Exec("DELETE FROM list_invitations WHERE list_id = $1;", listId)
	return err
}
