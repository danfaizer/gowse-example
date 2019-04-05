package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/danfaizer/gowse"
)

type Check struct {
	ID            string
	ChecktypeName string
}

func main() {
	s := gowse.New()
	t := s.CreateTopic("test")
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		gowse.TopicHandler(t, w, r)
	})

	go func() {
		c := Check{ID: "aaaa", ChecktypeName: "bbb"}
		time.Sleep(10 * time.Second)
		ticker := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-ticker.C:
				t.Broadcast(c)
			}
		}
	}()

	wss := http.Server{Addr: ":" + "9000", Handler: mux}
	err := wss.ListenAndServe()
	fmt.Println(err)
}
