package testcase

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	_account1 = "0x7c00f5a4312a6a3e458a07c2d650ce13c76b68b1"
)

func Test_Transfer(t *testing.T) {
	amount := 1234
	// Client test
	transferTo(t, CmdClient, _account1, amount)

	// Light Test
	transferTo(t, CmdLight, _account1, amount)
}

func transferTo(t *testing.T, command, to string, amount int) {
	nonce1, err1 := getNonce(t, command, AccountShard1_1, ServerAddr)
	assert.NoError(t, err1)

	balance1, err2 := getBalance(t, command, to, ServerAddr)
	assert.NoError(t, err2)

	txHash, _, err3 := SendTx(t, command, amount, nonce1, 0, KeyFileShard1_1, to, "", ServerAddr)
	assert.NoError(t, err3)
	if txHash == "" {
		t.Fatal("tx hash is empty")
	}

	timeoutC := time.After(120 * time.Second)
	for {
		receipt, err4 := GetReceipt(t, command, txHash, ServerAddr)
		if receipt != nil && receipt.Failed == false {
			break
		}

		select {
		case <-timeoutC:
			t.Fatalf("over time. err: %s", err4)
		default:
			time.Sleep(10 * time.Second)
		}
	}

	nonce2, err5 := getNonce(t, command, AccountShard1_1, ServerAddr)
	assert.NoError(t, err5)
	assert.Equal(t, nonce1+1, nonce2)

	balance2, err6 := getBalance(t, command, to, ServerAddr)
	assert.NoError(t, err6)
	assert.Equal(t, balance1+int64(amount), balance2)
}
