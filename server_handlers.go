package main

import (
	"net/http"

	"github.com/paked/gerrycode/communicator"
)

func defaultServerHandler(w http.ResponseWriter, r *http.Request) {
	coms := communicator.New(w)

	resp := struct {
		Server string `json:"server"`
	}{
		Server: *defaultIRCServer,
	}

	coms.OKWithData("server", resp)
}
