package repository

import (
	"context"

	"github.com/VAGRAMCHIC/vless_reality_agent/internal/domain"
	"github.com/VAGRAMCHIC/vless_reality_agent/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepo struct {
	db  *pgxpool.Pool
	log *utils.Logger
}

func NewUserRepository(db *pgxpool.Pool, log *utils.Logger) UserRepository {
	return &userRepo{db: db, log: log}
}

func (r *userRepo) Create(ctx context.Context, u *domain.User) error {
	u.ID = uuid.New()

	_, err := r.db.Exec(ctx,
		`INSERT INTO users (id, short_id, email, tag, server)
		 VALUES ($1,$2,$3,$4,$5)`,
		u.ID, u.ShortID, u.Email, u.Tag, u.Server,
	)
	if err != nil {
		r.log.Error(ctx, "failed_create_user", map[string]any{
			"error": err.Error(),
		})
	}
	return err
}

func (r *userRepo) DeleteUser(ctx context.Context, id *uuid.UUID) error {
	query := `DELETE FROM users WHERE id=$1`

	err := r.db.QueryRow(ctx, query, id).Scan()
	if err != nil {
		r.log.Error(ctx, "failed_to_delete_user", map[string]any{
			"error": err.Error(),
		})
		return err
	}
	return nil
}

func (r *userRepo) GetAll(ctx context.Context) ([]domain.User, error) {
	query := `SELECT id, short_id, email, tag, server FROM users`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		r.log.Error(ctx, "failed_get_users", map[string]any{
			"error": err.Error(),
		})
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User

		err := rows.Scan(
			&u.ID,
			&u.ShortID,
			&u.Email,
			&u.Tag,
			&u.Server,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return users, nil
}

func (r *userRepo) GetUserByID(ctx context.Context, id *uuid.UUID) (*domain.User, error) {
	query := `SELECT id,short_id,email,tag,server FROM users WHERE id=$1`
	var user domain.User
	err := r.db.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.ShortID, &user.Email, &user.Tag, &user.Server)
	if err != nil {
		r.log.Error(ctx, "failed_get_user", map[string]any{
			"error": err.Error(),
		})
		return nil, err
	}
	return &user, nil

}
