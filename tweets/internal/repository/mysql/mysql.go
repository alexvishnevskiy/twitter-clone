package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
	"github.com/alexvishnevskiy/twitter-clone/tweets/pkg/model"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"time"
)

// time layout
const layout = "2006-01-02 15:04:05"

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

// Put new tweet to database
func (r *Repository) Put(
	ctx context.Context,
	userId types.UserId,
	content string,
	mediaUrl *string,
	retweetId *types.TweetId,
) (*types.TweetId, error) {
	createdAt := time.Now()
	row, err := r.db.ExecContext(
		ctx,
		"INSERT INTO Tweets (user_id, retweet_id, content, media_url, created_at) VALUES (?, ?, ?, ?, ?)",
		userId, retweetId, content, mediaUrl, createdAt.Format(layout),
	)
	if err != nil {
		return nil, err
	}
	id, err := row.LastInsertId()
	tweetId := types.TweetId(id)
	return &tweetId, err
}

// helper function to retrieve tweets from database
func get(ctx context.Context, r *Repository, idName string, ids []interface{}) ([]model.Tweet, error) {
	placeholder := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	query := fmt.Sprintf("SELECT * FROM Tweets WHERE %s IN (%s)", idName, placeholder)
	rows, err := r.db.QueryContext(ctx, query, ids...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.Tweet
	// iterate over result
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

// GetByTweet Retrieve by tweet id
func (r *Repository) GetByTweet(ctx context.Context, tweetIds ...types.TweetId) ([]model.Tweet, error) {
	// TODO: probably add some filters
	interfaceTweetIds := make([]interface{}, len(tweetIds))
	for i, v := range tweetIds {
		interfaceTweetIds[i] = v
	}
	return get(ctx, r, "tweet_id", interfaceTweetIds)
}

// GetByUser Retieve by user id
func (r *Repository) GetByUser(ctx context.Context, userIds ...types.UserId) ([]model.Tweet, error) {
	// TODO: probably add some filters
	interfaceTweetIds := make([]interface{}, len(userIds))
	for i, v := range userIds {
		interfaceTweetIds[i] = v
	}
	return get(ctx, r, "user_id", interfaceTweetIds)
}

// DeletePost delete post by tweet id
func (r *Repository) DeletePost(ctx context.Context, postId types.TweetId) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM Tweets WHERE tweet_id = ?", postId)
	return err
}
