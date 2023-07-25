package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexvishnevskiy/twitter-clone/follow/internal/controller"
	"github.com/alexvishnevskiy/twitter-clone/follow/internal/repository/mysql"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
	"io/ioutil"
	"net/http"
	"strconv"
)

// handler to handle http requests
type Handler struct {
	ctrl *controller.Controller
}

type retrieveFunc func(context.Context, types.UserId) ([]types.UserId, error)

type PostRequest struct {
	UserId   string `json:"user_id"`
	FollowId string `json:"following_id"`
}

func New(ctrl *controller.Controller) *Handler {
	return &Handler{ctrl}
}

// Follow handle follow requests
//
//	@description	Follow specific user
//	@Param			user_id			body		int	true	"User ID"
//	@Param			following_id	body		int	true	"Following ID"
//	@Success		200				{object}	int
//	@Failure		400				{object}	int
//	@Failure		404				{object}	int
//	@Failure		405				{object}	int
//	@Failure		500				{object}	int
//	@Router			/follow [post]
func (h *Handler) Follow(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var requestData PostRequest

	// read and unmarshal data from request
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

	// user_id and follow_id are non nil fields
	if requestData.UserId == "" && requestData.FollowId == "" {
		http.Error(w, "user_id and following_id are empty", http.StatusBadRequest)
		return
	}

	user, err := strconv.Atoi(requestData.UserId)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse user_id: %s", requestData.UserId), http.StatusBadRequest)
		return
	}

	follower, err := strconv.Atoi(requestData.FollowId)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse following_id: %s", requestData.FollowId), http.StatusBadRequest)
		return
	}

	userID := types.UserId(user)
	followerID := types.UserId(follower)

	// user controller to make request
	err = h.ctrl.Follow(req.Context(), userID, followerID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to follow a tweet: %s", err), http.StatusInternalServerError)
	}
}

// Unfollow handle unfollow requests
//
//	@description	Unfollow specific user
//	@Param			user_id			query		int	true	"User ID"
//	@Param			following_id	query		int	true	"Following ID"
//	@Success		200				{object}	int
//	@Failure		400				{object}	int
//	@Failure		404				{object}	int
//	@Failure		405				{object}	int
//	@Failure		500				{object}	int
//	@Router			/unfollow [delete]
func (h *Handler) Unfollow(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// retrieve user_id and following_id
	user_id := req.FormValue("user_id")
	userid, err := strconv.Atoi(user_id)
	if err != nil {
		http.Error(w, fmt.Sprintf("user_id is invalid: %s", user_id), http.StatusBadRequest)
		return
	}

	following_id := req.FormValue("following_id")
	followingid, err := strconv.Atoi(following_id)
	if err != nil {
		http.Error(w, fmt.Sprintf("following_id is invalid: %s", following_id), http.StatusBadRequest)
		return
	}

	// handle request with controller
	userId := types.UserId(userid)
	followingId := types.UserId(followingid)
	err = h.ctrl.Unfollow(req.Context(), userId, followingId)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to follow a tweet: %s", err), http.StatusInternalServerError)
	}
}

// helper function to get users with function f
func getUsers(w http.ResponseWriter, req *http.Request, f retrieveFunc) {
	if req.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// get user_id
	user_id := req.FormValue("user_id")
	userid, err := strconv.Atoi(user_id)
	if err != nil {
		http.Error(w, fmt.Sprintf("user_id is invalid: %s", user_id), http.StatusBadRequest)
		return
	}

	userId := types.UserId(userid)
	// retrieve users with function f
	users, err := f(req.Context(), userId)
	if err != nil && !errors.Is(err, mysql.ErrNotFound) {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}
	if err != nil && errors.Is(err, mysql.ErrNotFound) {
		http.Error(w, "failed to find followers by this user_id", http.StatusNotFound)
		return
	}
	// encode users to json
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "failed to encode followers", http.StatusInternalServerError)
	}
}

// GetUserFollowers get user followers
//
//	@description	Retrieve all user followers
//	@Param			user_id	query		int	true	"User ID"
//	@Success		200		{object}	[]types.UserId
//	@Failure		400		{object}	int
//	@Failure		404		{object}	int
//	@Failure		405		{object}	int
//	@Failure		500		{object}	int
//	@Router			/user_followers [get]
func (h *Handler) GetUserFollowers(w http.ResponseWriter, req *http.Request) {
	getUsers(w, req, h.ctrl.GetUserFollowers)
}

// GetFollowingUser get users who are following user
//
//	@description	Retrieve all following user
//	@Param			user_id	query		int	true	"User ID"
//	@Success		200		{object}	[]types.UserId
//	@Failure		400		{object}	int
//	@Failure		404		{object}	int
//	@Failure		405		{object}	int
//	@Failure		500		{object}	int
//	@Router			/following_user [get]
func (h *Handler) GetFollowingUser(w http.ResponseWriter, req *http.Request) {
	getUsers(w, req, h.ctrl.GetFollowingUser)
}
