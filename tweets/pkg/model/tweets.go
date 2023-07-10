package model

import "time"
import "github.com/alexvishnevskiy/twitter-clone/internal/types"

// tweets data types
type Tweet struct {
	UserId    types.UserId   `json:"user_id"`
	TweetId   types.TweetId  `json:"tweet_id"`
	RetweetId *types.TweetId `json:"retweet_id"`
	Content   string         `json:"content"`
	MediaUrl  *string        `json:"media_url"`
	CreatedAt time.Time      `json:"created_at"`
}

// struct for media
type Media struct {
	Media     string    `json:"media"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
