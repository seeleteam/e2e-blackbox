/**
*  @file
*  @copyright defined in go-seele/LICENSE
 */

package testcase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

// gas is too low
func Test_HTLC_Create_Low_Gas(t *testing.T) {
	locktime := generateTime(5)
	cmd := exec.Command(CmdClient, "htlc", "create", "--from", KeyFileShard1_1, "--to", AccountShard1_2, "--amount", "1234", "--price", "15",
		"--gas", "100", "--hash", Secretehash, "--time", strconv.FormatInt(locktime, 10))

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}

	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Create: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "create transaction err intrinsic gas too low") {
		t.Fatalf("Test_HTLC_Create_Low_Gas Err:%s", errStr)
	}
}

// available gas
func Test_HTLC_Create_Available_Gas(t *testing.T) {
	amount := int64(1234)
	maxGas := int64(200000)
	beginBalance, err := getBalance(t, CmdClient, AccountShard1_1, ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Create_Available_Gas get balance err: %s", err)
	}

	locktime := generateTime(5)
	cmd := exec.Command(CmdClient, "htlc", "create", "--from", KeyFileShard1_1, "--to", AccountShard1_2, "--amount", strconv.FormatInt(amount, 10), "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", Secretehash, "--time", strconv.FormatInt(locktime, 10))

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}

	defer stdin.Close()

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Create_Available_Gas: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")

	cmd.Wait()

	output, errStr := out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("Test_HTLC_Create_Available_Gas cmd err: %s", errStr)
	}

	str := output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]
	// fmt.Println("str:", str)
	var createInfo HTLCCreateInfo
	if err := json.Unmarshal([]byte(str), &createInfo); err != nil {
		t.Fatalf("Test_HTLC_Create_Available_Gas unmarshal created htlc tx err: %s", err)
	}

	for {
		time.Sleep(10)
		number, err := getPoolCountTxs(t, CmdClient, ServerAddr)
		if err != nil {
			t.Fatalf("Test_HTLC_Create_Available_Gas get pool count err: %s", err)
		}

		if number == 0 {
			break
		}
	}

	time.Sleep(20)

	receipt, err := getReceipt(t, CmdClient, createInfo.Tx.Hash, ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Create_Available_Gas get receipt err: %s", err)
	}

	if receipt.Failed {
		t.Fatalf("Test_HTLC_Create_Available_Gas tx operation fault")
	}

	currentBalance, err := getBalance(t, CmdClient, AccountShard1_1, ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Create_Available_Gas get balance err: %s", err)
	}

	if (receipt.TotalFee + amount + currentBalance) != beginBalance {
		t.Fatalf("Test_HTLC_Create_Available_Gas balance is not equal")
	}

	htlcCreateResult, err := htlcDecode(t, CmdClient, receipt.Result)
	if err != nil {
		t.Fatalf("Test_HTLC_Create_Available_Gas htlc decode err: %s", err)
	}

	if htlcCreateResult.Withdrawed {
		t.Fatal("Test_HTLC_Create_Available_Gas htlc withdrawed")
	}

	if htlcCreateResult.Refunded {
		t.Fatal("Test_HTLC_Create_Available_Gas htlc refunded")
	}

	if htlcCreateResult.Preimage != "" {
		t.Fatal("Test_HTLC_Create_Available_Gas htlc preimage is not empty")
	}

	if amount != htlcCreateResult.Tx.TxData.Amount {
		t.Fatal("Test_HTLC_Create_Available_Gas htlc amount is not equal to what has been set")
	}

	if AccountShard1_1 != htlcCreateResult.Tx.TxData.From {
		t.Fatal("Test_HTLC_Create_Available_Gas htlc sender is not equal to what has been set")
	}

	if AccountShard1_2 != htlcCreateResult.To {
		t.Fatal("Test_HTLC_Create_Available_Gas htlc receiver is not equal to what has been set")
	}

	if Secretehash != htlcCreateResult.HashLock {
		t.Fatal("Test_HTLC_Create_Available_Gas htlc secrete hash is not equal to what has been set")
	}

	if locktime != htlcCreateResult.TimeLock {
		t.Fatal("Test_HTLC_Create_Available_Gas htlc locked time is not equal to what has been set")
	}

}

func Test_HTLC_Create_Low_Balance(t *testing.T) {

	beginBalance, err := getBalance(t, CmdClient, AccountShard1_1, ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Create_Low_Balance get balance err: %s", err)
	}

	locktime := generateTime(5)
	cmd := exec.Command(CmdClient, "htlc", "create", "--from", KeyFileShard1_1, "--to", AccountShard1_2, "--amount", strconv.FormatInt(beginBalance, 10), "--price", "15",
		"--gas", "200000", "--hash", Secretehash, "--time", strconv.FormatInt(locktime, 10))

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}

	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Create: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "balance is not enough") {
		t.Fatalf("Test_HTLC_Create_Low_Balance Err:%s", errStr)
	}
}
