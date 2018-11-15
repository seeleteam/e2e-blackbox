/**
*  @file
*  @copyright defined in go-seele/LICENSE
 */

package htlc

import (
	"bytes"
	"encoding/json"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/seeleteam/e2e-blackbox/testcase/common"
)

// gas is too low
func Test_HTLC_Create_Low_Gas(t *testing.T) {
	locktime := common.GenerateTime(5)
	cmd := exec.Command(common.CmdClient, "htlc", "create", "--from", common.KeyFileShard1_1, "--to", common.AccountShard1_2, "--amount", "1234", "--price", "15",
		"--gas", "100", "--hash", common.Secretehash, "--time", strconv.FormatInt(locktime, 10))

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Test_HTLC_Create_Low_Gas err: %s", err)
	}

	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Create_Low_Gas: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "create transaction err intrinsic gas too low") {
		t.Fatalf("Test_HTLC_Create_Low_Gas Err: %s", errStr)
	}
}

// available gas
func Test_HTLC_Create_Available_Gas(t *testing.T) {
	amount := int64(1234)
	maxGas := int64(200000)
	beginBalance, err := common.GetBalance(t, common.CmdClient, common.AccountShard1_1, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Create_Available_Gas get balance err: %s", err)
	}

	locktime := common.GenerateTime(5)
	cmd := exec.Command(common.CmdClient, "htlc", "create", "--from", common.KeyFileShard1_1, "--to", common.AccountShard1_2, "--amount", strconv.FormatInt(amount, 10), "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", common.Secretehash, "--time", strconv.FormatInt(locktime, 10))

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Test_HTLC_Create_Available_Gas err: %s", err)
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

	var createInfo common.HTLCCreateInfo
	if err := json.Unmarshal([]byte(str), &createInfo); err != nil {
		t.Fatalf("Test_HTLC_Create_Available_Gas unmarshal created htlc tx err: %s", err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_HTLC_Create_Available_Gas get pool count err: %s", err)
		}

		if number == 0 {
			break
		}
	}

	time.Sleep(20)

	receipt, err := common.GetReceipt(t, common.CmdClient, createInfo.Tx.Hash, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Create_Available_Gas get receipt err: %s", err)
	}

	if receipt.Failed {
		t.Fatalf("Test_HTLC_Create_Available_Gas tx operation fault")
	}

	currentBalance, err := common.GetBalance(t, common.CmdClient, common.AccountShard1_1, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Create_Available_Gas get balance err: %s", err)
	}

	if (receipt.TotalFee + amount + currentBalance) != beginBalance {
		t.Fatalf("Test_HTLC_Create_Available_Gas balance is not equal")
	}

	htlcCreateResult, err := common.HTLCDecode(t, common.CmdClient, receipt.Result)
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

	if common.AccountShard1_1 != htlcCreateResult.Tx.TxData.From {
		t.Fatal("Test_HTLC_Create_Available_Gas htlc sender is not equal to what has been set")
	}

	if common.AccountShard1_2 != htlcCreateResult.To {

		t.Fatal("Test_HTLC_Create_Available_Gas htlc receiver is not equal to what has been set")
	}

	if common.Secretehash != htlcCreateResult.HashLock {
		t.Fatal("Test_HTLC_Create_Available_Gas htlc secrete hash is not equal to what has been set")
	}

	if locktime != htlcCreateResult.TimeLock {
		t.Fatal("Test_HTLC_Create_Available_Gas htlc locked time is not equal to what has been set")
	}

}

func Test_HTLC_Create_Invalid_Time(t *testing.T) {
	amount := int64(1234)
	maxGas := int64(200000)
	beginBalance, err := common.GetBalance(t, common.CmdClient, common.AccountShard1_1, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Create_Invalid_Time get balance err: %s", err)
	}

	locktime := time.Now().Unix()
	cmd := exec.Command(common.CmdClient, "htlc", "create", "--from", common.KeyFileShard1_1, "--to", common.AccountShard1_2, "--amount", strconv.FormatInt(amount, 10), "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", common.Secretehash, "--time", strconv.FormatInt(locktime, 10))

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Test_HTLC_Create_Invalid_Time err: %s", err)
	}

	defer stdin.Close()

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Create_Invalid_Time: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")

	cmd.Wait()

	output, errStr := out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("Test_HTLC_Create_Invalid_Time cmd err: %s", errStr)
	}

	str := output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]

	var createInfo common.HTLCCreateInfo
	if err := json.Unmarshal([]byte(str), &createInfo); err != nil {
		t.Fatalf("Test_HTLC_Create_Invalid_Time unmarshal created htlc tx err: %s", err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_HTLC_Create_Invalid_Time get pool count err: %s", err)
		}

		if number == 0 {
			break
		}
	}

	time.Sleep(10)

	receipt, err := common.GetReceipt(t, common.CmdClient, createInfo.Tx.Hash, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Create_Invalid_Time get receipt err: %s", err)
	}

	if !receipt.Failed {
		t.Fatalf("Test_HTLC_Create_Invalid_Time tx operation fault")
	}

	currentBalance, err := common.GetBalance(t, common.CmdClient, common.AccountShard1_1, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Create_Invalid_Time get balance err: %s", err)
	}

	if (receipt.TotalFee + currentBalance) != beginBalance {
		t.Fatalf("Test_HTLC_Create_Invalid_Time balance is not equal")
	}

	if !strings.Contains(receipt.Result, "Failed to lock, time is not in future") {
		t.Fatalf("Test_HTLC_Create_Invalid_Time Err: %s", receipt.Result)
	}
}

func Test_HTLC_Create_Low_Balance(t *testing.T) {

	beginBalance, err := common.GetBalance(t, common.CmdClient, common.AccountShard1_1, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Create_Low_Balance get balance err: %s", err)
	}

	locktime := common.GenerateTime(5)
	cmd := exec.Command(common.CmdClient, "htlc", "create", "--from", common.KeyFileShard1_1, "--to", common.AccountShard1_2, "--amount", strconv.FormatInt(beginBalance, 10), "--price", "15",
		"--gas", "200000", "--hash", common.Secretehash, "--time", strconv.FormatInt(locktime, 10))

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Test_HTLC_Create_Low_Balance err: %s", err)
	}

	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Create_Low_Balance: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "balance is not enough") {
		t.Fatalf("Test_HTLC_Create_Low_Balance Err: %s", errStr)
	}
}

