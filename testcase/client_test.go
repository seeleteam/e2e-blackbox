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
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

type ResGetInfo struct {
	CurrentBlockHeight int64  `json:"CurrentBlockHeight"`
	HeaderHash         string `json:"HeaderHash"`
	MinerStatus        string `json:"MinerStatus"`
	Shard              int    `json:"Shard"`
	Coinbase           string `json:"Coinbase"`
}

func Test_Client_GetInfo(t *testing.T) {
	cmd := exec.Command(CmdClient, "getinfo")
	res, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Test_Client_GetInfo: GetInfo error, %s", err)
	}

	var r ResGetInfo
	err = json.Unmarshal(res, &r)
	if err != nil {
		t.Fatalf("Test_Client_GetInfo: decode return result error %s", err)
	}

	if r.MinerStatus != "Running" {
		t.Fatalf("Test_Client_GetInfo: Node not running!")
	}
}

func Test_Client_Key(t *testing.T) {
	cmd := exec.Command(CmdClient, "key")
	res, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Test_Client_Key: Key error, %s", err)
	}

	re, err := regexp.Compile("public(.+)")

	keyField := strings.Split(re.FindString(string(res)), "  ")

	if len(keyField[1]) == 0 {
		t.Fatal("Test_Client_Key: public key not found!")
	}

	re = regexp.MustCompile("private(.+)")
	keyField = strings.Split(re.FindString(string(res)), " ")

	if len(keyField[2]) == 0 {
		t.Fatal("Test_Client_Key: private key not found!")
	}
}

func Test_Client_DumpHeap(t *testing.T) {
	cmd := exec.Command(CmdClient, "dumpheap")
	res, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Test_Client_DumpHeap: dumpheap error, %s", err)
	}

	if _, err = os.Stat(strings.TrimSpace(string(res))); os.IsNotExist(err) {
		t.Fatalf("Test_Client_DumpHeap: file %s not found!", string(res))
	}
}

func Test_Client_Dumpheap_Default_Filename(t *testing.T) {
	userPath, err := user.Current()
	if err != nil {
		t.Fatalf("Test_Client_Dumpheap_Default_Filename get user path failed %s", err.Error())
	}
	defaultDataFolder := filepath.Join(userPath.HomeDir, ".seele")
	defaultFilePath := filepath.Join(defaultDataFolder, "heap.dump\n")

	cmd := exec.Command(CmdClient, "dumpheap", "--address", ServerAddr)
	file, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Test_Client_Dumpheap_Default_Filename: An error occured: %s", err.Error())
	}

	if defaultFilePath != string(file) {
		t.Fatal("Test_Client_Dumpheap_Default_Filename: The actual dumpheapPath is not equal to the expected path")
	}
}

func Test_Client_Dumpheap_Specified_Filename(t *testing.T) {
	userPath, err := user.Current()
	if err != nil {
		t.Fatalf("Test_Client_Dumpheap_Specified_Filename: Get user path failed %s", err.Error())
	}
	defaultDataFolder := filepath.Join(userPath.HomeDir, ".seele")
	defaultFilePath := filepath.Join(defaultDataFolder, "test.dump\n")

	cmd := exec.Command(CmdClient, "dumpheap", "--address", ServerAddr, "--file", "test.dump")
	file, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Test_Client_Dumpheap_Specified_Filename: An error occured: %s", err.Error())
	}

	if defaultFilePath != string(file) {
		t.Fatal("Test_Client_Dumpheap_Specified_Filename: The actual dumpheapPath is not equal to the expected path")
	}
}

func Test_Client_Payload_ValidParameter(t *testing.T) {
	cmd := exec.Command(CmdClient, "payload", "--abi", "./contract/simplestorage/SimpleStorage.abi", "--method", "set",
		"--args", "10")
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_Payload_ValidParameter returns error with valid parameter %s", err.Error())
	}
}

func Test_Client_Payload_InvalidParameter(t *testing.T) {
	cmd := exec.Command(CmdClient, "payload", "./contract/simplestorage/SimpleStorage.abi", "--method", "set",
		"--args", "10")
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatal("Test_Client_Payload_InvalidParameter returns ok with invalid parameter")
	}
}

