package redis

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/sparxfort1ano/notes-microservices-api/notes/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type NoteModel domain.Note

func modelToDomain(model NoteModel) domain.Note {
	return domain.Note(model)
}

func domainToModel(domain domain.Note) NoteModel {
	return NoteModel(domain)
}

func (m *NoteModel) serialize() ([]byte, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("serialize note: %w", err)
	}

	return bytes, nil
}

func (m *NoteModel) deserialize(bytes []byte) error {
	if err := json.Unmarshal(bytes, m); err != nil {
		return fmt.Errorf("deserialize note: %w", err)
	}

	return nil
}

func getNoteKey(id, authorID string) string {
	return fmt.Sprintf("note:%s:%s", id, authorID)
}

type NoteListModel []NoteModel

func modelsToDomains(models NoteListModel) []domain.Note {
	notes := make([]domain.Note, len(models))

	for i, model := range models {
		notes[i] = modelToDomain(model)
	}

	return notes
}

func domainsToModels(domain []domain.Note) NoteListModel {
	models := make(NoteListModel, len(domain))

	for i, note := range domain {
		models[i] = domainToModel(note)
	}

	return models
}

func (m *NoteListModel) serialize() ([]byte, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("serialize note list: %w", err)
	}

	return bytes, nil
}

func (m *NoteListModel) deserialize(bytes []byte) error {
	if err := json.Unmarshal(bytes, m); err != nil {
		return fmt.Errorf("deserialize note list: %w", err)
	}

	return nil
}

func getNoteListKey(authorID string) string {
	return fmt.Sprintf("notes:author:%s", authorID)
}

func objectIDToStr(id bson.ObjectID) string {
	return id.Hex()
}

func authorIDToStr(authorID int) string {
	return strconv.Itoa(authorID)
}
