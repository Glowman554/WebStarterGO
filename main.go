package main

import (
	"database/sql"
	"net/http"

	"github.com/Glowman554/WebStarterGO/database"
	"github.com/Glowman554/WebStarterGO/render"
	"github.com/Glowman554/WebStarterGO/templates/components/account"
	"github.com/Glowman554/WebStarterGO/templates/pages"
	"github.com/Glowman554/WebStarterGO/users"
	"golang.org/x/exp/slog"
)

func main() {
	slog.Info("Starting...")

	err := database.WithDatabase(func(db *sql.DB) error {
		err := database.ApplyMigration(db, "migrations/0000_users.sql")

		if err != nil {
			return err
		}
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

		http.HandleFunc("/", render.AsHandler(render.ApplyLayout(pages.Index, "Home")))
		http.HandleFunc("/account/login", render.AsRequestHandler(users.HandleUserLogin))
		http.HandleFunc("/account/create", render.AsRequestHandler(users.HandleUserCreate))
		http.HandleFunc("/account/logout", render.AsRequestHandler(users.ApplyUser(users.HandleUserLogout, account.NotLoggedIn)))

		return http.ListenAndServe(":8080", nil)
	})
	if err != nil {
		slog.Error(err.Error())
	}
}
