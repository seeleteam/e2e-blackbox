package contract

import (
	"os/exec"
	"testing"
	"time"

	"github.com/seeleteam/e2e-blackbox/testcase"
	"github.com/stretchr/testify/assert"
)

const (
	cmdClient = "../../bin/client"
	cmdLight  = "../../bin/light"

	keyFileShard1 = "../../config/keyfile/shard1-0x0a57a2714e193b7ac50475ce625f2dcfb483d741"
)

func handldTx(t *testing.T, command, contract, payload string) (receipt *testcase.ReceiptInfo) {
	txHash, _, err1 := testcase.SendTx(t, command, 0, 0, 0, keyFileShard1, contract, payload, testcase.ServerAddr)
	assert.NoError(t, err1)
	if txHash == "" {
		t.Fatal("tx hash is empty")
	}

	timeoutC := time.After(150 * time.Second)
	for {
		var err2 error
		receipt, err2 = testcase.GetReceipt(t, command, txHash, testcase.ServerAddr)
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

	return receipt
}

func generatePayload(t *testing.T, command, abi, method string, args ...string) (payload string) {
	cmd := exec.Command(command, "payload", "--abi", abi, "--method", method)
	for _, arg := range args {
		cmd.Args = append(cmd.Args, "--args", arg)
	}

	bytes, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	return string(bytes)
}
