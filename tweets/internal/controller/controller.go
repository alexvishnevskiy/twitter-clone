package controller

import (
	"context"
	"github.com/alexvishnevskiy/twitter-clone/tweets/pkg/model"
)

// defines abstract methods for business logic (ex. handling requests to database)
// message communication

type tweetsRepository interface {
	Put(ctx context.Context, userId model.UserId, content string, mediaUrl *string, retweetId *model.TweetId) (*model.TweetId, error)
	GetByTweet(ctx context.Context, tweetIds ...model.TweetId) ([]model.Tweet, error)
	GetByUser(ctx context.Context, userIds ...model.UserId) ([]model.Tweet, error)
	DeletePost(ctx context.Context, postId model.TweetId) error
}

// controller for tweets
type Controller struct {
	repo tweetsRepository
	// TODO: probably there will be message queue
}

// Creates new tweets controller
func New(repo tweetsRepository) *Controller {
	return &Controller{repo}
}

func (ctrl *Controller) PostNewTweet(
	ctx context.Context,
	userId model.UserId,
	content string,
	mediaUrl *string,
	retweetId *model.TweetId,
) (*model.TweetId, error) {
	tweetId, err := ctrl.repo.Put(ctx, userId, content, mediaUrl, retweetId)
	return tweetId, err
}

func (ctrl *Controller) RetrieveByTweetID(ctx context.Context, tweetIds ...model.TweetId) ([]model.Tweet, error) {
	tweets, err := ctrl.repo.GetByTweet(ctx, tweetIds...)
	return tweets, err
}

func (ctrl *Controller) RetrieveByUserID(ctx context.Context, userIds ...model.UserId) ([]model.Tweet, error) {
	tweets, err := ctrl.repo.GetByUser(ctx, userIds...)
	return tweets, err
}

func (ctrl *Controller) DeletePost(ctx context.Context, postId model.TweetId) error {
	err := ctrl.repo.DeletePost(ctx, postId)
	return err
}

// TODO: probably add more logic
