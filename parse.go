package main

import (
	"client/types"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Parser interface {
	Parse(types.WebsocketContact, *websocket.Conn) error
}

func parse(t types.WebsocketContact, c *websocket.Conn) {
	res, _ := json.MarshalIndent(t, "", " ")
	fmt.Println(string(res))

	var parser Parser

	if t.Type == "newjob" {
		parser = NewJobParser{}
	} else {
		log.Println("ELSE!")
		return
	}

	if err := parser.Parse(t, c); err != nil {

		(ClientErrorParser{
			Error: err,
		}.Parse(t, c))
	} else {
		if t.NewJob != nil {
			t.NewJob.Status = "done"
		}
		c.WriteJSON(types.WebsocketContact{
			Type: "done",
			Update: &types.WebsocketUpdateJobStatus{
				Job: *t.NewJob,
			},
		})
	}
}
