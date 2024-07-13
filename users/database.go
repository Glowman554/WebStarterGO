package users

import (
	"database/sql"
)

type User struct {
	Username     string
	PasswordHash string
}

func insertUser(db *sql.DB, username string, passwordHash string) error {
	_, err := db.Exec(`insert into users (username, password_hash) values ($1, $2)`, username, passwordHash)
	return err
}

func insertUserSession(db *sql.DB, username string, token string) error {
	_, err := db.Exec(`insert into sessions (username, token) values ($1, $2)`, username, token)
	return err
}

func LoadUser(db *sql.DB, username string) (*User, error) {
	user := &User{
		Username: username,
	}

	err := db.QueryRow(`select password_hash from users where username = $1`, username).Scan(&user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func LoadUserByToken(db *sql.DB, token string) (*User, error) {
	user := &User{}

	err := db.QueryRow(`select sessions.username, password_hash from users, sessions where users.username = sessions.username and token = $1`, token).Scan(&user.Username, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return user, nil
}
