package money

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

const scale = 100 // fixed 2 decimal places

type Money struct {
	AmountMinor int64
	Currency    string
}

func Zero(currency string) (Money, error) {
	c, err := normalizeCurrency(currency)
	if err != nil {
		return Money{}, err
	}
	return Money{AmountMinor: 0, Currency: c}, nil
}

func FromMinor(amountMinor int64, currency string) (Money, error) {
	c, err := normalizeCurrency(currency)
	if err != nil {
		return Money{}, err
	}
	return Money{AmountMinor: amountMinor, Currency: c}, nil
}

func FromDecimalString(amount, currency string, allowNegative bool) (Money, error) {
	c, err := normalizeCurrency(currency)
	if err != nil {
		return Money{}, err
	}

	s := strings.TrimSpace(amount)
	if s == "" {
		return Money{}, fmt.Errorf("%w: empty string", ErrInvalidAmount)
	}
	lower := strings.ToLower(s)
	if strings.ContainsAny(lower, "e") || lower == "nan" || strings.Contains(lower, "inf") {
		return Money{}, fmt.Errorf("%w: %q", ErrInvalidAmount, amount)
	}

	negative := false
	switch {
	case strings.HasPrefix(s, "-"):
		negative = true
		s = s[1:]
	case strings.HasPrefix(s, "+"):
		s = s[1:]
	}
	if s == "" || !isDigitsAndDot(s) {
		return Money{}, fmt.Errorf("%w: %q", ErrInvalidAmount, amount)
	}

	parts := strings.SplitN(s, ".", 2)
	intPart := parts[0]
	fracPart := ""
	if len(parts) == 2 {
		fracPart = parts[1]
	}
	if intPart == "" {
		return Money{}, fmt.Errorf("%w: %q", ErrInvalidAmount, amount)
	}
	if len(fracPart) > 2 {
		return Money{}, fmt.Errorf("%w: excess scale in %q", ErrInvalidAmount, amount)
	}
	for len(fracPart) < 2 {
		fracPart += "0"
	}

	intValue, err := strconv.ParseInt(intPart, 10, 64)
	if err != nil {
		return Money{}, fmt.Errorf("%w: %q", ErrInvalidAmount, amount)
	}
	fracValue, err := strconv.ParseInt(fracPart, 10, 64)
	if err != nil {
		return Money{}, fmt.Errorf("%w: %q", ErrInvalidAmount, amount)
	}

	if intValue > (math.MaxInt64-fracValue)/scale {
		return Money{}, ErrOverflow
	}
	minor := intValue*scale + fracValue
	if negative {
		if !allowNegative {
			return Money{}, ErrNegativeNotAllowed
		}
		minor = -minor
	}

	return Money{AmountMinor: minor, Currency: c}, nil
}

func isDigitsAndDot(s string) bool {
	dots := 0
	for _, r := range s {
		if r == '.' {
			dots++
			if dots > 1 {
				return false
			}
			continue
		}
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func normalizeCurrency(currency string) (string, error) {
	c := strings.ToUpper(strings.TrimSpace(currency))
	if len(c) != 3 {
		return "", fmt.Errorf("%w: %q", ErrInvalidCurrency, currency)
	}

	for _, r := range c {
		if r < 'A' || r > 'Z' {
			return "", fmt.Errorf("%w: %q", ErrInvalidCurrency, currency)
		}
	}

	return c, nil
}

func (m Money) IsZero() bool     { return m.AmountMinor == 0 }
func (m Money) IsNegative() bool { return m.AmountMinor < 0 }

func (m Money) sameCurrency(other Money) error {
	if m.Currency != other.Currency {
		return fmt.Errorf("%w: %s vs %s", ErrCurrencyMismatch, m.Currency, other.Currency)
	}
	return nil
}

func (m Money) Add(other Money) (Money, error) {
	if err := m.sameCurrency(other); err != nil {
		return Money{}, err
	}
	sum := m.AmountMinor + other.AmountMinor
	if (other.AmountMinor > 0 && sum < m.AmountMinor) || (other.AmountMinor < 0 && sum > m.AmountMinor) {
		return Money{}, ErrOverflow
	}
	return Money{AmountMinor: sum, Currency: m.Currency}, nil
}

func (m Money) Sub(other Money) (Money, error) {
	negative, err := other.Negate()
	if err != nil {
		return Money{}, err
	}

	return m.Add(negative)
}

func (m Money) Negate() (Money, error) {
	if m.AmountMinor == math.MinInt64 {
		return Money{}, ErrOverflow
	}

	return Money{
		AmountMinor: -m.AmountMinor,
		Currency:    m.Currency,
	}, nil
}

func (m Money) Compare(other Money) (int, error) {
	if err := m.sameCurrency(other); err != nil {
		return 0, err
	}
	switch {
	case m.AmountMinor < other.AmountMinor:
		return -1, nil
	case m.AmountMinor > other.AmountMinor:
		return 1, nil
	default:
		return 0, nil
	}
}

func (m Money) String() string {
	if m.AmountMinor < 0 {
		absolute := uint64(-(m.AmountMinor + 1)) + 1

		return fmt.Sprintf(
			"-%d.%02d",
			absolute/scale,
			absolute%scale,
		)
	}

	return fmt.Sprintf(
		"%d.%02d",
		m.AmountMinor/scale,
		m.AmountMinor%scale,
	)
}
