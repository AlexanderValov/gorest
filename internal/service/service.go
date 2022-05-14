package service

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"rag/internal"
	"rag/internal/models"
	"rag/internal/server"
)

type Repositories interface {
	StartOfTransaction() *sqlx.Tx
	CreateUser(username, bio string, tx *sqlx.Tx) (*int, error)
	GetUsername(username string) error
	UpdateUser(req server.UpdateUserRequest, tx *sqlx.Tx) (*models.User, error)
	DeleteUser(id int) error
	GetAllUsers() ([]*models.User, error)
	GetUser(id int) (*models.User, error)
}

type UserService struct {
	db       Repositories
	logger   *zap.Logger
	settings *internal.Settings
}

func NewUserService(db Repositories, stt *internal.Settings, l *zap.Logger) (*UserService, error) {
	return &UserService{db: db, settings: stt, logger: l}, nil
}

func (us *UserService) validateUsername(username string) error {
	// check if the username empty
	if username == "" {
		return internal.ErrEmptyUsername
	}
	// check the username field is between 3 and 120 chars
	if len(username) < 2 || len(username) > 40 {
		return internal.ErrInvalidSizeOfUsername
	}
	// check username in db
	if err := us.db.GetUsername(username); err != nil {
		return fmt.Errorf("us.db.GetUsername(): %w", err)
	}
	return nil
}

func (us *UserService) CreateUser(request server.CreateUserRequest) (*int, error) {
	// create a user
	userID := make(chan int, 1)
	defer close(userID)
	tx := us.db.StartOfTransaction()
	errs, _ := errgroup.WithContext(context.Background())
	errs.Go(func() error {
		newUserID, err := us.db.CreateUser(request.Username, request.Bio, tx)
		if err != nil {
			return fmt.Errorf("us.db.CreateUser(): %w", err)
		}
		select {
		case <-userID:
		default:
			userID <- *newUserID
		}
		return nil
	})
	// check that a user with this username is not exist, and that it's valid
	if err := us.validateUsername(request.Username); err != nil {
		return nil, fmt.Errorf("us.validateUsername(): %w", err)
	}
	// check errors
	if err := errs.Wait(); err != nil {
		err := tx.Rollback()
		if err != nil {
			return nil, fmt.Errorf("tx.Rollback(), err: %w", err)
		}
		return nil, fmt.Errorf("errs.Wait(), err: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("tx.Commit() err: %v", err)
	}
	id := <-userID
	go us.logger.Info(fmt.Sprintf("User with userID %d created", id))
	return &id, nil
}

func (us *UserService) UpdateUser(request server.UpdateUserRequest) (*models.User, error) {
	// check that a user with this username is not exist, and that it's valid
	errs, _ := errgroup.WithContext(context.Background())
	errs.Go(func() error {
		if err := us.validateUsername(request.Username); err != nil {
			return fmt.Errorf("us.validateUsername(): %w", err)
		}
		return nil
	})
	// update a user
	tx := us.db.StartOfTransaction()
	updatedUser, err := us.db.UpdateUser(request, tx)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return nil, fmt.Errorf("tx.Rollback(), err: %w", err)
		}
		return nil, fmt.Errorf("us.db.UpdateUser(): %v", err)
	}
	// check errors
	if err := errs.Wait(); err != nil {
		err := tx.Rollback()
		if err != nil {
			return nil, fmt.Errorf("tx.Rollback(), err: %w", err)
		}
		return nil, fmt.Errorf("errs.Wait(), err: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("tx.Commit() err: %v", err)
	}
	go us.logger.Info(fmt.Sprintf("User with userID %d updated", request.ID))
	return updatedUser, nil
}

func (us *UserService) DeleteUser(userID int) error {
	if err := us.db.DeleteUser(userID); err != nil {
		return fmt.Errorf("us.db.DeleteUser(): %w", err)
	}
	go us.logger.Info(fmt.Sprintf("User with userID %d deleted", userID))
	return nil
}

func (us *UserService) GetAllUsers() ([]*models.User, error) {
	users, err := us.db.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("us.db.GetAllUsers(): %w", err)
	}
	return users, nil
}

func (us *UserService) GetUser(id int) (*models.User, error) {
	user, err := us.db.GetUser(id)
	if err != nil {
		return nil, fmt.Errorf("us.db.GetUser(): %w", err)
	}
	return user, nil
}
