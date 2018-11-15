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
	"strings"
	"testing"
	"time"

	"github.com/seeleteam/e2e-blackbox/testcase/common"
)

/*
func common.GetPendingTxs(t *testing.T, command, serverAddr string) (infoL []*PoolTxInfo, err error) {
	var output []byte
	cmd := exec.Command(command, "getpendingtxs", "--address", serverAddr)
	if output, err = cmd.CombinedOutput(); err != nil {
		return
	}
	//fmt.Println("pendingtxs:", string(output))
	var curL []PoolTxInfo
	if err = json.Unmarshal(output, &curL); err == nil {
		for _, item := range curL {
			infoL = append(infoL, &item)
		}
	}
	return
}
*/

// account should ignore character case.
func Test_Light_AccountIgnoreCase(t *testing.T) {
	common.AccountCase(common.CmdLight, common.Account1_Aux2, common.Account1_Aux2Mix, t)
}

func Test_Light_GetBalance_InvalidAccount(t *testing.T) {
	if _, err := common.GetBalance(t, common.CmdLight, common.AccountErr, common.ServerAddr); err == nil {
		t.Fatalf("getbalance common.AccountErr success?")
	}
}

func Test_Light_GetBalance_InvalidAccountType(t *testing.T) {
	if _, err := common.GetBalance(t, common.CmdLight, common.InvalidAccountType, common.ServerAddr); err == nil {
		t.Fatalf("getbalance common.InvalidAccountType success? should return error")
	}
}

func Test_Light_GetBalance_AccountFromOtherShard(t *testing.T) {
	if _, err := common.GetBalance(t, common.CmdLight, common.AccountShard2_1, common.ServerAddr); err == nil {
		t.Fatalf("getbalance account from other shard success? should return error")
	}
}

func Test_Light_GetBlock_ByInvalidHeight(t *testing.T) {
	// invalid height
	cmd := exec.Command(common.CmdLight, "getblock", "--height", "100000000", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getblock error parameter success?")
	}
}

func Test_Light_GetBlock_ByInvalidHash(t *testing.T) {
	cmd := exec.Command(common.CmdLight, "getblock", "--hash", common.BlockHashErr, "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getblock error parameter success?")
	}
}

func Test_Light_GetBlock_ByInvalidHash0x(t *testing.T) {
	cmd := exec.Command(common.CmdLight, "getblock", "--hash", "0x", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getblock error parameter success?")
	}
}

func Test_Light_GetBlock_InvalidParameter(t *testing.T) {
	cmd := exec.Command(common.CmdLight, "getblock", "1", "--address", common.ServerAddr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("getblock error parameter success?")
	}

	if !strings.Contains(string(out), "flag is not specified for value") {
		t.Fatal("Test_Light_GetBlock_InvalidParameter is not ok")
	}
}

func Test_Light_GetBlock_ByHeight(t *testing.T) {
	cmd := exec.Command(common.CmdLight, "getblock", "--height", "0", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("getblock error, %s", err)
	}
}

// func Test_Light_GetBlock_ByHash(t *testing.T) {
// 	cmd := exec.Command(common.CmdLight, "getblock", "--hash", common.BlockHash, "--address", common.ServerAddr)
// 	if _, err := cmd.CombinedOutput(); err != nil {
// 		t.Fatalf("getblock error, %s", err)
// 	}
// }

func Test_Light_GetBlock_Fulltx(t *testing.T) {
	// getblock fulltx support.
	cmd := exec.Command(common.CmdLight, "getblock", "--height", "1", "--fulltx", "--address", common.ServerAddr)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("getblock error, %s", err)
	} else {
		var blockInfo common.BlockInfo
		//fmt.Println(string(output))
		if err = json.Unmarshal(output, &blockInfo); err != nil {
			t.Fatalf("Test_Light_GetBlock_Fulltx: %s", err)
		}
		if len(blockInfo.Transactions) <= 0 {
			t.Fatalf("Test_Light_GetBlock_Fulltx, block should contain one transaction at lease")
		}
	}
}

func Test_Light_GetBlockHeight(t *testing.T) {
	cmd := exec.Command(common.CmdLight, "getblockheight", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("getblockheight error, %s", err)
	}
}

