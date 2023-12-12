package etherkit

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

// EtherSigner ETH账户
type EtherSigner interface {
	GetAddress() common.Address
	GetPrivateKey() *ecdsa.PrivateKey
}

type Signer struct {
	pk      *ecdsa.PrivateKey
	address common.Address
}

// NewSigner 创建一个默认的账号信息
func NewSigner() (*Signer, error) {
	pk, err := GeneratePrivateKey()
	if err != nil {
		return nil, errors.Wrap(err, "generate private key error")
	}
	return NewSignerFromPrivateKey(pk)
}

// NewSignerFromPrivateKey 使用私钥创建一个账号信息
func NewSignerFromPrivateKey(pk *ecdsa.PrivateKey) (*Signer, error) {
	return &Signer{
		pk:      pk,
		address: PrivateKeyToAddress(pk),
	}, nil
}

// NewSignerFromMnemonic 使用助记词创建一个账号信息
func NewSignerFromMnemonic(mnemonic string) (*Signer, error) {
	return NewSignerFromMnemonicAndAccountId(mnemonic, 0)
}

// NewSignerFromMnemonicAndAccountId 使用助记词创建一个账号信息
func NewSignerFromMnemonicAndAccountId(mnemonic string, accountId uint32) (*Signer, error) {
	pk, err := BuildPrivateKeyFromMnemonic(mnemonic, accountId)
	if err != nil {
		return nil, err
	}
	return NewSignerFromPrivateKey(pk)
}

// NewSignerFromRawPrivateKey 使用私钥创建一个账号信息
func NewSignerFromRawPrivateKey(rawPk []byte) (*Signer, error) {
	pk, err := crypto.ToECDSA(rawPk)
	if err != nil {
		return nil, errors.Wrap(err, "invalid raw private key")
	}
	return NewSignerFromPrivateKey(pk)
}

// NewSignerFromHexPrivateKey 使用私钥创建一个账号信息
func NewSignerFromHexPrivateKey(hexPk string) (*Signer, error) {
	pk, err := BuildPrivateKeyFromHex(hexPk)
	if err != nil {
		return nil, errors.Wrap(err, "invalid hex private key")
	}
	return NewSignerFromPrivateKey(pk)
}

// GetAddress 获得地址
func (s *Signer) GetAddress() common.Address {
	return s.address
}

// GetPrivateKey 获得私钥
func (s *Signer) GetPrivateKey() *ecdsa.PrivateKey {
	return s.pk
}
