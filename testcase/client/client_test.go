/**
*  @file
*  @copyright defined in go-seele/LICENSE
 */

package client

import (
	"bytes"
	"encoding/json"
	"errors"
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
	"time"

	"github.com/seeleteam/e2e-blackbox/testcase/common"
)

func Test_Client_GetInfo(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getinfo")
	res, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Test_Client_GetInfo: GetInfo error, %s", err)
	}

	var r common.ResGetInfo
	err = json.Unmarshal(res, &r)
	if err != nil {
		t.Fatalf("Test_Client_GetInfo: decode return result error %s", err)
	}

	if r.MinerStatus != "Running" {
		t.Fatalf("Test_Client_GetInfo: Node not running!")
	}
}

func Test_Client_Key(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "key")
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
	cmd := exec.Command(common.CmdClient, "dumpheap")
	res, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Test_Client_DumpHeap: dumpheap error, %s", err)
	}

	if _, err = os.Stat(strings.TrimSpace(string(res))); os.IsNotExist(err) {
		t.Fatalf("Test_Client_DumpHeap:  File %s not found!", string(res))
	}
}

func Test_Client_Dumpheap_Default_Filename(t *testing.T) {
	userPath, err := user.Current()
	if err != nil {
		t.Fatalf("Test_Client_Dumpheap_Default_Filename get user path failed %s", err.Error())
	}
	defaultDataFolder := filepath.Join(userPath.HomeDir, ".seele")
	defaultFilePath := filepath.Join(defaultDataFolder, "heap.dump\n")

	cmd := exec.Command(common.CmdClient, "dumpheap", "--address", common.ServerAddr)
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

	cmd := exec.Command(common.CmdClient, "dumpheap", "--address", common.ServerAddr, "--file", "test.dump")
	file, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Test_Client_Dumpheap_Specified_Filename: An error occured: %s", err.Error())
	}

	if defaultFilePath != string(file) {
		t.Fatal("Test_Client_Dumpheap_Specified_Filename: The actual dumpheapPath is not equal to the expected path")
	}
}

func Test_Client_Payload_ValidParameter(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "payload", "--abi", "../contract/simplestorage/SimpleStorage.abi", "--method", "set",
		"--args", "10")
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_Payload_ValidParameter returns error with valid parameter %s", err.Error())
	}
}

func Test_Client_Payload_InvalidParameter(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "payload", "./contract/simplestorage/SimpleStorage.abi", "--method", "set",
		"--args", "10")
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatal("Test_Client_Payload_InvalidParameter returns ok with invalid parameter")
	}
}

func Test_Client_Payload_Method_InvalidParameter(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "payload", "--abi", "./contract/simplestorage/SimpleStorage.abi", "--method", "get",
		"--args", "10")
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatal("Test_Client_Payload_Method_InvalidParameter returns ok with method invalid parameter")
	}
}

