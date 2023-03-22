package location

import "context"

type repository interface {
	GetLocation(ctx context.Context) (*Location, error)
}

type UseCase struct {
	repository repository
}

func (u *UseCase) GetLocation(ctx context.Context) (*Location, error) {
	return u.repository.GetLocation(ctx)
}

func NewUseCase(repository repository) *UseCase {
	return &UseCase{repository: repository}
}
