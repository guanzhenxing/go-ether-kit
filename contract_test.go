package etherkit

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestGetABI(t *testing.T) {
	tests := []struct {
		name        string
		abiString   string
		shouldError bool
	}{
		{
			name: "Valid ERC20 ABI",
			abiString: `[
				{
					"inputs": [{"name": "account", "type": "address"}],
					"name": "balanceOf",
					"outputs": [{"name": "", "type": "uint256"}],
					"stateMutability": "view",
					"type": "function"
				}
			]`,
			shouldError: false,
		},
		{
			name: "Valid transfer ABI",
			abiString: `[
				{
					"inputs": [
						{"name": "to", "type": "address"},
						{"name": "amount", "type": "uint256"}
					],
					"name": "transfer",
					"outputs": [{"name": "", "type": "bool"}],
					"stateMutability": "nonpayable",
					"type": "function"
				}
			]`,
			shouldError: false,
		},
		{
			name:        "Empty ABI",
			abiString:   "[]",
			shouldError: false,
		},
		{
			name:        "Invalid JSON",
			abiString:   "invalid json",
			shouldError: true,
		},
		{
			name:        "Invalid ABI structure",
			abiString:   `[{"invalid": "structure"}]`,
			shouldError: true,
		},
		{
			name:        "Empty string",
			abiString:   "",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			abi, err := GetABI(tt.abiString)

			if tt.shouldError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// 验证ABI对象可用
				if len(abi.Methods) == 0 && len(abi.Events) == 0 && tt.abiString != "[]" {
					t.Error("ABI should contain methods or events")
				}
			}
		})
	}
}

func TestGetContractMethodId(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		expected string
	}{
		{"ERC20 transfer", "transfer(address,uint256)", "0xa9059cbb"},
		{"ERC20 approve", "approve(address,uint256)", "0x095ea7b3"},
		{"ERC20 transferFrom", "transferFrom(address,address,uint256)", "0x23b872dd"},
		{"ERC20 balanceOf", "balanceOf(address)", "0x70a08231"},
		{"ERC20 totalSupply", "totalSupply()", "0x18160ddd"},
		{"ERC20 allowance", "allowance(address,address)", "0xdd62ed3e"},
		{"Simple function", "test()", "0xf8a8fd6d"},
		{"Function with multiple params", "complex(uint256,string,bool)", "0xd820ad92"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetContractMethodId(tt.method)

			if result != tt.expected {
				t.Errorf("GetContractMethodId(%s) = %s, expected %s",
					tt.method, result, tt.expected)
			}

			// 验证结果格式
			if !strings.HasPrefix(result, "0x") {
				t.Error("Method ID should start with 0x")
			}

			if len(result) != 10 { // 0x + 8 hex chars
				t.Errorf("Method ID length = %d, expected 10", len(result))
			}
		})
	}
}

func TestGetEventTopic(t *testing.T) {
	tests := []struct {
		name     string
		event    string
		expected string
	}{
		{"ERC20 Transfer", "Transfer(address,address,uint256)", "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"},
		{"ERC20 Approval", "Approval(address,address,uint256)", "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925"},
		{"Simple event", "Test()", "0xa163a6249e860c278ef4049759a7f7c7e8c141d30fd634fda9b5a6a95d111a30"},
		{"Event with indexed params", "Mint(address,uint256)", "0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetEventTopic(tt.event)

			if result != tt.expected {
				t.Errorf("GetEventTopic(%s) = %s, expected %s",
					tt.event, result, tt.expected)
			}

			// 验证结果格式
			if !strings.HasPrefix(result, "0x") {
				t.Error("Event topic should start with 0x")
			}

			if len(result) != 66 { // 0x + 64 hex chars
				t.Errorf("Event topic length = %d, expected 66", len(result))
			}
		})
	}
}

