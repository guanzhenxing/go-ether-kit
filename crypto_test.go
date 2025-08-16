package etherkit

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestGeneratePrivateKey(t *testing.T) {
	// 测试生成私钥
	pk, err := GeneratePrivateKey()
	if err != nil {
		t.Fatalf("GeneratePrivateKey() failed: %v", err)
	}

	if pk == nil {
		t.Fatal("Generated private key is nil")
	}

	// 验证私钥不为空即可，类型已经由函数签名保证

	// 验证可以生成地址
	address := PrivateKeyToAddress(pk)
	if !IsValidAddress(address) {
		t.Error("Generated address from private key is invalid")
	}
}

func TestPrivateKeyToAddress(t *testing.T) {
	// 使用已知的私钥和地址对进行测试
	testPrivateKeyHex := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	expectedAddress := "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"

	pk, err := BuildPrivateKeyFromHex(testPrivateKeyHex)
	if err != nil {
		t.Fatalf("BuildPrivateKeyFromHex() failed: %v", err)
	}

	address := PrivateKeyToAddress(pk)

	// 比较地址（忽略大小写）
	if !strings.EqualFold(address.Hex(), expectedAddress) {
		t.Errorf("PrivateKeyToAddress() = %s, expected %s", address.Hex(), expectedAddress)
	}
}

func TestGetHexPrivateKey(t *testing.T) {
	// 生成私钥
	pk, err := GeneratePrivateKey()
	if err != nil {
		t.Fatalf("GeneratePrivateKey() failed: %v", err)
	}

	hexPk := GetHexPrivateKey(pk)

	// 验证十六进制格式
	if len(hexPk) != 64 {
		t.Errorf("Hex private key length = %d, expected 64", len(hexPk))
	}

	// 验证不包含0x前缀
	if strings.HasPrefix(hexPk, "0x") {
		t.Error("Hex private key should not have 0x prefix")
	}

	// 验证可以重新构建私钥
	rebuiltPk, err := BuildPrivateKeyFromHex(hexPk)
	if err != nil {
		t.Errorf("Failed to rebuild private key from hex: %v", err)
	}

	// 验证重建的私钥生成相同地址
	originalAddress := PrivateKeyToAddress(pk)
	rebuiltAddress := PrivateKeyToAddress(rebuiltPk)

	if originalAddress.Hex() != rebuiltAddress.Hex() {
		t.Error("Rebuilt private key generates different address")
	}
}

func TestGetHexPublicKey(t *testing.T) {
	pk, err := GeneratePrivateKey()
	if err != nil {
		t.Fatalf("GeneratePrivateKey() failed: %v", err)
	}

	hexPubKey := GetHexPublicKey(pk)

	// 验证公钥长度 (64字节 = 128个十六进制字符)
	if len(hexPubKey) != 128 {
		t.Errorf("Hex public key length = %d, expected 128", len(hexPubKey))
	}

	// 验证不包含0x和04前缀
	if strings.HasPrefix(hexPubKey, "0x") || strings.HasPrefix(hexPubKey, "04") {
		t.Error("Hex public key should not have 0x or 04 prefix")
	}
}

func TestBuildPrivateKeyFromHex(t *testing.T) {
	tests := []struct {
		name        string
		hexKey      string
		shouldError bool
	}{
		{"Valid hex key without 0x", "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", false},
		{"Valid hex key with 0x", "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", false},
		{"Invalid hex key - too short", "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff", true},
		{"Invalid hex key - too long", "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff8012", true},
		{"Invalid hex key - invalid characters", "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ffXX", true},
		{"Empty string", "", true},
		{"Invalid format", "not_a_hex_key", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pk, err := BuildPrivateKeyFromHex(tt.hexKey)

			if tt.shouldError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if pk != nil {
					t.Error("Expected nil private key when error occurs")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if pk == nil {
					t.Error("Expected valid private key but got nil")
				}
			}
		})
	}
}

