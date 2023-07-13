package controller

import (
	"context"
	"encoding/json"
	cachestorage "github.com/alexvishnevskiy/twitter-clone/internal/cache"
	"github.com/alexvishnevskiy/twitter-clone/internal/storage"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
	"github.com/alexvishnevskiy/twitter-clone/tweets/pkg/model"
	"mime/multipart"
	"sort"
	"time"
)

// defines abstract methods for business logic (ex. handling requests to database)
// message communication

type tweetsRepository interface {
	Put(ctx context.Context, userId types.UserId, content string, mediaUrl *string, retweetId *types.TweetId) (types.TweetId, time.Time, error)
	GetByTweet(ctx context.Context, tweetIds ...types.TweetId) ([]model.Tweet, error)
	GetByUser(ctx context.Context, userIds ...types.UserId) ([]model.Tweet, error)
	DeletePost(ctx context.Context, postId types.TweetId) error
}

// controller for tweets
type Controller struct {
	repo    tweetsRepository
	storage storage.Storage
	cache   cachestorage.Cache
}

// Creates new tweets controller
func New(repo tweetsRepository, storage storage.Storage, cache cachestorage.Cache) *Controller {
	return &Controller{repo, storage, cache}
}

// check one user_id from cache for tweets
// user_id -> tweets
func getUserIdFromCache(cache cachestorage.Cache, userId types.UserId) ([]model.Tweet, error) {
	var (
		tweet  model.Tweet
		tweets []model.Tweet
	)

	// return tweets for specific user_id
	genUserId := cachestorage.GenerateUserId(userId)
	if err, tweets_ := cache.StartsWith(genUserId); err == nil {
		for _, tweet_ := range tweets_ {
			err := json.Unmarshal(tweet_, &tweet)
			if err != nil {
				return nil, err
			}
			tweets = append(tweets, tweet)
		}
	}
	return tweets, cachestorage.CacheError
}

// retrieve all user ids from cache
func retrieveByUserIds(cache cachestorage.Cache, userIds ...types.UserId) (bool, []model.Tweet, []types.UserId) {
	var (
		tweets          []model.Tweet
		indicesToRemove []int
		checkDb         bool
	)
	// retrieve all tweets that we have in cache
	for i, userId := range userIds {
		tweets_, err := getUserIdFromCache(cache, userId)
		if err == nil {
			tweets = append(tweets, tweets_...)
			indicesToRemove = append(indicesToRemove, i)
		}
	}
	// condition to check db
	checkDb = len(indicesToRemove) != len(userIds)

	// truncate original tweetIds
	sort.Sort(sort.Reverse(sort.IntSlice(indicesToRemove)))
	for _, i := range indicesToRemove {
		userIds = append(userIds[:i], userIds[i+1:]...)
	}
	return checkDb, tweets, userIds
}

// put to cache
func putUserIdToCache(cache cachestorage.Cache, userId types.UserId, tweet model.Tweet) error {
	// put user_id to cache
	key := cachestorage.GenerateUserToTweetId(userId, tweet.TweetId)
	data, err := json.Marshal(tweet)
	if err != nil {
		return err
	}
	err = cache.Put(key, data)
	return err
}

// check one tweet_id from cache for tweets
// tweet_id -> tweet
func getTweetIdFromCache(cache cachestorage.Cache, tweetId types.TweetId) (model.Tweet, error) {
	var tweet model.Tweet

	genTweetId := cachestorage.GenerateTweetId(tweetId)
	if ok, val := cache.Get(genTweetId); ok {
		err := json.Unmarshal(val, &tweet)
		return tweet, err
	}
	return tweet, cachestorage.CacheError
}

// retrieve all tweet ids from cache
func retrieveByTweetIds(cache cachestorage.Cache, tweetIds ...types.TweetId) (bool, []model.Tweet, []types.TweetId) {
	var (
		tweets          []model.Tweet
		indicesToRemove []int
		checkDb         bool
	)
	// retrieve all tweets that we have in cache
	for i, tweetId := range tweetIds {
		tweet_, err := getTweetIdFromCache(cache, tweetId)
		if err == nil {
			tweets = append(tweets, tweet_)
			indicesToRemove = append(indicesToRemove, i)
		}
	}
	// condition to check db
	checkDb = len(indicesToRemove) != len(tweetIds)

	// truncate original tweetIds
	sort.Sort(sort.Reverse(sort.IntSlice(indicesToRemove)))
	for _, i := range indicesToRemove {
		tweetIds = append(tweetIds[:i], tweetIds[i+1:]...)
	}
	return checkDb, tweets, tweetIds
}

