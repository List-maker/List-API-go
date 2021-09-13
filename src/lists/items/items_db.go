package items

import (
	"database/sql"
	"fmt"
	"listes_back/src/database"
	"listes_back/src/utils"
)

func createItem(listId uint64, name string) (uint64, error) {
	conn, err := database.GetDb().GetConnection()
	if err != nil {
		return 0, err
	}
	defer database.CloseConnection(conn)

	var itemId uint64
	err = conn.QueryRow("INSERT INTO list_items(parent_id, name) VALUES($1, $2) returning id;", listId, name).Scan(&itemId)
	return itemId, err
}

func loadItemFromDb(itemId uint64) (Item, bool) {
	conn, err := database.GetDb().GetConnection()
	if err != nil {
		return Item{}, false
	}
	defer database.CloseConnection(conn)

	row := conn.QueryRow("SELECT id, parent_id, name, checked FROM list_items WHERE id = $1 LIMIT 1;", itemId)
	if err = row.Err(); err != nil {
		utils.PrintError(err)
		return Item{}, false
	}

	var item Item
	err = row.Scan(&item.Id, &item.ParentId, &item.Name, &item.Checked)
	if err != nil {
		if err != sql.ErrNoRows {
			utils.PrintError(err)
		}
		return Item{}, false
	}

	return item, true
}

func updateItem(itemId uint64, field string, value interface{}) error {
	conn, err := database.GetDb().GetConnection()
	if err != nil {
		return err
	}
	defer database.CloseConnection(conn)

	_, err = conn.Exec(fmt.Sprintf("UPDATE list_items SET %s = $1 WHERE id = $2;", field), value, itemId)
	return err
}

func deleteItem(itemId uint64) error {
	conn, err := database.GetDb().GetConnection()
	if err != nil {
		return err
	}
	defer database.CloseConnection(conn)

	_, err = conn.Exec("DELETE FROM list_items WHERE id = $1;", itemId)
	return err
}