func Test_Client_Payload_Method_InvalidParameter(t *testing.T) {
	cmd := exec.Command(CmdClient, "payload", "--abi", "./contract/simplestorage/SimpleStorage.abi", "--method", "get",
		"--args", "10")
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatal("Test_Client_Payload_Method_InvalidParameter returns ok with method invalid parameter")
	}
}

/*
func Test_Client_Miner_Status(t *testing.T) {
	cmd := exec.Command(CmdClient, "miner", "status", "--address", ServerAddr)
	status, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Test_Client_Miner_Status: An error occured: %s", err.Error())
	}
	if string(status) == "Running\n" {
		cmd = exec.Command(CmdClient, "miner", "stop", "--address", ServerAddr)
		if _, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("Stop mining returns error %s", err.Error())
		}
		cmd = exec.Command(CmdClient, "miner", "status", "--address", ServerAddr)
		status, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Test_Client_Miner_Status: An error occured: %s", err.Error())
		}
		if string(status) != "Stopped\n" {
			t.Fatal("Test_Client_Miner_Status returns error status")
		}
	} else if string(status) == "Stopped\n" {
		cmd := exec.Command(CmdClient, "miner", "start", "--address", ServerAddr)
		if _, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("Start mining returns error %s", err.Error())
		}
		cmd = exec.Command(CmdClient, "miner", "status", "--address", ServerAddr)
		status, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Test_Client_Miner_Status: An error occured: %s", err.Error())
		}
		if string(status) != "Running\n" {
			t.Fatal("Test_Client_Miner_Status returns error status")
		}
	} else {
		t.Fatalf("Test_Client_Miner_Status return error status %s", string(status))
	}
}

func Test_Client_Miner_Start_Multiply(t *testing.T) {
	cmd := exec.Command(CmdClient, "miner", "status", "--address", ServerAddr)
	status, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Get miner status returns error %s", err.Error())
	}
	if string(status) != "Running\n" {
		cmd := exec.Command(CmdClient, "miner", "start", "--threads", "3", "--address", ServerAddr)
		if _, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("Start mining returns error %s", err.Error())
		}
	}
	cmd = exec.Command(CmdClient, "miner", "start", "--threads", "3", "--address", ServerAddr)
	if _, err = cmd.CombinedOutput(); err == nil {
		t.Fatal("Test_Client_Miner_Start_Multiply returns ok")
	}
}

func Test_Client_Miner_Start_Invalid_Threads(t *testing.T) {
	cmd := exec.Command(CmdClient, "miner", "status", "--address", ServerAddr)
	status, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Get miner status returns error %s", err.Error())
	}
	if string(status) == "Running\n" {
		cmd = exec.Command(CmdClient, "miner", "stop", "--address", ServerAddr)
		if _, err = cmd.CombinedOutput(); err != nil {
			t.Fatalf("stop mining failed %s", err.Error())
		}
	}
	cmd = exec.Command(CmdClient, "miner", "start", "--threads", "-1", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatal("Test_Client_Miner_Start_Invalid_Threads returns ok")
	}
}

func Test_Client_Miner_Start_Valid_Threads(t *testing.T) {
	cmd := exec.Command(CmdClient, "miner", "status", "--address", ServerAddr)
	status, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Get miner status returns error %s", err.Error())
	}
	if string(status) == "Running\n" {
		cmd = exec.Command(CmdClient, "miner", "stop", "--address", ServerAddr)
		if _, err = cmd.CombinedOutput(); err != nil {
			t.Fatalf("stop mining failed %s", err.Error())
		}
	}
	cmd = exec.Command(CmdClient, "miner", "start", "--threads", "2", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatal("Test_Client_Miner_Start_Valid_Threads returns ok")
	}
	cmd = exec.Command(CmdClient, "miner", "threads", "--address", ServerAddr)
	n, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Test_Client_Miner_Start_Valid_Threads: An error occured: %s", err.Error())
	}
	if string(n) != "2\n" {
		t.Fatal("Test_Client_Miner_Start_Valid_Threads did not set the threads number")
	}
}

func Test_Client_Miner_Start_Default_Threads(t *testing.T) {
	cmd := exec.Command(CmdClient, "miner", "status", "--address", ServerAddr)
	status, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Get miner status returns error %s", err.Error())
	}
	if string(status) == "Running\n" {
		cmd = exec.Command(CmdClient, "miner", "stop", "--address", ServerAddr)
		if _, err = cmd.CombinedOutput(); err != nil {
			t.Fatalf("stop mining failed %s", err.Error())
		}
	}
	cmd = exec.Command(CmdClient, "miner", "start", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_Miner_Start_Default_Threads: An error occured: %s", err.Error())
	}
	cmd = exec.Command(CmdClient, "miner", "threads", "--address", ServerAddr)
	threads, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Test_Client_Miner_Start_Default_Threads: An error occured: %s", err.Error())
	}
	threads = bytes.TrimRight(threads, "\n")
	n, err := strconv.ParseInt(string(threads), 10, 64)
	if err != nil {
		t.Fatalf("Parse string to int failed %s", err.Error())
	}
	cpuNum := runtime.NumCPU()
	if int(n) != cpuNum {
		t.Fatal("Test_Client_Miner_Start_Default_Threads did not set default cpu nunber as threads number")
	}
}

func Test_Client_Miner_Stop_Multiply(t *testing.T) {
	cmd := exec.Command(CmdClient, "miner", "status", "--address", ServerAddr)
	status, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Get miner status returns error %s", err.Error())
	}
	if string(status) != "Stopped\n" {
		cmd := exec.Command(CmdClient, "miner", "stop", "--address", ServerAddr)
		if _, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("Stop mining returns error %s", err.Error())
		}
	}
	cmd = exec.Command(CmdClient, "miner", "stop", "--address", ServerAddr)
	if _, err = cmd.CombinedOutput(); err == nil {
		t.Fatal("Test_Client_Miner_Stop_Multiply returns ok")
	}
}

func Test_Client_Miner_Getcoinbase(t *testing.T) {
	cmd := exec.Command(CmdClient, "miner", "getcoinbase", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_Miner_Getcoinbase returns error %s", err.Error())
	}
}

func Test_Client_Miner_Hashrate(t *testing.T) {
	cmd := exec.Command(CmdClient, "miner", "hashrate", "--address", ServerAddr)
	result, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Test_Client_Miner_Hashrate: An error occured: %s", err.Error())
	}
	result = bytes.TrimRight(result, "\n")
	hashrate, err := strconv.ParseUint(string(result), 10, 64)
	if err != nil {
		t.Fatalf("Parse string to uint failed %s", err.Error())
	}
	if big.NewInt(int64(hashrate)).Cmp(big.NewInt(0)) < 0 {
		t.Fatal("Test_Client_Miner_Hashrate returns invalid result")
	}
}

func Test_Client_Miner_Setcoinbase_Valid(t *testing.T) {
	cmd := exec.Command(CmdClient, "miner", "setcoinbase", "--coinbase", AccountShard1_1, "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_Miner_Setcoinbase_Valid: An error occured: %s", err.Error())
	}
	cmd = exec.Command(CmdClient, "miner", "getcoinbase", "--address", ServerAddr)
	account, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Get coinbase error %s", err.Error())
	}
	account = bytes.TrimRight(account, "\n")
	if string(account) != AccountShard1_1 {
		t.Fatal("Test_Client_Miner_Setcoinbase_Valid did not set the coinbase successfully")
	}

	cmd = exec.Command(CmdClient, "miner", "status", "--address", ServerAddr)
	status, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Get miner status returns error %s", err.Error())
	}
	if string(status) == "Running\n" {
		cmd = exec.Command(CmdClient, "miner", "stop", "--address", ServerAddr)
		if _, err = cmd.CombinedOutput(); err != nil {
			t.Fatalf("stop mining failed %s", err.Error())
		}
	} else {
		cmd = exec.Command(CmdClient, "miner", "start", "--address", ServerAddr)
		if _, err = cmd.CombinedOutput(); err != nil {
			t.Fatalf("start mining failed %s", err.Error())
		}
	}
	cmd = exec.Command(CmdClient, "miner", "setcoinbase", "--coinbase", AccountShard1_2, "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_Miner_Setcoinbase_Valid: An error occured: %s", err.Error())
	}
	cmd = exec.Command(CmdClient, "miner", "getcoinbase", "--address", ServerAddr)
	account, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Get coinbase error %s", err.Error())
	}
	account = bytes.TrimRight(account, "\n")
	if string(account) != AccountShard1_2 {
		t.Fatal("Test_Client_Miner_Setcoinbase_Valid did not set the coinbase successfully")
	}
}

func Test_Client_Miner_Setcoinbase_InvalidParameter(t *testing.T) {
	cmd := exec.Command(CmdClient, "miner", "setcoinbase", AccountShard1_1, "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatal("Test_Client_Miner_Setcoinbase_InvalidParameter returns ok")
	}
}

func Test_Client_Miner_Setcoinbase_InvalidAccount(t *testing.T) {
	cmd := exec.Command(CmdClient, "miner", "setcoinbase", "--coinbase", AccountErr, "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatal("Test_Client_Miner_Setcoinbase_InvalidAccount return ok")
	}
}

func Test_Client_Miner_Setcoinbase_InvalidAccountType(t *testing.T) {
	cmd := exec.Command(CmdClient, "miner", "setcoinbase", "--coinbase", InvalidAccountType, "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatal("Test_Client_Miner_Setcoinbase_InvalidAccountType return ok")
	}
}

func Test_Client_Miner_Setcoinbase_AccountFromOtherShard(t *testing.T) {
	cmd := exec.Command(CmdClient, "miner", "setcoinbase", "--coinbase", AccountShard2_1, "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatal("Test_Client_Miner_Setcoinbase_AccountFromOtherShard return ok")
	}
}

func Test_Client_Miner_Threads(t *testing.T) {
	cmd := exec.Command(CmdClient, "miner", "threads", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_Miner_Threads: An error occured: %s", err.Error())
	}
	cmd = exec.Command(CmdClient, "miner", "status", "--address", ServerAddr)
	status, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Get miner status returns error %s", err.Error())
	}
	if string(status) == "Running\n" {
		cmd = exec.Command(CmdClient, "miner", "stop", "--address", ServerAddr)
		if _, err = cmd.CombinedOutput(); err != nil {
			t.Fatalf("stop mining failed %s", err.Error())
		}
	} else {
		cmd = exec.Command(CmdClient, "miner", "start", "--address", ServerAddr)
		if _, err = cmd.CombinedOutput(); err != nil {
			t.Fatalf("start mining failed %s", err.Error())
		}
	}
	cmd = exec.Command(CmdClient, "miner", "threads", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_Miner_Threads: An error occured: %s", err.Error())
	}
}

func Test_Client_Miner_SetThreads(t *testing.T) {
	cmd := exec.Command(CmdClient, "miner", "setthreads", "--threads", "10", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_Miner_SetThreads: An error occured: %s", err.Error())
	}
	cmd = exec.Command(CmdClient, "miner", "threads", "--address", ServerAddr)
	n1, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Test_Client_Miner_SetThreads: An error occured: %s", err.Error())
	}
	if string(n1) != "10\n" {
		t.Fatal("Test_Client_Miner_SetThreads did not set the threads number")
	}

	cmd = exec.Command(CmdClient, "miner", "status", "--address", ServerAddr)
	status, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("get miner status returns error %s", err.Error())
	}
	if string(status) == "Running\n" {
		cmd = exec.Command(CmdClient, "miner", "stop", "--address", ServerAddr)
		if _, err = cmd.CombinedOutput(); err != nil {
			t.Fatalf("stop mining failed %s", err.Error())
		}
	} else {
		cmd = exec.Command(CmdClient, "miner", "start", "--address", ServerAddr)
		if _, err = cmd.CombinedOutput(); err != nil {
			t.Fatalf("start mining failed %s", err.Error())
		}
	}
	cmd = exec.Command(CmdClient, "miner", "setthreads", "--threads", "5", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_Miner_SetThreads: An error occured: %s", err.Error())
	}
	cmd = exec.Command(CmdClient, "miner", "threads", "--address", ServerAddr)
	n2, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Test_Client_Miner_SetThreads: An error occured: %s", err.Error())
	}
	if string(n2) != "5\n" {
		t.Fatal("Test_Client_Miner_SetThreads did not set the threads number")
	}
}

func Test_Client_Miner_SetThreads_Default(t *testing.T) {
	cmd := exec.Command(CmdClient, "miner", "setthreads", "--threads", "10", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_Miner_SetThreads_Default: An error occured: %s", err.Error())
	}
	cmd = exec.Command(CmdClient, "miner", "setthreads", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_Miner_SetThreads_Default: An error occured: %s", err.Error())
	}
	cmd = exec.Command(CmdClient, "miner", "threads", "--address", ServerAddr)
	threads, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Test_Client_Miner_SetThreads_Default: An error occured: %s", err.Error())
	}
	threads = bytes.TrimRight(threads, "\n")
	n, err := strconv.ParseInt(string(threads), 10, 64)
	if err != nil {
		t.Fatalf("Parse string to int failed %s", err.Error())
	}
	cpuNum := runtime.NumCPU()
	if int(n) != cpuNum {
		t.Fatal("Test_Client_Miner_SetThreads_Default did not set default cpu nunber as threads number")
	}
}*/

