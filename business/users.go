package business

import (
	"context"

	"github.com/google/uuid"
)

type userStore interface {
	InsertUser(ctx context.Context, name string) (uuid.UUID, error)
	InsertScore(ctx context.Context, user, pick uuid.UUID, score int) error
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