func Test_HTLC_Create_Invalid_KeyFile(t *testing.T) {
	locktime := common.GenerateTime(5)
	cmd := exec.Command(common.CmdClient, "htlc", "create", "--from", "common.KeyFileShard1_1", "--to", common.AccountShard1_2, "--amount", "1", "--price", "15",
		"--gas", "200000", "--hash", common.Secretehash, "--time", strconv.FormatInt(locktime, 10))

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Test_HTLC_Create_Invalid_KeyFile err: %s", err)
	}

	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Create_Invalid_KeyFile: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "invalid sender key file") {
		t.Fatalf("Test_HTLC_Create_Invalid_KeyFile Err: %s", errStr)
	}
}

func Test_HTLC_Create_Invalid_To_Without_Prefix_0x(t *testing.T) {
	locktime := common.GenerateTime(5)
	cmd := exec.Command(common.CmdClient, "htlc", "create", "--from", common.KeyFileShard1_1, "--to", "common.AccountShard1_2", "--amount", "1", "--price", "15",
		"--gas", "200000", "--hash", common.Secretehash, "--time", strconv.FormatInt(locktime, 10))

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Test_HTLC_Create_Invalid_To_Without_Prefix_0x err: %s", err)
	}

	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Create_Invalid_To_Without_Prefix_0x: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "hex string without 0x prefix") {
		t.Fatalf("Test_HTLC_Create_Invalid_To_Without_Prefix_0x Err: %s", errStr)
	}
}

func Test_HTLC_Create_Invalid_To_With_Prefix_0x_Odd(t *testing.T) {
	locktime := common.GenerateTime(5)
	cmd := exec.Command(common.CmdClient, "htlc", "create", "--from", common.KeyFileShard1_1, "--to", "0x123", "--amount", "1", "--price", "15",
		"--gas", "200000", "--hash", common.Secretehash, "--time", strconv.FormatInt(locktime, 10))

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Test_HTLC_Create_Invalid_To_With_Prefix_0x_Odd err: %s", err)
	}

	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Create_Invalid_To_With_Prefix_0x_Odd: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "hex string of odd length") {
		t.Fatalf("Test_HTLC_Create_Invalid_To_With_Prefix_0x_Odd Err: %s", errStr)
	}
}

func Test_HTLC_Create_Invalid_To_With_Prefix_0x_Even(t *testing.T) {
	locktime := common.GenerateTime(5)
	cmd := exec.Command(common.CmdClient, "htlc", "create", "--from", common.KeyFileShard1_1, "--to", "0x1234", "--amount", "1", "--price", "15",
		"--gas", "200000", "--hash", common.Secretehash, "--time", strconv.FormatInt(locktime, 10))

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Test_HTLC_Create_Invalid_To_With_Prefix_0x_Even err: %s", err)
	}

	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Create_Invalid_To_With_Prefix_0x_Even: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "invalid address length") {
		t.Fatalf("Test_HTLC_Create_Invalid_To_With_Prefix_0x_Even Err: %s", errStr)
	}
}

func Test_HTLC_Create_Invalid_To_With_Prefix_0x_Long_Than_Address(t *testing.T) {
	locktime := common.GenerateTime(5)
	cmd := exec.Command(common.CmdClient, "htlc", "create", "--from", common.KeyFileShard1_1, "--to", "0x0a57a2714e193b7ac50475ce625f2dcfb483d74101", "--amount", "1", "--price", "15",
		"--gas", "200000", "--hash", common.Secretehash, "--time", strconv.FormatInt(locktime, 10))

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Test_HTLC_Create_Invalid_To_With_Prefix_0x_Long_Than_Address err: %s", err)
	}

	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Create_Invalid_To_With_Prefix_0x_Long_Than_Address: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "invalid address length") {
		t.Fatalf("Test_HTLC_Create_Invalid_To_With_Prefix_0x_Long_Than_Address Err: %s", errStr)
	}
}

