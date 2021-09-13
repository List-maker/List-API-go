package utils

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

var userNameRegex = regexp.MustCompile(`^([[:alnum:]]|[_,-]){3,16}$`)
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func IsUserNameValid(username string) error {
	if len(username) == 0 {
		return errors.New("aucun nom d'utilisateur fourni")
	}

	if len(username) < 3 || len(username) > 16 {
		return errors.New("la longueur du nom d'utilisateur doit être comprise entre 3 et 16 caractères")
	}

	if !userNameRegex.MatchString(username) {
		return errors.New("le nom d'utilisateur ne respecte pas le format demandé")
	}

	return nil
}

func IsEmailValid(email string) error {
	if len(email) == 0 {
		return errors.New("aucun email fourni")
	}

	if len(email) < 5 && len(email) > 254 {
		return errors.New("la longueur de l'email doit être comprise entre 5 et 254 caractères")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("l'email ne respecte pas le format demandé")
	}

	return nil
}

func IsPasswordValid(password string) error {

	if !regexp.MustCompile(`(?i)^.{8,64}$`).MatchString(password) {
		return errors.New("le mot de passe ne respecte pas le format demandé")
	}

	if !regexp.MustCompile(`^.*[A-Z]+.*$`).MatchString(password) {
		return errors.New("le mot de passe doit contenir au moins une lettre majuscule")
	}

	if !regexp.MustCompile(`^.*[a-z]+.*$`).MatchString(password) {
		return errors.New("le mot de passe doit contenir au moins une lettre miniscule")
	}

	if !regexp.MustCompile(`(?i)^.*[0-9]+.*$`).MatchString(password) {
		return errors.New("le mot de passe doit contenir au moins un chiffre")
	}

	if !regexp.MustCompile(`^.*([ -/]|[:-@]|[\[-\x60]|[{-~])+`).MatchString(password) {
		return errors.New("le mot de passe doit contenir au moins un caractère spécial")
	}

	return nil
}

func IsNameValid(name string) error {
	if len(name) < 1 || len(name) > 255 {
		return errors.New("the length of the name must be included between 0 and 255 characters")
	}
	return nil
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
