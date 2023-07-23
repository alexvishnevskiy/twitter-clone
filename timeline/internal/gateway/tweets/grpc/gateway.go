package grpc

import (
	"context"
	gen "github.com/alexvishnevskiy/twitter-clone/gen/api/tweets"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
	"github.com/alexvishnevskiy/twitter-clone/tweets/pkg/model"
	"google.golang.org/grpc"
)

type Gateway struct {
	Url string
}

func New(url string) *Gateway {
	return &Gateway{url}
}

// get tweets from tweets service using user_ids
func (g *Gateway) GetTweets(ctx context.Context, userId ...types.UserId) ([]model.Media, error) {
	// Set up the connection to the gRPC server.
	conn, err := grpc.Dial(g.Url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// create tweets grpc service client
	client := gen.NewTweetsServiceClient(conn)
	users := make([]int32, len(userId))
	for i, user := range userId {
		users[i] = int32(user)
	}
	// retrieve tweets
	response, err := client.Retrieve(ctx, &gen.RetrieveRequest{UserId: users})
	if err != nil {
		return nil, err
	}

	tweets := make([]model.Media, len(response.MediaContent))
	for i, media := range response.MediaContent {
		tweets[i] = *model.MediaFromProto(media)
	}
	return tweets, nil
}
