/**
*  @file
*  @copyright defined in go-seele/LICENSE
 */

package testcase

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"os"
// 	"os/exec"
// 	"strings"
// 	"sync"
// 	"testing"
// 	"time"

// 	"github.com/seeleteam/e2e-blackbox/testcase/common"
// )

// // ==============================begin template command===============================================
// func Test_Client_SubChain_template_Invalid_Name(t *testing.T) {
// 	cmd := exec.Command(common.CmdClient, "subchain", "template", "--file", "subchain", "--name", "seele.123_we")

// 	var outErr bytes.Buffer
// 	cmd.Stderr = &outErr

// 	if err := cmd.Start(); err != nil {
// 		t.Fatalf("%s: An error occured: %s", "Test_Client_SubChain_template_Invalid_Name", err)
// 	}

// 	cmd.Wait()

// 	errStr := outErr.String()
// 	if !strings.Contains(errStr, "invalid name, only numbers, letters, and dash lines are allowed") {
// 		t.Fatalf("%s Err:%s", "Test_Client_SubChain_template_Invalid_Name", errStr)
// 	}
// }

// func Test_Client_SubChain_template_Invalid_Name_Empty(t *testing.T) {
// 	cmd := exec.Command(common.CmdClient, "subchain", "template", "--file", "subchain", "--name", "")

// 	var outErr bytes.Buffer
// 	cmd.Stderr = &outErr

// 	if err := cmd.Start(); err != nil {
// 		t.Fatalf("%s: An error occured: %s", "Test_Client_SubChain_template_Invalid_Name_Empty", err)
// 	}

// 	cmd.Wait()

// 	errStr := outErr.String()
// 	if !strings.Contains(errStr, "name is empty") {
// 		t.Fatalf("%s Err:%s", "Test_Client_SubChain_template_Invalid_Name_Empty", errStr)
// 	}
// }

// func Test_Client_SubChain_template_Invalid_Name_Exceed_Max_Length(t *testing.T) {
// 	domainName := ""
// 	for i := 0; i < 33; i++ {
// 		domainName += "s"
// 	}

// 	cmd := exec.Command(common.CmdClient, "subchain", "template", "--file", "subchain", "--name", domainName)

// 	var outErr bytes.Buffer
// 	cmd.Stderr = &outErr

// 	if err := cmd.Start(); err != nil {
// 		t.Fatalf("%s: An error occured: %s", "Test_Client_SubChain_template_Invalid_Name_Exceed_Max_Length", err)
// 	}

// 	cmd.Wait()

// 	errStr := outErr.String()
// 	if !strings.Contains(errStr, "name too long") {
// 		t.Fatalf("%s Err:%s", "Test_Client_SubChain_template_Invalid_Name_Exceed_Max_Length", errStr)
// 	}
// }

// func Test_Client_SubChain_template_Invalid_FilePath(t *testing.T) {
// 	cmd := exec.Command(common.CmdClient, "subchain", "template", "--file", "subc<<<<^%hain", "--name", "testinvalidfilepath")

// 	var outErr bytes.Buffer
// 	cmd.Stderr = &outErr

// 	if err := cmd.Start(); err != nil {
// 		t.Fatalf("%s: An error occured: %s", "Test_Client_SubChain_template_Invalid_FilePath", err)
// 	}

// 	cmd.Wait()

// 	errStr := outErr.String()
// 	if !strings.Contains(errStr, "The filename, directory name, or volume label syntax is incorrect") {
// 		t.Fatalf("%s Err:%s", "Test_Client_SubChain_template_Invalid_FilePath", errStr)
// 	}
// }

// func Test_Client_SubChain_template(t *testing.T) {
// 	subChainTemplate(t, "Test_Client_SubChain_template", "Test_Client_SubChain_template", "subchaintest")
// 	os.RemoveAll("Test_Client_SubChain_template")
// }

// func subChainTemplate(t *testing.T, funcName, subChainFile, domainName string) {
// 	cmd := exec.Command(common.CmdClient, "subchain", "template", "--file", subChainFile, "--name", domainName)

