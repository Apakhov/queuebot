package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"queuebot/fakesite"
)

var token = "958848651:AAGMkNg7_tXzUNTRu4qt0LYuywQi1h5ntGs"

func main() {
	fakesite.Start(os.Getenv("PORT"))
	bot, err := newQueueBot(token, "")
	if err != nil {
		log.Fatal("ERROR: ", err)
	}

	//interuption save
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Println("recieved ^C, shuting down")
			bot.Stop()
		}
	}()

	//panic save
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			bot.Stop()
			panic(r)
		}
	}()

	bot.ListenUsers()

	bot.Wait()

}
