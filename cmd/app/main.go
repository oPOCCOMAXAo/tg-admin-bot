package main

import (
	"log"

	"github.com/opoccomaxao/tg-admin-bot/pkg/app"
)

func main() {
	err := app.Run()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
