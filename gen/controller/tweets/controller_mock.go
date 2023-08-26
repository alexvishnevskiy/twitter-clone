// Code generated by MockGen. DO NOT EDIT.
// Source: tweets/internal/controller/controller.go

// Package tweets is a generated GoMock package.
package tweets

import (
	context "context"
	reflect "reflect"
	time "time"

	types "github.com/alexvishnevskiy/twitter-clone/internal/types"
	model "github.com/alexvishnevskiy/twitter-clone/tweets/pkg/model"
	gomock "github.com/golang/mock/gomock"
)

// MocktweetsRepository is a mock of tweetsRepository interface.
type MocktweetsRepository struct {
	ctrl     *gomock.Controller
	recorder *MocktweetsRepositoryMockRecorder
}

// MocktweetsRepositoryMockRecorder is the mock recorder for MocktweetsRepository.
type MocktweetsRepositoryMockRecorder struct {
	mock *MocktweetsRepository
}

// NewMocktweetsRepository creates a new mock instance.
func NewMocktweetsRepository(ctrl *gomock.Controller) *MocktweetsRepository {
	mock := &MocktweetsRepository{ctrl: ctrl}
	mock.recorder = &MocktweetsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocktweetsRepository) EXPECT() *MocktweetsRepositoryMockRecorder {
	return m.recorder
}

// DeletePost mocks base method.
func (m *MocktweetsRepository) DeletePost(ctx context.Context, postId types.TweetId) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePost", ctx, postId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePost indicates an expected call of DeletePost.
func (mr *MocktweetsRepositoryMockRecorder) DeletePost(ctx, postId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePost", reflect.TypeOf((*MocktweetsRepository)(nil).DeletePost), ctx, postId)
}

// GetByTweet mocks base method.
func (m *MocktweetsRepository) GetByTweet(ctx context.Context, tweetIds ...types.TweetId) ([]model.Tweet, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range tweetIds {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetByTweet", varargs...)
	ret0, _ := ret[0].([]model.Tweet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByTweet indicates an expected call of GetByTweet.
func (mr *MocktweetsRepositoryMockRecorder) GetByTweet(ctx interface{}, tweetIds ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, tweetIds...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByTweet", reflect.TypeOf((*MocktweetsRepository)(nil).GetByTweet), varargs...)
}

// GetByUser mocks base method.
func (m *MocktweetsRepository) GetByUser(ctx context.Context, userIds ...types.UserId) ([]model.Tweet, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range userIds {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetByUser", varargs...)
	ret0, _ := ret[0].([]model.Tweet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUser indicates an expected call of GetByUser.
func (mr *MocktweetsRepositoryMockRecorder) GetByUser(ctx interface{}, userIds ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, userIds...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUser", reflect.TypeOf((*MocktweetsRepository)(nil).GetByUser), varargs...)
}

// Put mocks base method.
func (m *MocktweetsRepository) Put(ctx context.Context, userId types.UserId, content string, mediaUrl *string, retweetId *types.TweetId) (types.TweetId, time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", ctx, userId, content, mediaUrl, retweetId)
	ret0, _ := ret[0].(types.TweetId)
	ret1, _ := ret[1].(time.Time)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Put indicates an expected call of Put.
func (mr *MocktweetsRepositoryMockRecorder) Put(ctx, userId, content, mediaUrl, retweetId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MocktweetsRepository)(nil).Put), ctx, userId, content, mediaUrl, retweetId)
}
