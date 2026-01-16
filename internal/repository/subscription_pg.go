package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/MichaelShi-san/subscription-test/internal/model"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(ctx context.Context, s *model.Subscription) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`,
		s.Service,
		s.Price,
		s.UserID,
		s.StartDate,
		s.EndDate,
	).Scan(&s.ID, &s.CreatedAt)
}

func (r *SubscriptionRepository) List(ctx context.Context) ([]model.Subscription, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at
		FROM subscriptions
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Subscription
	for rows.Next() {
		var s model.Subscription
		if err := rows.Scan(
			&s.ID,
			&s.Service,
			&s.Price,
			&s.UserID,
			&s.StartDate,
			&s.EndDate,
			&s.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, s)
	}

	return result, nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM subscriptions WHERE id = $1`, id)
	return err
}

func (r *SubscriptionRepository) Update(ctx context.Context, s *model.Subscription) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE subscriptions
		SET service_name = $1, price = $2, end_date = $3
		WHERE id = $4
	`,
		s.Service,
		s.Price,
		s.EndDate,
		s.ID,
	)
	return err
}

func (r *SubscriptionRepository) TotalCost(
	ctx context.Context,
	userID string,
	service string,
	from, to time.Time,
) (int, error) {
	var sum int
	err := r.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(price), 0)
		FROM subscriptions
		WHERE user_id = $1
		  AND service_name = $2
		  AND start_date BETWEEN $3 AND $4
	`,
		userID, service, from, to,
	).Scan(&sum)

	return sum, err
}
