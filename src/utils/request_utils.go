package utils

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func ExtractUintFromRequest(field string, r *http.Request) (value uint64, found bool, valid bool) {
	vars := mux.Vars(r)

	strValue, found := vars[field]
	if !found {
		return 0, false, false
	}

	parsedValue, err := strconv.ParseUint(strValue, 10, 64)
	if err != nil {
		return 0, true, false
	}

	return parsedValue, true, true
}
