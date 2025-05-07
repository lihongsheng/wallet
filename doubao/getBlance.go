package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

// USDT合约地址（USDT - ERC20）
const usdtContractAddress = "0xdAC17F958D2ee523a2206206994597C13D831ec7"

func main() {
	// 连接到以太坊节点，使用 Cloudflare 的以太坊主网节点
	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	// 要查询的以太坊地址
	address := common.HexToAddress("0xYourEthAddress")

	// 获取ETH余额
	ethBalance, err := getETHBalance(client, address)
	if err != nil {
		log.Fatalf("Failed to get ETH balance: %v", err)
	}
	fmt.Printf("ETH balance: %.18f ETH\n", float64(ethBalance.Int64())/float64(params.Ether))

	// 获取USDT余额
	usdtBalance, err := getUSDTBalance(client, address)
	if err != nil {
		log.Fatalf("Failed to get USDT balance: %v", err)
	}
	fmt.Printf("USDT balance: %.6f USDT\n", float64(usdtBalance.Int64())/1e6)
}

// getETHBalance 获取指定地址的ETH余额
func getETHBalance(client *ethclient.Client, address common.Address) (*big.Int, error) {
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

// getUSDTBalance 获取指定地址的USDT余额
func getUSDTBalance(client *ethclient.Client, address common.Address) (*big.Int, error) {
	contractAddress := common.HexToAddress(usdtContractAddress)
	contract, err := NewTokenCaller(contractAddress, client)
	if err != nil {
		return nil, err
	}
	balance, err := contract.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		return nil, err
	}
	return balance, nil
}
