package model

import "github.com/alexvishnevskiy/twitter-clone/internal/types"

type User struct {
	UserId    types.UserId `json:"user_id"`
	Nickname  string       `json:"nickname"`
	FirstName string       `json:"first_name"`
	LastName  string       `json:"last_name"`
	Email     string       `json:"email"`
	Password  string       `json:"password"`
}
