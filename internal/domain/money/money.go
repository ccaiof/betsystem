package money

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

var (
	ErrInvalidAmount      = errors.New("money: invalid amount")
	ErrInvalidCurrency    = errors.New("money: invalid currency")
	ErrCurrencyMismatch   = errors.New("money: currency mismatch")
	ErrNegativeNotAllowed = errors.New("money: negative amount not allowed")
	ErrOverflow           = errors.New("money: overflow")
)

const scale = 100 // fixed 2 decimal places

type Money struct {
	amountMinor int64
	currency    string
}

func Zero(currency string) (Money, error) {
	c, err := normalizeCurrency(currency)
	if err != nil {
		return Money{}, err
	}
	return Money{amountMinor: 0, currency: c}, nil
}

func FromMinor(amountMinor int64, currency string) (Money, error) {
	c, err := normalizeCurrency(currency)
	if err != nil {
		return Money{}, err
	}
	return Money{amountMinor: amountMinor, currency: c}, nil
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

	return Money{amountMinor: minor, currency: c}, nil
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

func (m Money) AmountMinor() int64 { return m.amountMinor }
func (m Money) Currency() string   { return m.currency }
func (m Money) IsZero() bool       { return m.amountMinor == 0 }
func (m Money) IsNegative() bool   { return m.amountMinor < 0 }

func (m Money) sameCurrency(other Money) error {
	if m.currency != other.currency {
		return fmt.Errorf("%w: %s vs %s", ErrCurrencyMismatch, m.currency, other.currency)
	}
	return nil
}

func (m Money) Add(other Money) (Money, error) {
	if err := m.sameCurrency(other); err != nil {
		return Money{}, err
	}
	sum := m.amountMinor + other.amountMinor
	if (other.amountMinor > 0 && sum < m.amountMinor) || (other.amountMinor < 0 && sum > m.amountMinor) {
		return Money{}, ErrOverflow
	}
	return Money{amountMinor: sum, currency: m.currency}, nil
}

func (m Money) Sub(other Money) (Money, error) {
	return m.Add(other.Negate())
}

func (m Money) Negate() Money {
	return Money{amountMinor: -m.amountMinor, currency: m.currency}
}

func (m Money) Compare(other Money) (int, error) {
	if err := m.sameCurrency(other); err != nil {
		return 0, err
	}
	switch {
	case m.amountMinor < other.amountMinor:
		return -1, nil
	case m.amountMinor > other.amountMinor:
		return 1, nil
	default:
		return 0, nil
	}
}

func (m Money) String() string {
	sign := ""
	v := m.amountMinor
	if v < 0 {
		sign = "-"
		v = -v
	}
	return fmt.Sprintf("%s%d.%02d", sign, v/scale, v%scale)
}
