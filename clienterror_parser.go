package main

import (
	"client/types"
	"log"

	"github.com/gorilla/websocket"
)

type ClientErrorParser struct {
	Error error
}

func (c ClientErrorParser) Parse(t types.WebsocketContact, conn *websocket.Conn) error {
	log.Println("Client Error: ", c.Error)
	job := t.NewJob
	if job == nil {
		job.Status = "error"
	}

	return conn.WriteJSON(types.WebsocketContact{
		Type: "error",
		Error: &types.WebsocketError{
			Error: c.Error.Error(),
			Job:   job,
		},
	})

}
