package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/EgorTarasov/shopx/internal/domain"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/larek-tech/storage/postgres"
)

const usersTable = "users"

var userColumns = []string{"id", "email", "name", "password_hash", "created_at"}

type UserRepo struct {
	db *postgres.DB
}

func NewUserRepo(db *postgres.DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) Create(ctx context.Context, u domain.User) error {
	sql, args, err := qb.Insert(usersTable).
		Columns(userColumns...).
		Values(u.ID.String(), u.Email, u.Name, u.PasswordHash, u.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("build insert user: %w", err)
	}

	if _, err := r.db.GetPool().Exec(ctx, sql, args...); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return domain.ErrEmailTaken
		}
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	return r.getOne(ctx, sq.Eq{"email": email})
}

func (r *UserRepo) GetByID(ctx context.Context, id domain.UserID) (domain.User, error) {
	return r.getOne(ctx, sq.Eq{"id": id.String()})
}

func (r *UserRepo) getOne(ctx context.Context, where sq.Sqlizer) (domain.User, error) {
	sql, args, err := qb.Select(userColumns...).From(usersTable).Where(where).Limit(1).ToSql()
	if err != nil {
		return domain.User{}, fmt.Errorf("build select user: %w", err)
	}
	return scanUser(r.db.GetPool().QueryRow(ctx, sql, args...))
}

func scanUser(row pgx.Row) (domain.User, error) {
	var (
		u  domain.User
		id string
	)
	if err := row.Scan(&id, &u.Email, &u.Name, &u.PasswordHash, &u.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, fmt.Errorf("scan user: %w", err)
	}
	u.ID = domain.UserID(id)
	return u, nil
}
