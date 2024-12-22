package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"quiz-service/init/logger"
	"quiz-service/internal/entities"
	"quiz-service/internal/repository/postgres"
)

type QuestionsRepository interface {
	QuestionAdd(ctx context.Context, variantId int, question *entities.Question) error
	QuestionRemove(ctx context.Context, variantId int, question string) (int64, error)
	QuestionGet(ctx context.Context, variantId, questionId int) (*entities.Question, error)
	QuestionCount(ctx context.Context, variantId int) (int, error)
	QuestionAccept(ctx context.Context, answer string, userId, testId, variantId int) error
}

type RegisterRepository interface {
	Register(ctx context.Context, register *entities.Register) (*entities.User, error)
	Login(ctx context.Context, login *entities.Login) (*entities.User, error)
}

type TestingRepository interface {
	TestGet(ctx context.Context, userId, variantId int) (*entities.Testing, error)
}

type UserRepository interface {
	Quit(ctx context.Context, uuid string) (int64, error)
	Authenticated(ctx context.Context, uuid string) (*entities.User, error)
}

type VariantRepository interface {
	VariantAdd(ctx context.Context, name string) error
	VariantRemove(ctx context.Context, name string) (int64, error)
	VariantList(ctx context.Context) ([]*entities.Variant, error)
	VariantGet(ctx context.Context, name string) (*entities.Variant, error)
	VariantStart(ctx context.Context, variantId, userId int) error
	VariantResults(ctx context.Context, variantId, userId int) (*entities.Testing, error)
}

type Repository struct {
	QuestionsRepository
	RegisterRepository
	TestingRepository
	UserRepository
	VariantRepository
}

func NewRepository(db *sqlx.DB, logger *logger.Logger) *Repository {
	return &Repository{
		QuestionsRepository: postgres.NewQuestions(db, logger),
		RegisterRepository:  postgres.NewRegister(db, logger),
		TestingRepository:   postgres.NewTesting(db, logger),
		UserRepository:      postgres.NewUser(db, logger),
		VariantRepository:   postgres.NewVariant(db, logger),
	}
}
