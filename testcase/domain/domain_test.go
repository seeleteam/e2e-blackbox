/**
*  @file
*  @copyright defined in go-seele/LICENSE
 */

package domain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/seeleteam/e2e-blackbox/testcase/common"

	seele "github.com/seeleteam/go-seele/common"
)

func Test_Client_Domain_register_Invalid_KeyFile(t *testing.T) {
	validateInfo := `invalid sender key file`
	domainInvalidRegister(t, "Test_Client_Domain_register_Invalid_KeyFile", "common.KeyFileShard1_1",
		"123456", "15", "200000", "", "game", validateInfo)
}

func Test_Client_Domain_register_Unmatched_keyfile_And_Pass(t *testing.T) {
	validateInfo := `invalid sender key file`
	domainInvalidRegister(t, "Test_Client_Domain_register_Unmatched_keyfile_And_Pass", common.KeyFileShard1_1,
		"123456", "15", "200000", "", "game", validateInfo)
}

func Test_Client_Domain_register_Invalid_PriceValue(t *testing.T) {
	validateInfo := `invalid gas price value`
	domainInvalidRegister(t, "Test_Client_Domain_register_Invalid_PriceValue", common.KeyFileShard1_1,
		"123", "q2", "200000", "", "game", validateInfo)
}

func Test_Client_Domain_register_Invalid_Gas(t *testing.T) {
	validateInfo := `invalid value "qw" for flag -gas: strconv.ParseUint: parsing "qw"`
	domainInvalidRegister(t, "Test_Client_Domain_register_Invalid_Gas", common.KeyFileShard1_1,
		"123", "15", "qw", "", "game", validateInfo)
}

func Test_Client_Domain_register_Invalid_Nonce(t *testing.T) {
	validateInfo := `invalid value "er" for flag -nonce: strconv.ParseUint: parsing "er": invalid syntax`
	domainInvalidRegister(t, "Test_Client_Domain_register_Invalid_Nonce", common.KeyFileShard1_1,
		"123", "15", "200000", "er", "game", validateInfo)
}

// Name cannot contain special characters, such as /, \, `
// it may be a string consisting of numbers, letters, and middle lines, etc.
// func Test_Client_Domain_register_Invalid_Name(t *testing.T) {
// 	validateInfo := `invalid name, only numbers, letters, and dash lines are allowed`
// 	domainInvalidRegister(t, "Test_Client_Domain_register_Invalid_Name", common.KeyFileShard1_1,
// 		"123", "15", "200000", "", "seele.game_23", validateInfo)
// }

func Test_Client_Domain_register_Invalid_Name_Empty(t *testing.T) {
	validateInfo := `name is empty`
	domainInvalidRegister(t, "Test_Client_Domain_register_Invalid_Name_Empty", common.KeyFileShard1_1,
		"123", "15", "200000", "", "", validateInfo)
}

func Test_Client_Domain_register_Invalid_Name_Exceed_Max_Length(t *testing.T) {
	domainName := ""
	for i := 0; i < len(seele.EmptyHash)+1; i++ {
		domainName += "s"
	}

	validateInfo := `name too long`
	domainInvalidRegister(t, "Test_Client_Domain_register_Invalid_Name_Exceed_Max_Length", common.KeyFileShard1_1,
		"123", "15", "200000", "", domainName, validateInfo)
}

func Test_Client_Domain_register(t *testing.T) {
	receipt1 := domainRegister(t, "Test_Client_Domain_register", "seele")
	if receipt1.Failed {
		t.Fatalf("Test_Client_Domain_register recepit error, %s", receipt1.Result)
	}
}

func Test_Client_Domain_register_Invalid_Name_Existed(t *testing.T) {
	receipt1 := domainRegister(t, "Test_Client_Domain_register_Invalid_Name_Existed", "game")
	if receipt1.Failed {
		t.Fatalf("Test_Client_Domain_register_Invalid_Name_Existed recepit error, %s", receipt1.Result)
	}

	receipt2 := domainRegister(t, "Test_Client_Domain_register_Invalid_Name_Existed", "game")

	if !receipt2.Failed {
		t.Fatalf("Test_Client_Domain_register_Invalid_Name_Existed, Domain name repeated registration successully")
	}

	if !strings.Contains(receipt2.Result, "already exists") {
		t.Fatalf("Test_Client_Domain_register_Invalid_Name_Existed result does not contain already exists, Result:%s", receipt2.Result)
	}
}

