package etherkit

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/sha3"
)

//############ Address ############

const (
	ZeroAddress = "0x0000000000000000000000000000000000000000"
	EAddress    = "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
)

const ZeroHash = "0x0000000000000000000000000000000000000000000000000000000000000000"

// IsValidAddress 验证是否是十六进制地址
func IsValidAddress(iAddress interface{}) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	switch v := iAddress.(type) {
	case string:
		return re.MatchString(v)
	case common.Address:
		return re.MatchString(v.Hex())
	default:
		return false
	}
}

//############ Account ############

// GeneratePrivateKey 创建私钥
func GeneratePrivateKey() (*ecdsa.PrivateKey, error) {
	return crypto.GenerateKey()
}

// GetHexPrivateKey 创建十六进制的私钥(不以0x开头)
func GetHexPrivateKey(privateKey *ecdsa.PrivateKey) string {
	return hexutil.Encode(crypto.FromECDSA(privateKey))[2:]
}

// PrivateKeyToAddress 从私钥中获得地址
func PrivateKeyToAddress(privateKey *ecdsa.PrivateKey) common.Address {
	publicKey := privateKey.Public()
	publicKeyECDSA := publicKey.(*ecdsa.PublicKey)
	return crypto.PubkeyToAddress(*publicKeyECDSA)
}

// GetHexPublicKey 从私钥中获得十六进制的公钥(不以0x和04开头)
func GetHexPublicKey(privateKey *ecdsa.PrivateKey) string {
	publicKey := privateKey.Public()
	publicKeyECDSA := publicKey.(*ecdsa.PublicKey)

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	return hexutil.Encode(publicKeyBytes)[4:]
}

// BuildPrivateKeyFromHex 从字符串私钥构建私钥对象
func BuildPrivateKeyFromHex(privateKeyHex string) (*ecdsa.PrivateKey, error) {
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func BuildPrivateKeyFromMnemonic(mnemonic string) (*ecdsa.PrivateKey, error) {
	return BuildPrivateKeyFromMnemonicAndAccountId(mnemonic, 0)
}

// BuildPrivateKeyFromMnemonicAndAccountId 从助记词获得私钥
func BuildPrivateKeyFromMnemonicAndAccountId(mnemonic string, accountId uint32) (*ecdsa.PrivateKey, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create HD wallet from mnemonic")
	}
	path, err := accounts.ParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", accountId))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse derivation path")
	}
	account, err := wallet.Derive(path, true)
	if err != nil {
		return nil, errors.Wrap(err, "failed to derive account from HD wallet")
	}
	pk, err := wallet.PrivateKey(account)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account's private key from HD wallet")
	}
	return pk, nil
}

// PublicKeyBytesToAddress 公钥到地址
func PublicKeyBytesToAddress(publicKey []byte) common.Address {
	var buf []byte

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKey[1:]) // remove EC prefix 04
	buf = hash.Sum(nil)
	address := buf[12:]

	return common.HexToAddress(hex.EncodeToString(address))
}

//############ Cast ############

// ToDecimal wei to decimals
func ToDecimal(iValue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)
	switch v := iValue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}

// ToWei decimals to wei
func ToWei(iAmount interface{}, decimals int) *big.Int {
	amount := decimal.NewFromFloat(0)
	switch v := iAmount.(type) {
	case string:
		amount, _ = decimal.NewFromString(v)
	case float64:
		amount = decimal.NewFromFloat(v)
	case int64:
		amount = decimal.NewFromFloat(float64(v))
	case int:
		amount = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		amount = v
	case *decimal.Decimal:
		amount = *v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	result := amount.Mul(mul)

	wei := new(big.Int)
	wei.SetString(result.String(), 10)

	return wei
}

//############ Contract ############

// GetABI 从abi字符串中获得ABI对象
func GetABI(abiStr string) (abi.ABI, error) {
	abiContract, err := abi.JSON(strings.NewReader(abiStr))
	return abiContract, err
}

// GetContractMethodId 获得合约的methodId
// 参数method，如：transfer(address,uint256)
func GetContractMethodId(method string) string {
	methodId := hexutil.Encode(crypto.Keccak256([]byte(method))[:4])
	return methodId
}

// GetEventTopic 获得事件的topic。event字符串如：transfer(address,uint256)
func GetEventTopic(event string) string {
	return crypto.Keccak256Hash([]byte(event)).String()
}

// BuildContractInputData 构建合约的input data
func BuildContractInputData(contract abi.ABI, name string, args ...interface{}) ([]byte, error) {
	return contract.Pack(name, args...)
}

//############ Transaction ############

// NewTx 新建一个tx
func NewTx(to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, data []byte) (*types.Transaction, error) {
	return types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	}), nil
}

// NewTxWithHexData 基于hexData构建一个Tx
func NewTxWithHexData(to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, hexData string) (*types.Transaction, error) {
	data, err := hex.DecodeString(hexData)
	if err != nil {
		return nil, err
	}
	return types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	}), nil
}

// VerifySignature 验证签名
// address: 用于签名的地址
// digestHash: 用于验证的原始数据
// signature: 需要进行验证的签名数据
func VerifySignature(address string, data, signature []byte) bool {

	digestHash := crypto.Keccak256Hash(data)
	//returns the public key that created the given signature.
	sigPublicKeyECDSA, err := crypto.SigToPub(digestHash.Bytes(), signature)
	if err != nil {
		return false
	}

	sigAddress := crypto.PubkeyToAddress(*sigPublicKeyECDSA)
	return sigAddress.String() == address
}

// DecodeRawTxHex 解析rawTx
func DecodeRawTxHex(rawTx string) (*types.Transaction, error) {

	tx := new(types.Transaction)
	rawTxBytes, err := hex.DecodeString(rawTx)
	if err != nil {
		return nil, err
	}
	err = rlp.DecodeBytes(rawTxBytes, &tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// GetMaxUint256 获得合约中MaxUint256
func GetMaxUint256() *big.Int {
	return new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), big.NewInt(1))
}