// 	var out bytes.Buffer
// 	var outErr bytes.Buffer
// 	cmd.Stdout, cmd.Stderr = &out, &outErr

// 	if err := cmd.Start(); err != nil {
// 		t.Fatalf("%s: An error occured: %s", funcName, err)
// 	}

// 	cmd.Wait()

// 	output, errStr := out.String(), outErr.String()
// 	if errStr != "" {
// 		t.Fatalf("%s cmd err: %s", funcName, errStr)
// 	}

// 	if !strings.Contains(output, "generate template json file for sub chain register successfully") {
// 		t.Fatalf("%s Err:%s", funcName, errStr)
// 	}
// }

// // ==============================end template command===============================================

// // ==============================begin register command===============================================
// func Test_Client_SubChain_register_Invalid_KeyFile(t *testing.T) {
// 	validateInfo := "invalid name"
// 	subChainInvalidRegister(t, "Test_Client_SubChain_register_Invalid_KeyFile", "common.KeyFileShard1_1",
// 		"123", "15", "200000", "", "subChainTemplate.json", validateInfo)
// }

// func Test_Client_SubChain_register_Unmatched_keyfile_And_Pass(t *testing.T) {
// 	validateInfo := "invalid name"
// 	subChainInvalidRegister(t, "Test_Client_SubChain_register_Unmatched_keyfile_And_Pass", common.KeyFileShard1_1,
// 		"12345", "15", "200000", "", "subChainTemplate.json", validateInfo)
// }

// func Test_Client_SubChain_register_Invalid_PriceValue(t *testing.T) {
// 	validateInfo := "invalid name"
// 	subChainInvalidRegister(t, "Test_Client_SubChain_register_Invalid_PriceValue", common.KeyFileShard1_1,
// 		"123", "q3", "200000", "", "subChainTemplate.json", validateInfo)
// }

// func Test_Client_SubChain_register_Invalid_Gas(t *testing.T) {
// 	validateInfo := "invalid value"
// 	subChainInvalidRegister(t, "Test_Client_SubChain_register_Invalid_Gas", common.KeyFileShard1_1,
// 		"123", "15", "qw", "", "subChainTemplate.json", validateInfo)
// }

// func Test_Client_SubChain_register_Invalid_Nonce(t *testing.T) {
// 	validateInfo := "invalid value"
// 	subChainInvalidRegister(t, "Test_Client_SubChain_register_Invalid_Nonce", common.KeyFileShard1_1,
// 		"123", "15", "200000", "er", "subChainTemplate.json", validateInfo)
// }

// // Name cannot contain special characters, such as /, \, `
// // it may be a string consisting of numbers, letters, and middle lines, etc.
// func Test_Client_SubChain_register_Invalid_Name(t *testing.T) {
// 	validateInfo := "invalid name, only numbers, letters, and dash lines are allowed"
// 	subChainInvalidRegister(t, "Test_Client_SubChain_register_Invalid_Name", common.KeyFileShard1_1,
// 		"123", "15", "200000", "", "subChainRegisterInvalidName.json", validateInfo)
// }

// func Test_Client_SubChain_register_Invalid_Name_Empty(t *testing.T) {
// 	validateInfo := "name is empty"
// 	subChainInvalidRegister(t, "Test_Client_SubChain_register_Invalid_Name_Empty", common.KeyFileShard1_1,
// 		"123", "15", "200000", "", "subChainRegisterNameEmpty.json", validateInfo)
// }

// func Test_Client_SubChain_register_Invalid_Name_Exceed_Max_Length(t *testing.T) {
// 	validateInfo := "name too long"
// 	subChainInvalidRegister(t, "Test_Client_SubChain_register_Invalid_Name_Exceed_Max_Length", common.KeyFileShard1_1,
// 		"123", "15", "200000", "", "subChainRegisterNameTooLong.json", validateInfo)
// }

