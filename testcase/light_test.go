/**
*  @file
*  @copyright defined in go-seele/LICENSE
 */

package testcase

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"testing"
)

type BlockInfo struct {
	Hash         string        `json:"hash"`
	Transactions []interface{} `json:"transactions"`
}

func accountCase(command, account, accountMix string, t *testing.T) {
	cmd := exec.Command(CmdLight, "getbalance", "--account", account, "--address", ServerAddr)
	var output, outputMix []byte
	var err error
	if output, err = cmd.CombinedOutput(); err != nil {
		t.Fatalf("getbalance err: %s", err)
	} else {
		fmt.Println("light getbalance =", account, string(output))
	}

	cmd = exec.Command(CmdLight, "getbalance", "--account", accountMix, "--address", ServerAddr)
	if outputMix, err = cmd.CombinedOutput(); err != nil {
		t.Fatalf("getbalance err: %s", err)
	} else {
		fmt.Println("light getbalance =", accountMix, string(outputMix))
	}

	if string(output) != string(outputMix) {
		t.Fail()
	}
}

// account should ignore character case.
func Test_Light_AccountCase(t *testing.T) {
	cmd := exec.Command(CmdLight, "getbalance", "--account", AccountErr, "--address", ServerAddr)
	if output, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getbalance AccountErr success?")
	} else {
		fmt.Println("light getbalance =", string(output))
	}

	accountCase(CmdLight, Account2, AccountMix2, t)
}

func Test_Light_GetErrBlock(t *testing.T) {

	// invalid height
	cmd := exec.Command(CmdLight, "getblock", "--height", "100000000", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getblock error parameter success?")
	} else {
		fmt.Println("light getblock for invalid parameter. err =", err)
	}

	cmd = exec.Command(CmdLight, "getblock", "--hash", BlockHashErr, "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getblock error parameter success?")
	} else {
		fmt.Println("light getblock for invalid parameter. err =", err)
	}

	cmd = exec.Command(CmdLight, "getblock", "--hash", "0x", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getblock error parameter success?")
	} else {
		fmt.Println("light getblock for invalid parameter. err =", err)
	}

	cmd = exec.Command(CmdLight, "getblock", "1", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getblock error parameter success?")
	} else {
		fmt.Println("light getblock for invalid parameter. err =", err)
	}
}

func Test_Light_GetBlock(t *testing.T) {
	cmd := exec.Command(CmdLight, "getblock", "--height", "0", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("getblock error, %s", err)
	} else {
		fmt.Println("light getblock by height ok.")
	}

	cmd = exec.Command(CmdLight, "getblock", "--hash", BlockHash, "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("getblock error, %s", err)
	} else {
		fmt.Println("light getblock by hash ok.")
	}

	// getblock fulltx support.
	cmd = exec.Command(CmdLight, "getblock", "--height", "1", "--fulltx", "--address", ServerAddr)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("getblock error, %s", err)
	} else {
		var blockInfo BlockInfo
		//fmt.Println(string(output))
		if err = json.Unmarshal(output, &blockInfo); err != nil {
			t.Fatalf("%s", err)
		} else {
			fmt.Println("p2p len(blockInfos.transactions)=", len(blockInfo.Transactions))
		}
	}

}

func Test_Light_GetBlockHeight(t *testing.T) {
	cmd := exec.Command(CmdLight, "getblockheight", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("getblockheight error, %s", err)
	} else {
		fmt.Println("light getblockheight ok.")
	}

	cmd = exec.Command(CmdLight, "getblockheight", "100", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getblockheight returns ok with invalid parameter")
	} else {
		fmt.Println("light getblockheight invalid parameter test ok.")
	}
}

func Test_Light_GetErrBlockTXCount(t *testing.T) {
	// invalid height
	cmd := exec.Command(CmdLight, "getblocktxcount", "--height", "100000000", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getblocktxcount error parameter success?")
	} else {
		fmt.Println("light getblocktxcount for invalid parameter. err =", err)
	}

	cmd = exec.Command(CmdLight, "getblocktxcount", "--hash", BlockHashErr, "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getblocktxcount error parameter success?")
	} else {
		fmt.Println("light getblocktxcount for invalid parameter. err =", err)
	}

	cmd = exec.Command(CmdLight, "getblocktxcount", "--hash", "0x", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getblocktxcount error parameter success?")
	} else {
		fmt.Println("light getblocktxcount for invalid parameter. err =", err)
	}

	cmd = exec.Command(CmdLight, "getblocktxcount", "1", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("getblock error parameter success?")
	} else {
		fmt.Println("light getblock for invalid parameter. err =", err)
	}
}

func Test_Light_GetBlockTXCount(t *testing.T) {
	cmd := exec.Command(CmdLight, "getblocktxcount", "--height", "0", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("getblocktxcount error, %s", err)
	} else {
		fmt.Println("light getblocktxcount by height ok.")
	}

	cmd = exec.Command(CmdLight, "getblocktxcount", "--hash", BlockHash, "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("getblocktxcount error, %s", err)
	} else {
		fmt.Println("light getblocktxcount by hash ok.")
	}
}
