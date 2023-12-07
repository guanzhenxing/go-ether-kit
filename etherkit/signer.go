package etherkit

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

// EthSigner ETH账户
type EthSigner interface {
	GetAddress() common.Address
	GetPrivateKey() *ecdsa.PrivateKey
}

type DefaultEthSigner struct {
	pk      *ecdsa.PrivateKey
	address common.Address
}

func newEthSignerFromPk(pk *ecdsa.PrivateKey) (*DefaultEthSigner, error) {
	return &DefaultEthSigner{
		pk:      pk,
		address: PrivateKeyToAddress(pk),
	}, nil
}

// NewEthSigner 创建一个默认的账号信息
func NewEthSigner() (*DefaultEthSigner, error) {
	pk, err := GeneratePrivateKey()
	if err != nil {
		return nil, errors.Wrap(err, "generate private key error")
	}
	return newEthSignerFromPk(pk)
}

// NewEthSignerFromMnemonic 使用助记词创建一个账号信息
func NewEthSignerFromMnemonic(mnemonic string) (*DefaultEthSigner, error) {
	return NewEthSignerFromMnemonicAndAccountId(mnemonic, 0)
}

// NewEthSignerFromMnemonicAndAccountId 使用助记词创建一个账号信息
func NewEthSignerFromMnemonicAndAccountId(mnemonic string, accountId uint32) (*DefaultEthSigner, error) {
	pk, err := BuildPrivateKeyFromMnemonic(mnemonic, accountId)
	if err != nil {
		return nil, err
	}
	return newEthSignerFromPk(pk)
}

// NewEthSignerFromRawPrivateKey 使用私钥创建一个账号信息
func NewEthSignerFromRawPrivateKey(rawPk []byte) (*DefaultEthSigner, error) {
	pk, err := crypto.ToECDSA(rawPk)
	if err != nil {
		return nil, errors.Wrap(err, "invalid raw private key")
	}
	return newEthSignerFromPk(pk)
}

// NewEthSignerFromHexPrivateKey 使用私钥创建一个账号信息
func NewEthSignerFromHexPrivateKey(hexPk string) (*DefaultEthSigner, error) {
	pk, err := BuildPrivateKeyFromHex(hexPk)
	if err != nil {
		return nil, errors.Wrap(err, "invalid hex private key")
	}
	return newEthSignerFromPk(pk)
}

// GetAddress 获得地址
func (s *DefaultEthSigner) GetAddress() common.Address {
	return s.address
}

// GetPrivateKey 获得私钥
func (s *DefaultEthSigner) GetPrivateKey() *ecdsa.PrivateKey {
	return s.pk
}
