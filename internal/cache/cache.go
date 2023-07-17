package cache

import (
	"fmt"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
)

// interface to implement cache
type Cache interface {
	Put(key string, value []byte) error
	Get(key string) (bool, []byte)
	Remove(key string) error
	StartsWith(key string) (error, [][]byte)
}

// function to generate user_id keys for cache
func GenerateUserId(id types.UserId) string {
	return fmt.Sprintf("user_id_%d", id)
}

// function to generate user_id_tweet_id keys for cache
func GenerateUserToTweetId(user_id types.UserId, tweet_id types.TweetId) string {
	return fmt.Sprintf("user_id_%d_tweet_id%d", user_id, tweet_id)
}

// function to generate tweet_id keys for cache
func GenerateTweetId(id types.TweetId) string {
	return fmt.Sprintf("tweet_id_%d", id)
}
