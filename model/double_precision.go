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

func (d Decimal) String() string               { return fmt.Sprint(d.exponent(), ".", d.fraction()) }
func (d Decimal) MarshalText() ([]byte, error) { return []byte(d.String()), nil }
func (d Decimal) MarshalJSON() ([]byte, error) { return d.MarshalText() }

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
		if frac < 0 || frac > 99 {
			return fmt.Errorf("fraction must be value between 0 and 99")
		}
	}

	*d = Decimal(exp*100 + frac)

	return nil
}
func (d *Decimal) UnmarshalJSON(data []byte) error { return d.UnmarshalText(data) }
