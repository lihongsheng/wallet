package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"wallet/transfer/eth/abi" // 引入生成的 USDT 绑定代码
)

func testTransferUSDT() {
	// 连接到以太坊节点，这里使用 Infura 作为示例，你需要替换为自己的 Infura 项目 ID
	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	// USDT 合约地址
	usdtContractAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")

	// 加载 USDT 合约
	token, err := abi.NewUSDT(usdtContractAddress, client)
	if err != nil {
		log.Fatalf("Failed to load USDT contract: %v", err)
	}

	// 发送方私钥，这里需要替换为实际的私钥
	privateKeyStr := "YOUR_PRIVATE_KEY"
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// 创建交易凭证
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Failed to get public key")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1)) // 以太坊主网链 ID 为 1
	if err != nil {
		log.Fatalf("Failed to create transactor: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // 不发送 ETH，仅转账 USDT
	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(210000) // 可以根据实际情况调整

	// 接收方地址，需要替换为实际的接收方地址
	toAddress := common.HexToAddress("RECIPIENT_ADDRESS")

	// 转账数量，注意 USDT 有 6 位小数，这里以 1 USDT 为例
	amount := big.NewInt(1000000)

	// 发起转账
	tx, err := token.Transfer(auth, toAddress, amount)
	if err != nil {
		log.Fatalf("Failed to transfer USDT: %v", err)
	}

	fmt.Printf("Transaction hash: %s\n", tx.Hash().Hex())
}
