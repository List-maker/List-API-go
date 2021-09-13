package items

import (
	"errors"
	"listes_back/src/utils"
	"net/http"
)

func LoadItemFromRequest(r *http.Request) (Item, error) {
	itemId, found, valid := utils.ExtractUintFromRequest("id", r)
	if !found {
		return Item{}, errors.New("no id provided")
	}
	if !valid {
		return Item{}, errors.New("invalid item id")
	}

	item, found := loadItemFromDb(itemId)
	if !found {
		return Item{}, errors.New("item not found")
	}
	return item, nil
}
