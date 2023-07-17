package http

import (
	"encoding/json"
	"fmt"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
	"github.com/alexvishnevskiy/twitter-clone/timeline/internal/controller"
	"net/http"
	"strconv"
)

type Hanlder struct {
	ctrl *controller.Controller
}

func New(ctrl *controller.Controller) *Hanlder {
	return &Hanlder{ctrl}
}

// get all tweets from the users who this user is following
func (h *Hanlder) GetHomeTimeline(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	user := req.FormValue("user_id")
	user_id, err := strconv.Atoi(user)
	if err != nil {
		http.Error(w, fmt.Sprintf("user_id is invalid %s", user), http.StatusBadRequest)
		return
	}

	userId := types.UserId(user_id)
	tweets, err := h.ctrl.GetHomeTimeline(req.Context(), userId)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(tweets); err != nil {
		http.Error(w, "failed to encode tweets", http.StatusInternalServerError)
	}
}
