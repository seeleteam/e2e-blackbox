/**
*  @file
*  @copyright defined in go-seele/LICENSE
 */

package testcase

var (
	CurShard  int    = 2
	CmdClient string = "../bin/client"
	CmdLight  string = "../bin/light"

	ServerAddr string = "127.0.0.1:8027"
	AccountErr string = "0xaaaaaaaaaaaaaaaaa"

	Account1_Aux string = "0x7c00f5a4312a6a3e458a07c2d650ce13c76b68b1"
	Account1     string = "0xa00d22dc3624d4696eff8d1641b442f79c3379b1" // account for shard1

	AccountMix1        string = "0xA00D22dc3624d4696eff8d1641b442f79c3379b1"
	Account2           string = "0xc910e52e3a314c93fdf545b88d264f39becb8d41" // account for shard2
	AccountMix2        string = "0xc910e52e3a314c93fdf545b88d264f39becb8d41"
	AccountPrivateKey2 string = "0x9b9245066c57a5cd376a378b9edc69ea545a195771d5f55859180f1a2ff61240" //private key for account2

	InvalidAccountType string = "0xff0fb1e59e92e94fac74febec98cfd58b956fa6f" // account type == 15, invalid
	AccountShard2             = [...]string{"0xc910e52e3a314c93fdf545b88d264f39becb8d41", "0xff0fb1e59e92e94fac74febec98cfd58b956fa61"}

	BlockHash    string = "0x000000c80d93ee588e022dfa396357c6f1f77b4f8576d8188dcbb821ed742900"
	BlockHashErr string = "0x88aad2ac0921f7784d0d3f6d7865e48ec0e454dbd7dc60e4ecf6eaa08c548410"
)
