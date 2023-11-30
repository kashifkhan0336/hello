package main

import (
	"log"
	"os/exec"

	"gopkg.in/natefinch/npipe.v2"
)

var mpvPipeName string = ""

var title string = ""

func InitializeVideo() {
	println("From initvideo : ")
	println(host)
	if host {
		title = "Host Player"
		mpvPipeName = `\\.\pipe\mpvhost`
	} else {
		title = "Client Player"
		mpvPipeName = `\\.\pipe\mpvclient`
	}

	cmd := exec.Command("./lib/mpv.exe", "--input-ipc-server="+mpvPipeName, "--force-window", "--title="+title, "--idle", "--hwdec=auto", "file.mp4")
	err := cmd.Start()
	if err != nil {
		log.Fatal("Error starting MPV:", err)
	}

}

func InitializeVideoIpc() (*npipe.PipeConn, error) {
	pipe, err := npipe.Dial(mpvPipeName)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to MPV!")
	return pipe, nil
}