/*
func Test_Client_Miner_Status(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "miner", "status", "--address", common.ServerAddr)
	status, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Test_Client_Miner_Status: An error occured: %s", err.Error())
	}
	if string(status) == "Running\n" {
		cmd = exec.Command(common.CmdClient, "miner", "stop", "--address", common.ServerAddr)
		if _, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("Stop mining returns error %s", err.Error())
		}
		cmd = exec.Command(common.CmdClient, "miner", "status", "--address", common.ServerAddr)
		status, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Test_Client_Miner_Status: An error occured: %s", err.Error())
		}
		if string(status) != "Stopped\n" {
			t.Fatal("Test_Client_Miner_Status returns error status")
		}
	} else if string(status) == "Stopped\n" {
		cmd := exec.Command(common.CmdClient, "miner", "start", "--address", common.ServerAddr)
		if _, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("Start mining returns error %s", err.Error())
		}
		cmd = exec.Command(common.CmdClient, "miner", "status", "--address", common.ServerAddr)
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
	cmd := exec.Command(common.CmdClient, "miner", "status", "--address", common.ServerAddr)
	status, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Get miner status returns error %s", err.Error())
	}
	if string(status) != "Running\n" {
		cmd := exec.Command(common.CmdClient, "miner", "start", "--threads", "3", "--address", common.ServerAddr)
		if _, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("Start mining returns error %s", err.Error())
		}
	}
	cmd = exec.Command(common.CmdClient, "miner", "start", "--threads", "3", "--address", common.ServerAddr)
	if _, err = cmd.CombinedOutput(); err == nil {
		t.Fatal("Test_Client_Miner_Start_Multiply returns ok")
	}
}

func Test_Client_Miner_Start_Invalid_Threads(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "miner", "status", "--address", common.ServerAddr)
	status, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Get miner status returns error %s", err.Error())
	}
	if string(status) == "Running\n" {
		cmd = exec.Command(common.CmdClient, "miner", "stop", "--address", common.ServerAddr)
		if _, err = cmd.CombinedOutput(); err != nil {
			t.Fatalf("stop mining failed %s", err.Error())
		}
	}
	cmd = exec.Command(common.CmdClient, "miner", "start", "--threads", "-1", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatal("Test_Client_Miner_Start_Invalid_Threads returns ok")
	}
}

func Test_Client_Miner_Start_Valid_Threads(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "miner", "status", "--address", common.ServerAddr)
	status, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Get miner status returns error %s", err.Error())
	}
	if string(status) == "Running\n" {
		cmd = exec.Command(common.CmdClient, "miner", "stop", "--address", common.ServerAddr)
		if _, err = cmd.CombinedOutput(); err != nil {
			t.Fatalf("stop mining failed %s", err.Error())
		}
	}
	cmd = exec.Command(common.CmdClient, "miner", "start", "--threads", "2", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatal("Test_Client_Miner_Start_Valid_Threads returns ok")
	}
	cmd = exec.Command(common.CmdClient, "miner", "threads", "--address", common.ServerAddr)
	n, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Test_Client_Miner_Start_Valid_Threads: An error occured: %s", err.Error())
	}
	if string(n) != "2\n" {
		t.Fatal("Test_Client_Miner_Start_Valid_Threads did not set the threads number")
	}
}

func Test_Client_Miner_Start_Default_Threads(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "miner", "status", "--address", common.ServerAddr)
	status, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Get miner status returns error %s", err.Error())
	}
	if string(status) == "Running\n" {
		cmd = exec.Command(common.CmdClient, "miner", "stop", "--address", common.ServerAddr)
		if _, err = cmd.CombinedOutput(); err != nil {
			t.Fatalf("stop mining failed %s", err.Error())
		}
	}
	cmd = exec.Command(common.CmdClient, "miner", "start", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_Miner_Start_Default_Threads: An error occured: %s", err.Error())
	}
	cmd = exec.Command(common.CmdClient, "miner", "threads", "--address", common.ServerAddr)
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
	cmd := exec.Command(common.CmdClient, "miner", "status", "--address", common.ServerAddr)
	status, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Get miner status returns error %s", err.Error())
	}
	if string(status) != "Stopped\n" {
		cmd := exec.Command(common.CmdClient, "miner", "stop", "--address", common.ServerAddr)
		if _, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("Stop mining returns error %s", err.Error())
		}
	}
	cmd = exec.Command(common.CmdClient, "miner", "stop", "--address", common.ServerAddr)
	if _, err = cmd.CombinedOutput(); err == nil {
		t.Fatal("Test_Client_Miner_Stop_Multiply returns ok")
	}
}

func Test_Client_Miner_Getcoinbase(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "miner", "getcoinbase", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_Miner_Getcoinbase returns error %s", err.Error())
	}
}

func Test_Client_Miner_Hashrate(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "miner", "hashrate", "--address", common.ServerAddr)
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
	cmd := exec.Command(common.CmdClient, "miner", "setcoinbase", "--coinbase", common.AccountShard1_1, "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_Miner_Setcoinbase_Valid: An error occured: %s", err.Error())
	}
	cmd = exec.Command(common.CmdClient, "miner", "getcoinbase", "--address", common.ServerAddr)
	account, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Get coinbase error %s", err.Error())
	}
	account = bytes.TrimRight(account, "\n")
	if string(account) != common.AccountShard1_1 {
		t.Fatal("Test_Client_Miner_Setcoinbase_Valid did not set the coinbase successfully")
	}

	cmd = exec.Command(common.CmdClient, "miner", "status", "--address", common.ServerAddr)
	status, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Get miner status returns error %s", err.Error())
	}
	if string(status) == "Running\n" {
		cmd = exec.Command(common.CmdClient, "miner", "stop", "--address", common.ServerAddr)
		if _, err = cmd.CombinedOutput(); err != nil {
			t.Fatalf("stop mining failed %s", err.Error())
		}
	} else {
		cmd = exec.Command(common.CmdClient, "miner", "start", "--address", common.ServerAddr)
		if _, err = cmd.CombinedOutput(); err != nil {
			t.Fatalf("start mining failed %s", err.Error())
		}
	}
	cmd = exec.Command(common.CmdClient, "miner", "setcoinbase", "--coinbase", AccountShard1_2, "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_Miner_Setcoinbase_Valid: An error occured: %s", err.Error())
	}
	cmd = exec.Command(common.CmdClient, "miner", "getcoinbase", "--address", common.ServerAddr)
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
	cmd := exec.Command(common.CmdClient, "miner", "setcoinbase", common.AccountShard1_1, "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatal("Test_Client_Miner_Setcoinbase_InvalidParameter returns ok")
	}
}

func Test_Client_Miner_Setcoinbase_InvalidAccount(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "miner", "setcoinbase", "--coinbase", common.AccountErr, "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatal("Test_Client_Miner_Setcoinbase_InvalidAccount return ok")
	}
}

func Test_Client_Miner_Setcoinbase_testcase.InvalidAccountType(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "miner", "setcoinbase", "--coinbase", common.InvalidAccountType, "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatal("Test_Client_Miner_Setcoinbase_testcase.InvalidAccountType return ok")
	}
}

func Test_Client_Miner_Setcoinbase_AccountFromOtherShard(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "miner", "setcoinbase", "--coinbase", AccountShard2_1, "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatal("Test_Client_Miner_Setcoinbase_AccountFromOtherShard return ok")
	}
}

func Test_Client_Miner_Threads(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "miner", "threads", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_Miner_Threads: An error occured: %s", err.Error())
	}
	cmd = exec.Command(common.CmdClient, "miner", "status", "--address", common.ServerAddr)
	status, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Get miner status returns error %s", err.Error())
	}
	if string(status) == "Running\n" {
		cmd = exec.Command(common.CmdClient, "miner", "stop", "--address", common.ServerAddr)
		if _, err = cmd.CombinedOutput(); err != nil {
			t.Fatalf("stop mining failed %s", err.Error())
		}
	} else {
		cmd = exec.Command(common.CmdClient, "miner", "start", "--address", common.ServerAddr)
		if _, err = cmd.CombinedOutput(); err != nil {
			t.Fatalf("start mining failed %s", err.Error())
		}
	}
	cmd = exec.Command(common.CmdClient, "miner", "threads", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_Miner_Threads: An error occured: %s", err.Error())
	}
}

func Test_Client_Miner_SetThreads(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "miner", "setthreads", "--threads", "10", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_Miner_SetThreads: An error occured: %s", err.Error())
	}
	cmd = exec.Command(common.CmdClient, "miner", "threads", "--address", common.ServerAddr)
	n1, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Test_Client_Miner_SetThreads: An error occured: %s", err.Error())
	}
	if string(n1) != "10\n" {
		t.Fatal("Test_Client_Miner_SetThreads did not set the threads number")
	}

	cmd = exec.Command(common.CmdClient, "miner", "status", "--address", common.ServerAddr)
	status, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("get miner status returns error %s", err.Error())
	}
	if string(status) == "Running\n" {
		cmd = exec.Command(common.CmdClient, "miner", "stop", "--address", common.ServerAddr)
		if _, err = cmd.CombinedOutput(); err != nil {
			t.Fatalf("stop mining failed %s", err.Error())
		}
	} else {
		cmd = exec.Command(common.CmdClient, "miner", "start", "--address", common.ServerAddr)
		if _, err = cmd.CombinedOutput(); err != nil {
			t.Fatalf("start mining failed %s", err.Error())
		}
	}
	cmd = exec.Command(common.CmdClient, "miner", "setthreads", "--threads", "5", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_Miner_SetThreads: An error occured: %s", err.Error())
	}
	cmd = exec.Command(common.CmdClient, "miner", "threads", "--address", common.ServerAddr)
	n2, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Test_Client_Miner_SetThreads: An error occured: %s", err.Error())
	}
	if string(n2) != "5\n" {
		t.Fatal("Test_Client_Miner_SetThreads did not set the threads number")
	}
}

func Test_Client_Miner_SetThreads_Default(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "miner", "setthreads", "--threads", "10", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_Miner_SetThreads_Default: An error occured: %s", err.Error())
	}
	cmd = exec.Command(common.CmdClient, "miner", "setthreads", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_Miner_SetThreads_Default: An error occured: %s", err.Error())
	}
	cmd = exec.Command(common.CmdClient, "miner", "threads", "--address", common.ServerAddr)
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
	cmd := exec.Command(common.CmdClient, "savekey", "--privatekey", "123")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_SaveKey_Invalid_Privatekey_Without_Prefix_0x: An error occured: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "Input string not a valid ecdsa string") {
		t.Fatalf("Test_Client_SaveKey_Invalid_Privatekey_Without_Prefix_0x,savekey  should return error with privatekey without prefix 0x: %s", errStr)
	}
}

func Test_Client_SaveKey_Invalid_Privatekey_With_Prefix_Odd(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "savekey", "--privatekey", "0x123")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_SaveKey_Invalid_Privatekey_With_Prefix_Odd: An error occured: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "encoding/hex: odd length hex string") {
		t.Fatalf("Test_Client_SaveKey_Invalid_Privatekey_With_Prefix_Odd,savekey should return error with privatekey is odd length: %s", errStr)
	}
}

func Test_Client_SaveKey_Invalid_Privatekey_Syntax_Characeter(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "savekey", "--privatekey", "0x12345-")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_SaveKey_Invalid_Privatekey_Syntax_Characeter: An error occured: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "encoding/hex: invalid byte") {
		t.Fatalf("Test_Client_SaveKey_Invalid_Privatekey_Syntax_Characeter,savekey should return error with privatekey has syntax character: %s", errStr)
	}
}

func Test_Client_SaveKey_Invalid_FileNameValue_Empty(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "savekey", "--privatekey", common.AccountPrivateKey2, "--file", "")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_SaveKey_Invalid_FileNameValue_Empty: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "please specify the key file path") {
		t.Fatalf("Test_Client_SaveKey_Invalid_FileNameValue_Empty,savekey should return error with empty filename: %s", errStr)
	}
}

func Test_Client_SaveKey_Invalid_Privatekey_With_Invalid_length(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "savekey", "--privatekey", "0x")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_SaveKey_Invalid_Privatekey_With_Invalid_length: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "invalid length, need 256 bits") {
		t.Fatalf("Test_Client_SaveKey_Invalid_Privatekey_With_Invalid_length,savekey  should return error with privatekey of invalid length(less than 256 bits): %s", errStr)
	}

}

func Test_Client_SaveKey(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "savekey", "--privatekey", common.AccountPrivateKey2, "--file", ".test_keystore")
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

// --------------------test getbalance start-------------------
func Test_Client_GetBalance_Account_Invalid_With_Prefix_Odd(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getbalance", "--account", common.AccountErr, "--address", common.ServerAddr)
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_GetBalance_Account_Invalid_With_Prefix_Odd: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "hex string of odd length") {
		t.Fatalf("Test_Client_GetBalance_Account_Invalid_With_Prefix_Odd,getbalance should return error with account of odd length: %s", errStr)
	}
}

func Test_Client_GetBalance_Account_Invalid_Without_Prefix_0x(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getbalance", "--account", "aaaaaaaaaaaaaaaaa", "--address", common.ServerAddr)
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_GetBalance_Account_Invalid_Without_Prefix_0x: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "hex string without 0x prefix") {
		t.Fatalf("Test_Client_GetBalance_Account_Invalid_Without_Prefix_0x,getbalance should return error with account without prefix 0x: %s", errStr)
	}
}

func Test_Client_GetBalance_Account_Invalid_Syntax_Characeter(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getbalance", "--account", "0xaaaaaaaaaaaaaaaaa-", "--address", common.ServerAddr)
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_GetBalance_Account_Invalid_Syntax_Characeter: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "invalid hex string") {
		t.Fatalf("Test_Client_GetBalance_Account_Invalid_Syntax_Characeter,getbalance should return error with account has syntax character: %s", errStr)
	}
}

func Test_Client_GetBalance_Account_Invalid_empty(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getbalance", "--account", "", "--address", common.ServerAddr)
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_GetBalance_Account_Invalid_empty: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "invalid account") {
		t.Fatalf("Test_Client_GetBalance_Account_Invalid_empty,getbalance should return error with empty account: %s", errStr)
	}
}

func Test_Client_GetBalance_Account_Invalid_FromOtherShard(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getbalance", "--account", common.Account2, "--address", common.ServerAddr)
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_GetBalance_Account_Invalid_FromOtherShard: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "local shard is: 1, your shard is: 2, you need to change to shard 2 to get your balance") {
		t.Fatalf("Test_Client_GetBalance_Account_Invalid_FromOtherShard,getbalance should return error with from other shard: %s", errStr)
	}
}

func Test_Client_GetBalance_Account(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getbalance", "--account", "0x0a57a2714e193b7ac50475ce625f2dcfb483d741", "--address", common.ServerAddr)
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_GetBalance_Account: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if errStr != "" {
		t.Fatalf("Test_Client_GetBalance_Account,getbalance should return error with empty account: %s", errStr)
	}
}

// --------------------test getbalance end-------------------

// --------------------test getshardnum start-------------------
/*func Test_Client_GetShardNum_Account_Invalid_With_Invalid_Type(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getshardnum", "--account", "0xff0fb1e59e92e94fac74febec98cfd58b956fa6d")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_GetShardNum_Account_Invalid_With_Invalid_Type: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()
	fmt.Println("err:", errStr)
	if !strings.Contains(errStr, "invalid account type") {
		t.Fatalf("Test_Client_GetShardNum_Account_Invalid_With_Invalid_Type,getshardnum should return error with invalid account type: %s", errStr)
	}
}*/

