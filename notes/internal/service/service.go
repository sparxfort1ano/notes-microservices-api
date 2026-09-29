package service

import (
	"context"

	"github.com/sparxfort1ano/notes-microservices-api/notes/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service interface {
	Create(ctx context.Context, note domain.Note) (domain.Note, error)
	Get(ctx context.Context, objectID bson.ObjectID, authorID int) (domain.Note, error)
	GetAll(ctx context.Context, authorID int) ([]domain.Note, error)
	Update(ctx context.Context, note domain.Note) (domain.Note, error)
	Delete(ctx context.Context, objectID bson.ObjectID, authorID int) error
	Close() error
}
