package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	userentity "github.com/kitanoyoru/kgym/internal/apps/user/internal/entity/user"
	userrepo "github.com/kitanoyoru/kgym/internal/apps/user/internal/repository"
	usermodel "github.com/kitanoyoru/kgym/internal/apps/user/internal/repository/models/user"
)

type Repository struct {
	db *pgxpool.Pool
}

var _ userrepo.IRepository = (*Repository)(nil)

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db,
	}
}

func (r *Repository) GetByID(ctx context.Context, id string) (usermodel.User, error) {
	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(usermodel.Columns...).
		From(usermodel.Table).
		Where(sq.Eq{"id": id, "deleted_at": nil}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return usermodel.User{}, err
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return usermodel.User{}, err
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[usermodel.User])
	if err != nil {
		return usermodel.User{}, err
	}

	return user, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (usermodel.User, error) {
	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(usermodel.Columns...).
		From(usermodel.Table).
		Where(sq.Eq{"email": email, "deleted_at": nil}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return usermodel.User{}, err
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return usermodel.User{}, err
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[usermodel.User])
	if err != nil {
		return usermodel.User{}, err
	}

	return user, nil
}

func (r *Repository) Create(ctx context.Context, user userentity.User) error {
	model, err := usermodel.FromEntity(user)
	if err != nil {
		return err
	}

	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(usermodel.Table).
		Columns(usermodel.Columns...).
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

func (r *Repository) DeleteByID(ctx context.Context, id string) error {
	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update(usermodel.Table).
		Set("deleted_at", sq.Expr("now()")).
		Where(sq.Eq{"id": id})

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
