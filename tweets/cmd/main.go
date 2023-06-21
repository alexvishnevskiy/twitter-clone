package main

import (
	"flag"
	"fmt"
	"github.com/alexvishnevskiy/twitter-clone/internal/storage/local"
	"github.com/alexvishnevskiy/twitter-clone/tweets/internal/controller"
	httphandler "github.com/alexvishnevskiy/twitter-clone/tweets/internal/handler/http"
	"github.com/alexvishnevskiy/twitter-clone/tweets/internal/repository/mysql"
	"log"
	"net/http"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "API handler port")
	flag.Parse()
	log.Printf("Starting the tweets service on port %d", port)

	// probably use configs for driverName and dataSourceName
	repository, err := mysql.New("mysql", "root:root@tcp(localhost:3306)/twitter")
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	ctrl := controller.New(repository)
	storage := local.New("/Users/alexander/Downloads/tweets/")
	h := httphandler.New(ctrl, storage)
	http.Handle("/post_tweet", http.HandlerFunc(h.Post))
	http.Handle("/retrieve_tweet", http.HandlerFunc(h.Retrieve))
	http.Handle("/delete_tweet", http.HandlerFunc(h.Delete))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
