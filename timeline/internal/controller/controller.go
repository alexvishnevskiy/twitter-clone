package controller

import (
	"context"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
	"github.com/alexvishnevskiy/twitter-clone/tweets/pkg/model"
)

type tweetsGateway interface {
	GetTweets(ctx context.Context, userId ...types.UserId) ([]model.Media, error)
}

type followGateway interface {
	GetUsers(ctx context.Context, userId types.UserId) ([]types.UserId, error)
}

// controller for timeline
type Controller struct {
	TweetsService tweetsGateway
	FollowService followGateway
}

func New(tweets tweetsGateway, follow followGateway) *Controller {
	return &Controller{tweets, follow}
}

// get all tweets from the users who this user is following
func (ctrl *Controller) GetHomeTimeline(ctx context.Context, userId types.UserId) ([]model.Media, error) {
	users, err := ctrl.FollowService.GetUsers(ctx, userId)
	if err != nil {
		return nil, err
	}

	tweets, err := ctrl.TweetsService.GetTweets(ctx, users...)
	return tweets, err
}

func (ctrl *Controller) GetMentionsTimeline(ctx context.Context, userId types.UserId) ([]model.Tweet, error) {
	// TODO
	return nil, nil
}

func (ctrl *Controller) GetGlobalTimeline(ctx context.Context, userId types.UserId) ([]model.Tweet, error) {
	// TODO
	return nil, nil
}
