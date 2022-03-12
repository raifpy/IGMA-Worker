package types

import (
	"strings"
	"time"
)

type Job struct { //Client-Server
	Status  string    `json:"status"`
	JobID   int64     `json:"job_id"`
	Expired time.Time `json:"expired"`
	Exec    *Exec     `json:"exec"` // nilable
}

type Exec struct {
	Exec     string   // wget
	Args     []string // $mediain:3party=https://codeksion.net/images/project/phishdroid/phishdroid.png&filename=in.png $mediaup:uploadid=<JobID>&filename=out.mp4
	ShareSTD bool     // stdout:
}

type ExecArgs []string

func (e ExecArgs) IsMediaExec() bool {
	for _, args := range e {
		if strings.Contains(args, "$media") {
			return true
		}
	}
	return false
}
