package grpc

import (
	"context"
	gen "github.com/alexvishnevskiy/twitter-clone/gen/api/follow"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
	"google.golang.org/grpc"
)

type Gateway struct {
	Url string
}

func New(url string) *Gateway {
	return &Gateway{url}
}

func (g *Gateway) GetUsers(ctx context.Context, userId types.UserId) ([]types.UserId, error) {
	// Set up the connection to the gRPC server.
	conn, err := grpc.Dial(g.Url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := gen.NewFollowServiceClient(conn)
	response, err := client.GetUserFollowers(ctx, &gen.UserId{UserId: int32(userId)})
	if err != nil {
		return nil, err
	}

	var users []types.UserId
	for _, protoUser := range response.UserId {
		user := types.UserId(protoUser.GetUserId())
		users = append(users, user)
	}
	return users, nil
}
