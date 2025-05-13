package eth

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"testing"
)

func TestAccountNew(t *testing.T) {
	// 连接到以太坊节点，这里使用 Infura 作为示例，你需要替换为自己的 Infura 项目 ID
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

}
