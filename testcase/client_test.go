/**
*  @file
*  @copyright defined in go-seele/LICENSE
 */

package testcase

import (
	"encoding/json"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"
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

func Test_Client_GetBlock_ByNormalHeight(t *testing.T) {
	// Normal height
ErrContinue:
	cmd := exec.Command(CmdClient, "getblock", "--height", "1", "--address", ServerAddr)
	if res, err := cmd.CombinedOutput(); err != nil {
		t.Logf("Test_Client_GetBlock_ByNormalHeight: Node not running!")
		time.Sleep(5 * time.Second)
		goto ErrContinue
	} else {
		var blockInfo BlockInfo
		if err = json.Unmarshal(res, &blockInfo); err != nil {
			t.Fatalf("Test_Client_GetBlock_ByNormalHeight: %s", err)
		}

		headerMp := blockInfo.Header.(map[string]interface{})
		height := uint64(headerMp["Height"].(float64))
		if height != 1 {
			t.Fatalf("Test_Client_GetBlock_ByNormalHeight: Expect the return value is not correct!")
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

		headerMp := blockInfo.Header.(map[string]interface{})
		height := uint64(headerMp["Height"].(float64))
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
		if blockInfo.Hash != "0x000000c80d93ee588e022dfa396357c6f1f77b4f8576d8188dcbb821ed742900" {
			t.Fatalf("Test_Client_GetBlock_ByHashFulltx: Expect the return value is not correct!")
		}
		if len(blockInfo.Transactions) <= 0 {
			t.Fatalf("Test_Client_GetBlock_ByHashFulltx, hash should contain one transaction at lease")
		}
	}
}
