package main

import (
	"log"
	"net"
	"restaurants/src/database"
	"restaurants/src/proto"
	"restaurants/src/servers"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Connect to the mongodb database
	db, err := database.ConnectDB()

	if err != nil {
		println("ERROR")
		log.Fatalf("%f", err)
	}

	// Define the mongo collections
	db_name := "restaurant"
	collection := db.Database(db_name).Collection("restaurant")

	// Create the grpc server
	server := grpc.NewServer()
	restaurant := servers.NewRestaurantServer(collection)
	proto.RegisterRestaurantServiceServer(server, restaurant)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	println("gRPC order server running")
}