// // func Test_Client_SubChain_register_Invalid_Name_Existed(t *testing.T) {
// // 	receipt := subChainRegister(t, "Test_Client_SubChain_register_Invalid_Name_Existed", "subChainTemplate.json")
// // 	if receipt.Failed {
// // 		t.Fatalf("Test_Client_SubChain_register_Invalid_Name_Existed recepit error, %s", receipt.Result)
// // 	}

// // 	receipt1 := subChainRegister(t, "Test_Client_SubChain_register_Invalid_Name_Existed", "subChainTemplate.json")
// // 	if !receipt1.Failed {
// // 		t.Fatalf("Test_Client_SubChain_register_Invalid_Name_Existed SubChain repeated registration successully")
// // 	}
// // 	if !strings.Contains(receipt1.Result, "already exists") {
// // 		t.Fatalf("Test_Client_SubChain_register_Invalid_Name_Existed result does not contain already exists, Result:%s", receipt1.Result)
// // 	}
// // }

// func Test_Client_SubChain_register_Version_Empty(t *testing.T) {
// 	validateInfo := "invalid subchain version"
// 	subChainInvalidRegister(t, "Test_Client_SubChain_register_Version_Empty", common.KeyFileShard1_1,
// 		"123", "15", "200000", "", "subChainRegisterVersionEmpty.json", validateInfo)
// }

// func Test_Client_SubChain_register_TokenFullName_Empty(t *testing.T) {
// 	validateInfo := "invalid subchain token full name"
// 	subChainInvalidRegister(t, "Test_Client_SubChain_register_TokenFullName_Empty", common.KeyFileShard1_1,
// 		"123", "15", "200000", "", "subChainRegisterTokenFullNameEmpty.json", validateInfo)
// }

// func Test_Client_SubChain_register_TokenFullName_Equal_defaultTokenFullName(t *testing.T) {
// 	validateInfo := "invalid subchain token full name"
// 	subChainInvalidRegister(t, "Test_Client_SubChain_register_TokenFullName_Equal_defaultTokenFullName", common.KeyFileShard1_1,
// 		"123", "15", "200000", "", "subChainRegisterDefaultTokenFullName.json", validateInfo)
// }

// func Test_Client_SubChain_register_TokenShortName_Empty(t *testing.T) {
// 	validateInfo := "invalid subchain token short name"
// 	subChainInvalidRegister(t, "Test_Client_SubChain_register_TokenShortName_Empty", common.KeyFileShard1_1,
// 		"123", "15", "200000", "", "subChainRegisterTokenShortNameEmpty.json", validateInfo)
// }

// func Test_Client_SubChain_register_TokenShortName_Equal_defaultTokenShortName(t *testing.T) {
// 	validateInfo := "invalid subchain token short name"
// 	subChainInvalidRegister(t, "Test_Client_SubChain_register_TokenShortName_Equal_defaultTokenShortName", common.KeyFileShard1_1,
// 		"123", "15", "200000", "", "subChainRegisterDefaultTokenShortName.json", validateInfo)
// }

// func Test_Client_SubChain_register_Invalid_TokenAmount(t *testing.T) {
// 	validateInfo := "invalid subchain token amount"
// 	subChainInvalidRegister(t, "Test_Client_SubChain_register_Invalid_TokenAmount", common.KeyFileShard1_1,
// 		"123", "15", "200000", "", "subChainRegisterTokenAmount.json", validateInfo)
// }

// func Test_Client_SubChain_register(t *testing.T) {
// 	receipt := subChainRegister(t, "Test_Client_SubChain_register", "subChainTemplate1.json")
// 	if receipt.Failed {
// 		t.Fatalf("Test_Client_SubChain_register recepit error, %s", receipt.Result)
// 	}
// }

// func subChainInvalidRegister(t *testing.T, funcName, keyFile, passWord, price, gas, nonce, subChainFile, validateInfo string) {
// 	if len(nonce) == 0 {
// 		accountNonce, err := common.GetNonce(t, common.CmdClient, common.AccountShard1_1, common.ServerAddr)
// 		if err != nil {
// 			t.Fatalf("%s, err:%s", funcName, err)
// 		}