func Test_Client_GetShardNum_Account_Invalid_Without_Prefix_0x(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getshardnum", "--account", "123")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_GetShardNum_Account_Invalid_Without_Prefix_0x: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "the account is invalid for: hex string without 0x prefix") {
		t.Fatalf("Test_Client_GetShardNum_Account_Invalid_Without_Prefix_0x,getshardnum should return error with account without prefix 0x: %s", errStr)
	}
}

func Test_Client_GetShardNum_Account_Invalid_With_Prefix_0x_Odd(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getshardnum", "--account", "0x123")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_GetShardNum_Account_Invalid_With_Prefix_0x_Odd: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "the account is invalid for: hex string of odd length") {
		t.Fatalf("Test_Client_GetShardNum_Account_Invalid_With_Prefix_0x_Odd,getshardnum should return error with account of odd length: %s", errStr)
	}
}

func Test_Client_GetShardNum_Account_Invalid_Empty(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getshardnum", "--account", "")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_GetShardNum_Account_Invalid_Empty: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "the account is invalid for: empty hex string") {
		t.Fatalf("Test_Client_GetShardNum_Account_Invalid_Empty,getshardnum should return error with empty account: %s", errStr)
	}
}
func Test_Client_GetShardNum_Account_Invalid_Syntax_Character(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getshardnum", "--account", "0x12345-")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_GetShardNum_Account_Invalid_Syntax_Character: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "the account is invalid for: invalid hex string") {
		t.Fatalf("Test_Client_GetShardNum_Account_Invalid_Syntax_Character,getshardnum should return error with account has syntax character: %s", errStr)
	}
}

func Test_Client_GetShardNum_Account(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getshardnum", "--account", common.Account2)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("getshardnum error:%s", err)
	} else {
		if !strings.Contains(string(output), "2") {
			t.Fatalf("Test_Client_GetShardNum_Account,getshardnum returns error shardnum")
		}
	}
}

func Test_Client_GetShardNum_PrivateKey_Invalid_Without_Prefix_0x(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getshardnum", "--privatekey", "1234")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_GetShardNum_PrivateKey_Invalid_Without_Prefix_0x: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "failed to load the private key: Input string not a valid ecdsa string") {
		t.Fatalf("Test_Client_GetShardNum_PrivateKey_Invalid_Without_Prefix_0x,getshardnum should return error with privatekey without prefix 0x: %s", errStr)
	}
}

func Test_Client_GetShardNum_PrivateKey_Invalid_With_Prefix_0x_Odd(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getshardnum", "--privatekey", "0x123")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_GetShardNum_PrivateKey_Invalid_With_Prefix_0x_Odd: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "failed to load the private key: encoding/hex: odd length hex string") {
		t.Fatalf("Test_Client_GetShardNum_PrivateKey_Invalid_With_Prefix_0x_Odd,getshardnum should return error with privatekey of odd length: %s", errStr)
	}
}

func Test_Client_GetShardNum_PrivateKey_Invalid_Syntax_Character(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getshardnum", "--privatekey", "0x12345-")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_GetShardNum_PrivateKey_Invalid_Syntax_Character: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "failed to load the private key: encoding/hex: invalid byte: U+002D '-'") {
		t.Fatalf("Test_Client_GetShardNum_PrivateKey_Invalid_Syntax_Character,getshardnum should return error with privatekey has syntax character: %s", errStr)
	}
}

func Test_Client_GetShardNum_PrivateKey(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getshardnum", "--privatekey", common.AccountPrivateKey2)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_GetShardNum_PrivateKey,getshardnum error:%s", err)
	} else {
		if !strings.Contains(string(output), "2") {
			t.Fatalf("Test_Client_GetShardNum_PrivateKey,getshardnum returns error shardnum")
		}
	}
}

// --------------------test getshardnum end-------------------

// --------------------test key start-------------------
func Test_Client_Key_Invalid_Shard_Greater_Than_2(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "key", "--shard", "3")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Key_Invalid_Shard_Greater_Than_2: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "not supported shard number, shard number should be [0, 2]") {
		t.Fatalf("Test_Client_Key_Invalid_Shard_Greater_Than_2,key should return error with shard greater than 2: %s", errStr)
	}
}

func Test_Client_Key_Invalid_Shard_Non_Numerical(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "key", "--shard", "a")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Key_Invalid_Shard_Non_Numerical: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "parsing \"a\": invalid syntax") {
		t.Fatalf("Test_Client_Key_Invalid_Shard_Non_Numerical,key should return error with shard non-numerical: %s", errStr)
	}
}

func Test_Client_Key_Invalid_Shard_Empty(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "key", "--shard", "")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Key_Invalid_Shard_Empty: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "parsing \"\": invalid syntax") {
		t.Fatalf("Test_Client_Key_Invalid_Shard_Empty,key should return error with empty shard: %s", errStr)
	}
}

// --------------------test key end-------------------

// --------------------test sign start-------------------
func Test_Client_Sign_Invalid_privatekey_Without_Prefix_0x(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sign", "--privatekey", "123")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Sign_Invalid_privatekey_Without_Prefix_0x: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "failed to load key Input string not a valid ecdsa string") {
		t.Fatalf("Test_Client_Sign_Invalid_privatekey_Without_Prefix_0x,sign should return error with privatekey without prefix 0x: %s", errStr)
	}
}

func Test_Client_Sign_Invalid_privatekey_With_Prefix_Odd(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sign", "--privatekey", "0x123")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Sign_Invalid_privatekey_With_Prefix_Odd: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "failed to load key encoding/hex: odd length hex string") {
		t.Fatalf("Test_Client_Sign_Invalid_privatekey_With_Prefix_Odd,sign should return error with privatekey of odd length: %s", errStr)
	}
}

