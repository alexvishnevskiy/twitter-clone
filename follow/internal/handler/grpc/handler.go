package grpc

import (
	"context"
	"github.com/alexvishnevskiy/twitter-clone/follow/internal/controller"
	gen "github.com/alexvishnevskiy/twitter-clone/gen/api/follow"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	gen.UnimplementedFollowServiceServer
	ctrl *controller.Controller
}

func New(ctrl *controller.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func usersToProto(users []types.UserId) []*gen.UserId {
	var protoUsers []*gen.UserId
	for _, user := range users {
		protoUser := gen.UserId{UserId: int32(user)}
		protoUsers = append(protoUsers, &protoUser)
	}
	return protoUsers
}

func (h *Handler) GetUserFollowers(ctx context.Context, req *gen.UserId) (*gen.GetResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil req")
	}

	users, err := h.ctrl.GetUserFollowers(ctx, types.UserId(req.UserId))
	if err != nil {
		return nil, err
	}

	protoUsers := usersToProto(users)
	return &gen.GetResponse{
		UserId: protoUsers,
	}, nil
}

func (h *Handler) GetFollowingUser(ctx context.Context, req *gen.UserId) (*gen.GetResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil req")
	}

	users, err := h.ctrl.GetFollowingUser(ctx, types.UserId(req.UserId))
	if err != nil {
		return nil, err
	}

	protoUsers := usersToProto(users)
	return &gen.GetResponse{
		UserId: protoUsers,
	}, nil
}
