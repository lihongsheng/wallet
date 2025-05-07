package main

import (
	"fmt"
	"log"
	"wallet/pkg/enum"
	"wallet/pkg/hdkey"
)

// 示例使用
func main() {
	kg := hdkey.NewKeyGenerator()

	// 生成助记词（示例使用12个单词）
	err := kg.ImportMnemonic("trick horror alert giggle egg share tree wing favorite quarter squeeze lawsuit")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generated Mnemonic: %s\n", kg.GetMnemonic())
	// 生成种子
	if err := kg.GenerateSeed(""); err != nil { // 空密码
		log.Fatal(err)
	}

	btcWIF, btcAddr, err := kg.GenerateKey(enum.BTC)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("BTC WIF: %s\nAddress: %s\n", btcWIF, btcAddr)

	ethPriv, ethAddr, err := kg.GenerateKey(enum.ETH)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ETH Private: 0x%s\nAddress: %s\n", ethPriv, ethAddr)
}