// --------------------test savekey start-------------------
func Test_Client_SaveKey_Invalid_Privatekey_Without_Prefix_0x(t *testing.T) {
	cmd := exec.Command(CmdClient, "savekey", "--privatekey", "123")
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_SaveKey_Invalid_Privatekey_Without_Prefix_0x,savekey  should return error with privatekey without prefix 0x")
	}
}

func Test_Client_SaveKey_Invalid_Privatekey_With_Prefix_Odd(t *testing.T) {
	cmd := exec.Command(CmdClient, "savekey", "--privatekey", "0x123")
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_SaveKey_Invalid_Privatekey_With_Prefix_Odd,savekey should return error with privatekey is odd length")
	}
}

func Test_Client_SaveKey_Invalid_Privatekey_Syntax_Characeter(t *testing.T) {
	cmd := exec.Command(CmdClient, "savekey", "--privatekey", "0x1234-")
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_SaveKey_Invalid_Privatekey_Syntax_Characeter,savekey should return error with privatekey has syntax character")
	}
}

func Test_Client_SaveKey_Invalid_FileNameValue_Empty(t *testing.T) {
	cmd := exec.Command(CmdClient, "savekey", "--privatekey", AccountPrivateKey2, "--file", "")
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_SaveKey_Invalid_FileNameValue_Empty,savekey should return error with empty filename")
	}
}

