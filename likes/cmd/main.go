package main

import (
	"flag"
	"fmt"
	"github.com/alexvishnevskiy/twitter-clone/likes/internal/controller"
	httphandler "github.com/alexvishnevskiy/twitter-clone/likes/internal/handler/http"
	"github.com/alexvishnevskiy/twitter-clone/likes/internal/repository/mysql"
	"log"
	"net/http"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8081, "API handler port")
	flag.Parse()
	log.Printf("Starting the tweets service on port %d", port)

	// probably use configs for driverName and dataSourceName
	repository, err := mysql.New("mysql", "root:root@tcp(localhost:3306)/twitter")
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	ctrl := controller.New(repository)
	h := httphandler.New(ctrl)

	http.Handle("/like_tweet", http.HandlerFunc(h.Like))
	http.Handle("/unlike_tweet", http.HandlerFunc(h.Unlike))
	http.Handle("/users_tweet", http.HandlerFunc(h.GetUsersByTweet))
	http.Handle("/tweets_user", http.HandlerFunc(h.GetTweetsByUser))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
