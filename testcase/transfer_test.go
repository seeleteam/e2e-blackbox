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
	nonce1, err1 := getNonce(t, CmdClient, AccountShard1_1, ServerAddr)
	assert.NoError(t, err1)

	balance1, err2 := getBalance(t, CmdClient, _account1, ServerAddr)
	assert.NoError(t, err2)

	txHash, _, err3 := sentTX(t, CmdClient, amount, nonce1, KeyFileShard1_1, _account1, ServerAddr)
	assert.NoError(t, err3)
	if txHash == "" {
		t.Fatal("tx hash is empty")
	}

	timeoutC := time.After(120 * time.Second)
	for {
		receipt, err4 := getReceipt(t, CmdClient, txHash, ServerAddr)
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

	nonce2, err5 := getNonce(t, CmdClient, AccountShard1_1, ServerAddr)
	assert.NoError(t, err5)
	assert.Equal(t, nonce1+1, nonce2)

	balance2, err6 := getBalance(t, CmdClient, _account1, ServerAddr)
	assert.NoError(t, err6)
	assert.Equal(t, balance1+int64(amount), balance2)

	// Light Test
	nonce11, err11 := getNonce(t, CmdLight, AccountShard1_1, ServerAddr)
	assert.NoError(t, err11)
	assert.Equal(t, nonce11, nonce2)

	balance11, err12 := getBalance(t, CmdLight, _account1, ServerAddr)
	assert.NoError(t, err12)
	assert.Equal(t, balance11, balance2)

	txHash1, _, err13 := sentTX(t, CmdLight, amount, nonce2, KeyFileShard1_1, _account1, ServerAddr)
	assert.NoError(t, err13)
	if txHash1 == "" {
		t.Fatal("tx hash is empty")
	}

	timeoutC1 := time.After(120 * time.Second)
	for {
		receipt1, err14 := getReceipt(t, CmdLight, txHash1, ServerAddr)
		if receipt1 != nil && receipt1.Failed == false {
			break
		}

		select {
		case <-timeoutC1:
			t.Fatalf("over time. err: %s", err14)
		default:
			time.Sleep(10 * time.Second)
		}
	}

	nonce12, err15 := getNonce(t, CmdLight, AccountShard1_1, ServerAddr)
	assert.NoError(t, err15)
	assert.Equal(t, nonce11+1, nonce12)

	balance12, err16 := getBalance(t, CmdLight, _account1, ServerAddr)
	assert.NoError(t, err16)
	assert.Equal(t, balance11+int64(amount), balance12)
}
