package service

import (
	"context"
	"database/sql"
	"errors"
	"quiz-service/internal/entities"
	"quiz-service/internal/repository"
	"quiz-service/pkg/constants"
)

type User struct {
	repo repository.UserRepository
}

func NewUser(repo repository.UserRepository) *User {
	return &User{repo: repo}
}

func (u *User) Quit(ctx context.Context, uuid string) error {
	rowsAffected, err := u.repo.Quit(ctx, uuid)
	if err != nil {
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
		return nil, err
	}
	return user, nil
}
