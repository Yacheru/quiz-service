package service

import (
	"context"
	"database/sql"
	"errors"
	"quiz-service/init/logger"
	"quiz-service/internal/entities"
	"quiz-service/internal/repository"
	"quiz-service/pkg/constants"
	"strings"
)

type Variant struct {
	repo repository.VariantRepository

	log logger.Logging
}

func NewVariant(repo repository.VariantRepository, log logger.Logging) *Variant {
	return &Variant{repo: repo, log: log}
}

func (v *Variant) VariantAdd(ctx context.Context, name string) error {
	if err := v.repo.VariantAdd(ctx, name); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return constants.ErrorVariantAlreadyExists
		}
		if strings.Contains(err.Error(), "value too long for type character") {
			return constants.ErrorVariantTooLong
		}
		v.log.ErrorF("VariantAdd failed: %v", err)
		return err
	}

	return nil
}

func (v *Variant) VariantRemove(ctx context.Context, name string) error {
	num, err := v.repo.VariantRemove(ctx, name)
	if err != nil {
		v.log.ErrorF("VariantRemove failed: %v", err)
		return err
	}
	if num == 0 {
		return constants.ErrorVariantNotFound
	}

	return nil
}

func (v *Variant) VariantList(ctx context.Context) ([]*entities.Variant, error) {
	variants, err := v.repo.VariantList(ctx)
	if err != nil {
		v.log.ErrorF("VariantList failed: %v", err)
		return nil, err
	}

	if len(variants) == 0 {
		return nil, constants.ErrorNoVariantsYet
	}

	return variants, nil
}

func (v *Variant) VariantGet(ctx context.Context, variantName string) (*entities.Variant, error) {
	variant, err := v.repo.VariantGet(ctx, variantName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrorVariantNotFound
		}
		v.log.ErrorF("VariantGet failed: %v", err)
		return nil, err
	}
	return variant, nil
}

func (v *Variant) VariantStart(ctx context.Context, variantId, userId int) error {
	if err := v.repo.VariantStart(ctx, variantId, userId); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil
		}
		v.log.ErrorF("VariantStart failed: %v", err)
		return err
	}
	return nil
}

func (v *Variant) VariantResults(ctx context.Context, variantId, userId int) (*entities.Testing, error) {
	testing, err := v.repo.VariantResults(ctx, variantId, userId)
	if err != nil {
		v.log.ErrorF("VariantResults failed: %v", err)
		return nil, err
	}
	return testing, nil
}
