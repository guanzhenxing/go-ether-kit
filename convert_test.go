package etherkit

import (
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
)

func TestToWei(t *testing.T) {
	tests := []struct {
		name     string
		amount   interface{}
		decimals int
		expected string // 使用字符串比较避免大数精度问题
	}{
		{"String amount - 1 ETH", "1", 18, "1000000000000000000"},
		{"String amount - 0.5 ETH", "0.5", 18, "500000000000000000"},
		{"String amount - 1.5 ETH", "1.5", 18, "1500000000000000000"},
		{"Float64 amount - 1.0", 1.0, 18, "1000000000000000000"},
		{"Float64 amount - 0.1", 0.1, 18, "100000000000000000"},
		{"Int64 amount - 5", int64(5), 18, "5000000000000000000"},
		{"Int amount - 10", 10, 18, "10000000000000000000"},
		{"Decimal amount", decimal.NewFromFloat(2.5), 18, "2500000000000000000"},
		{"USDC amount - 100", "100", 6, "100000000"},
		{"USDT amount - 50.5", "50.5", 6, "50500000"},
		{"Zero amount", "0", 18, "0"},
		{"Small amount", "0.000000000000000001", 18, "1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToWei(tt.amount, tt.decimals)
			if result.String() != tt.expected {
				t.Errorf("ToWei(%v, %d) = %s, expected %s",
					tt.amount, tt.decimals, result.String(), tt.expected)
			}
		})
	}
}

func TestToDecimal(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		decimals int
		expected string
	}{
		{"1 ETH in Wei", "1000000000000000000", 18, "1"},
		{"0.5 ETH in Wei", "500000000000000000", 18, "0.5"},
		{"Big Int - 1.5 ETH", big.NewInt(1500000000000000000), 18, "1.5"},
		{"100 USDC in Wei", "100000000", 6, "100"},
		{"50.5 USDT in Wei", "50500000", 6, "50.5"},
		{"Zero value", "0", 18, "0"},
		{"Small value", "1", 18, "0"},
		{"Large value", "1000000000000000000000", 18, "1000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToDecimal(tt.value, tt.decimals)
			if result.String() != tt.expected {
				t.Errorf("ToDecimal(%v, %d) = %s, expected %s",
					tt.value, tt.decimals, result.String(), tt.expected)
			}
		})
	}
}

func TestToWeiToDecimalRoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		amount   string
		decimals int
	}{
		{"1 ETH", "1", 18},
		{"0.5 ETH", "0.5", 18},
		{"1.5 ETH", "1.5", 18},
		{"100 USDC", "100", 6},
		{"50.5 USDT", "50.5", 6},
		{"0.000001 ETH", "0.000001", 18},
		{"1000 tokens", "1000", 18},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Amount -> Wei -> Amount
			wei := ToWei(tt.amount, tt.decimals)
			result := ToDecimal(wei, tt.decimals)

			if result.String() != tt.amount {
				t.Errorf("Round trip failed for %s: got %s", tt.amount, result.String())
			}
		})
	}
}

func TestToWeiEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		amount   interface{}
		decimals int
		expected string
	}{
		{"Very small amount", "0.000000000000000001", 18, "1"},
		{"Very large amount", "999999999999999999", 18, "999999999999999999000000000000000000"},
		{"Zero decimals", "123", 0, "123"},
		{"High decimals", "1", 30, "1000000000000000000000000000000"},
		{"Negative-like string", "-0", 18, "0"}, // decimal包会处理这种情况
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToWei(tt.amount, tt.decimals)
			if result.String() != tt.expected {
				t.Errorf("ToWei(%v, %d) = %s, expected %s",
					tt.amount, tt.decimals, result.String(), tt.expected)
			}
		})
	}
}

func TestToDecimalEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		decimals int
		expected string
	}{
		{"Very large Wei value", "999999999999999999000000000000000000", 18, "999999999999999999"},
		{"Zero decimals", "123", 0, "123"},
		{"High decimals", "1000000000000000000000000000000", 30, "1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToDecimal(tt.value, tt.decimals)
			if result.String() != tt.expected {
				t.Errorf("ToDecimal(%v, %d) = %s, expected %s",
					tt.value, tt.decimals, result.String(), tt.expected)
			}
		})
	}
}

// 性能测试
func BenchmarkToWei(b *testing.B) {
	amount := "1.5"
	decimals := 18

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToWei(amount, decimals)
	}
}

func BenchmarkToDecimal(b *testing.B) {
	value := big.NewInt(1500000000000000000) // 1.5 ETH in Wei
	decimals := 18

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToDecimal(value, decimals)
	}
}

func BenchmarkToWeiFloat64(b *testing.B) {
	amount := 1.5
	decimals := 18

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToWei(amount, decimals)
	}
}

func BenchmarkToWeiString(b *testing.B) {
	amount := "1.5"
	decimals := 18

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToWei(amount, decimals)
	}
}

func BenchmarkToWeiDecimal(b *testing.B) {
	amount := decimal.NewFromFloat(1.5)
	decimals := 18

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToWei(amount, decimals)
	}
}
