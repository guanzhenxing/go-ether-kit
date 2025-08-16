package etherkit

import "math/big"

// 网络链ID常量
const (
	MainnetChainID   = 1
	GoerliChainID    = 5
	SepoliaChainID   = 11155111
	PolygonChainID   = 137
	BSCChainID       = 56
	ArbitrumChainID  = 42161
	OptimismChainID  = 10
	AvalancheChainID = 43114
	FantomChainID    = 250
)

// Gas 相关常量
const (
	DefaultGasLimit  = 21000         // ETH 转账默认 gas limit
	ContractGasLimit = 200000        // 合约调用默认 gas limit
	DefaultGasPrice  = 20000000000   // 20 Gwei
	MaxGasPrice      = 1000000000000 // 1000 Gwei
	MinGasPrice      = 1000000000    // 1 Gwei
)

// 常用 Gas Limit
const (
	ERC20TransferGasLimit = 65000
	ERC20ApproveGasLimit  = 50000
	UniswapSwapGasLimit   = 300000
)

// 时间相关常量 (秒)
const (
	DefaultTimeout        = 30
	BlockConfirmationTime = 12  // 以太坊平均出块时间
	FastConfirmationTime  = 15  // 快速确认超时
	SafeConfirmationTime  = 180 // 安全确认超时
)

// 以太坊单位常量
const (
	Wei   = 1
	GWei  = 1000000000
	Ether = 1000000000000000000
)

// 默认小数位数
const (
	EthDecimals     = 18
	USDCDecimals    = 6
	USDTDecimals    = 6
	DefaultDecimals = 18
)

// 常用地址
const (
	// 零地址
	ZeroAddress = "0x0000000000000000000000000000000000000000"
	// 原生代币地址 (用于表示 ETH/BNB/MATIC 等)
	NativeTokenAddress = "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
)

// 常用哈希
const (
	ZeroHash = "0x0000000000000000000000000000000000000000000000000000000000000000"
)

// ERC20 常用方法选择器
const (
	ERC20TransferMethodID     = "0xa9059cbb"
	ERC20TransferFromMethodID = "0x23b872dd"
	ERC20ApproveMethodID      = "0x095ea7b3"
	ERC20BalanceOfMethodID    = "0x70a08231"
	ERC20TotalSupplyMethodID  = "0x18160ddd"
)

// ERC20 常用事件主题
const (
	ERC20TransferEventTopic = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	ERC20ApprovalEventTopic = "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925"
)

// 常用 big.Int 常量
var (
	BigInt0   = big.NewInt(0)
	BigInt1   = big.NewInt(1)
	BigInt2   = big.NewInt(2)
	BigInt10  = big.NewInt(10)
	BigInt100 = big.NewInt(100)

	// Gas 价格常量
	DefaultGasPriceBig = big.NewInt(DefaultGasPrice)
	MinGasPriceBig     = big.NewInt(MinGasPrice)
	MaxGasPriceBig     = big.NewInt(MaxGasPrice)

	// Gas limit 常量
	DefaultGasLimitBig       = big.NewInt(DefaultGasLimit)
	ContractGasLimitBig      = big.NewInt(ContractGasLimit)
	ERC20TransferGasLimitBig = big.NewInt(ERC20TransferGasLimit)
	ERC20ApproveGasLimitBig  = big.NewInt(ERC20ApproveGasLimit)

	// 单位转换常量
	WeiPerEther = big.NewInt(Ether)
	WeiPerGWei  = big.NewInt(GWei)
)

// 网络配置
type NetworkConfig struct {
	ChainID       int64
	Name          string
	Symbol        string
	BlockTime     int // 秒
	Confirmations int
}

// 预定义网络配置
var NetworkConfigs = map[int64]NetworkConfig{
	MainnetChainID: {
		ChainID:       MainnetChainID,
		Name:          "Ethereum Mainnet",
		Symbol:        "ETH",
		BlockTime:     12,
		Confirmations: 12,
	},
	GoerliChainID: {
		ChainID:       GoerliChainID,
		Name:          "Goerli Testnet",
		Symbol:        "ETH",
		BlockTime:     12,
		Confirmations: 3,
	},
	SepoliaChainID: {
		ChainID:       SepoliaChainID,
		Name:          "Sepolia Testnet",
		Symbol:        "ETH",
		BlockTime:     12,
		Confirmations: 3,
	},
	PolygonChainID: {
		ChainID:       PolygonChainID,
		Name:          "Polygon",
		Symbol:        "MATIC",
		BlockTime:     2,
		Confirmations: 20,
	},
	BSCChainID: {
		ChainID:       BSCChainID,
		Name:          "Binance Smart Chain",
		Symbol:        "BNB",
		BlockTime:     3,
		Confirmations: 15,
	},
}
