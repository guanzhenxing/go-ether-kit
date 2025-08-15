# go-ether-kit

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**go-ether-kit** 是一个功能强大的以太坊及 EVM 兼容网络开发工具包，提供简洁易用的 API 来进行链上交互、钱包管理和智能合约操作。

## ✨ 特性

- 🔐 **钱包管理**：支持私钥、助记词、随机生成等多种方式创建账户
- 🌐 **网络连接**：轻松连接以太坊主网、测试网及其他 EVM 兼容网络  
- 💰 **交易操作**：完整的交易构建、签名、发送流程
- 📄 **智能合约**：合约调用、事件监听、ABI 处理
- 🪙 **代币支持**：内置 ERC20 代币操作支持
- 🔧 **实用工具**：单位转换、地址验证、签名验证等
- ⚡ **自动化**：自动计算 nonce、gas price 等参数
- 🔍 **链上查询**：区块、交易、余额等数据查询

## 📦 安装

```bash
go get github.com/guanzhenxing/go-ether-kit
```

## 🚀 快速开始

### 创建钱包连接

```go
package main

import (
    "fmt"
    "log"
    "github.com/guanzhenxing/go-ether-kit/etherkit"
)

func main() {
    // 使用私钥创建钱包
    privateKey := "your_private_key_here"
    rpcURL := "https://eth-mainnet.g.alchemy.com/v2/your-api-key"
    
    wallet, err := etherkit.NewWallet(privateKey, rpcURL)
    if err != nil {
        log.Fatal(err)
    }
    defer wallet.CloseWallet()
    
    // 获取账户地址
    address := wallet.GetAddress()
    fmt.Printf("钱包地址: %s\n", address.Hex())
    
    // 获取余额
    balance, err := wallet.GetBalance()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("ETH 余额: %s\n", etherkit.ToDecimal(balance, 18))
}
```

### 发送 ETH 转账

```go
func sendETH(wallet *etherkit.Wallet) {
    toAddress := common.HexToAddress("0x742F35Cc6634C0532925a3b8D6dA2e")
    amount := etherkit.ToWei("0.1", 18) // 0.1 ETH
    
    txHash, err := wallet.SendTx(
        toAddress,     // 收款地址
        0,             // nonce (0 表示自动计算)
        0,             // gasLimit (0 表示自动估算)
        nil,           // gasPrice (nil 表示自动获取)
        amount,        // 转账金额
        nil,           // 交易数据
    )
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("交易哈希: %s\n", txHash.Hex())
}
```

### ERC20 代币操作

```go
import (
    "github.com/guanzhenxing/go-ether-kit/contracts/erc20"
)

func transferToken(wallet *etherkit.Wallet) {
    tokenAddress := common.HexToAddress("0xA0b86a33E6411b6dE9C80e7F8DeD6c") // USDC 地址
    
    // 创建 ERC20 合约实例
    token, err := erc20.NewIERC20(tokenAddress, wallet.GetClient())
    if err != nil {
        log.Fatal(err)
    }
    
    // 构建交易选项
    opts, err := wallet.BuildTxOpts(
        big.NewInt(0),    // value
        nil,              // nonce (自动计算)
        nil,              // gasPrice (自动获取)
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // 转账代币
    toAddress := common.HexToAddress("0x742F35Cc6634C0532925a3b8D6dA2e")
    amount := etherkit.ToWei("100", 6) // 100 USDC (6 位小数)
    
    tx, err := token.Transfer(opts, toAddress, amount)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("代币转账交易: %s\n", tx.Hash().Hex())
}
```

### 智能合约调用

```go
func callContract(wallet *etherkit.Wallet) {
    contractAddress := common.HexToAddress("0x...")
    abiString := `[{"inputs":[],"name":"totalSupply","outputs":[{"type":"uint256"}],"type":"function"}]`
    
    // 获取合约 ABI
    contractAbi, err := etherkit.GetABI(abiString)
    if err != nil {
        log.Fatal(err)
    }
    
    // 调用合约方法 (只读)
    result, err := wallet.CallContract(contractAddress, contractAbi, "totalSupply")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("总供应量: %v\n", result[0])
}
```

## 📚 API 文档

### Provider (网络提供者)