func Test_Light_GetBlockHeight_InvalidParameter(t *testing.T) {
	cmd := exec.Command(common.CmdLight, "getblockheight", "100", "--address", common.ServerAddr)
	out, err := cmd.CombinedOutput()

	if !strings.Contains(string(out), "flag is not specified for value") {
		t.Fatal("Test_Light_GetBlockHeight_InvalidParameter is not ok")
	}

	if err != nil {
		t.Fatalf("getblockheight returns ok with invalid parameter")
	}
}

func Test_Light_GetBlockTXCount_ByInvalidHeight(t *testing.T) {
	cmd := exec.Command(common.CmdLight, "getblocktxcount", "--height", "100000000", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getblocktxcount error parameter success?")
	}
}

func Test_Light_GetBlockTXCount_ByInvalidHash(t *testing.T) {
	cmd := exec.Command(common.CmdLight, "getblocktxcount", "--hash", common.BlockHashErr, "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getblocktxcount error parameter success?")
	}
}

func Test_Light_GetBlockTXCount_ByInvalidHash0x(t *testing.T) {
	cmd := exec.Command(common.CmdLight, "getblocktxcount", "--hash", "0x", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getblocktxcount error parameter success?")
	}
}

func Test_Light_GetBlockTXCount_InvalidParameter(t *testing.T) {
	cmd := exec.Command(common.CmdLight, "getblocktxcount", "1", "--address", common.ServerAddr)

	out, err := cmd.CombinedOutput()

	if !strings.Contains(string(out), "flag is not specified for value") {
		t.Fatal("Test_Light_GetBlockTXCount_InvalidParameter is not ok")
	}

	if err != nil {
		t.Fatalf("getblocktxcount error parameter success?")
	}
}

func Test_Light_GetBlockTXCount_ByHeight(t *testing.T) {
	cmd := exec.Command(common.CmdLight, "getblocktxcount", "--height", "0", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("getblocktxcount error, %s", err)
	}
}

// func Test_Light_GetBlockTXCount_ByHash(t *testing.T) {
// 	cmd := exec.Command(common.CmdLight, "getblocktxcount", "--hash", common.BlockHash, "--address", common.ServerAddr)
// 	if _, err := cmd.CombinedOutput(); err != nil {
// 		t.Fatalf("getblocktxcount error, %s %s", err, cmd.Args)
// 	}
// }

/*
func Test_Light_SendTx(t *testing.T) {
	cmd := exec.Command(common.CmdLight, "sendtx", "--amount", "10000", "--price", "1", "--from", common.KeyFileShard1_3, "--to", common.Account1_Aux)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_Light_SendTx: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()

	outStr, errStr := out.String(), outErr.String()
	fmt.Println(outStr, errStr)
	if strings.Contains(errStr, "Failed to call rpc") {
		t.Fatalf("Test_Light_SendTx Err:%s", errStr)
	}
}

func Test_Light_SendTx_RemoveTimestamp(t *testing.T) {
	curNonce, err := common.GetNonce(t, common.CmdLight, Account1, common.ServerAddr)
	if err != nil {
		t.Fatalf("getnonce returns with error input", err)
	}

	cmd := exec.Command(common.CmdLight, "sendtx", "--amount", "10000", "--price", "1", "--from", common.KeyFileShard1_3, "--to", common.Account1_Aux, "--nonce", strconv.Itoa(curNonce+1))
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_Light_SendTx: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()

	outStr, errStr := out.String(), outErr.String()
	fmt.Println(outStr, errStr)
	if len(string(errStr)) > 0 {
		t.Fatalf("Test_Light_SendTx_RemoveTimestamp sendtx error. %s %s", errStr, cmd.Args)
	}

	if strings.Contains(outStr, "Timestamp") {
		t.Fatalf("Test_Light_SendTx_RemoveTimestamp should remove Timestamp item from json")
	}
}
*/

// func Test_light_SendTx_CrossShard(t *testing.T) {
// 	curNonce, err := common.GetNonce(t, common.CmdLight, common.AccountShard1_4, common.ServerAddr)
// 	if err != nil {
// 		t.Fatalf("getnonce returns with error input %s", err)
// 	}

// 	for cnt := 0; cnt < 100; cnt++ {
// 		itemNonce := curNonce + 2 + cnt
// 		txHash, debtHash, err := common.SendTx(t, common.CmdLight, 10000, itemNonce, 21000, common.KeyFileShard1_4, common.AccountShard2_4, "", common.ServerAddr)
// 		if err != nil {
// 			t.Fatalf("Test_Light_SendTx: An error occured: %s", err)
// 		}
// 		fmt.Println("txHash=", txHash, " debtHash=", debtHash)
// 	}
// }

