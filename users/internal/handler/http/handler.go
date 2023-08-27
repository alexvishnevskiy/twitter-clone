package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alexvishnevskiy/twitter-clone/internal/jwt"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
	"github.com/alexvishnevskiy/twitter-clone/users/internal/controller"
	"github.com/alexvishnevskiy/twitter-clone/users/pkg/model"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Handler struct {
	ctrl *controller.Controller
}

// type for handler functions
type handlerMethod func(http.ResponseWriter, *http.Request)

// key for user_id field
type ctxKey struct{ name string }

var idCtxKey = &ctxKey{"user_id"}

func New(ctrl *controller.Controller) *Handler {
	return &Handler{ctrl}
}

// handler to generate jwt token
func JwtHandler(h handlerMethod) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		h(w, req)

		ctx := req.Context()
		userId, ok := ctx.Value(idCtxKey).(types.UserId)
		if !ok {
			http.Error(w, "failed to get id", http.StatusInternalServerError)
			return
		}

		token, err := jwt.GenerateJWT(userId)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to generate jwt: %s", err), http.StatusInternalServerError)
			return
		}

		http.SetCookie(
			w, &http.Cookie{
				Name:  "token",
				Value: token,
			},
		)
	}
}

// Register handle register method
//
//	@description	Register new user
//	@Param			password	body		string	true	"Password"
//	@Param			email		body		string	true	"Email"
//	@Param			nickname	body		string	true	"Nickname"
//	@Param			first_name	body		string	false	"First name"
//	@Param			last_name	body		string	false	"Last name"
//	@Success		200			{object}	int
//	@Failure		400			{object}	int
//	@Failure		404			{object}	int
//	@Failure		405			{object}	int
//	@Failure		500			{object}	int
//	@Router			/register       [post]
func (h *Handler) Register(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var requestData model.User

	// read all body data
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

	if requestData.Password == "" || requestData.Email == "" || requestData.Nickname == "" ||
		requestData.FirstName == "" || requestData.LastName == "" {
		http.Error(w, fmt.Sprintf("all fields should be present for User struct"), http.StatusBadRequest)
		return
	}

	// register user
	id, err := h.ctrl.Register(req.Context(), requestData)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to register user: %s", err), http.StatusInternalServerError)
	}

	// Set id in request context.
	ctx := context.WithValue(req.Context(), idCtxKey, id)
	*req = *req.WithContext(ctx)
}

// Login handle login method
//
//	@description	Login for user
//	@Param			password	body		string	true	"Password"
//	@Param			email		body		string	true	"Email"
//	@Success		200			{object}	int
//	@Failure		400			{object}	int
//	@Failure		404			{object}	int
//	@Failure		405			{object}	int
//	@Failure		500			{object}	int
//	@Router			/login       [post]
func (h *Handler) Login(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	bodyBytes, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// retrieve all data
	err = json.Unmarshal(bodyBytes, &requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId, err := h.ctrl.Login(req.Context(), requestData.Email, requestData.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid email or password: %s", err), http.StatusForbidden)
		return
	}

	// Set id in request context.
	ctx := context.WithValue(req.Context(), idCtxKey, userId)
	*req = *req.WithContext(ctx)
}

// Delete handle delete method
//
//	@description	Delete user
//	@Param			user_id	query		string	true	"User id"
//	@Success		200		{object}	int
//	@Failure		400		{object}	int
//	@Failure		404		{object}	int
//	@Failure		405		{object}	int
//	@Failure		500		{object}	int
//	@Router			/delete       [delete]
func (h *Handler) Delete(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userId, err := strconv.Atoi(req.FormValue("user_id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid user_id :%s", err), http.StatusBadRequest)
		return
	}

	err = h.ctrl.Delete(req.Context(), types.UserId(userId))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to delete: %s", err), http.StatusBadRequest)
	}
}

// Update handle update method
//
//	@description	Update user data
//	@Param			user_id		body		int		true	"User id"
//	@Param			nickname	body		string	false	"Nickname"
//	@Param			first_name	body		string	false	"First name"
//	@Param			last_name	body		string	false	"Last name"
//	@Param			email		body		string	false	"Email"
//	@Param			password	body		string	false	"Password"
//	@Success		200			{object}	int
//	@Failure		400			{object}	int
//	@Failure		404			{object}	int
//	@Failure		405			{object}	int
//	@Failure		500			{object}	int
//	@Router			/update       [put]
func (h *Handler) Update(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestData model.User

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

	err = h.ctrl.Update(req.Context(), requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
