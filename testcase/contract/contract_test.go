package contract

import (
	"strings"
	"testing"

	"github.com/seeleteam/go-seele/common"
	"github.com/stretchr/testify/assert"
)

var bins = []string{
	// testcase\contract_test\solidity\simple_storage.sol
	"0x6060604052341561000f57600080fd5b6040516020806100f383398101604052808051505060056000555060bb806100386000396000f30060606040526004361060485763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166360fe47b18114604d5780636d4ce63c146062575b600080fd5b3415605757600080fd5b60606004356084565b005b3415606c57600080fd5b60726089565b60405190815260200160405180910390f35b600055565b600054905600a165627a7a7230582077f07b169365d1e17ab595dd1a1577bc1871dba4c8657151e8724876055cf0000029",
}

var abiFiles = []string{
	"./solabi/SimpleStorage.abi",
}

func Test_DeployAndCallContract(t *testing.T) {
	for _, bin := range bins {
		// deploy contract
		receipt1 := handldTx(t, cmdClient, "", bin)
		CallSimpleStorage(t, cmdClient, receipt1.Contract)
		// deploy contract
		receipt2 := handldTx(t, cmdLight, "", bin)
		CallSimpleStorage(t, cmdLight, receipt2.Contract)
	}
}

func CallSimpleStorage(t *testing.T, command, contract string) {
	// filepath
	filepath := "./solabi/SimpleStorage.abi"
	if !common.FileOrFolderExists(filepath) {
		t.Fatal("abi file not found")
	}

	// call get contract
	payload := generatePayload(t, command, filepath, "get")
	payload = payload[strings.IndexAny(payload, "0x"):strings.IndexAny(payload, "\n")]

	callRec1 := handldTx(t, command, contract, payload)
	assert.NotNil(t, callRec1)
	assert.Equal(t, false, callRec1.Failed)
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000005", callRec1.Result)
	// call set contract
	payload = generatePayload(t, command, filepath, "set", "23")
	payload = payload[strings.IndexAny(payload, "0x"):strings.IndexAny(payload, "\n")]

	callRec2 := handldTx(t, command, contract, payload)
	assert.NotNil(t, callRec2)
	assert.Equal(t, false, callRec2.Failed)
	// call get contract
	payload = generatePayload(t, command, filepath, "get")
	payload = payload[strings.IndexAny(payload, "0x"):strings.IndexAny(payload, "\n")]

	callRec3 := handldTx(t, command, contract, payload)
	assert.NotNil(t, callRec3)
	assert.Equal(t, false, callRec3.Failed)
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000017", callRec3.Result)
}
