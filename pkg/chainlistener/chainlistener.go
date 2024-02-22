package chainlistener

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ChainListener struct holds the configuration and state for the Ethereum chain listener.
type ChainListener struct {
	rpcURL            string
	contractAddress   common.Address
	eventHash         common.Hash
	startHeight       *big.Int
	eventCallbackFunc func(types.Log)
}

// NewChainListener creates a new ChainListener instance.
func NewChainListener(rpcURL string, contractAddress string) *ChainListener {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal("Failed to connect to the Ethereum client:", err)
	}

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(header.Number)

	fmt.Println(client.Client())

	return &ChainListener{
		rpcURL:          rpcURL,
		contractAddress: common.HexToAddress(contractAddress),
		startHeight:     header.Number,
	}
}

// SubscribeEvent subscribes to the specified event.
func (cl *ChainListener) SubscribeEvent(eventHash string, callbackFunc func(types.Log)) {
	cl.eventHash = common.HexToHash(eventHash)
	cl.eventCallbackFunc = callbackFunc
}

// Start starts the Ethereum chain listener.
func (cl *ChainListener) Start() {
	client, err := ethclient.Dial(cl.rpcURL)
	if err != nil {
		log.Fatal("Failed to connect to the Ethereum client:", err)
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{cl.contractAddress},
		FromBlock: cl.startHeight,
	}

	logsCh := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logsCh)
	if err != nil {
		log.Fatal("Failed to subscribe to logs:", err)
	}

	log.Println("Ethereum chain listener started.")

	for {
		select {
		case err := <-sub.Err():
			log.Fatal("Subscription error:", err)
		case log := <-logsCh:
			if log.Topics[0] == cl.eventHash {
				go cl.handleEvent(log)
			}
		}
	}
}

// handleEvent handles the incoming Ethereum event.
func (cl *ChainListener) handleEvent(log types.Log) {
	cl.eventCallbackFunc(log)
}
