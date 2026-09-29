package errors

import "errors"

var (
	ErrNoteNotFound = errors.New("note not found")
	ErrNoteCreation = errors.New("note creating error")
	ErrNoteUpdate   = errors.New("update note error")

	ErrDatabaseConnection = errors.New("ошибка подключения к базе данных")
	ErrDatabaseOperation  = errors.New("database operation error")
	ErrDatabaseClose      = errors.New("ошибка закрытия соединения с базой данных")
	ErrDatabaseNotInit    = errors.New("база данных не инициализирована")
	ErrCacheConnection    = errors.New("ошибка подключения к кэшу")
	ErrCacheClose         = errors.New("ошибка закрытия соединения с кэшем")
	ErrCacheSet           = errors.New("ошибка записи в кэш")
	ErrCacheGet           = errors.New("ошибка чтения из кэша")
	ErrCacheSerialization = errors.New("ошибка сериализации данных для кэша")
	ErrIterationNotes     = errors.New("notes iteration error")
	ErrDecodeNote         = errors.New("note decoding error")

	ErrInvalidData = errors.New("invalid data format")

	ErrUnknown = errors.New("unknown error")
)

const (
	MsgNoteNotFound  = "note not found"
	MsgInvalidNoteID = "invalid note ID"
	MsgNoteCreation  = "creating note error"
	MsgNoteUpdate    = "note update error"
	MsgNoteDeletion  = "note delete error"

	MsgMissingUserID = "user ID not found in token"

	MsgDatabaseConnection = "Ошибка подключения к базе данных"
	MsgDatabaseOperation  = "database operation error"
	MsgDatabaseClose      = "Ошибка закрытия соединения с базой данных"
	MsgDatabaseNotInit    = "База данных не инициализирована"
	MsgIterationNotes     = "Ошибка итерации по заметкам"
	MsgDecodeNote         = "Ошибка декодирования заметки"

	MsgInvalidData = "invalid data format"

	MsgNoteCreated = "note successfully created"
	MsgNoteUpdated = "note successfully updated"
	MsgNoteDeleted = "note successfully deleted"
	MsgNoteFound   = "note successfully found"
	MsgNotesFound  = "notes successfully found"
)
