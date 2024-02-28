package logic

import (
	"github.com/ethereum/go-ethereum/core/types"
)

// CallbackFunc represents the type of the callback function for handling events.
type CallbackFunc func(string, string, string, types.Log)

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

// MongoDBConfig holds the configuration for MongoDB.
type MongoDBConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

// Config struct holds the overall configuration for the application.
type Config struct {
	Networks []NetworkConfig
	MongoDB  MongoDBConfig
}

// GetConfig loads the configuration from environment variables or a config file.
func GetConfig() (*Config, error) {
	cfg := &Config{}
	cfg.Networks = []NetworkConfig{
		{
			RPCURL:      "wss://l1ws.internal.taiko.xyz",
			BeaconURL:   "https://l1beacon.internal.taiko.xyz/eth/v1/beacon/blob_sidecars/", // Set the beacon URL here
			NetworkName: "L2A",
			IndexedEvents: []IndexedEvent{
				{
					Contract:  "0x1FD3Df7E9C15390c8589D2E4d43757eA692ae256",
					EventHash: "0xa62cea5af360b010ef0d23472a2a7493b54175fd9fd2f9c2aa2bb427d2f4d3ca",
					EventName: "BlockProposed",
					Callback:  BlockProposedCallback,
				},
			},
		},
		{
			RPCURL:      "wss://l1ws.internal.taiko.xyz",
			BeaconURL:   "https://l1beacon.internal.taiko.xyz/eth/v1/beacon/blob_sidecars/", // Set the beacon URL here
			NetworkName: "L2B",
			IndexedEvents: []IndexedEvent{
				{
					Contract:  "0x1670010000000000000000000000000000010001",
					EventHash: "0xa62cea5af360b010ef0d23472a2a7493b54175fd9fd2f9c2aa2bb427d2f4d3ca",
					EventName: "BlockProposed",
					Callback:  BlockProposedCallback,
				},
			},
		},
	}

	// Same DB for the blobs (L2A, L2B, etc.)
	cfg.MongoDB = MongoDBConfig{
		Host:     "localhost",
		Port:     27017,
		Username: "",             // Add your MongoDB username if needed
		Password: "",             // Add your MongoDB password if needed
		Database: "blob_storage", // Choose your MongoDB database name
	}

	return cfg, nil
}