func Test_Client_Sign_Invalid_privatekey_With_Syntax_Character(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sign", "--privatekey", "0x12345-")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Sign_Invalid_privatekey_With_Syntax_Character: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "failed to load key encoding/hex: invalid byte: U+002D '-'") {
		t.Fatalf("Test_Client_Sign_Invalid_privatekey_With_Syntax_Character,sign should return error with privatekey has syntax character: %s", errStr)
	}
}

func Test_Client_Sign_Invalid_To_Address_Empty(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sign", "--privatekey", common.AccountPrivateKey2, "--to", "")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Sign_Invalid_To_Address_Empty: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "invalid amount value") {
		t.Fatalf("Test_Client_Sign_Invalid_To_Address_Empty,sign should return error with empty to address: %s", errStr)
	}
}

func Test_Client_Sign_Invalid_To_Address_Without_Prefix_0x(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sign", "--privatekey", common.AccountPrivateKey2, "--to", "123")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Sign_Invalid_To_Address_Without_Prefix_0x: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "invalid receiver address: hex string without 0x prefix") {
		t.Fatalf("Test_Client_Sign_Invalid_To_Address_Without_Prefix_0x,sign should return error with the to address without prefix 0x: %s", errStr)
	}
}

func Test_Client_Sign_Invalid_To_Address_With_Prefix_Odd(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sign", "--privatekey", common.AccountPrivateKey2, "--to", "0x123")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Sign_Invalid_To_Address_With_Prefix_Odd: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "invalid receiver address: hex string of odd length") {
		t.Fatalf("Test_Client_Sign_Invalid_To_Address_With_Prefix_Odd,sign should return error with the to address of odd length: %s", errStr)
	}
}

func Test_Client_Sign_Invalid_To_Address_With_Syntax_Character(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sign", "--privatekey", common.AccountPrivateKey2, "--to", "0x12345-")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Sign_Invalid_To_Address_With_Syntax_Character: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "invalid receiver address: invalid hex string") {
		t.Fatalf("Test_Client_Sign_Invalid_To_Address_With_Syntax_Character,sign should return error with the to address has syntax character: %s", errStr)
	}
}

func Test_Client_Sign_Invalid_Amount_With_Empty(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sign", "--privatekey", common.AccountPrivateKey2, "--amount", "")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Sign_Invalid_Amount_With_Empty: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "invalid amount value") {
		t.Fatalf("Test_Client_Sign_Invalid_Amount_With_Empty,sign should return error with empty amount: %s", errStr)
	}
}

func Test_Client_Sign_Invalid_Amount_With_Non_Numerical(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sign", "--privatekey", common.AccountPrivateKey2, "--amount", "a")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Sign_Invalid_Amount_With_Non_Numerical: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "invalid amount value") {
		t.Fatalf("Test_Client_Sign_Invalid_Amount_With_Non_Numerical,sign should return error with  amount non-numerical: %s", errStr)
	}
}

func Test_Client_Sign_Invalid_Price_With_Empty(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sign", "--privatekey", common.AccountPrivateKey2, "--amount", "2", "--price", "")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Sign_Invalid_Price_With_Empty: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "invalid gas price value") {
		t.Fatalf("Test_Client_Sign_Invalid_Price_With_Empty,sign should return error with empty price: %s", errStr)
	}
}

func Test_Client_Sign_Invalid_Price_With_Non_Numerical(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sign", "--privatekey", common.AccountPrivateKey2, "--amount", "2", "--price", "a")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Sign_Invalid_Price_With_Non_Numerical: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "invalid gas price value") {
		t.Fatalf("Test_Client_Sign_Invalid_Price_With_Non_Numerical,sign should return error with price non-numerical: %s", errStr)
	}
}

func Test_Client_Sign_Invalid_Gaslimit_With_Non_Numerical(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sign", "--privatekey", common.AccountPrivateKey2, "--amount", "2", "--price", "1", "--gas", "a")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Sign_Invalid_Price_With_Non_Numerical: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "parsing \"a\": invalid syntax") {
		t.Fatalf("Test_Client_Sign_Invalid_Price_With_Non_Numerical,sign should return error with gas non-numerical: %s", errStr)
	}
}

func Test_Client_Sign_Invalid_Gaslimit_With_Empty(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sign", "--privatekey", common.AccountPrivateKey2, "--amount", "2", "--price", "2", "--gas", "")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Sign_Invalid_Gaslimit_With_Empty: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "parsing \"\": invalid syntax") {
		t.Fatalf("Test_Client_Sign_Invalid_Gaslimit_With_Empty,sign should return error with empty gas: %s", errStr)
	}
}

func Test_Client_Sign_Invalid_Gaslimit_With_Non_Integer(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sign", "--privatekey", common.AccountPrivateKey2, "--amount", "2", "--price", "2", "--gas", "17.5")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Sign_Invalid_Gaslimit_With_Non_Integer: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "parsing \"17.5\": invalid syntax") {
		t.Fatalf("Test_Client_Sign_Invalid_Gaslimit_With_Non_Int,sign should return error with gas non-integer: %s", errStr)
	}
}

func Test_Client_Sign_Invalid_Gaslimit_With_Negative_Integer(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sign", "--privatekey", common.AccountPrivateKey2, "--amount", "2", "--price", "2", "--gas", "-17")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Sign_Invalid_Gaslimit_With_Negative_Integer: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "invalid value \"-17\" for flag -gas: strconv.ParseUint: parsing \"-17\": invalid syntax") {
		t.Fatalf("Test_Client_Sign_Invalid_Gaslimit_With_Negative_Integer,sign should return error with gas negative integer: %s", errStr)
	}
}

func Test_Client_Sign_Invalid_Nonce_With_Negative_Integer(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sign", "--privatekey", common.AccountPrivateKey2, "--amount", "2", "--price", "2", "--gas", "1000000000000", "--nonce", "-1")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Sign_Invalid_Nonce_With_Negative_Integer: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "strconv.ParseUint: parsing \"-1\": invalid syntax") {
		t.Fatalf("Test_Client_Sign_Invalid_Nonce_With_Negative_Integer,sign should return error with nonce negative integer: %s", errStr)
	}
}

func Test_Client_Sign_Invalid_Nonce_With_Non_Integer(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sign", "--privatekey", common.AccountPrivateKey2, "--amount", "2", "--price", "2", "--gas", "1000000000000", "--nonce", "17.5")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Sign_Invalid_Nonce_With_Non_Integer: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "parsing \"17.5\": invalid syntax") {
		t.Fatalf("Test_Client_Sign_Invalid_Nonce_With_Non_Integer,sign should return error with nonce non-integer: %s", errStr)
	}
}

func Test_Client_Sign_Invalid_Nonce_With_Non_Numeric(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sign", "--privatekey", common.AccountPrivateKey2, "--amount", "2", "--price", "2", "--gas", "1000000000000", "--nonce", "a")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Sign_Invalid_Nonce_With_Non_Numeric: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "parsing \"a\": invalid syntax") {
		t.Fatalf("Test_Client_Sign_Invalid_Nonce_With_Non_Numeric,sign should return error with nonce non-numeric: %s", errStr)
	}
}

func Test_Client_Sign_Invalid_Nonce_With_Empty(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sign", "--privatekey", common.AccountPrivateKey2, "--amount", "2", "--price", "2", "--gas", "1000000000000", "--nonce", "")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Sign_Invalid_Nonce_With_Empty: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "parsing \"\": invalid syntax") {
		t.Fatalf("Test_Client_Sign_Invalid_Nonce_With_Empty,sign should return error with empty nonce: %s", errStr)
	}
}

func Test_Client_Sign_Invalid_Payload_With_Empty(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sign", "--privatekey", common.AccountPrivateKey2, "--amount", "2", "--price", "2", "--gas", "1000000000000", "--nonce", "")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Sign_Invalid_Payload_With_Empty: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "parsing \"\": invalid syntax") {
		t.Fatalf("Test_Client_Sign_Invalid_Payload_With_Empty,sign should return error with empty payload: %s", errStr)
	}
}

