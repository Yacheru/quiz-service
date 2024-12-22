package service

import (
	"context"
	"quiz-service/pkg/hash"

	"quiz-service/internal/entities"
	"quiz-service/internal/repository"
	"quiz-service/internal/service/services"
)

type QuestionsService interface {
	QuestionAdd(ctx context.Context, variantId int, question *entities.Question) error
	QuestionRemove(ctx context.Context, variantId int, question *entities.QuestionRemove) error
	QuestionGet(ctx context.Context, variantId, questionId int) (*entities.Question, error)
	QuestionAccept(ctx context.Context, variantId, userId int, answer string) error
}

type TestingService interface {
}

type UserService interface {
	Quit(ctx context.Context, uuid string) error
	Authenticated(ctx context.Context, uuid string) (*entities.User, error)
}

type RegisterService interface {
	Register(ctx context.Context, register *entities.Register) (*entities.User, error)
	Login(ctx context.Context, login *entities.Login) (*entities.User, error)
}

type VariantService interface {
	VariantAdd(ctx context.Context, name string) error
	VariantRemove(ctx context.Context, name string) error
	VariantList(ctx context.Context) ([]*entities.Variant, error)
	VariantStart(ctx context.Context, variantId, userId int) error
	VariantGet(ctx context.Context, variantName string) (*entities.Variant, error)
	VariantResults(ctx context.Context, variantId, userId int) (*entities.Testing, error)
}

type Service struct {
	QuestionsService
	TestingService
	UserService
	RegisterService
	VariantService
}

func NewService(repo *repository.Repository, hasher hash.Hasher) *Service {
	return &Service{
		QuestionsService: service.NewQuestions(repo.QuestionsRepository, repo.VariantRepository, repo.TestingRepository),
		TestingService:   service.NewTesting(repo.TestingRepository),
		UserService:      service.NewUser(repo.UserRepository),
		RegisterService:  service.NewRegister(repo.RegisterRepository, hasher),
		VariantService:   service.NewVariant(repo.VariantRepository),
	}
}
