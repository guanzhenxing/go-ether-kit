package etherkit

import (
	"encoding/hex"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

//############ Address ############

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

// PublicKeyBytesToAddress 公钥到地址
func PublicKeyBytesToAddress(publicKey []byte) common.Address {
	var buf []byte

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKey[1:]) // remove EC prefix 04
	buf = hash.Sum(nil)
	address := buf[12:]

	return common.HexToAddress(hex.EncodeToString(address))
}