func test_Light_TxInPool(t *testing.T) {
	curNonce, err := common.GetNonce(t, common.CmdLight, common.AccountShard1_3, common.ServerAddr)
	if err != nil {
		t.Fatalf("getnonce returns with error input %s", err)
	}
	fmt.Println("nonce=", curNonce)
	var beginBalance, dstBeginBalance int64
	beginBalance, err = common.GetBalance(t, common.CmdLight, common.AccountShard1_3, common.ServerAddr)
	if err != nil {
		t.Fatalf("common.GetBalance returns with error input %s", err)
	}

	dstBeginBalance, err = common.GetBalance(t, common.CmdLight, common.Account1_Aux, common.ServerAddr)
	if err != nil {
		t.Fatalf("common.GetBalance returns with error input %s", err)
	}

	fmt.Println("fromAccount=", beginBalance, "dstAccount=", dstBeginBalance)
	var txHash string
	itemNonce := curNonce + 1
	txHash, _, err = common.SendTx(t, common.CmdLight, 10000, itemNonce, 21000, common.KeyFileShard1_3, common.Account1_Aux, "", common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_Light_SendTx: An error occured: %s", err)
	}
	time.Sleep(2 * time.Second)
	info, err3 := common.GetTxByHash(t, common.CmdLight, txHash, common.ServerAddr)

	if err3 != nil {
		t.Fatalf("Test_Light_SendTx: An error occured when GetTxByHash:%v, err=%s", info, err3)
	}
}

func test_Light_SendManyTx(t *testing.T) {
	curNonce, err := common.GetNonce(t, common.CmdLight, common.AccountShard1_5, common.ServerAddr)
	if err != nil {
		t.Fatalf("getnonce returns with error input %s", err)
	}

	var beginBalance, dstBeginBalance int64
	beginBalance, err = common.GetBalance(t, common.CmdLight, common.AccountShard1_5, common.ServerAddr)
	if err != nil {
		t.Fatalf("common.GetBalance returns with error input %s", err)
	}

	dstBeginBalance, err = common.GetBalance(t, common.CmdLight, common.Account1_Aux2, common.ServerAddr)
	if err != nil {
		t.Fatalf("common.GetBalance returns with error input %s", err)
	}

	fmt.Println("fromAccount=", beginBalance, "dstAccount=", dstBeginBalance)
	var txHash string
	var sendTxL []*common.SendTxInfo

	var maxSendNonce, curNonceAfter int
	for cnt := 0; cnt < 100; cnt++ {
		itemNonce := curNonce + 2 + cnt
		txHash, _, err = common.SendTx(t, common.CmdLight, 10000, itemNonce, 21000, common.KeyFileShard1_5, common.Account1_Aux2, "", common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_Light_SendTx: An error occured: %s", err)
		}
		maxSendNonce = itemNonce
		info := &common.SendTxInfo{
			Nonce:  itemNonce,
			Hash:   txHash,
			BMined: false,
		}
		sendTxL = append(sendTxL, info)
		//time.Sleep(8 * time.Second)
	}

	time.Sleep(8 * time.Second)
	cnt := 0
	for {
		pendingL, err1 := common.GetPendingTxs(t, common.CmdLight, common.ServerAddr)
		if err1 != nil {
			t.Fatalf("common.GetPendingTxs err:%s", err1)
		}
		contentM, err2 := common.GetPoolContentTxs(t, common.CmdLight, common.ServerAddr)
		if err2 != nil {
			t.Fatalf("common.GetPoolContentTxs err:%s", err1)
		}

		if len(pendingL)+len(contentM) == 0 || cnt > 10 {
			break
		}
		cnt++
		time.Sleep(3 * time.Second)
	}

	time.Sleep(8 * time.Second)
	validCnt := 0
	for _, sendTxInfo := range sendTxL {
		//var receiptInfo *ReceiptInfo
		info, err3 := common.GetReceipt(t, common.CmdLight, sendTxInfo.Hash, common.ServerAddr)
		if err3 == nil {
			//t.Fatalf("getReceipt err:%s", err3)
			if info.Hash != sendTxInfo.Hash {
				fmt.Println("XXXXXXX Receipt Hash not match with tx")
			}
			validCnt++
			sendTxInfo.BMined = true
		} else {
			fmt.Println("getReceipt err. nonce=", sendTxInfo.Nonce, err3)
		}
	}

	var endBalance, dstEndBalance int64
	endBalance, err = common.GetBalance(t, common.CmdLight, common.AccountShard1_5, common.ServerAddr)
	if err != nil {
		t.Fatalf("common.GetBalance returns with error input %s", err)
	}

	dstEndBalance, err = common.GetBalance(t, common.CmdLight, common.Account1_Aux2, common.ServerAddr)
	if err != nil {
		t.Fatalf("common.GetBalance returns with error input %s", err)
	}

	curNonceAfter, err = common.GetNonce(t, common.CmdLight, common.AccountShard1_5, common.ServerAddr)
	if err != nil {
		t.Fatalf("getnonce returns with error input %s", err)
	}

	fmt.Println("account1=", endBalance, "dstAccount=", dstEndBalance)
	fmt.Println("diff account1=", beginBalance-endBalance, "dstAccount=", dstEndBalance-dstBeginBalance)
	fmt.Println("sendMaxNonce=", maxSendNonce, " nonce from chain=", curNonceAfter)
	fmt.Println("validTx=", validCnt, "account1_times=", (beginBalance-endBalance)/31000)
	for _, sendTxInfo := range sendTxL {
		fmt.Println("./client gettxbyhash --hash ", sendTxInfo.Hash)
	}

	for _, sendTxInfo := range sendTxL {
		fmt.Println("./client getreceipt --hash ", sendTxInfo.Hash)

	}
}

