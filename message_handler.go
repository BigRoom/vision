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

	pageString := r.FormValue("page")
	page, err := strconv.ParseInt(pageString, 10, 64)
	if err != nil {
		coms.Fail("Could not parse page as a number")
		return
	}

	msgs, err := models.Messages(channelKey, page)
	if err != nil {
		coms.Fail("Could not fetch messages")
		return
	}

	coms.OKWithData("Here are your messages", msgs)
}