func TestBuildContractInputData(t *testing.T) {
	// 创建测试ABI
	abiString := `[
		{
			"inputs": [
				{"name": "to", "type": "address"},
				{"name": "amount", "type": "uint256"}
			],
			"name": "transfer",
			"outputs": [{"name": "", "type": "bool"}],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [{"name": "account", "type": "address"}],
			"name": "balanceOf",
			"outputs": [{"name": "", "type": "uint256"}],
			"stateMutability": "view",
			"type": "function"
		}
	]`

	abi, err := GetABI(abiString)
	if err != nil {
		t.Fatalf("Failed to parse ABI: %v", err)
	}

	tests := []struct {
		name        string
		method      string
		args        []interface{}
		shouldError bool
		expectedLen int
	}{
		{
			name:        "Valid transfer call",
			method:      "transfer",
			args:        []interface{}{common.HexToAddress("0x742F35C6dB4634C0532925a3b8D6dA2E"), big.NewInt(1000)},
			shouldError: false,
			expectedLen: 68, // 4 bytes method ID + 32 bytes address + 32 bytes uint256
		},
		{
			name:        "Valid balanceOf call",
			method:      "balanceOf",
			args:        []interface{}{common.HexToAddress("0x742F35C6dB4634C0532925a3b8D6dA2E")},
			shouldError: false,
			expectedLen: 36, // 4 bytes method ID + 32 bytes address
		},
		{
			name:        "Invalid method name",
			method:      "nonExistentMethod",
			args:        []interface{}{},
			shouldError: true,
			expectedLen: 0,
		},
		{
			name:        "Wrong number of arguments",
			method:      "transfer",
			args:        []interface{}{common.HexToAddress("0x742F35C6dB4634C0532925a3b8D6dA2E")}, // Missing amount
			shouldError: true,
			expectedLen: 0,
		},
		{
			name:        "Wrong argument type",
			method:      "transfer",
			args:        []interface{}{"invalid_address", big.NewInt(1000)},
			shouldError: true,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := BuildContractInputData(abi, tt.method, tt.args...)

			if tt.shouldError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if len(data) != 0 {
					t.Error("Expected empty data when error occurs")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if len(data) != tt.expectedLen {
					t.Errorf("Data length = %d, expected %d", len(data), tt.expectedLen)
				}

				// 验证方法ID是否正确
				if len(data) >= 4 {
					expectedMethodID := GetContractMethodId(getMethodSignature(abi, tt.method))
					actualMethodID := "0x" + common.Bytes2Hex(data[:4])

					if actualMethodID != expectedMethodID {
						t.Errorf("Method ID = %s, expected %s", actualMethodID, expectedMethodID)
					}
				}
			}
		})
	}
}

// 辅助函数：获取方法签名
func getMethodSignature(abi interface{}, methodName string) string {
	// 这里简化处理，实际项目中可以从ABI中提取完整签名
	switch methodName {
	case "transfer":
		return "transfer(address,uint256)"
	case "balanceOf":
		return "balanceOf(address)"
	default:
		return ""
	}
}

// 性能测试
func BenchmarkGetContractMethodId(b *testing.B) {
	method := "transfer(address,uint256)"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetContractMethodId(method)
	}
}

func BenchmarkGetEventTopic(b *testing.B) {
	event := "Transfer(address,address,uint256)"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetEventTopic(event)
	}
}

func BenchmarkGetABI(b *testing.B) {
	abiString := `[
		{
			"inputs": [{"name": "account", "type": "address"}],
			"name": "balanceOf",
			"outputs": [{"name": "", "type": "uint256"}],
			"stateMutability": "view",
			"type": "function"
		}
	]`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetABI(abiString)
	}
}

func BenchmarkBuildContractInputData(b *testing.B) {
	abiString := `[
		{
			"inputs": [
				{"name": "to", "type": "address"},
				{"name": "amount", "type": "uint256"}
			],
			"name": "transfer",
			"outputs": [{"name": "", "type": "bool"}],
			"stateMutability": "nonpayable",
			"type": "function"
		}
	]`

	abi, _ := GetABI(abiString)
	toAddress := common.HexToAddress("0x742F35C6dB4634C0532925a3b8D6dA2E")
	amount := big.NewInt(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = BuildContractInputData(abi, "transfer", toAddress, amount)
	}
}
