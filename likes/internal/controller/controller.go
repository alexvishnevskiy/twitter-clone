package controller

import (
	"context"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
)

type likesRepository interface {
	Like(ctx context.Context, userId types.UserId, tweetId types.TweetId) error
	Unlike(ctx context.Context, userId types.UserId, tweetId types.TweetId) error
	GetUsersByTweet(ctx context.Context, tweetId types.TweetId) ([]types.UserId, error)
	GetTweetsByUser(ctx context.Context, userId types.UserId) ([]types.TweetId, error)
}

type Controller struct {
	repo likesRepository
}

func New(repo likesRepository) *Controller {
	return &Controller{repo: repo}
}

func (ctrl *Controller) LikeTweet(
	ctx context.Context,
	userId types.UserId,
	tweetId types.TweetId,
) error {
	err := ctrl.repo.Like(ctx, userId, tweetId)
	return err
}

func (ctrl *Controller) UnlikeTweet(
	ctx context.Context,
	userId types.UserId,
	tweetId types.TweetId,
) error {
	err := ctrl.repo.Unlike(ctx, userId, tweetId)
	return err
}

func (ctrl *Controller) GetUsersByTweet(
	ctx context.Context,
	tweetId types.TweetId,
) ([]types.UserId, error) {
	users, err := ctrl.repo.GetUsersByTweet(ctx, tweetId)
	return users, err
}

func (ctrl *Controller) GetTweetsByUser(
	ctx context.Context,
	userId types.UserId,
) ([]types.TweetId, error) {
	tweets, err := ctrl.repo.GetTweetsByUser(ctx, userId)
	return tweets, err
}
