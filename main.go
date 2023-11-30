package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-json"

	"github.com/centrifugal/centrifuge-go"
	"gopkg.in/natefinch/npipe.v2"
)

var host bool = false

type Message struct {
	Input string `json:"input"`
}
type MPVRequest struct {
	Command []interface{} `json:"command"`
}

func GetPositionValue() (float64, error) {
	var position map[string]interface{}

	pipe, err := npipe.Dial(mpvPipeName)
	if err != nil {
		log.Printf("error occured while connecting from different goroutine %s", err)
	}
	sendRequest(pipe, GetPosition)
	msg, err := bufio.NewReader(pipe).ReadString('\n')
	if err != nil {
		println("error in goroutine")
		return 0, err
	}
	err = json.Unmarshal([]byte(msg), &position)
	if err != nil {
		fmt.Printf("Error : %s", err)
		return 0, err
	}
	if timePos, ok := position["data"].(float64); ok {
		return timePos, nil
		//pubText("{event: sync}", Subscription)
	}
	return 0, err
}
func sendRequest(pipe *npipe.PipeConn, request MPVRequest) error {
	// fmt.Print("from sendRequest : ")
	// fmt.Println(request)
	encoder := json.NewEncoder(pipe)
	err := encoder.Encode(request)
	return err
}
func ReadingRoutine(pipe *npipe.PipeConn, msgChannel chan string) {
	for {
		msg, err := bufio.NewReader(pipe).ReadString('\n')
		if err != nil {
			close(msgChannel)
			println("Player closed!")
			return
		}
		msgChannel <- msg

	}
}
func pubText(text string, sub *centrifuge.Subscription) error {
	msg := &Message{
		Input: text,
	}
	data, _ := json.Marshal(msg)
	_, err := sub.Publish(context.Background(), data)
	return err
}

type playbackEvent struct {
	Event  string `json:"event"`
	Sender string `json:"sender"`
	Data   string `json:"data"`
}

func createEvent(event, sender, dataValue string) string {
	eventStruct := playbackEvent{
		Event:  event,
		Sender: sender,
		Data:   dataValue,
	}

	// Marshal the Event struct into a JSON string
	jsonBytes, err := json.Marshal(eventStruct)
	if err != nil {
		return ""
	}

	return string(jsonBytes)
}

func handlePlaybackEvent(eventName string, sub *centrifuge.Subscription) {
	pos, err := GetPositionValue()
	s := strconv.FormatFloat(pos, 'f', -1, 64)
	if err != nil {
		fmt.Println("failed to get po")
		return
	}
	switch eventName {
	case "play":
		pubText(createEvent("play", "host", s), sub)
	case "pause":
		pubText(createEvent("pause", "host", s), sub)
	case "seek":
		pubText(createEvent("seek", "host", s), sub)
	case "sync":
		pubText(createEvent("sync", "host", s), sub)
	}

}

func main() {

	args := os.Args[1:]
	fmt.Println(args)
	fmt.Println(len(args))
	if len(args) == 1 {
		if args[0] == "host" {
			fmt.Println("hosting!")
			host = true
		}
	}
	println(host)
	InitializeVideo()
	pipe, err := InitializeVideoIpc()
	if err != nil {
		log.Fatal("Can't connect mpv ipc")
	}
	client := InitializeClient()
	defer client.Close()
	configureClientEvents(client)

	error_ := client.Connect()
	if error_ != nil {
		log.Fatalln(error_)
	}

	sub, err := client.NewSubscription("party", centrifuge.SubscriptionConfig{
		Recoverable: true,
		JoinLeave:   true,
	})
	go func() {
		for range time.Tick(time.Second * 60) {
			print("Executing syncing....")
			pos, err := GetPositionValue()
			s := strconv.FormatFloat(pos, 'f', -1, 64)
			if err != nil {
				fmt.Println("failed to get position")
				return
			}
			pubText(createEvent("sync", "host", s), sub)
		}
	}()

	configureSubscriptionEvents(sub, pipe)
	if err != nil {
		log.Fatalln(err)
	}
	//sendRequest(pipe, ObserveTimePositionProperty)

	sendRequest(pipe, ObserveSeekProperty)
	msgChannel := make(chan string)
	go ReadingRoutine(pipe, msgChannel)

	go func() {
		for {
			select {
			case msg, ok := <-msgChannel:
				if !ok {
					fmt.Println("Reading routine closed the channel.")
					os.Exit(0)
					return
				}
				if host {
					var raw map[string]interface{}
					err := json.Unmarshal([]byte(msg), &raw)
					if err != nil {
						print(err)
					}
					//fmt.Printf("raw msg -> %s", msg)
					//fmt.Print(raw)
					if _, ok := raw["event"]; ok {
						fmt.Println(raw["event"], raw["data"], raw["name"])
						if raw["name"] == "pause" {
							if raw["data"] == true {
								handlePlaybackEvent("pause", sub)

							} else {
								handlePlaybackEvent("play", sub)

							}
						}
						if raw["event"] == "playback-restart" {
							handlePlaybackEvent("seek", sub)
						}
					}
				}
			}
		}

	}()
	sendRequest(pipe, ObserveStateProperty)
	sendRequest(pipe, ObserveVolumeProperty)
	err = sub.Subscribe()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Print something and press ENTER to send\n")
	err = pubText("hello", sub)
	if err != nil {
		log.Printf("Error publish: %s", err)
	}
	go func(sub *centrifuge.Subscription) {
		reader := bufio.NewReader(os.Stdin)
		for {
			text, _ := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			switch text {
			case "#subscribe":
				err := sub.Subscribe()
				if err != nil {
					log.Println(err)
				}
			case "#volume":
				//sendRequest(pipe, GetVolumeRequest)
				log.Println("Observing volume")
			case "#unsubscribe":
				err := sub.Unsubscribe()
				if err != nil {
					log.Println(err)
				}
			case "#disconnect":
				err := client.Disconnect()
				if err != nil {
					log.Println(err)
				}
			case "#pos":
				go func() {
					var position map[string]interface{}

					print("hello from goroutine!")
					pipe, err := npipe.Dial(mpvPipeName)
					if err != nil {
						log.Printf("error occured while connecting from different goroutine %s", err)
					}
					log.Println("Connected to MPV from different goroutine")

					sendRequest(pipe, GetPosition)
					msg, err := bufio.NewReader(pipe).ReadString('\n')
					if err != nil {
						println("error in goroutine")
					}
					err = json.Unmarshal([]byte(msg), &position)
					if err != nil {
						fmt.Printf("Error : %s", err)
					}
					fmt.Printf("From seprate goroutine : %f", position["data"].(float64))
				}()
			case "#connect":
				err := client.Connect()
				if err != nil {
					log.Println(err)
				}
			case "#close":
				client.Close()
			default:
				err = pubText(text, sub)
				if err != nil {
					log.Printf("Error publish: %s", err)
				}
			}
		}
	}(sub)
	select {}

}
