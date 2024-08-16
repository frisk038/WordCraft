package business

import (
	"context"

	"github.com/frisk038/wordcraft/business/models"
	"github.com/google/uuid"
)

type userStore interface {
	InsertUser(ctx context.Context, name string) (uuid.UUID, error)
	InsertScore(ctx context.Context, user, pick uuid.UUID, score int) error
	GetLeaderBoard(ctx context.Context) ([]models.UserScore, error)
}

type BUsers struct {
	store userStore
}

func NewBUsers(store userStore) BUsers {
	return BUsers{store: store}
}

func (bp *BUsers) InsertUser(ctx context.Context, username string) (uuid.UUID, error) {
	return bp.store.InsertUser(ctx, username)
}

func (bp *BUsers) InsertScore(ctx context.Context, user, pick uuid.UUID, score int) error {
	return bp.store.InsertScore(ctx, user, pick, score)
}

func (bp *BUsers) GetLeaderBoard(ctx context.Context) ([]models.UserScore, error) {
	return bp.store.GetLeaderBoard(ctx)
}
