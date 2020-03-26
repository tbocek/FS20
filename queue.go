package main

import (
	"fmt"
	"time"
)

func main() {
	messages := make(chan string)
	go func() {
		fmt.Println("waiting 1 sec")
		time.Sleep(time.Second)
		messages <- "ping"
	}()
	msg := <-messages
	fmt.Println(msg)
}