// 		nonce = fmt.Sprintf("%d", accountNonce)
// 	}

// 	cmd := exec.Command(common.CmdClient, "subchain", "register", "--from", keyFile, "--price", price, "--gas", gas,
// 		"--nonce", nonce, "--file", subChainFile)

// 	stdin, err := cmd.StdinPipe()
// 	if err != nil {
// 		t.Fatalf("%s err: %s", funcName, err)
// 	}

// 	defer stdin.Close()

// 	var outErr bytes.Buffer
// 	cmd.Stderr = &outErr

// 	if err = cmd.Start(); err != nil {
// 		t.Fatalf("%s: An error occured: %s", funcName, err)
// 	}

// 	io.WriteString(stdin, passWord+"\n")
// 	cmd.Wait()

// 	errStr := outErr.String()
// 	if !strings.Contains(errStr, validateInfo) {
// 		t.Fatalf("%s err=:%s, should be %s", funcName, errStr, validateInfo)
// 	}
// }

// func subChainRegister(t *testing.T, funcName, domainName string) *common.ReceiptInfo {
// 	cmd := exec.Command(common.CmdClient, "subchain", "register", "--from", common.KeyFileShard1_1, "--file", domainName)

// 	stdin, err := cmd.StdinPipe()
// 	if err != nil {
// 		t.Fatalf("%s err: %s", funcName, err)
// 	}

// 	defer stdin.Close()

// 	var out bytes.Buffer
// 	var outErr bytes.Buffer
// 	cmd.Stdout, cmd.Stderr = &out, &outErr

// 	if err = cmd.Start(); err != nil {
// 		t.Fatalf("%s: An error occured: %s", funcName, err)
// 	}

// 	io.WriteString(stdin, "123\n")
// 	cmd.Wait()

// 	output, errStr := out.String(), outErr.String()
// 	if errStr != "" {
// 		t.Fatalf("%s cmd err:%s", funcName, errStr)
// 	}

// 	str := output[strings.Index(output, `"Tx":`)+5 : strings.LastIndex(output, "}")]

// 	var tx common.TxInfo
// 	if err := json.Unmarshal([]byte(str), &tx); err != nil {
// 		t.Fatalf("%s unmarshal register domain tx err: %s", funcName, err)
// 	}

// 	for {
// 		time.Sleep(10)
// 		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
// 		if err != nil {
// 			t.Fatalf("%s get pool count err: %s", funcName, err)
// 		}

// 		if number == 0 {
// 			break
// 		}
// 	}

// 	time.Sleep(20)

// 	receipt, err := common.GetReceipt(t, common.CmdClient, tx.Hash, common.ServerAddr)
// 	if err != nil {
// 		t.Fatalf("%s get receipt err: %s", funcName, err)
// 	}

// 	return receipt
// }

// // ==============================end register command===============================================

// // ==============================begin query command===============================================
// func Test_Client_SubChain_query_Invalid_KeyFile(t *testing.T) {
// 	validateInfo := "invalid sender key file"
// 	subChainInvalidQuery(t, "Test_Client_SubChain_query_Invalid_KeyFile", "common.KeyFileShard1_1",
// 		"123456", "15", "200000", "", "game", validateInfo)
// }

// func Test_Client_SubChain_query_Unmatched_keyfile_And_Pass(t *testing.T) {
// 	validateInfo := "invalid sender key file"
// 	subChainInvalidQuery(t, "Test_Client_SubChain_query_Unmatched_keyfile_And_Pass", common.KeyFileShard1_1,
// 		"123456", "15", "200000", "", "game", validateInfo)
// }

// func Test_Client_SubChain_query_Invalid_PriceValue(t *testing.T) {
// 	validateInfo := "invalid gas price value"
// 	subChainInvalidQuery(t, "Test_Client_SubChain_query_Invalid_PriceValue", common.KeyFileShard1_1,
// 		"123", "q2", "200000", "", "game", validateInfo)
// }

