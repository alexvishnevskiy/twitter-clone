package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	mockCache "github.com/alexvishnevskiy/twitter-clone/gen/cache"
	mockcontroller "github.com/alexvishnevskiy/twitter-clone/gen/controller/tweets"
	mockStorage "github.com/alexvishnevskiy/twitter-clone/gen/storage"
	localcache "github.com/alexvishnevskiy/twitter-clone/internal/cache/local"
	"github.com/alexvishnevskiy/twitter-clone/internal/storage/local"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
	"github.com/alexvishnevskiy/twitter-clone/tweets/internal/controller"
	"github.com/alexvishnevskiy/twitter-clone/tweets/pkg/model"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_Delete(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// mock tweet repo
	mockTweetRepo := mockcontroller.NewMocktweetsRepository(mockCtrl)
	mockTweetRepo.EXPECT().GetByTweet(ctx, types.TweetId(1)).Return(
		[]model.Tweet{
			{
				UserId:    types.UserId(1),
				TweetId:   types.TweetId(1),
				Content:   "",
				CreatedAt: time.Now(),
			},
		}, nil,
	)
	mockTweetRepo.EXPECT().GetByTweet(ctx, types.TweetId(2)).Return(
		[]model.Tweet{
			{
				UserId:    types.UserId(2),
				TweetId:   types.TweetId(2),
				Content:   "",
				CreatedAt: time.Now(),
			},
		}, nil,
	)
	mockTweetRepo.EXPECT().DeletePost(ctx, types.TweetId(1)).Return(nil)
	mockTweetRepo.EXPECT().DeletePost(ctx, types.TweetId(2)).Return(errors.New(""))

	// mock storage and cache
	mockstorage := mockStorage.NewMockStorage(mockCtrl)
	mockcache := mockCache.NewMockCache(mockCtrl)
	mockcache.EXPECT().Remove("user_id_1_tweet_id_1").Return(nil)
	mockcache.EXPECT().Remove("tweet_id_1").Return(nil)

	// tweet controller
	tweetCtrl := controller.New(mockTweetRepo, mockstorage, mockcache)
	tweetHandler := New(tweetCtrl)

	testCases := []struct {
		name    string
		tweetId int
		method  string
	}{
		{
			name:    "DELETE1",
			tweetId: 1,
			method:  "DELETE",
		},
		{
			name:    "DELETE2",
			tweetId: 2,
			method:  "DELETE",
		},
		{
			name:    "PUT1",
			tweetId: 1,
			method:  "PUT",
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				// Create a request to pass to our handler.
				req, err := http.NewRequest(tc.method, fmt.Sprintf("/delete_tweet?tweet_id=%d", tc.tweetId), nil)
				if err != nil {
					t.Fatal(err)
				}

				// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(tweetHandler.Delete)
				handler.ServeHTTP(rr, req)

				switch tc.method {
				case "PUT":
					if status := rr.Code; tc.method == "PUT" && status != http.StatusMethodNotAllowed {
						t.Errorf(
							"handler returned wrong status code: got %v want %v",
							status, http.StatusMethodNotAllowed,
						)
					}
				case "DELETE":
					if status := rr.Code; tc.tweetId == 1 && tc.method == "DELETE" && status != http.StatusOK {
						t.Errorf(
							"handler returned wrong status code: got %v want %v",
							status, http.StatusOK,
						)
					}
					if status := rr.Code; tc.tweetId != 1 && tc.method == "DELETE" && status != http.StatusInternalServerError {
						t.Errorf(
							"handler returned wrong status code: got %v want %v",
							status, http.StatusInternalServerError,
						)
					}
				}
			},
		)
	}
}

func TestHandler_Post(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// mock tweet controller
	mockTweetRepo := mockcontroller.NewMocktweetsRepository(mockCtrl)
	tweetCtrl := controller.New(mockTweetRepo, local.New("./"), localcache.New(10))
	tweetHandler := New(tweetCtrl)

	want := types.TweetId(1)
	// mock tweet controller
	mockTweetRepo.EXPECT().
		Put(ctx, types.UserId(1), "content", nil, nil).
		Return(want, time.Now(), nil)

	// make json for body request
	payloadBytes, err := json.Marshal(
		struct {
			UserId  types.UserId `json:"user_id"`
			Content string       `json:"content"`
		}{1, "content"},
	)
	if err != nil {
		log.Fatalf("Failed to marshal payload: %v", err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "/post_tweet", body)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tweetHandler.Post)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status == http.StatusMethodNotAllowed || status == http.StatusBadRequest {
		t.Errorf(
			"handler returned wrong status code got %v want %v", status, http.StatusOK,
		)
	}

	var res *types.TweetId
	decoder := json.NewDecoder(rr.Body)
	err = decoder.Decode(&res)
	if err != nil {
		t.Errorf("failed to unmarshal result request")
	}
	if diff := cmp.Diff(want, *res); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestHandler_Retrieve(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// mock tweet controller
	mockTweetRepo := mockcontroller.NewMocktweetsRepository(mockCtrl)
	tweetCtrl := controller.New(mockTweetRepo, local.New("./"), localcache.New(10))
	tweetHandler := New(tweetCtrl)

	// expected output
	timeNow := time.Now()
	wantRepo := []model.Tweet{
		{
			UserId:    types.UserId(1),
			TweetId:   types.TweetId(1),
			Content:   "content",
			CreatedAt: timeNow,
		},
		{
			UserId:    types.UserId(2),
			TweetId:   types.TweetId(2),
			Content:   "content",
			CreatedAt: timeNow,
		},
	}
	wantHandler := []model.Media{
		{
			Media:     "",
			Content:   "content",
			CreatedAt: timeNow,
		},
		{
			Media:     "",
			Content:   "content",
			CreatedAt: timeNow,
		},
	}

	// expected behaviour
	mockTweetRepo.EXPECT().GetByUser(ctx, types.UserId(1), types.UserId(2)).Return(wantRepo, nil)
	mockTweetRepo.EXPECT().GetByTweet(ctx, types.TweetId(1), types.TweetId(2)).Return(wantRepo, nil)

	testCases := []struct {
		name     string
		user_id  []int
		tweet_id []int
	}{
		{
			name:     "getByUser",
			user_id:  []int{1, 2},
			tweet_id: []int{},
		},
		{
			name:     "getByTweet",
			user_id:  []int{},
			tweet_id: []int{1, 2},
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				var body string
				if len(tc.user_id) > 0 {
					body = fmt.Sprintf("/retrieve_tweet?user_id=%d&user_id=%d", tc.user_id[0], tc.user_id[1])
				} else {
					body = fmt.Sprintf("/retrieve_tweet?tweet_id=%d&tweet_id=%d", tc.tweet_id[0], tc.tweet_id[1])
				}

				req, err := http.NewRequest("GET", body, nil)
				if err != nil {
					t.Fatal(err)
				}
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(tweetHandler.Retrieve)
				handler.ServeHTTP(rr, req)

				if status := rr.Code; status == http.StatusMethodNotAllowed || status == http.StatusBadRequest {
					t.Errorf(
						"handler returned wrong status code",
					)
				}

				// check output
				var res []model.Media
				decoder := json.NewDecoder(rr.Body)
				err = decoder.Decode(&res)
				if err != nil {
					t.Errorf("failed to unmarshal result request")
				}
				if diff := cmp.Diff(wantHandler, res); diff != "" {
					t.Errorf("mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}
