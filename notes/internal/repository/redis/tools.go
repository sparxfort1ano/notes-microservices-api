package redis

import (
	"context"
	"errors"

	goredis "github.com/redis/go-redis/v9"
	"github.com/sparxfort1ano/notes-microservices-api/notes/internal/domain"
	"github.com/sparxfort1ano/notes-microservices-api/notes/internal/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/zap"
)

func (c *CachedDatabase) NoteFromCache(
	ctx context.Context,
	id bson.ObjectID,
	authorID int,
) (domain.Note, bool) {
	log := logger.FromContext(ctx)

	key := getNoteKey(objectIDToStr(id), authorIDToStr(authorID))

	bytes, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if !errors.Is(err, goredis.Nil) {
			log.Error("read note from cache", zap.Error(err))
		}

		return domain.Note{}, false
	}

	var noteModel NoteModel
	if err := noteModel.deserialize(bytes); err != nil {
		log.Error("deserialize cached note", zap.Error(err))
		return domain.Note{}, false
	}

	note := modelToDomain(noteModel)

	return note, true
}

func (c *CachedDatabase) NotesFromCache(
	ctx context.Context,
	authorID int,
) ([]domain.Note, bool) {
	log := logger.FromContext(ctx)

	key := getNoteListKey(authorIDToStr(authorID))

	bytes, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if !errors.Is(err, goredis.Nil) {
			log.Error("read note list", zap.Error(err))
		}

		return nil, false
	}

	var noteModels NoteListModel
	if err := noteModels.deserialize(bytes); err != nil {
		log.Error("deserialize note list", zap.Error(err))
		return nil, false
	}

	notes := modelsToDomains(noteModels)

	return notes, true
}

func (c *CachedDatabase) CacheNote(
	ctx context.Context,
	note domain.Note,
) {
	log := logger.FromContext(ctx)

	key := getNoteKey(objectIDToStr(note.ID), authorIDToStr(note.AuthorID))

	model := domainToModel(note)
	bytes, err := model.serialize()
	if err != nil {
		log.Error("serialize note", zap.Error(err))
		return
	}

	if err = c.client.Set(ctx, key, bytes, c.TTL()).Err(); err != nil {
		log.Error("set note in cache", zap.Error(err))
		return
	}
}

func (c *CachedDatabase) CacheNotes(
	ctx context.Context,
	notes []domain.Note,
	authorID int,
) {
	log := logger.FromContext(ctx)

	key := getNoteListKey(authorIDToStr(authorID))

	models := domainsToModels(notes)
	bytes, err := models.serialize()
	if err != nil {
		log.Error("serialize note list", zap.Error(err))
		return
	}

	if err = c.client.Set(ctx, key, bytes, c.TTL()).Err(); err != nil {
		log.Error("set note list", zap.Error(err))
		return
	}
}

func (c *CachedDatabase) InvalidateNotes(
	ctx context.Context,
	id bson.ObjectID,
	authorID int,
) {
	log := logger.FromContext(ctx)

	invalidateKeys := []string{
		getNoteListKey(authorIDToStr(authorID)),
	}
	if !id.IsZero() {
		invalidateKeys = append(invalidateKeys, getNoteKey(objectIDToStr(id), authorIDToStr(authorID)))
	}

	if err := c.client.Del(ctx, invalidateKeys...).Err(); err != nil {
		log.Error("invalidate notes", zap.Error(err))
		return
	}
}
