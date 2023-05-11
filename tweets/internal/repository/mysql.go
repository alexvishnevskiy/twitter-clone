package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/alexvishnevskiy/twitter-clone/tweets/pkg/model"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"time"
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

func (r *Repository) Put(
	ctx context.Context,
	userId model.UserId,
	content string,
	mediaUrl *string,
	retweetId *model.TweetId,
) (*int64, error) {
	createdAt := time.Now()
	row, err := r.db.ExecContext(
		ctx,
		"INSERT INTO Tweets (user_id, retweet_id, content, media_url, created_at) VALUES (?, ?, ?, ?, ?)",
		userId, retweetId, content, mediaUrl, createdAt,
	)
	if err != nil {
		return nil, err
	}
	id, err := row.LastInsertId()
	return &id, err
}

func get(ctx context.Context, r *Repository, idName string, ids []interface{}) ([]model.Tweet, error) {
	placeholder := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	query := fmt.Sprintf("SELECT * FROM Tweets WHERE %s IN (%s)", idName, placeholder)
	rows, err := r.db.QueryContext(ctx, query, ids...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.Tweet
	for rows.Next() {
		var tweet model.Tweet
		var createdAtStr string

		if err := rows.Scan(
			&tweet.TweetId, &tweet.UserId,
			&tweet.RetweetId, &tweet.Content,
			&tweet.MediaUrl, &createdAtStr,
		); err != nil {
			return nil, err
		}

		const layout = "2006-01-02 15:04:05"
		createdAt, err := time.Parse(layout, createdAtStr)
		if err != nil {
			return nil, err
		}
		tweet.CreatedAt = createdAt
		res = append(res, tweet)
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("there are no tweets for your request")
	}
	return res, nil
}

func (r *Repository) GetByTweet(ctx context.Context, tweetIds ...model.TweetId) ([]model.Tweet, error) {
	// TODO: probably add some filters
	interfaceTweetIds := make([]interface{}, len(tweetIds))
	for i, v := range tweetIds {
		interfaceTweetIds[i] = v
	}
	return get(ctx, r, "tweet_id", interfaceTweetIds)
}

func (r *Repository) GetByUser(ctx context.Context, userIds ...model.UserId) ([]model.Tweet, error) {
	// TODO: probably add some filters
	interfaceTweetIds := make([]interface{}, len(userIds))
	for i, v := range userIds {
		interfaceTweetIds[i] = v
	}
	return get(ctx, r, "user_id", interfaceTweetIds)
}

func (r *Repository) DeletePost(ctx context.Context, postId model.TweetId) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM Tweets WHERE tweet_id = ?", postId)
	return err
}
