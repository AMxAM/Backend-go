package persistence

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/alexander/go-api-hex/internal/application/ports"
	"github.com/alexander/go-api-hex/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// userRow es el modelo de persistencia (con tags de DB).
// Está separado de domain.User para no contaminar el dominio.
type userRow struct {
	ID        uuid.UUID `db:"id"`
	Nombre    string    `db:"nombre"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r userRow) toDomain() *domain.User {
	return &domain.User{
		ID:        r.ID,
		Nombre:    r.Nombre,
		Email:     r.Email,
		Password:  r.Password,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

// userRepository es el adaptador concreto que implementa ports.UserRepository.
type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) ports.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, nombre, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Nombre, user.Email, user.Password, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var row userRow
	query := `SELECT id, nombre, email, password, created_at, updated_at FROM users WHERE id = $1`
	if err := r.db.GetContext(ctx, &row, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return row.toDomain(), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var row userRow
	query := `SELECT id, nombre, email, password, created_at, updated_at FROM users WHERE email = $1`
	if err := r.db.GetContext(ctx, &row, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return row.toDomain(), nil
}

func (r *userRepository) List(ctx context.Context) ([]*domain.User, error) {
	var rows []userRow
	query := `SELECT id, nombre, email, password, created_at, updated_at FROM users ORDER BY created_at DESC`
	if err := r.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, err
	}
	users := make([]*domain.User, 0, len(rows))
	for _, row := range rows {
		users = append(users, row.toDomain())
	}
	return users, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now().UTC()
	query := `
		UPDATE users
		SET nombre = $1, email = $2, password = $3, updated_at = $4
		WHERE id = $5
	`
	res, err := r.db.ExecContext(ctx, query,
		user.Nombre, user.Email, user.Password, user.UpdatedAt, user.ID,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}
