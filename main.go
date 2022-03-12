package main

import (
	"client/types"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var (
	Token = ""
	Id    = ""

	Host   = "<HOST>"
	şema   = "https://"
	wsşema = "wss://"

	/*şema   = "http://"
	wsşema = "ws://"
	Host   = "localhost:2095"*/
)

func init() {
	rand.Seed(time.Now().Unix())
}

func atmain() {
	conn, res, err := websocket.DefaultDialer.Dial(wsşema+Host+"/worker/wsworker", http.Header{
		"Token": {Token},
		"Id":    {Id},
	})
	if err != nil {

		log.Println(err)

		if res != nil {
			fmt.Println(res.Status)
			//fmt.Printf("res.Request.Header: %v\n", res.Request.Header)
			res.Body.Close()

		}

		return
	}

	defer res.Body.Close()
	defer conn.Close()

	fmt.Println("Ayaktayım")
	var online = true

	go func() {
		for {
			if !online {
				break
			}
			conn.WriteMessage(websocket.PingMessage, nil)
			time.Sleep(time.Second * 5)
		}
	}()

	for {
		var t types.WebsocketContact
		if err := conn.ReadJSON(&t); err != nil {
			log.Println(err)

			break
		}

		go parse(t, conn)

	}
	online = false
}

func main() {
	j, err := os.ReadFile(".json")
	if err != nil {
		log.Fatalln(err, " where is meta?")
	}
	var meta struct {
		Token string `json:"token"`
		Id    int64  `json:"id"`
	}
	if err := json.Unmarshal(j, &meta); err != nil {
		log.Fatalln(err)
	}
	Token = meta.Token
	Id = fmt.Sprint(meta.Id)

	for { //!! !!
		atmain()
		log.Println("re-connect: 2+")
		time.Sleep(time.Second * 2)
	}
}
