// callbacks.go

package config

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/taikoxyz/blob-storage/internal/taikol1"
)

// Callback functions

// BlockProposedCallback is a callback for the "BlockProposed" event.
func BlockProposedCallback(log types.Log) {

	contractAbi, err := abi.JSON(strings.NewReader(taikol1.TaikoL1ABI))
	if err != nil {
		fmt.Println("No bueno")
	}

	eventData := taikol1.TaikoL1BlockProposed{}

	// Some debug logs for now.
	fmt.Println("Block id is:")
	fmt.Println(log.Topics[1])
	fmt.Println("Assigned prover is:")
	fmt.Println(log.Topics[2])

	err = contractAbi.UnpackIntoInterface(&eventData, "BlockProposed", log.Data)
	if err != nil {
		fmt.Println("No bueno")
	}

	fmt.Println("blobUsed:")
	fmt.Println(eventData.Meta.BlobUsed)
	fmt.Println("blobUsed:")
	fmt.Println(eventData.Meta.BlobHash)
	// ToDo: We need to determine some things
	// - mock a blobHash now and retrieve as written in document and compare
}

// Add more callback functions as needed
