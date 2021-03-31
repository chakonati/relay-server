package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"server/handlers"
	"server/messaging"
	"server/persistence"
	"server/session"
	"strconv"

	"nhooyr.io/websocket"
)

var listenPort = 4560
var listenAddr = "0.0.0.0"

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

func useEnv() {
	if listenPortEnv, exists := os.LookupEnv("PORT"); exists {
		port, err := strconv.Atoi(listenPortEnv)
		if err != nil {
			log.Fatalf("Invalid port: %s", listenPortEnv)
		}
		listenPort = port
	}
	if listenAddrEnv, exists := os.LookupEnv("LISTEN_ADDR"); exists {
		_, _, err := net.SplitHostPort(fmt.Sprintf("%s:1", listenAddrEnv))
		if err != nil {
			log.Fatalf("Invalid listen address: %s", listenAddrEnv)
		}
		listenAddr = listenAddrEnv
	}
}

func startHTTP() {
	http.HandleFunc("/", handler)

	listenOn := fmt.Sprintf("%s:%d", listenAddr, listenPort)
	log.Println("Listening on", listenOn)
	log.Fatal(http.ListenAndServe(listenOn, nil))
}

var relay = messaging.Relay{}

func startWorkers() {
	go relay.StartWorking()
}

func main() {
	useEnv()
	if err := persistence.InitDatabases(); err != nil {
		log.Fatal("Failed to initialize databases:", err)
	}
	startWorkers()
	startHTTP()
}
