package application

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(
	ctx context.Context,
	req CreateApplicationRequest,
) (*Application, error) {
	app := &Application{
		Name:      req.Name,
		Code:      req.Code,
		URL:       req.URL,
		IconID:    req.IconID,
		IsDefault: req.IsDefault,
		IsActive:  true,
	}

	if err := s.repo.Create(ctx, app); err != nil {
		return nil, err
	}

	return app, nil
}

func (s *Service) FindAll(
	ctx context.Context,
) ([]Application, error) {
	return s.repo.FindAll(ctx)
}

func (s *Service) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*Application, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) Update(
	ctx context.Context,
	id uuid.UUID,
	req UpdateApplicationRequest,
) (*Application, error) {
	if err := s.repo.Update(ctx, id, req); err != nil {
		return nil, err
	}

	return s.repo.FindByID(ctx, id)
}

func (s *Service) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	return s.repo.Delete(ctx, id)
}
