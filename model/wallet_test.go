package model_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/mustan989/wallet/model"
)

func TestWallet_Equals(t *testing.T) {
	subtests := [...]struct {
		name        string
		left, right model.Wallet
		expect      bool
	}{
		{
			"Full",
			wallet(1, "name", stringp("desc"), "curr", 101, true, time.Time{}, time.Time{}, timep(time.Time{})),
			wallet(1, "name", stringp("desc"), "curr", 101, true, time.Time{}, time.Time{}, timep(time.Time{})),
			true,
		},
		{
			"ID",
			wallet(2, "name", stringp("desc"), "curr", 101, true, time.Time{}, time.Time{}, timep(time.Time{})),
			wallet(1, "name", stringp("desc"), "curr", 101, true, time.Time{}, time.Time{}, timep(time.Time{})),
			true,
		},
		{
			"Name",
			wallet(1, "nm", stringp("desc"), "curr", 101, true, time.Time{}, time.Time{}, timep(time.Time{})),
			wallet(1, "name", stringp("desc"), "curr", 101, true, time.Time{}, time.Time{}, timep(time.Time{})),
			false,
		},
		{
			"Description",
			wallet(1, "name", stringp("des"), "curr", 101, true, time.Time{}, time.Time{}, timep(time.Time{})),
			wallet(1, "name", stringp("desc"), "curr", 101, true, time.Time{}, time.Time{}, timep(time.Time{})),
			false,
		},
		{
			"Description nil",
			wallet(1, "name", nil, "curr", 101, true, time.Time{}, time.Time{}, timep(time.Time{})),
			wallet(1, "name", stringp("desc"), "curr", 101, true, time.Time{}, time.Time{}, timep(time.Time{})),
			false,
		},
		{
			"Currency",
			wallet(1, "name", stringp("desc"), "cu", 101, true, time.Time{}, time.Time{}, timep(time.Time{})),
			wallet(1, "name", stringp("desc"), "curr", 101, true, time.Time{}, time.Time{}, timep(time.Time{})),
			false,
		},
		{
			"Amount",
			wallet(1, "name", stringp("desc"), "curr", 1011, true, time.Time{}, time.Time{}, timep(time.Time{})),
			wallet(1, "name", stringp("desc"), "curr", 101, true, time.Time{}, time.Time{}, timep(time.Time{})),
			false,
		},
		{
			"Personal",
			wallet(1, "name", stringp("desc"), "curr", 101, false, time.Time{}, time.Time{}, timep(time.Time{})),
			wallet(1, "name", stringp("desc"), "curr", 101, true, time.Time{}, time.Time{}, timep(time.Time{})),
			false,
		},
		{
			"CreatedAt",
			wallet(1, "name", stringp("desc"), "curr", 101, true, time.Now(), time.Time{}, timep(time.Time{})),
			wallet(1, "name", stringp("desc"), "curr", 101, true, time.Time{}, time.Time{}, timep(time.Time{})),
			true,
		},
		{
			"UpdatedAt",
			wallet(1, "name", stringp("desc"), "curr", 101, true, time.Time{}, time.Now(), timep(time.Time{})),
			wallet(1, "name", stringp("desc"), "curr", 101, true, time.Time{}, time.Time{}, timep(time.Time{})),
			true,
		},
		{
			"DeletedAt",
			wallet(1, "name", stringp("desc"), "curr", 101, true, time.Time{}, time.Time{}, timep(time.Now())),
			wallet(1, "name", stringp("desc"), "curr", 101, true, time.Time{}, time.Time{}, timep(time.Time{})),
			false,
		},
		{
			"DeletedAt nil",
			wallet(1, "name", stringp("desc"), "curr", 101, true, time.Time{}, time.Time{}, nil),
			wallet(1, "name", stringp("desc"), "curr", 101, true, time.Time{}, time.Time{}, timep(time.Time{})),
			false,
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			require.Equal(t, subtest.expect, subtest.left.Equals(subtest.right))
		})
	}
}

func stringp(s string) *string     { return &s }
func timep(t time.Time) *time.Time { return &t }

func wallet(id uint64, name string, description *string, currency string, amount model.Decimal, personal bool, createdAt, updatedAt time.Time, deletedAt *time.Time) model.Wallet {
	return model.Wallet{id, name, description, currency, amount, personal, createdAt, updatedAt, deletedAt}
}