func Test_HTLC_Create_Invalid_Amount_Less_Than_Zero(t *testing.T) {
	locktime := common.GenerateTime(5)
	cmd := exec.Command(common.CmdClient, "htlc", "create", "--from", common.KeyFileShard1_1, "--to", common.AccountShard1_2, "--amount", "-1", "--price", "15",
		"--gas", "200000", "--hash", common.Secretehash, "--time", strconv.FormatInt(locktime, 10))
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Test_HTLC_Create_Invalid_Amount_Less_Than_Zero err: %s", err)
	}

	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Create_Invalid_Amount_Less_Than_Zero: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()
	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "amount is negative") {
		t.Fatalf("Test_HTLC_Create_Invalid_Amount_Less_Than_Zero Err: %s", errStr)
	}
}
func Test_HTLC_Create_Invalid_Amount(t *testing.T) {
	locktime := common.GenerateTime(5)
	cmd := exec.Command(common.CmdClient, "htlc", "create", "--from", common.KeyFileShard1_1, "--to", common.AccountShard1_2, "--amount", "0op", "--price", "15",
		"--gas", "200000", "--hash", common.Secretehash, "--time", strconv.FormatInt(locktime, 10))
	stdin, err := cmd.StdinPipe()

	if err != nil {
		t.Fatalf("Test_HTLC_Create_Invalid_Amount err: %s", err)
	}

	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Create_Invalid_Amount: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()
	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "invalid amount value") {
		t.Fatalf("Test_HTLC_Create_Invalid_Amount Err: %s", errStr)
	}
}
func Test_HTLC_Withdraw_Available_Gas(t *testing.T) {
	amount := int64(1234)
	maxGas := int64(200000)
	locktime := common.GenerateTime(5)
	cmd := exec.Command(common.CmdClient, "htlc", "create", "--from", common.KeyFileShard1_1, "--to", common.AccountShard1_2, "--amount", strconv.FormatInt(amount, 10), "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", common.Secretehash, "--time", strconv.FormatInt(locktime, 10))

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Available_Gas err: %s", err)
	}

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Available_Gas: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()
	stdin.Close()

	output, errStr := out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("Test_HTLC_Withdraw_Available_Gas cmd err: %s", errStr)
	}

	str := output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]
	var createInfo common.HTLCCreateInfo

	if err := json.Unmarshal([]byte(str), &createInfo); err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Available_Gas unmarshal created htlc tx err: %s", err)
	}

	beginBalance, err := common.GetBalance(t, common.CmdClient, common.AccountShard1_2, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Available_Gas get balance err: %s", err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_HTLC_Withdraw_Available_Gas get pool count err: %s", err)
		}

		if number == 0 {
			break
		}
	}

	time.Sleep(10)

	cmd = exec.Command(common.CmdClient, "htlc", "withdraw", "--from", common.KeyFileShard1_2, "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", createInfo.Tx.Hash, "--preimage", common.Secret)
	out.Reset()
	outErr.Reset()
	cmd.Stdout, cmd.Stderr = &out, &outErr
	stdin, err = cmd.StdinPipe()

	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Available_Gas err: %s", err)
	}

	defer stdin.Close()

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Available_Gas: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()
	output, errStr = out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("Test_HTLC_Withdraw_Available_Gas cmd err: %s", errStr)
	}

	str = output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]
	var withdrawInfo common.HTLCWithDrawInfo

	if err := json.Unmarshal([]byte(str), &withdrawInfo); err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Available_Gas unmarshal created htlc tx err: %s", err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_HTLC_Withdraw_Available_Gas get pool count err: %s", err)
		}
		if number == 0 {
			break
		}
	}

	time.Sleep(10)
	receipt, err := common.GetReceipt(t, common.CmdClient, withdrawInfo.Tx.Hash, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Available_Gas get receipt err: %s", err)
	}

	if receipt.Failed {
		t.Fatalf("Test_HTLC_Withdraw_Available_Gas tx operation fault")
	}

	currentBalance, err := common.GetBalance(t, common.CmdClient, common.AccountShard1_2, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Available_Gas get balance err: %s", err)
	}

	if (receipt.TotalFee + currentBalance - amount) != beginBalance {
		t.Fatalf("Test_HTLC_Withdraw_Available_Gas balance is not equal")
	}

	htlcWithdrawResult, err := common.HTLCDecode(t, common.CmdClient, receipt.Result)
	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Available_Gas htlc decode err: %s", err)
	}

	if !htlcWithdrawResult.Withdrawed {
		t.Fatal("Test_HTLC_Withdraw_Available_Gas htlc withdrawed")
	}

	if htlcWithdrawResult.Refunded {
		t.Fatal("Test_HTLC_Withdraw_Available_Gas htlc refunded")
	}

	if htlcWithdrawResult.Preimage != common.Secret {
		t.Fatal("Test_HTLC_Withdraw_Available_Gas htlc preimage is not equal")
	}

	if amount != htlcWithdrawResult.Tx.TxData.Amount {
		t.Fatal("Test_HTLC_Withdraw_Available_Gas htlc amount is not equal to what has been set")
	}

	if common.AccountShard1_1 != htlcWithdrawResult.Tx.TxData.From {
		t.Fatal("Test_HTLC_Withdraw_Available_Gas htlc sender is not equal to what has been set")
	}

	if common.AccountShard1_2 != htlcWithdrawResult.To {
		t.Fatal("Test_HTLC_Withdraw_Available_Gas htlc receiver is not equal to what has been set")
	}

	if common.Secretehash != htlcWithdrawResult.HashLock {
		t.Fatal("Test_HTLC_Withdraw_Available_Gas htlc secrete hash is not equal to what has been set")
	}

	if locktime != htlcWithdrawResult.TimeLock {
		t.Fatal("Test_HTLC_Withdraw_Available_Gas htlc locked time is not equal to what has been set")
	}
}

func Test_HTLC_Withdraw_Forged_Receiver(t *testing.T) {
	amount := int64(1234)
	maxGas := int64(200000)
	locktime := common.GenerateTime(5)
	cmd := exec.Command(common.CmdClient, "htlc", "create", "--from", common.KeyFileShard1_1, "--to", common.AccountShard1_2, "--amount", strconv.FormatInt(amount, 10), "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", common.Secretehash, "--time", strconv.FormatInt(locktime, 10))

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Receiver err: %s", err)
	}

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Receiver: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()
	stdin.Close()

	output, errStr := out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Receiver cmd err: %s", errStr)
	}

	str := output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]
	var createInfo common.HTLCCreateInfo

	if err := json.Unmarshal([]byte(str), &createInfo); err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Receiver unmarshal created htlc tx err: %s", err)
	}

	beginBalance, err := common.GetBalance(t, common.CmdClient, common.AccountShard1_3, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Receiver get balance err: %s", err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_HTLC_Withdraw_Forged_Receiver get pool count err: %s", err)
		}

		if number == 0 {
			break
		}
	}

	time.Sleep(10)

	cmd = exec.Command(common.CmdClient, "htlc", "withdraw", "--from", common.KeyFileShard1_3, "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", createInfo.Tx.Hash, "--preimage", common.Secret)
	out.Reset()
	outErr.Reset()
	cmd.Stdout, cmd.Stderr = &out, &outErr
	stdin, err = cmd.StdinPipe()

	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Receiver err: %s", err)
	}

	defer stdin.Close()

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Receiver: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()
	output, errStr = out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Receiver cmd err: %s", errStr)
	}

	str = output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]
	var withdrawInfo common.HTLCWithDrawInfo

	if err := json.Unmarshal([]byte(str), &withdrawInfo); err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Receiver unmarshal created htlc tx err: %s", err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_HTLC_Withdraw_Forged_Receiver get pool count err: %s", err)
		}
		if number == 0 {
			break
		}
	}

	time.Sleep(10)
	receipt, err := common.GetReceipt(t, common.CmdClient, withdrawInfo.Tx.Hash, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Receiver get receipt err: %s", err)
	}

	if !receipt.Failed {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Receiver forged receiver withdrawed")
	}

	currentBalance, err := common.GetBalance(t, common.CmdClient, common.AccountShard1_3, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Receiver get balance err: %s", err)
	}

	if (receipt.TotalFee + currentBalance) != beginBalance {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Receiver balance is not equal")
	}

	if !strings.Contains(receipt.Result, "Failed to withdraw, only receiver is allowed") {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Receiver err: %s", receipt.Result)
	}

}

