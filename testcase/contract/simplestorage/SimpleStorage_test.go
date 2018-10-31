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
	abiFile = "./simplestorage.abi"
	binFile = "./SimpleStorage.bin"
)

func Test_DeployAndCallContract_client(t *testing.T) {
	// deploy contract
	receipt := contract.HandleTx(t, 0, contract.CmdClient, contract.KeyFileShard1, "", contract.ParseBinFile(t, binFile))
	callSimpleStorage(t, contract.CmdClient, contract.KeyFileShard1, receipt.Contract)
}

func Test_DeployAndCallContract_light(t *testing.T) {
	// deploy contract
	receipt := contract.HandleTx(t, 0, contract.CmdLight, contract.KeyFileShard2, "", contract.ParseBinFile(t, binFile))
	callSimpleStorage(t, contract.CmdLight, contract.KeyFileShard2, receipt.Contract)
}

func callSimpleStorage(t *testing.T, command, from, contractAddr string) {
	if !common.FileOrFolderExists(abiFile) {
		t.Fatal("abi file not found")
	}

	// call get contract
	payload := contract.GeneratePayload(t, command, abiFile, "get")
	payload = payload[strings.IndexAny(payload, "0x"):strings.IndexAny(payload, "\n")]

	callRec1 := contract.HandleTx(t, 0, command, from, contractAddr, payload)
	assert.NotNil(t, callRec1)
	assert.Equal(t, false, callRec1.Failed)
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000005", callRec1.Result)
	// call set contract
	payload = contract.GeneratePayload(t, command, abiFile, "set", "23")
	payload = payload[strings.IndexAny(payload, "0x"):strings.IndexAny(payload, "\n")]

	callRec2 := contract.HandleTx(t, 0, command, from, contractAddr, payload)
	assert.NotNil(t, callRec2)
	assert.Equal(t, false, callRec2.Failed)
	// call get contract
	payload = contract.GeneratePayload(t, command, abiFile, "get")
	payload = payload[strings.IndexAny(payload, "0x"):strings.IndexAny(payload, "\n")]

	callRec3 := contract.HandleTx(t, 0, command, from, contractAddr, payload)
	assert.NotNil(t, callRec3)
	assert.Equal(t, false, callRec3.Failed)
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000017", callRec3.Result)
}
