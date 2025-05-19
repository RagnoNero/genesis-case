package sql

import (
	"database/sql"
	"fmt"
	"time"
	"weather-subscription/sql/dto"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) CreateSubscription(sub *dto.SubscriptionDto) error {
	query := `INSERT INTO subscriptions (email, city, frequency, confirmed, token) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query, sub.Email, sub.City, sub.Frequency, false, sub.Token)
	return err
}

func (r *SubscriptionRepository) IsSubscribed(email string) bool {
	query := `SELECT COUNT(*) FROM subscriptions WHERE email = $1`
	var count int
	err := r.db.QueryRow(query, email).Scan(&count)

	if err != nil {
		return false
	}

	return count > 0
}

func (r *SubscriptionRepository) GetByToken(token string) (*dto.SubscriptionDto, error) {
	query := `SELECT email, city, frequency FROM subscriptions WHERE token = $1`
	row := r.db.QueryRow(query, token)

	var sub dto.SubscriptionDto
	err := row.Scan(&sub.Email, &sub.City, &sub.Frequency)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *SubscriptionRepository) ConfirmToken(token string) error {
	now := time.Now().UTC()
	query := `UPDATE subscriptions 
	          SET confirmed = true, 
	              subscribed_at = $1, 
	              last_sent_at = $2 
	          WHERE token = $3`
	_, err := r.db.Exec(query, now, now, token)
	return err
}

func (r *SubscriptionRepository) GetConfirmedSubscriptions() ([]dto.SubscriptionDto, error) {
	query := `
		SELECT 
			email, city, frequency, token, confirmed, subscribed_at, last_sent_at
		FROM subscriptions
		WHERE confirmed = true
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []dto.SubscriptionDto

	for rows.Next() {
		var sub dto.SubscriptionDto

		err := rows.Scan(
			&sub.Email,
			&sub.City,
			&sub.Frequency,
			&sub.Token,
			&sub.Confirmed,
			&sub.SubscribedAt,
			&sub.LastSentAt,
		)

		if err != nil {
			return nil, fmt.Errorf("scan subscriptions: %w", err)
		}

		subscriptions = append(subscriptions, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate subscriptions: %w", err)
	}

	return subscriptions, nil
}

func (r *SubscriptionRepository) UpdateLastSent(email string) error {
	now := time.Now().UTC()
	query := `UPDATE subscriptions SET last_sent_at = $1 WHERE email = $2`
	_, err := r.db.Exec(query, now, email)
	return err
}

func (r *SubscriptionRepository) Unsubscribe(token string) error {
	query := `DELETE FROM subscriptions WHERE token = $1`
	_, err := r.db.Exec(query, token)
	return err
}