func Test_HTLC_Withdraw_Forged_Preimage(t *testing.T) {
	amount := int64(1234)
	maxGas := int64(200000)
	locktime := common.GenerateTime(5)
	cmd := exec.Command(common.CmdClient, "htlc", "create", "--from", common.KeyFileShard1_1, "--to", common.AccountShard1_2, "--amount", strconv.FormatInt(amount, 10), "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", common.Secretehash, "--time", strconv.FormatInt(locktime, 10))

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Preimage err: %s", err)
	}

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Preimage: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()
	stdin.Close()

	output, errStr := out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Preimage cmd err: %s", errStr)
	}

	str := output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]
	var createInfo common.HTLCCreateInfo

	if err := json.Unmarshal([]byte(str), &createInfo); err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Preimage unmarshal created htlc tx err: %s", err)
	}

	beginBalance, err := common.GetBalance(t, common.CmdClient, common.AccountShard1_2, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Preimage get balance err: %s", err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_HTLC_Withdraw_Forged_Preimage get pool count err: %s", err)
		}

		if number == 0 {
			break
		}
	}

	time.Sleep(10)

	cmd = exec.Command(common.CmdClient, "htlc", "withdraw", "--from", common.KeyFileShard1_2, "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", createInfo.Tx.Hash, "--preimage", common.ForgedSecret)
	out.Reset()
	outErr.Reset()
	cmd.Stdout, cmd.Stderr = &out, &outErr
	stdin, err = cmd.StdinPipe()

	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Preimage err: %s", err)
	}

	defer stdin.Close()

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Preimage: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()
	output, errStr = out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Preimage cmd err: %s", errStr)
	}

	str = output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]
	var withdrawInfo common.HTLCWithDrawInfo

	if err := json.Unmarshal([]byte(str), &withdrawInfo); err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Preimage unmarshal created htlc tx err: %s", err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_HTLC_Withdraw_Forged_Preimage get pool count err: %s", err)
		}
		if number == 0 {
			break
		}
	}

	time.Sleep(10)
	receipt, err := common.GetReceipt(t, common.CmdClient, withdrawInfo.Tx.Hash, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Preimage get receipt err: %s", err)
	}

	if !receipt.Failed {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Preimage forged preimage withdrawed")
	}

	currentBalance, err := common.GetBalance(t, common.CmdClient, common.AccountShard1_2, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Preimage get balance err: %s", err)
	}

	if (receipt.TotalFee + currentBalance) != beginBalance {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Preimage balance is not equal")
	}

	if !strings.Contains(receipt.Result, "Failed to use preimage to match hash") {
		t.Fatalf("Test_HTLC_Withdraw_Forged_Preimage err : %s", receipt.Result)
	}
}

func Test_HTLC_Withdraw_After_Withdrawed(t *testing.T) {
	amount := int64(1234)
	maxGas := int64(200000)
	locktime := common.GenerateTime(5)
	cmd := exec.Command(common.CmdClient, "htlc", "create", "--from", common.KeyFileShard1_1, "--to", common.AccountShard1_2, "--amount", strconv.FormatInt(amount, 10), "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", common.Secretehash, "--time", strconv.FormatInt(locktime, 10))

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_After_Withdrawed err: %s", err)
	}

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Withdraw_After_Withdrawed: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()
	stdin.Close()

	output, errStr := out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("Test_HTLC_Withdraw_After_Withdrawed cmd err: %s", errStr)
	}

	str := output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]
	var createInfo common.HTLCCreateInfo

	if err := json.Unmarshal([]byte(str), &createInfo); err != nil {
		t.Fatalf("Test_HTLC_Withdraw_After_Withdrawed unmarshal created htlc tx err: %s", err)
	}

	beginBalance, err := common.GetBalance(t, common.CmdClient, common.AccountShard1_2, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_After_Withdrawed get balance err: %s", err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_HTLC_Withdraw_After_Withdrawed get pool count err: %s", err)
		}

		if number == 0 {
			break
		}
	}

	time.Sleep(10)

	cmd = exec.Command(common.CmdClient, "htlc", "withdraw", "--from", common.KeyFileShard1_2, "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", createInfo.Tx.Hash, "--preimage", common.Secret)
	out.Reset()
	outErr.Reset()
	cmd.Stdout, cmd.Stderr = &out, &outErr
	stdin, err = cmd.StdinPipe()

	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_After_Withdrawed err: %s", err)
	}

	defer stdin.Close()

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Withdraw_After_Withdrawed: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()
	output, errStr = out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("Test_HTLC_Withdraw_After_Withdrawed cmd err: %s", errStr)
	}

	str = output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]
	var withdrawInfo common.HTLCWithDrawInfo

	if err := json.Unmarshal([]byte(str), &withdrawInfo); err != nil {
		t.Fatalf("Test_HTLC_Withdraw_After_Withdrawed unmarshal created htlc tx err: %s", err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_HTLC_Withdraw_After_Withdrawed get pool count err: %s", err)
		}
		if number == 0 {
			break
		}
	}

	time.Sleep(10)
	receipt, err := common.GetReceipt(t, common.CmdClient, withdrawInfo.Tx.Hash, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_After_Withdrawed get receipt err: %s", err)
	}

	if receipt.Failed {
		t.Fatalf("Test_HTLC_Withdraw_After_Withdrawed tx operation fault")
	}

	currentBalance, err := common.GetBalance(t, common.CmdClient, common.AccountShard1_2, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_After_Withdrawed get balance err: %s", err)
	}

	if (receipt.TotalFee + currentBalance - amount) != beginBalance {
		t.Fatalf("Test_HTLC_Withdraw_After_Withdrawed balance is not equal")
	}

	htlcWithdrawResult, err := common.HTLCDecode(t, common.CmdClient, receipt.Result)
	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_After_Withdrawed htlc decode err: %s", err)
	}

	if !htlcWithdrawResult.Withdrawed {
		t.Fatal("Test_HTLC_Withdraw_After_Withdrawed htlc withdrawed")
	}

	if htlcWithdrawResult.Refunded {
		t.Fatal("Test_HTLC_Withdraw_After_Withdrawed htlc refunded")
	}

	if htlcWithdrawResult.Preimage != common.Secret {
		t.Fatal("Test_HTLC_Withdraw_After_Withdrawed htlc preimage is not equal")
	}

	if amount != htlcWithdrawResult.Tx.TxData.Amount {
		t.Fatal("Test_HTLC_Withdraw_After_Withdrawed htlc amount is not equal to what has been set")
	}

	if common.AccountShard1_1 != htlcWithdrawResult.Tx.TxData.From {
		t.Fatal("Test_HTLC_Withdraw_After_Withdrawed htlc sender is not equal to what has been set")
	}

	if common.AccountShard1_2 != htlcWithdrawResult.To {
		t.Fatal("Test_HTLC_Withdraw_After_Withdrawed htlc receiver is not equal to what has been set")
	}

	if common.Secretehash != htlcWithdrawResult.HashLock {
		t.Fatal("Test_HTLC_Withdraw_After_Withdrawed htlc secrete hash is not equal to what has been set")
	}

	if locktime != htlcWithdrawResult.TimeLock {
		t.Fatal("Test_HTLC_Withdraw_After_Withdrawed htlc locked time is not equal to what has been set")
	}

	cmd = exec.Command(common.CmdClient, "htlc", "withdraw", "--from", common.KeyFileShard1_2, "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", createInfo.Tx.Hash, "--preimage", common.Secret)
	out.Reset()
	outErr.Reset()
	cmd.Stdout, cmd.Stderr = &out, &outErr
	stdin, err = cmd.StdinPipe()

	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_After_Withdrawed err: %s", err)
	}

	defer stdin.Close()

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Withdraw_After_Withdrawed: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()
	output, errStr = out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("Test_HTLC_Withdraw_After_Withdrawed cmd err: %s", errStr)
	}

	str = output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]

	if err := json.Unmarshal([]byte(str), &withdrawInfo); err != nil {
		t.Fatalf("Test_HTLC_Withdraw_After_Withdrawed unmarshal created htlc tx err: %s", err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_HTLC_Withdraw_After_Withdrawed get pool count err: %s", err)
		}
		if number == 0 {
			break
		}
	}

	time.Sleep(10)
	receipt, err = common.GetReceipt(t, common.CmdClient, withdrawInfo.Tx.Hash, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_After_Withdrawed get receipt err: %s", err)
	}

	if !receipt.Failed {
		t.Fatalf("Test_HTLC_Withdraw_After_Withdrawed withdraw again success")
	}

	if !strings.Contains(receipt.Result, "Failed to withdraw, receiver have withdrawed") {
		t.Fatalf("Test_HTLC_Withdraw_After_Withdrawed err : %s", receipt.Result)
	}
}

