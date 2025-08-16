package etherkit

import (
	"crypto/ecdsa"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/pkg/errors"
)

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