func Test_Light_GetReceipt_old(t *testing.T) {
	curNonce, err := common.GetNonce(t, common.CmdLight, common.AccountShard1_5, common.ServerAddr)
	if err != nil {
		t.Fatalf("getnonce returns with error input err: %s", err)
	}

	var beginBalance, dstBeginBalance int64
	beginBalance, err = common.GetBalance(t, common.CmdLight, common.AccountShard1_5, common.ServerAddr)
	if err != nil {
		t.Fatalf("common.GetBalance returns with error input err: %s", err)
	}

	dstBeginBalance, err = common.GetBalance(t, common.CmdLight, common.Account1_Aux2, common.ServerAddr)
	if err != nil {
		t.Fatalf("common.GetBalance returns with error input err: %s", err)
	}

	fmt.Println("account1=", beginBalance, "dstAccount=", dstBeginBalance)
	var txHash string
	var sendTxL []*common.SendTxInfo

	for cnt := 0; cnt < 100; cnt++ {
		itemNonce := curNonce + 2 + cnt
		txHash, _, err = common.SendTx(t, common.CmdLight, 10000, itemNonce, 21000, common.KeyFileShard1_5, common.Account1_Aux2, "", common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_Light_SendTx: An error occured: %s", err)
		}
		info := &common.SendTxInfo{
			Nonce:  itemNonce,
			Hash:   txHash,
			BMined: false,
		}
		sendTxL = append(sendTxL, info)
		//time.Sleep(8 * time.Second)
	}

	for {
		pendingL, err1 := common.GetPendingTxs(t, common.CmdLight, common.ServerAddr)
		if err1 != nil {
			t.Fatalf("common.GetPendingTxs err:%s", err1)
		}
		contentM, err2 := common.GetPoolContentTxs(t, common.CmdLight, common.ServerAddr)
		if err2 != nil {
			t.Fatalf("common.GetPoolContentTxs err:%s", err1)
		}

		bAllMined := true
		for _, sendTxInfo := range sendTxL {
			if sendTxInfo.BMined {
				continue
			}

			bPending, bContent := common.FindTxHashFromPool(sendTxInfo.Hash, &pendingL, &contentM)
			if bPending || bContent {
				bAllMined = false
				continue
			}

			//
			//var receiptInfo *ReceiptInfo
			_, err3 := common.GetReceipt(t, common.CmdLight, sendTxInfo.Hash, common.ServerAddr)
			if err3 == nil {
				//t.Fatalf("getReceipt err:%s", err3)
				sendTxInfo.BMined = true
			} else {
				bAllMined = false
			}

			//sendTxInfo.amount = receiptInfo.
		}

		if bAllMined {
			break
		}
		time.Sleep(5 * time.Second)
	}

	var endBalance, dstEndBalance int64
	endBalance, err = common.GetBalance(t, common.CmdLight, common.AccountShard1_5, common.ServerAddr)
	if err != nil {
		t.Fatalf("common.GetBalance returns with error input err: %s", err)
	}

	dstEndBalance, err = common.GetBalance(t, common.CmdLight, common.Account1_Aux2, common.ServerAddr)
	if err != nil {
		t.Fatalf("common.GetBalance returns with error input err: %s", err)
	}

	fmt.Println("account1=", endBalance, "dstAccount=", dstEndBalance)
	fmt.Println("diff account1=", beginBalance-endBalance, "dstAccount=", dstEndBalance-dstBeginBalance)

	// for _, sendTxInfo := range sendTxL {
	// 	fmt.Println("./client gettxbyhash --hash ", sendTxInfo.Hash)
	// }
}

