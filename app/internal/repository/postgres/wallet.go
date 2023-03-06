package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/mustan989/wallet/model"
	"github.com/mustan989/wallet/repository"
)

func NewWallet(pool Pool) repository.Wallet {
	return &wallet{pool: pool}
}

type wallet struct{ pool Pool }

const (
	walletsTable   = "wallets"
	walletsBuilder = sqlbuilder.PostgreSQL
)

func (w *wallet) CountAll(ctx context.Context, filter *model.WalletFilter) (count uint64, err error) {
	sb := walletsBuilder.NewSelectBuilder().
		Select("COUNT(*)").
		From(walletsTable)

	if filter.NameLike != "" {
		sb.Where(sb.Like("name", fmt.Sprint("%", filter.NameLike, "%")))
	}
	if filter.DescriptionLike != "" {
		sb.Where(sb.Like("description", fmt.Sprint("%", filter.DescriptionLike, "%")))
	}
	if filter.Currency != "" {
		sb.Where(sb.Equal("currency", filter.Currency))
	}
	if filter.Personal != nil {
		sb.Where(sb.Equal("personal", filter.Personal))
	}

	sql, args := sb.Build()

	err = w.pool.QueryRow(ctx, sql, args...).Scan(&count)

	return
}

func (w *wallet) FindAll(ctx context.Context, filter *model.WalletFilter) (data []*model.Wallet, err error) {
	sb := walletsBuilder.NewSelectBuilder().
		Select("id", "name", "description", "currency", "amount", "personal", "created_at", "updated_at", "deleted_at").
		From(walletsTable)

	if filter.NameLike != "" {
		sb.Where(sb.Like("name", fmt.Sprint("%", filter.NameLike, "%")))
	}
	if filter.DescriptionLike != "" {
		sb.Where(sb.Like("description", fmt.Sprint("%", filter.DescriptionLike, "%")))
	}
	if filter.Currency != "" {
		sb.Where(sb.Equal("currency", filter.Currency))
	}
	if filter.Personal != nil {
		sb.Where(sb.Equal("personal", filter.Personal))
	}

	if filter.Limit != 0 {
		sb.Limit(int(filter.Limit))
	}
	if filter.Offset != 0 {
		sb.Offset(int(filter.Offset))
	}

	sql, args := sb.OrderBy("id DESC").Build()

	rows, err := w.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data = []*model.Wallet{}
	for rows.Next() {
		elem := &model.Wallet{}

		if err = rows.Scan(
			&elem.ID, &elem.Name, &elem.Description, &elem.Currency, &elem.Amount, &elem.Personal, &elem.CreatedAt, &elem.UpdatedAt, &elem.DeletedAt,
		); err != nil {
			return nil, err
		}

		data = append(data, elem)
	}

	return
}

func (w *wallet) FindByID(ctx context.Context, id uint64) (data *model.Wallet, err error) {
	sb := walletsBuilder.NewSelectBuilder().
		Select("id", "name", "description", "currency", "amount", "personal", "created_at", "updated_at", "deleted_at").
		From(walletsTable)
	sb.Where(sb.E("id", id)).Limit(1)

	sql, args := sb.Build()

	data = &model.Wallet{}
	err = w.pool.QueryRow(ctx, sql, args...).Scan(
		&data.ID, &data.Name, &data.Description, &data.Currency, &data.Amount, &data.Personal, &data.CreatedAt, &data.UpdatedAt, &data.DeletedAt,
	)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrWalletNotFound
	}
	if err != nil {
		return nil, err
	}

	return
}

func (w *wallet) Create(ctx context.Context, data *model.Wallet) error {
	ib := walletsBuilder.NewInsertBuilder().
		InsertInto(walletsTable).
		Cols("name", "description", "currency", "amount", "personal").
		Values(data.Name, data.Description, data.Currency, data.Amount, data.Personal)

	sql, args := sqlbuilder.Build(
		`$? RETURNING "id", "name", "description", "currency", "amount", "personal", "created_at", "updated_at", "deleted_at"`, ib,
	).Build()

	if err := w.pool.QueryRow(ctx, sql, args...).Scan(
		&data.ID, &data.Name, &data.Description, &data.Currency, &data.Amount, &data.Personal, &data.CreatedAt, &data.UpdatedAt, &data.DeletedAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return repository.ErrWalletConflict
		}

		return err
	}

	return nil
}

func (w *wallet) Update(ctx context.Context, data *model.Wallet) error {
	ub := walletsBuilder.NewUpdateBuilder().
		Update(walletsTable)
	ub.Set(
		ub.Assign("name", data.Name),
		ub.Assign("description", data.Description),
		ub.Assign("currency", data.Currency),
		ub.Assign("amount", data.Amount),
		ub.Assign("personal", data.Personal),
		"updated_at = default",
		ub.Assign("deleted_at", data.DeletedAt),
	).Where(ub.E("id", data.ID))

	sql, args := sqlbuilder.Build(
		`$? RETURNING "id", "name", "description", "currency", "amount", "personal", "created_at", "updated_at", "deleted_at"`, ub,
	).Build()

	if err := w.pool.QueryRow(ctx, sql, args...).Scan(
		&data.ID, &data.Name, &data.Description, &data.Currency, &data.Amount, &data.Personal, &data.CreatedAt, &data.UpdatedAt, &data.DeletedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrWalletNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return repository.ErrWalletConflict
		}

		return err
	}

	return nil
}

func (w *wallet) DeleteByID(ctx context.Context, id uint64) (deleted *model.Wallet, err error) {
	db := walletsBuilder.NewDeleteBuilder().
		DeleteFrom(walletsTable)
	db.Where(db.E("id", id))

	sql, args := sqlbuilder.Build(
		`$? RETURNING "id", "name", "description", "currency", "amount", "personal", "created_at", "updated_at", "deleted_at"`, db,
	).Build()

	deleted = &model.Wallet{}
	err = w.pool.QueryRow(ctx, sql, args...).Scan(
		&deleted.ID, &deleted.Name, &deleted.Description, &deleted.Currency, &deleted.Amount, &deleted.Personal, &deleted.CreatedAt, &deleted.UpdatedAt, &deleted.DeletedAt,
	)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrWalletNotFound
	}
	if err != nil {
		return nil, err
	}

	return
}
