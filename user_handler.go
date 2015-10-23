package main

import (
	"net/http"

	"github.com/bigroom/vision/models"
	"github.com/paked/gerrycode/communicator"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	coms := communicator.New(w)

	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	u, err := models.NewUser(username, password, email)
	if err != nil {
		coms.Error("Unable to create user")
		return
	}

	coms.OKWithData("user", u)
}