// func Test_Client_SubChain_query_Invalid_Gas(t *testing.T) {
// 	validateInfo := "invalid value"
// 	subChainInvalidQuery(t, "Test_Client_SubChain_query_Invalid_Gas", common.KeyFileShard1_1,
// 		"123", "15", "qw", "", "game", validateInfo)
// }

// func Test_Client_SubChain_query_Invalid_Nonce(t *testing.T) {
// 	validateInfo := "invalid value "
// 	subChainInvalidQuery(t, "Test_Client_SubChain_query_Invalid_Nonce", common.KeyFileShard1_1,
// 		"123", "15", "200000", "er", "game", validateInfo)
// }

// // Name cannot contain special characters, such as /, \, `
// // it may be a string consisting of numbers, letters, and middle lines, etc.
// func Test_Client_SubChain_query_Invalid_Name(t *testing.T) {
// 	validateInfo := "invalid name, only numbers, letters, and dash lines are allowed"
// 	subChainInvalidQuery(t, "Test_Client_SubChain_query_Invalid_Name", common.KeyFileShard1_1,
// 		"123", "15", "200000", "", "seele.game_23", validateInfo)
// }

// func Test_Client_SubChain_query_Invalid_Name_Empty(t *testing.T) {
// 	validateInfo := "name is empty"
// 	subChainInvalidQuery(t, "Test_Client_SubChain_query_Invalid_Name_Empty", common.KeyFileShard1_1,
// 		"123", "15", "200000", "", "", validateInfo)
// }

// func Test_Client_SubChain_query_Invalid_Name_Exceed_Max_Length(t *testing.T) {
// 	domainName := ""
// 	for i := 0; i < 33; i++ {
// 		domainName += "s"
// 	}

// 	validateInfo := "name too long"
// 	subChainInvalidQuery(t, "Test_Client_SubChain_query_Invalid_Name_Exceed_Max_Length", common.KeyFileShard1_1,
// 		"123", "15", "200000", "", domainName, validateInfo)
// }

// func Test_Client_SubChain_query(t *testing.T) {
// 	receipt := subChainRegister(t, "Test_Client_SubChain_register", "subChainTemplate_query.json")
// 	if receipt.Failed {
// 		t.Fatalf("Test_Client_SubChain_register recepit error, %s", receipt.Result)
// 	}

// 	receipt1 := subChainQuery(t, "Test_Client_SubChain_register", "testsubchaintemplatequery")
// 	if receipt1.Failed {
// 		t.Fatalf("Test_Client_SubChain_register recepit error, %s", receipt.Result)
// 	}
// }

// func subChainInvalidQuery(t *testing.T, funcName, keyFile, passWord, price, gas, nonce, domainName, validateInfo string) {
// 	if len(nonce) == 0 {
// 		accountNonce, err := common.GetNonce(t, common.CmdClient, common.AccountShard1_1, common.ServerAddr)
// 		if err != nil {
// 			t.Fatalf("%s, err:%s", funcName, err)
// 		}

// 		nonce = fmt.Sprintf("%d", accountNonce)
// 	}

// 	cmd := exec.Command(common.CmdClient, "subchain", "query", "--from", keyFile, "--price", price, "--gas", gas,
// 		"--nonce", nonce, "--name", domainName)

// 	stdin, err := cmd.StdinPipe()
// 	if err != nil {
// 		t.Fatalf("%s err: %s", funcName, err)
// 	}

// 	defer stdin.Close()

// 	var outErr bytes.Buffer
// 	cmd.Stderr = &outErr

// 	if err = cmd.Start(); err != nil {
// 		t.Fatalf("%s: An error occured: %s", funcName, err)
// 	}

// 	io.WriteString(stdin, passWord+"\n")
// 	cmd.Wait()

// 	errStr := outErr.String()
// 	if !strings.Contains(errStr, validateInfo) {
// 		t.Fatalf("%s Err:%s", funcName, errStr)
// 	}
// }

// func subChainQuery(t *testing.T, funcName, domainName string) *common.ReceiptInfo {
// 	cmd := exec.Command(common.CmdClient, "subchain", "query", "--from", common.KeyFileShard1_1, "--name", domainName)

