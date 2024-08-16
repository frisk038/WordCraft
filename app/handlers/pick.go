package handlers

import (
	"context"
	"net/http"

	"github.com/frisk038/wordcraft/business/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type businessPick interface {
	GetDailyWord(ctx context.Context) (models.Pick, error)
	CheckWordExists(ctx context.Context, word string) (models.Score, error)
}

func GetDailyLetters(b businessPick) gin.HandlerFunc {
	return func(c *gin.Context) {
		p, err := b.GetDailyWord(c.Request.Context())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		st := struct {
			ID      uuid.UUID `json:"id"`
			Letters []string  `json:"letters"`
		}{
			ID:      p.ID,
			Letters: p.Letters,
		}
		c.JSON(http.StatusOK, st)
	}
}

func CheckWordExists(b businessPick) gin.HandlerFunc {
	return func(c *gin.Context) {
		rqt := struct {
			Word string `json:"word"`
		}{}
		err := c.ShouldBindJSON(&rqt)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		score, err := b.CheckWordExists(c.Request.Context(), rqt.Word)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		st := struct {
			Exists bool   `json:"exists"`
			Word   string `json:"word"`
			Score  int    `json:"score"`
		}{
			Exists: score.Exist,
			Word:   score.Word,
			Score:  score.Score,
		}
		c.JSON(http.StatusOK, st)
	}
}
