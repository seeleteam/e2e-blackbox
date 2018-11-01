/**
*  @file
*  @copyright defined in go-seele/LICENSE
 */

package testcase

import (
	"bytes"
	"encoding/json"
	"math/big"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
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
		t.Fatal("Test_Client_Dumpheap_Default_Filename get user path failed %s", err.Error())
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
		t.Fatal("Test_Client_Dumpheap_Specified_Filename: Get user path failed %s", err.Error())
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
		t.Fatalf("Test_Client_Payload_InvalidParameter returns ok with invalid parameter")
	}
}

func Test_Client_Payload_Method_InvalidParameter(t *testing.T) {
	cmd := exec.Command(CmdClient, "payload", "--abi", "./contract/simplestorage/SimpleStorage.abi", "--method", "get",
		"--args", "10")
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_Payload_Method_InvalidParameter returns ok with method invalid parameter")
	}
}

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
			t.Fatalf("Test_Client_Miner_Status returns error status")
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
			t.Fatalf("Test_Client_Miner_Status returns error status")
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
		t.Fatalf("Test_Client_Miner_Start_Multiply returns ok")
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
		t.Fatalf("Test_Client_Miner_Start_Invalid_Threads returns ok")
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
		t.Fatalf("Test_Client_Miner_Start_Valid_Threads returns ok")
	}
	cmd = exec.Command(CmdClient, "miner", "threads", "--address", ServerAddr)
	n, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Test_Client_Miner_Start_Valid_Threads: An error occured: %s", err.Error())
	}
	if string(n) != "2\n" {
		t.Fatalf("Test_Client_Miner_Start_Valid_Threads did not set the threads number")
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
		t.Fatalf("Test_Client_Miner_Start_Default_Threads did not set default cpu nunber as threads number")
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
		t.Fatalf("Test_Client_Miner_Stop_Multiply returns ok")
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
		t.Fatalf("Test_Client_Miner_Setcoinbase_Valid did not set the coinbase successfully")
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
		t.Fatalf("Test_Client_Miner_Setcoinbase_Valid did not set the coinbase successfully")
	}
}

func Test_Client_Miner_Setcoinbase_InvalidParameter(t *testing.T) {
	cmd := exec.Command(CmdClient, "miner", "setcoinbase", AccountShard1_1, "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_Miner_Setcoinbase_InvalidParameter returns ok")
	}
}

func Test_Client_Miner_Setcoinbase_InvalidAccount(t *testing.T) {
	cmd := exec.Command(CmdClient, "miner", "setcoinbase", "--coinbase", AccountErr, "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_Miner_Setcoinbase_InvalidAccount return ok")
	}
}

func Test_Client_Miner_Setcoinbase_InvalidAccountType(t *testing.T) {
	cmd := exec.Command(CmdClient, "miner", "setcoinbase", "--coinbase", InvalidAccountType, "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_Miner_Setcoinbase_InvalidAccountType return ok")
	}
}

func Test_Client_Miner_Setcoinbase_AccountFromOtherShard(t *testing.T) {
	cmd := exec.Command(CmdClient, "miner", "setcoinbase", "--coinbase", AccountShard2_1, "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_Miner_Setcoinbase_AccountFromOtherShard return ok")
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
		t.Fatalf("Test_Client_Miner_SetThreads did not set the threads number")
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
		t.Fatalf("Test_Client_Miner_SetThreads did not set the threads number")
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
		t.Fatalf("Test_Client_Miner_SetThreads_Default did not set default cpu nunber as threads number")
	}
}
