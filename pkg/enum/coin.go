package enum

// 定义支持的币种类型
type CoinType uint32

const (
	BTC CoinType = 0  // BIP44 比特币路径 m/44'/0'/0'/0/...
	ETH CoinType = 60 // BIP44 以太坊路径 m/44'/60'/0'/0/...
)
