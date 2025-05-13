package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TestTransferEthTransaction(t *testing.T) {
	// 连接到以太坊节点，这里使用 Infura 作为示例，你需要替换为自己的 Infura 项目 ID
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	// 发送方私钥，这里需要替换为实际的私钥
	privateKeyStr := "4937011b55f0d01f79ac51b414d22aafba329190c2622c6f3e52172ba08916e5"
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
	// 获取发送方地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Failed to get public key")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 接收方地址，需要替换为实际的接收方地址
	toAddress := common.HexToAddress("0xf93E22f8763f34875B1A2cC7631e899A0c4A9Cef")

	// 转账数量，以 Wei 为单位，1 ETH = 1e18 Wei
	amount := big.NewInt(1e18)

	// 创建交易凭证
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1))
	if err != nil {
		log.Fatalf("Failed to create transactor: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = amount
	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(21000)

	// 创建交易
	tx := types.NewTransaction(nonce, toAddress, amount, auth.GasLimit, auth.GasPrice, nil)

	// 签名交易
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get network ID: %v", err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	// 发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	// 获取交易号（交易哈希）
	txHash := signedTx.Hash().Hex()
	fmt.Printf("Transaction hash: %s\n", txHash)
}
