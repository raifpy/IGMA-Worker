package types

import "time"

type DatabaseAuth struct {
	Token string `json:"token" mongo:"token"`
	IP    string `json:"ip" mongo:"ip"`
	Id    int64  `json:"id" mongo:"id"`
	//Baerer     string        `json:"baerer" mongo:"baerer"` //??Belki
	Until time.Time `json:"until" mongo:"until"`
}
