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

const (
	// shard file name
	KeyFileShard1_1 = "../bin/config/keyfile/shard1-0x0a57a2714e193b7ac50475ce625f2dcfb483d741"
	KeyFileShard1_2 = "../bin/config/keyfile/shard1-0x2a87b6504cd00af95a83b9887112016a2a991cf1"
	KeyFileShard1_3 = "../bin/config/keyfile/shard1-0x3b691130ec4166bfc9ec7240217fc8d08903cf21"
	KeyFileShard1_4 = "../bin/config/keyfile/shard1-0x4fb7c8b0287378f0cf8b5a9262bf3ef7e101f8d1"
	KeyFileShard1_5 = "../bin/config/keyfile/shard1-0xec759db47a65f6537d630517f6cd3ca39c6f93d1"
	KeyFileShard2_1 = "../bin/config/keyfile/shard2-0x2a23825407740fa7163069257c57452c4d4fc3d1"
	KeyFileShard2_2 = "../bin/config/keyfile/shard2-0x4eea165e9266f20bf6e5e08e0c11d38e8fc02661"
	KeyFileShard2_3 = "../bin/config/keyfile/shard2-0x007d1b1ea335e8e4a74c0be781d828dc7db934b1"
	KeyFileShard2_4 = "../bin/config/keyfile/shard2-0xfaf78f23293cc537154c275c874ede0f8c8b8801"
	KeyFileShard2_5 = "../bin/config/keyfile/shard2-0xfbe506bdaf256682551873290d0a794d51bac4d1"

	// accounts corresponding to keyFileShard above
	AccountShard1_1 = "0x0a57a2714e193b7ac50475ce625f2dcfb483d741"
	AccountShard1_2 = "0x2a87b6504cd00af95a83b9887112016a2a991cf1"
	AccountShard1_3 = "0x3b691130ec4166bfc9ec7240217fc8d08903cf21"
	AccountShard1_4 = "0x4fb7c8b0287378f0cf8b5a9262bf3ef7e101f8d1"
	AccountShard1_5 = "0xec759db47a65f6537d630517f6cd3ca39c6f93d1"
	AccountShard2_1 = "0x2a23825407740fa7163069257c57452c4d4fc3d1"
	AccountShard2_2 = "0x4eea165e9266f20bf6e5e08e0c11d38e8fc02661"
	AccountShard2_3 = "0x007d1b1ea335e8e4a74c0be781d828dc7db934b1"
	AccountShard2_4 = "0xfaf78f23293cc537154c275c874ede0f8c8b8801"
	AccountShard2_5 = "0xfbe506bdaf256682551873290d0a794d51bac4d1"

	// config path
	ConfigPath = "../bin/config"

	// account base balance
	BaseBalance = 1000000000000

	// HTLC secret and secret hash
	Secret       = "0x31aa0be185fbc89048a0381cc5136969565be9d9c13048f7c2ee9322811b99fc"
	ForgedSecret = "0xc5543fa77c58024c27879360b1fcd3fa67f546c3ebdc5f3598c30d10266e2830"
	Secretehash  = "0x57e685963f607851af252e7922483a61fbceced12accd745444f412295517768"
)
