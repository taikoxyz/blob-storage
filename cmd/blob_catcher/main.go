package main

import (
	"log"

	"github.com/taikoxyz/blob-storage/internal/indexer"
	"github.com/taikoxyz/blob-storage/internal/logic"
)

func main() {
	cfg, err := logic.GetConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	if err := indexer.InitFromConfig(cfg); err != nil {
		log.Fatal("Error running indexer:", err)
	}
}