// 	stdin, err := cmd.StdinPipe()
// 	if err != nil {
// 		t.Fatalf("%s err: %s", funcName, err)
// 	}

// 	defer stdin.Close()

// 	var out bytes.Buffer
// 	var outErr bytes.Buffer
// 	cmd.Stdout, cmd.Stderr = &out, &outErr

// 	if err = cmd.Start(); err != nil {
// 		t.Fatalf("%s: An error occured: %s", funcName, err)
// 	}

// 	io.WriteString(stdin, "123\n")
// 	cmd.Wait()

// 	output, errStr := out.String(), outErr.String()
// 	if errStr != "" {
// 		t.Fatalf("%s cmd err: %s", funcName, errStr)
// 	}

// 	str := output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]

// 	var tx common.TxInfo
// 	if err := json.Unmarshal([]byte(str), &tx); err != nil {
// 		t.Fatalf("%s unmarshal register domain tx err: %s", funcName, err)
// 	}

// 	for {
// 		time.Sleep(10)
// 		number, err := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
// 		if err != nil {
// 			t.Fatalf("%s get pool count err: %s", funcName, err)
// 		}

// 		if number == 0 {
// 			break
// 		}
// 	}

// 	time.Sleep(20)

// 	receipt, err := common.GetReceipt(t, common.CmdClient, tx.Hash, common.ServerAddr)
// 	if err != nil {
// 		t.Fatalf("%s get receipt err: %s", funcName, err)
// 	}

// 	return receipt
// }

// // ==============================end query command===============================================

// // ==============================begin config command===============================================
// var (
// 	isExecuteQuerySubchain = 0

// 	lock sync.Mutex
// )

// func querySubChain(t *testing.T, funcName string) {
// 	if isExecuteQuerySubchain == 0 {
// 		lock.Lock()
// 		if isExecuteQuerySubchain == 0 {
// 			receipt := subChainQuery(t, funcName, "testsubchaintemplateconfig")
// 			if receipt.Failed {
// 				t.Fatalf("%s recepit error, %s", funcName, receipt.Result)
// 			}

// 			if receipt.Result == "0x" {
// 				receipt1 := subChainRegister(t, funcName, "subChainTemplate_config.json")
// 				if receipt1.Failed {
// 					t.Fatalf("%s recepit error, %s", funcName, receipt.Result)
// 				}
// 			}
// 			isExecuteQuerySubchain = 1
// 		}
// 		lock.Unlock()
// 	}
// }

// func subChainInvalidConfig(t *testing.T, funcName, coinbase, algorithm, privatekey, shard, node, name, output, validateInfo string) {
// 	cmd := exec.Command(common.CmdClient, "subchain", "config", "--coinbase", coinbase, "--algorithm", algorithm,
// 		"--privatekey", privatekey, "--shard", shard, "--node", node, "--output", output, "--name", name)

// 	var out bytes.Buffer
// 	var outErr bytes.Buffer
// 	cmd.Stdout, cmd.Stderr = &out, &outErr

// 	if err := cmd.Start(); err != nil {
// 		t.Fatalf("%s: An error occured: %s", funcName, err)
// 	}

// 	cmd.Wait()

// 	outStr, errStr := out.String(), outErr.String()
// 	fmt.Println(outStr)
// 	if !strings.Contains(errStr, validateInfo) {
// 		t.Fatalf("%s Err:%s", funcName, errStr)
// 	}
// }

