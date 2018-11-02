/**
*  @file
*  @copyright defined in go-seele/LICENSE
 */

package testcase

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strconv"
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
		t.Fatalf("Test_Client_P2P_NetVersion: An error occured: %s", err.Error())
	}
}

func Test_Client_P2P_NetworkID(t *testing.T) {
	cmd := exec.Command(CmdClient, "p2p", "networkid", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_P2P_NetworkID: An error occured: %s", err.Error())
	}
}

func Test_Client_P2P_Peers(t *testing.T) {
	cmd := exec.Command(CmdClient, "p2p", "peers", "--address", ServerAddr)
	peers, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Test_Client_P2P_Peers: An error occured: %s", err.Error())
	}
	peers = bytes.TrimRight(peers, "\n")
	n, err := strconv.ParseInt(string(peers), 10, 64)
	if err != nil {
		t.Fatalf("parse string to int64 failed %s", err.Error())
	}
	if n < 0 {
		t.Fatalf("Test_Client_P2P_Peers returns invalid peer number")
	}
}

func Test_Client_P2P_ProtocolVersion(t *testing.T) {
	cmd := exec.Command(CmdClient, "p2p", "protocolversion", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_P2P_ProtocolVersion: An error occured: %s", err.Error())
	}
}

func Test_Client_P2P_PeersInfo(t *testing.T) {
	cmd := exec.Command(CmdClient, "p2p", "peersinfo", "--address", ServerAddr)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Client_P2P_PeersInfo: An error occured: %s", err.Error())
	} else {
		var peersInfo []PeerInfo
		if err = json.Unmarshal(output, &peersInfo); err != nil {
			t.Fatalf("Test_Client_P2P_PeersInfo ummarshal failed %s", err)
		}
	}
}

func Test_Light_P2P_NetVersion(t *testing.T) {
	cmd := exec.Command(CmdLight, "p2p", "netversion", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Light_P2P_NetVersion: An error occured: %s", err.Error())
	}
}

func Test_Light_P2P_NetworkID(t *testing.T) {
	cmd := exec.Command(CmdLight, "p2p", "networkid", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Light_P2P_NetworkID: An error occured: %s", err.Error())
	}
}

func Test_Light_P2P_Peers(t *testing.T) {
	cmd := exec.Command(CmdLight, "p2p", "peers", "--address", ServerAddr)
	peers, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Test_Client_P2P_Peers: An error occured: %s", err.Error())
	}
	peers = bytes.TrimRight(peers, "\n")
	n, err := strconv.ParseInt(string(peers), 10, 64)
	if err != nil {
		t.Fatalf("parse string to int64 failed %s", err.Error())
	}
	if n < 0 {
		t.Fatalf("Test_Light_P2P_Peers returns invalid peer number")
	}
}

func Test_Light_P2P_ProtocolVersion(t *testing.T) {
	cmd := exec.Command(CmdLight, "p2p", "protocolversion", "--address", ServerAddr)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Light_P2P_ProtocolVersion: An error occured: %s", err.Error())
	}
}

func Test_Light_P2P_PeersInfo(t *testing.T) {
	cmd := exec.Command(CmdLight, "p2p", "peersinfo", "--address", ServerAddr)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Test_Light_P2P_PeersInfo: An error occured: %s", err.Error())
	} else {
		var peersInfo []PeerInfo
		if err = json.Unmarshal(output, &peersInfo); err != nil {
			t.Fatalf("Test_Client_P2P_PeersInfo ummarshal failed %s", err)
		}
	}
}
