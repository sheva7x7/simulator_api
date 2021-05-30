package main

import (
	ESClient "wenle/elasticsearch/esclient"
	"wenle/elasticsearch/mqtt"
	"wenle/elasticsearch/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	ESClient.StartClient()
	mqtt.StartClient()
	server.HandleRequest()
}