func Test_Client_Sign_Invalid_Payload_Without_Prefix_0x(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sign", "--privatekey", common.AccountPrivateKey2, "--amount", "2", "--price", "2", "--gas", "1", "--nonce", "1", "--payload", "aaa")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Sign_Invalid_Payload_Without_Prefix_0x: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "hex string without 0x prefix") {
		t.Fatalf("Test_Client_Sign_Invalid_Payload_Without_Prefix_0x,sign should return error with the to address without prefix 0x: %s", errStr)
	}
}

func Test_Client_Sign_Invalid_Payload_With_Prefix_Odd(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sign", "--privatekey", common.AccountPrivateKey2, "--amount", "2", "--price", "2", "--gas", "1", "--nonce", "1", "--payload", "0x123")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Sign_Invalid_To_Address_With_Prefix_Odd: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "hex string of odd length") {
		t.Fatalf("Test_Client_Sign_Invalid_To_Address_With_Prefix_Odd,sign should return error with the to address of odd length: %s", errStr)
	}
}

func Test_Client_Sign_Invalid_Payload_With_Syntax_Characeter(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sign", "--privatekey", common.AccountPrivateKey2, "--amount", "2", "--price", "2", "--gas", "1", "--nonce", "1", "--payload", "0x12345-")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Sign_Invalid_To_Address_With_Syntax_Characeter: %s", err)
	}
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "invalid hex string") {
		t.Fatalf("Test_Client_Sign_Invalid_To_Address_With_Syntax_Characeter,sign should return error with the to address has syntax character: %s", errStr)
	}
}

// --------------------test sign end-------------------

// --------------------test sendtx start-------------------
func Test_Client_SendTx_InvalidAccountLength(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sendtx", "--amount", "10000", "--price", "1", "--from", common.KeyFileShard1_1, "--to", "0x")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_Client_SendTx_InvalidAccountLength: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "invalid address") {
		t.Fatalf("Test_Client_SendTx_InvalidAccountLength Err:%s", errStr)
	}
}

func Test_Client_SendTx_InvalidAccountType(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sendtx", "--amount", "10000", "--price", "1", "--from", common.KeyFileShard2_1, "--to", common.InvalidAccountType)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_Client_SendTx_testcase.InvalidAccountType: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "invalid address type") {
		t.Fatalf("Test_Client_SendTx_testcase.InvalidAccountType Err:%s", errStr)
	}
}

func Test_Client_SendTx_InvalidAmountValue(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sendtx", "--amount", "", "--price", "1", "--from", common.KeyFileShard1_1, "--to", common.AccountShard1_2)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_Client_SendTx_InvalidAccountValue: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "invalid amount value") {
		t.Fatalf("Test_Client_SendTx_InvalidAccountValue Err:%s", errStr)
	}
}

func Test_Client_SendTx_InvalidPriceValue(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sendtx", "--amount", "10000", "--price", "", "--from", common.KeyFileShard1_1, "--to", common.AccountShard1_2)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_Client_SendTx_testcase.InvalidAccountType: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "invalid gas price value") {
		t.Fatalf("Test_Client_SendTx_testcase.InvalidAccountType Err:%s", errStr)
	}
}

func Test_Client_SendTx_Unmatched_keyfile_And_Pass(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sendtx", "--amount", "10000", "--price", "1", "--from", common.KeyFileShard2_1, "--to", common.InvalidAccountType)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_Client_SendTx_Unmatched_keyfile_And_Pass: An error occured: %s", err)
	}

	io.WriteString(stdin, "123456\n")
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "could not decrypt key with given passphrase") {
		t.Fatalf("Test_Client_SendTx_Unmatched_keyfile_And_Pass Err:%s", errStr)
	}
}

func Test_Client_SendTx_Invalid_Gas(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sendtx", "--amount", "10000", "--price", "", "--from", common.KeyFileShard2_1, "--to", common.Account2, "--gas", "")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_Client_SendTx_Invalid_Gas: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "parsing \"\": invalid syntax") {
		t.Fatalf("Test_Client_SendTx_Invalid_Gas Err:%s", errStr)
	}
}

func Test_Client_SendTx_Invalid_Payload_Without_Prefix_0x(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sendtx", "--amount", "10000", "--price", "1", "--from", common.KeyFileShard2_1, "--to", common.Account2, "--gas", "1", "--payload", "-1")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_Client_SendTx_Invalid_Payload_Without_Prefix_0x: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "hex string without 0x prefix") {
		t.Fatalf("Test_Client_SendTx_Invalid_Payload_Without_Prefix_0x Err:%s", errStr)
	}
}

func Test_Client_SendTx_Invalid_Payload_With_Prefix_Odd(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sendtx", "--amount", "10000", "--price", "1", "--from", common.KeyFileShard2_1, "--to", common.Account2, "--gas", "1", "--payload", "0x123")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_Client_SendTx_Invalid_Payload_With_Prefix_Odd: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "hex string of odd length") {
		t.Fatalf("Test_Client_SendTx_Invalid_Payload_With_Prefix_Odd Err:%s", errStr)
	}
}

func Test_Client_SendTx_Invalid_Payload_With_Syntax_Characeter(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sendtx", "--amount", "10000", "--price", "1", "--from", common.KeyFileShard2_1, "--to", common.Account2, "--gas", "1", "--payload", "0x12345-")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_Client_SendTx_Invalid_Payload_With_Syntax_Characeter: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "invalid hex string") {
		t.Fatalf("Test_Client_SendTx_Invalid_Payload_With_Syntax_Characeter Err:%s", errStr)
	}
}

func Test_Client_SendTx_Invalid_Nonce(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sendtx", "--amount", "10000", "--price", "1", "--from", common.KeyFileShard2_1, "--to", common.Account2, "--gas", "1", "--payload", "1", "--nonce", "")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_Client_SendTx_Invalid_Nonce: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "parsing \"\": invalid syntax") {
		t.Fatalf("Test_Client_SendTx_Invalid_Nonce Err:%s", errStr)
	}
}

// --------------------test sendtx end-------------------

// --------------------test deckeyfile start-------------------
func Test_Client_Deckeyfile_Invalid_Pass(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "deckeyfile", "--file", common.KeyFileShard2_1)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Deckeyfile_Invalid_Pass: An error occured: %s", err)
	}

	io.WriteString(stdin, "1234\n")
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "could not decrypt key with given passphrase") {
		t.Fatalf("Test_Client_Deckeyfile_Invalid_Pass Err:%s", errStr)
	}
}

func Test_Client_Deckeyfile_Invalid_Keyfile(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "deckeyfile", "--file", "../config/keyfile/shard1-0x1234567890")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Deckeyfile_Invalid_Keyfile: An error occured: %s", err)
	}

	io.WriteString(stdin, "1234\n")
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if !strings.Contains(errStr, "invalid key file") {
		t.Fatalf("Test_Client_Deckeyfile_Invalid_Keyfile Err:%s", errStr)
	}
}

func Test_Client_Deckeyfiles(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "deckeyfile", "--file", common.KeyFileShard2_1)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	defer stdin.Close()

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &outErr

	if err = cmd.Start(); err != nil {
		t.Fatalf("Test_Client_Deckeyfiles: An error occured: %s", err)
	}

	io.WriteString(stdin, "123\n")
	cmd.Wait()

	_, errStr := out.String(), outErr.String()

	if errStr != "" {
		t.Fatalf("Test_Client_Deckeyfiles Err:%s", errStr)
	}
}

// --------------------test deckeyfile end-------------------
func Test_Client_GetBlockHeight_NodeStop(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getblockheight", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_GetBlockHeight error, %s", err)
	}
}

func Test_Client_GetBlockHeight_NodeStart(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getblockheight", "--address", common.ServerAddr)
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
	cmd := exec.Command(common.CmdClient, "getblockheight", "--height", "1000000000", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_GetBlockHeight_InvalidParameter returns ok with invalid parameter")
	}
}