func Test_HTLC_Withdraw_After_TimeLock(t *testing.T) {
	amount := int64(1234)
	maxGas := int64(200000)
	locktime := time.Now().Unix() + 60
	cmd := exec.Command(common.CmdClient, "htlc", "create", "--from", common.KeyFileShard1_1, "--to", common.AccountShard1_2, "--amount", strconv.FormatInt(amount, 10), "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", common.Secretehash, "--time", strconv.FormatInt(locktime, 10))

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_After_TimeLock err: %s", err)
	}

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Withdraw_After_TimeLock: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()
	stdin.Close()

	output, errStr := out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("Test_HTLC_Withdraw_After_TimeLock cmd err: %s", errStr)
	}

	str := output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]
	var createInfo common.HTLCCreateInfo

	if err := json.Unmarshal([]byte(str), &createInfo); err != nil {
		t.Fatalf("Test_HTLC_Withdraw_After_TimeLock unmarshal created htlc tx err: %s", err)
	}

	timer := time.After(60 * time.Second)
	<-timer

	cmd = exec.Command(common.CmdClient, "htlc", "withdraw", "--from", common.KeyFileShard1_2, "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", createInfo.Tx.Hash, "--preimage", common.Secret)
	out.Reset()
	outErr.Reset()
	cmd.Stdout, cmd.Stderr = &out, &outErr
	stdin, err = cmd.StdinPipe()

	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_After_TimeLock err: %s", err)
	}

	defer stdin.Close()

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Withdraw_After_TimeLock: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()
	output, errStr = out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("Test_HTLC_Withdraw_After_TimeLock cmd err: %s", errStr)
	}

	str = output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]
	var withdrawInfo common.HTLCWithDrawInfo

	if err := json.Unmarshal([]byte(str), &withdrawInfo); err != nil {
		t.Fatalf("Test_HTLC_Withdraw_After_TimeLock unmarshal created htlc tx err: %s", err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_HTLC_Withdraw_After_TimeLock get pool count err: %s", err)
		}
		if number == 0 {
			break
		}
	}

	time.Sleep(10)
	receipt, err := common.GetReceipt(t, common.CmdClient, withdrawInfo.Tx.Hash, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Withdraw_After_TimeLock get receipt err: %s", err)
	}

	if !receipt.Failed {
		t.Fatalf("Test_HTLC_Withdraw_After_TimeLock tx operation fault")
	}

	if !strings.Contains(receipt.Result, "Failed to withraw, time lock is over") {
		t.Fatalf("Test_HTLC_Withdraw_After_TimeLock err : %s", receipt.Result)
	}
}

