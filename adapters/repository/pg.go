package repository

import (
	"context"

	"github.com/frisk038/wordcraft/business/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	errDupl = `23505`

	checkWordExists = `SELECT word FROM words WHERE word=$1;`
	pickWord        = `SELECT id,word FROM words WHERE len > 4 ORDER BY RANDOM() LIMIT 1;`
	insertWord      = `INSERT INTO picks(pick) VALUES($1) RETURNING id;`
	insertUser      = `INSERT INTO users(name) VALUES($1) RETURNING id;`
	selectUser      = `SELECT id FROM users WHERE name = $1;`
	insertScore     = `INSERT INTO scores(userid,picks,score) VALUES($1, $2, $3);`
	getDailyWord    = `SELECT picks.id ,words.word 
					   FROM picks LEFT JOIN words on picks.pick=words.id 
					   WHERE DATE_trunc('day',dt)=DATE_TRUNC('day', NOW())`
)

type Client struct {
	conn *pgxpool.Pool
}

func NewClient(url string) (*Client, error) {
	conn, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

func (c *Client) CheckWordExists(ctx context.Context, word string) (bool, error) {
	rows := c.conn.QueryRow(ctx, checkWordExists, word)
	var text string
	err := rows.Scan(&text)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c *Client) getUser(ctx context.Context, name string) (uuid.UUID, error) {
	rows := c.conn.QueryRow(ctx, selectUser, name)
	var userID uuid.UUID
	if err := rows.Scan(&userID); err != nil {
		return uuid.UUID{}, err
	}
	return userID, nil
}

func (c *Client) InsertUser(ctx context.Context, name string) (uuid.UUID, error) {
	rows := c.conn.QueryRow(ctx, insertUser, name)
	var userID uuid.UUID
	if err := rows.Scan(&userID); err != nil {
		pgErr := err.(*pgconn.PgError)
		if err == pgx.ErrNoRows || pgErr.Code == errDupl {
			userID, err = c.getUser(ctx, name)
			if err != nil {
				return uuid.UUID{}, err
			}
			return userID, nil
		}
		return uuid.UUID{}, err
	}
	return userID, nil
}

func (c *Client) InsertScore(ctx context.Context, user, pick uuid.UUID, score int) error {
	if _, err := c.conn.Exec(ctx, insertScore, user, pick, score); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetDailyWord(ctx context.Context) (models.Pick, error) {
	rows := c.conn.QueryRow(ctx, getDailyWord)
	var pick models.Pick
	if err := rows.Scan(&pick.ID, &pick.Word); err != nil {
		if err == pgx.ErrNoRows {
			return models.Pick{}, models.ErrNoDailyPick
		}
		return models.Pick{}, err
	}
	return pick, nil
}

func (c *Client) PickDailyWord(ctx context.Context) (models.Pick, error) {
	tx, err := c.conn.Begin(ctx)
	if err != nil {
		return models.Pick{}, err
	}
	defer tx.Rollback(ctx)

	rows := tx.QueryRow(ctx, pickWord)
	var (
		id   uuid.UUID
		pick models.Pick
	)
	if err = rows.Scan(&id, &pick.Word); err != nil {
		return models.Pick{}, err
	}

	rows = tx.QueryRow(ctx, insertWord, id)
	if err = rows.Scan(&pick.ID); err != nil {
		return models.Pick{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		return models.Pick{}, err
	}

	return pick, nil
}
