package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type businessUser interface {
	InsertScore(ctx context.Context, user, pick uuid.UUID, score int) error
	InsertUser(ctx context.Context, username string) (uuid.UUID, error)
}

func InsertUser(business businessUser) gin.HandlerFunc {
	return func(c *gin.Context) {
		rqt := struct {
			Name string `json:"name"`
		}{}
		err := c.ShouldBindJSON(&rqt)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		id, err := business.InsertUser(c.Request.Context(), rqt.Name)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		st := struct {
			ID uuid.UUID `json:"id"`
		}{
			ID: id,
		}
		c.JSON(http.StatusOK, st)
	}
}

func InsertScore(business businessUser) gin.HandlerFunc {
	return func(c *gin.Context) {
		rqt := struct {
			User  uuid.UUID `json:"user"`
			Pick  uuid.UUID `json:"pick"`
			Score int       `json:"score"`
		}{}
		err := c.ShouldBindJSON(&rqt)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		err = business.InsertScore(c.Request.Context(), rqt.User, rqt.Pick, rqt.Score)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Status(http.StatusNoContent)
	}
}
