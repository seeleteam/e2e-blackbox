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
	"regexp"
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

// --------------------test savekey start-------------------
func Test_Client_SaveKey_Invalid_To_Privatekey_Without_Prefix_0x(t *testing.T) {
	cmd := exec.Command(CmdClient, "savekey", "--privatekey", "123")
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_SaveKey_Invalid_To_Privatekey_Without_Prefix_0x,savekey  should return error with privatekey without prefix 0x")
	}
}

func Test_Client_SaveKey_Invalid_To_Privatekey_With_Prefix_Odd(t *testing.T) {
	cmd := exec.Command(CmdClient, "savekey", "--privatekey", "0x123")
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_SaveKey_Invalid_To_Privatekey_With_Prefix_Odd,savekey should return error with privatekey is odd length")
	}
}

func Test_Client_SaveKey_Invalid_To_Privatekey_Syntax_Characeter(t *testing.T) {
	cmd := exec.Command(CmdClient, "savekey", "--privatekey", "0x1234-")
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_SaveKey_Invalid_To_Privatekey_Syntax_Characeter,savekey should return error with privatekey has syntax character")
	}
}

func Test_Client_SaveKey_Invalid_To_FileNameValue_Empty(t *testing.T) {
	cmd := exec.Command(CmdClient, "savekey", "--privatekey", AccountPrivateKey2, "--file", "")
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Test_Client_SaveKey_Invalid_To_FileNameValue_Empty,savekey should return error with empty filename")
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