func domainInvalidRegister(t *testing.T, funcName, keyFile, passWord, price, gas, nonce, domainName, validateInfo string) {
	if len(nonce) == 0 {
		accountNonce, err := common.GetNonce(t, common.CmdClient, common.AccountShard1_1, common.ServerAddr)
		if err != nil {
			t.Fatalf("%s, err:%s", funcName, err)
		}

		nonce = fmt.Sprintf("%d", accountNonce)
	}

	cmd := exec.Command(common.CmdClient, "domain", "register", "--from", keyFile, "--price", price, "--gas", gas,
		"--nonce", nonce, "--name", domainName)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("%s err: %s", funcName, err)
	}

	defer stdin.Close()

	var outErr bytes.Buffer
	cmd.Stderr = &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("%s: An error occured: %s", funcName, err)
	}

	io.WriteString(stdin, passWord+"\n")
	cmd.Wait()

	errStr := outErr.String()
	if !strings.Contains(errStr, validateInfo) {
		t.Fatalf("%s get err=:%s, should be %s", funcName, errStr, validateInfo)
	}
}

func domainRegister(t *testing.T, funcName, domainName string) *common.ReceiptInfo {
	cmd := exec.Command(common.CmdClient, "domain", "register", "--from", common.KeyFileShard1_1, "--price", "15", "--gas", "200000",
		"--name", domainName)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("%s err: %s", funcName, err)
	}

	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("%s: An error occured: %s", funcName, err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()

	output, errStr := out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("%s cmd err: %s", funcName, errStr)
	}

	str := output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]

	var tx common.TxInfo
	if err := json.Unmarshal([]byte(str), &tx); err != nil {
		t.Fatalf("%s unmarshal register domain tx err: %s", funcName, err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("%s get pool count err: %s", funcName, err)
		}

		if number == 0 {
			break
		}
	}

	time.Sleep(20)

	receipt, err := common.GetReceipt(t, common.CmdClient, tx.Hash, common.ServerAddr)
	if err != nil {
		t.Fatalf("%s get receipt err: %s", funcName, err)
	}

	return receipt
}

func Test_Client_Domain_owner_Invalid_KeyFile(t *testing.T) {
	validateInfo := `invalid sender key file`
	domainInvalidOwner(t, "Test_Client_Domain_owner_Invalid_KeyFile", "common.KeyFileShard1_1",
		"123456", "15", "200000", "", "game", validateInfo)
}

func Test_Client_Domain_owner_Unmatched_keyfile_And_Pass(t *testing.T) {
	validateInfo := `invalid sender key file`
	domainInvalidOwner(t, "Test_Client_Domain_owner_Unmatched_keyfile_And_Pass", common.KeyFileShard1_1,
		"123456", "15", "200000", "", "game", validateInfo)
}

func Test_Client_Domain_owner_Invalid_PriceValue(t *testing.T) {
	validateInfo := `invalid gas price value`
	domainInvalidOwner(t, "Test_Client_Domain_owner_Invalid_PriceValue", common.KeyFileShard1_1,
		"123", "q2", "200000", "", "game", validateInfo)
}

func Test_Client_Domain_owner_Invalid_Gas(t *testing.T) {
	validateInfo := `invalid value "qw" for flag -gas: strconv.ParseUint: parsing "qw"`
	domainInvalidOwner(t, "Test_Client_Domain_owner_Invalid_Gas", common.KeyFileShard1_1,
		"123", "15", "qw", "", "game", validateInfo)
}

func Test_Client_Domain_owner_Invalid_Nonce(t *testing.T) {
	validateInfo := `invalid value "er" for flag -nonce: strconv.ParseUint: parsing "er": invalid syntax`
	domainInvalidOwner(t, "Test_Client_Domain_owner_Invalid_Nonce", common.KeyFileShard1_1,
		"123", "15", "200000", "er", "game", validateInfo)
}

// Name cannot contain special characters, such as /, \, `
// it may be a string consisting of numbers, letters, and middle lines, etc.
// func Test_Client_Domain_owner_Invalid_Name(t *testing.T) {
// 	validateInfo := `invalid name, only numbers, letters, and dash lines are allowed`
// 	domainInvalidOwner(t, "Test_Client_Domain_owner_Invalid_Name", common.KeyFileShard1_1,
// 		"123", "15", "200000", "", "seele.game_23", validateInfo)
// }