func Test_HTLC_Refund(t *testing.T) {
	amount := int64(1234)
	maxGas := int64(200000)
	locktime := time.Now().Unix() + 60
	cmd := exec.Command(common.CmdClient, "htlc", "create", "--from", common.KeyFileShard1_1, "--to", common.AccountShard1_2, "--amount", strconv.FormatInt(amount, 10), "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", common.Secretehash, "--time", strconv.FormatInt(locktime, 10))

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Test_HTLC_Refund err: %s", err)
	}

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Refund: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()
	stdin.Close()

	output, errStr := out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("Test_HTLC_Refund cmd err: %s", errStr)
	}

	str := output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]
	var createInfo common.HTLCCreateInfo

	if err := json.Unmarshal([]byte(str), &createInfo); err != nil {
		t.Fatalf("Test_HTLC_Refund unmarshal created htlc tx err: %s", err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_HTLC_Refund get pool count err: %s", err)
		}
		if number == 0 {
			break
		}
	}

	time.Sleep(10)

	receipt, err := common.GetReceipt(t, common.CmdClient, createInfo.Tx.Hash, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Refund get receipt err: %s", err)
	}

	beginBalance, err := common.GetBalance(t, common.CmdClient, common.AccountShard1_1, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Refund get balance err: %s", err)
	}

	timer := time.After(60 * time.Second)
	<-timer
	cmd = exec.Command(common.CmdClient, "htlc", "refund", "--from", common.KeyFileShard1_1, "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", createInfo.Tx.Hash)
	out.Reset()
	outErr.Reset()
	cmd.Stdout, cmd.Stderr = &out, &outErr
	stdin, err = cmd.StdinPipe()

	if err != nil {
		t.Fatalf("Test_HTLC_Refund err: %s", err)
	}

	defer stdin.Close()

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Refund: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()
	output, errStr = out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("Test_HTLC_Refund cmd err: %s", errStr)
	}

	str = output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]
	var refundInfo common.HTLCRefundInfo

	if err := json.Unmarshal([]byte(str), &refundInfo); err != nil {
		t.Fatalf("Test_HTLC_Refund unmarshal refund htlc tx err: %s", err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_HTLC_Refund get pool count err: %s", err)
		}
		if number == 0 {
			break
		}
	}

	time.Sleep(10)
	receipt, err = common.GetReceipt(t, common.CmdClient, refundInfo.Tx.Hash, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Refund get receipt err: %s", err)
	}

	if receipt.Failed {
		t.Fatalf("Test_HTLC_Refund tx operation fault")
	}

	currentBalance, err := common.GetBalance(t, common.CmdClient, common.AccountShard1_1, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Refund get balance err: %s", err)
	}

	if (receipt.TotalFee + currentBalance - amount) != beginBalance {
		t.Fatalf("Test_HTLC_Refund balance is not equal")
	}

	htlcRefundResult, err := common.HTLCDecode(t, common.CmdClient, receipt.Result)
	if err != nil {
		t.Fatalf("Test_HTLC_Refund htlc decode err: %s", err)
	}

	if htlcRefundResult.Withdrawed {
		t.Fatal("Test_HTLC_Refund htlc withdrawed")
	}

	if !htlcRefundResult.Refunded {
		t.Fatal("Test_HTLC_Refund htlc refunded")
	}

	if htlcRefundResult.Preimage != "" {
		t.Fatal("Test_HTLC_Refund htlc preimage is not equal")
	}

	if amount != htlcRefundResult.Tx.TxData.Amount {
		t.Fatal("Test_HTLC_Refund htlc amount is not equal to what has been set")
	}

	if common.AccountShard1_1 != htlcRefundResult.Tx.TxData.From {
		t.Fatal("Test_HTLC_Refund htlc sender is not equal to what has been set")
	}

	if common.AccountShard1_2 != htlcRefundResult.To {
		t.Fatal("Test_HTLC_Refund htlc receiver is not equal to what has been set")
	}

	if common.Secretehash != htlcRefundResult.HashLock {
		t.Fatal("Test_HTLC_Refund htlc secrete hash is not equal to what has been set")
	}

	if locktime != htlcRefundResult.TimeLock {
		t.Fatal("Test_HTLC_Refund htlc locked time is not equal to what has been set")
	}
}

func Test_HTLC_Refund_After_Refund(t *testing.T) {
	amount := int64(1234)
	maxGas := int64(200000)
	locktime := time.Now().Unix() + 60
	cmd := exec.Command(common.CmdClient, "htlc", "create", "--from", common.KeyFileShard1_1, "--to", common.AccountShard1_2, "--amount", strconv.FormatInt(amount, 10), "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", common.Secretehash, "--time", strconv.FormatInt(locktime, 10))

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Refund err: %s", err)
	}

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Refund: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()
	stdin.Close()

	output, errStr := out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("Test_HTLC_Refund_After_Refund cmd err: %s", errStr)
	}

	str := output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]
	var createInfo common.HTLCCreateInfo

	if err := json.Unmarshal([]byte(str), &createInfo); err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Refund unmarshal created htlc tx err: %s", err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_HTLC_Refund_After_Refund get pool count err: %s", err)
		}
		if number == 0 {
			break
		}
	}

	time.Sleep(10)

	receipt, err := common.GetReceipt(t, common.CmdClient, createInfo.Tx.Hash, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Refund get receipt err: %s", err)
	}

	beginBalance, err := common.GetBalance(t, common.CmdClient, common.AccountShard1_1, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Refund get balance err: %s", err)
	}

	timer := time.After(60 * time.Second)
	<-timer
	cmd = exec.Command(common.CmdClient, "htlc", "refund", "--from", common.KeyFileShard1_1, "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", createInfo.Tx.Hash)
	out.Reset()
	outErr.Reset()
	cmd.Stdout, cmd.Stderr = &out, &outErr
	stdin, err = cmd.StdinPipe()

	if err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Refund err: %s", err)
	}

	defer stdin.Close()

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Refund: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()
	output, errStr = out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("Test_HTLC_Refund_After_Refund cmd err: %s", errStr)
	}

	str = output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]
	var refundInfo common.HTLCRefundInfo

	if err := json.Unmarshal([]byte(str), &refundInfo); err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Refund unmarshal refund htlc tx err: %s", err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_HTLC_Refund_After_Refund get pool count err: %s", err)
		}
		if number == 0 {
			break
		}
	}

	time.Sleep(10)
	receipt, err = common.GetReceipt(t, common.CmdClient, refundInfo.Tx.Hash, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Refund get receipt err: %s", err)
	}

	if receipt.Failed {
		t.Fatalf("Test_HTLC_Refund_After_Refund tx operation fault")
	}

	currentBalance, err := common.GetBalance(t, common.CmdClient, common.AccountShard1_1, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Refund get balance err: %s", err)
	}

	if (receipt.TotalFee + currentBalance - amount) != beginBalance {
		t.Fatalf("Test_HTLC_Refund_After_Refund balance is not equal")
	}

	htlcRefundResult, err := common.HTLCDecode(t, common.CmdClient, receipt.Result)
	if err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Refund htlc decode err: %s", err)
	}

	if htlcRefundResult.Withdrawed {
		t.Fatal("Test_HTLC_Refund_After_Refund htlc withdrawed")
	}

	if !htlcRefundResult.Refunded {
		t.Fatal("Test_HTLC_Refund_After_Refund htlc refunded")
	}

	if htlcRefundResult.Preimage != "" {
		t.Fatal("Test_HTLC_Refund_After_Refund htlc preimage is not equal")
	}

	if amount != htlcRefundResult.Tx.TxData.Amount {
		t.Fatal("Test_HTLC_Refund_After_Refund htlc amount is not equal to what has been set")
	}

	if common.AccountShard1_1 != htlcRefundResult.Tx.TxData.From {
		t.Fatal("Test_HTLC_Refund_After_Refund htlc sender is not equal to what has been set")
	}

	if common.AccountShard1_2 != htlcRefundResult.To {
		t.Fatal("Test_HTLC_Refund_After_Refund htlc receiver is not equal to what has been set")
	}

	if common.Secretehash != htlcRefundResult.HashLock {
		t.Fatal("Test_HTLC_Refund_After_Refund htlc secrete hash is not equal to what has been set")
	}

	if locktime != htlcRefundResult.TimeLock {
		t.Fatal("Test_HTLC_Refund_After_Refund htlc locked time is not equal to what has been set")
	}

	beginBalance, err = common.GetBalance(t, common.CmdClient, common.AccountShard1_1, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Refund get balance err: %s", err)
	}

	cmd = exec.Command(common.CmdClient, "htlc", "refund", "--from", common.KeyFileShard1_1, "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", createInfo.Tx.Hash)
	out.Reset()
	outErr.Reset()
	cmd.Stdout, cmd.Stderr = &out, &outErr
	stdin, err = cmd.StdinPipe()

	if err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Refund err: %s", err)
	}

	defer stdin.Close()

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Refund: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()
	output, errStr = out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("Test_HTLC_Refund_After_Refund cmd err: %s", errStr)
	}

	str = output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]

	if err := json.Unmarshal([]byte(str), &refundInfo); err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Refund unmarshal refund htlc tx err: %s", err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_HTLC_Refund_After_Refund get pool count err: %s", err)
		}
		if number == 0 {
			break
		}
	}

	time.Sleep(10)
	receipt, err = common.GetReceipt(t, common.CmdClient, refundInfo.Tx.Hash, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Refund get receipt err: %s", err)
	}

	if !receipt.Failed {
		t.Fatalf("Test_HTLC_Refund_After_Refund tx operation fault")
	}

	currentBalance, err = common.GetBalance(t, common.CmdClient, common.AccountShard1_1, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Refund get balance err: %s", err)
	}

	if (receipt.TotalFee + currentBalance) != beginBalance {
		t.Fatalf("Test_HTLC_Refund_After_Refund balance is not equal")
	}

	if !strings.Contains(receipt.Result, "Failed to refund, owner have refunded") {
		t.Fatalf("Test_HTLC_Refund_After_Refund Err: %s", receipt.Result)
	}
}

