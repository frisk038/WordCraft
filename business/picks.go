package business

import (
	"context"
	"strings"

	"github.com/frisk038/wordcraft/business/models"
)

type pickStore interface {
	CheckWordExists(ctx context.Context, word string) (bool, error)
	GetDailyWord(ctx context.Context) (models.Pick, error)
	PickDailyWord(ctx context.Context) (models.Pick, error)
}

type BPicks struct {
	store pickStore
}

func NewBPicks(store pickStore) BPicks {
	return BPicks{store: store}
}

func (bp *BPicks) GetDailyWord(ctx context.Context) (models.Pick, error) {
	p, err := bp.store.GetDailyWord(ctx)
	switch err {
	case models.ErrNoDailyPick:
		p, err = bp.store.PickDailyWord(ctx)
		if err != nil {
			return models.Pick{}, err
		}
		fallthrough
	case nil:
		p.Letters = strings.Split(p.Word, "")
		return p, nil
	default:
		return models.Pick{}, err
	}
}

func (bp *BPicks) CheckWordExists(ctx context.Context, word string) (bool, error) {
	return bp.store.CheckWordExists(ctx, word)
}
