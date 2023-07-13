package main

import (
	"flag"
	"fmt"
	localcache "github.com/alexvishnevskiy/twitter-clone/internal/cache/local"
	"github.com/alexvishnevskiy/twitter-clone/internal/storage/local"
	"github.com/alexvishnevskiy/twitter-clone/tweets/internal/controller"
	httphandler "github.com/alexvishnevskiy/twitter-clone/tweets/internal/handler/http"
	"github.com/alexvishnevskiy/twitter-clone/tweets/internal/repository/mysql"
	"log"
	"net/http"
)

func main() {
	var (
		port     int
		capacity int
	)
	flag.IntVar(&port, "port", 8080, "API handler port")
	flag.IntVar(&capacity, "capacity", 5000, "Capacity of cache")
	flag.Parse()
	log.Printf("Starting the tweets service on port %d", port)

	// probably use configs for driverName and dataSourceName
	repository, err := mysql.New("mysql", "root:root@tcp(localhost:3306)/twitter")
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	storage := local.New("/Users/alexander/Downloads/tweets/")
	cache := localcache.New(capacity)
	ctrl := controller.New(repository, storage, cache)
	h := httphandler.New(ctrl)
	http.Handle("/post_tweet", http.HandlerFunc(h.Post))
	http.Handle("/retrieve_tweet", http.HandlerFunc(h.Retrieve))
	http.Handle("/delete_tweet", http.HandlerFunc(h.Delete))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
