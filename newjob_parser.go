package main

import (
	"client/types"
	"errors"

	"github.com/gorilla/websocket"
)

type NewJobParser struct {
}

func (n NewJobParser) Parse(t types.WebsocketContact, conn *websocket.Conn) error {
	job := t.NewJob
	if job.JobID == 0 {
		return errors.New("job_id cannot 0")
	}

	switch {
	case job.Exec != nil:
		return ExecParser{}.Parse(t, conn)

	}
	return errors.New("unexcepted response")
}
