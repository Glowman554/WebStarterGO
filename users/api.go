package users

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/Glowman554/WebStarterGO/database"
	"github.com/Glowman554/WebStarterGO/render"
	"github.com/Glowman554/WebStarterGO/templates/components"
	"github.com/Glowman554/WebStarterGO/templates/components/account"
	"github.com/a-h/templ"
)

type Payload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func handleCommon(r *http.Request, w http.ResponseWriter, handler func(*sql.DB, string, string) (string, error), field render.Component, success string) templ.Component {
	switch r.Method {
	case "POST":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return components.Error(err)
		}

		var payload Payload
		err = json.Unmarshal(body, &payload)
		if err != nil {
			return components.Error(err)
		}

		token, err := handler(database.DB, payload.Username, payload.Password)
		if err != nil {
			return components.Error(err)
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "Authentication",
			Value: token,
			Path:  "/",
		})
		return account.Success(success)

	case "GET":
		return ApplyUser(func(user *User, r *http.Request, w http.ResponseWriter) templ.Component {
			return account.AlreadyLoggedIn(user.Username)
		}, field)(r, w)
	default:
		return components.Error(errors.New("Invalid method"))
	}
}

func HandleUserCreate(r *http.Request, w http.ResponseWriter) templ.Component {
	return handleCommon(r, w, CreateUser, account.CreateField, "Created user successfully")
}

func HandleUserLogin(r *http.Request, w http.ResponseWriter) templ.Component {
	return handleCommon(r, w, LoginUser, account.LoginField, "Login successful")
}

func HandleUserLogout(user *User, r *http.Request, w http.ResponseWriter) templ.Component {
	http.SetCookie(w, &http.Cookie{
		Name:    "Authentication",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	})
	return account.LogoutSuccess()
}
