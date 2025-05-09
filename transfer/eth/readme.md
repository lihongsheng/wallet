### 创建测试环境
1. 创建目录
```shell
mdkir -p ~/ethereum-multinode/{node1,node2,node3}
```
2. 创建账户
```shell
geth --datadir ~/ethereum-multinode/node1 account new
```
3. 初始化节点
    1. 创建创世文件 genesis.json
```json
{
  "config": {
    "chainId": 1337,
    "homesteadBlock": 0,
    "eip150Block": 0,
    "eip155Block": 0,
    "eip158Block": 0,
    "byzantiumBlock": 0,
    "constantinopleBlock": 0,
    "petersburgBlock": 0,
    "istanbulBlock": 0,
    "clique": {
      "period": 5,
      "epoch": 30000
    }
  },
  "extradata": "0x00000000000000000000000000000000000000000000000000000000000000005a197390748F9aB2617Ff49939680B5a984d89630000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
  "difficulty": "1",
  "gasLimit": "8000000",
  "alloc": {
    "5a197390748F9aB2617Ff49939680B5a984d8963": {
      "balance": "1000000000000000000000"
    }
  }
}
```
    2. 初始化节点
```shell
geth --datadir ~/ethereum-multinode/node1 init ~/ethereum-multinode/genesis.json
```
4. 启动节点
```shell
geth --datadir ~/ethereum-multinode/node1 --networkid 1234 --rpc --rpcaddr 0.0.0.0 --rpcport 8545 --port 30303 --nodiscover --maxpeers 0 --nat none --targetgaslimit 10000000 --allow-insecure-unlock --unlock 0 --password ~/ethereum-multinode/password.txt
```
5. 启动另一个节点
```shell
geth --datadir ~/ethereum-multinode/node2 --networkid 1234 --rpc --rpcaddr 0.0.0.0 --rpcport 8546 --port 30304 --nodiscover --maxpeers 0 --nat none --targetgaslimit 10000000 --allow-insecure-unlock --unlock 0 --password ~/ethereum-multinode/password.txt
```
6. 启动第三个节点
```shell
geth --datadir ~/ethereum-multinode/node3 --networkid 1234 --rpc --rpcaddr 0.0.0.0 --rpcport 8547 --port 30305 --nodiscover --maxpeers 0 --nat none --targetgaslimit 10000000 --allow-insecure-unlock --unlock 0 --password ~/ethereum-multinode/password.txt
```