package model_test

import (
	"encoding"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mustan989/wallet/model"
)

func TestDecimal_String(t *testing.T) {
	subtests := []struct {
		name   string
		input  model.Decimal
		expect string
	}{
		{"Zero", 0, `0.00`},
		{"Precision only", 99, `0.99`},
		{"Exponent only", 9900, `99.00`},
		{"Positive", 9999, `99.99`},
		{"Max value", 9223372036854775799, `92233720368547757.99`},
		{"Negative precision only", -99, `-0.99`},
		{"Negative exponent only", -9900, `-99.00`},
		{"Negative", -9999, `-99.99`},
		{"Min value", -9223372036854775799, `-92233720368547757.99`},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			require.Equal(t, subtest.expect, fmt.Stringer(subtest.input).String())
		})
	}
}

func TestDecimal_MarshalText(t *testing.T) {
	subtests := [...]struct {
		name   string
		input  model.Decimal
		expect []byte
	}{
		{"Zero", 0, []byte(`0.00`)},
		{"Precision only", 99, []byte(`0.99`)},
		{"Exponent only", 9900, []byte(`99.00`)},
		{"Positive", 9999, []byte(`99.99`)},
		{"Max value", 9223372036854775799, []byte(`92233720368547757.99`)},
		{"Negative precision only", -99, []byte(`-0.99`)},
		{"Negative exponent only", -9900, []byte(`-99.00`)},
		{"Negative", -9999, []byte(`-99.99`)},
		{"Min value", -9223372036854775799, []byte(`-92233720368547757.99`)},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			text, err := encoding.TextMarshaler(subtest.input).MarshalText()
			require.NoError(t, err)
			require.Equal(t, subtest.expect, text)
		})
	}
}

func TestDecimal_UnmarshalTextNoError(t *testing.T) {
	subtests := [...]struct {
		name   string
		input  []byte
		expect model.Decimal
	}{
		{"Zero", []byte(`0.00`), 0},
		{"Precision only", []byte(`0.99`), 99},
		{"Exponent only", []byte(`99.00`), 9900},
		{"Integer", []byte(`99`), 9900},
		{"Positive", []byte(`99.99`), 9999},
		{"Max value", []byte(`92233720368547757.99`), 9223372036854775799},
		{"Negative precision only", []byte(`-0.99`), -99},
		{"Negative exponent only", []byte(`-99.00`), -9900},
		{"Negative", []byte(`-99.99`), -9999},
		{"Min value", []byte(`-92233720368547757.99`), -9223372036854775799},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			var val model.Decimal
			err := encoding.TextUnmarshaler(&val).UnmarshalText(subtest.input)
			require.NoError(t, err)
			require.Equal(t, subtest.expect, val)
		})
	}
}

func TestDecimal_UnmarshalTextError(t *testing.T) {
	subtests := [...]struct {
		name  string
		input []byte
	}{
		{"Three dots", []byte(`99.99.99`)},
		{"String precision", []byte(`99.zero`)},
		{"String fraction", []byte(`zero.99`)},
		{"Fraction over limit", []byte(`99.100`)},
		{"Fraction below limit", []byte(`99.-1`)},
		{"Quoted string", []byte(`"99.99"`)},
		{"Over top limit", []byte(`92233720368547758.00`)},
		{"Below bottom limit", []byte(`-92233720368547758.00`)},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			var val model.Decimal
			require.Error(t, encoding.TextUnmarshaler(&val).UnmarshalText(subtest.input))
		})
	}
}

func TestDecimal_MarshalJSON(t *testing.T) {
	subtests := [...]struct {
		name   string
		input  model.Decimal
		expect []byte
	}{
		{"Zero", 0, []byte(`0.00`)},
		{"Precision only", 99, []byte(`0.99`)},
		{"Exponent only", 9900, []byte(`99.00`)},
		{"Positive", 9999, []byte(`99.99`)},
		{"Max value", 9223372036854775799, []byte(`92233720368547757.99`)},
		{"Negative precision only", -99, []byte(`-0.99`)},
		{"Negative exponent only", -9900, []byte(`-99.00`)},
		{"Negative", -9999, []byte(`-99.99`)},
		{"Min value", -9223372036854775799, []byte(`-92233720368547757.99`)},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			text, err := json.Marshal(subtest.input)
			require.NoError(t, err)
			require.Equal(t, subtest.expect, text)
		})
	}
}

func TestDecimal_UnmarshalJSONNoError(t *testing.T) {
	subtests := [...]struct {
		name   string
		input  []byte
		expect model.Decimal
	}{
		{"Zero", []byte(`0.00`), 0},
		{"Precision only", []byte(`0.99`), 99},
		{"Exponent only", []byte(`99.00`), 9900},
		{"Integer", []byte(`99`), 9900},
		{"Positive", []byte(`99.99`), 9999},
		{"Max value", []byte(`92233720368547757.99`), 9223372036854775799},
		{"Negative precision only", []byte(`-0.99`), -99},
		{"Negative exponent only", []byte(`-99.00`), -9900},
		{"Negative", []byte(`-99.99`), -9999},
		{"Min value", []byte(`-92233720368547757.99`), -9223372036854775799},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			var val model.Decimal
			err := json.Unmarshal(subtest.input, &val)
			require.NoError(t, err)
			require.Equal(t, subtest.expect, val)
		})
	}
}

func TestDecimal_UnmarshalJSONError(t *testing.T) {
	subtests := [...]struct {
		name  string
		input []byte
	}{
		{"Three dots", []byte(`99.99.99`)},
		{"String precision", []byte(`99.zero`)},
		{"String fraction", []byte(`zero.99`)},
		{"Fraction over limit", []byte(`99.100`)},
		{"Fraction below limit", []byte(`99.-1`)},
		{"Quoted string", []byte(`"99.99"`)},
		{"Over top limit", []byte(`92233720368547758.00`)},
		{"Below bottom limit", []byte(`-92233720368547758.00`)},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			var val model.Decimal
			require.Error(t, json.Unmarshal(subtest.input, &val))
		})
	}
}