// put to cache
func putTweetIdToCache(cache cachestorage.Cache, tweetId types.TweetId, tweet model.Tweet) error {
	tweetKey := cachestorage.GenerateTweetId(tweetId)

	data, err := json.Marshal(tweet)
	if err != nil {
		return err
	}
	err = cache.Put(tweetKey, data)
	return err
}

func (ctrl *Controller) PostNewTweet(
	ctx context.Context,
	file multipart.File,
	handler *multipart.FileHeader,
	userId types.UserId,
	content string,
	retweetId *types.TweetId,
) (*types.TweetId, error) {
	var (
		MediaUrl *string = nil
		url      string
		err      error
	)

	// save to storage
	if handler != nil {
		url, err = ctrl.storage.SaveImageFromRequest(file, handler)
		if err != nil {
			return nil, err
		}
	}

	// save to db
	MediaUrl = &url
	tweetId, time, err := ctrl.repo.Put(ctx, userId, content, MediaUrl, retweetId)
	// tweet metadata
	tweet := model.Tweet{
		userId,
		tweetId,
		retweetId,
		content,
		MediaUrl,
		time,
	}

	// save to cache
	if ctrl.cache != nil && err == nil {
		err = putTweetIdToCache(ctrl.cache, tweetId, tweet)
		if err != nil {
			return nil, err
		}
		err = putUserIdToCache(ctrl.cache, userId, tweet)
		if err != nil {
			return nil, err
		}
	}
	return &tweetId, err
}

func (ctrl *Controller) RetrieveByTweetID(ctx context.Context, tweetIds ...types.TweetId) ([]model.Media, error) {
	var (
		checkDb = true
		tweets  []model.Tweet
	)

	// retrieve from cache
	if ctrl.cache != nil {
		checkDb, tweets, tweetIds = retrieveByTweetIds(ctrl.cache, tweetIds...)
	}

	// retrieve all remaining tweets
	if checkDb {
		repoTweets, err := ctrl.repo.GetByTweet(ctx, tweetIds...)
		if err != nil {
			return nil, err
		}
		for _, repoTweet := range repoTweets {
			tweets = append(tweets, repoTweet)
		}
		// put to cache
		if ctrl.cache != nil {
			for _, tweet := range repoTweets {
				err = putTweetIdToCache(ctrl.cache, tweet.TweetId, tweet)
				if err != nil {
					return nil, err
				}
			}
		}
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
	return tweetsMedia, nil
}

func (ctrl *Controller) RetrieveByUserID(ctx context.Context, userIds ...types.UserId) ([]model.Media, error) {
	var (
		checkDb = true
		tweets  []model.Tweet
	)

	// retrieve from cache
	if ctrl.cache != nil {
		checkDb, tweets, userIds = retrieveByUserIds(ctrl.cache, userIds...)
	}

	// retrieve all remaining tweets
	if checkDb {
		repoTweets, err := ctrl.repo.GetByUser(ctx, userIds...)
		if err != nil {
			return nil, err
		}
		for _, repoTweet := range repoTweets {
			tweets = append(tweets, repoTweet)
		}
		// put to cache
		if ctrl.cache != nil {
			for _, tweet := range repoTweets {
				err = putUserIdToCache(ctrl.cache, tweet.UserId, tweet)
				if err != nil {
					return nil, err
				}
			}
		}
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
	return tweetsMedia, nil
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
	// remove from cache
	if ctrl.cache != nil {
		tweetId := cachestorage.GenerateTweetId(postId)
		userTweetKey := cachestorage.GenerateUserToTweetId(tweetData[0].UserId, postId)
		err = ctrl.cache.Remove(tweetId)
		if err != nil {
			return err
		}
		err = ctrl.cache.Remove(userTweetKey)
		if err != nil {
			return err
		}
	}

	// delete from storage
	err = ctrl.storage.Delete(*tweetData[0].MediaUrl)
	return err
}