func Test_Client_SaveKey_Invalid_Privatekey_With_Invalid_length(t *testing.T) {
	cmd := exec.Command(CmdClient, "savekey", "--privatekey", "0x")
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_SaveKey_Invalid_Privatekey_With_Invalid_length,savekey  should return error with privatekey of invalid length(less than 256 bits)")
	}
}

func Test_Client_SaveKey(t *testing.T) {
	cmd := exec.Command(CmdClient, "savekey", "--privatekey", AccountPrivateKey2, "--file", ".test_keystore")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_Client_SaveKey: An error occured: %s", err)
	}
	// need input twice
	io.WriteString(stdin, "12345\n")
	io.WriteString(stdin, "12345\n")
	cmd.Wait()

	outStr, errStr := out.String(), outErr.String()
	fmt.Println(outStr, errStr)
	if len(string(errStr)) > 0 {
		t.Fatalf("Test_Client_SaveKey savekey error. %s %s", errStr, cmd.Args)
	}
}

// --------------------test savekey end-------------------

func Test_Client_GetBlockHeight_NodeStop(t *testing.T) {
	cmd := exec.Command(CmdClient, "getblockheight", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_GetBlockHeight error, %s", err)
	}
}

func Test_Client_GetBlockHeight_NodeStart(t *testing.T) {
	cmd := exec.Command(CmdClient, "getblockheight", "--address", ServerAddr)
	if res, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_GetBlockHeight error, %s", err)
	} else {
		res = bytes.TrimRight(res, "\n")
		height, _ := strconv.ParseInt(string(res), 10, 64)
		if height < 0 {
			t.Fatalf("Test_Client_GetBlockHeight_NodeStart: The return value is not correct!")
		}
	}
}

