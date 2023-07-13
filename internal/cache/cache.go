package cache

import (
	"fmt"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
)

type Cache interface {
	Put(key string, value []byte) error
	Get(key string) (bool, []byte)
	Remove(key string) error
	StartsWith(key string) (error, [][]byte)
}

// function to generate keys for cache
func GenerateUserId(id types.UserId) string {
	return fmt.Sprintf("user_id_%d", id)
}

func GenerateUserToTweetId(user_id types.UserId, tweet_id types.TweetId) string {
	return fmt.Sprintf("user_id_%d_tweet_id%d", user_id, tweet_id)
}

func GenerateTweetId(id types.TweetId) string {
	return fmt.Sprintf("tweet_id_%d", id)
}
