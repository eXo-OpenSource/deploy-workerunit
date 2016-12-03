package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Api struct {
	MTAServer *MTAServer
}

func NewApi(mtaServer *MTAServer) *Api {
	api := new(Api)
	api.MTAServer = mtaServer

	// Bind routes
	api.BindRoutes()

	return api
}

func (api *Api) Listen() {
	http.ListenAndServe(":8080", nil)
}

func (api *Api) BindRoutes() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "This is the API manuel (TODO)!")
	})

	http.HandleFunc("/start", func(res http.ResponseWriter, req *http.Request) {
		err := api.MTAServer.Start()

		api.SendErrorOrOkMessage(&res, err)
	})

	http.HandleFunc("/stop", func(res http.ResponseWriter, req *http.Request) {
		err := api.MTAServer.Stop()

		api.SendErrorOrOkMessage(&res, err)
	})

	http.HandleFunc("/output", func(res http.ResponseWriter, req *http.Request) {
		output := api.MTAServer.TailBuffer()

		json.NewEncoder(res).Encode(ConsoleOutputMessage{ApiMessage: ApiMessage{Status: "OK"}, Output: output})
	})

	http.HandleFunc("/command", func(res http.ResponseWriter, req *http.Request) {
		req.ParseForm()

		if req.Method != "POST" {
			api.SendErrorOrOkMessage(&res, errors.New("Bad method"))
		} else {
			command := req.Form.Get("command")
			err := api.MTAServer.ExecCommand(command)

			api.SendErrorOrOkMessage(&res, err)
		}
	})
}

func (api *Api) SendOkMessage(res *http.ResponseWriter) {
	json.NewEncoder(*res).Encode(ApiMessage{Status: "OK"})
}

func (api *Api) SendErrorOrOkMessage(res *http.ResponseWriter, err error) {
	if err != nil {
		json.NewEncoder(*res).Encode(ApiMessage{Status: err.Error()})
	} else {
		json.NewEncoder(*res).Encode(ApiMessage{Status: "OK"})
	}
}