func TestBuildPrivateKeyFromMnemonic(t *testing.T) {
	// 测试助记词（这是一个测试助记词，不要在生产环境使用）
	testMnemonic := "test test test test test test test test test test test junk"

	pk, err := BuildPrivateKeyFromMnemonic(testMnemonic)
	if err != nil {
		t.Fatalf("BuildPrivateKeyFromMnemonic() failed: %v", err)
	}

	if pk == nil {
		t.Fatal("Generated private key is nil")
	}

	// 验证地址生成
	address := PrivateKeyToAddress(pk)
	if !IsValidAddress(address) {
		t.Error("Generated address from mnemonic is invalid")
	}
}

func TestBuildPrivateKeyFromMnemonicAndAccountId(t *testing.T) {
	testMnemonic := "test test test test test test test test test test test junk"

	// 测试不同的账户ID应该生成不同的私钥
	pk1, err := BuildPrivateKeyFromMnemonicAndAccountId(testMnemonic, 0)
	if err != nil {
		t.Fatalf("BuildPrivateKeyFromMnemonicAndAccountId(0) failed: %v", err)
	}

	pk2, err := BuildPrivateKeyFromMnemonicAndAccountId(testMnemonic, 1)
	if err != nil {
		t.Fatalf("BuildPrivateKeyFromMnemonicAndAccountId(1) failed: %v", err)
	}

	// 验证不同账户ID生成不同私钥
	addr1 := PrivateKeyToAddress(pk1)
	addr2 := PrivateKeyToAddress(pk2)

	if addr1.Hex() == addr2.Hex() {
		t.Error("Different account IDs should generate different addresses")
	}

	// 验证相同账户ID生成相同私钥
	pk3, err := BuildPrivateKeyFromMnemonicAndAccountId(testMnemonic, 0)
	if err != nil {
		t.Fatalf("BuildPrivateKeyFromMnemonicAndAccountId(0) second call failed: %v", err)
	}

	addr3 := PrivateKeyToAddress(pk3)
	if addr1.Hex() != addr3.Hex() {
		t.Error("Same account ID should generate same address")
	}
}

func TestVerifySignature(t *testing.T) {
	// 生成测试私钥
	pk, err := GeneratePrivateKey()
	if err != nil {
		t.Fatalf("GeneratePrivateKey() failed: %v", err)
	}

	address := PrivateKeyToAddress(pk)
	testData := []byte("Hello, Ethereum!")

	// 生成签名
	hash := crypto.Keccak256Hash(testData)
	signature, err := crypto.Sign(hash.Bytes(), pk)
	if err != nil {
		t.Fatalf("crypto.Sign() failed: %v", err)
	}

	// 测试正确的签名验证
	isValid := VerifySignature(address.Hex(), testData, signature)
	if !isValid {
		t.Error("Valid signature verification failed")
	}

	// 测试错误的数据
	wrongData := []byte("Wrong data")
	isValid = VerifySignature(address.Hex(), wrongData, signature)
	if isValid {
		t.Error("Invalid signature verification should fail")
	}

	// 测试错误的地址
	wrongAddress := "0x742F35C6dB4634C0532925a3b8D6dA2E"
	isValid = VerifySignature(wrongAddress, testData, signature)
	if isValid {
		t.Error("Wrong address signature verification should fail")
	}

	// 测试错误的签名
	wrongSignature := make([]byte, len(signature))
	copy(wrongSignature, signature)
	wrongSignature[0] ^= 0xFF // 修改第一个字节
	isValid = VerifySignature(address.Hex(), testData, wrongSignature)
	if isValid {
		t.Error("Wrong signature verification should fail")
	}
}

// 性能测试
func BenchmarkGeneratePrivateKey(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GeneratePrivateKey()
	}
}

func BenchmarkPrivateKeyToAddress(b *testing.B) {
	pk, _ := GeneratePrivateKey()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PrivateKeyToAddress(pk)
	}
}

func BenchmarkBuildPrivateKeyFromHex(b *testing.B) {
	hexKey := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = BuildPrivateKeyFromHex(hexKey)
	}
}

func BenchmarkVerifySignature(b *testing.B) {
	pk, _ := GeneratePrivateKey()
	address := PrivateKeyToAddress(pk)
	testData := []byte("Hello, Ethereum!")
	hash := crypto.Keccak256Hash(testData)
	signature, _ := crypto.Sign(hash.Bytes(), pk)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		VerifySignature(address.Hex(), testData, signature)
	}
}
