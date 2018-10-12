/**
*  @file
*  @copyright defined in go-seele/LICENSE
 */

package testcase

var (
	CurShard  int    = 2
	CmdClient string = "../bin/client"
	CmdLight  string = "../bin/light"

	ServerAddr  string = "127.0.0.1:8027"
	AccountErr  string = "0xaaaaaaaaaaaaaaaaa"
	Account1    string = "0x7c00f5a4312a6a3e458a07c2d650ce13c76b68b1" // account for shard1
	AccountMix1 string = "0x7C00F5a4312a6a3e458a07c2d650ce13c76b68B1"
	Account2    string = "0xc910e52e3a314c93fdf545b88d264f39becb8d41" // account for shard2
	AccountMix2 string = "0xc910e52e3a314c93fdf545b88d264f39becb8d41"

	AccountShard2 = [...]string{"0xc910e52e3a314c93fdf545b88d264f39becb8d41", "0xff0fb1e59e92e94fac74febec98cfd58b956fa61"}

	BlockHash    string = "0x02ffd2ac0921f7784d0d3f6d7865e48ec0e454dbd7dc60e4ecf6eaa08c548410"
	BlockHashErr string = "0x88aad2ac0921f7784d0d3f6d7865e48ec0e454dbd7dc60e4ecf6eaa08c548410"
)

/*
func main() {
	fmt.Println("hello owr")

	cmd := exec.Command("./node", "--help")
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out), err)

}
*/
