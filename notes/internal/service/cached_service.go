package service

import (
	"context"
	"errors"

	"github.com/sparxfort1ano/notes-microservices-api/notes/internal/domain"
	"github.com/sparxfort1ano/notes-microservices-api/notes/internal/repository/redis"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CachedService struct {
	cache *redis.CachedDatabase
	db    Service
}

func NewCachedService(cache *redis.CachedDatabase, db Service) Service {
	return &CachedService{
		cache: cache,
		db:    db,
	}
}

func (c *CachedService) Create(ctx context.Context, note domain.Note) (domain.Note, error) {
	createdNote, err := c.db.Create(ctx, note)
	if err != nil {
		return domain.Note{}, err
	}

	c.cache.CacheNote(ctx, createdNote)
	c.cache.InvalidateNotes(ctx, bson.NilObjectID, note.AuthorID)

	return createdNote, nil
}

func (c *CachedService) Get(ctx context.Context, objectID bson.ObjectID, authorID int) (domain.Note, error) {
	if note, hit := c.cache.NoteFromCache(ctx, objectID, authorID); hit {
		return note, nil
	}

	note, err := c.db.Get(ctx, objectID, authorID)
	if err != nil {
		return domain.Note{}, err
	}

	c.cache.CacheNote(ctx, note)

	return note, nil
}

func (c *CachedService) GetAll(ctx context.Context, authorID int) ([]domain.Note, error) {
	if notes, hit := c.cache.NotesFromCache(ctx, authorID); hit {
		return notes, nil
	}

	notes, err := c.db.GetAll(ctx, authorID)
	if err != nil {
		return nil, err
	}

	c.cache.CacheNotes(ctx, notes, authorID)

	return notes, nil
}

func (c *CachedService) Update(ctx context.Context, note domain.Note) (domain.Note, error) {
	updatedNote, err := c.db.Update(ctx, note)
	if err != nil {
		return domain.Note{}, err
	}

	c.cache.CacheNote(ctx, updatedNote)
	c.cache.InvalidateNotes(ctx, bson.NilObjectID, updatedNote.AuthorID)

	return updatedNote, nil
}

func (c *CachedService) Delete(ctx context.Context, objectID bson.ObjectID, authorID int) error {
	if err := c.db.Delete(ctx, objectID, authorID); err != nil {
		return err
	}

	c.cache.InvalidateNotes(ctx, objectID, authorID)

	return nil
}

func (c *CachedService) Close() error {
	return errors.Join(c.cache.Close(), c.db.Close())
}
