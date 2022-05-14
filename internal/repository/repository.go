package repository

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"rag/internal"
	"rag/internal/models"
	"rag/internal/server"
)

type Repository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewRepository(settings *internal.Settings, logger *zap.Logger) (*Repository, error) {
	// connect to database
	// NOTE: we are using pgx as SQL driver: https://github.com/jackc/pgx
	db, err := sqlx.Connect("pgx", settings.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("repository.NewRepository() sqlx.Connect error: %w", err)
	}
	// ping database to check that connection is OK
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("repository.NewRepository() ping db error: %w", err)
	}

	r := &Repository{
		db:     db,
		logger: logger,
	}
	return r, nil
}

func (r *Repository) StartOfTransaction() *sqlx.Tx {
	tx := r.db.MustBegin()
	return tx
}

func (r *Repository) GetUsername(username string) error {
	query := `
	SELECT username
	FROM api_user
	WHERE username = $1
	LIMIT 1`
	var nameFromDB *string
	if err := r.db.Get(&nameFromDB, query, username); err != nil {
		// if no rows == username isn't exist
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	// if no errors == username exists
	return internal.ErrUsernameExist
}

func (r *Repository) CreateUser(username, bio string, tx *sqlx.Tx) (*int, error) {
	query := `
	INSERT INTO api_user(username, bio, created_at, updated_at)
	VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	RETURNING id`
	var id *int
	err := tx.QueryRow(query, username, bio).Scan(&id)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (r *Repository) UpdateUser(req server.UpdateUserRequest, tx *sqlx.Tx) (*models.User, error) {
	query := `
	UPDATE api_user
    SET username = $2, bio = $3, is_active = $4, updated_at = CURRENT_TIMESTAMP
    WHERE id = $1
    RETURNING *`
	var updatedUser models.User
	err := tx.QueryRow(query, req.ID, req.Username, req.Bio, req.IsActive).Scan(
		&updatedUser.ID, &updatedUser.Username, &updatedUser.Bio,
		&updatedUser.IsActive, &updatedUser.CreatedAt, &updatedUser.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &updatedUser, nil
}

func (r *Repository) DeleteUser(id int) error {
	query := `DELETE FROM api_user WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *Repository) GetAllUsers() ([]*models.User, error) {
	query := `
	SELECT *
	FROM api_user`
	var users []*models.User
	if err := r.db.Select(&users, query); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) GetUser(id int) (*models.User, error) {
	query := `
	SELECT *
	FROM api_user
	WHERE id = $1
	LIMIT 1`
	var user models.User
	if err := r.db.Get(&user, query, id); err != nil {
		// if no rows == username isn't exist
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internal.ErrNoUser
		}
		return nil, err
	}
	return &user, nil
}
