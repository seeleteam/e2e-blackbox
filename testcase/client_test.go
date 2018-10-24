/**
*  @file
*  @copyright defined in go-seele/LICENSE
 */

package testcase

import (
	"encoding/json"
	"os/exec"
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
		t.Fatalf("getinfo error, %s", err)
	}

	var r ResGetInfo
	r_err := json.Unmarshal(res, &r)
	if r_err != nil {
		t.Fatal("decode return result error")
	}

	if r.MinerStatus != "Running" {
		t.Fatal("Node not running!")
	}
}
