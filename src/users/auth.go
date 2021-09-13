package users

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"listes_back/src/utils"
	"net/http"
	"time"
)

var privKey *rsa.PrivateKey
var keySet jwk.Set

const expirationRefreshToken = time.Hour * 168
const expirationAccessToken = time.Hour * 3
const expirationResetToken = time.Hour * 24
const expirationActivateToken = time.Hour * 168

func init() {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("private key generation failed: %s\n", err)
		return
	}
	privKey = key

	pubKey, err := jwk.New(privKey.PublicKey)
	if err != nil {
		fmt.Printf("JWK creation failed: %s\n", err)
		return
	}
	_ = pubKey.Set(jwk.AlgorithmKey, jwa.RS256)

	// This JWKS can *only* have 1 key.
	set := jwk.NewSet()
	set.Add(pubKey)
	keySet = set
}

func ExtractCurrentUserFromRequest(r *http.Request) (User, error) {
	userId, err := ExtractIdFromRequest(r)
	if err != nil {
		return User{}, err
	}

	user, found := LoadUserById(userId)
	if !found {
		return User{}, err
	}

	return user, nil
}

func ExtractUserFromRequest(r *http.Request) (User, error) {
	userId, found, valid := utils.ExtractUintFromRequest("id", r)
	if !found {
		return User{}, errors.New("no id provided")
	}
	if !valid {
		return User{}, errors.New("invalid id")
	}

	user, found := LoadUserById(userId)
	if !found {
		return User{}, errors.New("invalid id")
	}
	return user, nil
}
