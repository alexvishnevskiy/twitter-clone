package main

import (
	"flag"
	"fmt"
	"github.com/alexvishnevskiy/twitter-clone/users/internal/controller"
	httphandler "github.com/alexvishnevskiy/twitter-clone/users/internal/handler/http"
	"github.com/alexvishnevskiy/twitter-clone/users/internal/repository/mysql"
	"log"
	"net/http"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8084, "API handler port")
	flag.Parse()
	log.Printf("Starting users service on port %d", port)

	repo, err := mysql.New("mysql", "root:root@tcp(localhost:3306)/twitter")
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	ctrl := controller.New(repo)
	h := httphandler.New(ctrl)

	http.Handle("/login", http.HandlerFunc(h.Login))
	http.Handle("/update", http.HandlerFunc(h.Update))
	http.Handle("/register", http.HandlerFunc(h.Register))
	http.Handle("/delete", http.HandlerFunc(h.Delete))
	
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
