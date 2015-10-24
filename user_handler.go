package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bigroom/vision/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/paked/gerrycode/communicator"
	"github.com/paked/restrict"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	coms := communicator.New(w)

	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	fmt.Println(username, password, email)

	u, err := models.NewUser(username, password, email)
	if err != nil {
		log.Println(err)
		coms.Error("Unable to create user")
		return
	}

	coms.OKWithData("user", u)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	coms := communicator.New(w)

	username := r.FormValue("username")
	password := r.FormValue("password")

	u, err := models.FetchUser("username", username)
	if err != nil {
		coms.Error("Unable to login user")
		return
	}

	if !u.Login(password) {
		coms.Error("Unable to login")
		return
	}

	claims := make(map[string]interface{})
	claims["id"] = u.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	ts, err := restrict.Token(claims)
	if err != nil {
		coms.Fail("Failure signing the token")
		return
	}

	coms.OKWithData("token", ts)
}

func secretHandler(w http.ResponseWriter, r *http.Request, t *jwt.Token) {
	coms := communicator.New(w)

	u, err := models.FetchUser("id", t.Claims["id"])
	if err != nil {
		coms.Error("That user does not exist!")
		return
	}

	coms.OKWithData("ID", u)
}
