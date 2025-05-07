package hdkey

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
)

func GenerateBtc(seed []byte) (privateKey string, address string, err error) {
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return "", "", err
	}
	// 构建 BIP44 路径 m/44'/{coin}'/{account}'/{change}/{index}
	// m / 44'
	purposeKey, err := masterKey.Derive(hdkeychain.HardenedKeyStart + 44)
	if err != nil {
		return "", "", err
	}
	// m / 44' / 0'
	coinTypeKey, err := purposeKey.Derive(hdkeychain.HardenedKeyStart + 0)
	if err != nil {
		return "", "", err
	}
	// m / 44' / 0' / account'
	accountKey, err := coinTypeKey.Derive(hdkeychain.HardenedKeyStart + 0)
	if err != nil {
		return "", "", err
	}
	// m / 44' / 0' / account' / 0
	changeKey, err := accountKey.Derive(0)
	if err != nil {
		return "", "", err
	}
	// m / 44' / 0' / account' / 0 / index
	addressKey, err := changeKey.Derive(0)
	if err != nil {
		return "", "", err
	}
	// 从派生的密钥中获取 ECDSA 私钥
	privateKeyObj, err := addressKey.ECPrivKey()
	if err != nil {
		// 如果获取私钥时出错，返回错误信息
		return "", "", err
	}
	// 从私钥中获取对应的公钥
	publicKey := privateKeyObj.PubKey()
	// 根据公钥生成比特币地址
	addressObj, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(publicKey.SerializeCompressed()), &chaincfg.MainNetParams)
	if err != nil {
		// 如果生成地址时出错，返回错误信息
		return "", "", err
	}
	// 将私钥转换为 WIF（Wallet Import Format）格式
	wif, err := btcutil.NewWIF(privateKeyObj, &chaincfg.MainNetParams, true)
	if err != nil {
		// 如果转换为 WIF 格式时出错，返回错误信息
		return "", "", err
	}
	return wif.String(), addressObj.EncodeAddress(), nil
}

// GetBTCKeyPair() (wif string, address string, err error) {
//   	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
//	// 导出私钥
//	privKey, err := masterKey.ECPrivKey()
//	if err != nil {
//		return "", "", err
//	}
//	// 生成 WIF 格式私钥
//	wifObj, err := btcutil.NewWIF(privKey, &chaincfg.MainNetParams, true)
//	if err != nil {
//		return "", "", err
//	}
//
//	// 生成 P2PKH 地址
//	pubKey, err := kg.masterKey.ECPubKey()
//	if err != nil {
//		return "", "", err
//	}
//	hash160 := btcutil.Hash160(pubKey.SerializeCompressed())
//	addressObj, err := btcutil.NewAddressPubKeyHash(hash160, &chaincfg.MainNetParams)
//	if err != nil {
//		return "", "", err
//	}
//
//	return wifObj.String(), addressObj.String(), nil
//}
