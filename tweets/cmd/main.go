package main

import (
	"flag"
	"fmt"
	gen "github.com/alexvishnevskiy/twitter-clone/gen/api/tweets"
	localcache "github.com/alexvishnevskiy/twitter-clone/internal/cache/local"
	"github.com/alexvishnevskiy/twitter-clone/internal/storage/local"
	_ "github.com/alexvishnevskiy/twitter-clone/tweets/docs"
	"github.com/alexvishnevskiy/twitter-clone/tweets/internal/controller"
	grpchandler "github.com/alexvishnevskiy/twitter-clone/tweets/internal/handler/grpc"
	httphandler "github.com/alexvishnevskiy/twitter-clone/tweets/internal/handler/http"
	"github.com/alexvishnevskiy/twitter-clone/tweets/internal/repository/mysql"
	"github.com/soheilhy/cmux"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
)

// @title			Tweets API documentation
// @version		1.0.0
// @host			localhost:8080
// @description	This is API for tweets service
func main() {
	var (
		port        int
		capacity    int
		storagePath string
	)
	flag.IntVar(&port, "port", 8080, "API handler port")
	flag.IntVar(&capacity, "capacity", 5000, "Capacity of cache")
	flag.StringVar(&storagePath, "storage_path", getStoragePath(), "storage path")
	flag.Parse()
	log.Printf("Starting the tweets service on port %d", port)

	// probably use configs for driverName and dataSourceName
	repository, err := mysql.New("mysql", "root:root@tcp(localhost:3306)/twitter")
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	storage := local.New(storagePath)
	cache := localcache.New(capacity)
	ctrl := controller.New(repository, storage, cache)

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
	// TODO: make multiple ports for grpc and http
	go srv.Serve(grpcL)
	go httpS.Serve(httpL)

	// grpc handler
	grpch := grpchandler.New(ctrl)
	gen.RegisterTweetsServiceServer(srv, grpch)
	// http handler
	httph := httphandler.New(ctrl)
	http.Handle("/post_tweet", http.HandlerFunc(httph.Post))
	http.Handle("/retrieve_tweet", http.HandlerFunc(httph.Retrieve))
	http.Handle("/delete_tweet", http.HandlerFunc(httph.Delete))
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	// Start serving!
	log.Fatal(m.Serve())
}
