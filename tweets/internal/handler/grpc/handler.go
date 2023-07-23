package grpc

import (
	"context"
	gen "github.com/alexvishnevskiy/twitter-clone/gen/api/tweets"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
	"github.com/alexvishnevskiy/twitter-clone/tweets/internal/controller"
	"github.com/alexvishnevskiy/twitter-clone/tweets/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	gen.UnimplementedTweetsServiceServer
	ctrl *controller.Controller
}

func New(ctrl *controller.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// Retrieve tweet either by tweet_id or user_id
func (h *Handler) Retrieve(ctx context.Context, req *gen.RetrieveRequest) (*gen.RetrieveResponse, error) {
	if req == nil || (req.TweetId == nil && req.UserId == nil) {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	var (
		tweetsData []model.Media
		err        error
	)

	if req.TweetId != nil {
		tweetIds := make([]types.TweetId, len(req.TweetId))
		for i, id := range req.TweetId {
			tweetIds[i] = types.TweetId(id)
		}
		tweetsData, err = h.ctrl.RetrieveByTweetID(ctx, tweetIds...)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}
	if req.UserId != nil {
		userIds := make([]types.UserId, len(req.UserId))
		for i, id := range req.UserId {
			userIds[i] = types.UserId(id)
		}
		tweetsData, err = h.ctrl.RetrieveByUserID(ctx, userIds...)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}

	if tweetsData == nil {
		return nil, status.Errorf(codes.InvalidArgument, "something wrong with either user_id or tweet_id")
	}

	// convert to proto and write response
	var protoResponse []*gen.Media
	for _, id := range tweetsData {
		mediaProto := model.MediaToProto(&id)
		protoResponse = append(protoResponse, mediaProto)
	}

	return &gen.RetrieveResponse{
		MediaContent: protoResponse,
	}, nil
}
