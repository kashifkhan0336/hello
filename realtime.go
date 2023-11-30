package main

import "github.com/centrifugal/centrifuge-go"

func InitializeClient() *centrifuge.Client {
	client := centrifuge.NewJsonClient(
		"ws://26.128.237.142:8000/connection/websocket",
		centrifuge.Config{
			// Sending token makes it work with Centrifugo JWT auth (with `secret` HMAC key).
			Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxIiwiZXhwIjoyMDYxMDIwNDcwLCJpYXQiOjE3MDEwMjA0NzB9.f8JJnMM8a-_ftMlh4VrEZkMUKonNNEobf_1rbZDX2PQ",
		},
	)

	return client
}
