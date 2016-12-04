package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var apiSecret string

type Api struct {
	MTAServer *MTAServer
}

func NewApi(mtaServer *MTAServer) *Api {
	api := new(Api)
	api.MTAServer = mtaServer

	// Bind routes
	api.BindRoutes()

	// Get API secret
	apiSecret = os.Getenv("API_SECRET")

	return api
}

func (api *Api) Listen() {
	http.ListenAndServe(":8080", nil)
}

func (api *Api) BindRoutes() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, `
		Help:
		- /start : Starts the MTA server
		- /stop : Stops the MTA server
		- /restart : Restarts the MTA server (waits until stopped and starts then)
		- /logs : Retrieves a the last n lines of the standard output (uses a ring buffer internally)
		- /command : Execute a command on the server's console
		`)
	})

	http.HandleFunc("/start", func(res http.ResponseWriter, req *http.Request) {
		if !api.CheckAPISecret(req) {
			api.SendStatusMessage(&res, "Wrong API secret")
			return
		}
		err := api.MTAServer.Start()

		api.SendStatusError(&res, err)
	})

	http.HandleFunc("/stop", func(res http.ResponseWriter, req *http.Request) {
		if !api.CheckAPISecret(req) {
			api.SendStatusMessage(&res, "Wrong API secret")
			return
		}
		err := api.MTAServer.Stop()

		api.SendStatusError(&res, err)
	})

	http.HandleFunc("/restart", func(res http.ResponseWriter, req *http.Request) {
		if !api.CheckAPISecret(req) {
			api.SendStatusMessage(&res, "Wrong API secret")
			return
		}
		err := api.MTAServer.Restart()

		api.SendStatusError(&res, err)
	})

	http.HandleFunc("/logs", func(res http.ResponseWriter, req *http.Request) {
		if !api.CheckAPISecret(req) {
			api.SendStatusMessage(&res, "Wrong API secret")
			return
		}
		output := api.MTAServer.TailBuffer()

		json.NewEncoder(res).Encode(ConsoleOutputMessage{ApiMessage: ApiMessage{Status: "OK"}, Output: output})
	})

	http.HandleFunc("/command", func(res http.ResponseWriter, req *http.Request) {
		if !api.CheckAPISecret(req) {
			api.SendStatusMessage(&res, "Wrong API secret")
			return
		}

		// Parse POST parameters
		req.ParseForm()

		if req.Method != "POST" {
			api.SendStatusMessage(&res, "Bad method")
		} else {
			command := req.Form.Get("command")
			err := api.MTAServer.ExecCommand(command)

			api.SendStatusError(&res, err)
		}
	})
}

func (api *Api) SendOkMessage(res *http.ResponseWriter) {
	json.NewEncoder(*res).Encode(ApiMessage{Status: "OK"})
}

func (api *Api) SendStatusMessage(res *http.ResponseWriter, message string) {
	json.NewEncoder(*res).Encode(ApiMessage{Status: message})
}

func (api *Api) SendStatusError(res *http.ResponseWriter, err error) {
	if err != nil {
		json.NewEncoder(*res).Encode(ApiMessage{Status: err.Error()})
	} else {
		json.NewEncoder(*res).Encode(ApiMessage{Status: "OK"})
	}
}

func (api *Api) CheckAPISecret(req *http.Request) bool {
	return req.Header.Get("API_SECRET") == apiSecret
}
