package etherkit

import "errors"

// 标准错误定义
var (
	// 网络相关错误
	ErrNetworkConnection = errors.New("failed to connect to ethereum network")
	ErrInvalidRPCURL     = errors.New("invalid RPC URL")
	ErrNetworkTimeout    = errors.New("network request timeout")

	// 地址相关错误
	ErrInvalidAddress = errors.New("invalid ethereum address")
	ErrZeroAddress    = errors.New("address cannot be zero address")

	// 私钥相关错误
	ErrInvalidPrivateKey = errors.New("invalid private key")
	ErrInvalidMnemonic   = errors.New("invalid mnemonic phrase")
	ErrInvalidKeyFormat  = errors.New("invalid key format")

	// 交易相关错误
	ErrInsufficientFunds = errors.New("insufficient funds for transaction")
	ErrInvalidGasPrice   = errors.New("invalid gas price")
	ErrInvalidGasLimit   = errors.New("invalid gas limit")
	ErrInvalidNonce      = errors.New("invalid nonce")
	ErrTransactionFailed = errors.New("transaction execution failed")

	// 合约相关错误
	ErrContractCall           = errors.New("contract call failed")
	ErrInvalidABI             = errors.New("invalid contract ABI")
	ErrInvalidContractAddress = errors.New("invalid contract address")

	// 签名相关错误
	ErrSignatureFailed             = errors.New("signature generation failed")
	ErrInvalidSignature            = errors.New("invalid signature")
	ErrSignatureVerificationFailed = errors.New("signature verification failed")

	// 钱包相关错误
	ErrWalletClosed        = errors.New("wallet connection is closed")
	ErrInvalidWalletConfig = errors.New("invalid wallet configuration")
)
