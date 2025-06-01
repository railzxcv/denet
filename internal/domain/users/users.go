package users

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type LeaderboardEntry struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Points   int    `json:"points"`
}

type StatusInfo struct {
	ID           uuid.UUID   `db:"id" json:"id"`
	Username     string      `db:"username" json:"username"`
	Email        string      `db:"email" json:"email"`
	Points       int         `db:"points" json:"points"`
	ReferralCode pgtype.Text `db:"referral_code" json:"referral_code"`
	ReferredBy   pgtype.Text `db:"referred_by" json:"referred_by"`
	CreatedAt    time.Time   `db:"created_at" json:"created_at"`
}
