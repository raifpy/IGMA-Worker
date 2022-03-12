package main

import (
	"bytes"
	"client/types"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/gorilla/websocket"
)

type ExecParser struct {
}

func (n ExecParser) Parse(t types.WebsocketContact, conn *websocket.Conn) error {
	if !execExist(t.NewJob.Exec.Exec) {
		return fmt.Errorf("%s not exists on client", t.NewJob.Exec.Exec)
	}

	ex := t.NewJob.Exec

	conn.WriteJSON(types.WebsocketContact{
		Type: "downloading",
		Update: &types.WebsocketUpdateJobStatus{
			Job: *t.NewJob,
		},
	})
	args, err := argsParser(t)
	if err != nil {
		return err
	}

	fmt.Println(args)
	var buf = &bytes.Buffer{}
	var errbuf = &bytes.Buffer{}
	cmd := exec.Command(ex.Exec, args.NewArgs...)
	//cmd.Stderr = io.MultiWriter(os.Stderr, errbuf)
	//cmd.Stdout = io.MultiWriter(os.Stdout, buf)
	cmd.Stderr = errbuf
	cmd.Stdout = buf
	cmd.Stdin = nil

	conn.WriteJSON(types.WebsocketContact{
		Type: "rendering",
		Update: &types.WebsocketUpdateJobStatus{
			Job: *t.NewJob,
		},
	})

	if err := cmd.Run(); err != nil {
		return errors.New(errbuf.String())
	}

	body := &bytes.Buffer{} // !! Memory

	multiwriter := multipart.NewWriter(body)

	mediawriter, err := multiwriter.CreateFormFile("media", args.UploadFile)
	if err != nil {
		return err
	}
	defer multiwriter.Close() //Yey

	upfile, err := os.Open(path.Join(args.WorkDir, args.UploadFile))
	if err != nil {
		return err
	}
	defer upfile.Close()
	size, err := io.Copy(mediawriter, upfile)
	if err != nil {
		return err
	}

	if stdout, err := multiwriter.CreateFormFile("stdout", "stdout"); err == nil {
		io.Copy(stdout, buf)
	}
	multiwriter.Close()
	fmt.Printf("%d  size copied\n", size)

	request, err := http.NewRequest("POST", şema+Host+"/worker/upload?id="+args.UploadID, body)
	if err != nil {
		return err
	}
	request.Header.Set("token", Token)
	request.Header.Set("id", Id)
	//request.Header.Set("Content-Type", multiwriter.FormDataContentType())

	request.Header.Set("Content-Type", multiwriter.FormDataContentType())

	fmt.Println("request headers: ", request.Header)

	conn.WriteJSON(types.WebsocketContact{
		Type: "uploading",
		Update: &types.WebsocketUpdateJobStatus{
			Job: *t.NewJob,
		},
	})

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}

	fmt.Println(res.Header)

	return nil
}
