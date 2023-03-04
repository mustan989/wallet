package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Pool interface to wrap *pgxpool.Pool to interface
type Pool interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Ping(ctx context.Context) error
	// Close()
	// Acquire(ctx context.Context) (*pgxpool.Conn, error)
	// AcquireFunc(ctx context.Context, f func(*pgxpool.Conn) error) error
	// AcquireAllIdle(ctx context.Context) []*pgxpool.Conn
	// Reset()
	// Config() *pgxpool.Config
	// Stat() *pgxpool.Stat
	// CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
}
