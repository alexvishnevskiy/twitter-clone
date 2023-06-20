package main

import (
	"flag"
	"fmt"
	"github.com/alexvishnevskiy/twitter-clone/timeline/internal/controller"
	followGateway "github.com/alexvishnevskiy/twitter-clone/timeline/internal/gateway/follow/http"
	tweetsGateway "github.com/alexvishnevskiy/twitter-clone/timeline/internal/gateway/tweets/http"
	httphandler "github.com/alexvishnevskiy/twitter-clone/timeline/internal/handler/http"
	"log"
	"net/http"
)

func main() {
	var (
		port        int
		tweets_port int
		follow_port int
	)
	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.IntVar(&tweets_port, "tweets_port", 8080, "tweets API handler port")
	flag.IntVar(&follow_port, "follow_port", 8082, "follow API handler port")
	flag.Parse()
	log.Printf("Starting timeline service on port %d", port)

	tweetsService := tweetsGateway.New(fmt.Sprintf("http://localhost:%d", tweets_port))
	followService := followGateway.New(fmt.Sprintf("http://localhost:%d", follow_port))
	ctrl := controller.New(tweetsService, followService)
	h := httphandler.New(ctrl)

	http.Handle("/home_timeline", http.HandlerFunc(h.GetHomeTimeline))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
