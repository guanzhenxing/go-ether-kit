package etherkit

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

// Signer ETH账户
type Signer interface {
	GetAddress() common.Address
	GetPrivateKey() *ecdsa.PrivateKey
}

type EtherSigner struct {
	pk      *ecdsa.PrivateKey
	address common.Address
}

// NewEtherSigner 创建一个默认的账号信息
func NewEtherSigner() (*EtherSigner, error) {
	pk, err := GeneratePrivateKey()
	if err != nil {
		return nil, errors.Wrap(err, "generate private key error")
	}
	return NewEtherSignerFromPrivateKey(pk)
}

// NewEtherSignerFromPrivateKey 使用私钥创建一个账号信息
func NewEtherSignerFromPrivateKey(pk *ecdsa.PrivateKey) (*EtherSigner, error) {
	return &EtherSigner{
		pk:      pk,
		address: PrivateKeyToAddress(pk),
	}, nil
}

// NewEtherSignerFromMnemonic 使用助记词创建一个账号信息
func NewEtherSignerFromMnemonic(mnemonic string) (*EtherSigner, error) {
	return NewEtherSignerFromMnemonicAndAccountId(mnemonic, 0)
}

// NewEtherSignerFromMnemonicAndAccountId 使用助记词创建一个账号信息
func NewEtherSignerFromMnemonicAndAccountId(mnemonic string, accountId uint32) (*EtherSigner, error) {
	pk, err := BuildPrivateKeyFromMnemonic(mnemonic, accountId)
	if err != nil {
		return nil, err
	}
	return NewEtherSignerFromPrivateKey(pk)
}

// NewEthSignerFromRawPrivateKey 使用私钥创建一个账号信息
func NewEthSignerFromRawPrivateKey(rawPk []byte) (*EtherSigner, error) {
	pk, err := crypto.ToECDSA(rawPk)
	if err != nil {
		return nil, errors.Wrap(err, "invalid raw private key")
	}
	return NewEtherSignerFromPrivateKey(pk)
}

// NewEtherSignerFromHexPrivateKey 使用私钥创建一个账号信息
func NewEtherSignerFromHexPrivateKey(hexPk string) (*EtherSigner, error) {
	pk, err := BuildPrivateKeyFromHex(hexPk)
	if err != nil {
		return nil, errors.Wrap(err, "invalid hex private key")
	}
	return NewEtherSignerFromPrivateKey(pk)
}

// GetAddress 获得地址
func (s *EtherSigner) GetAddress() common.Address {
	return s.address
}

// GetPrivateKey 获得私钥
func (s *EtherSigner) GetPrivateKey() *ecdsa.PrivateKey {
	return s.pk
}
