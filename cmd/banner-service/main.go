package main

import (
	"banner-service/internal/app"
	"log"
)

func main() {
	a, err := app.NewApp()
	if err != nil {
		log.Fatalf("failed to init app %v", err)
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app %v", err)
	}
}
