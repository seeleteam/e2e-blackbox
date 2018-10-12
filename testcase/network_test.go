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

type PeerInfo struct {
	ID      string   `json:"id"`   // Unique of the node
	Caps    []string `json:"caps"` // Sum-protocols advertised by this particular peer
	Network struct {
		LocalAddress  string `json:"localAddress"`  // Local endpoint of the TCP data connection
		RemoteAddress string `json:"remoteAddress"` // Remote endpoint of the TCP data connection
	} `json:"network"`
	Protocols map[string]interface{} `json:"protocols"` // Sub-protocol specific metadata fields
	Shard     uint                   `json:"shard"`     // shard id of the node
}

func Test_Client_P2P(t *testing.T) {
	cmd := exec.Command(CmdClient, "p2p", "netversion", "--address", ServerAddr)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s", err)
	} else {
		fmt.Println("p2p netversion=", string(output))
	}

	cmd = exec.Command(CmdClient, "p2p", "networkid", "--address", ServerAddr)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s", err)
	} else {
		fmt.Println("p2p networkid=", string(output))
	}

	cmd = exec.Command(CmdClient, "p2p", "peers", "--address", ServerAddr)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s", err)
	} else {
		fmt.Println("p2p peers=", string(output))
	}

	cmd = exec.Command(CmdClient, "p2p", "protocolversion", "--address", ServerAddr)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s", err)
	} else {
		fmt.Println("p2p protocolversion=", string(output))
	}

	cmd = exec.Command(CmdClient, "p2p", "peersinfo", "--address", ServerAddr)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s", err)
	} else {
		var peersInfo []PeerInfo
		if err = json.Unmarshal(output, &peersInfo); err != nil {
			t.Fatalf("%s", err)
		} else {
			fmt.Println("p2p len(peersinfo)=", len(peersInfo))
		}
	}
}