func Test_Client_Domain_owner_Invalid_Name_Empty(t *testing.T) {
	validateInfo := `name is empty`
	domainInvalidOwner(t, "Test_Client_Domain_owner_Invalid_Name_Empty", common.KeyFileShard1_1,
		"123", "15", "200000", "", "", validateInfo)
}

func Test_Client_Domain_owner_Invalid_Name_Exceed_Max_Length(t *testing.T) {
	domainName := ""
	for i := 0; i < len(seele.EmptyHash)+1; i++ {
		domainName += "s"
	}

	validateInfo := `name too long`
	domainInvalidOwner(t, "Test_Client_Domain_owner_Invalid_Name_Exceed_Max_Length", common.KeyFileShard1_1,
		"123", "15", "200000", "", domainName, validateInfo)
}

func Test_Client_Domain_owner_Invalid_Name_Not_Found(t *testing.T) {
	receipt1 := domainOwner(t, "Test_Client_Domain_owner_Invalid_Name_Not_Found", "testownernotfound")
	if !receipt1.Failed {
		t.Fatalf("Test_Client_Domain_owner_Invalid_Name_Not_Found get domain, result:%s", receipt1.Result)
	}

	if !strings.Contains(receipt1.Result, "Failed to get data with key") {
		t.Fatalf("Test_Client_Domain_owner_Invalid_Name_Not_Found, result:%s", receipt1.Result)
	}
}

func Test_Client_Domain_owner(t *testing.T) {
	receipt1 := domainRegister(t, "Test_Client_Domain_owner", "testowner")
	if receipt1.Failed {
		t.Fatalf("Test_Client_Domain_owner register domain error, result:%s", receipt1.Result)
	}

	receipt2 := domainOwner(t, "Test_Client_Domain_owner", "testowner")
	if receipt2.Failed {
		t.Fatalf("Test_Client_Domain_owner get domain, result:%s", receipt2.Result)
	}
}

func domainInvalidOwner(t *testing.T, funcName, keyFile, passWord, price, gas, nonce, domainName, validateInfo string) {
	if len(nonce) == 0 {
		accountNonce, err := common.GetNonce(t, common.CmdClient, common.AccountShard1_1, common.ServerAddr)
		if err != nil {
			t.Fatalf("%s, err:%s", funcName, err)
		}

		nonce = fmt.Sprintf("%d", accountNonce)
	}

	cmd := exec.Command(common.CmdClient, "domain", "owner", "--from", keyFile, "--price", price, "--gas", gas,
		"--nonce", nonce, "--name", domainName)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("%s err: %s", funcName, err)
	}

	defer stdin.Close()

	var outErr bytes.Buffer
	cmd.Stderr = &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("%s: An error occured: %s", funcName, err)
	}

	io.WriteString(stdin, passWord+"\n")
	cmd.Wait()

	errStr := outErr.String()
	if !strings.Contains(errStr, validateInfo) {
		t.Fatalf("%s err:%s, should be %s", funcName, errStr, validateInfo)
	}
}

func domainOwner(t *testing.T, funcName, domainName string) *common.ReceiptInfo {
	cmd := exec.Command(common.CmdClient, "domain", "owner", "--from", common.KeyFileShard1_1, "--price", "15", "--gas", "200000",
		"--name", domainName)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("%s err: %s", funcName, err)
	}

	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("%s: An error occured: %s", funcName, err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()

	output, errStr := out.String(), outErr.String()
	if errStr != "" {
		t.Fatalf("%s cmd err: %s", funcName, errStr)
	}

	str := output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]

	var tx common.TxInfo
	if err := json.Unmarshal([]byte(str), &tx); err != nil {
		t.Fatalf("%s unmarshal register domain tx err: %s", funcName, err)
	}

	for {
		time.Sleep(10)
		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err != nil {
			t.Fatalf("%s get pool count err: %s", funcName, err)
		}

		if number == 0 {
			break
		}
	}

	time.Sleep(20)

	receipt, err := common.GetReceipt(t, common.CmdClient, tx.Hash, common.ServerAddr)
	if err != nil {
		t.Fatalf("%s get receipt err: %s", funcName, err)
	}

	return receipt
}
