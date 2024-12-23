package service

import (
	"context"
	"database/sql"
	"errors"
	"quiz-service/init/logger"
	"quiz-service/internal/entities"
	"quiz-service/internal/repository"
	"quiz-service/pkg/constants"
)

type User struct {
	repo repository.UserRepository

	log logger.Logging
}

func NewUser(repo repository.UserRepository, log logger.Logging) *User {
	return &User{repo: repo, log: log}
}

func (u *User) Quit(ctx context.Context, uuid string) error {
	rowsAffected, err := u.repo.Quit(ctx, uuid)
	if err != nil {
		u.log.ErrorF("Quit failed: %v", err)
		return err
	}
	if rowsAffected == 0 {
		return constants.ErrorUserNotFound
	}

	return nil
}

func (u *User) Authenticated(ctx context.Context, uuid string) (*entities.User, error) {
	user, err := u.repo.Authenticated(ctx, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrorUserNotFound
		}
		u.log.ErrorF("Authenticated failed: %v", err)
		return nil, err
	}
	return user, nil
}
