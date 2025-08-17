package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	etherkit "github.com/guanzhenxing/go-evm-kit"
)

func main() {
	// 配置信息 - 请使用你自己的私钥和RPC URL
	privateKey := "your_private_key_here"                        // 请替换为真实的私钥
	rpcURL := "https://eth-goerli.g.alchemy.com/v2/your-api-key" // 请替换为真实的RPC URL

	fmt.Println("🚀 go-ether-kit 基础示例")
	fmt.Println(strings.Repeat("=", 50))

	// 1. 创建钱包连接
	fmt.Println("\n1. 创建钱包连接...")
	wallet, err := etherkit.NewWallet(privateKey, rpcURL)
	if err != nil {
		log.Fatalf("创建钱包失败: %v", err)
	}
	defer wallet.CloseWallet()

	// 获取钱包地址
	address := wallet.GetAddress()
	fmt.Printf("✅ 钱包地址: %s\n", address.Hex())

	// 2. 获取网络信息
	fmt.Println("\n2. 获取网络信息...")
	provider := wallet.GetEthProvider()

	chainID, err := provider.GetChainID()
	if err != nil {
		log.Printf("❌ 获取Chain ID失败: %v", err)
	} else {
		fmt.Printf("✅ Chain ID: %s\n", chainID.String())
	}

	blockNumber, err := provider.GetBlockNumber()
	if err != nil {
		log.Printf("❌ 获取最新区块号失败: %v", err)
	} else {
		fmt.Printf("✅ 最新区块号: %d\n", blockNumber)
	}

	gasPrice, err := provider.GetSuggestGasPrice()
	if err != nil {
		log.Printf("❌ 获取Gas价格失败: %v", err)
	} else {
		gasPriceGwei := etherkit.ToDecimal(gasPrice, 9) // 转换为Gwei
		fmt.Printf("✅ 当前Gas价格: %s Gwei\n", gasPriceGwei.String())
	}

	// 3. 查询账户余额
	fmt.Println("\n3. 查询账户余额...")
	balance, err := wallet.GetBalance()
	if err != nil {
		log.Printf("❌ 获取余额失败: %v", err)
	} else {
		ethBalance := etherkit.ToDecimal(balance, etherkit.EthDecimals)
		fmt.Printf("✅ ETH 余额: %s ETH\n", ethBalance.String())
	}

	// 4. 地址验证示例
	fmt.Println("\n4. 地址验证示例...")
	testAddresses := []string{
		"0x742F35C6dB4634C0532925a3b8D6dA2E",
		"0x742f35c6db4634c0532925a3b8d6da2e", // 小写
		"invalid-address",
		"0x742F35C6dB4634C0532925a3b8D6dA2E123", // 长度错误
	}

	for _, addr := range testAddresses {
		isValid := etherkit.IsValidAddress(addr)
		if isValid {
			fmt.Printf("✅ 有效地址: %s\n", addr)
		} else {
			fmt.Printf("❌ 无效地址: %s\n", addr)
		}
	}

	// 5. 单位转换示例
	fmt.Println("\n5. 单位转换示例...")

	// ETH 转 Wei
	ethAmount := "1.5"
	weiAmount := etherkit.ToWei(ethAmount, etherkit.EthDecimals)
	fmt.Printf("✅ %s ETH = %s Wei\n", ethAmount, weiAmount.String())

	// Wei 转 ETH
	ethBack := etherkit.ToDecimal(weiAmount, etherkit.EthDecimals)
	fmt.Printf("✅ %s Wei = %s ETH\n", weiAmount.String(), ethBack.String())

	// Gwei 转换
	gweiAmount := etherkit.ToWei("20", 9) // 20 Gwei
	fmt.Printf("✅ 20 Gwei = %s Wei\n", gweiAmount.String())

	// 6. 交易构建示例（不发送）
	fmt.Println("\n6. 构建交易示例...")
	toAddress := common.HexToAddress("0x742F35C6dB4634C0532925a3b8D6dA2E")
	value := etherkit.ToWei("0.01", etherkit.EthDecimals) // 0.01 ETH

	tx, err := wallet.NewTx(
		toAddress, // 收款地址
		0,         // nonce (0 = 自动计算)
		0,         // gasLimit (0 = 自动估算)
		nil,       // gasPrice (nil = 自动获取)
		value,     // 转账金额
		nil,       // 交易数据
	)
	if err != nil {
		log.Printf("❌ 构建交易失败: %v", err)
	} else {
		fmt.Printf("✅ 交易构建成功\n")
		fmt.Printf("   - To: %s\n", tx.To().Hex())
		fmt.Printf("   - Value: %s ETH\n", etherkit.ToDecimal(tx.Value(), etherkit.EthDecimals).String())
		fmt.Printf("   - Gas Limit: %d\n", tx.Gas())
		fmt.Printf("   - Gas Price: %s Gwei\n", etherkit.ToDecimal(tx.GasPrice(), 9).String())
		fmt.Printf("   - Nonce: %d\n", tx.Nonce())
	}

	// 7. 签名示例
	fmt.Println("\n7. 数据签名示例...")
	message := "Hello, Ethereum!"
	signature, err := wallet.Signature([]byte(message))
	if err != nil {
		log.Printf("❌ 签名失败: %v", err)
	} else {
		fmt.Printf("✅ 消息签名成功\n")
		fmt.Printf("   - 原始消息: %s\n", message)
		fmt.Printf("   - 签名长度: %d bytes\n", len(signature))

		// 验证签名
		isValid := etherkit.VerifySignature(address.Hex(), []byte(message), signature)
		if isValid {
			fmt.Printf("✅ 签名验证成功\n")
		} else {
			fmt.Printf("❌ 签名验证失败\n")
		}
	}

	fmt.Println("\n🎉 基础示例完成！")
	fmt.Println("\n注意：")
	fmt.Println("- 请替换示例中的私钥和RPC URL为真实值")
	fmt.Println("- 在主网上操作前，请先在测试网上测试")
	fmt.Println("- 私钥请妥善保管，不要泄露")
}
