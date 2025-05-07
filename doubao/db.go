package doubao

// import (
//
//	"crypto/ecdsa"
//	"encoding/hex"
//	"fmt"
//	"github.com/btcsuite/btcd/chaincfg"
//	"log"
//
//	"github.com/btcsuite/btcutil"
//	"github.com/btcsuite/btcutil/hdkeychain"
//	"github.com/ethereum/go-ethereum/crypto"
//	"github.com/miguelmota/go-ethereum-hdwallet"
//	"github.com/tyler-smith/go-bip39"
//
// )
//
// // GenerateMnemonic 生成指定长度的助记词
// // wordCount 参数为助记词的单词数量，支持 12、15、18、21、24
//
//	func GenerateMnemonic(wordCount int) (string, error) {
//		// 根据单词数量计算熵的字节数，每个单词对应 32/3 位熵
//		entropy, err := bip39.NewEntropy(wordCount * 32 / 3)
//		if err != nil {
//			// 如果生成熵时出错，返回错误信息
//			return "", err
//		}
//		// 根据熵生成助记词
//		mnemonic, err := bip39.NewMnemonic(entropy)
//		if err != nil {
//			// 如果生成助记词时出错，返回错误信息
//			return "", err
//		}
//		return mnemonic, nil
//	}
//
// // GenerateSeedFromMnemonic 根据助记词生成种子
// // mnemonic 是助记词，password 是可选的密码
//
//	func GenerateSeedFromMnemonic(mnemonic string, password string) []byte {
//		// 使用助记词和密码生成种子
//		return bip39.NewSeed(mnemonic, password)
//	}
//
// // GenerateBTCKeyAndAddress 生成比特币密钥和地址
// // seed 是生成密钥和地址的种子
//
//	func GenerateBTCKeyAndAddress(seed []byte) (string, string, error) {
//		// 根据种子生成主密钥
//		masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
//		if err != nil {
//			// 如果生成主密钥时出错，返回错误信息
//			return "", "", err
//		}
//		// 派生 BIP44 路径中的 purpose 层，固定为 44'
//		purpose, err := masterKey.Child(hdkeychain.HardenedKeyStart + 44)
//		if err != nil {
//			// 如果派生 purpose 层时出错，返回错误信息
//			return "", "", err
//		}
//		// 派生 BIP44 路径中的 coin_type 层，比特币为 0'
//		coinType, err := purpose.Child(hdkeychain.HardenedKeyStart + 0)
//		if err != nil {
//			// 如果派生 coin_type 层时出错，返回错误信息
//			return "", "", err
//		}
//		// 派生 BIP44 路径中的 account 层，这里使用 0'
//		account, err := coinType.Child(hdkeychain.HardenedKeyStart + 0)
//		if err != nil {
//			// 如果派生 account 层时出错，返回错误信息
//			return "", "", err
//		}
//		// 派生 BIP44 路径中的 change 层，这里使用 0 表示外部地址
//		change, err := account.Child(0)
//		if err != nil {
//			// 如果派生 change 层时出错，返回错误信息
//			return "", "", err
//		}
//		// 派生 BIP44 路径中的 address_index 层，这里使用 0 表示第一个地址
//		addressIndex, err := change.Child(0)
//		if err != nil {
//			// 如果派生 address_index 层时出错，返回错误信息
//			return "", "", err
//		}
//		// 从派生的密钥中获取 ECDSA 私钥
//		privateKey, err := addressIndex.ECPrivKey()
//		if err != nil {
//			// 如果获取私钥时出错，返回错误信息
//			return "", "", err
//		}
//		// 从私钥中获取对应的公钥
//		publicKey := privateKey.PubKey()
//		// 根据公钥生成比特币地址
//		address, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(publicKey.SerializeCompressed()), btcutil.MainNetParams)
//		if err != nil {
//			// 如果生成地址时出错，返回错误信息
//			return "", "", err
//		}
//		// 将私钥转换为 WIF（Wallet Import Format）格式
//		wif, err := btcutil.NewWIF(privateKey, btcutil.MainNetParams, true)
//		if err != nil {
//			// 如果转换为 WIF 格式时出错，返回错误信息
//			return "", "", err
//		}
//		return wif.String(), address.EncodeAddress(), nil
//	}
//
// // GenerateETHKeyAndAddress 生成以太坊密钥和地址
// // seed 是生成密钥和地址的种子
//func GenerateETHKeyAndAddress(seed []byte) (string, string, error) {
//	// 根据种子创建以太坊 HD 钱包
//	wallet, err := hdwallet.NewFromSeed(seed)
//	if err != nil {
//		// 如果创建钱包时出错，返回错误信息
//		return "", "", err
//	}
//	// 解析以太坊的 BIP44 派生路径
//	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
//	// 从钱包中派生指定路径的账户
//	account, err := wallet.Derive(path, false)
//	if err != nil {
//		// 如果派生账户时出错，返回错误信息
//		return "", "", err
//	}
//	// 从账户中获取私钥
//	privateKey, err := wallet.PrivateKey(account)
//	if err != nil {
//		// 如果获取私钥时出错，返回错误信息
//		return "", "", err
//	}
//	// 从私钥中获取对应的公钥
//	publicKey := privateKey.Public()
//	// 将公钥转换为 ECDSA 公钥类型
//	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
//	if !ok {
//		// 如果转换失败，返回错误信息
//		return "", "", fmt.Errorf("error casting public key to ECDSA")
//	}
//	// 根据公钥生成以太坊地址
//	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
//	// 将私钥转换为十六进制字符串
//	privateKeyHex := hex.EncodeToString(crypto.FromECDSA(privateKey))
//	return privateKeyHex, address, nil
//}

//
//func main() {
//	// 示例：生成 12 个单词的助记词
//	wordCount := 12
//	mnemonic, err := GenerateMnemonic(wordCount)
//	if err != nil {
//		// 如果生成助记词时出错，记录错误信息并终止程序
//		log.Fatal(err)
//	}
//	fmt.Printf("生成的助记词: %s\n", mnemonic)
//
//	// 示例：使用助记词生成种子
//	password := ""
//	seed := GenerateSeedFromMnemonic(mnemonic, password)
//
//	// 示例：生成比特币密钥和地址
//	btcPrivateKey, btcAddress, err := GenerateBTCKeyAndAddress(seed)
//	if err != nil {
//		// 如果生成比特币密钥和地址时出错，记录错误信息并终止程序
//		log.Fatal(err)
//	}
//	fmt.Printf("比特币私钥: %s\n", btcPrivateKey)
//	fmt.Printf("比特币地址: %s\n", btcAddress)
//
//	// 示例：生成以太坊密钥和地址
//	ethPrivateKey, ethAddress, err := GenerateETHKeyAndAddress(seed)
//	if err != nil {
//		// 如果生成以太坊密钥和地址时出错，记录错误信息并终止程序
//		log.Fatal(err)
//	}
//	fmt.Printf("以太坊私钥: %s\n", ethPrivateKey)
//	fmt.Printf("以太坊地址: %s\n", ethAddress)
//}
