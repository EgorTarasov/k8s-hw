package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/EgorTarasov/shopx/internal/domain"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/larek-tech/storage/postgres"
)

const ordersTable = "orders"

var orderColumns = []string{"id", "user_id", "status", "items", "created_at", "updated_at"}

type OrderRepo struct {
	db *postgres.DB
}

func NewOrderRepo(db *postgres.DB) *OrderRepo { return &OrderRepo{db: db} }

func (r *OrderRepo) Create(ctx context.Context, o domain.Order) error {
	itemsJSON, err := json.Marshal(o.Items)
	if err != nil {
		return fmt.Errorf("marshal items: %w", err)
	}
	sql, args, err := qb.Insert(ordersTable).
		Columns(orderColumns...).
		Values(o.ID.String(), o.UserID.String(), string(o.Status), itemsJSON, o.CreatedAt, o.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("build insert order: %w", err)
	}

	if _, err := r.db.GetPool().Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("insert order: %w", err)
	}
	return nil
}

func (r *OrderRepo) GetByID(ctx context.Context, id domain.OrderID) (domain.Order, error) {
	sql, args, err := qb.Select(orderColumns...).
		From(ordersTable).
		Where(sq.Eq{"id": id.String()}).
		Limit(1).
		ToSql()
	if err != nil {
		return domain.Order{}, fmt.Errorf("build select order: %w", err)
	}

	var (
		o         domain.Order
		idStr     string
		userIDStr string
		status    string
		itemsJSON []byte
	)
	if err := r.db.GetPool().QueryRow(ctx, sql, args...).
		Scan(&idStr, &userIDStr, &status, &itemsJSON, &o.CreatedAt, &o.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Order{}, domain.ErrOrderNotFound
		}
		return domain.Order{}, fmt.Errorf("scan order: %w", err)
	}
	o.ID = domain.OrderID(idStr)
	o.UserID = domain.UserID(userIDStr)
	o.Status = domain.OrderStatus(status)
	if err := json.Unmarshal(itemsJSON, &o.Items); err != nil {
		return domain.Order{}, fmt.Errorf("unmarshal items: %w", err)
	}
	return o, nil
}

func (r *OrderRepo) UpdateStatus(ctx context.Context, id domain.OrderID, status domain.OrderStatus) error {
	sql, args, err := qb.Update(ordersTable).
		Set("status", string(status)).
		Set("updated_at", time.Now().UTC()).
		Where(sq.Eq{"id": id.String()}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build update order status: %w", err)
	}

	tag, err := r.db.GetPool().Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("update order status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrOrderNotFound
	}
	return nil
}
