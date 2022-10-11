//go:build windows
// +build windows

package main

import (
	"fmt"
	winlog "github.com/ofcoursedude/gowinlog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	watcher, err := winlog.NewWinLogWatcher()
	if err != nil {
		fmt.Printf("Couldn't create watcher: %v\n", err)
		return
	}
	err = watcher.SubscribeFromBeginning("Application", "*")
	if err != nil {
		fmt.Printf("Couldn't subscribe to Application: %v", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case evt := <-watcher.Event():
			fmt.Printf("Event: %v\n", evt)
			bookmark := evt.Bookmark
			fmt.Printf("Bookmark: %v\n", bookmark)
		case err := <-watcher.Error():
			fmt.Printf("Error: %v\n\n", err)
		default:
		}

		select {
		case <-c:
			fmt.Println("Shutting down")
			watcher.Shutdown()
			return
		default:
			// If no event is waiting, need to wait or do something else, otherwise
			// the the app fails on deadlock.
			<-time.After(1 * time.Millisecond)
		}
	}
}
