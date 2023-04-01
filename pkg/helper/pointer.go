package helper

import "time"

func Boolp(b bool) *bool                   { return &b }
func Unit8(u uint8) *uint8                 { return &u }
func Uint16(u uint16) *uint16              { return &u }
func Uint32(u uint32) *uint32              { return &u }
func Uint64(u uint64) *uint64              { return &u }
func Int8(i int8) *int8                    { return &i }
func Int16(i int16) *int16                 { return &i }
func Int32(i int32) *int32                 { return &i }
func Int64(i int64) *int64                 { return &i }
func Float32p(f float32) *float32          { return &f }
func Float64p(f float64) *float64          { return &f }
func Complex64p(c complex64) *complex64    { return &c }
func Complex128p(c complex128) *complex128 { return &c }
func Stringp(s string) *string             { return &s }
func Intp(i int) *int                      { return &i }
func Uint(u uint) *uint                    { return &u }
func Timep(t time.Time) *time.Time         { return &t }
