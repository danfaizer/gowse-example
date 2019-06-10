package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danfaizer/gowse"
)

type Check struct {
	ID            string
	ChecktypeName string
}

func main() {
	l := log.New(os.Stdout, "", log.LstdFlags)
	s := gowse.NewServer(l)
	t := s.CreateTopic("test")
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := t.TopicHandler(w, r); err != nil {
			l.Printf("error handling subcriber request: %+v", err)
		}
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

	wss := http.Server{Addr: ":" + "9001", Handler: mux}
	done := make(chan error)
	go func() {
		err := wss.ListenAndServe()
		done <- err
	}()
	sg := make(chan os.Signal)
	signal.Notify(sg, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sg
		wss.Shutdown(context.Background())
	}()
	err := <-done
	if err != nil && err != http.ErrServerClosed {
		fmt.Printf("error stoping http server: %+v", err)
	}
	fmt.Printf("waiting gowse to stop")
	s.Stop()
}