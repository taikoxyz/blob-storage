package main

import (
	"log"

	"github.com/taikoxyz/blob-storage/internal/config"
	"github.com/taikoxyz/blob-storage/internal/indexer"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	idx := indexer.NewIndexer(cfg)
	if err := idx.Run(); err != nil {
		log.Fatal("Error running indexer:", err)
	}
}
