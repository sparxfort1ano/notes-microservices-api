package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sparxfort1ano/notes-microservices-api/notes/internal/domain"
	errs "github.com/sparxfort1ano/notes-microservices-api/notes/internal/errors"
	"github.com/sparxfort1ano/notes-microservices-api/notes/internal/jwt"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	messageStr = "message"
	noteStr    = "note"
	errorStr   = "error"
	detailsStr = "details"
)

func (h *HTTPHandler) CreateNote(c *gin.Context) {
	authorID, err := jwt.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			errorStr:   errs.MsgMissingUserID,
			detailsStr: err.Error(),
		})
		return
	}

	var note domain.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			errorStr:   errs.MsgInvalidData,
			detailsStr: err.Error(),
		})
		return
	}

	note.AuthorID = authorID

	createdNote, err := h.service.Create(c.Request.Context(), note)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			errorStr:   errs.MsgNoteCreation,
			detailsStr: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		messageStr: errs.MsgNoteCreated,
		noteStr:    createdNote,
	})
}

func (h *HTTPHandler) GetNote(c *gin.Context) {
	authorID, err := jwt.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			errorStr:   errs.MsgMissingUserID,
			detailsStr: err.Error(),
		})
		return
	}

	objectID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			errorStr: errs.MsgInvalidNoteID,
		})
		return
	}

	note, err := h.service.Get(c.Request.Context(), objectID, authorID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNoteNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				errorStr:   errs.MsgNoteNotFound,
				detailsStr: err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				errorStr:   errs.MsgNoteDeletion,
				detailsStr: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		messageStr: errs.MsgNoteFound,
		noteStr:    note,
	})
}

func (h *HTTPHandler) GetAllNotes(c *gin.Context) {
	authorID, err := jwt.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			errorStr:   errs.MsgMissingUserID,
			detailsStr: err.Error(),
		})
		return
	}

	notes, err := h.service.GetAll(c.Request.Context(), authorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			errorStr:   errs.MsgDatabaseOperation,
			detailsStr: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		messageStr: errs.MsgNotesFound,
		"notes":    notes,
		"count":    len(notes),
	})
}

func (h *HTTPHandler) UpdateNote(c *gin.Context) {
	authorID, err := jwt.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			errorStr:   errs.MsgMissingUserID,
			detailsStr: err.Error(),
		})
		return
	}

	objectID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			errorStr: errs.MsgInvalidNoteID,
		})
		return
	}

	var req struct {
		Name    string  `json:"name" binding:"omitempty,min=1"`
		Content *string `json:"content" binding:"omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			errorStr:   errs.MsgInvalidData,
			detailsStr: err.Error(),
		})
		return
	}

	note := domain.Note{
		ID:       objectID,
		AuthorID: authorID,
		Name:     req.Name,
		Content:  req.Content,
	}

	updatedNote, err := h.service.Update(c.Request.Context(), note)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNoteNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				errorStr:   errs.MsgNoteNotFound,
				detailsStr: err.Error(),
			})
		case errors.Is(err, errs.ErrInvalidData):
			c.JSON(http.StatusBadRequest, gin.H{
				errorStr:   errs.MsgInvalidData,
				detailsStr: err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				errorStr:   errs.MsgNoteUpdate,
				detailsStr: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		messageStr: errs.MsgNoteUpdated,
		noteStr:    updatedNote,
	})
}

func (h *HTTPHandler) DeleteNote(c *gin.Context) {
	authorID, err := jwt.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			errorStr:   errs.MsgMissingUserID,
			detailsStr: err.Error(),
		})
		return
	}

	objectID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			errorStr: errs.MsgInvalidNoteID,
		})
		return
	}

	if err := h.service.Delete(c.Request.Context(), objectID, authorID); err != nil {
		switch {
		case errors.Is(err, errs.ErrNoteNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				errorStr:   errs.MsgNoteNotFound,
				detailsStr: err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				errorStr:   errs.MsgNoteDeletion,
				detailsStr: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		messageStr: errs.MsgNoteDeleted,
	})
}
