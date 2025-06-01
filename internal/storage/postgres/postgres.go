package postgres

import (
	"context"
	"denet-app/internal/domain/users"
	errx "denet-app/internal/errx"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(URL string) (*PostgresStorage, error) {
	pool, err := pgxpool.New(context.Background(), URL)
	if err != nil {
		return nil, fmt.Errorf("can't create pgxpool.Pool for URL %s: %w", URL, err)
	}
	if pool.Ping(context.Background()) != nil {
		return nil, fmt.Errorf("can't connect to DB with URL %s: %w", URL, pool.Ping(context.Background()))
	}
	return &PostgresStorage{pool: pool}, nil
}

func (pg *PostgresStorage) GetLeaderboard() ([]users.LeaderboardEntry, error) {
	rows, err := pg.pool.Query(context.Background(),
		`SELECT u.id, u.username, u.email, u.points
		FROM users u
		ORDER BY u.points DESC, u.username ASC
		LIMIT 10;`,
	)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to query leaderboard: %w", err)
	}

	var leaderboard []users.LeaderboardEntry
	for rows.Next() {
		var user users.LeaderboardEntry
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Points); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		leaderboard = append(leaderboard, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return leaderboard, nil
}
func (pg *PostgresStorage) GetUserStatus(id string) (users.StatusInfo, error) {
	var user users.StatusInfo

	err := pg.pool.QueryRow(context.Background(),
		`SELECT u.id, u.username, u.email, u.points, u.referral_code, ref.username AS referred_by, u.created_at
		FROM users u
		LEFT JOIN users ref ON u.referred_by = ref.id
		WHERE u.id = $1;`,
		id,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Points,
		&user.ReferralCode,
		&user.ReferredBy,
		&user.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return user, fmt.Errorf("failed to get user status: %w", err)
	}

	return user, nil
}

func (pg *PostgresStorage) PostTaskComplete(id, taskType string) error {
	var exists bool
	err := pg.pool.QueryRow(context.Background(),
		`SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`, id,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return errx.ErrUserNotFound
	}
	var TaskCompleted bool
	err = pg.pool.QueryRow(context.Background(),
		`SELECT EXISTS (
			SELECT 1 
			FROM completed_tasks ct
			JOIN task_types tt
			ON ct.task_type_id = tt.id
			WHERE ct.user_id = $1 AND tt.task_type = $2
		)`,
		id,
		taskType,
	).Scan(&TaskCompleted)
	if err != nil {
		return fmt.Errorf("failed to check if task is already completed: %w", err)
	}
	if TaskCompleted {
		return errx.ErrNoChange
	}

	var TaskPoints int
	err = pg.pool.QueryRow(context.Background(),
		`SELECT points 
		FROM task_types
		WHERE task_type = $1;`,
		taskType,
	).Scan(&TaskPoints)
	if err != nil {
		return fmt.Errorf("failed to check task points: %w", err)
	}
	if err == pgx.ErrNoRows {
		return errx.ErrTaskNotFound
	}
	_, err = pg.pool.Exec(context.Background(),
		`UPDATE users 
		SET points = points + $2
		WHERE id = $1;`,
		id,
		TaskPoints,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update user points: %w", err)
	}
	_, err = pg.pool.Exec(context.Background(),
		`INSERT INTO completed_tasks (user_id, task_type_id)
		VALUES (
			$1,
			(SELECT id 
			FROM task_types 
			WHERE task_type = $2)
		)`,
		id,
		taskType,
	)

	if err != nil {
		return fmt.Errorf("failed to insert entry to completed_tasks: %w", err)
	}

	return nil

}

func (pg *PostgresStorage) PostReferrer(id, referralCode, referredBy string) error {
	var exists bool
	err := pg.pool.QueryRow(context.Background(),
		`SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`, id,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return errx.ErrUserNotFound
	}
	err = pg.pool.QueryRow(context.Background(),
		`SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`, referredBy,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check referrer existence: %w", err)
	}
	if !exists {
		return errx.ErrReferrerNotFound
	}

	_, err = pg.pool.Exec(context.Background(),
		`UPDATE users 
		SET referral_code = $2, referred_by = $3
		WHERE id = $1;`,
		id,
		referralCode,
		referredBy,
	)
	if err != nil {
		return fmt.Errorf("failed to update referrer: %w", err)
	}
	return nil
}

