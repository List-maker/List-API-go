package users

import (
	"errors"
	"github.com/lestrrat-go/jwx/jwt"
	"listes_back/src/utils"
	"net/http"
	"strings"
	"time"
)

var blackListAccessToken []string
var blackListRefreshToken []string

func AuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := ExtractToken(r)
		if err != nil {
			utils.Prettier(w, "invalid token (missing bearer)", nil, http.StatusUnauthorized)
			return
		}
		for _, val := range blackListAccessToken {
			if token == val {
				utils.Prettier(w, "invalid token !", nil, http.StatusUnauthorized)
				return
			}
		}

		_, err = ExtractIdFromRequest(r)
		if err != nil {
			utils.Prettier(w, "invalid token !", nil, http.StatusUnauthorized)
			return
		}
		handler.ServeHTTP(w, r)
	}
}

func UserRequired(handler func(http.ResponseWriter, *http.Request, User)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := ExtractCurrentUserFromRequest(r)
		if err != nil {
			utils.Prettier(w, "invalid token !", nil, http.StatusUnauthorized)
			return
		}
		handler(w, r, user)
	}
}

func ExtractIdFromRequest(r *http.Request) (uint64, error) {
	headerToken := r.Header.Get("Authorization")
	if !strings.HasPrefix(headerToken, "Bearer ") {
		return 0, errors.New("invalid token (missing bearer)")
	}
	headerToken = strings.TrimPrefix(headerToken, "Bearer ")
	token, err := jwt.ParseString(headerToken, jwt.WithKeySet(keySet), jwt.UseDefaultKey(true))
	if err != nil {
		return 0, err
	}
	id, ok := token.Get("id")
	if !ok {
		return 0, errors.New("invalid token")
	}

	if !time.Now().UTC().Before(token.Expiration()) {
		return 0, errors.New("expired token")
	}

	floatId := id.(float64)
	if floatId < 1 {
		return 0, errors.New("invalid token")
	}

	return uint64(floatId), nil
}

func ExtractToken(r *http.Request) (string, error) {
	headerToken := r.Header.Get("Authorization")
	if !strings.HasPrefix(headerToken, "Bearer ") {
		return "", errors.New("invalid token (missing bearer)")
	}
	headerToken = strings.TrimPrefix(headerToken, "Bearer ")

	return headerToken, nil
}
