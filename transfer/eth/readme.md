### 创建测试环境
1. 创建目录
```shell
mdkir -p ~/ethereum-multinode/{node1,node2,node3}
```
2. 创建账户
```shell
geth --datadir /Volumes/Lenovo/bit/ethereum-multinode/node1 account new
0xb9041E3A912A6dB7bD1f1AC959003f6Dde062C34
0x750977Cf178a6627E7584Ab27C8040Df27ADcd9f
0x16c49231F471f6446C0883e1907566FfD495af21
```
3. 初始化节点
    1. 创建创世文件 genesis.json
```json
{
  "config": {
    "chainId": 1337,
    "terminalTotalDifficulty": 0,
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
  "extradata": "0x0000000000000000000000000000000000000000000000000000000000000000b9041E3A912A6dB7bD1f1AC959003f6Dde062C340000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
  "difficulty": "1",
  "gasLimit": "8000000",
  "alloc": {
    "b9041E3A912A6dB7bD1f1AC959003f6Dde062C34": {
      "balance": "1000000000000000000000"
    }
  }
}
```
    2. 初始化节点
```shell
geth --datadir /Volumes/Lenovo/bit/ethereum-multinode/node1 init ~/ethereum-multinode/genesis.json
```
  3. dd
```shell
geth --datadir node1 --port 30303 console
> admin.nodeInfo.enode
# 输出示例：enode://a1b2...@127.0.0.1:30303
# 在每个节点的数据目录中创建 static-nodes.json：
# 节点1的 static-nodes.json（包含其他节点enode）
echo '[
  "enode://<节点2的enode>@127.0.0.1:30304",
  "enode://<节点3的enode>@127.0.0.1:30305"
]' > node1/geth/static-nodes.json

./build/bin/geth --datadir ~/bin/ethereum-multinode/node1 \
     --networkid 1337 \         # 与 genesis.json 的 chainId 一致
     --port 30303 \              # 节点间通信端口
     --http \                    # 启用 HTTP-RPC
     --http.addr 0.0.0.0 \
     --http.port 8545 \
     --http.api "eth,net,web3,personal,miner" \
     --unlock "b9041e3a912a6db7bd1f1ac959003f6dde062c34" \    # 解锁账户以挖矿
     --password < "1" \
     --allow-insecure-unlock \
     --mine \                    # 启动挖矿
     --miner.threads 1
```
4. 启动节点1
```shell

     ./build/bin/geth --datadir ~/bit/ethereum-multinode/node1 \
     --networkid 1337 \
     --port 30303 \
     --http \
     --http.addr 0.0.0.0 \
     --http.port 8545 \
     --http.api "eth,net,web3,personal,miner,engine" \
     --http.corsdomain "*" \
     --http.vhosts "*" \
     --syncmode full \
     --unlock "b9041e3a912a6db7bd1f1ac959003f6dde062c34" \
     --password <(echo -n "1") \
     --allow-insecure-unlock \
     --authrpc.port 8551 \
     --authrpc.vhosts localhost \
     --authrpc.jwtsecret ~/bit/ethereum-multinode/jwt/jwtsecret.hex \
     --mine

     PRYSM_ALLOW_UNVERIFIED_BINARIES=1 ./prysm.sh beacon-chain \
    --datadir=/Users/lhs/bit/ethereum-multinode/prysm \
    --jwt-secret=/Users/lhs/bit/ethereum-multinode/jwt/jwtsecret.hex \
    --execution-endpoint=http://localhost:8551 \
    --chain-id=1337
```
5. 启动另一个节点2
```shell
./build/bin/geth --datadir ~/bit/ethereum-multinode/node2 \
     --networkid 1337 \
     --port 30304 \
     --http \
     --http.addr 0.0.0.0 \
     --http.port 8546 \
     --http.api "eth,net,web3,personal,miner,engine" \
     --http.corsdomain "*" \
     --http.vhosts "*" \
     --bootnodes "enode://449988454dec135fcecda7433045beddb646629a8e4b212e85d1992212af173fba2b808448815c0dc016b53e09e574d2651d7f8cb53d33fa9836388c0c8567d0@127.0.0.1:30303?discport=21256" \
     --unlock "750977cf178a6627e7584ab27c8040df27adcd9f" \
     --password <(echo -n "1") \
     --syncmode full \
     --allow-insecure-unlock \
     --authrpc.port 8552 \
     --authrpc.vhosts localhost \
     --authrpc.jwtsecret ~/bit/ethereum-multinode/jwt/jwtsecret.hex \
     --mine
```
6. 启动第三个节点3 跟随节点
```shell
./build/bin/geth --datadir ~/bit/ethereum-multinode/node3 \
     --networkid 1337 \
     --port 30305 \
     --authrpc.port 8553 \
     --authrpc.vhosts localhost \
     --authrpc.jwtsecret ~/bit/ethereum-multinode/jwt/jwtsecret.hex \
     --syncmode full \
     --bootnodes "enode://449988454dec135fcecda7433045beddb646629a8e4b212e85d1992212af173fba2b808448815c0dc016b53e09e574d2651d7f8cb53d33fa9836388c0c8567d0@127.0.0.1:30303?discport=21256" \
     --http \
     --http.addr 0.0.0.0 \
     --http.api "eth,net,web3,personal,miner,engine" \
     --http.corsdomain "*" \
     --http.vhosts "*" \
     --http.port 8547
```
