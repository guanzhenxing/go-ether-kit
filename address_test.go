package etherkit

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestIsValidAddress(t *testing.T) {
	tests := []struct {
		name     string
		address  interface{}
		expected bool
	}{
		{"Valid address string", "0x742F35C6dB4634C0532925a3b8D6dA2E12345678", true},
		{"Valid address lowercase", "0x742f35c6db4634c0532925a3b8d6da2e12345678", true},
		{"Valid address mixed case", "0x742F35c6dB4634C0532925a3b8D6dA2E12345678", true},
		{"Valid zero address", ZeroAddress, true},
		{"Valid native token address", NativeTokenAddress, true},
		{"Valid common.Address", common.HexToAddress("0x742F35C6dB4634C0532925a3b8D6dA2E12345678"), true},
		{"Invalid address - no 0x prefix", "742F35C6dB4634C0532925a3b8D6dA2E12345678", false},
		{"Invalid address - wrong length", "0x742F35C6dB4634C0532925a3b8D6dA2E123", false},
		{"Invalid address - too short", "0x742F35C6dB4634C0532925a3b8D6dA", false},
		{"Invalid address - empty string", "", false},
		{"Invalid address - random string", "not_an_address", false},
		{"Invalid address - contains invalid chars", "0x742F35C6dB4634C0532925a3b8D6dA2G12345678", false},
		{"Invalid type - integer", 123, false},
		{"Invalid type - nil", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidAddress(tt.address)
			if result != tt.expected {
				t.Errorf("IsValidAddress(%v) = %v, expected %v", tt.address, result, tt.expected)
			}
		})
	}
}

func TestPublicKeyBytesToAddress(t *testing.T) {
	// 测试用的公钥 (不包含04前缀)
	publicKeyBytes := []byte{
		0x04, // EC prefix
		0x79, 0xbe, 0x66, 0x7e, 0xf9, 0xdc, 0xbb, 0xac, 0x55, 0xa0, 0x62, 0x95, 0xce, 0x87, 0x0b, 0x07,
		0x02, 0x9b, 0xfc, 0xdb, 0x2d, 0xce, 0x28, 0xd9, 0x59, 0xf2, 0x81, 0x5b, 0x16, 0xf8, 0x17, 0x98,
		0x48, 0x3a, 0xda, 0x77, 0x26, 0xa3, 0xc4, 0x65, 0x5d, 0xa4, 0xfb, 0xfc, 0x0e, 0x11, 0x08, 0xa8,
		0xfd, 0x17, 0xb4, 0x48, 0xa6, 0x85, 0x54, 0x19, 0x9c, 0x47, 0xd0, 0x8f, 0xfb, 0x10, 0xd4, 0xb8,
	}

	address := PublicKeyBytesToAddress(publicKeyBytes)

	// 验证生成的地址格式是否正确
	if !IsValidAddress(address) {
		t.Errorf("Generated address is invalid: %s", address.Hex())
	}

	// 验证地址不是零地址
	if address.Hex() == ZeroAddress {
		t.Error("Generated address should not be zero address")
	}
}

// 性能测试
func BenchmarkIsValidAddress(b *testing.B) {
	address := "0x742F35C6dB4634C0532925a3b8D6dA2E"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsValidAddress(address)
	}
}

func BenchmarkIsValidAddressInvalid(b *testing.B) {
	address := "invalid_address"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsValidAddress(address)
	}
}

func BenchmarkPublicKeyBytesToAddress(b *testing.B) {
	publicKeyBytes := []byte{
		0x04,
		0x79, 0xbe, 0x66, 0x7e, 0xf9, 0xdc, 0xbb, 0xac, 0x55, 0xa0, 0x62, 0x95, 0xce, 0x87, 0x0b, 0x07,
		0x02, 0x9b, 0xfc, 0xdb, 0x2d, 0xce, 0x28, 0xd9, 0x59, 0xf2, 0x81, 0x5b, 0x16, 0xf8, 0x17, 0x98,
		0x48, 0x3a, 0xda, 0x77, 0x26, 0xa3, 0xc4, 0x65, 0x5d, 0xa4, 0xfb, 0xfc, 0x0e, 0x11, 0x08, 0xa8,
		0xfd, 0x17, 0xb4, 0x48, 0xa6, 0x85, 0x54, 0x19, 0x9c, 0x47, 0xd0, 0x8f, 0xfb, 0x10, 0xd4, 0xb8,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PublicKeyBytesToAddress(publicKeyBytes)
	}
}
