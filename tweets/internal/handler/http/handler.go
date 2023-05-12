package http

// how to handle api requests

import (
	"encoding/json"
	"errors"
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

func (h *Handler) Retrieve(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, userOk := req.Form["user_id"]
	tweets, tweetsOk := req.Form["tweet_id"]

	if !userOk && !tweetsOk {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if userOk {
		userIds := make([]model.UserId, len(users))

		for i, user := range users {
			userId, err := strconv.Atoi(user)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			userIds[i] = model.UserId(userId)
		}
		tweets, err := h.ctrl.RetrieveByUserID(req.Context(), userIds...)
		if err != nil && errors.Is(err, mysql.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err := json.NewEncoder(w).Encode(tweets); err != nil {
			log.Printf("Response encode error: %v\n", err)
		}
	}
	if tweetsOk {
		tweetIds := make([]model.TweetId, len(tweets))

		for i, tweet := range tweets {
			tweetId, err := strconv.Atoi(tweet)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			tweetIds[i] = model.TweetId(tweetId)
		}
		tweets, err := h.ctrl.RetrieveByTweetID(req.Context(), tweetIds...)
		if err != nil && errors.Is(err, mysql.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err := json.NewEncoder(w).Encode(tweets); err != nil {
			log.Printf("Response encode error: %v\n", err)
		}
	}
}

func (h *Handler) Delete(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	tweet, err := strconv.Atoi(req.FormValue("tweet_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tweetID := model.TweetId(tweet)
	err = h.ctrl.DeletePost(req.Context(), tweetID)
	if err != nil {
		log.Printf("Failed to delete post: %v\n", err)
	}
}

func (h *Handler) Post(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var (
		userId    model.UserId
		content   string
		mediaUrl  *string
		retweetId *model.TweetId
	)

	user_id := req.FormValue("user_id")
	content = req.FormValue("content")
	media_url := req.FormValue("media_url")   // optional
	retweet_id := req.FormValue("retweet_id") //optional

	if user_id == "" && content == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := strconv.Atoi(user_id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId = model.UserId(user)

	if media_url == "" {
		mediaUrl = nil
	} else {
		mediaUrl = &media_url
	}

	if retweet, err := strconv.Atoi(retweet_id); retweet_id != "" && err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if retweet_id != "" {
		retweet_id := model.TweetId(retweet)
		retweetId = &retweet_id
	} else {
		retweetId = nil
	}

	_, err = h.ctrl.PostNewTweet(req.Context(), userId, content, mediaUrl, retweetId)
	if err != nil {
		log.Printf("Failed to post tweet: %v\n", err)
	}
}
