package hdkey

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

func GenerateEth(seed []byte) (privateKey string, address string, err error) {
	// 构建 BIP44 路径 m/44'/{coin}'/{account}'/{change}/{index}
	path := accounts.DefaultBaseDerivationPath
	path[0] = 44 + hdkeychain.HardenedKeyStart
	path[1] = uint32(60) + hdkeychain.HardenedKeyStart
	// account
	path[2] = 0 + hdkeychain.HardenedKeyStart
	// change
	path[3] = 0
	// index
	path[4] = 0
	// by seed generate private key and address
	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		// if there is an error generating the wallet, return the error
		return "", "", err
	}
	// bip44  path m/44'/60'/0'/0/0
	// by wallet and path generate account
	accountObj, err := wallet.Derive(path, false)
	if err != nil {
		// if there is an error deriving the account, return the error
		return "", "", err
	}
	// by account get private key
	privateKeyObj, err := wallet.PrivateKey(accountObj)
	if err != nil {
		// 如果获取私钥时出错，返回错误信息
		return "", "", err
	}
	// 从私钥中获取对应的公钥
	publicKey := privateKeyObj.Public()
	// try to convert public key to ECDSA
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		// if public key cannot be cast to ECDSA, return error
		return "", "", errors.New("error casting public key to ECDSA")
	}
	// by public key get address
	address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	// private key to hex string
	privateKeyHex := hex.EncodeToString(crypto.FromECDSA(privateKeyObj))
	return privateKeyHex, address, nil
}

//
// GetETHKeyPair 获取以太坊密钥对
//func (kg *KeyGenerator) GetETHKeyPair() (privateKey string, address string, err error) {
//   masterKey, err := hdkeychain.NewMaster(kg.seed, &chaincfg.MainNetParams)
//	// 获取扩展密钥的私钥（返回 btcec.PrivateKey 类型）
//	privKey, err := masterKey.ECPrivKey()
//	if err != nil {
//		return "", "", err
//	}
//	// 将 btcec.PrivateKey 转换为标准 ECDSA 私钥
//	ethPrivKey := privKey.ToECDSA()
//	// 生成以太坊地址
//	publicKey := ethPrivKey.PublicKey
//	address = crypto.PubkeyToAddress(publicKey).Hex()
//
//	// 将私钥编码为 HEX 字符串
//	privKeyBytes := crypto.FromECDSA(ethPrivKey)
//	return hex.EncodeToString(privKeyBytes), address, nil
//}
