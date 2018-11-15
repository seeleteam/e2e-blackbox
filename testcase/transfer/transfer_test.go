package testcase

import (
	"testing"
	"time"

	"github.com/seeleteam/e2e-blackbox/testcase/common"
)

const (
	_account1 = "0x7c00f5a4312a6a3e458a07c2d650ce13c76b68b1"
)

func Test_Transfer(t *testing.T) {
	amount := 1234
	// Client test
	transferTo(t, common.CmdClient, _account1, amount)

	// Light Test
	transferTo(t, common.CmdLight, _account1, amount)
}

func transferTo(t *testing.T, command, to string, amount int) {
	nonce1, err := common.GetNonce(t, command, common.AccountShard1_1, common.ServerAddr)
	if err != nil {
		t.Fatalf("transferTo GetNonce err=%s", err)
	}

	balance1, err := common.GetBalance(t, command, to, common.ServerAddr)
	if err != nil {
		t.Fatalf("transferTo GetBalance err=%s", err)
	}

	txHash, _, err := common.SendTx(t, command, amount, nonce1, 0, common.KeyFileShard1_1, to, "", common.ServerAddr)
	if err != nil {
		t.Fatalf("transferTo SendTx err=%s", err)
	}

	if txHash == "" {
		t.Fatal("tx hash is empty")
	}

	timeoutC := time.After(120 * time.Second)
	for {
		receipt, err := common.GetReceipt(t, command, txHash, common.ServerAddr)
		if receipt != nil && receipt.Failed == false {
			break
		}

		select {
		case <-timeoutC:
			t.Fatalf("over time. err: %s", err)
		default:
			time.Sleep(10 * time.Second)
		}
	}

	nonce2, err := common.GetNonce(t, command, common.AccountShard1_1, common.ServerAddr)
	if err != nil {
		t.Fatalf("transferTo GetNonce err=%s", err)
	}
	if nonce1+1 != nonce2 {
		t.Fatalf("transferTo nonce err %d ---> %d", nonce1+1, nonce2)
	}

	balance2, err := common.GetBalance(t, command, to, common.ServerAddr)
	if err != nil {
		t.Fatalf("transferTo GetBalance err=%s", err)
	}

	if balance1+int64(amount) != balance2 {
		t.Fatalf("transferTo GetBalance amount err %d ---> %d", balance1+int64(amount), balance2)
	}

}
