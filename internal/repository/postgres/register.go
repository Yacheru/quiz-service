package postgres

import (
	"context"
	"quiz-service/init/logger"
	"time"

	"github.com/jmoiron/sqlx"

	"quiz-service/internal/entities"
)

type Register struct {
	db     *sqlx.DB
	logger logger.Logging
}

func NewRegister(db *sqlx.DB, logger logger.Logging) *Register {
	return &Register{db: db, logger: logger}
}

func (r Register) Register(ctx context.Context, register *entities.Register) (*entities.User, error) {
	r.logger.InfoF("Register received | %+v", register)

	var userEntity = new(entities.User)

	query := `
		INSERT INTO auth (uuid, login, password, authorized) 
		VALUES ($1, $2, $3, $4)
		RETURNING id, uuid, login, authorized, authorized_at, quit_at
	`
	if err := r.db.GetContext(ctx, userEntity, query, register.UUID, register.Login, register.Password, true); err != nil {
		return nil, err
	}

	r.logger.InfoF("Register success | %+v", register)

	return userEntity, nil
}

func (r Register) Login(ctx context.Context, login *entities.Login) (*entities.User, error) {
	r.logger.InfoF("Login received | %+v", login)

	var userEntity = new(entities.User)

	query := `
		UPDATE auth 
		SET authorized = true, authorized_at = $1
		WHERE login = $2 AND password = $3
		RETURNING uuid, login, authorized, authorized_at, quit_at
	`
	if err := r.db.GetContext(ctx, userEntity, query, time.Now(), login.Login, login.Password); err != nil {
		return nil, err
	}

	r.logger.InfoF("Login success | %+v", login)

	return userEntity, nil
}
