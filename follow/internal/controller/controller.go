package controller

import (
	"context"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
)

type followRepository interface {
	Follow(ctx context.Context, userId types.UserId, followId types.UserId) error
	Unfollow(ctx context.Context, userId types.UserId, followId types.UserId) error
	GetUserFollowers(ctx context.Context, userId types.UserId) ([]types.UserId, error)
	GetFollowingUser(ctx context.Context, userId types.UserId) ([]types.UserId, error)
}

type Controller struct {
	repo followRepository
}

func New(repo followRepository) *Controller {
	return &Controller{repo: repo}
}

func (ctrl *Controller) Follow(
	ctx context.Context,
	userId types.UserId,
	followId types.UserId,
) error {
	err := ctrl.repo.Follow(ctx, userId, followId)
	return err
}

func (ctrl *Controller) Unfollow(
	ctx context.Context,
	userId types.UserId,
	followId types.UserId,
) error {
	err := ctrl.repo.Unfollow(ctx, userId, followId)
	return err
}

func (ctrl *Controller) GetUserFollowers(
	ctx context.Context,
	userId types.UserId,
) ([]types.UserId, error) {
	followers, err := ctrl.repo.GetUserFollowers(ctx, userId)
	return followers, err
}

func (ctrl *Controller) GetFollowingUser(
	ctx context.Context,
	userId types.UserId,
) ([]types.UserId, error) {
	followers, err := ctrl.repo.GetFollowingUser(ctx, userId)
	return followers, err
}
