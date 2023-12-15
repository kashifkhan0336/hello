package main

import (
	"fmt"

	"github.com/centrifugal/centrifuge-go"
)

func InitializeClient(ipAddress string) *centrifuge.Client {
	client := centrifuge.NewJsonClient(
		fmt.Sprintf("ws://%s:8000/connection/websocket", ipAddress),
		centrifuge.Config{
			// Sending token makes it work with Centrifugo JWT auth (with `secret` HMAC key).
			Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxIiwiZXhwIjoyMDYxMDIwNDcwLCJpYXQiOjE3MDEwMjA0NzB9.f8JJnMM8a-_ftMlh4VrEZkMUKonNNEobf_1rbZDX2PQ",
		},
	)

	return client
}