func Test_Client_GetBlockHeight_ByInvalidHeight(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getblockheight", "--height", "1", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_GetBlockHeight_ByInvalidHeight returns error not defined: -height")
	}
}

func Test_Client_GetBlockHeight_Parameter(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getblockheight", "1", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_GetBlockHeight_Parameter error, %s", err)
	}
}

func Test_Client_GetBlockTXCount_ByInvalidHeight(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getblocktxcount", "--height", "100000000", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_GetBlockTXCount_ByInvalidHeight error parameter success?")
	}
}

func Test_Client_GetBlockTXCount_ByInvalidHeight0x(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getblocktxcount", "--height", "0x", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_GetBlockTXCount_ByInvalidHeight0x return error invalid value")
	}
}

func Test_Client_GetBlockTXCount_ByHeight_NodeStart(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getblocktxcount", "--height", "1", "--address", common.ServerAddr)
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
	cmd := exec.Command(common.CmdClient, "getblocktxcount", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_GetBlockTXCount_DefaultParameter:error, %s", err)
	}
}

func Test_Client_GetBlockTXCount_InvalidParameter(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getblocktxcount", "1", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_GetBlockTXCount_InvalidParameter： error parameter success?")
	}
}

func Test_Client_GetBlockTXCount_ByInvalidHash(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getblocktxcount", "--hash", common.BlockHashErr, "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_GetBlockTXCount_ByInvalidHash： error parameter success?")
	}
}

func Test_Client_GetBlockTXCount_ByHash(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sendtx", "--amount", "900", "--price", "1", "--gas", "2", "--from", common.KeyFileShard1_5, "--to", common.Account1_Aux2)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
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
	outStr = outStr[strings.Index(outStr, "{"):]
	outStr = strings.Trim(outStr, "\n")
	outStr = strings.Trim(outStr, " ")
	var txInfo common.TxInfo
	if err = json.Unmarshal([]byte(outStr), &txInfo); err != nil {
		return
	}
	blockHash := Gettxbyhash(txInfo.Hash)
	cmd = exec.Command(common.CmdClient, "getblocktxcount", "--hash", blockHash, "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_GetBlockTXCount_ByHash: getblocktxcount error, %s", err)
	}
}

func Test_Client_GetBlockTXCount_ByInvalidHash0x(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getblocktxcount", "--hash", "0x", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Light_GetBlockTXCount_ByInvalidHash0x error parameter success?")
	}
}

func Test_Client_GetBlockTXCount_ByInvalidHash0x12(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getblocktxcount", "--hash", "0x12-", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_GetBlockTXCount_ByInvalidHash0x12： return error syntax character")
	}
}

func Test_Client_GetBlockTXCount_ByInvalidHash123(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getblocktxcount", "--hash", "123", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_GetBlockTXCount_ByInvalidHash123： return error hex string without 0x prefix")
	}
}

func Test_Client_GetBlock_ByHeight_NodeStop(t *testing.T) {
	// Normal height
	cmd := exec.Command(common.CmdClient, "getblock", "--height", "1", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_GetBlock_ByHeight_NodeStop: error, %s", err)
	}
}

func Test_Client_GetBlock_ByHeight_NodeStart(t *testing.T) {
	// Normal height
	cmd := exec.Command(common.CmdClient, "getblock", "--height", "1", "--address", common.ServerAddr)
	if res, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_GetBlock_ByHeight_NodeStart:Node to run returns error: %s", err)
	} else {
		var blockInfo common.BlockInfo
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
	cmd := exec.Command(common.CmdClient, "getblock", "--height", "10000000000", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_GetBlock_ByInvalidHeight: error parameter success?")
	}
}

func Test_Client_GetBlock_InvalidParameter(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getblock", "1", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_GetBlock_InvalidParameter error, %s", err)
	}
}

func Test_Client_GetBlock_ByNormalHash(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sendtx", "--amount", "900", "--price", "1", "--gas", "2", "--from", common.KeyFileShard1_5, "--to", common.Account1_Aux2)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
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
	outStr = outStr[strings.Index(outStr, "{"):]
	outStr = strings.Trim(outStr, "\n")
	outStr = strings.Trim(outStr, " ")
	var txInfo common.TxInfo
	if err = json.Unmarshal([]byte(outStr), &txInfo); err != nil {
		return
	}
	blockHash := Gettxbyhash(txInfo.Hash)
	cmd = exec.Command(common.CmdClient, "getblock", "--hash", blockHash, "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_GetBlock_ByNormalHash error, %s", err)
	}
}

func Test_Client_GetBlock_ByInvalidHash(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getblock", "--hash", common.BlockHashErr, "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_GetBlock_ByInvalidHash error parameter success?")
	}
}

