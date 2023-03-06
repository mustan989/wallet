package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Decimal int64

func (d Decimal) exponent() int64 { return int64(d / 100) }
func (d Decimal) fraction() int64 { return int64(d % 100) }

func (d Decimal) String() string {
	frac := d.fraction()
	if d < 0 {
		frac = -frac
	}
	if -100 < d && d < 0 {
		return fmt.Sprintf("-%d.%02d", d.exponent(), frac)
	}
	return fmt.Sprintf("%d.%02d", d.exponent(), frac)
}

func (d Decimal) MarshalText() ([]byte, error) { return []byte(d.String()), nil }

func (d *Decimal) UnmarshalText(data []byte) error {
	parts := strings.Split(string(data), ".")

	if len(parts) > 2 {
		return errors.New("number must be delimited by one point only")
	}

	var frac int64

	exp, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return fmt.Errorf("exponent: %w", err)
	}

	if len(parts) == 2 {
		frac, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return fmt.Errorf("fraction: %w", err)
		}
		if frac < 0 || 99 < frac {
			return fmt.Errorf("fraction must be value between 0 and 99")
		}
	}

	var int64max int64 = 2 << 62 / 100
	if exp > int64max-1 || exp < -(int64max)+1 {
		return fmt.Errorf("exponent must be in range between -92233720368547757 and 92233720368547757")
	}

	if exp < 0 || exp == 0 && []rune(parts[0])[0] == '-' {
		*d = Decimal(exp*100 - frac)
	} else {
		*d = Decimal(exp*100 + frac)
	}

	return nil
}
func (d Decimal) MarshalJSON() ([]byte, error)     { return d.MarshalText() }
func (d *Decimal) UnmarshalJSON(data []byte) error { return d.UnmarshalText(data) }
