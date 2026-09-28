package money

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestZero(t *testing.T) {
	t.Run("creates zero money", func(t *testing.T) {
		value, err := Zero("brl")

		require.NoError(t, err)
		require.Equal(t, int64(0), value.AmountMinor())
		require.Equal(t, "BRL", value.Currency())
		require.True(t, value.IsZero())
		require.False(t, value.IsNegative())
		require.Equal(t, "0.00", value.String())
	})

	t.Run("rejects invalid currency", func(t *testing.T) {
		_, err := Zero("BR")

		require.ErrorIs(t, err, ErrInvalidCurrency)
	})
}

func TestFromMinor(t *testing.T) {
	t.Run("creates money from minor units", func(t *testing.T) {
		value, err := FromMinor(12345, "BRL")

		require.NoError(t, err)
		require.Equal(t, int64(12345), value.AmountMinor())
		require.Equal(t, "BRL", value.Currency())
		require.Equal(t, "123.45", value.String())
	})

	t.Run("accepts negative value", func(t *testing.T) {
		value, err := FromMinor(-12345, "BRL")

		require.NoError(t, err)
		require.Equal(t, int64(-12345), value.AmountMinor())
		require.True(t, value.IsNegative())
		require.Equal(t, "-123.45", value.String())
	})

	t.Run("accepts zero", func(t *testing.T) {
		value, err := FromMinor(0, "BRL")

		require.NoError(t, err)
		require.True(t, value.IsZero())
	})
}

func TestFromDecimalString(t *testing.T) {
	tests := []struct {
		name       string
		amount     string
		currency   string
		allowNeg   bool
		wantMinor  int64
		wantString string
	}{
		{
			name:       "integer",
			amount:     "100",
			currency:   "BRL",
			wantMinor:  10000,
			wantString: "100.00",
		},
		{
			name:       "one decimal",
			amount:     "100.5",
			currency:   "BRL",
			wantMinor:  10050,
			wantString: "100.50",
		},
		{
			name:       "two decimals",
			amount:     "100.50",
			currency:   "BRL",
			wantMinor:  10050,
			wantString: "100.50",
		},
		{
			name:       "leading plus",
			amount:     "+100.50",
			currency:   "BRL",
			wantMinor:  10050,
			wantString: "100.50",
		},
		{
			name:       "leading and trailing spaces",
			amount:     " 100.50 ",
			currency:   "BRL",
			wantMinor:  10050,
			wantString: "100.50",
		},
		{
			name:       "zero",
			amount:     "0",
			currency:   "BRL",
			wantMinor:  0,
			wantString: "0.00",
		},
		{
			name:       "negative when allowed",
			amount:     "-100.50",
			currency:   "BRL",
			allowNeg:   true,
			wantMinor:  -10050,
			wantString: "-100.50",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := FromDecimalString(
				tt.amount,
				tt.currency,
				tt.allowNeg,
			)

			require.NoError(t, err)
			require.Equal(t, tt.wantMinor, value.AmountMinor())
			require.Equal(t, tt.wantString, value.String())
			require.Equal(t, "BRL", value.Currency())
		})
	}
}

func TestFromDecimalString_Negative(t *testing.T) {
	t.Run("rejects negative when not allowed", func(t *testing.T) {
		_, err := FromDecimalString(
			"-10.50",
			"BRL",
			false,
		)

		require.ErrorIs(t, err, ErrNegativeNotAllowed)
	})

	t.Run("accepts negative when allowed", func(t *testing.T) {
		value, err := FromDecimalString(
			"-10.50",
			"BRL",
			true,
		)

		require.NoError(t, err)
		require.Equal(t, int64(-1050), value.AmountMinor())
	})
}

func TestFromDecimalString_InvalidAmount(t *testing.T) {
	tests := []string{
		"",
		"abc",
		"1.234",
		"1.2.3",
		"1,23",
		"1e2",
		"1E2",
		"NaN",
		"Inf",
		"+Inf",
		"-Inf",
		".50",
		".",
		"-",
		"+",
	}

	for _, amount := range tests {
		t.Run(amount, func(t *testing.T) {
			_, err := FromDecimalString(
				amount,
				"BRL",
				true,
			)

			require.ErrorIs(t, err, ErrInvalidAmount)
		})
	}
}

func TestFromDecimalString_RejectsExcessScale(t *testing.T) {
	tests := []string{
		"10.001",
		"10.123",
		"0.001",
		"999.999",
	}

	for _, amount := range tests {
		t.Run(amount, func(t *testing.T) {
			_, err := FromDecimalString(
				amount,
				"BRL",
				false,
			)

			require.ErrorIs(t, err, ErrInvalidAmount)
		})
	}
}

func TestCurrency(t *testing.T) {
	t.Run("normalizes lowercase currency", func(t *testing.T) {
		value, err := FromMinor(100, "brl")

		require.NoError(t, err)
		require.Equal(t, "BRL", value.Currency())
	})

	t.Run("trims currency", func(t *testing.T) {
		value, err := FromMinor(100, " BRL ")

		require.NoError(t, err)
		require.Equal(t, "BRL", value.Currency())
	})

	t.Run("rejects invalid currencies", func(t *testing.T) {
		tests := []string{
			"",
			"BR",
			"BRLL",
			"123",
			"R$L",
			"BR1",
			"b r l",
		}

		for _, currency := range tests {
			t.Run(currency, func(t *testing.T) {
				_, err := FromMinor(100, currency)

				require.ErrorIs(t, err, ErrInvalidCurrency)
			})
		}
	})
}