// getblock fulltx support.
func Test_Client_GetBlock_ByHeightFulltx(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getblock", "--height", "1", "--fulltx", "--address", common.ServerAddr)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_GetBlock_ByHeightFulltx error, %s", err)
	} else {
		var blockInfo common.BlockInfo
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
// func Test_Client_GetBlock_ByHashFulltx(t *testing.T) {
// 	cmd := exec.Command(common.CmdClient, "sendtx", "--amount", "90", "--price", "1", "--gas", "2", "--from", common.KeyFileShard1_5, "--to", common.Account1_Aux2)
// 	stdin, err := cmd.StdinPipe()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	var out bytes.Buffer
// 	var outErr bytes.Buffer
// 	cmd.Stdout, cmd.Stderr = &out, &outErr
// 	if err = cmd.Start(); err != nil {
// 		return
// 	}
// 	io.WriteString(stdin, "123\n")
// 	cmd.Wait()
// 	outStr, errStr := out.String(), outErr.String()
// 	if len(string(errStr)) > 0 {
// 		err = errors.New(string(errStr))
// 		return
// 	}
// 	outStr = outStr[strings.Index(outStr, "{"):]
// 	outStr = strings.Trim(outStr, "\n")
// 	outStr = strings.Trim(outStr, " ")
// 	var txInfo common.TxInfo
// 	if err = json.Unmarshal([]byte(outStr), &txInfo); err != nil {
// 		return
// 	}

// 	blockHash := Gettxbyhash(txInfo.Hash)
// 	cmd = exec.Command(common.CmdClient, "getblock", "--hash", blockHash, "--fulltx", "--address", common.ServerAddr)
// 	if output, err := cmd.CombinedOutput(); err != nil {
// 		t.Fatalf("Test_Client_GetBlock_ByHashFulltx error, %s", err)
// 	} else {
// 		var blockInfo common.BlockInfo
// 		if err = json.Unmarshal(output, &blockInfo); err != nil {
// 			t.Fatalf("Test_Client_GetBlock_ByHashFulltx: %s", err)
// 		}
// 		if blockInfo.Hash != blockHash {
// 			t.Fatalf("Test_Client_GetBlock_ByHashFulltx: Expect the return value is not correct!")
// 		}
// 		if len(blockInfo.Transactions) <= 0 {
// 			t.Fatalf("Test_Client_GetBlock_ByHashFulltx, hash should contain one transaction at lease")
// 		}
// 	}
// }

// func Test_Client_GetLogs_ValidParameter(t *testing.T) {
// 	contract, height, topics, err := common.DeployContractAndSendTx(t)
// 	if err != nil {
// 		t.Fatalf("Test_Client_GetLogs_ValidParameter err %s", err.Error())
// 	}
// 	for _, topic := range topics {
// 		cmd := exec.Command(common.CmdClient, "getlogs", "--height", height, "--contract", contract, "--topic", topic, "--address", common.ServerAddr)
// 		if result, err := cmd.CombinedOutput(); err != nil {
// 			t.Fatalf("Test_Client_GetLogs_ValidParameter: An error occured: %s", err)
// 		} else {
// 			var logs []common.LogByTopic
// 			if err = json.Unmarshal(result, &logs); err != nil {
// 				t.Fatalf("Test_Client_GetLogs_ValidParameter getlogs unmarshal err %s", err)
// 			}
// 			if len(logs) != 1 {
// 				t.Fatal("Test_Client_GetLogs_ValidParameter returns log number is not 1")
// 			}
// 		}
// 	}
// }

// func Test_Client_GetLogs_Invalid_Topic(t *testing.T) {
// 	contract, height, _, err := common.DeployContractAndSendTx(t)
// 	if err != nil {
// 		t.Fatalf("Test_Client_GetLogs_Invalid_Topic err %s", err.Error())
// 	}
// 	errTopic := "0xaaaaaa"
// 	cmd := exec.Command(common.CmdClient, "getlogs", "--height", height, "--contract", contract, "--topic", errTopic, "--address", common.ServerAddr)
// 	a, err := cmd.CombinedOutput()
// 	if err != nil {
// 		t.Fatalf("Test_Client_GetLogs_Invalid_Topic: An error occured: %s", err)
// 	}
// 	if string(a) != "[]\n" {
// 		t.Fatalf("Test_Client_GetLogs_Invalid_Topic returns a log")
// 	}
// }

// func Test_Client_GetLogs_Invalid_Contract(t *testing.T) {
// 	_, height, topics, err := common.DeployContractAndSendTx(t)
// 	if err != nil {
// 		t.Fatalf("Test_Client_GetLogs_Invalid_Contract err %s", err.Error())
// 	}
// 	errContract := "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
// 	for _, topic := range topics {
// 		cmd := exec.Command(common.CmdClient, "getlogs", "--height", height, "--contract", errContract, "--topic", topic, "--address", common.ServerAddr)
// 		if result, err := cmd.CombinedOutput(); err != nil {
// 			t.Fatalf("Test_Client_GetLogs_Invalid_Contract: An error occured: %s", err)
// 		} else {
// 			var logs []common.LogByTopic
// 			if err = json.Unmarshal(result, &logs); err != nil {
// 				t.Fatalf("Test_Client_GetLogs_Invalid_Contract getlogs unmarshal err %s", err)
// 			}
// 			if len(logs) != 0 {
// 				t.Fatal("Test_Client_GetLogs_Invalid_Contract returns log number is not 0")
// 			}
// 		}
// 	}
// }

func Test_Client_GetLogs_Invalid_Length_Contract(t *testing.T) {
	_, height, topics, err := common.DeployContractAndSendTx(t)
	if err != nil {
		t.Fatalf("Test_Client_GetLogs_Invalid_Length_Contract err %s", err.Error())
	}
	errContract := "0xaaaaaaaaaaaaaaaaaaaaa"
	for _, topic := range topics {
		cmd := exec.Command(common.CmdClient, "getlogs", "--height", height, "--contract", errContract, "--topic", topic, "--address", common.ServerAddr)
		if _, err := cmd.CombinedOutput(); err == nil {
			t.Fatalf("Test_Client_GetLogs_Invalid_Length_Contract return ok")
		}
	}
}

// func Test_Client_GetLogs_Invalid_height(t *testing.T) {
// 	contract, _, topics, err := common.DeployContractAndSendTx(t)
// 	if err != nil {
// 		t.Fatalf("Test_Client_GetLogs_Invalid_height err %s", err.Error())
// 	}
// 	errHeight := "1.5"
// 	for _, topic := range topics {
// 		cmd := exec.Command(common.CmdClient, "getlogs", "--height", errHeight, "--contract", contract, "--topic", topic, "--address", common.ServerAddr)
// 		if _, err := cmd.CombinedOutput(); err == nil {
// 			t.Fatal("Test_Client_GetLogs_Invalid_height returns ok")
// 		}
// 	}
// }

func Test_Client_GetNonce_ByAccount(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getnonce", "--account", common.Account1_Aux, "--address", common.ServerAddr)
	if res, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("getnonce returns with error input err: %s", err)
	} else {
		res = bytes.TrimRight(res, "\n")
		nonce, _ := strconv.ParseInt(string(res), 10, 64)
		if nonce < 0 {
			t.Fatalf("Test_Client_testcase.GetNonce_InvalidAccount: Nonce value is not correct!")
		}
	}
}

func Test_Client_GetNonce_InvalidAccount0x(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getnonce", "--account", "0x", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_testcase.GetNonce_InvalidAccount0x returns err: %s", err)
	}
}

func Test_Client_GetNonce_InvalidAccount(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getnonce", "--account", common.AccountErr, "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_testcase.GetNonce_InvalidAccount returns error： hex string of odd length")
	}
}

func Test_Client_GetNonce_NoParameter(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getnonce", "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_testcase.GetNonce_NoParameter returns error： invalid account")
	}
}

func Test_Client_GetNonce_invalidParameter(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getnonce", common.Account1_Aux, "--address", common.ServerAddr)
	out, err := cmd.CombinedOutput()

	if !strings.Contains(string(out), "flag is not specified for value") {
		t.Fatal("Test_Client_GetNonce_invalidParameter failed is not ok")
	}

	if err != nil {
		t.Fatalf("Test_Client_testcase.GetNonce_invalidParameter returns error： invalid account")
	}
}

func Test_Client_GetNonce_AccountFromOtherShard(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getnonce", "--account", common.AccountShard1_1, "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_testcase.GetNonce_AccountFromOtherShard:getnonce returns successfully for other shard account")
	}
}

func Test_Client_GetReceipt_ByInvalidHash0x(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getreceipt", "0x", "--address", common.ServerAddr)

	out, err := cmd.CombinedOutput()
	if !strings.Contains(string(out), "flag is not specified for value") {
		t.Fatal("Test_Client_GetReceipt_ByInvalidHash0x failed is not ok")
	}

	if err != nil {
		t.Fatalf("Test_Client_GetReceipt_ByInvalidHash0x error: empty hex string")
	}
}

func Test_Client_GetReceipt_InvalidParameter(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sendtx", "--amount", "900", "--price", "1", "--gas", "2", "--from", common.KeyFileShard1_5, "--to", common.Account1_Aux2)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
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
	outStr = outStr[strings.Index(outStr, "{"):]
	outStr = strings.Trim(outStr, "\n")
	outStr = strings.Trim(outStr, " ")
	var txInfo common.TxInfo
	if err = json.Unmarshal([]byte(outStr), &txInfo); err != nil {
		return
	}

	cmd = exec.Command(common.CmdClient, "getreceipt", txInfo.Hash, "--address", common.ServerAddr)

	output, err := cmd.CombinedOutput()
	if !strings.Contains(string(output), "flag is not specified for value") {
		t.Fatal("Test_Client_GetReceipt_InvalidParameter failed is not ok")
	}

	if err != nil {
		t.Fatalf("Test_Client_GetReceipt  error ：Grammar is not correct")
	}

}
func Test_Client_SendManyTx(t *testing.T) {
	curNonce, err := common.GetNonce(t, common.CmdClient, common.AccountShard1_5, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_Client_SendManyTx : getnonce returns with error input %s", err)
	}

	beginBalance, err := common.GetBalance(t, common.CmdClient, common.AccountShard1_5, common.ServerAddr)
	if err != nil {
		t.Fatalf("Test_Client_SendManyTx : common.GetBalance returns with error input  %s", err)
	}
	if beginBalance == 0 {
		t.Fatalf("Test_Client_SendManyTx : common.GetBalance Insufficient amount of account")
	}

	var txHash string
	var sendTxL []*common.SendTxInfo

	for cnt := 0; cnt < 5; cnt++ {
		itemNonce := curNonce + 2 + cnt
		txHash, _, err = common.SendTx(t, common.CmdClient, 100, itemNonce, 2100, common.KeyFileShard1_5, common.Account1_Aux2, "", common.ServerAddr)
		if err != nil {
			t.Fatalf("Test_Client_SendManyTx: An error occured: %s", err)
		}

		info := &common.SendTxInfo{
			Nonce:  itemNonce,
			Hash:   txHash,
			BMined: false,
		}
		sendTxL = append(sendTxL, info)
	}

	cnt := 0
	for {
		pendingL, err1 := common.GetPendingTxs(t, common.CmdClient, common.ServerAddr)
		if err1 != nil {
			t.Fatalf("Test_Client_SendManyTx : common.GetPendingTxs err:%s", err1)
		}
		contentM, err2 := common.GetPoolContentTxs(t, common.CmdClient, common.ServerAddr)
		if err2 != nil {
			t.Fatalf("Test_Client_SendManyTx : common.GetPoolContentTxs err:%s", err2)
		}
		_, err3 := common.GetPoolCountTxs(t, common.CmdClient, common.ServerAddr)
		if err3 != nil {
			t.Fatalf("Test_Client_SendManyTx : common.GetPoolCountTxs err:%s", err3)
		}
		if len(pendingL)+len(contentM) == 0 {
			break
		}
		cnt++
	}
	time.Sleep(8 * time.Second)
	validCnt := 0
	for _, sendTxInfo := range sendTxL {
		info, err3 := common.GetReceipt(t, common.CmdClient, sendTxInfo.Hash, common.ServerAddr)
		if err3 == nil {
			if info.Hash != sendTxInfo.Hash {
				fmt.Println("Test_Client_SendManyTx :  Receipt Hash not match with tx")
			}
			validCnt++
			sendTxInfo.BMined = true
		} else {
			fmt.Println("Test_Client_SendManyTx : getReceipt err. nonce=", sendTxInfo.Nonce, err3)
		}
	}
}

