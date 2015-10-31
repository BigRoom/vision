package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// User is someone with an account on Big Room
type User struct {
	ID       int64  `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"-"`
	Email    string `db:"email" json:"email"`
}

// NewUser registers a new user
func NewUser(username, password, email string) (User, error) {
	var u User

	err := DB.
		Select("*").
		From("users").
		Where("username = $1", username).
		QueryStruct(&u)

	if err == nil || u != (User{}) {
		return u, errors.New("That user already exists")
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return u, err
	}

	err = DB.
		InsertInto("users").
		Columns("username", "password", "email").
		Values(username, string(pass), email).
		Returning("*").
		QueryStruct(&u)

	if err != nil {
		return u, err
	}

	if u.Username != username {
		return u, errors.New("User did not sync with database")
	}

	return u, err
}

// FetchUser retrieves a user given a value and key
func FetchUser(key string, value interface{}) (User, error) {
	var u User

	err := DB.
		Select("*").
		From("users").
		Where(key+" = $1", value).
		QueryStruct(&u)

	return u, err
}

// Login verifies that the provided password is correct
func (u User) Login(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