func Test_HTLC_Refund_After_Withdrawed(t *testing.T) {
	amount := int64(1234)
	maxGas := int64(200000)
	locktime := time.Now().Unix() + 60
	cmd := exec.Command(common.CmdClient, "htlc", "create", "--from", common.KeyFileShard1_1, "--to", common.AccountShard1_2, "--amount", strconv.FormatInt(amount, 10), "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", common.Secretehash, "--time", strconv.FormatInt(locktime, 10))

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Withdrawed err: %s", err)
	}

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Withdrawed: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()
	stdin.Close()

	output, errStr := out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("Test_HTLC_Refund_After_Withdrawed cmd err: %s", errStr)
	}

	str := output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]
	var createInfo common.HTLCCreateInfo

	if err := json.Unmarshal([]byte(str), &createInfo); err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Withdrawed unmarshal created htlc tx err: %s", err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_HTLC_Refund_After_Withdrawed get pool count err: %s", err)
		}
		if number == 0 {
			break
		}
	}

	time.Sleep(10)

	cmd = exec.Command(common.CmdClient, "htlc", "withdraw", "--from", common.KeyFileShard1_2, "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", createInfo.Tx.Hash, "--preimage", common.Secret)
	out.Reset()
	outErr.Reset()
	cmd.Stdout, cmd.Stderr = &out, &outErr
	stdin, err = cmd.StdinPipe()

	if err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Withdrawed err: %s", err)
	}

	defer stdin.Close()

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Withdrawed: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()
	output, errStr = out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("Test_HTLC_Refund_After_Withdrawed cmd err: %s", errStr)
	}

	str = output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]
	var withdrawInfo common.HTLCWithDrawInfo

	if err := json.Unmarshal([]byte(str), &withdrawInfo); err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Withdrawed unmarshal created htlc tx err: %s", err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_HTLC_Refund_After_Withdrawed get pool count err: %s", err)
		}
		if number == 0 {
			break
		}
	}

	time.Sleep(10)
	receipt, err := common.GetReceipt(t, common.CmdClient, withdrawInfo.Tx.Hash, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Withdrawed get receipt err: %s", err)
	}

	if receipt.Failed {
		t.Fatalf("Test_HTLC_Refund_After_Withdrawed tx operation fault")
	}

	beginBalance, err := common.GetBalance(t, common.CmdClient, common.AccountShard1_1, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Withdrawed get balance err: %s", err)
	}

	timer := time.After(60 * time.Second)
	<-timer
	cmd = exec.Command(common.CmdClient, "htlc", "refund", "--from", common.KeyFileShard1_1, "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", createInfo.Tx.Hash)
	out.Reset()
	outErr.Reset()
	cmd.Stdout, cmd.Stderr = &out, &outErr
	stdin, err = cmd.StdinPipe()

	if err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Withdrawed err: %s", err)
	}

	defer stdin.Close()

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Withdrawed: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()
	output, errStr = out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("Test_HTLC_Refund_After_Withdrawed cmd err: %s", errStr)
	}

	str = output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]
	var refundInfo common.HTLCRefundInfo

	if err := json.Unmarshal([]byte(str), &refundInfo); err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Withdrawed unmarshal refund htlc tx err: %s", err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_HTLC_Refund_After_Withdrawed get pool count err: %s", err)
		}
		if number == 0 {
			break
		}
	}

	time.Sleep(10)

	beginBalance, err = common.GetBalance(t, common.CmdClient, common.AccountShard1_1, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Withdrawed get balance err: %s", err)
	}

	cmd = exec.Command(common.CmdClient, "htlc", "refund", "--from", common.KeyFileShard1_1, "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", createInfo.Tx.Hash)
	out.Reset()
	outErr.Reset()
	cmd.Stdout, cmd.Stderr = &out, &outErr
	stdin, err = cmd.StdinPipe()

	if err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Withdrawed err: %s", err)
	}

	defer stdin.Close()

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Withdrawed: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()
	output, errStr = out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("Test_HTLC_Refund_After_Withdrawed cmd err: %s", errStr)
	}

	str = output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]

	if err := json.Unmarshal([]byte(str), &refundInfo); err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Withdrawed unmarshal refund htlc tx err: %s", err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_HTLC_Refund_After_Withdrawed get pool count err: %s", err)
		}
		if number == 0 {
			break
		}
	}

	time.Sleep(10)
	receipt, err = common.GetReceipt(t, common.CmdClient, refundInfo.Tx.Hash, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Withdrawed get receipt err: %s", err)
	}

	if !receipt.Failed {
		t.Fatalf("Test_HTLC_Refund_After_Withdrawed tx operation fault")
	}

	currentBalance, err := common.GetBalance(t, common.CmdClient, common.AccountShard1_1, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Refund_After_Withdrawed get balance err: %s", err)
	}

	if (receipt.TotalFee + currentBalance) != beginBalance {
		t.Fatalf("Test_HTLC_Refund_After_Withdrawed balance is not equal")
	}

	if !strings.Contains(receipt.Result, "Failed to refund, receiver have withdrawed") {
		t.Fatalf("Test_HTLC_Refund_After_Withdrawed Err: %s", receipt.Result)
	}
}

