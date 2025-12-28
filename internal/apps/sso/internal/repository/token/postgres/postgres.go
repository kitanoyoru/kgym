package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	tokenentity "github.com/kitanoyoru/kgym/internal/apps/sso/internal/entity/token"
	tokenrepo "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/token"
	tokenmodel "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/token/models/token"
)

type Repository struct {
	db *pgxpool.Pool
}

var _ tokenrepo.IRepository = (*Repository)(nil)

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

func (r *Repository) GetByTokenHash(ctx context.Context, tokenHash string) (tokenmodel.Token, error) {
	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(tokenmodel.Columns...).
		From(tokenmodel.Table).
		Where(sq.Eq{
			"token_hash": tokenHash,
			"deleted_at": nil,
		}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return tokenmodel.Token{}, err
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return tokenmodel.Token{}, err
	}
	defer rows.Close()

	token, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[tokenmodel.Token])
	if err != nil {
		return tokenmodel.Token{}, err
	}

	return token, nil
}

func (r *Repository) Revoke(ctx context.Context, tokenHash string) error {
	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update(tokenmodel.Table).
		Set("revoked", true).
		Set("updated_at", sq.Expr("now()")).
		Where(sq.Eq{"token_hash": tokenHash, "deleted_at": nil})

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
