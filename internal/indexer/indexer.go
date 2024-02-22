package indexer

import (
	"github.com/taikoxyz/blob-storage/internal/config"
	"github.com/taikoxyz/blob-storage/pkg/chainlistener"
)

// Indexer struct holds the indexer configuration and state.
type Indexer struct {
	cfg *config.Config
	// Add other fields as needed
}

// NewIndexer creates a new Indexer instance.
func NewIndexer(cfg *config.Config) *Indexer {
	return &Indexer{
		cfg: cfg,
	}
}

// Run starts the indexer and listens for events.
func (idx *Indexer) Run() error {
	// Initialize listeners for each network
	for _, network := range idx.cfg.Networks {
		for _, event := range network.IndexedEvents {
			listener := chainlistener.NewChainListener(network.RPCURL, event.Contract)

			listener.SubscribeEvent(event.EventHash, event.Callback)

			// Run the event listener for each network
			go listener.Start()
		}

	}
	select {}
}
