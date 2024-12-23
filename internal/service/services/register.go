package service

import (
	"context"
	"database/sql"
	"errors"
	"quiz-service/init/logger"
	"quiz-service/pkg/constants"
	"strings"

	"github.com/google/uuid"

	"quiz-service/internal/entities"
	"quiz-service/internal/repository"
	"quiz-service/pkg/hash"
)

type Register struct {
	repo repository.RegisterRepository

	log logger.Logging

	hasher hash.Hasher
}

func NewRegister(repo repository.RegisterRepository, hasher hash.Hasher, log logger.Logging) *Register {
	return &Register{repo: repo, hasher: hasher, log: log}
}

func (r *Register) Register(ctx context.Context, register *entities.Register) (*entities.User, error) {
	register.UUID = uuid.NewString()
	register.Password = r.hasher.Hash(register.Password)

	user, err := r.repo.Register(ctx, register)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, constants.ErrorUserAlreadyExists
		}
		r.log.ErrorF("Register failed: %v", err)
		return nil, err
	}
	return user, nil
}

func (r *Register) Login(ctx context.Context, login *entities.Login) (*entities.User, error) {
	login.Password = r.hasher.Hash(login.Password)

	user, err := r.repo.Login(ctx, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrorUserNotFound
		}
		r.log.ErrorF("Login failed: %v", err)
		return nil, err
	}

	return user, nil
}
