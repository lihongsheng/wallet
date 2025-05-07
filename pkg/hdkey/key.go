package hdkey

import (
	"errors"
	"fmt"
	"github.com/tyler-smith/go-bip39"
	"wallet/pkg/enum"
)

type KeyGenerator struct {
	mnemonic string // BIP39 助记词
	seed     []byte // BIP32 种子
}

// NewKeyGenerator 创建新实例
func NewKeyGenerator() *KeyGenerator {
	return &KeyGenerator{}
}

// GenerateMnemonic 生成助记词 (支持12/15/18/21/24个单词)
// GenerateMnemonic 生成指定长度的助记词
// wordCount 参数为助记词的单词数量，支持 12、15、18、21、24
func (kg *KeyGenerator) GenerateMnemonic(wordCount int) (string, error) {
	// 根据单词数量计算熵的字节数，每个单词对应 32/3 位熵
	entropy, err := bip39.NewEntropy(wordCount * 32 / 3)
	if err != nil {
		// 如果生成熵时出错，返回错误信息
		return "", err
	}
	// 根据熵生成助记词
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		// 如果生成助记词时出错，返回错误信息
		return "", err
	}
	kg.mnemonic = mnemonic
	return mnemonic, nil
}

// ImportMnemonic 导入已有助记词
func (kg *KeyGenerator) ImportMnemonic(mnemonic string) error {
	if !bip39.IsMnemonicValid(mnemonic) {
		return errors.New("invalid mnemonic")
	}
	kg.mnemonic = mnemonic
	return nil
}
func (kg *KeyGenerator) GetMnemonic() string {
	return kg.mnemonic
}

// GenerateSeed 生成种子 (BIP39)
// passphrase 是可选的密码
func (kg *KeyGenerator) GenerateSeed(passphrase string) error {
	seed := bip39.NewSeed(kg.mnemonic, passphrase)
	kg.seed = seed
	return nil
}

func (kg *KeyGenerator) GenerateKey(coin enum.CoinType) (privateKey string, address string, err error) {
	switch coin {
	case enum.ETH:
		return GenerateEth(kg.seed)
	case enum.BTC:
		return GenerateBtc(kg.seed)
	}
	return "", "", errors.New(fmt.Sprintf("not find coin[%s]", coin))
}
