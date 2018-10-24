package testcase

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var bins = []string{
	// testcase\contract_test\solidity\simple_storage.sol
	"0x6060604052341561000f57600080fd5b6040516020806100f383398101604052808051505060056000555060bb806100386000396000f30060606040526004361060485763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166360fe47b18114604d5780636d4ce63c146062575b600080fd5b3415605757600080fd5b60606004356084565b005b3415606c57600080fd5b60726089565b60405190815260200160405180910390f35b600055565b600054905600a165627a7a7230582077f07b169365d1e17ab595dd1a1577bc1871dba4c8657151e8724876055cf0000029",
}

func Test_DeployContract(t *testing.T) {
	for _, bin := range bins {
		deployContracts(t, CmdClient, bin)
		deployContracts(t, CmdLight, bin)
	}
}

func deployContracts(t *testing.T, command, contract string) {
	txHash, _, err1 := SendTx(t, command, 0, 0, 0, KeyFileShard1_1, "", contract, ServerAddr)
	assert.NoError(t, err1)
	if txHash == "" {
		t.Fatal("tx hash is empty")
	}

	timeoutC := time.After(120 * time.Second)
	for {
		receipt, err2 := GetReceipt(t, command, txHash, ServerAddr)
		if receipt != nil && receipt.Failed == false {
			break
		}

		select {
		case <-timeoutC:
			t.Fatalf("over time. err: %s", err2)
		default:
			time.Sleep(10 * time.Second)
		}
	}
}
