package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
	_ "github.com/go-sql-driver/mysql"
)

// Repository defines a MySQL-based repository.
type Repository struct {
	db *sql.DB
}

// New creates a new MySQL-based repository.
func New(driverName string, dataSourceName string) (*Repository, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &Repository{db}, nil
}

// Like specific post
func (r *Repository) Like(ctx context.Context, userId types.UserId, tweetId types.TweetId) error {
	row, err := r.db.ExecContext(
		ctx,
		"INSERT INTO Likes (user_id, tweet_id) VALUES (?, ?)", userId, tweetId,
	)
	if err != nil {
		return err
	}

	_, err = row.RowsAffected()
	if err != nil {
		return ErrNotFound
	}
	return nil
}

// Unlike specific post
func (r *Repository) Unlike(ctx context.Context, userId types.UserId, tweetId types.TweetId) error {
	row, err := r.db.ExecContext(ctx, "DELETE FROM Likes WHERE user_id = ? AND tweet_id = ?", userId, tweetId)
	if err != nil {
		return err
	}

	_, err = row.RowsAffected()
	if err != nil {
		return ErrNotFound
	}
	return nil
}

func get(ctx context.Context, r *Repository, colName string, idName string, id interface{}) ([]int, error) {
	query := fmt.Sprintf("SELECT %s FROM Likes WHERE %s = ?", colName, idName)
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []int
	// iterate over result
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		res = append(res, id)
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("there is no data for your request")
	}
	return res, nil
}

func (r *Repository) GetUsersByTweet(ctx context.Context, tweetId types.TweetId) ([]types.UserId, error) {
	users, err := get(ctx, r, "user_id", "tweet_id", tweetId)
	if err != nil {
		return nil, err
	}

	userIds := make([]types.UserId, len(users))
	for i, v := range users {
		id := types.UserId(v)
		userIds[i] = id
	}
	return userIds, nil
}

func (r *Repository) GetTweetsByUser(ctx context.Context, userId types.UserId) ([]types.TweetId, error) {
	tweets, err := get(ctx, r, "tweet_id", "user_id", userId)
	if err != nil {
		return nil, err
	}

	tweetIds := make([]types.TweetId, len(tweets))
	for i, v := range tweets {
		id := types.TweetId(v)
		tweetIds[i] = id
	}
	return tweetIds, nil
}
