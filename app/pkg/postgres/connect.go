package postgres

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mustan989/wallet/app/config"
)

func Connect(ctx context.Context, database *config.Database) (*pgxpool.Pool, error) {
	pass, err := base64.StdEncoding.DecodeString(database.Pass)
	if err != nil {
		return nil, fmt.Errorf("decode pass: %w", err)
	}
	return pgxpool.New(ctx, fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s", database.User, pass, database.Host, database.Port, database.Name,
	))
}