func Test_Light_SendTx_InvalidAccountLength(t *testing.T) {
	cmd := exec.Command(common.CmdLight, "sendtx", "--amount", "10000", "--price", "1", "--from", common.KeyFileShard1_1, "--to", "0x")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_Light_SendTx_InvalidAccountLength: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()

	_, errStr := out.String(), outErr.String()
	//fmt.Println(outStr, errStr)
	if !strings.Contains(errStr, "invalid address") {
		t.Fatalf("Test_Light_SendTx_InvalidAccountLength Err:%s", errStr)
	}
}

func Test_Light_SendTx_InvalidAccountType(t *testing.T) {
	cmd := exec.Command(common.CmdLight, "sendtx", "--amount", "10000", "--price", "1", "--from", common.KeyFileShard1_3, "--to", common.InvalidAccountType)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_Light_SendTx_InvalidATest_Light_SendTx_InvalidAccountTypeccountLength: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()

	_, errStr := out.String(), outErr.String()
	//fmt.Println(outStr, errStr)
	if !strings.Contains(errStr, " invalid address type") {
		t.Fatalf("Test_Light_SendTx_InvalidAccountType Err:%s", errStr)
	}
}

func Test_Light_GetShardNum_InvalidAccountType(t *testing.T) {
	cmd := exec.Command(common.CmdLight, "getshardnum", "--account", common.InvalidAccountType)
	out, err := cmd.CombinedOutput()

	if !strings.Contains(string(out), "nvalid address type") {
		t.Fatal("Test_Light_GetShardNum_InvalidAccountType is not ok")
	}

	if err == nil {
		t.Fatalf("Test_Light_GetShardNum_InvalidAccountType getshardnum should return error with invalid account type")
	}
}

func Test_Light_GetShardNum_ByPrivateKey(t *testing.T) {
	cmd := exec.Command(common.CmdLight, "getshardnum", "--privatekey", common.AccountPrivateKey2)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("getshardnum error:%s", err)
	} else {
		if !strings.Contains(string(output), "2") {
			t.Fatalf("getshardnum returns error shardnum")
		}
	}
}

func Test_Light_GetShardNum(t *testing.T) {
	cmd := exec.Command(common.CmdLight, "getshardnum", "--account", common.AccountShard2_1)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("getshardnum error:%s", err)
	} else {
		if !strings.Contains(string(output), "2") {
			t.Fatalf("getshardnum returns error shardnum")
		}
	}
}

func Test_Light_GetNonce_InvalidAccount0x(t *testing.T) {
	cmd := exec.Command(common.CmdLight, "getnonce", "--account", "0x", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getnonce returns with error input err: %s", err)
	}
}

func Test_Light_GetNonce_InvalidAccount(t *testing.T) {
	cmd := exec.Command(common.CmdLight, "getnonce", "--account", common.AccountErr, "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getnonce returns with error input")
	}
}

// func Test_Light_GetTxInBlock_ByHeight(t *testing.T) {
// 	cmd := exec.Command(common.CmdLight, "gettxinblock", "--height", "4839", "--index", "5")
// 	if output, err := cmd.CombinedOutput(); err == nil {
// 		var info common.TxInfoInBlock
// 		if err := json.Unmarshal(output, &info); err != nil {
// 			t.Fatalf("Test_Light_GetTxInBlock_ByHeight unmarshal Failed. err=%s", err)
// 		}
// 	} else {
// 		t.Fatalf("Test_Light_GetTxInBlock_ByHeight Failed. err=%s", err)
// 	}
// }

