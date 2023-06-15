package model

import "github.com/alexvishnevskiy/twitter-clone/internal/types"

type Follow struct {
	UserId      types.UserId `json:"user_id"`
	FollowingId types.UserId `json:"following_id"`
}
