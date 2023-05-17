package controller

import (
	"context"
	"github.com/alexvishnevskiy/twitter-clone/likes/pkg/model"
)

type likesRepository interface {
	Like(ctx context.Context, userId model.UserId, tweetId model.TweetId) error
	Unlike(ctx context.Context, userId model.UserId, tweetId model.TweetId) error
	GetUsersByTweet(ctx context.Context, tweetId model.TweetId) ([]model.UserId, error)
	GetTweetsByUser(ctx context.Context, userId model.UserId) ([]model.TweetId, error)
}

type Controller struct {
	repo likesRepository
}

func New(repo likesRepository) *Controller {
	return &Controller{repo: repo}
}

func (ctrl *Controller) LikeTweet(
	ctx context.Context,
	userId model.UserId,
	tweetId model.TweetId,
) error {
	err := ctrl.repo.Like(ctx, userId, tweetId)
	return err
}

func (ctrl *Controller) UnlikeTweet(
	ctx context.Context,
	userId model.UserId,
	tweetId model.TweetId,
) error {
	err := ctrl.repo.Unlike(ctx, userId, tweetId)
	return err
}

func (ctrl *Controller) GetUsersByTweet(
	ctx context.Context,
	tweetId model.TweetId,
) ([]model.UserId, error) {
	users, err := ctrl.repo.GetUsersByTweet(ctx, tweetId)
	return users, err
}

func (ctrl *Controller) GetTweetsByUser(
	ctx context.Context,
	userId model.UserId,
) ([]model.TweetId, error) {
	tweets, err := ctrl.repo.GetTweetsByUser(ctx, userId)
	return tweets, err
}