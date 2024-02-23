// callbacks.go

package config

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/taikoxyz/blob-storage/internal/taikol1"
)

type Response struct {
	Data []struct {
		Index            string `json:"index"`
		KzgCommitment    string `json:"kzg_commitment"`
		KzgCommitmentHex []byte `json:"-"`
	} `json:"data"`
}

// Callback functions

// BlockProposedCallback is a callback for the "BlockProposed" event.
func BlockProposedCallback(log types.Log) {

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

		storeBlob(strconv.Itoa(int(l1BlobHeight)), blobHash)
	}
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

func storeBlob(blockID, blobHashInMeta string) error {

	url := fmt.Sprintf("https://l1beacon.internal.taiko.xyz/eth/v1/beacon/blob_sidecars/%s", blockID)
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

	blobFound := false
	for _, data := range responseData.Data {

		data.KzgCommitmentHex, err = hex.DecodeString(data.KzgCommitment[2:])
		if err != nil {
			return err
		}

		// Comparing the hex strings of meta.blobHash (blobHash)
		if calculateBlobHash(data.KzgCommitment) == blobHashInMeta {
			blobFound = true
			fmt.Println("BLOB found")
		}
	}
	if !blobFound {
		return errors.New("BLOB not found")
	}
	return nil
}

// Add more callback functions as needed
