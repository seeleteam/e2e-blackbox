/**
*  @file
*  @copyright defined in go-seele/LICENSE
 */

package testcase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

// ResGetinfo
type ResGetInfo struct {
	CurrentBlockHeight int64  `json:"CurrentBlockHeight"`
	HeaderHash         string `json:"HeaderHash"`
	MinerStatus        string `json:"MinerStatus"`
	Shard              int    `json:"Shard"`
	Coinbase           string `json:"Coinbase"`
}

// BalanceInfo balance
type BalanceInfo struct {
	Account string
	Balance int64
}

type BlockHeader struct {
	CreateTimestamp   uint32
	Difficulty        uint64
	Height            uint64
	PreviousBlockHash string
}

type BlockInfo struct {
	Hash         string        `json:"hash"`
	Transactions []interface{} `json:"transactions"`
	Header       BlockHeader   `json:"header"`
}

type TxInfoInBlock struct {
	Hash     string `json:"hash"`
	From     string `json:"from"`
	To       string `json:"to"`
	Amount   int64  `json:"amount"`
	GasPrice int64  `json:"gasPrice"`
	GasLimit int64  `json:"gasLimit"`
}

// TxDataInfo tx data
type TxDataInfo struct {
	From     string
	To       string
	Amount   int64
	GasPrice int64
	GasLimit int64
}

// TxInfo tx
type TxInfo struct {
	Hash   string     `json:"hash"`
	TxData TxDataInfo `json:"Data"`
}

// HTLCSystemInfo htlc system
type HTLCSystemInfo struct {
	Tx         TxInfo `json:"Tx"`
	HashLock   string `json:"HashLock"`
	TimeLock   int64  `json:"TimeLock"`
	To         string `json:"To"`
	Refunded   bool   `json:"Refunded"`
	Withdrawed bool   `json:"Withdrawed"`
	Preimage   string `json:"Preimage"`
}

// HTLCCreateInfo htlc create
type HTLCCreateInfo struct {
	Tx       TxInfo `json:"Tx"`
	HashLock string `json:"HashLock"`
	TimeLock int64  `json:"TimeLock"`
}

// HTLCWithDrawInfo HTLC withdraw info
type HTLCWithDrawInfo struct {
	Tx       TxInfo `json:"Tx"`
	Hash     string `json:"hash"`
	PreImage string `json:"preimage"`
}

// HTLCRefundInfo refund info
type HTLCRefundInfo struct {
	Tx   TxInfo `json:"Tx"`
	Hash string `json:"hash"`
}

// ReceiptInfo receipt
type ReceiptInfo struct {
	Contract string        `json:"contract"`
	Failed   bool          `json:"failed"`
	TotalFee int64         `json:"totalFee"`
	UsedGas  int64         `json:"usedGas"`
	Result   string        `json:"result"`
	Hash     string        `json:"txhash"`
	Logs     []interface{} `json:"logs"`
}

// PoolTxInfo tx
type PoolTxInfo struct {
	Hash   string `json:"hash"`
	Nonce  int    `json:"accountNonce"`
	Amount int64  `json:"amount"`
	//From     string `json:"from"`
	//To       string `json:"to"`
	//GasLimit int    `json:"gasLimit"`
}

// SendTxInfo send tx
type SendTxInfo struct {
	nonce   int
	hash    string
	amount  int64
	gasUsed int
	bMined  bool
}

// LogByTopic contains Seele log
type LogByTopic struct {
	Log      Log    `json:"log"`
	LogIndex int    `json:"logIndex"`
	Txhash   string `json:"txhash"`
}

// Log contains Seele log of topic
type Log struct {
	Address          string   `json:"address"`
	BklockNumber     uint     `json:"blockNumber"`
	Data             string   `json:"data"`
	Topics           []string `json:"topics"`
	TransactionIndex int      `json:"transactionIndex"`
}

