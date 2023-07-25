package main

import (
	"flag"
	"fmt"
	_ "github.com/alexvishnevskiy/twitter-clone/follow/docs"
	"github.com/alexvishnevskiy/twitter-clone/follow/internal/controller"
	grpchandler "github.com/alexvishnevskiy/twitter-clone/follow/internal/handler/grpc"
	httphandler "github.com/alexvishnevskiy/twitter-clone/follow/internal/handler/http"
	"github.com/alexvishnevskiy/twitter-clone/follow/internal/repository/mysql"
	gen "github.com/alexvishnevskiy/twitter-clone/gen/api/follow"
	"github.com/soheilhy/cmux"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
)

// @title			Follow API documentation
// @version		1.0.0
// @host			localhost:8082
// @description	This is API for follow service
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

	// setup the main listener
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new cmux instance.
	m := cmux.New(lis)
	// Match connections in order: first gRPC, then HTTP.
	grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpL := m.Match(cmux.HTTP1Fast())

	// grpc and http server
	srv := grpc.NewServer()
	reflection.Register(srv)
	httpS := &http.Server{}

	// Use the servers in goroutines.
	go srv.Serve(grpcL)
	go httpS.Serve(httpL)

	// grpc handler
	grpch := grpchandler.New(ctrl)
	gen.RegisterFollowServiceServer(srv, grpch)
	// http handler
	httph := httphandler.New(ctrl)
	http.Handle("/follow", http.HandlerFunc(httph.Follow))
	http.Handle("/unfollow", http.HandlerFunc(httph.Unfollow))
	http.Handle("/user_followers", http.HandlerFunc(httph.GetUserFollowers))
	http.Handle("/following_user", http.HandlerFunc(httph.GetFollowingUser))
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	// Start serving!
	log.Fatal(m.Serve())
}
