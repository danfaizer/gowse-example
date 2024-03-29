package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"golang.org/x/net/websocket"
)

var (
	port = flag.String("port", "9001", "port used for ws connection")
)

func main() {
	flag.Parse()

	// connect
	ws, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// receive
	var m map[string]interface{}
	done := make(chan error)
	go func() {
		var err error
		for {
			err = websocket.JSON.Receive(ws, &m)
			if err != nil {
				fmt.Println("Error receiving message: ", err.Error())
				break
			}
			fmt.Println("Message: ", m)
		}
		done <- err
	}()
	<-done
}

// connect connects to the local chat server at port <port>
func connect() (*websocket.Conn, error) {
	return websocket.Dial(fmt.Sprintf("ws://localhost:%s", *port), "", mockedIP())
}

// mockedIP is a demo-only utility that generates a random IP address for this client
func mockedIP() string {
	var arr [4]int
	for i := 0; i < 4; i++ {
		rand.Seed(time.Now().UnixNano())
		arr[i] = rand.Intn(256)
	}
	return fmt.Sprintf("http://%d.%d.%d.%d", arr[0], arr[1], arr[2], arr[3])
}
