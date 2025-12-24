package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	usermodel "github.com/kitanoyoru/kgym/internal/apps/user/internal/repository/models/user"
)

type Repository struct {
	db *pgxpool.Pool
}

var _ IRepository = (*Repository)(nil)

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db,
	}
}

func (r *Repository) Create(ctx context.Context, user usermodel.User) error {
	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(usermodel.Table).
		Columns(usermodel.Columns...).
		Values(user.Values()...)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) List(ctx context.Context, filters ...Filter) ([]usermodel.User, error) {
	var dbFilters Filters
	for _, f := range filters {
		f(&dbFilters)
	}

	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(usermodel.Columns...).
		From(usermodel.Table)

	whereClause := dbFilters.SQL()
	whereClause["deleted_at"] = nil

	query = query.Where(whereClause)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[usermodel.User])
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) Delete(ctx context.Context, filters ...Filter) error {
	var dbFilters Filters
	for _, f := range filters {
		f(&dbFilters)
	}

	whereClause := dbFilters.SQL()
	if len(whereClause) == 0 {
		return nil
	}

	whereClause["deleted_at"] = nil

	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update(usermodel.Table).
		Set("deleted_at", sq.Expr("now()")).
		Where(whereClause)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
