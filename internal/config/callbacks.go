// callbacks.go

package config

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/taikoxyz/blob-storage/internal/taikol1"
)

type Response struct {
	Data []struct {
		Index            string `json:"index"`
		Blob             string `json:"blob"`
		KzgCommitment    string `json:"kzg_commitment"`
		KzgCommitmentHex []byte `json:"-"`
	} `json:"data"`
}

// Callback functions

// BlockProposedCallback is a callback for the "BlockProposed" event.
func BlockProposedCallback(rpcURL, beaconURL string, log types.Log) {

	contractAbi, err := abi.JSON(strings.NewReader(taikol1.TaikoL1ABI))
	if err != nil {
		fmt.Println("Could not initiate reader")
	}

	eventData := taikol1.TaikoL1BlockProposed{}

	err = contractAbi.UnpackIntoInterface(&eventData, "BlockProposed", log.Data)
	if err != nil {
		fmt.Println("Could not unpack log.Data")
	}

	// Some debug logs for now.
	// fmt.Println("Block id is:")
	// fmt.Println(log.Topics[1])

	if eventData.Meta.BlobUsed {
		// in LibPropose we assign block.height-1 to l1Height, which is the parent block.
		l1BlobHeight := eventData.Meta.L1Height + 1
		blobHash := hex.EncodeToString(eventData.Meta.BlobHash[:])

		storeBlob(rpcURL, beaconURL, strconv.Itoa(int(l1BlobHeight)), blobHash)
	}
}

func getBlockTimestamp(rpcURL string, blockNumber *big.Int) (uint64, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return 0, err
	}
	defer client.Close()

	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		return 0, err
	}

	return block.Time(), nil
}

func calculateBlobHash(commitmentStr string) string {
	// As per: https://eips.ethereum.org/EIPS/eip-4844
	commitment := kzg4844.Commitment(common.FromHex(commitmentStr))
	blobHash := kzg4844.CalcBlobHashV1(
		sha256.New(),
		&commitment)
	blobHashString := hex.EncodeToString(blobHash[:])
	return blobHashString
}

func storeBlob(rpcURL, beaconURL, blockID, blobHashInMeta string) error {

	// url := fmt.Sprintf("https://l1beacon.internal.taiko.xyz/eth/v1/beacon/blob_sidecars/%s", blockID)
	url := fmt.Sprintf("%s/%s", beaconURL, blockID)
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var responseData Response
	if err := json.Unmarshal(body, &responseData); err != nil {
		return err
	}

	for _, data := range responseData.Data {

		data.KzgCommitmentHex, err = hex.DecodeString(data.KzgCommitment[2:])
		if err != nil {
			return err
		}

		// Comparing the hex strings of meta.blobHash (blobHash)
		if calculateBlobHash(data.KzgCommitment) == blobHashInMeta {
			fmt.Println("BLOB found")
			n := new(big.Int)

			blockNrBig, ok := n.SetString(blockID, 10)
			if !ok {
				fmt.Println("SetString: error")
				return errors.New("SetString: error")
			}
			blockTs, err := getBlockTimestamp(rpcURL, blockNrBig)

			if err != nil {
				fmt.Println("TIMESTAMP issue")
				return errors.New("TIMESTAMP issue")
			}
			fmt.Println("The blobHash:", blobHashInMeta)
			fmt.Println("The block:", blockNrBig)
			fmt.Println("The kzg commitment:", data.KzgCommitment)
			fmt.Println("The corresponding timestamp:", blockTs)
			fmt.Println("The blob:", data.Blob[0:100])

			return nil
		}
	}

	return errors.New("BLOB not found")
}

// Add more functions as needed
