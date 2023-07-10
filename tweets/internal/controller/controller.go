package controller

import (
	"context"
	"github.com/alexvishnevskiy/twitter-clone/internal/storage"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
	"github.com/alexvishnevskiy/twitter-clone/tweets/pkg/model"
	"mime/multipart"
)

// defines abstract methods for business logic (ex. handling requests to database)
// message communication

type tweetsRepository interface {
	Put(ctx context.Context, userId types.UserId, content string, mediaUrl *string, retweetId *types.TweetId) (*types.TweetId, error)
	GetByTweet(ctx context.Context, tweetIds ...types.TweetId) ([]model.Tweet, error)
	GetByUser(ctx context.Context, userIds ...types.UserId) ([]model.Tweet, error)
	DeletePost(ctx context.Context, postId types.TweetId) error
}

// controller for tweets
type Controller struct {
	repo    tweetsRepository
	storage storage.Storage
}

// Creates new tweets controller
func New(repo tweetsRepository, storage storage.Storage) *Controller {
	return &Controller{repo, storage}
}

func (ctrl *Controller) PostNewTweet(
	ctx context.Context,
	file multipart.File,
	handler *multipart.FileHeader,
	userId types.UserId,
	content string,
	retweetId *types.TweetId,
) (*types.TweetId, error) {
	// save to storage
	url, err := ctrl.storage.SaveImageFromRequest(file, handler)
	if err != nil {
		return nil, err
	}

	// save to db
	MediaUrl := &url
	tweetId, err := ctrl.repo.Put(ctx, userId, content, MediaUrl, retweetId)
	return tweetId, err
}

func (ctrl *Controller) RetrieveByTweetID(ctx context.Context, tweetIds ...types.TweetId) ([]model.Media, error) {
	tweets, err := ctrl.repo.GetByTweet(ctx, tweetIds...)
	if err != nil {
		return nil, err
	}

	tweetsMedia := make([]model.Media, len(tweets))
	// converting to response object
	for i, tweet := range tweets {
		media, _ := ctrl.storage.ConvertImageFromStorage(*tweet.MediaUrl)
		tweetsMedia[i].Content = tweet.Content
		tweetsMedia[i].CreatedAt = tweet.CreatedAt
		tweetsMedia[i].Media = media
	}
	return tweetsMedia, err
}

func (ctrl *Controller) RetrieveByUserID(ctx context.Context, userIds ...types.UserId) ([]model.Media, error) {
	tweets, err := ctrl.repo.GetByUser(ctx, userIds...)
	if err != nil {
		return nil, err
	}

	var media string = ""
	tweetsMedia := make([]model.Media, len(tweets))
	// converting to response object
	for i, tweet := range tweets {
		if tweet.MediaUrl != nil {
			media, _ = ctrl.storage.ConvertImageFromStorage(*tweet.MediaUrl)
		}

		tweetsMedia[i].Content = tweet.Content
		tweetsMedia[i].CreatedAt = tweet.CreatedAt
		tweetsMedia[i].Media = media
	}
	return tweetsMedia, err
}

func (ctrl *Controller) DeletePost(ctx context.Context, postId types.TweetId) error {
	// get media url
	tweetData, err := ctrl.repo.GetByTweet(ctx, postId)
	if err != nil {
		return err
	}
	// delete from db
	err = ctrl.repo.DeletePost(ctx, postId)
	if err != nil {
		return err
	}
	// delete from storage
	err = ctrl.storage.Delete(*tweetData[0].MediaUrl)
	return err
}

// TODO: probably add more logic
