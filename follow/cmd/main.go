package main

import (
	"flag"
	"fmt"
	"github.com/alexvishnevskiy/twitter-clone/follow/internal/controller"
	grpchandler "github.com/alexvishnevskiy/twitter-clone/follow/internal/handler/grpc"
	"github.com/alexvishnevskiy/twitter-clone/follow/internal/repository/mysql"
	gen "github.com/alexvishnevskiy/twitter-clone/gen/api/follow"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8082, "API handler port")
	flag.Parse()
	log.Printf("Starting the tweets service on port %d", port)

	// probably use configs for driverName and dataSourceName
	repository, err := mysql.New("mysql", "root:root@tcp(localhost:3306)/twitter")
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	ctrl := controller.New(repository)
	h := grpchandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterFollowServiceServer(srv, h)
	srv.Serve(lis)

	//h := httphandler.New(ctrl)
	//
	//http.Handle("/follow", http.HandlerFunc(h.Follow))
	//http.Handle("/unfollow", http.HandlerFunc(h.Unfollow))
	//http.Handle("/user_followers", http.HandlerFunc(h.GetUserFollowers))
	//http.Handle("/following_user", http.HandlerFunc(h.GetFollowingUser))
	//
	//if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
	//	panic(err)
	//}
}
