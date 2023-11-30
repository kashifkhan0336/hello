package main

import (
	"fmt"
	"log"

	"github.com/centrifugal/centrifuge-go"
	"github.com/goccy/go-json"
	"gopkg.in/natefinch/npipe.v2"
)

type Event struct {
	Event  string      `json:"event,omitempty"`
	ID     int         `json:"id,omitempty"`
	Sender string      `json:"sender,omitempty"`
	Name   string      `json:"name,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func configureSubscriptionEvents(Subscription *centrifuge.Subscription, pipe *npipe.PipeConn) {
	Subscription.OnSubscribing(func(e centrifuge.SubscribingEvent) {
		log.Printf("Subscribing on channel %s - %d (%s)", Subscription.Channel, e.Code, e.Reason)
	})
	Subscription.OnSubscribed(func(e centrifuge.SubscribedEvent) {
		log.Printf("Subscribed on channel %s, (was recovering: %v, recovered: %v)", Subscription.Channel, e.WasRecovering, e.Recovered)
	})
	Subscription.OnUnsubscribed(func(e centrifuge.UnsubscribedEvent) {
		log.Printf("Unsubscribed from channel %s - %d (%s)", Subscription.Channel, e.Code, e.Reason)
	})

	Subscription.OnError(func(e centrifuge.SubscriptionErrorEvent) {
		log.Printf("Subscription error %s: %s", Subscription.Channel, e.Error)
	})
	Subscription.OnPublication(func(e centrifuge.PublicationEvent) {
		var event *Event
		// // err := json.Unmarshal(e.Data, &chatMessage)
		var jsonMap map[string]interface{}
		err := json.Unmarshal([]byte(e.Data), &jsonMap)
		if err != nil {
			log.Printf("error occured! %s", err)
		}
		//

		err = json.Unmarshal([]byte(jsonMap["input"].(string)), &event)
		if err != nil {
			log.Printf("error occured! %s", err)
			return
		}
		fmt.Println(event)
		if event.Event == "sync_request" {
			println("Sync request received sending sync response")
		}

		if event.Event == "play" && !host {
			fmt.Printf("Video played at %s", event.Data)
			sendRequest(pipe, MPVRequest{
				Command: []interface{}{"set_property", "pause", false},
			})
			sendRequest(pipe, MPVRequest{
				Command: []interface{}{"set_property", "time-pos", event.Data},
			})
		}
		if event.Event == "pause" && !host {
			fmt.Printf("Video paused at %s", event.Data)

			sendRequest(pipe, MPVRequest{
				Command: []interface{}{"set_property", "pause", true},
			})
			sendRequest(pipe, MPVRequest{
				Command: []interface{}{"set_property", "time-pos", event.Data},
			})
		}
		if event.Event == "seek" && !host {
			fmt.Printf("Video seeked at %s", event.Data)
			sendRequest(pipe, MPVRequest{
				Command: []interface{}{"set_property", "time-pos", event.Data},
			})
		}
		if event.Event == "sync" && !host {
			fmt.Printf("Syncing video at %s", event.Data)
			sendRequest(pipe, MPVRequest{
				Command: []interface{}{"show_text", "Syncing"},
			})
			sendRequest(pipe, MPVRequest{
				Command: []interface{}{"set_property", "time-pos", event.Data},
			})
		}
		// fmt.Println(event.Event)
		// fmt.Println(event.Data)
		// fmt.Println(event.Name)

		//fmt.Printf("Name : %s, Event : %s, Data : %s", event.Name, event.Event, event.Data)

	})
	Subscription.OnJoin(func(e centrifuge.JoinEvent) {
		log.Printf("Someone joined %s: user id %s, client id %s", Subscription.Channel, e.User, e.Client)
	})
	Subscription.OnLeave(func(e centrifuge.LeaveEvent) {
		log.Printf("Someone left %s: user id %s, client id %s", Subscription.Channel, e.User, e.Client)
	})
}
func configureClientEvents(client *centrifuge.Client) {
	client.OnConnecting(func(e centrifuge.ConnectingEvent) {
		log.Printf("Connecting - %d (%s)", e.Code, e.Reason)
	})
	client.OnConnected(func(e centrifuge.ConnectedEvent) {
		log.Printf("Connected with ID %s", e.ClientID)
	})
	client.OnDisconnected(func(e centrifuge.DisconnectedEvent) {
		log.Printf("Disconnected: %d (%s)", e.Code, e.Reason)
	})

	client.OnError(func(e centrifuge.ErrorEvent) {
		log.Printf("Error: %s", e.Error.Error())
	})

	client.OnMessage(func(e centrifuge.MessageEvent) {
		log.Printf("Message from server: %s", string(e.Data))
	})

	client.OnSubscribed(func(e centrifuge.ServerSubscribedEvent) {
		log.Printf("Subscribed to server-side channel %s: (was recovering: %v, recovered: %v)", e.Channel, e.WasRecovering, e.Recovered)
	})
	client.OnSubscribing(func(e centrifuge.ServerSubscribingEvent) {
		log.Printf("Subscribing to server-side channel %s", e.Channel)
	})
	client.OnUnsubscribed(func(e centrifuge.ServerUnsubscribedEvent) {
		log.Printf("Unsubscribed from server-side channel %s", e.Channel)
	})

	client.OnPublication(func(e centrifuge.ServerPublicationEvent) {
		log.Printf("Publication from server-side channel %s: %s (offset %d)", e.Channel, e.Data, e.Offset)
	})
	client.OnJoin(func(e centrifuge.ServerJoinEvent) {
		log.Printf("Join to server-side channel %s: %s (%s)", e.Channel, e.User, e.Client)
	})
	client.OnLeave(func(e centrifuge.ServerLeaveEvent) {
		log.Printf("Leave from server-side channel %s: %s (%s)", e.Channel, e.User, e.Client)
	})
}
