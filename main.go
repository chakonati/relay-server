package main

import (
	"fmt"
	"log"
	"net/http"
	"server/configuration"
	"server/handlers"
	"server/messaging"
	"server/persistence"
	"server/session"

	"nhooyr.io/websocket"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received connection")
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols:         nil,
		InsecureSkipVerify:   false,
		OriginPatterns:       []string{"*"},
		CompressionMode:      0,
		CompressionThreshold: 0,
	})
	if err != nil {
		log.Println(err)
		return
	}

	go (&handlers.Handler{
		Session: session.Session{
			State: session.Created,
			Conn:  c,
		},
	}).Handle()
}

func startHTTP() {
	http.HandleFunc("/", handler)

	listenOn := fmt.Sprintf("%s:%d", configuration.Config().ListenAddr, configuration.Config().ListenPort)
	log.Println("Listening on", listenOn)
	log.Fatal(http.ListenAndServe(listenOn, nil))
}

var relay = messaging.Relay{}

func startWorkers() {
	go relay.StartWorking()
}

func main() {
	if err := configuration.LoadConfigFromEnv(); err != nil {
		log.Fatal("Failed to load configuration from environment: ", err)
	}
	if err := persistence.InitDatabases(); err != nil {
		log.Fatal("Failed to initialize databases: ", err)
	}
	startWorkers()
	startHTTP()
}
