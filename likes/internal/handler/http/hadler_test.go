package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	mock_controller "github.com/alexvishnevskiy/twitter-clone/gen/controller/likes"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
	"github.com/alexvishnevskiy/twitter-clone/likes/internal/controller"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Like(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// mock likes controller
	mockLikesRepo := mock_controller.NewMocklikesRepository(mockCtrl)
	mockLikesRepo.EXPECT().Like(ctx, types.UserId(1), types.TweetId(1)).Return(nil)
	tweetCtrl := controller.New(mockLikesRepo)
	likesHandler := New(tweetCtrl)

	// make json for body request
	payloadBytes, err := json.Marshal(
		struct {
			UserId  string `json:"user_id"`
			TweetId string `json:"tweet_id"`
		}{"1", "1"},
	)
	if err != nil {
		log.Fatalf("Failed to marshal payload: %v", err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "/like_tweet", body)
	if err != nil {
		t.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(likesHandler.Like)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf(
			"handler returned wrong status code",
		)
	}
}

func TestHandler_Unlike(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// mock likes controller
	mockLikesRepo := mock_controller.NewMocklikesRepository(mockCtrl)
	mockLikesRepo.EXPECT().Unlike(ctx, types.UserId(1), types.TweetId(1)).Return(nil)
	mockLikesRepo.EXPECT().Unlike(ctx, types.UserId(1), types.TweetId(3)).Return(errors.New(""))
	tweetCtrl := controller.New(mockLikesRepo)
	likesHandler := New(tweetCtrl)

	req, err := http.NewRequest("DELETE", "/unlike_tweet?user_id=1&tweet_id=1", nil)
	req1, err1 := http.NewRequest("DELETE", "/unlike_tweet?user_id=1&tweet_id=3", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err1 != nil {
		t.Fatal(err1)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(likesHandler.Unlike)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf(
			"handler returned wrong status code",
		)
	}

	handler.ServeHTTP(rr, req1)
	if status := rr.Code; status == http.StatusOK {
		t.Errorf(
			"handler returned wrong status code",
		)
	}
}

func TestHandler_GetTweetsByUser(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// mock tweet controller
	mockLikesRepo := mock_controller.NewMocklikesRepository(mockCtrl)
	tweetCtrl := controller.New(mockLikesRepo)
	likesHandler := New(tweetCtrl)

	want := []types.TweetId{
		1, 2, 3,
	}
	mockLikesRepo.EXPECT().GetTweetsByUser(ctx, types.UserId(1)).Return(want, nil)
	mockLikesRepo.EXPECT().GetTweetsByUser(ctx, types.UserId(2)).Return(nil, errors.New(""))

	req, err := http.NewRequest("GET", "/users_tweet?user_id=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req1, err1 := http.NewRequest("GET", "/users_tweet?user_id=2", nil)
	if err1 != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(likesHandler.GetTweetsByUser)
	handler.ServeHTTP(rr, req1)
	if status := rr.Code; status == http.StatusOK {
		t.Errorf(
			"handler returned wrong status code",
		)
	}

	handler.ServeHTTP(rr, req)
	var res []types.TweetId
	decoder := json.NewDecoder(rr.Body)
	err = decoder.Decode(&res)
	if err != nil {
		t.Errorf("failed to unmarshal result request")
	}
	if diff := cmp.Diff(want, res); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestHandler_GetUsersByTweet(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// mock tweet controller
	mockLikesRepo := mock_controller.NewMocklikesRepository(mockCtrl)
	tweetCtrl := controller.New(mockLikesRepo)
	likesHandler := New(tweetCtrl)

	want := []types.UserId{
		1, 2, 3,
	}
	mockLikesRepo.EXPECT().GetUsersByTweet(ctx, types.TweetId(1)).Return(want, nil)
	mockLikesRepo.EXPECT().GetUsersByTweet(ctx, types.TweetId(2)).Return(nil, errors.New(""))

	req, err := http.NewRequest("GET", "/tweets_user?tweet_id=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req1, err1 := http.NewRequest("GET", "/tweets_user?tweet_id=2", nil)
	if err1 != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(likesHandler.GetUsersByTweet)
	handler.ServeHTTP(rr, req1)
	if status := rr.Code; status == http.StatusOK {
		t.Errorf(
			"handler returned wrong status code",
		)
	}

	handler.ServeHTTP(rr, req)
	var res []types.UserId
	decoder := json.NewDecoder(rr.Body)
	err = decoder.Decode(&res)
	if err != nil {
		t.Errorf("failed to unmarshal result request")
	}
	if diff := cmp.Diff(want, res); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}
