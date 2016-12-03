package main

type ApiMessage struct {
	Status string `json:"status"`
}

type ConsoleOutputMessage struct {
	ApiMessage

	Output string `json:"output"`
}
