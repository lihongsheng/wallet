package eth

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"io/ioutil"
	"log"
	"math/big"
	"strings"
	"testing"
)

func TestGetBlockNumberByTxHash(t *testing.T) {
	// 连接到以太坊节点，这里使用 Infura 作为示例，你需要替换为自己的 Infura 项目 ID
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	// 替换为实际的交易哈希
	txHash := common.HexToHash("0xdc57570d9db9b7be4f28432b186de77727dad56609cd182af87e965cf99d4d66")
	// 通过交易哈希获取交易收据
	receipt, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Fatalf("Failed to get transaction receipt: %v", err)
	}
	// 从交易收据中获取区块号
	blockNumber := receipt.BlockNumber
	fmt.Printf("Transaction is in block number: %d\n", blockNumber.Int64())
}

func TestGetLatestBlock(t *testing.T) {
	// 连接到以太坊节点，这里使用 Infura 作为示例，你需要替换为自己的 Infura 项目 ID
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	// 获取最新区块
	//blockInt := big.NewInt(22194332)
	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to get the latest block: %v", err)
	}

	fmt.Printf("Latest block number: %d\n", block.NumberU64())
	fmt.Printf("Block hash: %s\n", block.Hash().Hex())
	fmt.Printf("Block timestamp: %d\n", block.Time())
	fmt.Printf("Block transactions count: %d\n", len(block.Transactions()))

	// 遍历区块中的交易
	for _, tx := range block.Transactions() {
		fmt.Println("--- Transaction ---")
		fmt.Printf("Transaction hash: %s\n", tx.Hash().Hex())
		fmt.Printf("From: %s\n", getTransactionSender(client, tx))
		if tx.To() != nil {
			fmt.Printf("To: %s\n", tx.To().Hex())
		} else {
			fmt.Println("To: Contract creation")
		}
		fmt.Printf("Value: %s ETH\n", convertWeiToEth(tx.Value()))
		fmt.Printf("Gas: %d\n", tx.Gas())
		if tx.Type() == types.DynamicFeeTxType {
			feeCap := tx.GasFeeCap()
			fmt.Printf("Gas Fee Cap: %s Gwei\n", convertWeiToGwei(feeCap))
			tip := tx.GasTipCap()
			fmt.Printf("Gas Tip: %s Gwei\n", convertWeiToGwei(tip))
		} else {
			fmt.Printf("Gas Price: %s Gwei\n", convertWeiToGwei(tx.GasPrice()))
		}
		fmt.Printf("Nonce: %d\n", tx.Nonce())
		fmt.Printf("Nonce: %s\n", tx.Time().String())
	}
}

// 获取交易发送方地址
func getTransactionSender(client *ethclient.Client, tx *types.Transaction) string {
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Printf("Failed to get network ID: %v", err)
		return "Unknown"
	}
	signer := types.MakeSigner(params.MainnetChainConfig, chainID, 1000)
	from, err := types.Sender(signer, tx)
	if err != nil {
		log.Printf("Failed to get transaction sender: %v", err)
		return "Unknown"
	}
	return from.Hex()
}

// 将 Wei 转换为 ETH
func convertWeiToEth(wei *big.Int) string {
	eth := new(big.Int).Div(wei, big.NewInt(1e18))
	return eth.String()
}

// 将 Wei 转换为 Gwei
func convertWeiToGwei(wei *big.Int) string {
	gwei := new(big.Int).Div(wei, big.NewInt(1e9))
	return gwei.String()
}

func TestRead(t *testing.T) {
	// 1. 设置 Keystore 文件路径
	keystorePath := "/Users/lhs/bit/ethereum-multinode/node1/keystore/UTC--2025-05-09T14-54-10.460567000Z--b9041e3a912a6db7bd1f1ac959003f6dde062c34" // 替换为你的 Keystore 文件路径

	// 2. 读取 Keystore 文件内容
	keyJson, err := ioutil.ReadFile(keystorePath)
	if err != nil {
		fmt.Printf("读取 Keystore 文件失败: %v\n", err)
		return
	}

	// 3. 输入密码
	password := "1"

	// 4. 解密私钥
	key, err := keystore.DecryptKey(keyJson, strings.TrimSpace(string(password)))
	if err != nil {
		fmt.Printf("解密失败: %v\n", err)
		return
	}

	// 5. 导出私钥（十六进制格式）
	privateKey := fmt.Sprintf("0x%x", crypto.FromECDSA(key.PrivateKey))
	fmt.Println("私钥:", privateKey)

}
