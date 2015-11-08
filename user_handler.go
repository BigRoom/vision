package main

import (
	"net/http"
	"time"

	"github.com/bigroom/communicator"
	"github.com/bigroom/vision/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/paked/restrict"
	log "github.com/sirupsen/logrus"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	coms := communicator.New(w)

	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	log.WithFields(log.Fields{
		"username": username,
		"email":    email,
	}).Info("User is tryingt to register")

	u, err := models.NewUser(username, password, email)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to connect to user")

		coms.Error("Unable to create user")
		return
	}

	coms.With(u).OK("user")
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

	if err := u.Login(password); err != nil {
		coms.Errorf("Unable to login %v", err)
		return
	}

	claims := make(map[string]interface{})
	claims["id"] = u.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	ts, err := restrict.Token(claims)
	if err != nil {
		coms.Fail("Failure signing the token")
		sentry.CaptureError(err, nil)
		return
	}

	coms.With(ts).OK()
}

func secretHandler(w http.ResponseWriter, r *http.Request, t *jwt.Token) {
	coms := communicator.New(w)

	u, err := models.FetchUser("id", t.Claims["id"])
	if err != nil {
		coms.Error("That user does not exist!")
		return
	}

	coms.With(u).OK()
}
