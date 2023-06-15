package model

type UserId int

type Follow struct {
	UserId      UserId `json:"user_id"`
	FollowingId UserId `json:"following_id"`
}
