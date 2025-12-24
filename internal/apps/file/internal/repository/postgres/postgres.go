package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	filemodel "github.com/kitanoyoru/kgym/internal/apps/file/internal/repository/models/file"
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

func (r *Repository) Create(ctx context.Context, file filemodel.File) error {
	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(filemodel.Table).
		Columns(filemodel.Columns...).
		Values(file.Values()...)

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

func (r *Repository) List(ctx context.Context, filters ...Filter) ([]filemodel.File, error) {
	var dbFilters Filters
	for _, f := range filters {
		f(&dbFilters)
	}

	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(filemodel.Columns...).
		From(filemodel.Table)

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

	files, err := pgx.CollectRows(rows, pgx.RowToStructByName[filemodel.File])
	if err != nil {
		return nil, err
	}

	return files, nil
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
		Update(filemodel.Table).
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
