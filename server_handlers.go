package main

import (
	"net/http"

	"github.com/bigroom/communicator"
)

func defaultServerHandler(w http.ResponseWriter, r *http.Request) {
	coms := communicator.New(w)

	resp := struct {
		Server string `json:"server"`
	}{
		Server: *defaultIRCServer,
	}

	coms.With(resp).
		OK("Here is your thing!")
}
