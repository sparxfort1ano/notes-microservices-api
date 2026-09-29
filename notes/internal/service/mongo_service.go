package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/sparxfort1ano/notes-microservices-api/notes/internal/domain"
	errs "github.com/sparxfort1ano/notes-microservices-api/notes/internal/errors"
	mongodb "github.com/sparxfort1ano/notes-microservices-api/notes/internal/repository/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const idStr = "_id"
const authorIDStr = "author_id"

type MongoService struct {
	collection *mongo.Collection
}

func NewService(db *mongo.Database) Service {
	collection := db.Collection(mongodb.Collection())

	return &MongoService{
		collection: collection,
	}
}

func (s *MongoService) Create(ctx context.Context, note domain.Note) (domain.Note, error) {
	result, err := s.collection.InsertOne(ctx, note)
	if err != nil {
		return domain.Note{}, fmt.Errorf("%w: %v", errs.ErrNoteCreation, err)
	}

	insertedID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return domain.Note{}, fmt.Errorf("%w: expected bson.ObjectID", errs.ErrUnknown)
	}
	note.ID = insertedID

	return note, nil
}

func (s *MongoService) Get(ctx context.Context, objectID bson.ObjectID, authorID int) (domain.Note, error) {
	filter := bson.M{idStr: objectID, authorIDStr: authorID}

	var note domain.Note
	if err := s.collection.FindOne(ctx, filter).Decode(&note); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Note{}, fmt.Errorf("%w: id=%v", errs.ErrNoteNotFound, objectID)
		}
		return domain.Note{}, fmt.Errorf("%w: %v", errs.ErrDatabaseOperation, err)
	}

	return note, nil
}

func (s *MongoService) GetAll(ctx context.Context, authorID int) ([]domain.Note, error) {
	filter := bson.M{authorIDStr: authorID}

	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrDatabaseOperation, err)
	}
	defer cursor.Close(ctx)

	notes := make([]domain.Note, 0)
	for cursor.Next(ctx) {
		var note domain.Note
		if err := cursor.Decode(&note); err != nil {
			return nil, fmt.Errorf("%w: %v", errs.ErrDecodeNote, err)
		}
		notes = append(notes, note)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrIterationNotes, err)
	}

	return notes, nil
}

func (s *MongoService) Update(ctx context.Context, note domain.Note) (domain.Note, error) {
	filter := bson.M{idStr: note.ID, authorIDStr: note.AuthorID}

	setMap := bson.M{}
	if note.Name != "" {
		setMap["name"] = note.Name
	}
	if note.Content != nil {
		setMap["content"] = note.Content
	}

	if len(setMap) == 0 {
		return domain.Note{}, fmt.Errorf("%w: no fields to update", errs.ErrInvalidData)
	}

	update := bson.M{
		"$set": setMap,
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedNote domain.Note
	if err := s.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedNote); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Note{}, fmt.Errorf("%w: id=%v", errs.ErrNoteNotFound, note.ID)
		}
		return domain.Note{}, fmt.Errorf("%w: %v", errs.ErrNoteUpdate, err)
	}

	return updatedNote, nil
}

func (s *MongoService) Delete(ctx context.Context, objectID bson.ObjectID, authorID int) error {
	filter := bson.M{idStr: objectID, authorIDStr: authorID}

	result, err := s.collection.DeleteOne((ctx), filter)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrDatabaseOperation, err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("%w: id=%v", errs.ErrNoteNotFound, objectID)
	}

	return nil
}

func (s *MongoService) Close() error {
	return mongodb.CloseDB(s.collection.Database().Client())
}
