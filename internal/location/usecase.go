package location

import (
	"context"
)

type repository interface {
	GetLocation(ctx context.Context, locChan chan Location, errChan chan error)
}

type UseCase struct {
	repository repository
}

func (u *UseCase) GetLocation(ctx context.Context, locChan chan Location, err chan error) {
	u.repository.GetLocation(ctx, locChan, err)
}

func NewUseCase(repository repository) *UseCase {
	return &UseCase{repository: repository}
}
