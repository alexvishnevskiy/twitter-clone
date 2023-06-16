// Code generated by MockGen. DO NOT EDIT.
// Source: likes/internal/controller/controller.go

// Package controller is a generated GoMock package.
package controller

import (
	context "context"
	reflect "reflect"

	types "github.com/alexvishnevskiy/twitter-clone/internal/types"
	gomock "github.com/golang/mock/gomock"
)

// MocklikesRepository is a mock of likesRepository interface.
type MocklikesRepository struct {
	ctrl     *gomock.Controller
	recorder *MocklikesRepositoryMockRecorder
}

// MocklikesRepositoryMockRecorder is the mock recorder for MocklikesRepository.
type MocklikesRepositoryMockRecorder struct {
	mock *MocklikesRepository
}

// NewMocklikesRepository creates a new mock instance.
func NewMocklikesRepository(ctrl *gomock.Controller) *MocklikesRepository {
	mock := &MocklikesRepository{ctrl: ctrl}
	mock.recorder = &MocklikesRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocklikesRepository) EXPECT() *MocklikesRepositoryMockRecorder {
	return m.recorder
}

// GetTweetsByUser mocks base method.
func (m *MocklikesRepository) GetTweetsByUser(ctx context.Context, userId types.UserId) ([]types.TweetId, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTweetsByUser", ctx, userId)
	ret0, _ := ret[0].([]types.TweetId)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTweetsByUser indicates an expected call of GetTweetsByUser.
func (mr *MocklikesRepositoryMockRecorder) GetTweetsByUser(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTweetsByUser", reflect.TypeOf((*MocklikesRepository)(nil).GetTweetsByUser), ctx, userId)
}

// GetUsersByTweet mocks base method.
func (m *MocklikesRepository) GetUsersByTweet(ctx context.Context, tweetId types.TweetId) ([]types.UserId, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersByTweet", ctx, tweetId)
	ret0, _ := ret[0].([]types.UserId)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersByTweet indicates an expected call of GetUsersByTweet.
func (mr *MocklikesRepositoryMockRecorder) GetUsersByTweet(ctx, tweetId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersByTweet", reflect.TypeOf((*MocklikesRepository)(nil).GetUsersByTweet), ctx, tweetId)
}

// Like mocks base method.
func (m *MocklikesRepository) Like(ctx context.Context, userId types.UserId, tweetId types.TweetId) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Like", ctx, userId, tweetId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Like indicates an expected call of Like.
func (mr *MocklikesRepositoryMockRecorder) Like(ctx, userId, tweetId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Like", reflect.TypeOf((*MocklikesRepository)(nil).Like), ctx, userId, tweetId)
}

// Unlike mocks base method.
func (m *MocklikesRepository) Unlike(ctx context.Context, userId types.UserId, tweetId types.TweetId) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unlike", ctx, userId, tweetId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unlike indicates an expected call of Unlike.
func (mr *MocklikesRepositoryMockRecorder) Unlike(ctx, userId, tweetId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unlike", reflect.TypeOf((*MocklikesRepository)(nil).Unlike), ctx, userId, tweetId)
}