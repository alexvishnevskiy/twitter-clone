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

func (h *Handler) Like(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var requestData PostRequest

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

func (h *Handler) Unlike(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

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
	err = h.ctrl.UnlikeTweet(req.Context(), userID, tweetID)
	if err != nil {
		http.Error(w, "failed to unlike a tweet", http.StatusInternalServerError)
	}
}

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
