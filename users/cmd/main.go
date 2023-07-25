package main

import (
	"flag"
	"fmt"
	"github.com/alexvishnevskiy/twitter-clone/internal/jwt"
	_ "github.com/alexvishnevskiy/twitter-clone/users/docs"
	"github.com/alexvishnevskiy/twitter-clone/users/internal/controller"
	httphandler "github.com/alexvishnevskiy/twitter-clone/users/internal/handler/http"
	"github.com/alexvishnevskiy/twitter-clone/users/internal/repository/mysql"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

// @title Users API documentation
// @version 1.0.0
// @host localhost:8084
// @description This is API for users service
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

	updateHandler := jwt.ValidateMiddleware(http.HandlerFunc(h.Update))
	deleteHandler := jwt.ValidateMiddleware(http.HandlerFunc(h.Delete))
	loginHandler := http.HandlerFunc(httphandler.JwtHandler(h.Login))
	registerHandler := http.HandlerFunc(httphandler.JwtHandler(h.Register))

	http.Handle("/login", loginHandler)
	http.Handle("/register", registerHandler)
	http.Handle("/update", updateHandler)
	http.Handle("/delete", deleteHandler)
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
