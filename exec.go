package main

import (
	"client/types"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

func execExist(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

type Args struct {
	NewArgs []string

	DownloadFile string
	UploadFile   string

	DownloadLocalPath string
	UploadLocalPath   string

	DownloadID string
	UploadID   string

	WorkDir string
}

func argsParser(t types.WebsocketContact) (a Args, err error) {
	a = Args{
		NewArgs: []string{},
	}

	fmt.Println("orj args: ", t.NewJob.Exec.Args)

	for _, value := range t.NewJob.Exec.Args {

		if !strings.Contains(value, "$media") {
			a.NewArgs = append(a.NewArgs, value)
			continue
		}
		spl := strings.Split(value, ":")
		if len(spl) <= 1 {

			continue
		}
		fmt.Printf("1 value: %v\n", value)
		values, err := url.ParseQuery(spl[1])
		if err != nil {
			log.Println("url:parsequery error: ", err)
		}

		var filename string = values.Get("filename")

		if filename == "" {
			filename = RandStringRunes(20)
		}

		a.WorkDir = path.Join("mediajob", fmt.Sprint(t.NewJob.JobID))
		if err := os.MkdirAll(a.WorkDir, os.ModePerm); err != nil {
			return a, err
		}
		f, err := os.Create(path.Join(a.WorkDir, filename))
		if err != nil {
			return a, err
		}
		time.AfterFunc(time.Hour, func() {
			log.Println("File removing: ", path.Join(a.WorkDir, filename))
			os.Remove(path.Join(a.WorkDir, filename))
		})

		defer f.Close()

		if spl[0] == "$mediain" {
			log.Println("Media in")
			party3 := values.Get("3party")
			downloadid := values.Get("downloadid")

			fmt.Printf("party3: %v\n", party3)
			fmt.Printf("downloadid: %v\n", downloadid)

			if party3 != "" {
				log.Println("part3 not empty;")
				res, err := http.Get(party3)
				if err != nil {
					return a, err
				}
				fmt.Println(res.StatusCode)
				defer res.Body.Close()
				_, err = io.Copy(f, res.Body)
				if err != nil {
					return a, err
				}

				a.DownloadFile = party3
				a.DownloadLocalPath = filename

			}
			if downloadid != "" {
				req, err := http.NewRequest("GET", şema+Host+"/worker/download?id="+downloadid+"&token="+values.Get("token"), nil)
				if err != nil {
					return a, err
				}

				req.Header.Set("Token", Token)
				req.Header.Set("Id", Id)

				res, err := http.DefaultClient.Do(req)
				if err != nil {
					return a, err
				}

				defer res.Body.Close()

				log.Println("DownloadMedia Response header: ", res.Header)

				buffer := make([]byte, 1024*1024*10)

				for {
					size, err := res.Body.Read(buffer)
					f.Write(buffer[0:size])
					if err != nil {
						fmt.Printf("err: %v\n", err)
						break
					}
				}

				/*_, err = io.Copy(io.MultiWriter(f, progressbar.Default(res.ContentLength, "downloading: "+downloadid)), res.Body)
				if err != nil {
					return a, err
				}*/

				fdel := f.Name()

				time.AfterFunc(time.Hour*2, func() {
					os.Remove(fdel)
				})

				a.DownloadID = downloadid
				a.DownloadLocalPath = filename
			}

			f.Close()
			//value = path.Join(a.WorkDir, filename)
		}
		if spl[0] == "$mediaout" {

			a.UploadFile = filename
			a.UploadID = values.Get("uploadid")
		}
		value = path.Join(a.WorkDir, filename)
		a.NewArgs = append(a.NewArgs, value)

	}

	fmt.Println("\n\n ")
	fmt.Println("exec: ", t.NewJob.Exec.Exec)
	fmt.Println("args: ", a.NewArgs)

	return
}