func Test_Client_GetBlockHeight_InvalidParameter(t *testing.T) {
	cmd := exec.Command(CmdClient, "getblockheight", "--height", "1000000000", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_GetBlockHeight_InvalidParameter returns ok with invalid parameter")
	}
}

func Test_Client_GetBlockHeight_ByInvalidHeight(t *testing.T) {
	cmd := exec.Command(CmdClient, "getblockheight", "--height", "1", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_GetBlockHeight_ByInvalidHeight returns error not defined: -height")
	}
}

func Test_Client_GetBlockHeight_Parameter(t *testing.T) {
	cmd := exec.Command(CmdClient, "getblockheight", "1", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_GetBlockHeight_Parameter error, %s", err)
	}
}

func Test_Client_GetBlockTXCount_ByInvalidHeight(t *testing.T) {
	cmd := exec.Command(CmdClient, "getblocktxcount", "--height", "100000000", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_GetBlockTXCount_ByInvalidHeight error parameter success?")
	}
}

func Test_Client_GetBlockTXCount_ByInvalidHeight0x(t *testing.T) {
	cmd := exec.Command(CmdClient, "getblocktxcount", "--height", "0x", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_GetBlockTXCount_ByInvalidHeight0x return error invalid value")
	}
}

func Test_Client_GetBlockTXCount_ByHeight_NodeStart(t *testing.T) {
	cmd := exec.Command(CmdClient, "getblocktxcount", "--height", "1", "--address", ServerAddr)
	if res, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_GetBlockTXCount_ByHeight: error, %s", err)
	} else {
		res = bytes.TrimRight(res, "\n")
		TXCount, _ := strconv.ParseInt(string(res), 10, 64)
		if TXCount < 0 {
			t.Fatalf("Test_Client_GetBlockHeight_NodeStart: The return value is not correct!")
		}
	}
}