// func Test_Light_GetTxInBlock_ByHash(t *testing.T) {
// 	cmd := exec.Command(common.CmdLight, "gettxinblock", "--hash", common.BlockHash, "--index", "0")
// 	if output, err := cmd.CombinedOutput(); err == nil {
// 		var info common.TxInfoInBlock
// 		if err := json.Unmarshal(output, &info); err != nil {
// 			t.Fatalf("Test_Light_GetTxInBlock_ByHash unmarshal Failed. err=%s", err)
// 		}
// 	} else {
// 		t.Fatalf("Test_Light_GetTxInBlock_ByHash Failed. err=%s", err)
// 	}
// }

func Test_Light_GetTxInBlock_ByHash0x(t *testing.T) {
	cmd := exec.Command(common.CmdLight, "gettxinblock", "--hash", "0x", "--index", "0")
	if output, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Light_GetTxInBlock_ByHash0x failed.")
	} else {
		//t.Fatalf("Test_Light_GetTxInBlock_ByHash0x Failed. err=%s", err)
		outStr := string(output)
		if strings.Index(outStr, "max index is -1") > 0 {
			t.Fatalf("Test_Light_GetTxInBlock_ByHash0x failed. comment is wrong.")
		}
	}
}

// func Test_Light_GetNonce_AccountFromOtherShard(t *testing.T) {
// 	cmd := exec.Command(common.CmdLight, "getnonce", "--account", common.AccountShard2_1, "--address", common.ServerAddr)
// 	_, err := cmd.CombinedOutput()

// 	if err == nil {
// 		t.Fatalf("getnonce returns successfully for other shard account")
// 	}
// }

func Test_CheckChain_Consistent(t *testing.T) {
	block, err := common.GetBlock(t, common.CmdClient, -1, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_CheckChain_Consistent getBlock err. %s", err)
	}

	toHeight := block.Header.Height
	if toHeight > 10 {
		toHeight = toHeight - 6
	}

	block, err = common.GetBlock(t, common.CmdClient, 0, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_CheckChain_Consistent getBlock err. %s", err)
	}

	//	toHeight = 2001
	preHash, preTimestamp, preHeight := block.Hash, block.Header.CreateTimestamp, block.Header.Height

	allTime := uint32(0)
	maxTime := uint32(0)
	intL := make([]int, 10)
	allTXs := 0
	blockCnt := 0
	for cur := uint64(1); cur <= toHeight; cur++ {
		block, err = common.GetBlock(t, common.CmdClient, int64(cur), common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_CheckChain_Consistent getBlock err. %s", err)
		}

		if block.Header.PreviousBlockHash != preHash || block.Header.Height != cur {
			t.Fatalf("Test_CheckChain_Consistent preHash not match: preHash=%s curHash=%s curHeight=%d cur=%d", preHash, block.Header.PreviousBlockHash, block.Header.Height, cur)
		}

		if block.Header.CreateTimestamp < preTimestamp {
			t.Fatalf("Test_CheckChain_Consistent timestamp not match. curHeight=%d cur=%d", block.Header.Height, cur)
		}

		if block.Header.Height != preHeight+1 {
			t.Fatalf("Test_CheckChain_Consistent height not match. curHeight=%d cur=%d", block.Header.Height, cur)
		}

		diffTime := block.Header.CreateTimestamp - preTimestamp
		blockCnt = blockCnt + 1
		if diffTime < 10000 {
			allTime = allTime + diffTime
			idx := diffTime / 10
			if idx > 9 {
				idx = 9
			}

			intL[idx] = intL[idx] + 1
			if maxTime < diffTime {
				maxTime = diffTime
			}
		}
		preHash, preTimestamp, preHeight = block.Hash, block.Header.CreateTimestamp, block.Header.Height
		allTXs = allTXs + len(block.Transactions)
		if cur%1000 == 0 {
			fmt.Println("Test_CheckChain_Consistent checked. cur=", cur, " txs=", allTXs)
		}
	}

	// fmt.Println("Test_CheckChain_Consistent average block-creation time=", allTime/uint32(toHeight), " maxTime =", maxTime)
	// fmt.Println("Test_CheckChain_Consistent txs=", allTXs, ", blockCnt=", blockCnt, " avg_txs/block =", allTXs/blockCnt)
	// fmt.Println("Test_CheckChain_Consistent txs=", allTXs)
	// //for _, cnt := range intL {
	// fmt.Println("Test_CheckChain_Consistent:", intL)
	//}
}