func TestMoney_Add(t *testing.T) {
	t.Run("adds values with same currency", func(t *testing.T) {
		left, err := FromMinor(1000, "BRL")
		require.NoError(t, err)

		right, err := FromMinor(250, "BRL")
		require.NoError(t, err)

		result, err := left.Add(right)

		require.NoError(t, err)
		require.Equal(t, int64(1250), result.AmountMinor())
		require.Equal(t, "12.50", result.String())
	})

	t.Run("rejects different currencies", func(t *testing.T) {
		left, err := FromMinor(1000, "BRL")
		require.NoError(t, err)

		right, err := FromMinor(250, "USD")
		require.NoError(t, err)

		_, err = left.Add(right)

		require.ErrorIs(t, err, ErrCurrencyMismatch)
	})
}

func TestMoney_Sub(t *testing.T) {
	t.Run("subtracts values with same currency", func(t *testing.T) {
		left, err := FromMinor(1000, "BRL")
		require.NoError(t, err)

		right, err := FromMinor(250, "BRL")
		require.NoError(t, err)

		result, err := left.Sub(right)

		require.NoError(t, err)
		require.Equal(t, int64(750), result.AmountMinor())
		require.Equal(t, "7.50", result.String())
	})

	t.Run("supports negative result", func(t *testing.T) {
		left, err := FromMinor(100, "BRL")
		require.NoError(t, err)

		right, err := FromMinor(250, "BRL")
		require.NoError(t, err)

		result, err := left.Sub(right)

		require.NoError(t, err)
		require.Equal(t, int64(-150), result.AmountMinor())
		require.Equal(t, "-1.50", result.String())
	})

	t.Run("rejects different currencies", func(t *testing.T) {
		left, err := FromMinor(1000, "BRL")
		require.NoError(t, err)

		right, err := FromMinor(250, "USD")
		require.NoError(t, err)

		_, err = left.Sub(right)

		require.ErrorIs(t, err, ErrCurrencyMismatch)
	})
}

func TestMoney_Compare(t *testing.T) {
	tests := []struct {
		name     string
		left     int64
		right    int64
		expected int
	}{
		{
			name:     "less",
			left:     100,
			right:    200,
			expected: -1,
		},
		{
			name:     "equal",
			left:     100,
			right:    100,
			expected: 0,
		},
		{
			name:     "greater",
			left:     200,
			right:    100,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left, err := FromMinor(tt.left, "BRL")
			require.NoError(t, err)

			right, err := FromMinor(tt.right, "BRL")
			require.NoError(t, err)

			result, err := left.Compare(right)

			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}

	t.Run("rejects different currencies", func(t *testing.T) {
		left, err := FromMinor(100, "BRL")
		require.NoError(t, err)

		right, err := FromMinor(100, "USD")
		require.NoError(t, err)

		_, err = left.Compare(right)

		require.ErrorIs(t, err, ErrCurrencyMismatch)
	})
}

func TestMoney_Negate(t *testing.T) {
	t.Run("negates positive value", func(t *testing.T) {
		value, err := FromMinor(12345, "BRL")
		require.NoError(t, err)

		negative, err := value.Negate()

		require.NoError(t, err)
		require.Equal(t, int64(-12345), negative.AmountMinor())
		require.Equal(t, "BRL", negative.Currency())
		require.Equal(t, "-123.45", negative.String())
	})

	t.Run("negates negative value", func(t *testing.T) {
		value, err := FromMinor(-12345, "BRL")
		require.NoError(t, err)

		positive, err := value.Negate()

		require.NoError(t, err)
		require.Equal(t, int64(12345), positive.AmountMinor())
		require.Equal(t, "BRL", positive.Currency())
		require.Equal(t, "123.45", positive.String())
	})

	t.Run("rejects MinInt64", func(t *testing.T) {
		value, err := FromMinor(math.MinInt64, "BRL")
		require.NoError(t, err)

		_, err = value.Negate()

		require.ErrorIs(t, err, ErrOverflow)
	})
}

func TestMoney_AddOverflow(t *testing.T) {
	t.Run("positive overflow", func(t *testing.T) {
		left, err := FromMinor(math.MaxInt64, "BRL")
		require.NoError(t, err)

		right, err := FromMinor(1, "BRL")
		require.NoError(t, err)

		_, err = left.Add(right)

		require.ErrorIs(t, err, ErrOverflow)
	})

	t.Run("negative overflow", func(t *testing.T) {
		left, err := FromMinor(math.MinInt64, "BRL")
		require.NoError(t, err)

		right, err := FromMinor(-1, "BRL")
		require.NoError(t, err)

		_, err = left.Add(right)

		require.ErrorIs(t, err, ErrOverflow)
	})
}

func TestFromDecimalString_Overflow(t *testing.T) {
	_, err := FromDecimalString(
		"92233720368547758.08",
		"BRL",
		false,
	)

	require.ErrorIs(t, err, ErrOverflow)
}

func TestMoney_MinInt64(t *testing.T) {
	value, err := FromMinor(math.MinInt64, "BRL")
	require.NoError(t, err)

	require.Equal(t, int64(math.MinInt64), value.AmountMinor())

	require.Equal(
		t,
		"-92233720368547758.08",
		value.String(),
	)
}

func TestMoney_NegateMinInt64(t *testing.T) {
	value, err := FromMinor(math.MinInt64, "BRL")
	require.NoError(t, err)

	_, err = value.Negate()

	require.ErrorIs(t, err, ErrOverflow)
}

func TestMoney_SubOverflow(t *testing.T) {
	left, err := FromMinor(0, "BRL")
	require.NoError(t, err)

	right, err := FromMinor(math.MinInt64, "BRL")
	require.NoError(t, err)

	_, err = left.Sub(right)

	require.ErrorIs(t, err, ErrOverflow)
}
