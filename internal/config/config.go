package config

import "github.com/ethereum/go-ethereum/core/types"

// CallbackFunc represents the type of the callback function for handling events.
type CallbackFunc func(types.Log)

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
			RPCURL:      "wss://ethereum-holesky.publicnode.com",
			NetworkName: "holesky",
			IndexedEvents: []IndexedEvent{
				{
					Contract:  "0xB20BB9105e007Bd3E0F73d63D4D3dA2c8f736b77",
					EventHash: "0xa62cea5af360b010ef0d23472a2a7493b54175fd9fd2f9c2aa2bb427d2f4d3ca",
					EventName: "BlockProposed",
					Callback:  BlockProposedCallback,
				},
			},
		},
	}

	return cfg, nil
}