func Test_Client_GetBlockTXCount_DefaultParameter(t *testing.T) {
	cmd := exec.Command(CmdClient, "getblocktxcount", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_GetBlockTXCount_DefaultParameter:error, %s", err)
	}
}

func Test_Client_GetBlockTXCount_InvalidParameter(t *testing.T) {
	cmd := exec.Command(CmdClient, "getblocktxcount", "1", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_GetBlockTXCount_InvalidParameter： error parameter success?")
	}
}

func Test_Client_GetBlockTXCount_ByInvalidHash(t *testing.T) {
	cmd := exec.Command(CmdClient, "getblocktxcount", "--hash", BlockHashErr, "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_GetBlockTXCount_ByInvalidHash： error parameter success?")
	}
}

func Test_Client_GetBlockTXCount_ByHash(t *testing.T) {
	cmd := exec.Command(CmdClient, "getblocktxcount", "--hash", BlockHash, "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_GetBlockTXCount_ByHash: getblocktxcount error, %s", err)
	}
}

func Test_Client_GetBlockTXCount_ByInvalidHash0x(t *testing.T) {
	cmd := exec.Command(CmdClient, "getblocktxcount", "--hash", "0x", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Light_GetBlockTXCount_ByInvalidHash0x error parameter success?")
	}
}

func Test_Client_GetBlockTXCount_ByInvalidHash0x12(t *testing.T) {
	cmd := exec.Command(CmdClient, "getblocktxcount", "--hash", "0x12-", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_GetBlockTXCount_ByInvalidHash0x12： return error syntax character")
	}
}

