package contract

import (
	"strings"
	"testing"

	"github.com/seeleteam/e2e-blackbox/testcase/contract"
	"github.com/seeleteam/go-seele/common"
	"github.com/stretchr/testify/assert"
)

var (
	// testcase\contract\simplestorage\simplestorage.sol
	abiFile = "./SimpleStorage.abi"
	binFile = "./SimpleStorage.bin"
)

// func Test_DeployAndCallContract_client(t *testing.T) {
// 	// deploy contract
// 	receipt := contract.HandleTx(t, 0, contract.CmdClient, contract.KeyFileShard11, "", contract.ParseBinFile(t, binFile))
// 	// fmt.Println("receipt:", receipt)
// 	callSimpleStorage(t, contract.CmdClient, contract.KeyFileShard11, receipt.Contract)
// }

// func Test_DeployAndCallContract_light(t *testing.T) {
// 	// deploy contract
// 	receipt := contract.HandleTx(t, 0, contract.CmdLight, contract.KeyFileShard12, "", contract.ParseBinFile(t, binFile))
// 	// fmt.Println("receipt:", receipt)
// 	callSimpleStorage(t, contract.CmdLight, contract.KeyFileShard12, receipt.Contract)
// }

func callSimpleStorage(t *testing.T, command, from, contractAddr string) {
	if !common.FileOrFolderExists(abiFile) {
		t.Fatal("abi file not found")
	}

	// call get contract
	payload := contract.GeneratePayload(t, command, abiFile, "get")
	payload = payload[strings.IndexAny(payload, "0x"):strings.IndexAny(payload, "\n")]

	callRec1 := contract.HandleTx(t, 0, command, from, contractAddr, payload)
	// fmt.Println("callRec1:", callRec1)
	assert.NotNil(t, callRec1)
	assert.Equal(t, false, callRec1.Failed)
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000005", callRec1.Result)
	// call set contract
	payload = contract.GeneratePayload(t, command, abiFile, "set", "23")
	payload = payload[strings.IndexAny(payload, "0x"):strings.IndexAny(payload, "\n")]

	callRec2 := contract.HandleTx(t, 0, command, from, contractAddr, payload)
	// fmt.Println("callRec2:", callRec2)
	assert.NotNil(t, callRec2)
	assert.Equal(t, false, callRec2.Failed)
	// call get contract
	payload = contract.GeneratePayload(t, command, abiFile, "get")
	payload = payload[strings.IndexAny(payload, "0x"):strings.IndexAny(payload, "\n")]

	callRec3 := contract.HandleTx(t, 0, command, from, contractAddr, payload)
	// fmt.Println("callRec3:", callRec3)
	assert.NotNil(t, callRec3)
	assert.Equal(t, false, callRec3.Failed)
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000017", callRec3.Result)
}
