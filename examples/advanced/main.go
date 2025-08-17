package main

import (
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	etherkit "github.com/guanzhenxing/go-evm-kit"
)

func main() {
	// 配置信息
	privateKey := "your_private_key_here"
	rpcURL := "https://eth-goerli.g.alchemy.com/v2/your-api-key"

	fmt.Println("🚀 go-ether-kit 高级功能示例")
	fmt.Println(strings.Repeat("=", 60))

	// 1. 创建多个网络连接
	fmt.Println("\n1. 多网络连接示例...")

	// 主网提供者
	mainnetProvider, err := etherkit.NewProviderWithChainId(
		"https://eth-mainnet.g.alchemy.com/v2/your-api-key",
		etherkit.MainnetChainID)
	if err != nil {
		log.Printf("❌ 创建主网连接失败: %v", err)
	} else {
		defer mainnetProvider.Close()
		fmt.Printf("✅ 主网连接成功 (Chain ID: %d)\n", etherkit.MainnetChainID)
	}

	// 测试网提供者
	goerliProvider, err := etherkit.NewProviderWithChainId(rpcURL, etherkit.GoerliChainID)
	if err != nil {
		log.Fatalf("创建测试网连接失败: %v", err)
	}
	defer goerliProvider.Close()
	fmt.Printf("✅ Goerli测试网连接成功 (Chain ID: %d)\n", etherkit.GoerliChainID)

	// 2. 创建钱包并演示高级功能
	fmt.Println("\n2. 创建钱包...")
	signer, err := etherkit.NewSignerFromHexPrivateKey(privateKey)
	if err != nil {
		log.Fatalf("创建签名器失败: %v", err)
	}

	wallet, err := etherkit.NewWalletWithComponents(signer, goerliProvider)
	if err != nil {
		log.Fatalf("创建钱包失败: %v", err)
	}
	defer wallet.CloseWallet()

	address := wallet.GetAddress()
	fmt.Printf("✅ 钱包地址: %s\n", address.Hex())

	// 3. 网络状态监控
	fmt.Println("\n3. 网络状态监控...")

	// 获取多个网络配置信息
	networks := map[string]int64{
		"Ethereum Mainnet": etherkit.MainnetChainID,
		"Polygon":          etherkit.PolygonChainID,
		"BSC":              etherkit.BSCChainID,
		"Arbitrum":         etherkit.ArbitrumChainID,
	}

	for name, chainID := range networks {
		if config, exists := etherkit.NetworkConfigs[chainID]; exists {
			fmt.Printf("✅ %s:\n", name)
			fmt.Printf("   - Chain ID: %d\n", config.ChainID)
			fmt.Printf("   - Symbol: %s\n", config.Symbol)
			fmt.Printf("   - Block Time: %ds\n", config.BlockTime)
			fmt.Printf("   - Confirmations: %d\n", config.Confirmations)
		}
	}

	// 4. 高级交易构建
	fmt.Println("\n4. 高级交易构建...")

	// 构建带有自定义参数的交易
	toAddress := common.HexToAddress("0x742F35C6dB4634C0532925a3b8D6dA2E")
	value := etherkit.ToWei("0.001", etherkit.EthDecimals)

	// 获取当前nonce
	nonce, err := wallet.GetNonce()
	if err != nil {
		log.Printf("❌ 获取nonce失败: %v", err)
	} else {
		fmt.Printf("✅ 当前nonce: %d\n", nonce)
	}

	// 使用自定义Gas参数
	gasPrice := big.NewInt(20 * etherkit.GWei) // 20 Gwei
	gasLimit := uint64(21000)

	tx, err := wallet.NewTx(toAddress, nonce, gasLimit, gasPrice, value, nil)
	if err != nil {
		log.Printf("❌ 构建交易失败: %v", err)
	} else {
		fmt.Printf("✅ 高级交易构建成功:\n")
		fmt.Printf("   - Nonce: %d\n", tx.Nonce())
		fmt.Printf("   - Gas Limit: %d\n", tx.Gas())
		fmt.Printf("   - Gas Price: %s Gwei\n",
			etherkit.ToDecimal(tx.GasPrice(), 9).String())
		fmt.Printf("   - Value: %s ETH\n",
			etherkit.ToDecimal(tx.Value(), etherkit.EthDecimals).String())
	}

	// 5. 交易签名和序列化
	fmt.Println("\n5. 交易签名和序列化...")

	if tx != nil {
		// 签名交易
		signedTx, err := wallet.SignTx(tx)
		if err != nil {
			log.Printf("❌ 签名交易失败: %v", err)
		} else {
			fmt.Printf("✅ 交易签名成功\n")
			fmt.Printf("   - 交易哈希: %s\n", signedTx.Hash().Hex())

			// 序列化交易（用于离线传输）
			rawTx, err := signedTx.MarshalBinary()
			if err != nil {
				log.Printf("❌ 序列化交易失败: %v", err)
			} else {
				fmt.Printf("✅ 交易序列化成功\n")
				fmt.Printf("   - 原始交易长度: %d bytes\n", len(rawTx))
				fmt.Printf("   - 原始交易(前64字符): %x...\n", rawTx[:min(32, len(rawTx))])
			}
		}
	}

	// 6. 批量地址验证
	fmt.Println("\n6. 批量地址验证...")

	testAddresses := map[string]string{
		"Valid ETH Address":    "0x742F35C6dB4634C0532925a3b8D6dA2E",
		"Valid Checksum":       "0x742f35c6dB4634C0532925a3b8D6dA2E",
		"Zero Address":         etherkit.ZeroAddress,
		"Native Token Address": etherkit.NativeTokenAddress,
		"Invalid Format":       "invalid_address",
		"Wrong Length":         "0x742F35C6dB4634C0532925a3b8D6dA2E123",
	}

	for desc, addr := range testAddresses {
		isValid := etherkit.IsValidAddress(addr)
		status := "❌"
		if isValid {
			status = "✅"
		}
		fmt.Printf("%s %s: %s\n", status, desc, addr)
	}

	// 7. 合约方法ID和事件主题计算
	fmt.Println("\n7. 合约方法ID和事件主题...")

	methods := map[string]string{
		"transfer(address,uint256)":             "ERC20 Transfer",
		"approve(address,uint256)":              "ERC20 Approve",
		"transferFrom(address,address,uint256)": "ERC20 TransferFrom",
		"balanceOf(address)":                    "ERC20 BalanceOf",
		"mint(address,uint256)":                 "Custom Mint",
	}

	for method, desc := range methods {
		methodID := etherkit.GetContractMethodId(method)
		fmt.Printf("✅ %s:\n", desc)
		fmt.Printf("   - 方法签名: %s\n", method)
		fmt.Printf("   - 方法ID: %s\n", methodID)
	}

	events := map[string]string{
		"Transfer(address,address,uint256)": "ERC20 Transfer Event",
		"Approval(address,address,uint256)": "ERC20 Approval Event",
		"Mint(address,uint256)":             "Custom Mint Event",
	}

	for event, desc := range events {
		topic := etherkit.GetEventTopic(event)
		fmt.Printf("✅ %s:\n", desc)
		fmt.Printf("   - 事件签名: %s\n", event)
		fmt.Printf("   - 事件主题: %s\n", topic)
	}

	// 8. 单位转换和计算
	fmt.Println("\n8. 单位转换和数学计算...")

	// 各种单位转换
	amounts := []string{"0.001", "1", "100", "1000.5"}
	for _, amount := range amounts {
		wei := etherkit.ToWei(amount, etherkit.EthDecimals)
		gwei := etherkit.ToWei(amount, 9)

		fmt.Printf("✅ %s ETH:\n", amount)
		fmt.Printf("   - Wei: %s\n", wei.String())
		fmt.Printf("   - Gwei: %s\n", gwei.String())

		// 转回ETH验证
		ethBack := etherkit.ToDecimal(wei, etherkit.EthDecimals)
		fmt.Printf("   - 转回ETH: %s\n", ethBack.String())
	}

	// 9. 错误处理示例
	fmt.Println("\n9. 错误处理示例...")

	// 演示各种错误情况
	fmt.Println("✅ 测试各种错误情况:")

	// 无效地址
	if !etherkit.IsValidAddress("invalid") {
		fmt.Printf("   - 检测到无效地址: %v\n", etherkit.ErrInvalidAddress)
	}

	// 无效私钥格式
	_, err = etherkit.BuildPrivateKeyFromHex("invalid_key")
	if err != nil {
		fmt.Printf("   - 检测到无效私钥: %v\n", err)
	}

	// 10. 性能测试
	fmt.Println("\n10. 性能测试...")

	// 测试地址验证性能
	start := time.Now()
	validCount := 0
	testCount := 1000

	for i := 0; i < testCount; i++ {
		testAddr := fmt.Sprintf("0x742F35C6dB4634C0532925a3b8D6dA2E%04d", i%10000)
		if etherkit.IsValidAddress(testAddr) {
			validCount++
		}
	}

	duration := time.Since(start)
	fmt.Printf("✅ 地址验证性能测试:\n")
	fmt.Printf("   - 测试数量: %d 个地址\n", testCount)
	fmt.Printf("   - 耗时: %v\n", duration)
	fmt.Printf("   - 平均每次: %v\n", duration/time.Duration(testCount))
	fmt.Printf("   - 有效地址: %d\n", validCount)

	// 测试单位转换性能
	start = time.Now()
	for i := 0; i < testCount; i++ {
		amount := fmt.Sprintf("%d.%d", i%1000, i%100)
		wei := etherkit.ToWei(amount, etherkit.EthDecimals)
		etherkit.ToDecimal(wei, etherkit.EthDecimals)
	}
	duration = time.Since(start)
	fmt.Printf("✅ 单位转换性能测试:\n")
	fmt.Printf("   - 测试数量: %d 次转换\n", testCount*2)
	fmt.Printf("   - 耗时: %v\n", duration)
	fmt.Printf("   - 平均每次: %v\n", duration/time.Duration(testCount*2))

	fmt.Println("\n🎉 高级功能示例完成！")
	fmt.Println("\n高级功能总结：")
	fmt.Println("- ✅ 多网络连接管理")
	fmt.Println("- ✅ 高级交易构建和签名")
	fmt.Println("- ✅ 批量地址验证")
	fmt.Println("- ✅ 合约方法ID和事件主题计算")
	fmt.Println("- ✅ 灵活的单位转换")
	fmt.Println("- ✅ 完善的错误处理")
	fmt.Println("- ✅ 性能监控和测试")
}

// 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
