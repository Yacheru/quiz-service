package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"quiz-service/init/logger"
	"quiz-service/internal/entities"
	"quiz-service/pkg/constants"
	"time"
)

type User struct {
	db     *sqlx.DB
	logger logger.Logging
}

func NewUser(db *sqlx.DB, logger logger.Logging) *User {
	return &User{db: db, logger: logger}
}

func (u *User) Authenticated(ctx context.Context, uuid string) (*entities.User, error) {
	u.logger.InfoF("Authenticated received | %s", uuid)

	var userEntity = new(entities.User)

	query := `
		SELECT id, uuid, login, authorized, authorized_at, quit_at FROM auth WHERE uuid = $1
	`
	if err := u.db.GetContext(ctx, userEntity, query, uuid); err != nil {
		return nil, err
	}

	if !userEntity.Authorized {
		return nil, constants.ErrorUserNotAuthorized
	}

	u.logger.InfoF("Authenticated success | %s", uuid)

	return userEntity, nil
}

func (u *User) Quit(ctx context.Context, uuid string) (int64, error) {
	u.logger.InfoF("Quit received | %s", uuid)

	now := time.Now()

	query := `
		UPDATE auth 
		SET authorized = false, quit_at = $2
		WHERE uuid = $1;
	`
	res, err := u.db.ExecContext(ctx, query, uuid, now)
	if err != nil {
		return 0, err
	}

	u.logger.InfoF("Quit success | %s", uuid)

	return res.RowsAffected()
}
