package controller

import (
	"context"
	"fmt"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
	"github.com/alexvishnevskiy/twitter-clone/users/pkg/model"
	"github.com/stretchr/objx"
)

type usersRepository interface {
	Register(
		ctx context.Context,
		nickname string,
		firstname string,
		lastname string,
		email string,
		password string) (types.UserId, error)
	RetrievePassword(
		ctx context.Context,
		email string,
	) (types.UserId, string, error)
	Delete(
		ctx context.Context,
		userid types.UserId,
	) error
	Update(
		ctx context.Context,
		userData model.User,
	) error
}

type Controller struct {
	repo usersRepository
}

func New(repo usersRepository) *Controller {
	return &Controller{repo}
}

func encodePassword(password string) string {
	// TODO: replace with decode algorithm
	return objx.HashWithKey(password, "password")[:5]
}

func checkPassword(enteredPassword string, databasePassword string) bool {
	// TODO: add logic to check encrypted password and entered password
	return enteredPassword == databasePassword
}

func (ctrl *Controller) Register(
	ctx context.Context,
	userData model.User,
) (types.UserId, error) {
	decodedPassword := encodePassword(userData.Password)
	id, err := ctrl.repo.Register(
		ctx, userData.Nickname, userData.FirstName, userData.LastName, userData.Email, decodedPassword,
	)
	return id, err
}

func (ctrl *Controller) Login(ctx context.Context, email string, password string) (types.UserId, error) {
	userId, databasePassword, err := ctrl.repo.RetrievePassword(ctx, email)
	if err != nil {
		return types.UserId(0), err
	}

	check := checkPassword(password, databasePassword)
	if check {
		return userId, nil
	}
	return types.UserId(0), fmt.Errorf("password is incorrect")
}

func (ctrl *Controller) Delete(ctx context.Context, userid types.UserId) error {
	err := ctrl.repo.Delete(ctx, userid)
	return err
}

func (ctrl *Controller) Update(ctx context.Context, userData model.User) error {
	if userData.Password != "" {
		userData.Password = encodePassword(userData.Password)
	}
	err := ctrl.repo.Update(ctx, userData)
	return err
}
