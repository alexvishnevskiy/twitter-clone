package model

import "github.com/alexvishnevskiy/twitter-clone/internal/types"

// likes data types
type Like struct {
	UserId  types.UserId  `json:"user_id"`
	TweetId types.TweetId `json:"tweet_id"`
}
