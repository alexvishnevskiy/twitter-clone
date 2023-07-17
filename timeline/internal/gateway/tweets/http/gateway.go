package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alexvishnevskiy/twitter-clone/internal/types"
	"github.com/alexvishnevskiy/twitter-clone/tweets/pkg/model"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Gateway struct {
	Url string
}

func New(url string) *Gateway {
	return &Gateway{url}
}

// get tweets from tweets service
func (g *Gateway) GetTweets(ctx context.Context, userId ...types.UserId) ([]model.Tweet, error) {
	base, _ := url.Parse(g.Url)
	newURL, _ := url.Parse(path.Join(base.Path, "/retrieve_tweet"))
	base = base.ResolveReference(newURL)
	g.Url = base.String()

	req, err := http.NewRequest(http.MethodGet, g.Url, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	for _, user := range userId {
		values.Add("user_id", strconv.Itoa(int(user)))
	}
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx response: %v", resp)
	}
	var tweets []model.Tweet
	if err := json.NewDecoder(resp.Body).Decode(&tweets); err != nil {
		return nil, err
	}
	return tweets, nil
}