func Test_Client_GetBlockTXCount_ByInvalidHash123(t *testing.T) {
	cmd := exec.Command(CmdClient, "getblocktxcount", "--hash", "123", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_GetBlockTXCount_ByInvalidHash123： return error hex string without 0x prefix")
	}
}

func Test_Client_GetBlock_ByHeight_NodeStop(t *testing.T) {
	// Normal height
	cmd := exec.Command(CmdClient, "getblock", "--height", "1", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_GetBlock_ByHeight_NodeStop: Node not running!")
	}
}

func Test_Client_GetBlock_ByHeight_NodeStart(t *testing.T) {
	// Normal height
	cmd := exec.Command(CmdClient, "getblock", "--height", "1", "--address", ServerAddr)
	if res, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_GetBlock_ByHeight_NodeStart:Node to run returns error: %s", err)
	} else {
		var blockInfo BlockInfo
		if err = json.Unmarshal(res, &blockInfo); err != nil {
			t.Fatalf("Test_Client_GetBlock_ByHeight_NodeStart: %s", err)
		}

		height := blockInfo.Header.Height
		if height != 1 {
			t.Fatalf("Test_Client_GetBlock_ByHeight_NodeStart: Expect the return value is not correct!")
		}

	}
}

func Test_Client_GetBlock_ByInvalidHeight(t *testing.T) {
	// invalid height
	cmd := exec.Command(CmdClient, "getblock", "--height", "10000000000", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_GetBlock_ByInvalidHeight: error parameter success?")
	}
}

func Test_Client_GetBlock_InvalidParameter(t *testing.T) {
	cmd := exec.Command(CmdClient, "getblock", "1", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_GetBlock_InvalidParameter error, %s", err)
	}
}

func Test_Client_GetBlock_ByNormalHash(t *testing.T) {
	// invalid height
	cmd := exec.Command(CmdClient, "getblock", "--hash", BlockHash, "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_GetBlock_ByNormalHash error, %s", err)
	}
}

func Test_Client_GetBlock_ByInvalidHash(t *testing.T) {
	// invalid height
	cmd := exec.Command(CmdClient, "getblock", "--hash", BlockHashErr, "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_GetBlock_ByInvalidHash error parameter success?")
	}
}

// getblock fulltx support.
func Test_Client_GetBlock_ByHeightFulltx(t *testing.T) {
	cmd := exec.Command(CmdClient, "getblock", "--height", "1", "--fulltx", "--address", ServerAddr)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_GetBlock_ByHeightFulltx error, %s", err)
	} else {
		var blockInfo BlockInfo
		if err = json.Unmarshal(output, &blockInfo); err != nil {
			t.Fatalf("Test_Client_GetBlock_ByHeightFulltx: %s", err)
		}

		height := blockInfo.Header.Height
		if height != 1 {
			t.Fatalf("Test_Client_GetBlock_ByHeightFulltx: Expect the return value is not correct!")
		}
		if len(blockInfo.Transactions) <= 0 {
			t.Fatalf("Test_Client_GetBlock_ByHeightFulltx, block should contain one transaction at lease")
		}
	}
}

// getblock fulltx support.
func Test_Client_GetBlock_ByHashFulltx(t *testing.T) {
	cmd := exec.Command(CmdClient, "getblock", "--hash", BlockHash, "--fulltx", "--address", ServerAddr)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_GetBlock_ByHashFulltx error, %s", err)
	} else {
		var blockInfo BlockInfo
		if err = json.Unmarshal(output, &blockInfo); err != nil {
			t.Fatalf("Test_Client_GetBlock_ByHashFulltx: %s", err)
		}
		if blockInfo.Hash != BlockHash {
			t.Fatalf("Test_Client_GetBlock_ByHashFulltx: Expect the return value is not correct!")
		}
		if len(blockInfo.Transactions) <= 0 {
			t.Fatalf("Test_Client_GetBlock_ByHashFulltx, hash should contain one transaction at lease")
		}
	}
}
