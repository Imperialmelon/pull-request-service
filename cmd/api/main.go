package main

import (
	"log"

	"github.com/Imperialmelon/AvitoTechTest/cmd/api/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("application error: %v", err)
	}
}
