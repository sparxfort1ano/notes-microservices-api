package service

import (
	"context"

	"github.com/sparxfort1ano/notes-microservices-api/auth/internal/domain"
)

type Service interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	Read(ctx context.Context, id int) (domain.User, error)
	Update(ctx context.Context, user domain.User) (domain.User, error)
	Delete(ctx context.Context, id int) error
	Authenticate(ctx context.Context, username, password string) (domain.User, error)
	Close() error
}
