package btc

import (
	"fmt"
	"log"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
)

func main() {
	// 连接到比特币节点，需替换为你的节点信息
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:8332",
		User:         "your_username",
		Pass:         "your_password",
		HTTPPostMode: true,
		DisableTLS:   true,
	}
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatalf("Failed to connect to Bitcoin node: %v", err)
	}
	defer client.Shutdown()

	// 替换为实际的交易哈希
	txHashStr := "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	txHash, err := chainhash.NewHashFromStr(txHashStr)
	if err != nil {
		log.Fatalf("Failed to parse transaction hash: %v", err)
	}

	// 查询交易详细信息
	tx, err := client.GetTransaction(txHash)
	if err != nil {
		log.Fatalf("Failed to get transaction details: %v", err)
	}

	// 验证交易是否成功
	if tx.Confirmations > 0 {
		fmt.Printf("Transaction is confirmed with %d confirmations.\n", tx.Confirmations)
		fmt.Println("Transaction is successful.")
	} else {
		fmt.Println("Transaction is not confirmed yet.")
	}
}
