package main

import (
	"banner-service/internal/app"
	"log"
)

func main() {
	if err := app.Start(); err != nil {
		log.Fatalf("failed to start app %v", err)
	}
}
