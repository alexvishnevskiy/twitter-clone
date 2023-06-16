package http

// how to handle api requests

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
	"github.com/alexvishnevskiy/twitter-clone/tweets/internal/controller"
	"github.com/alexvishnevskiy/twitter-clone/tweets/internal/repository/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	ctrl *controller.Controller
}

// Define the structure of the request body data.
type PostRequest struct {
	UserId    string  `json:"user_id"`
	Content   string  `json:"content"`
	MediaUrl  *string `json:"media_url"`  // optional
	RetweetId *string `json:"retweet_id"` // optional
}

func New(ctrl *controller.Controller) *Handler {
	return &Handler{ctrl}
}

// Retrieve either by tweet id or user id
func (h *Handler) Retrieve(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := req.ParseForm()
	if err != nil {
		// TODO: write more concise error to w
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, userOk := req.Form["user_id"]
	tweets, tweetsOk := req.Form["tweet_id"]

	if !userOk && !tweetsOk {
		// TODO: write more concise error to w
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// can retrieve both by user_id and tweet_id
	if userOk {
		userIds := make([]types.UserId, len(users))

		for i, user := range users {
			userId, err := strconv.Atoi(user)
			if err != nil {
				// TODO: write more concise error to w
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			userIds[i] = types.UserId(userId)
		}
		tweets, err := h.ctrl.RetrieveByUserID(req.Context(), userIds...)
		if err != nil && errors.Is(err, mysql.ErrNotFound) {
			// TODO: write more concise error to w
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err := json.NewEncoder(w).Encode(tweets); err != nil {
			http.Error(w, "failed to encode tweets", http.StatusInternalServerError)
		}
	}
	if tweetsOk {
		tweetIds := make([]types.TweetId, len(tweets))

		for i, tweet := range tweets {
			tweetId, err := strconv.Atoi(tweet)
			if err != nil {
				// TODO: write more concise error to w
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			tweetIds[i] = types.TweetId(tweetId)
		}
		tweets, err := h.ctrl.RetrieveByTweetID(req.Context(), tweetIds...)
		if err != nil && errors.Is(err, mysql.ErrNotFound) {
			// TODO: write more concise error to w
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err := json.NewEncoder(w).Encode(tweets); err != nil {
			http.Error(w, "failed to encode tweets", http.StatusInternalServerError)
		}
	}
}

// Delete by tweet id
func (h *Handler) Delete(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	tweet, err := strconv.Atoi(req.FormValue("tweet_id"))
	if err != nil {
		// TODO: write more concise error to w
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tweetID := types.TweetId(tweet)
	err = h.ctrl.DeletePost(req.Context(), tweetID)
	if err != nil {
		// TODO: write more concise error to w
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Failed to delete post: %v\n", err)
	}
}

func (h *Handler) Post(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var (
		requestData PostRequest
		retweetId   *types.TweetId = nil
	)

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

	if requestData.UserId == "" && requestData.Content == "" {
		http.Error(w, "user_id and content are empty", http.StatusBadRequest)
		return
	}

	userId, err := strconv.Atoi(requestData.UserId)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse user_id: %s", requestData.UserId), http.StatusBadRequest)
		return
	}

	if requestData.RetweetId != nil {
		tweetId, err := strconv.Atoi(*requestData.RetweetId)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to parse retweetId: %d", requestData.RetweetId), http.StatusBadRequest)
			return
		}
		tweet := types.TweetId(tweetId)
		retweetId = &tweet
	}
	// Assuming you have a matching function PostNewTweet in your ctrl:
	tweetId, err := h.ctrl.PostNewTweet(
		req.Context(),
		types.UserId(userId),
		requestData.Content,
		requestData.MediaUrl,
		retweetId,
	)
	if err != nil {
		http.Error(w, "failed to post tweet", http.StatusBadRequest)
	}
	if err := json.NewEncoder(w).Encode(tweetId); err != nil {
		http.Error(w, "response encode error", http.StatusInternalServerError)
		log.Printf("Response encode error: %v\n", err)
	}
}
