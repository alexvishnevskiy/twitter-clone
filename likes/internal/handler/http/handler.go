package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
	"github.com/alexvishnevskiy/twitter-clone/likes/internal/controller"
	"github.com/alexvishnevskiy/twitter-clone/likes/internal/repository/mysql"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Handler struct {
	ctrl *controller.Controller
}

// Define the structure of the request body data.
type PostRequest struct {
	UserId  string `json:"user_id"`
	TweetId string `json:"tweet_id"`
}

func New(ctrl *controller.Controller) *Handler {
	return &Handler{ctrl}
}

// Like handle like request
//
//	@description	Like specific tweet
//	@Param			user_id		body		int	true	"User ID"
//	@Param			tweet_id	body		int	true	"Tweet ID"
//	@Success		200			{object}	int
//	@Failure		400			{object}	int
//	@Failure		404			{object}	int
//	@Failure		405			{object}	int
//	@Failure		500			{object}	int
//	@Router			/like_tweet [post]
func (h *Handler) Like(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var requestData PostRequest

	// read all request data
	bodyBytes, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(bodyBytes, &requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if requestData.UserId == "" && requestData.TweetId == "" {
		http.Error(w, "user_id and tweet_id are empty", http.StatusBadRequest)
		return
	}

	// preprocess data
	tweet, err := strconv.Atoi(requestData.TweetId)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse tweet_id: %s", requestData.TweetId), http.StatusBadRequest)
		return
	}

	user, err := strconv.Atoi(requestData.UserId)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse user_id: %s", requestData.UserId), http.StatusBadRequest)
		return
	}

	tweetID := types.TweetId(tweet)
	userID := types.UserId(user)
	err = h.ctrl.LikeTweet(req.Context(), userID, tweetID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to like a tweet: %s", err), http.StatusInternalServerError)
	}
}

// Unlike handle unlike request
//
//	@description	Unlike specific tweet
//	@Param			user_id		body		int	true	"User ID"
//	@Param			tweet_id	body		int	true	"Tweet ID"
//	@Success		200			{object}	int
//	@Failure		400			{object}	int
//	@Failure		404			{object}	int
//	@Failure		405			{object}	int
//	@Failure		500			{object}	int
//	@Router			/unlike_tweet [delete]
func (h *Handler) Unlike(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// read body data
	req_tweet := req.FormValue("tweet_id")
	tweet, err := strconv.Atoi(req_tweet)
	if err != nil {
		http.Error(w, fmt.Sprintf("tweet_id is invalid: %s", req_tweet), http.StatusBadRequest)
		return
	}

	req_user := req.FormValue("user_id")
	user, err := strconv.Atoi(req_user)
	if err != nil {
		http.Error(w, fmt.Sprintf("user_id is invalid: %s", req_user), http.StatusBadRequest)
		return
	}

	tweetID := types.TweetId(tweet)
	userID := types.UserId(user)
	// make request to controller
	err = h.ctrl.UnlikeTweet(req.Context(), userID, tweetID)
	if err != nil {
		http.Error(w, "failed to unlike a tweet", http.StatusInternalServerError)
	}
}

// GetUsersByTweet all users who like tweet
//
//	@description	Retrieve all users who liked tweet
//	@Param			tweet_id	query		int	true	"Tweet ID"
//	@Success		200			{object}	[]types.UserId
//	@Failure		400			{object}	int
//	@Failure		404			{object}	int
//	@Failure		405			{object}	int
//	@Failure		500			{object}	int
//	@Router			/users_tweet [get]
func (h *Handler) GetUsersByTweet(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	req_tweet := req.FormValue("tweet_id")
	tweet, err := strconv.Atoi(req_tweet)
	if err != nil {
		http.Error(w, fmt.Sprintf("tweet_id is invalid: %s", req_tweet), http.StatusBadRequest)
		return
	}

	tweetID := types.TweetId(tweet)
	users, err := h.ctrl.GetUsersByTweet(req.Context(), tweetID)

	if err != nil && !errors.Is(err, mysql.ErrNotFound) {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}
	if err != nil && errors.Is(err, mysql.ErrNotFound) {
		http.Error(w, "failed to find users by this tweet_id", http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "failed to encode users", http.StatusInternalServerError)
	}
}

// GetTweetsByUser get all tweets liked by user
//
//	@description	Retrieve all tweet liked by user
//	@Param			user_id	query		int	true	"User ID"
//	@Success		200		{object}	[]types.TweetId
//	@Failure		400		{object}	int
//	@Failure		404		{object}	int
//	@Failure		405		{object}	int
//	@Failure		500		{object}	int
//	@Router			/tweets_user    [get]
func (h *Handler) GetTweetsByUser(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	req_user := req.FormValue("user_id")
	user, err := strconv.Atoi(req_user)
	if err != nil {
		http.Error(w, fmt.Sprintf("user_id is invalidL %s", req_user), http.StatusBadRequest)
		return
	}

	userID := types.UserId(user)
	tweets, err := h.ctrl.GetTweetsByUser(req.Context(), userID)

	if err != nil && !errors.Is(err, mysql.ErrNotFound) {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}
	if err != nil && errors.Is(err, mysql.ErrNotFound) {
		http.Error(w, "failed to find tweets by this user_id", http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(tweets); err != nil {
		http.Error(w, "failed to encode tweets", http.StatusInternalServerError)
	}
}
