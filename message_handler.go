package main

import (
	"net/http"
	"strconv"

	"github.com/bigroom/vision/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/paked/gerrycode/communicator"
)

func messagesHandler(w http.ResponseWriter, r *http.Request, t *jwt.Token) {
	coms := communicator.New(w)

	vars := mux.Vars(r)
	channel := vars["channel"]
	host := vars["host"]

	channelKey := host + "/#" + channel

	offsetString := r.FormValue("offset")
	offset, err := strconv.ParseInt(offsetString, 10, 64)
	if err != nil {
		coms.Fail("Could not parse offset as a number")
		return
	}

	msgs, err := models.Messages(channelKey, offset)
	if err != nil {
		coms.Fail("Could not fetch messages")
		return
	}

	coms.OKWithData("Here are your messages", msgs)
}