func Test_HTLC_Refund_Forged_Sender(t *testing.T) {
	amount := int64(1234)
	maxGas := int64(200000)
	locktime := time.Now().Unix() + 60
	cmd := exec.Command(common.CmdClient, "htlc", "create", "--from", common.KeyFileShard1_1, "--to", common.AccountShard1_2, "--amount", strconv.FormatInt(amount, 10), "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", common.Secretehash, "--time", strconv.FormatInt(locktime, 10))

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Test_HTLC_Refund_Forged_Sender err: %s", err)
	}

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Refund_Forged_Sender: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()
	stdin.Close()

	output, errStr := out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("Test_HTLC_Refund_Forged_Sender cmd err: %s", errStr)
	}

	str := output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]
	var createInfo common.HTLCCreateInfo

	if err := json.Unmarshal([]byte(str), &createInfo); err != nil {
		t.Fatalf("Test_HTLC_Refund_Forged_Sender unmarshal created htlc tx err: %s", err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_HTLC_Refund_Forged_Sender get pool count err: %s", err)
		}
		if number == 0 {
			break
		}
	}

	time.Sleep(10)

	receipt, err := common.GetReceipt(t, common.CmdClient, createInfo.Tx.Hash, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Refund_Forged_Sender get receipt err: %s", err)
	}

	timer := time.After(60 * time.Second)
	<-timer
	cmd = exec.Command(common.CmdClient, "htlc", "refund", "--from", common.KeyFileShard1_3, "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", createInfo.Tx.Hash)
	out.Reset()
	outErr.Reset()
	cmd.Stdout, cmd.Stderr = &out, &outErr
	stdin, err = cmd.StdinPipe()

	if err != nil {
		t.Fatalf("Test_HTLC_Refund_Forged_Sender err: %s", err)
	}

	defer stdin.Close()

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Refund_Forged_Sender: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()
	output, errStr = out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("Test_HTLC_Refund_Forged_Sender cmd err: %s", errStr)
	}

	str = output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]
	var refundInfo common.HTLCRefundInfo

	if err := json.Unmarshal([]byte(str), &refundInfo); err != nil {
		t.Fatalf("Test_HTLC_Refund_Forged_Sender unmarshal refund htlc tx err: %s", err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_HTLC_Refund_Forged_Sender get pool count err: %s", err)
		}
		if number == 0 {
			break
		}
	}

	time.Sleep(10)
	receipt, err = common.GetReceipt(t, common.CmdClient, refundInfo.Tx.Hash, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Refund_Forged_Sender get receipt err: %s", err)
	}

	if !receipt.Failed {
		t.Fatalf("Test_HTLC_Refund_Forged_Sender tx operation fault")
	}

	if !strings.Contains(receipt.Result, "Failed to refund, only owner is allowed") {
		t.Fatalf("Test_HTLC_Refund_Forged_Sender Err: %s", receipt.Result)
	}
}

func Test_HTLC_Refund_Forged_Hash(t *testing.T) {
	amount := int64(1234)
	maxGas := int64(200000)
	locktime := time.Now().Unix() + 60
	cmd := exec.Command(common.CmdClient, "htlc", "create", "--from", common.KeyFileShard1_1, "--to", common.AccountShard1_2, "--amount", strconv.FormatInt(amount, 10), "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", common.Secretehash, "--time", strconv.FormatInt(locktime, 10))

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Test_HTLC_Refund_Forged_Hash err: %s", err)
	}

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Refund_Forged_Hash: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()
	stdin.Close()

	output, errStr := out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("Test_HTLC_Refund_Forged_Hash cmd err: %s", errStr)
	}

	str := output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]
	var createInfo common.HTLCCreateInfo

	if err := json.Unmarshal([]byte(str), &createInfo); err != nil {
		t.Fatalf("Test_HTLC_Refund_Forged_Hash unmarshal created htlc tx err: %s", err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_HTLC_Refund_Forged_Hash get pool count err: %s", err)
		}
		if number == 0 {
			break
		}
	}

	time.Sleep(10)

	receipt, err := common.GetReceipt(t, common.CmdClient, createInfo.Tx.Hash, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Refund_Forged_Hash get receipt err: %s", err)
	}

	timer := time.After(60 * time.Second)
	<-timer
	cmd = exec.Command(common.CmdClient, "htlc", "refund", "--from", common.KeyFileShard1_1, "--price", "15",
		"--gas", strconv.FormatInt(maxGas, 10), "--hash", "0x1234567890")
	out.Reset()
	outErr.Reset()
	cmd.Stdout, cmd.Stderr = &out, &outErr
	stdin, err = cmd.StdinPipe()

	if err != nil {
		t.Fatalf("Test_HTLC_Refund_Forged_Hash err: %s", err)
	}

	defer stdin.Close()

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_HTLC_Refund_Forged_Hash: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()
	output, errStr = out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("Test_HTLC_Refund_Forged_Hash cmd err: %s", errStr)
	}

	str = output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]
	var refundInfo common.HTLCRefundInfo

	if err := json.Unmarshal([]byte(str), &refundInfo); err != nil {
		t.Fatalf("Test_HTLC_Refund_Forged_Hash unmarshal refund htlc tx err: %s", err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_HTLC_Refund_Forged_Hash get pool count err: %s", err)
		}
		if number == 0 {
			break
		}
	}

	time.Sleep(10)
	receipt, err = common.GetReceipt(t, common.CmdClient, refundInfo.Tx.Hash, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Refund_Forged_Hash get receipt err: %s", err)
	}

	if !receipt.Failed {
		t.Fatalf("Test_HTLC_Refund_Forged_Hash tx operation fault")
	}

	if !strings.Contains(receipt.Result, "Failed to get data with key") {
		t.Fatalf("Test_HTLC_Refund_Forged_Hash Err: %s", receipt.Result)
	}
}
