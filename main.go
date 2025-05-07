package main

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
	"log"
)

// 定义支持的币种类型
type CoinType uint32

const (
	BTC CoinType = 0  // BIP44 比特币路径 m/44'/0'/0'/0/...
	ETH CoinType = 60 // BIP44 以太坊路径 m/44'/60'/0'/0/...
)

// KeyGenerator 密钥生成器结构体
type KeyGenerator struct {
	mnemonic       string                  // BIP39 助记词
	seed           []byte                  // BIP32 种子
	masterKey      *hdkeychain.ExtendedKey // 主私钥
	derivationPath string                  // 派生路径
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
	return mnemonic, nil
}

// ImportMnemonic 导入已有助记词
func (kg *KeyGenerator) ImportMnemonic(mnemonic string) error {
	if !bip39.IsMnemonicValid(mnemonic) {
		return fmt.Errorf("invalid mnemonic")
	}
	kg.mnemonic = mnemonic
	return nil
}

// GenerateSeed 生成种子 (BIP39)
// passphrase 是可选的密码
func (kg *KeyGenerator) GenerateSeed(passphrase string) error {
	seed := bip39.NewSeed(kg.mnemonic, passphrase)
	kg.seed = seed
	return nil
}

// CreateMasterKey 创建主私钥 (BIP32)
func (kg *KeyGenerator) CreateMasterKey() error {
	masterKey, err := hdkeychain.NewMaster(kg.seed, &chaincfg.MainNetParams)
	if err != nil {
		return err
	}
	kg.masterKey = masterKey
	return nil
}

// DeriveKey 派生指定路径密钥
func (kg *KeyGenerator) DeriveKey(coin CoinType, account, change, index uint32) error {
	// 构建 BIP44 路径 m/44'/{coin}'/{account}'/{change}/{index}
	path := accounts.DefaultBaseDerivationPath
	path[0] = 44 + hdkeychain.HardenedKeyStart
	path[1] = uint32(coin) + hdkeychain.HardenedKeyStart
	path[2] = account + hdkeychain.HardenedKeyStart
	path[3] = change
	path[4] = index

	// 逐层派生密钥
	currentKey := kg.masterKey
	for _, n := range path {
		childKey, err := currentKey.Derive(n)
		if err != nil {
			return err
		}
		currentKey = childKey
	}

	kg.derivationPath = fmt.Sprintf("%s", path.String())
	return nil
}

// GetBTCKeyPair 获取比特币密钥对
func (kg *KeyGenerator) GetBTCKeyPair() (wif string, address string, err error) {
	// 导出私钥
	privKey, err := kg.masterKey.ECPrivKey()
	if err != nil {
		return "", "", err
	}
	// 生成 WIF 格式私钥
	wifObj, err := btcutil.NewWIF(privKey, &chaincfg.MainNetParams, true)
	if err != nil {
		return "", "", err
	}

	// 生成 P2PKH 地址
	pubKey, err := kg.masterKey.ECPubKey()
	if err != nil {
		return "", "", err
	}
	hash160 := btcutil.Hash160(pubKey.SerializeCompressed())
	addressObj, err := btcutil.NewAddressPubKeyHash(hash160, &chaincfg.MainNetParams)
	if err != nil {
		return "", "", err
	}

	return wifObj.String(), addressObj.String(), nil
}

// GetETHKeyPair 获取以太坊密钥对
func (kg *KeyGenerator) GetETHKeyPair() (privateKey string, address string, err error) {
	// 获取扩展密钥的私钥（返回 btcec.PrivateKey 类型）
	privKey, err := kg.masterKey.ECPrivKey()
	if err != nil {
		return "", "", err
	}
	// 将 btcec.PrivateKey 转换为标准 ECDSA 私钥
	ethPrivKey := privKey.ToECDSA()
	// 生成以太坊地址
	publicKey := ethPrivKey.PublicKey
	address = crypto.PubkeyToAddress(publicKey).Hex()

	// 将私钥编码为 HEX 字符串
	privKeyBytes := crypto.FromECDSA(ethPrivKey)
	return hex.EncodeToString(privKeyBytes), address, nil
}

// 示例使用
func main() {
	kg := NewKeyGenerator()

	// 生成助记词（示例使用12个单词）
	mnemonic, err := kg.GenerateMnemonic(12)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generated Mnemonic: %s\n", mnemonic)

	// 生成种子
	if err := kg.GenerateSeed(""); err != nil { // 空密码
		log.Fatal(err)
	}

	// 创建主密钥
	if err := kg.CreateMasterKey(); err != nil {
		log.Fatal(err)
	}

	// 派生比特币密钥
	if err := kg.DeriveKey(BTC, 0, 0, 0); err != nil {
		log.Fatal(err)
	}
	btcWIF, btcAddr, err := kg.GetBTCKeyPair()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("BTC WIF: %s\nAddress: %s\n %s", btcWIF, btcAddr, kg.derivationPath)

	// 派生以太坊密钥
	if err := kg.DeriveKey(ETH, 0, 0, 0); err != nil {
		log.Fatal(err)
	}
	ethPriv, ethAddr, err := kg.GetETHKeyPair()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ETH Private: 0x%s\nAddress: %s\n %s", ethPriv, ethAddr, kg.derivationPath)
}
