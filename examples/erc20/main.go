package main

import (
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	etherkit "github.com/guanzhenxing/go-ether-kit"
	"github.com/guanzhenxing/go-ether-kit/contracts/erc20"
)

func main() {
	// 配置信息 - 请使用你自己的私钥和RPC URL
	privateKey := "your_private_key_here"                        // 请替换为真实的私钥
	rpcURL := "https://eth-goerli.g.alchemy.com/v2/your-api-key" // 请替换为真实的RPC URL

	// USDC 合约地址 (Goerli 测试网)
	usdcAddress := common.HexToAddress("0x07865c6E87B9F70255377e024ace6630C1Eaa37F")

	fmt.Println("🪙 go-ether-kit ERC20 代币操作示例")
	fmt.Println(strings.Repeat("=", 60))

	// 1. 创建钱包连接
	fmt.Println("\n1. 创建钱包连接...")
	wallet, err := etherkit.NewWallet(privateKey, rpcURL)
	if err != nil {
		log.Fatalf("创建钱包失败: %v", err)
	}
	defer wallet.CloseWallet()

	address := wallet.GetAddress()
	fmt.Printf("✅ 钱包地址: %s\n", address.Hex())

	// 2. 创建 ERC20 合约实例
	fmt.Println("\n2. 连接 USDC 合约...")
	token, err := erc20.NewIERC20(usdcAddress, wallet.GetClient())
	if err != nil {
		log.Fatalf("创建ERC20合约实例失败: %v", err)
	}

	// 3. 查询代币基本信息
	fmt.Println("\n3. 查询代币基本信息...")

	// 代币名称
	name, err := token.Name(nil)
	if err != nil {
		log.Printf("❌ 获取代币名称失败: %v", err)
	} else {
		fmt.Printf("✅ 代币名称: %s\n", name)
	}

	// 代币符号
	symbol, err := token.Symbol(nil)
	if err != nil {
		log.Printf("❌ 获取代币符号失败: %v", err)
	} else {
		fmt.Printf("✅ 代币符号: %s\n", symbol)
	}

	// 小数位数
	decimals, err := token.Decimals(nil)
	if err != nil {
		log.Printf("❌ 获取小数位数失败: %v", err)
	} else {
		fmt.Printf("✅ 小数位数: %d\n", decimals)
	}

	// 总供应量
	totalSupply, err := token.TotalSupply(nil)
	if err != nil {
		log.Printf("❌ 获取总供应量失败: %v", err)
	} else {
		totalSupplyDecimal := etherkit.ToDecimal(totalSupply, int(decimals))
		fmt.Printf("✅ 总供应量: %s %s\n", totalSupplyDecimal.String(), symbol)
	}

	// 4. 查询账户代币余额
	fmt.Println("\n4. 查询账户代币余额...")
	balance, err := token.BalanceOf(nil, address)
	if err != nil {
		log.Printf("❌ 获取代币余额失败: %v", err)
	} else {
		balanceDecimal := etherkit.ToDecimal(balance, int(decimals))
		fmt.Printf("✅ %s 余额: %s %s\n", symbol, balanceDecimal.String(), symbol)
	}

	// 5. 查询授权额度
	fmt.Println("\n5. 查询授权额度...")
	spenderAddress := common.HexToAddress("0x742F35C6dB4634C0532925a3b8D6dA2E") // 示例授权地址
	allowance, err := token.Allowance(nil, address, spenderAddress)
	if err != nil {
		log.Printf("❌ 获取授权额度失败: %v", err)
	} else {
		allowanceDecimal := etherkit.ToDecimal(allowance, int(decimals))
		fmt.Printf("✅ 对 %s 的授权额度: %s %s\n",
			spenderAddress.Hex(), allowanceDecimal.String(), symbol)
	}

	// 6. 代币转账示例（构建交易但不发送）
	fmt.Println("\n6. 构建代币转账交易...")
	toAddress := common.HexToAddress("0x742F35C6dB4634C0532925a3b8D6dA2E")
	transferAmount := etherkit.ToWei("10", int(decimals)) // 转账 10 USDC

	// 构建交易选项
	txOpts, err := wallet.BuildTxOpts(
		big.NewInt(0), // value (ERC20转账value为0)
		nil,           // nonce (自动计算)
		nil,           // gasPrice (自动获取)
	)
	if err != nil {
		log.Printf("❌ 构建交易选项失败: %v", err)
	} else {
		fmt.Printf("✅ 交易选项构建成功\n")
		fmt.Printf("   - From: %s\n", address.Hex())
		fmt.Printf("   - To: %s\n", toAddress.Hex())
		fmt.Printf("   - Amount: %s %s\n",
			etherkit.ToDecimal(transferAmount, int(decimals)).String(), symbol)
		fmt.Printf("   - Gas Price: %s Gwei\n",
			etherkit.ToDecimal(txOpts.GasPrice, 9).String())
		fmt.Printf("   - Nonce: %s\n", txOpts.Nonce.String())

		// 注意：这里只是演示，不实际发送交易
		fmt.Println("   ⚠️  交易已构建但未发送（仅演示）")
	}

	// 7. 授权操作示例（构建交易但不发送）
	fmt.Println("\n7. 构建代币授权交易...")
	approveAmount := etherkit.ToWei("100", int(decimals)) // 授权 100 USDC

	approveTxOpts, err := wallet.BuildTxOpts(
		big.NewInt(0), // value
		nil,           // nonce
		nil,           // gasPrice
	)
	if err != nil {
		log.Printf("❌ 构建授权交易选项失败: %v", err)
	} else {
		fmt.Printf("✅ 授权交易构建成功\n")
		fmt.Printf("   - Spender: %s\n", spenderAddress.Hex())
		fmt.Printf("   - Amount: %s %s\n",
			etherkit.ToDecimal(approveAmount, int(decimals)).String(), symbol)
		fmt.Printf("   - Gas Price: %s Gwei\n",
			etherkit.ToDecimal(approveTxOpts.GasPrice, 9).String())
		fmt.Printf("   ⚠️  交易已构建但未发送（仅演示）")
	}

	// 8. 使用合约调用接口查询余额
	fmt.Println("\n8. 使用合约调用接口查询余额...")

	// 获取 ERC20 ABI
	abiString := `[{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`
	contractAbi, err := etherkit.GetABI(abiString)
	if err != nil {
		log.Printf("❌ 解析ABI失败: %v", err)
	} else {
		result, err := wallet.CallContract(usdcAddress, contractAbi, "balanceOf", address)
		if err != nil {
			log.Printf("❌ 合约调用失败: %v", err)
		} else {
			if len(result) > 0 {
				if balance, ok := result[0].(*big.Int); ok {
					balanceDecimal := etherkit.ToDecimal(balance, int(decimals))
					fmt.Printf("✅ 通过合约调用查询余额: %s %s\n",
						balanceDecimal.String(), symbol)
				}
			}
		}
	}

	// 9. 计算交易费用
	fmt.Println("\n9. 估算交易费用...")

	// 估算转账交易的 Gas
	transferData, err := etherkit.BuildContractInputData(
		contractAbi, "transfer", toAddress, transferAmount)
	if err != nil {
		log.Printf("❌ 构建交易数据失败: %v", err)
	} else {
		gasLimit, err := wallet.GetEthProvider().EstimateGas(
			address, usdcAddress, 0, nil, big.NewInt(0), transferData)
		if err != nil {
			log.Printf("❌ 估算Gas失败: %v", err)
		} else {
			gasPrice, _ := wallet.GetEthProvider().GetSuggestGasPrice()
			totalFee := new(big.Int).Mul(big.NewInt(int64(gasLimit)), gasPrice)
			totalFeeEth := etherkit.ToDecimal(totalFee, etherkit.EthDecimals)

			fmt.Printf("✅ 转账交易费用估算:\n")
			fmt.Printf("   - Gas Limit: %d\n", gasLimit)
			fmt.Printf("   - Gas Price: %s Gwei\n",
				etherkit.ToDecimal(gasPrice, 9).String())
			fmt.Printf("   - 总费用: %s ETH\n", totalFeeEth.String())
		}
	}

	// 10. 代币相关常量
	fmt.Println("\n10. 常用代币操作常量...")
	fmt.Printf("✅ ERC20 Transfer 方法ID: %s\n", etherkit.ERC20TransferMethodID)
	fmt.Printf("✅ ERC20 Approve 方法ID: %s\n", etherkit.ERC20ApproveMethodID)
	fmt.Printf("✅ ERC20 Transfer 事件主题: %s\n", etherkit.ERC20TransferEventTopic)
	fmt.Printf("✅ ERC20 Approval 事件主题: %s\n", etherkit.ERC20ApprovalEventTopic)

	fmt.Println("\n🎉 ERC20 示例完成！")
	fmt.Println("\n提示：")
	fmt.Println("- 要执行真实交易，请取消注释相关代码并设置正确的参数")
	fmt.Println("- 建议先在测试网上进行测试")
	fmt.Println("- 转账前请确保有足够的代币余额和ETH手续费")
	fmt.Println("- 授权操作请谨慎，避免授权过大金额")
}
