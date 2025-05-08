package btc

import (
	"fmt"
	"log"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
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

	// 上一次拉取的区块号
	lastKnownBlockNumber := int64(123456)

	// 获取最新的区块号
	latestBlockNumber, err := client.GetBlockCount()
	if err != nil {
		log.Fatalf("Failed to get latest block number: %v", err)
	}

	fmt.Printf("Last known block number: %d\n", lastKnownBlockNumber)
	fmt.Printf("Latest block number: %d\n", latestBlockNumber)

	// 要查询余额的比特币地址，需替换为实际地址
	addressStr := "1ExampleAddress1234567890"
	address, err := btcutil.DecodeAddress(addressStr, &chaincfg.MainNetParams)
	if err != nil {
		log.Fatalf("Failed to decode address: %v", err)
	}

	// 获取该地址在指定区块范围内的未花费交易输出（UTXO）
	utxos, err := client.ListUnspentMinMaxAddresses(int64(lastKnownBlockNumber), latestBlockNumber, []btcutil.Address{address})
	if err != nil {
		log.Fatalf("Failed to list unspent outputs: %v", err)
	}

	// 计算余额
	balance := btcutil.Amount(0)
	for _, utxo := range utxos {
		balance += utxo.Amount
	}

	fmt.Printf("Address %s balance: %s BTC\n", addressStr, balance.String())
}
