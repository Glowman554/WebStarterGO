package users

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/Glowman554/WebStarterGO/database"
	"github.com/Glowman554/WebStarterGO/render"
	"github.com/a-h/templ"
	"golang.org/x/crypto/bcrypt"
)

func LoginUser(db *sql.DB, username string, password string) (string, error) {
	if strings.TrimSpace(password) == "" || strings.TrimSpace(username) == "" {
		return "", errors.New("neither password nor username should be empty")
	}

	user, err := LoadUser(db, username)
	if err != nil {
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", err
	}
	return complete(db, username)
}

func CreateUser(db *sql.DB, username string, password string) (string, error) {
	if strings.TrimSpace(password) == "" || strings.TrimSpace(username) == "" {
		return "", errors.New("neither password nor username should be empty")
	}
	err := isValidPassword(password)
	if err != nil {
		return "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	err = insertUser(db, username, string(hash))
	if err != nil {
		return "", err
	}

	return complete(db, username)
}

func complete(db *sql.DB, username string) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	err = insertUserSession(db, username, token)
	if err != nil {
		return "", err
	}

	return token, nil
}

type UserComponent func(user *User, r *http.Request, w http.ResponseWriter) templ.Component

func ApplyUser(component UserComponent, notLoggedIn render.Component) render.RequestComponent {
	return func(r *http.Request, w http.ResponseWriter) templ.Component {
		cookie, err := r.Cookie("Authentication")
		if err != nil {
			return notLoggedIn()
		}

		user, err := LoadUserByToken(database.DB, cookie.Value)
		if err != nil {
			return notLoggedIn()
		}

		return component(user, r, w)
	}
}