```go
// 创建网络连接
provider, err := etherkit.NewProvider("https://eth-mainnet.g.alchemy.com/v2/your-api-key")
provider, err := etherkit.NewProviderWithChainId("https://polygon-rpc.com", 137)

// 基本查询
chainID, err := provider.GetChainID()
blockNumber, err := provider.GetBlockNumber() 
gasPrice, err := provider.GetSuggestGasPrice()
block, err := provider.GetBlockByNumber(big.NewInt(123456))
receipt, err := provider.GetTransactionReceipt(txHash)
```

### Signer (签名器)

```go
// 多种创建方式
signer, err := etherkit.NewSigner()                              // 随机生成
signer, err := etherkit.NewSignerFromHexPrivateKey("0x...")      // 私钥
signer, err := etherkit.NewSignerFromMnemonic("word1 word2...")  // 助记词

// 获取账户信息  
address := signer.GetAddress()
privateKey := signer.GetPrivateKey()
```

### Wallet (钱包)

```go
// 创建钱包
wallet, err := etherkit.NewWallet(privateKey, rpcURL)

// 账户操作
address := wallet.GetAddress()
balance, err := wallet.GetBalance()
nonce, err := wallet.GetNonce()

// 交易操作
tx, err := wallet.NewTx(toAddr, nonce, gasLimit, gasPrice, value, data)
txHash, err := wallet.SendTx(toAddr, nonce, gasLimit, gasPrice, value, data)
signedTx, err := wallet.SignTx(tx)
```

### 工具函数

```go
// 单位转换
wei := etherkit.ToWei("1.5", 18)        // 1.5 ETH 转 wei
eth := etherkit.ToDecimal(wei, 18)      // wei 转 ETH

// 地址验证
isValid := etherkit.IsValidAddress("0x...")

// 签名验证  
isValid := etherkit.VerifySignature(address, data, signature)

// 合约工具
methodID := etherkit.GetContractMethodId("transfer(address,uint256)")
eventTopic := etherkit.GetEventTopic("Transfer(address,address,uint256)")
```

## 📁 项目结构

```
go-ether-kit/
├── etherkit/           # 核心包
│   ├── provider.go     # 网络连接和查询
│   ├── signer.go       # 账户和签名
│   ├── wallet.go       # 钱包操作
│   └── utils.go        # 工具函数
├── contracts/          # 智能合约绑定
│   └── erc20/         # ERC20 合约
│       └── erc20.go
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```

## 🌐 支持的网络

- **以太坊主网** - Chain ID: 1
- **Goerli 测试网** - Chain ID: 5  
- **Sepolia 测试网** - Chain ID: 11155111
- **Polygon 主网** - Chain ID: 137
- **BSC 主网** - Chain ID: 56
- **Arbitrum One** - Chain ID: 42161
- **Optimism** - Chain ID: 10
- 以及其他 EVM 兼容网络

## 🔧 高级用法

### 批量操作

```go
// 批量查询余额
addresses := []common.Address{addr1, addr2, addr3}
for _, addr := range addresses {
    balance, _ := provider.GetEthClient().BalanceAt(context.Background(), addr, nil)
    fmt.Printf("地址 %s 余额: %s ETH\n", addr.Hex(), etherkit.ToDecimal(balance, 18))
}
```

### 事件监听

```go
// 监听 ERC20 Transfer 事件
query := ethereum.FilterQuery{
    Addresses: []common.Address{tokenAddress},
    Topics: [][]common.Hash{
        {common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")},
    },
}

logs := make(chan types.Log)
sub, err := provider.GetEthClient().SubscribeFilterLogs(context.Background(), query, logs)
if err != nil {
    log.Fatal(err)
}

for {
    select {
    case err := <-sub.Err():
        log.Fatal(err)
    case vLog := <-logs:
        fmt.Printf("发现 Transfer 事件: %s\n", vLog.TxHash.Hex())
    }
}
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目基于 MIT 许可证开源 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🔗 相关资源

- [以太坊官方文档](https://ethereum.org/developers/)
- [go-ethereum 文档](https://geth.ethereum.org/docs/)
- [Web3 开发指南](https://web3.guide/)

## 📞 支持

如有问题或建议，请：

- 提交 [Issue](https://github.com/guanzhenxing/go-ether-kit/issues)
- 发送邮件至 [your-email@example.com]
- 加入我们的讨论群组

---

⭐ 如果这个项目对你有帮助，请给个 Star！