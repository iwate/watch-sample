package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT)

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Start watch ...", dir)
	log.Println("[Ctrl+c]: Finish watch")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	go func() {
		for {
			select {
			case evt := <-watcher.Events:
				log.Println("event:", evt)
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(`C:\Users\iwate\Go\src\github.com\hackm\watch-sample`)
	if err != nil {
		log.Fatal(err)
	}
	<-done
	log.Println("Bye!")
}
