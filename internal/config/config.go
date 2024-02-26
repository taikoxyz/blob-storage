package config

import "github.com/ethereum/go-ethereum/core/types"

// CallbackFunc represents the type of the callback function for handling events.
type CallbackFunc func(string, string, types.Log)

// IndexedEvent struct represents the configuration for an indexed event.
type IndexedEvent struct {
	Contract  string
	EventName string
	EventHash string
	Callback  CallbackFunc
}

// NetworkConfig struct represents the configuration for a network.
type NetworkConfig struct {
	RPCURL        string
	BeaconURL     string // Add this field
	NetworkName   string
	IndexedEvents []IndexedEvent
}

// Config struct holds the overall configuration for the application.
type Config struct {
	Networks []NetworkConfig
}

// GetConfig loads the configuration from environment variables or a config file.
func GetConfig() (*Config, error) {
	cfg := &Config{}
	cfg.Networks = []NetworkConfig{
		{
			RPCURL:      "wss://l1ws.internal.taiko.xyz",
			BeaconURL:   "https://l1beacon.internal.taiko.xyz/eth/v1/beacon/blob_sidecars/", // Set the beacon URL here
			NetworkName: "taiko_internal_l1",
			IndexedEvents: []IndexedEvent{
				{
					Contract:  "0xbE71D121291517c85Ab4d3ac65d70F6b1FD57118",
					EventHash: "0xa62cea5af360b010ef0d23472a2a7493b54175fd9fd2f9c2aa2bb427d2f4d3ca",
					EventName: "BlockProposed",
					Callback:  BlockProposedCallback,
				},
			},
		},
	}

	return cfg, nil
}
