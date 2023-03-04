package model

import "time"

type Wallet struct {
	ID          uint64     `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Currency    string     `json:"currency"`
	Amount      Decimal    `json:"amount"`
	Personal    bool       `json:"personal"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

func (w Wallet) Equals(wallet Wallet) bool {
	return w.Name == wallet.Name &&
		(w.Description == nil && wallet.Description == nil || w.Description != nil && wallet.Description != nil && *w.Description == *wallet.Description) &&
		w.Currency == wallet.Currency &&
		w.Amount == wallet.Amount &&
		w.Personal == wallet.Personal &&
		(w.DeletedAt == nil && wallet.DeletedAt == nil || w.DeletedAt != nil && wallet.DeletedAt != nil && w.DeletedAt.Equal(*wallet.DeletedAt))
}

type WalletFilter struct {
	Filter
	NameLike        string `query:"name_like"`
	DescriptionLike string `query:"description_like"`
	Currency        string `query:"currency"`
	Personal        *bool  `query:"personal"`
}
