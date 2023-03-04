package model

type Filter struct {
	Limit  uint64 `json:"limit" query:"limit"`
	Offset uint64 `json:"offset" query:"offset"`
}