// func Test_Client_SubChain_config_Invalid_coinbase(t *testing.T) {
// 	querySubChain(t, "Test_Client_SubChain_config_Invalid_coinbase")
// 	subChainInvalidConfig(t, "Test_Client_SubChain_config_Invalid_coinbase", "2323", "sha256",
// 		"0xf65e40c6809643b25ce4df33153da2f3338876f181f83d2281c6ac4a987b1479", "1", "",
// 		"testsubchaintemplateconfig", "config", "invalid coinbase, err:hex string without 0x prefix")
// 	subChainInvalidConfig(t, "Test_Client_SubChain_config_Invalid_coinbase", "0x2323", "sha256",
// 		"0xf65e40c6809643b25ce4df33153da2f3338876f181f83d2281c6ac4a987b1479", "1", "",
// 		"testsubchaintemplateconfig", "config", "invalid coinbase, err:invalid address length 2, expected length is 20")
// 	subChainInvalidConfig(t, "Test_Client_SubChain_config_Invalid_coinbase", "0x4c10f2cd2159b", "sha256",
// 		"0xf65e40c6809643b25ce4df33153da2f3338876f181f83d2281c6ac4a987b1479", "1", "",
// 		"testsubchaintemplateconfig", "config", "hex string of odd length")
// }

// func Test_Client_SubChain_config_ShardofCoinbase_NotEqual_ShardValue(t *testing.T) {
// 	querySubChain(t, "Test_Client_SubChain_config_ShardofCoinbase_NotEqual_ShardValue")
// 	subChainInvalidConfig(t, "Test_Client_SubChain_config_ShardofCoinbase_NotEqual_ShardValue", "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21", "sha256",
// 		"0xf65e40c6809643b25ce4df33153da2f3338876f181f83d2281c6ac4a987b1479", "2", "",
// 		"testsubchaintemplateconfig", "config", "input shard(2) is not equal to shard nubmer(1) obtained from the input coinbase:")
// }

// func Test_Client_SubChain_config_Invalid_PrivateKey(t *testing.T) {
// 	querySubChain(t, "Test_Client_SubChain_config_Invalid_PrivateKey")
// 	subChainInvalidConfig(t, "Test_Client_SubChain_config_Invalid_PrivateKey", "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21", "sha256",
// 		"sdsd", "1", "",
// 		"testsubchaintemplateconfig", "config", "Input string not a valid ecdsa string")
// 	subChainInvalidConfig(t, "Test_Client_SubChain_config_Invalid_PrivateKey", "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21", "sha256",
// 		"0x4c10f2cd2159bb4", "1", "",
// 		"testsubchaintemplateconfig", "config", "invalid key: encoding/hex: odd length hex string")
// 	subChainInvalidConfig(t, "Test_Client_SubChain_config_Invalid_PrivateKey", "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21", "sha256",
// 		"0x2323", "1", "",
// 		"testsubchaintemplateconfig", "config", "invalid key: invalid length, need 256 bits")
// }

// // static node can be empty
// func Test_Client_SubChain_config_Invalid_StaticNode(t *testing.T) {
// 	querySubChain(t, "Test_Client_SubChain_config_Invalid_StaticNode")
// 	subChainInvalidConfig(t, "Test_Client_SubChain_config_Invalid_StaticNode", "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21", "sha256",
// 		"0xf65e40c6809643b25ce4df33153da2f3338876f181f83d2281c6ac4a987b1479", "1", "we23",
// 		"testsubchaintemplateconfig", "config", "address we23: missing port in address")

// 	subChainInvalidConfig(t, "Test_Client_SubChain_config_Invalid_StaticNode", "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21", "sha256",
// 		"0xf65e40c6809643b25ce4df33153da2f3338876f181f83d2281c6ac4a987b1479", "1", "we23:we",
// 		"testsubchaintemplateconfig", "config", "lookup udp/we: getaddrinfow: The specified class was not found.")

// 	subChainInvalidConfig(t, "Test_Client_SubChain_config_Invalid_StaticNode", "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21", "sha256",
// 		"0xf65e40c6809643b25ce4df33153da2f3338876f181f83d2281c6ac4a987b1479", "1", "we23:2323",
// 		"testsubchaintemplateconfig", "config", "lookup we23: no such host")
// }

// // Name cannot contain special characters, such as /, \, `
// // it may be a string consisting of numbers, letters, and middle lines, etc.
// func Test_Client_SubChain_config_Invalid_Name(t *testing.T) {
// 	subChainInvalidConfig(t, "Test_Client_SubChain_config_Invalid_Name", "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21", "sha256",
// 		"0xf65e40c6809643b25ce4df33153da2f3338876f181f83d2281c6ac4a987b1479", "1", "",
// 		"testsubchaintemplate.config", "config", "invalid name, only numbers, letters, and dash lines are allowed")
// }

