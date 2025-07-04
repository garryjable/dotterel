package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gplover/config"
	"gplover/dictionary"
	"gplover/engine"
	"gplover/machine"
)

func main() {

	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Load your dictionary
	dict, err := dictionary.LoadDictionaries("dictionaries")
	if err != nil {
		log.Fatalf("Error loading dictionary: %v", err)
	}
	e := engine.NewEngine(dict)

	gemini := machine.NewGeminiPrMachine(cfg.Port, cfg.Baud, func(keys []string) {
		word := e.TranslateSteno(keys)
		fmt.Print(word + " ")
	})

	// Start machine capture
	err = gemini.StartCapture()
	if err != nil {
		log.Fatalf("Failed to start Gemini PR machine: %v", err)
	}
	defer gemini.StopCapture()

	fmt.Println("Dotterel now running. Press Ctrl+C to quit.")

	// Handle Ctrl+C
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	<-sigs
	fmt.Println("\n[dotterel] Quit with Ctrl+C")
}