// TxByHashInfo output of gettxbyhash
type TxByHashInfo struct {
	BlockHash   string     `json:"blockHash"`
	Height      int        `json:"blockHeight"`
	TxIndex     int        `json:"txIndex"`
	Transaction PoolTxInfo `json:"transaction"`
}

func accountCase(command, account, accountMix string, t *testing.T) {
	cmd := exec.Command(CmdLight, "getbalance", "--account", account, "--address", ServerAddr)
	var output, outputMix []byte
	var err error
	if output, err = cmd.CombinedOutput(); err != nil {
		t.Fatalf("getbalance err: %s", err)
	}

	cmd = exec.Command(CmdLight, "getbalance", "--account", accountMix, "--address", ServerAddr)
	if outputMix, err = cmd.CombinedOutput(); err != nil {
		t.Fatalf("getbalance err: %s", err)
	}

	if string(output) != string(outputMix) {
		t.Fail()
	}
}

func getBalance(t *testing.T, command, account, serverAddr string) (int64, error) {
	cmd := exec.Command(command, "getbalance", "--account", account, "--address", serverAddr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	var info BalanceInfo
	if err = json.Unmarshal(output, &info); err != nil {
		return 0, err
	}

	return info.Balance, nil
}

func GetBlock(t *testing.T, command string, height int64, serverAddr string) (ret *BlockInfo, err error) {
	cmd := exec.Command(command, "getblock", "--height", strconv.FormatInt(height, 10), "--address", serverAddr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var info BlockInfo
	if err = json.Unmarshal(output, &info); err != nil {
		return nil, err
	}

	ret = &info
	return
}

func getNonce(t *testing.T, command, account, serverAddr string) (int, error) {
	cmd := exec.Command(command, "getnonce", "--account", account, "--address", serverAddr)
	//var curNonce int
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	output = bytes.Trim(output, "\n")
	//fmt.Println(string(output))
	return strconv.Atoi(string(output))
}

// SendTx send a tx
func SendTx(t *testing.T, command string, amount, nonce, gaslimit int, keystore, to, payload, serverAddr string) (txHash, debtHash string, err error) {
	if gaslimit <= 0 {
		gaslimit = 3000000
	}

	var cmd *exec.Cmd
	if payload == "" || payload == "0x" {
		cmd = exec.Command(command, "sendtx", "--amount", strconv.Itoa(amount), "--price", "1", "--gas", strconv.Itoa(gaslimit), "--from", keystore, "--to", to, "--nonce", strconv.Itoa(nonce))
	} else {
		cmd = exec.Command(command, "sendtx", "--amount", strconv.Itoa(amount), "--price", "1", "--gas", strconv.Itoa(gaslimit), "--from", keystore, "--to", to, "--nonce", strconv.Itoa(nonce), "--payload", payload)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		return
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()

	outStr, errStr := out.String(), outErr.String()

	if len(string(errStr)) > 0 {
		err = errors.New(string(errStr))
		return
	}
	fmt.Println("sendtx nonce=", nonce)

	outStr = outStr[strings.Index(outStr, "{"):]
	outStr = strings.Trim(outStr, "\n")
	outStr = strings.Trim(outStr, " ")
	// fmt.Println("sendtx out:[", outStr, "]")
	var txInfo TxInfo
	if err = json.Unmarshal([]byte(outStr), &txInfo); err != nil {
		return
	}

	txHash = txInfo.Hash
	return
}

func getPendingTxs(t *testing.T, command, serverAddr string) (infoL []PoolTxInfo, err error) {
	var output []byte
	cmd := exec.Command(command, "getpendingtxs", "--address", serverAddr)
	if output, err = cmd.CombinedOutput(); err != nil {
		return
	}
	//fmt.Println("pendingtxs:", string(output))
	var curL []PoolTxInfo
	if err = json.Unmarshal(output, &curL); err == nil {
		fmt.Println("getPendingTx:", curL)
		infoL = curL
	}
	return
}

func getPoolContentTxs(t *testing.T, command, serverAddr string) (infoM map[string][]PoolTxInfo, err error) {
	var output []byte
	cmd := exec.Command(command, "gettxpoolcontent", "--address", serverAddr)
	if output, err = cmd.CombinedOutput(); err != nil {
		return
	}
	//fmt.Println("gettxpoolcontent:", string(output))
	//var curM map[string][]PoolTxInfo
	if err = json.Unmarshal(output, &infoM); err == nil {
		/*for key, item := range curM {
			var infoL []*PoolTxInfo
			//infoL = append(infoL, &item)
			fmt.Println(key,item)
		}*/
	}
	return
}

func getPoolCountTxs(t *testing.T, command, serverAddr string) (int64, error) {
	var output []byte
	cmd := exec.Command(command, "gettxpoolcount", "--address", serverAddr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	str := string(output)
	str = strings.Trim(str, "\n")
	tmp, err := strconv.Atoi(str)

	return int64(tmp), nil
}

func GetTxByHash(t *testing.T, command, txHash, serverAddr string) (*TxByHashInfo, error) {
	cmd := exec.Command(command, "gettxbyhash", "--hash", txHash, "--address", serverAddr)

	output, errStr := cmd.CombinedOutput()
	if errStr != nil {
		return nil, errors.New(string(output))
	}

	//fmt.Println(string(output), errStr)

	var info TxByHashInfo
	if err := json.Unmarshal([]byte(output), &info); err != nil {
		return nil, err
	}

	return &info, nil
}

func GetReceipt(t *testing.T, command, txHash, serverAddr string) (*ReceiptInfo, error) {
	cmd := exec.Command(command, "getreceipt", "--hash", txHash, "--address", serverAddr)

	output, errStr := cmd.CombinedOutput()
	if errStr != nil {
		return nil, errors.New(string(output))
	}

	var info ReceiptInfo
	if err := json.Unmarshal([]byte(output), &info); err != nil {
		return nil, err
	}

	return &info, nil
}

func htlcDecode(t *testing.T, command, hexResult string) (*HTLCSystemInfo, error) {
	var output []byte
	cmd := exec.Command(command, "htlc", "decode", "--payload", hexResult)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var info HTLCSystemInfo
	if err := json.Unmarshal(output, &info); err != nil {
		return nil, err
	}

	return &info, nil
}

func findTxHashFromPool(txHash string, infoL *[]PoolTxInfo, infoM *map[string][]PoolTxInfo) (bPending, bContentPool bool) {
	if infoL != nil {
		for _, info := range *infoL {
			if info.Hash == txHash {
				bPending = true
			}
		}
	}

	if infoM != nil {
		for _, itemL := range *infoM {
			for _, info := range itemL {
				if info.Hash == txHash {
					bContentPool = true
				}
			}
		}
	}

	return
}

// generateTime generate time
func generateTime(minutes int64) int64 {
	return time.Now().Unix() + minutes*60
}

func deployContractAndSendTx(t *testing.T) (string, string, []string, error) {
	contract, err := ioutil.ReadFile("./contract/simplestorage/SimpleEvent.bin")
	if err != nil {
		return "", "", nil, fmt.Errorf("deployContractAndSendTx read contract failed %s", err.Error())
	}
	cmd := exec.Command(CmdClient, "sendtx", "--from", KeyFileShard1_1, "--amount", "0", "--payload", string(contract), "--address", ServerAddr)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", "", nil, fmt.Errorf("deployContractAndSendTx create contract err %s", err.Error())
	}
	defer stdin.Close()

	var (
		out    bytes.Buffer
		outErr bytes.Buffer
	)
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		return "", "", nil, fmt.Errorf("deployContractAndSendTx: An error occured: %s", err.Error())
	}
	io.WriteString(stdin, "123\n")
	cmd.Wait()

	output, errStr := out.String(), outErr.String()
	if errStr != "" {
		return "", "", nil, fmt.Errorf("deployContractAndSendTx cmd err: %s", errStr)
	}
	str := output[strings.Index(output, "{") : strings.LastIndex(output, "}")+1]
	var txInfo TxInfo
	if err = json.Unmarshal([]byte(str), &txInfo); err != nil {
		return "", "", nil, fmt.Errorf("deployContractAndSendTx create contract unmarshal err: %s", err)
	}
	for {
		time.Sleep(10)
		number, err := getPoolCountTxs(t, CmdClient, ServerAddr)
		if err != nil {
			return "", "", nil, fmt.Errorf("deployContractAndSendTx get pool count err: %s", err)
		}

		if number == 0 {
			break
		}
	}

	time.Sleep(20)

	receipt, err := GetReceipt(t, CmdClient, txInfo.Hash, ServerAddr)
	if err != nil {
		return "", "", nil, fmt.Errorf("deployContractAndSendTx get receipt err: %s", err)
	}
	if receipt.Failed {
		return "", "", nil, errors.New("deployContractAndSendTx tx operation fault")
	}

	cmd = exec.Command(CmdLight, "payload", "--abi", "./contract/simplestorage/SimpleEvent.abi", "--method", "get")
	method, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", nil, errors.New("deployContractAndSendTx returns false with valid parameter")
	}
	method = method[9 : len(method)-1]
	cmd = exec.Command(CmdClient, "sendtx", "--from", KeyFileShard1_1, "--to", receipt.Contract,
		"--amount", "0", "--payload", string(method), "--address", ServerAddr)
	stdin, err = cmd.StdinPipe()
	if err != nil {
		return "", "", nil, fmt.Errorf("deployContractAndSendTx call contract err: %s", err)
	}

	out.Reset()
	outErr.Reset()
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		return "", "", nil, fmt.Errorf("deployContractAndSendTx: An error occured: %s", err)
	}
	io.WriteString(stdin, "123\n")
	cmd.Wait()

	output1, errStr := out.String(), outErr.String()
	if errStr != "" {
		return "", "", nil, fmt.Errorf("deployContractAndSendTx cmd err: %s", errStr)
	}
	str1 := output1[strings.Index(output1, "{") : strings.LastIndex(output1, "}")+1]
	var tx TxInfo
	if err = json.Unmarshal([]byte(str1), &tx); err != nil {
		return "", "", nil, fmt.Errorf("deployContractAndSendTx call contract tx unmarshal err: %s", err)
	}
	for {
		time.Sleep(10)
		number, err := getPoolCountTxs(t, CmdClient, ServerAddr)
		if err != nil {
			return "", "", nil, fmt.Errorf("deployContractAndSendTx get pool count err: %s", err)
		}

		if number == 0 {
			break
		}
	}

	time.Sleep(20)

	receipt1, err := GetReceipt(t, CmdClient, tx.Hash, ServerAddr)
	if err != nil {
		return "", "", nil, fmt.Errorf("deployContractAndSendTx get receipt err: %s", err)
	}

	if receipt1.Failed {
		return "", "", nil, errors.New("deployContractAndSendTx tx operation fault")
	}
	var topics []string
	for _, log := range receipt1.Logs {
		l := log.(map[string]interface{})
		topic := l["topic"].(string)
		topics = append(topics, topic)
	}
	if len(topics) != 1 {
		return "", "", nil, errors.New("deployContractAndSendTx returns log number is not 1")
	}
	cmd = exec.Command(CmdClient, "gettxbyhash", "--hash", tx.Hash, "--address", ServerAddr)
	result, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", nil, fmt.Errorf("deployContractAndSendTx: An error occured: %s", err)
	}

	var txByHash map[string]interface{}
	if err = json.Unmarshal([]byte(result), &txByHash); err != nil {
		return "", "", nil, fmt.Errorf("deployContractAndSendTx get tx by hash unmarshal err: %s", err)
	}
	blockHeight := uint64(txByHash["blockHeight"].(float64))
	height := strconv.FormatUint(blockHeight, 10)
	return receipt.Contract, height, topics, nil
}