// //
// func Test_Client_SubChain_config_Invalid_Name_Empty(t *testing.T) {
// 	querySubChain(t, "Test_Client_SubChain_config_Invalid_Name_Empty")
// 	subChainInvalidConfig(t, "Test_Client_SubChain_config_Invalid_Name_Empty", "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21", "sha256",
// 		"0xf65e40c6809643b25ce4df33153da2f3338876f181f83d2281c6ac4a987b1479", "1", "",
// 		"", "config", "name is empty")
// }

// func Test_Client_SubChain_config_Invalid_Name_Exceed_Max_Length(t *testing.T) {
// 	querySubChain(t, "Test_Client_SubChain_config_Invalid_Name_Exceed_Max_Length")
// 	domainName := ""
// 	for i := 0; i < 33; i++ {
// 		domainName += "s"
// 	}
// 	subChainInvalidConfig(t, "Test_Client_SubChain_config_Invalid_Name_Exceed_Max_Length", "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21", "sha256",
// 		"0xf65e40c6809643b25ce4df33153da2f3338876f181f83d2281c6ac4a987b1479", "1", "",
// 		domainName, "config", "name too long")
// }

// func Test_Client_SubChain_config_Invalid_Name_NotFound(t *testing.T) {
// 	subChainInvalidConfig(t, "Test_Client_SubChain_config_Invalid_Name_NotFound", "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21", "sha256",
// 		"0xf65e40c6809643b25ce4df33153da2f3338876f181f83d2281c6ac4a987b1479", "1", "",
// 		"testsubchaintempla-teconfig", "config", "sub-chain testsubchaintempla-teconfig does not exist")
// }

// func Test_Client_SubChain_config_Invalid_Shard(t *testing.T) {
// 	querySubChain(t, "Test_Client_SubChain_config_Invalid_Shard")
// 	subChainInvalidConfig(t, "Test_Client_SubChain_config_Invalid_Shard", "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21", "sha256",
// 		"0xf65e40c6809643b25ce4df33153da2f3338876f181f83d2281c6ac4a987b1479", "123", "",
// 		"testsubchaintemplateconfig", "config", "input shard(123) is not equal to shard nubmer(1) obtained from the input coinbase")
// }

// func Test_Client_SubChain_config_Invalid_OutPut(t *testing.T) {
// 	querySubChain(t, "Test_Client_SubChain_config_Invalid_OutPut")
// 	subChainInvalidConfig(t, "Test_Client_SubChain_config_Invalid_OutPut", "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21", "sha256",
// 		"0xf65e40c6809643b25ce4df33153da2f3338876f181f83d2281c6ac4a987b1479", "1", "",
// 		"testsubchaintemplateconfig", "conf<<<<^%ig", "mkdir conf<<<<^%ig: The filename, directory name, or volume label syntax is incorrect")
// }

// func Test_Client_SubChain_config(t *testing.T) {
// 	funcName := "Test_Client_SubChain_config"
// 	querySubChain(t, funcName)
// 	cmd := exec.Command(common.CmdClient, "subchain", "config", "--coinbase", "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21",
// 		"--name", "testsubchaintemplateconfig")

// 	var out bytes.Buffer
// 	var outErr bytes.Buffer
// 	cmd.Stdout, cmd.Stderr = &out, &outErr

// 	if err := cmd.Start(); err != nil {
// 		t.Fatalf("%s: An error occured: %s", funcName, err)
// 	}

// 	cmd.Wait()

// 	output, errStr := out.String(), outErr.String()
// 	if errStr != "" {
// 		t.Fatalf("%s cmd err: %s", funcName, errStr)
// 	}
// 	if !strings.Contains(output, "generate sub chain config files successfully") {
// 		t.Fatalf("%s run err: %s", funcName, output)
// 	}
// }

// // ==============================end config command===============================================
