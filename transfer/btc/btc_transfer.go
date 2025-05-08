package btc

import (
	"fmt"
	"log"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func main() {
	// 连接到比特币节点，需替换为你的节点信息
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:8332",
		User:         "your_username",
		Pass:         "your_password",
		HTTPPostMode: true,
		DisableTLS:   true,
	}
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatalf("Failed to connect to Bitcoin node: %v", err)
	}
	defer client.Shutdown()

	// 生成或加载私钥，这里简单示例，实际应用需妥善管理私钥
	seed := make([]byte, hdkeychain.RecommendedSeedLen)
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		log.Fatalf("Failed to create master key: %v", err)
	}
	childKey, err := masterKey.Derive(hdkeychain.HardenedKeyStart + 0)
	if err != nil {
		log.Fatalf("Failed to derive child key: %v", err)
	}
	privateKey, err := childKey.ECPrivKey()
	if err != nil {
		log.Fatalf("Failed to get private key: %v", err)
	}
	publicKey := privateKey.PubKey()

	// 生成发送方地址
	addressPubKeyHash, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(publicKey.SerializeCompressed()), &chaincfg.MainNetParams)
	if err != nil {
		log.Fatalf("Failed to create address: %v", err)
	}

	// 接收方地址，需替换为实际接收地址
	toAddressStr := "1ExampleAddress1234567890"
	toAddress, err := btcutil.DecodeAddress(toAddressStr, &chaincfg.MainNetParams)
	if err != nil {
		log.Fatalf("Failed to decode address: %v", err)
	}

	// 获取未花费交易输出（UTXO）
	utxos, err := client.ListUnspentMinMaxAddresses(1, 9999999, []btcutil.Address{addressPubKeyHash})
	if err != nil {
		log.Fatalf("Failed to list unspent outputs: %v", err)
	}
	if len(utxos) == 0 {
		log.Fatal("No unspent outputs available")
	}

	// 创建新交易
	tx := wire.NewMsgTx(wire.TxVersion)

	// 添加输入
	for _, utxo := range utxos {
		outPoint := wire.NewOutPoint(&utxo.TxID, utxo.Vout)
		txIn := wire.NewTxIn(outPoint, nil, nil)
		tx.AddTxIn(txIn)
	}

	// 转账金额（以聪为单位）
	amountToSend := btcutil.Amount(10000)

	// 添加输出
	pkScript, err := txscript.PayToAddrScript(toAddress)
	if err != nil {
		log.Fatalf("Failed to create pay-to-address script: %v", err)
	}
	txOut := wire.NewTxOut(int64(amountToSend), pkScript)
	tx.AddTxOut(txOut)

	// 计算手续费，这里简单设置，实际应用需根据网络情况调整
	fee := btcutil.Amount(1000)
	changeAmount := btcutil.Amount(0)
	for _, utxo := range utxos {
		changeAmount += utxo.Amount
	}
	changeAmount -= amountToSend + fee

	// 添加找零输出
	if changeAmount > 0 {
		changePkScript, err := txscript.PayToAddrScript(addressPubKeyHash)
		if err != nil {
			log.Fatalf("Failed to create change script: %v", err)
		}
		changeTxOut := wire.NewTxOut(int64(changeAmount), changePkScript)
		tx.AddTxOut(changeTxOut)
	}

	// 签名交易
	for i, txIn := range tx.TxIn {
		utxo := utxos[i]
		pkScript, err := client.GetTxOutScriptPubKey(&utxo.TxID, utxo.Vout)
		if err != nil {
			log.Fatalf("Failed to get script pub key: %v", err)
		}
		sigScript, err := txscript.SignatureScript(tx, i, pkScript, txscript.SigHashAll, privateKey, true)
		if err != nil {
			log.Fatalf("Failed to sign transaction: %v", err)
		}
		txIn.SignatureScript = sigScript
	}

	// 发送交易
	txHash, err := client.SendRawTransaction(tx, false)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	fmt.Printf("Transaction sent. TxID: %s\n", txHash.String())
}
