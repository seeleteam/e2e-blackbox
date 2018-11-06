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
)

/*
func getPendingTxs(t *testing.T, command, serverAddr string) (infoL []*PoolTxInfo, err error) {
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
	accountCase(CmdLight, Account1_Aux2, Account1_Aux2Mix, t)
}

func Test_Light_GetBalance_InvalidAccount(t *testing.T) {
	if _, err := getBalance(t, CmdLight, AccountErr, ServerAddr); err == nil {
		t.Fatalf("getbalance AccountErr success?")
	}
}

func Test_Light_GetBalance_InvalidAccountType(t *testing.T) {
	if _, err := getBalance(t, CmdLight, InvalidAccountType, ServerAddr); err == nil {
		t.Fatalf("getbalance InvalidAccountType success? should return error")
	}
}

func Test_Light_GetBalance_AccountFromOtherShard(t *testing.T) {
	if _, err := getBalance(t, CmdLight, AccountShard2_1, ServerAddr); err == nil {
		t.Fatalf("getbalance account from other shard success? should return error")
	}
}

func Test_Light_GetBlock_ByInvalidHeight(t *testing.T) {
	// invalid height
	cmd := exec.Command(CmdLight, "getblock", "--height", "100000000", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getblock error parameter success?")
	}
}

func Test_Light_GetBlock_ByInvalidHash(t *testing.T) {
	cmd := exec.Command(CmdLight, "getblock", "--hash", BlockHashErr, "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getblock error parameter success?")
	}
}

func Test_Light_GetBlock_ByInvalidHash0x(t *testing.T) {
	cmd := exec.Command(CmdLight, "getblock", "--hash", "0x", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getblock error parameter success?")
	}
}

func Test_Light_GetBlock_InvalidParameter(t *testing.T) {
	cmd := exec.Command(CmdLight, "getblock", "1", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getblock error parameter success?")
	}
}

func Test_Light_GetBlock_ByHeight(t *testing.T) {
	cmd := exec.Command(CmdLight, "getblock", "--height", "0", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("getblock error, %s", err)
	}
}

func Test_Light_GetBlock_ByHash(t *testing.T) {
	cmd := exec.Command(CmdLight, "getblock", "--hash", BlockHash, "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("getblock error, %s", err)
	}
}

func Test_Light_GetBlock_Fulltx(t *testing.T) {
	// getblock fulltx support.
	cmd := exec.Command(CmdLight, "getblock", "--height", "1", "--fulltx", "--address", ServerAddr)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("getblock error, %s", err)
	} else {
		var blockInfo BlockInfo
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
	cmd := exec.Command(CmdLight, "getblockheight", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("getblockheight error, %s", err)
	}
}

func Test_Light_GetBlockHeight_InvalidParameter(t *testing.T) {
	cmd := exec.Command(CmdLight, "getblockheight", "100", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getblockheight returns ok with invalid parameter")
	}
}

func Test_Light_GetBlockTXCount_ByInvalidHeight(t *testing.T) {
	cmd := exec.Command(CmdLight, "getblocktxcount", "--height", "100000000", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getblocktxcount error parameter success?")
	}
}

func Test_Light_GetBlockTXCount_ByInvalidHash(t *testing.T) {
	cmd := exec.Command(CmdLight, "getblocktxcount", "--hash", BlockHashErr, "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getblocktxcount error parameter success?")
	}
}

func Test_Light_GetBlockTXCount_ByInvalidHash0x(t *testing.T) {
	cmd := exec.Command(CmdLight, "getblocktxcount", "--hash", "0x", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getblocktxcount error parameter success?")
	}
}

func Test_Light_GetBlockTXCount_InvalidParameter(t *testing.T) {
	cmd := exec.Command(CmdLight, "getblocktxcount", "1", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getblocktxcount error parameter success?")
	}
}

func Test_Light_GetBlockTXCount_ByHeight(t *testing.T) {
	cmd := exec.Command(CmdLight, "getblocktxcount", "--height", "0", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("getblocktxcount error, %s", err)
	}
}

func Test_Light_GetBlockTXCount_ByHash(t *testing.T) {
	cmd := exec.Command(CmdLight, "getblocktxcount", "--hash", BlockHash, "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("getblocktxcount error, %s %s", err, cmd.Args)
	}
}

/*
func Test_Light_SendTx(t *testing.T) {
	cmd := exec.Command(CmdLight, "sendtx", "--amount", "10000", "--price", "1", "--from", KeyFileShard1_3, "--to", Account1_Aux)
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
	curNonce, err := getNonce(t, CmdLight, Account1, ServerAddr)
	if err != nil {
		t.Fatalf("getnonce returns with error input", err)
	}

	cmd := exec.Command(CmdLight, "sendtx", "--amount", "10000", "--price", "1", "--from", KeyFileShard1_3, "--to", Account1_Aux, "--nonce", strconv.Itoa(curNonce+1))
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

func Test_light_SendTx_CrossShard(t *testing.T) {
	curNonce, err := getNonce(t, CmdLight, AccountShard1_4, ServerAddr)
	if err != nil {
		t.Fatalf("getnonce returns with error input %s", err)
	}

	for cnt := 0; cnt < 100; cnt++ {
		itemNonce := curNonce + 2 + cnt
		txHash, debtHash, err := SendTx(t, CmdLight, 10000, itemNonce, 21000, KeyFileShard1_4, AccountShard2_4, "", ServerAddr)
		if err != nil {
			t.Fatalf("Test_Light_SendTx: An error occured: %s", err)
		}
		fmt.Println("txHash=", txHash, " debtHash=", debtHash)
	}
}

func Test_Light_SendManyTx(t *testing.T) {
	curNonce, err := getNonce(t, CmdLight, AccountShard1_5, ServerAddr)
	if err != nil {
		t.Fatalf("getnonce returns with error input %s", err)
	}

	var beginBalance, dstBeginBalance int64
	beginBalance, err = getBalance(t, CmdLight, AccountShard1_5, ServerAddr)
	if err != nil {
		t.Fatalf("getBalance returns with error input %s", err)
	}

	dstBeginBalance, err = getBalance(t, CmdLight, Account1_Aux2, ServerAddr)
	if err != nil {
		t.Fatalf("getBalance returns with error input %s", err)
	}

	fmt.Println("fromAccount=", beginBalance, "dstAccount=", dstBeginBalance)
	var txHash string
	var sendTxL []*SendTxInfo

	var maxSendNonce, curNonceAfter int
	for cnt := 0; cnt < 100; cnt++ {
		itemNonce := curNonce + 2 + cnt
		txHash, _, err = SendTx(t, CmdLight, 10000, itemNonce, 21000, KeyFileShard1_5, Account1_Aux2, "", ServerAddr)
		if err != nil {
			t.Fatalf("Test_Light_SendTx: An error occured: %s", err)
		}
		maxSendNonce = itemNonce
		info := &SendTxInfo{
			nonce:  itemNonce,
			hash:   txHash,
			bMined: false,
		}
		sendTxL = append(sendTxL, info)
		//time.Sleep(8 * time.Second)
	}

	time.Sleep(8 * time.Second)
	cnt := 0
	for {
		pendingL, err1 := getPendingTxs(t, CmdLight, ServerAddr)
		if err1 != nil {
			t.Fatalf("getPendingTxs err:%s", err1)
		}
		contentM, err2 := getPoolContentTxs(t, CmdLight, ServerAddr)
		if err2 != nil {
			t.Fatalf("getPoolContentTxs err:%s", err1)
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
		info, err3 := GetReceipt(t, CmdLight, sendTxInfo.hash, ServerAddr)
		if err3 == nil {
			//t.Fatalf("getReceipt err:%s", err3)
			if info.Hash != sendTxInfo.hash {
				fmt.Println("XXXXXXX Receipt Hash not match with tx")
			}
			validCnt++
			sendTxInfo.bMined = true
		} else {
			fmt.Println("getReceipt err. nonce=", sendTxInfo.nonce, err3)
		}
	}

	var endBalance, dstEndBalance int64
	endBalance, err = getBalance(t, CmdLight, AccountShard1_5, ServerAddr)
	if err != nil {
		t.Fatalf("getBalance returns with error input %s", err)
	}

	dstEndBalance, err = getBalance(t, CmdLight, Account1_Aux2, ServerAddr)
	if err != nil {
		t.Fatalf("getBalance returns with error input %s", err)
	}

	curNonceAfter, err = getNonce(t, CmdLight, AccountShard1_5, ServerAddr)
	if err != nil {
		t.Fatalf("getnonce returns with error input %s", err)
	}

	fmt.Println("account1=", endBalance, "dstAccount=", dstEndBalance)
	fmt.Println("diff account1=", beginBalance-endBalance, "dstAccount=", dstEndBalance-dstBeginBalance)
	fmt.Println("sendMaxNonce=", maxSendNonce, " nonce from chain=", curNonceAfter)
	fmt.Println("validTx=", validCnt, "account1_times=", (beginBalance-endBalance)/31000)
	for _, sendTxInfo := range sendTxL {
		fmt.Println("./client gettxbyhash --hash ", sendTxInfo.hash)
	}
}

func Test_Light_GetReceipt_old(t *testing.T) {
	curNonce, err := getNonce(t, CmdLight, AccountShard1_5, ServerAddr)
	if err != nil {
		t.Fatalf("getnonce returns with error input err: %s", err)
	}

	var beginBalance, dstBeginBalance int64
	beginBalance, err = getBalance(t, CmdLight, AccountShard1_5, ServerAddr)
	if err != nil {
		t.Fatalf("getBalance returns with error input err: %s", err)
	}

	dstBeginBalance, err = getBalance(t, CmdLight, Account1_Aux2, ServerAddr)
	if err != nil {
		t.Fatalf("getBalance returns with error input err: %s", err)
	}

	fmt.Println("account1=", beginBalance, "dstAccount=", dstBeginBalance)
	var txHash string
	var sendTxL []*SendTxInfo

	for cnt := 0; cnt < 100; cnt++ {
		itemNonce := curNonce + 2 + cnt
		txHash, _, err = SendTx(t, CmdLight, 10000, itemNonce, 21000, KeyFileShard1_5, Account1_Aux2, "", ServerAddr)
		if err != nil {
			t.Fatalf("Test_Light_SendTx: An error occured: %s", err)
		}
		info := &SendTxInfo{
			nonce:  itemNonce,
			hash:   txHash,
			bMined: false,
		}
		sendTxL = append(sendTxL, info)
		//time.Sleep(8 * time.Second)
	}

	for {
		pendingL, err1 := getPendingTxs(t, CmdLight, ServerAddr)
		if err1 != nil {
			t.Fatalf("getPendingTxs err:%s", err1)
		}
		contentM, err2 := getPoolContentTxs(t, CmdLight, ServerAddr)
		if err2 != nil {
			t.Fatalf("getPoolContentTxs err:%s", err1)
		}

		bAllMined := true
		for _, sendTxInfo := range sendTxL {
			if sendTxInfo.bMined {
				continue
			}

			bPending, bContent := findTxHashFromPool(sendTxInfo.hash, &pendingL, &contentM)
			if bPending || bContent {
				bAllMined = false
				continue
			}

			//
			//var receiptInfo *ReceiptInfo
			_, err3 := GetReceipt(t, CmdLight, sendTxInfo.hash, ServerAddr)
			if err3 == nil {
				//t.Fatalf("getReceipt err:%s", err3)
				sendTxInfo.bMined = true
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
	endBalance, err = getBalance(t, CmdLight, AccountShard1_5, ServerAddr)
	if err != nil {
		t.Fatalf("getBalance returns with error input err: %s", err)
	}

	dstEndBalance, err = getBalance(t, CmdLight, Account1_Aux2, ServerAddr)
	if err != nil {
		t.Fatalf("getBalance returns with error input err: %s", err)
	}

	fmt.Println("account1=", endBalance, "dstAccount=", dstEndBalance)
	fmt.Println("diff account1=", beginBalance-endBalance, "dstAccount=", dstEndBalance-dstBeginBalance)

	for _, sendTxInfo := range sendTxL {
		fmt.Println("./client gettxbyhash --hash ", sendTxInfo.hash)
	}
}

func Test_Light_SendTx_InvalidAccountLength(t *testing.T) {
	cmd := exec.Command(CmdLight, "sendtx", "--amount", "10000", "--price", "1", "--from", KeyFileShard1_1, "--to", "0x")
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
	cmd := exec.Command(CmdLight, "sendtx", "--amount", "10000", "--price", "1", "--from", KeyFileShard1_3, "--to", InvalidAccountType)
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
	if !strings.Contains(errStr, " unsupported address type") {
		t.Fatalf("Test_Light_SendTx_InvalidAccountType Err:%s", errStr)
	}
}

func Test_Light_GetShardNum_InvalidAccountType(t *testing.T) {
	cmd := exec.Command(CmdLight, "getshardnum", "--account", InvalidAccountType)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getshardnum should return error with invalid account type")
	}
}

func Test_Light_GetShardNum_ByPrivateKey(t *testing.T) {
	cmd := exec.Command(CmdLight, "getshardnum", "--privatekey", AccountPrivateKey2)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("getshardnum error:%s", err)
	} else {
		if !strings.Contains(string(output), "2") {
			t.Fatalf("getshardnum returns error shardnum")
		}
	}
}

func Test_Light_GetShardNum(t *testing.T) {
	cmd := exec.Command(CmdLight, "getshardnum", "--account", AccountShard2_1)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("getshardnum error:%s", err)
	} else {
		if !strings.Contains(string(output), "2") {
			t.Fatalf("getshardnum returns error shardnum")
		}
	}
}

func Test_Light_GetNonce_InvalidAccount0x(t *testing.T) {
	cmd := exec.Command(CmdLight, "getnonce", "--account", "0x", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getnonce returns with error input err: %s", err)
	}
}

func Test_Light_GetNonce_InvalidAccount(t *testing.T) {
	cmd := exec.Command(CmdLight, "getnonce", "--account", AccountErr, "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getnonce returns with error input")
	}
}

func Test_Light_GetTxInBlock_ByHeight(t *testing.T) {
	//cmd := exec.Command(CmdLight, "gettxinblock", "--height", "1", "--index", "0")
	cmd := exec.Command(CmdLight, "gettxinblock", "--height", "4839", "--index", "5")
	if output, err := cmd.CombinedOutput(); err == nil {
		var info TxInfoInBlock
		if err := json.Unmarshal(output, &info); err != nil {
			t.Fatalf("Test_Light_GetTxInBlock_ByHeight unmarshal Failed. err=%s", err)
		}
	} else {
		t.Fatalf("Test_Light_GetTxInBlock_ByHeight Failed. err=%s", err)
	}
}

func Test_Light_GetTxInBlock_ByHash(t *testing.T) {
	cmd := exec.Command(CmdLight, "gettxinblock", "--hash", BlockHash, "--index", "0")
	if output, err := cmd.CombinedOutput(); err == nil {
		var info TxInfoInBlock
		if err := json.Unmarshal(output, &info); err != nil {
			t.Fatalf("Test_Light_GetTxInBlock_ByHash unmarshal Failed. err=%s", err)
		}
	} else {
		t.Fatalf("Test_Light_GetTxInBlock_ByHash Failed. err=%s", err)
	}
}

func Test_Light_GetTxInBlock_ByHash0x(t *testing.T) {
	cmd := exec.Command(CmdLight, "gettxinblock", "--hash", "0x", "--index", "0")
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

func Test_Light_GetNonce_AccountFromOtherShard(t *testing.T) {
	cmd := exec.Command(CmdLight, "getnonce", "--account", AccountShard1_1, "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getnonce returns successfully for other shard account")
	}
}

func Test_CheckChain_Consistent(t *testing.T) {
	block, err := GetBlock(t, CmdClient, -1, ServerAddr)
	if err != nil {
		t.Fatalf("Test_CheckChain_Consistent getBlock err. %s", err)
	}

	toHeight := block.Header.Height
	if toHeight > 10 {
		toHeight = toHeight - 6
	}

	block, err = GetBlock(t, CmdClient, 0, ServerAddr)
	if err != nil {
		t.Fatalf("Test_CheckChain_Consistent getBlock err. %s", err)
	}

	//	toHeight = 2001
	preHash, preTimestamp, preHeight := block.Hash, block.Header.CreateTimestamp, block.Header.Height

	allTime := uint32(0)
	maxTime := uint32(0)
	intL := make([]int, 10)
	for cur := uint64(1); cur <= toHeight; cur++ {
		block, err = GetBlock(t, CmdClient, int64(cur), ServerAddr)
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

		if cur%1000 == 0 {
			fmt.Println("Test_CheckChain_Consistent checked. cur=", cur)
		}
	}

	fmt.Println("Test_CheckChain_Consistent average block-creation time=", allTime/uint32(toHeight), " maxTime =", maxTime)
	//for _, cnt := range intL {
	fmt.Println("Test_CheckChain_Consistent:", intL)
	//}
}
