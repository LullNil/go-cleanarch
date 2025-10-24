package main

import (
	"log"
	"os"

	"github.com/LullNil/go-cleanarch/config"
	"github.com/LullNil/go-cleanarch/internal/app"
)

func main() {
	// Init config
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Run application
	if err := app.Run(cfg); err != nil {
		log.Printf("application stopped with error: %v", err)
		os.Exit(1)
	}

	log.Println("application exited gracefully")
}
