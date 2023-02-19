package main

import (
	"context"
	"fmt"
	"log"
	"main/uc/pb"
	"main/ucGrpcServer"
	"net"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Server is Running ...")
	//grpclog.Println("Server is Running ...")
	port := ":8080"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln(" Failed to listen ...")
	}
	//grpclog.Println("listening on 127.0.0.1:8080")
	fmt.Println("listening on 127.0.0.1:8080")
	var opts []grpc.ServerOption
	server := grpc.NewServer(opts...)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		//log.Fatalln(err)
		fmt.Println(err.Error())
	}
	defer client.Disconnect(ctx)

	IUCGrpc := ucGrpcServer.NewUCGrpcServerStruct(client)
	if err != nil {
		//log.Fatalln(err)
		fmt.Println(err.Error())
	}

	pb.RegisterUCServiceServer(server, IUCGrpc)
	log.Println("register now...")
	err = server.Serve(listener)

	if err != nil {
		//log.Fatalln(err)
		fmt.Println(err.Error())
	}
}
