package main

import (
	"log"

	"github.com/imperialmelon/avito/cmd/api/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("application error: %v", err)
	}
}
