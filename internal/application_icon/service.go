package application_icon

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
	name string,
	fileData []byte,
	mimeType string,
) (*ApplicationIcon, error) {
	icon := &ApplicationIcon{
		Name:     name,
		FileData: fileData,
		MimeType: mimeType,
	}

	if err := s.repo.Create(ctx, icon); err != nil {
		return nil, err
	}

	return icon, nil
}

func (s *Service) FindAll(
	ctx context.Context,
) ([]ApplicationIcon, error) {
	return s.repo.FindAll(ctx)
}

func (s *Service) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*ApplicationIcon, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) GetFile(
	ctx context.Context,
	id uuid.UUID,
) (*ApplicationIcon, error) {
	return s.repo.FindFileByID(ctx, id)
}

func (s *Service) Update(
	ctx context.Context,
	id uuid.UUID,
	name string,
	fileData []byte,
	mimeType string,
) (*ApplicationIcon, error) {
	icon := &ApplicationIcon{
		Name:     name,
		FileData: fileData,
		MimeType: mimeType,
	}

	if err := s.repo.Update(ctx, id, icon); err != nil {
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
