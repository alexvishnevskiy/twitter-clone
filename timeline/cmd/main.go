package main

import (
	"flag"
	"fmt"
	"github.com/alexvishnevskiy/twitter-clone/timeline/internal/controller"
	followGateway "github.com/alexvishnevskiy/twitter-clone/timeline/internal/gateway/follow/grpc"
	tweetsGateway "github.com/alexvishnevskiy/twitter-clone/timeline/internal/gateway/tweets/grpc"
	httphandler "github.com/alexvishnevskiy/twitter-clone/timeline/internal/handler/http"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

//	@title			Timeline API documentation
//	@version		1.0.0
//	@host			localhost:8083
//	@description	This is API for timeline service
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

	tweetsService := tweetsGateway.New(fmt.Sprintf("localhost:%d", tweets_port))
	followService := followGateway.New(fmt.Sprintf("localhost:%d", follow_port))
	ctrl := controller.New(tweetsService, followService)
	h := httphandler.New(ctrl)

	http.Handle("/home_timeline", http.HandlerFunc(h.GetHomeTimeline))
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
