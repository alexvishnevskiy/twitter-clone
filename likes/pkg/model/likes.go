package model

type UserId int
type TweetId int

// likes data model
type Like struct {
	UserId  UserId  `json:"user_id"`
	TweetId TweetId `json:"tweet_id"`
}
