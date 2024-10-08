package business

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

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

var consonne = []string{"b", "c", "d", "f", "g", "h", "j",
	"k", "l", "m", "n", "p", "q", "r", "s", "t",
	"v", "w", "x", "z"}

var voyelle = []string{"a", "e", "i", "o", "u", "y"}

var values = map[string]int{"a": 1, "b": 3, "c": 3, "d": 2, "e": 1, "f": 4,
	"g": 2, "h": 4, "i": 1, "j": 8, "k": 5, "l": 1, "m": 3, "n": 1, "o": 1, "p": 3,
	"q": 10, "r": 1, "s": 1, "t": 1, "u": 1, "v": 4, "w": 4, "x": 8, "y": 4, "z": 10}

func NewBPicks(store pickStore) BPicks {
	return BPicks{store: store}
}

func (bp *BPicks) GetDailyWord(ctx context.Context) (models.Pick, error) {
	p, err := bp.store.GetDailyWord(ctx)
	if err != nil {
		if err == models.ErrNoDailyPick {
			rand.Shuffle(len(voyelle), func(i, j int) { voyelle[i], voyelle[j] = voyelle[j], voyelle[i] })
			rand.Shuffle(len(consonne), func(i, j int) { consonne[i], consonne[j] = consonne[j], consonne[i] })

			p.Letters = consonne[:7]
			p.Letters = append(p.Letters, voyelle[:2]...)
			rand.Shuffle(len(p.Letters), func(i, j int) { p.Letters[i], p.Letters[j] = p.Letters[j], p.Letters[i] })
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

func (bp *BPicks) CheckWordExists(ctx context.Context, word string) (models.Score, error) {
	score := models.Score{
		Word: word,
	}
	var err error
	score.Exist, err = bp.store.CheckWordExists(ctx, word)
	if err != nil {
		return models.Score{}, err
	}

	letters := strings.Split(word, "")
	for _, v := range letters {
		fmt.Println(values[v])
		score.Score += values[v]
	}

	fmt.Println(letters, score)
	return score, nil
}
