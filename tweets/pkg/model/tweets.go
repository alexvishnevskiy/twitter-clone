package model

import "time"

type UserId int
type TweetId int

type Tweet struct {
	UserId    UserId    `json:"user_id"`
	TweetId   TweetId   `json:"tweet_id"`
	RetweetId *TweetId  `json:"retweet_id"`
	Content   string    `json:"content"`
	MediaUrl  *string   `json:"media_url"`
	CreatedAt time.Time `json:"created_at"`
}
