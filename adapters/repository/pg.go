package repository

import (
	"context"
	"strings"

	"github.com/frisk038/wordcraft/business/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	errDupl = `23505`

	checkWordExists = `SELECT word FROM words WHERE word=UPPER($1);`
	insertLetters   = `INSERT INTO picks(letters) VALUES($1) RETURNING id;`
	insertUser      = `INSERT INTO users(name) VALUES($1) RETURNING id;`
	selectUser      = `SELECT id FROM users WHERE name = $1;`
	insertScore     = `INSERT INTO scores(userid,picks,score) VALUES($1, $2, $3);`
	getDailyWord    = `SELECT id ,letters FROM picks WHERE DATE_trunc('day',dt)=DATE_TRUNC('day', NOW());`
	getLeaderBoard  = `SELECT users.name, scores.score 
					   FROM scores 
					   LEFT JOIN users ON users.id=scores.userid 
					   WHERE picks =(SELECT id  FROM picks WHERE DATE_trunc('day',dt)=DATE_TRUNC('day', NOW()))
					   ORDER BY score DESC LIMIT 5;`
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
	var (
		pick    models.Pick
		letters string
	)
	if err := rows.Scan(&pick.ID, &letters); err != nil {
		if err == pgx.ErrNoRows {
			return models.Pick{}, models.ErrNoDailyPick
		}
		return models.Pick{}, err
	}
	l := strings.Split(letters, "")
	pick.Letters = make([]string, len(l))
	for i, v := range l {
		pick.Letters[i] = strings.ToLower(v)
	}

	return pick, nil
}

func (c *Client) InsertLetters(ctx context.Context, letters []string) (uuid.UUID, error) {
	rows := c.conn.QueryRow(ctx, insertLetters, strings.Join(letters, ""))
	var pickID uuid.UUID
	if err := rows.Scan(&pickID); err != nil {
		return uuid.UUID{}, err
	}

	return pickID, nil
}

func (c *Client) GetLeaderBoard(ctx context.Context) ([]models.UserScore, error) {
	rows, err := c.conn.Query(ctx, getLeaderBoard)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userScores []models.UserScore
	for rows.Next() {
		var us models.UserScore
		err := rows.Scan(&us.User, &us.Score)
		if err != nil {
			return nil, err
		}
		userScores = append(userScores, us)
	}

	return userScores, rows.Err()
}
