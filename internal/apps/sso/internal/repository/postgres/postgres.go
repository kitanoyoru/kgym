package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	tokenentity "github.com/kitanoyoru/kgym/internal/apps/sso/internal/entity/token"
	tokenmodel "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/models/token"
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

func (r *Repository) Create(ctx context.Context, token tokenentity.Token) error {
	model := tokenmodel.TokenFromEntity(token)

	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(tokenmodel.Table).
		Columns(tokenmodel.Columns...).
		Values(model.Values()...)

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

func (r *Repository) List(ctx context.Context, filters ...Filter) ([]tokenmodel.Token, error) {
	var dbFilters Filters
	for _, f := range filters {
		f(&dbFilters)
	}

	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(tokenmodel.Columns...).
		From(tokenmodel.Table)

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

	tokens, err := pgx.CollectRows(rows, pgx.RowToStructByName[tokenmodel.Token])
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (r *Repository) Update(ctx context.Context, id, token string) error {
	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update(tokenmodel.Table).
		Set("token", token).
		Set("updated_at", sq.Expr("now()")).
		Where(sq.Eq{"id": id, "deleted_at": nil})

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

func (r *Repository) Delete(ctx context.Context, filters ...Filter) error {
	var dbFilters Filters
	for _, f := range filters {
		f(&dbFilters)
	}

	whereClause := dbFilters.SQL()
	whereClause["deleted_at"] = nil

	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update(tokenmodel.Table).
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