func Test_Client_GetTxInBlock_ByHeightindex(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "gettxinblock", "--height", "1", "--index", "0")
	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Test_Client_GetTxInBlock_ByHeightindex err=%s", err)
	}
}

func Test_Client_GetTxInBlock_ByHeight(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "gettxinblock", "--height", "1")
	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Test_Client_GetTxInBlock_ByHeight err=%s", err)
	}
}

func Test_Client_GetTxInBlock_ByInvalidHeight(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "gettxinblock", "--height", "1000000000", "--index", "0")
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Test_Client_GetTxInBlock_ByInvalidHeight err=leveldb: not found")
	}
}

func Test_Client_GetTxInBlock_ByHashindex(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sendtx", "--amount", "900", "--price", "1", "--gas", "2", "--from", common.KeyFileShard1_5, "--to", common.Account1_Aux2)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
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
	outStr = outStr[strings.Index(outStr, "{"):]
	outStr = strings.Trim(outStr, "\n")
	outStr = strings.Trim(outStr, " ")
	var txInfo common.TxInfo
	if err = json.Unmarshal([]byte(outStr), &txInfo); err != nil {
		return
	}
	blockHash := Gettxbyhash(txInfo.Hash)
	cmd = exec.Command(common.CmdClient, "gettxinblock", "--hash", blockHash, "--index", "0")
	_, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Test_Client_GetTxInBlock_ByHashindex err=%s", err)
	}
}

func Gettxbyhash(txhash string) (blockHash string) {
ErrContinue:
	cmd := exec.Command(common.CmdClient, "gettxbyhash", "--hash", txhash, "--address", common.ServerAddr)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
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
	outStr = outStr[strings.Index(outStr, "{"):]
	outStr = strings.Trim(outStr, "\n")
	outStr = strings.Trim(outStr, " ")
	blocktxInfo := make(map[string]interface{})
	if err = json.Unmarshal([]byte(outStr), &blocktxInfo); err != nil {
		return
	}
	status := blocktxInfo["status"].(string)
	if status == "pool" {
		goto ErrContinue
	}
	blockHash = blocktxInfo["blockHash"].(string)
	return blockHash
}

func Test_Client_GetTxInBlock_ByHash(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sendtx", "--amount", "900", "--price", "1", "--gas", "2", "--from", common.KeyFileShard1_5, "--to", common.Account1_Aux2)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
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
	outStr = outStr[strings.Index(outStr, "{"):]
	outStr = strings.Trim(outStr, "\n")
	outStr = strings.Trim(outStr, " ")
	var txInfo common.TxInfo
	if err = json.Unmarshal([]byte(outStr), &txInfo); err != nil {
		return
	}

	blockHash := Gettxbyhash(txInfo.Hash)
	cmd = exec.Command(common.CmdClient, "gettxinblock", "--hash", blockHash)
	_, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Test_Client_GetTxInBlock_ByHash err=%s", err)
	}
}

func Test_Client_GetTxInBlock_ByHashErr(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "gettxinblock", "--hash", common.BlockHashErr, "--index", "0")
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Test_Client_GetTxInBlock_ByHashErr err=leveldb: not found")
	}
}

func Test_Client_GetTxInBlock_ByHash0x(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "gettxinblock", "--hash", "0x", "--index", "0")
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Test_Client_GetTxInBlock_ByHash0x err=empty hex string")
	}
}

func Test_Client_GettxByHash(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sendtx", "--amount", "900", "--price", "1", "--gas", "2", "--from", common.KeyFileShard1_5, "--to", common.Account1_Aux2)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}

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

	outStr = outStr[strings.Index(outStr, "{"):]
	outStr = strings.Trim(outStr, "\n")
	outStr = strings.Trim(outStr, " ")
	var txInfo common.TxInfo
	if err = json.Unmarshal([]byte(outStr), &txInfo); err != nil {
		return
	}

	cmd = exec.Command(common.CmdClient, "gettxbyhash", "--hash", txInfo.Hash, "--address", common.ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_GettxByHash  error ：%s", err)
	}

	cmd = exec.Command(common.CmdClient, "gettxbyhash", txInfo.Hash, "--address", common.ServerAddr)

	output, err := cmd.CombinedOutput()
	if !strings.Contains(string(output), "flag is not specified for value") {
		t.Fatal("Test_Client_GettxByHash failed is not ok")
	}

	if err != nil {
		t.Fatalf("Test_Client_GettxByHash error ： empty hex string")
	}
}

func Test_Client_GettxByHash0x(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "gettxbyhash", "--hash", "0x", "--index", "0")
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Test_Client_GetTxInBlock_ByHash0x err=empty hex string")
	}
}

func Test_Client_Getdebtbyhash(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "sendtx", "--amount", "101", "--price", "1", "--gas", "2", "--from", common.KeyFileShard1_1, "--to", common.AccountShard2_2)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}

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

	outStr = outStr[strings.Index(outStr, "It is a cross shard transaction, its debt is:"):]
	outStr = outStr[strings.Index(outStr, "{"):]
	outStr = strings.Trim(outStr, "\n")
	outStr = strings.Trim(outStr, " ")
	var txInfo common.TxInfo

	if err = json.Unmarshal([]byte(outStr), &txInfo); err != nil {
		return
	}
	timeout := time.After(time.Second * 300)
	finish := make(chan bool)
	go func() {
		for {
			select {
			case <-timeout:
				finish <- true
				return
			default:
				cmd = exec.Command(common.CmdClient, "getdebtbyhash", "--hash", txInfo.Hash, "--address", common.ServertwoAddr)
				if _, err := cmd.CombinedOutput(); err == nil {
					finish <- true
				}
			}
			time.Sleep(time.Second * 2)
		}
	}()
	<-finish

	cmd = exec.Command(common.CmdClient, "getdebtbyhash", txInfo.Hash, "--address", common.ServertwoAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_Getdebtbyhash  flag is not specified for value")
	}
}

func Test_Client_Getdebtbyhash0x(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getdebtbyhash", "--hash", "0x")
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Test_Client_Getdebtbyhash0x err=empty hex string")
	}
}

func Test_Client_GetdebtbyhashAddr(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getdebtbyhash", "--hash", "0x", "--address", common.ServerAddr)
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Test_Client_GetdebtbyhashAddr err=empty hex string")
	}
}

func Test_Client_GetdebtbyhashtwoAddr(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getdebtbyhash", "--hash", "0x", "--address", common.ServertwoAddr)
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Test_Client_GetdebtbyhashAddr err=empty hex string")
	}
}

func Test_Client_Getdebts(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getdebts")
	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Test_Client_GetdebtbyhashAddr err=empty hex string")
	}
}

func Test_Client_GetdebtstwoAddr(t *testing.T) {
	cmd := exec.Command(common.CmdClient, "getdebts", "--address", common.ServertwoAddr)
	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Test_Client_GetdebtbyhashAddr err=empty hex string")
	}
}
