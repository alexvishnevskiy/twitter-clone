package http

// how to handle api requests

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
	"github.com/alexvishnevskiy/twitter-clone/tweets/internal/controller"
	"github.com/alexvishnevskiy/twitter-clone/tweets/internal/repository/mysql"
	"github.com/alexvishnevskiy/twitter-clone/tweets/pkg/model"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	ctrl *controller.Controller
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
		//TODO: write more concise error to w
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var tweetsData []model.Media
	users, userOk := req.Form["user_id"]
	tweets, tweetsOk := req.Form["tweet_id"]

	if !userOk && !tweetsOk {
		//TODO: write more concise error to w
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// can retrieve both by user_id and tweet_id
	if userOk {
		userIds := make([]types.UserId, len(users))

		for i, user := range users {
			userId, err := strconv.Atoi(user)
			if err != nil {
				//TODO: write more concise error to w
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			userIds[i] = types.UserId(userId)
		}
		tweetsData, err = h.ctrl.RetrieveByUserID(req.Context(), userIds...)
		if err != nil && errors.Is(err, mysql.ErrNotFound) {
			//TODO: write more concise error to w
			w.WriteHeader(http.StatusNotFound)
			return
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
		tweetsData, err = h.ctrl.RetrieveByTweetID(req.Context(), tweetIds...)
		if err != nil && errors.Is(err, mysql.ErrNotFound) {
			// TODO: write more concise error to w
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	jsonData, err := json.Marshal(tweetsData)
	if err != nil {
		http.Error(w, "Could not convert data to JSON", http.StatusInternalServerError)
		return
	}
	// Set the content type header.
	w.Header().Set("Content-Type", "application/json")
	// Write the JSON data to the response.
	w.Write(jsonData)
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
	// 1 << 16 is the maximum size you can read from the request
	req.ParseMultipartForm(1 << 16)

	var retweetId *types.TweetId = nil

	UserId := req.FormValue("user_id")
	Content := req.FormValue("content")
	RetweetId := req.FormValue("retweet_id")

	if UserId == "" && Content == "" {
		http.Error(w, "user_id and content are empty", http.StatusBadRequest)
		return
	}

	userId, err := strconv.Atoi(UserId)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse user_id: %s", UserId), http.StatusBadRequest)
		return
	}

	if RetweetId != "" {
		tweetId, err := strconv.Atoi(RetweetId)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to parse retweetId: %d", RetweetId), http.StatusBadRequest)
			return
		}
		tweet := types.TweetId(tweetId)
		retweetId = &tweet
	}

	file, handler, err := req.FormFile("media")
	if handler != nil {
		defer file.Close()
	}

	tweetId, err := h.ctrl.PostNewTweet(
		req.Context(),
		file,
		handler,
		types.UserId(userId),
		Content,
		retweetId,
	)

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to post tweet: %s", err), http.StatusBadRequest)
	}
	if err := json.NewEncoder(w).Encode(tweetId); err != nil {
		http.Error(w, "response encode error", http.StatusInternalServerError)
		log.Printf("Response encode error: %v\n", err)
	}
}
