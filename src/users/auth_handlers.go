package users

import (
	"fmt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"golang.org/x/crypto/bcrypt"
	"listes_back/src/utils"
	"net/http"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var foundUsers []User
	var err error

	_ = r.ParseForm()

	// check login (username or email) validity
	login := r.FormValue("login")
	if invalidEmailError := utils.IsEmailValid(login); invalidEmailError == nil {
		foundUsers, err = loadUsersBy("email", login)
	} else if invalidUsernameError := utils.IsUserNameValid(login); invalidUsernameError == nil {
		foundUsers, err = loadUsersBy("username", login)
	} else {
		utils.Prettier(w, "bad login or password !", nil, http.StatusUnauthorized)
		return
	}

	if err != nil {
		utils.PrintError(err)
		utils.Prettier(w, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	// check if user exist
	if len(foundUsers) == 0 {
		utils.Prettier(w, "user doesnt exist!", nil, http.StatusUnauthorized)
		return
	}

	// get password
	password := r.FormValue("password")

	// get user
	user := foundUsers[0]

	// check hash password
	compareErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if compareErr != nil {
		utils.Prettier(w, "invalid password !", nil, http.StatusUnauthorized)
		return
	}

	// Create the RefreshToken
	refreshToken := jwt.New()
	_ = refreshToken.Set("id", user.Id)
	_ = refreshToken.Set(jwt.ExpirationKey, time.Now().Add(expirationRefreshToken))

	// Sign the token and generate a payload
	signed, err := jwt.Sign(refreshToken, jwa.RS256, privKey)
	if err != nil {
		utils.PrintError(err)
		utils.Prettier(w, fmt.Sprintf("échec de la génération de la payload chiffrée: %s\n", err), nil, http.StatusInternalServerError)
		return
	}
	refreshTokenStr := string(signed)

	// Create the accessToken
	accessToken := jwt.New()
	_ = accessToken.Set("id", user.Id)
	_ = accessToken.Set("username", user.Username)
	_ = accessToken.Set("email", user.Email)
	_ = accessToken.Set(jwt.ExpirationKey, time.Now().Add(expirationAccessToken))

	// Sign the token and generate a payload
	signed, err = jwt.Sign(accessToken, jwa.RS256, privKey)
	if err != nil {
		utils.PrintError(err)
		utils.Prettier(w, fmt.Sprintf("échec de la génération de la payload chiffrée: %s\n", err), nil, http.StatusInternalServerError)
		return
	}
	accessTokenStr := string(signed)

	utils.Prettier(w, "Token generated !", struct {
		RefreshToken string `json:"refresh_token"`
		AccessToken  string `json:"access_token"`
	}{refreshTokenStr, accessTokenStr}, http.StatusOK)
}

func Register(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()

	// Check username validity
	username := r.FormValue("username")
	fmt.Println("Username: ", username)
	if invalidUserNameError := utils.IsUserNameValid(username); invalidUserNameError != nil {
		utils.Prettier(w, invalidUserNameError.Error(), nil, http.StatusBadRequest)
		return
	}

	// Check email validity
	email := r.FormValue("email")
	if invalidEmailError := utils.IsEmailValid(email); invalidEmailError != nil {
		utils.Prettier(w, invalidEmailError.Error(), nil, http.StatusBadRequest)
		return
	}

	// Check password validity
	password := r.FormValue("password")
	if invalidPasswordError := utils.IsPasswordValid(password); invalidPasswordError != nil {
		utils.Prettier(w, invalidPasswordError.Error(), nil, http.StatusBadRequest)
		return
	}

	takenField, alreadyExist, err := checkUserExistence(username, email)
	if err != nil {
		utils.Prettier(w, takenField+err.Error(), nil, http.StatusInternalServerError)
		return
	}
	if alreadyExist {
		utils.Prettier(w, "This "+takenField+" is already taken", nil, http.StatusBadRequest) // BadRequest ?
		return
	}

	// Check hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	user, err := createUser(username, hashedPassword, email)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	// Create the RefreshToken
	refreshToken := jwt.New()
	_ = refreshToken.Set("id", user.Id)
	_ = refreshToken.Set(jwt.ExpirationKey, time.Now().Add(expirationRefreshToken))

	// Sign the token and generate a payload
	signed, err := jwt.Sign(refreshToken, jwa.RS256, privKey)
	if err != nil {
		utils.PrintError(err)
		utils.Prettier(w, fmt.Sprintf("échec de la génération de la payload chiffrée: %s\n", err), nil, http.StatusInternalServerError)
		return
	}
	refreshTokenStr := string(signed)

	// Create the accessToken
	accessToken := jwt.New()
	_ = accessToken.Set("id", user.Id)
	_ = accessToken.Set("username", user.Username)
	_ = accessToken.Set("email", user.Email)
	_ = accessToken.Set(jwt.ExpirationKey, time.Now().Add(expirationAccessToken))

	// Sign the token and generate a payload
	signed, err = jwt.Sign(accessToken, jwa.RS256, privKey)
	if err != nil {
		utils.PrintError(err)
		utils.Prettier(w, fmt.Sprintf("échec de la génération de la payload chiffrée: %s\n", err), nil, http.StatusInternalServerError)
		return
	}
	accessTokenStr := string(signed)

	utils.Prettier(w, "user created !", struct {
		RefreshToken string `json:"refresh_token"`
		AccessToken  string `json:"access_token"`
	}{refreshTokenStr, accessTokenStr}, http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		utils.Prettier(w, "invalid token", nil, http.StatusBadRequest)
		return
	}

	accessToken, err := ExtractToken(r)
	if err != nil {
		utils.Prettier(w, "invalid token (missing bearer)", nil, http.StatusBadRequest)
	}

	refreshToken := r.FormValue("refresh_token")

	blackListAccessToken = append(blackListAccessToken, accessToken)
	blackListRefreshToken = append(blackListRefreshToken, refreshToken)

	utils.Prettier(w, "logout success", nil, http.StatusOK)
}
