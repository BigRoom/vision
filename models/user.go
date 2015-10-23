package models

import (
	"errors"
	"fmt"
)

// User is someone with an account on Big Room
type User struct {
	ID       int64  `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Email    string `db:"email"`
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

	err = DB.
		InsertInto("users").
		Columns("username", "password", "email").
		Values(username, password, email).
		Returning("*").
		QueryStruct(&u)

	if err != nil {
		return u, err
	}

	if u.Username != username {
		fmt.Println(u)
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
//	TODO Implement bcrypt authentication
func (u User) Login(password string) bool {
	if u.Password == password {
		return true
	}

	return false
}
