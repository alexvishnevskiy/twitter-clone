package http

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/alexvishnevskiy/twitter-clone/follow/internal/controller"
	mock_controller "github.com/alexvishnevskiy/twitter-clone/gen/controller/follow"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Follow(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFollowRepo := mock_controller.NewMockfollowRepository(mockCtrl)
	mockFollowRepo.EXPECT().Follow(ctx, types.UserId(1), types.UserId(2)).Return(nil)
	followCtrl := controller.New(mockFollowRepo)
	followHandler := New(followCtrl)

	// make json for body request
	payloadBytes, err := json.Marshal(
		struct {
			UserId   string `json:"user_id"`
			FollowId string `json:"following_id"`
		}{"1", "2"},
	)
	if err != nil {
		log.Fatalf("Failed to marshal payload: %v", err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "/follow", body)
	if err != nil {
		t.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(followHandler.Follow)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf(
			"handler returned wrong status code",
		)
	}
}

func TestHandler_Unfollow(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFollowRepo := mock_controller.NewMockfollowRepository(mockCtrl)
	mockFollowRepo.EXPECT().Unfollow(ctx, types.UserId(1), types.UserId(2)).Return(nil)
	followCtrl := controller.New(mockFollowRepo)
	followHandler := New(followCtrl)

	req, err := http.NewRequest("DELETE", "/unfollow?user_id=1&following_id=2", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(followHandler.Unfollow)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf(
			"handler returned wrong status code",
		)
	}
}

func testGetFunction(t *testing.T, funcType int) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	var (
		rr  *httptest.ResponseRecorder
		res []types.UserId
	)

	want := []types.UserId{
		1, 2, 3,
	}

	mockFollowRepo := mock_controller.NewMockfollowRepository(mockCtrl)
	followCtrl := controller.New(mockFollowRepo)
	followHandler := New(followCtrl)

	switch funcType {
	case 0:
		mockFollowRepo.EXPECT().GetFollowingUser(ctx, types.UserId(1)).Return(want, nil)
		req, err := http.NewRequest("GET", "/user_followers?user_id=1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr = httptest.NewRecorder()
		handler := http.HandlerFunc(followHandler.GetFollowingUser)
		handler.ServeHTTP(rr, req)
	case 1:
		mockFollowRepo.EXPECT().GetUserFollowers(ctx, types.UserId(1)).Return(want, nil)
		req, err := http.NewRequest("GET", "/following_user?user_id=1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr = httptest.NewRecorder()
		handler := http.HandlerFunc(followHandler.GetUserFollowers)
		handler.ServeHTTP(rr, req)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf(
			"handler returned wrong status code",
		)
	}

	decoder := json.NewDecoder(rr.Body)
	err := decoder.Decode(&res)
	if err != nil {
		t.Errorf("failed to unmarshal result request")
	}
	if diff := cmp.Diff(want, res); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestHandler_GetFollowingUser(t *testing.T) {
	testGetFunction(t, 0)
}

func TestHandler_GetUserFollowers(t *testing.T) {
	testGetFunction(t, 1)
}
