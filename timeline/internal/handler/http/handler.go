package http

import (
	"encoding/json"
	"fmt"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
	_ "github.com/alexvishnevskiy/twitter-clone/timeline/docs"
	"github.com/alexvishnevskiy/twitter-clone/timeline/internal/controller"
	"github.com/alexvishnevskiy/twitter-clone/tweets/pkg/model"
	"net/http"
	"strconv"
)

type Hanlder struct {
	ctrl *controller.Controller
}

func New(ctrl *controller.Controller) *Hanlder {
	return &Hanlder{ctrl}
}

// GetHomeTimeline get all tweets from the users who this user is following
//
//	    @description    Retrieve home timeline
//		@Param			user_id query int true "User ID"
//		@Success		200	{object}	[]model.Media
//		@Failure		400	{object}	int
//		@Failure		404	{object}	int
//		@Failure		405	{object}	int
//		@Failure		500	{object}	int
//		@Router			/home_timeline [get]
func (h *Hanlder) GetHomeTimeline(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var tweets []model.Media

	// get user_id
	user := req.FormValue("user_id")
	user_id, err := strconv.Atoi(user)
	if err != nil {
		http.Error(w, fmt.Sprintf("user_id is invalid %s", user), http.StatusBadRequest)
		return
	}

	// retrieve timeline
	userId := types.UserId(user_id)
	tweets, err = h.ctrl.GetHomeTimeline(req.Context(), userId)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(tweets); err != nil {
		http.Error(w, "failed to encode tweets", http.StatusInternalServerError)
	}
}
