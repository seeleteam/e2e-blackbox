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

func Test_Client_P2P_NetVersion(t *testing.T) {
	cmd := exec.Command(CmdClient, "p2p", "netversion", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s", err)
	}
}

func Test_Client_P2P_NetworkID(t *testing.T) {
	cmd := exec.Command(CmdClient, "p2p", "networkid", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s", err)
	}
}

func Test_Client_P2P_Peers(t *testing.T) {
	cmd := exec.Command(CmdClient, "p2p", "peers", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s", err)
	}
}

func Test_Client_P2P_ProtocolVersiont(t *testing.T) {
	cmd := exec.Command(CmdClient, "p2p", "protocolversion", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s", err)
	}
}

func Test_Client_P2P_PeersInfo(t *testing.T) {
	cmd := exec.Command(CmdClient, "p2p", "peersinfo", "--address", ServerAddr)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s", err)
	} else {
		var peersInfo []PeerInfo
		if err = json.Unmarshal(output, &peersInfo); err != nil {
			t.Fatalf("%s", err)
		}
		//fmt.Println("peersinfo:", string(output))
		//fmt.Println("infoL:", peersInfo)
	}
}
