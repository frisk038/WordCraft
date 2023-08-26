package business

import (
	"context"
	"math/rand"

	"github.com/frisk038/wordcraft/business/models"
	"github.com/google/uuid"
)

type pickStore interface {
	CheckWordExists(ctx context.Context, word string) (bool, error)
	GetDailyWord(ctx context.Context) (models.Pick, error)
	InsertLetters(ctx context.Context, letters []string) (uuid.UUID, error)
}

type BPicks struct {
	store pickStore
}

var alphabet = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j",
	"k", "l", "m", "n", "o", "p", "q", "r", "s", "t",
	"u", "v", "w", "x", "y", "z"}

func NewBPicks(store pickStore) BPicks {
	return BPicks{store: store}
}

func (bp *BPicks) GetDailyWord(ctx context.Context) (models.Pick, error) {
	p, err := bp.store.GetDailyWord(ctx)
	if err != nil {
		if err == models.ErrNoDailyPick {
			rand.Shuffle(26, func(i, j int) { alphabet[i], alphabet[j] = alphabet[j], alphabet[i] })
			p.Letters = alphabet[:9]
			p.ID, err = bp.store.InsertLetters(ctx, p.Letters)
			if err != nil {
				return models.Pick{}, err
			}
			return p, nil
		}
		return models.Pick{}, err
	}

	return p, err
}

func (bp *BPicks) CheckWordExists(ctx context.Context, word string) (bool, error) {
	return bp.store.CheckWordExists(ctx, word)
}
