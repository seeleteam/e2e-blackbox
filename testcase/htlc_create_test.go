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
	"math/big"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/seeleteam/go-seele/api"
	"github.com/seeleteam/go-seele/core/types"
)

// getBalance get account balance
func getBalance(command, account, address string) (*big.Int, error) {
	cmd := exec.Command(command, "getbalance", "--account", account, "--address", address)

	info, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var balanceInfo api.GetBalanceResponse
	if err := json.Unmarshal(info, &balanceInfo); err != nil {
		return nil, err
	}

	return balanceInfo.Balance, nil
}

// get receipt
func getReceipt(command, hash, address string) (map[string]interface{}, error) {
	cmd := exec.Command(command, "getreceipt", "--hash", hash, "--address", address)

	info, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	if err := json.Unmarshal(info, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// generateTime generate time
func generateTime(minutes int64) string {
	num := time.Now().Unix() + minutes*60
	return strconv.FormatInt(num, 10)
}

// gas is too low
func Test_HTLC_Create_Low_Gas(t *testing.T) {
	cmd := exec.Command(CmdClient, "htlc", "create", "--from", KeyFileShard1_1, "--to", AccountShard1_2, "--amount", "1234", "--price", "15",
		"--gas", "100", "--hash", Secretehash, "--time", generateTime(5))

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
	amount := "1234"
	maxGas := "200000"
	beginBalance, err := getBalance(CmdClient, AccountShard1_1, ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Create_Available_Gas get balance err: %s", err)
	}

	fmt.Println("currentBalance:", beginBalance)

	cmd := exec.Command(CmdClient, "htlc", "create", "--from", KeyFileShard1_1, "--to", AccountShard1_2, "--amount", amount, "--price", "15",
		"--gas", maxGas, "--hash", Secretehash, "--time", generateTime(5))

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
	assert.Equal(t, "", errStr)

	first := strings.Index(output, "{")
	end := strings.LastIndex(output, "}")

	htlcInfo := make(map[string]interface{})
	if err := json.Unmarshal([]byte(output[first:end+1]), &htlcInfo); err != nil {
		t.Fatalf("Test_HTLC_Create_Available_Gas unmarshal created htlc tx err: %s", err)
	}
	fmt.Println("htlcinfo:", htlcInfo)

	tx, ok := htlcInfo["Tx"].(*types.Transaction)
	fmt.Println("tx:", tx)
	if !ok {
		t.Fatalf("Test_HTLC_Create_Available_Gas htlc tx is not types.Transaction")
	}

	time.Sleep(5)
	receipt, err := getReceipt(CmdClient, tx.Hash.ToHex(), ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Create_Available_Gas get receipt err")
	}

	tag := receipt["failed"]
	// result := receipt["result"]
	totalFee := receipt["totalFee"]
	// txHash := receipt["txhash"]
	// usedGas := receipt["usedGas"]
	if v, ok := tag.(bool); ok && v {
		t.Fatalf("Test_HTLC_Create_Available_Gas receipt failed is true")
	}

	currentBalance, err := getBalance(CmdClient, AccountShard1_1, ServerAddr)
	if err != nil {
		t.Fatalf("Test_HTLC_Create_Available_Gas get balance err: %s", err)
	}

	amountNum, err := strconv.Atoi(amount)
	if err != nil {
		t.Fatalf("Test_HTLC_Create_Available_Gas change amount string err: %s", err)
	}

	num := reflect.ValueOf(totalFee).Int() + int64(amountNum)
	tmp := big.NewInt(num)

	if tmp.Add(currentBalance, tmp).Cmp(beginBalance) != 0 {
		t.Fatalf("Test_HTLC_Create_Available_Gas balance not equal")
	}

}
