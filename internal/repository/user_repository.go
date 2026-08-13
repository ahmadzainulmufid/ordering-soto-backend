package repository

import (
	"context"
	"errors"
	"fmt"

	"SotoAyam/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailAlreadyInUse = errors.New("email already in use")
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id int64) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	UpdateLastLogin(ctx context.Context, ID int64) error
	FindAll(ctx context.Context) ([]models.User, error)
	Update(ctx context.Context, id int64, fullName, phone, role string) error
	Delete(ctx context.Context, id int64) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {
	const query = `
		SELECT
			id,
			full_name,
			email,
			phone,
			password_hash,
			role,
			is_active,
			last_login_at,
			created_at,
			updated_at
		FROM users
		WHERE LOWER(email) = LOWER($1)
		LIMIT 1
	`

	var user models.User

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("gagal mencari user berdasarkan email: %w", err)
	}

	return &user, nil
}

func (r *userRepository) FindByID(
	ctx context.Context,
	id int64,
) (*models.User, error) {
	const query = `
		SELECT
			id,
			full_name,
			email,
			phone,
			password_hash,
			role,
			is_active,
			last_login_at,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
		LIMIT 1
	`

	var user models.User

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("gagal mencari user berdasarkan id: %w", err)
	}

	return &user, nil
}

func (r *userRepository) Create(
	ctx context.Context,
	user *models.User,
) error {
	const query = `
		INSERT INTO users (
			full_name,
			email,
			phone,
			password_hash,
			role,
			is_active
		)
		VALUES ($1, LOWER($2), $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		user.FullName,
		user.Email,
		user.Phone,
		user.PasswordHash,
		user.Role,
		user.IsActive,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("gagal membuat user: %w", err)
	}

	return nil
}

func (r *userRepository) UpdateLastLogin(
	ctx context.Context,
	id int64,
) error {
	const query = `
		UPDATE users
		SET
			last_login_at = CURRENT_TIMESTAMP,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("gagal memperbarui waktu login: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *userRepository) FindAll(ctx context.Context) ([]models.User, error) {
	const query = `
		SELECT
			id,
			full_name,
			email,
			phone,
			password_hash,
			role,
			is_active,
			last_login_at,
			created_at,
			updated_at
		FROM users
		ORDER BY id ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil daftar users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.FullName,
			&user.Email,
			&user.Phone,
			&user.PasswordHash,
			&user.Role,
			&user.IsActive,
			&user.LastLoginAt,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal scan data user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error pada iterasi rows user: %w", err)
	}

	return users, nil
}

func (r *userRepository) Update(ctx context.Context, id int64, fullName, phone, role string) error {
	const query = `
		UPDATE users 
		SET full_name = $1, phone = $2, role = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND role != 'owner'
	`
	_, err := r.db.Exec(ctx, query, fullName, phone, role, id)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM users WHERE id = $1 AND role != 'owner'`
	_, err := r.db.Exec(ctx, query, id)
	return err
}