package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"v0/deletepb"
	"v0/domain"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/mgo.v2/bson"
)

type server struct{}

func (*server) DeleteUser(ctx context.Context, req *deletepb.DeleteUserRequest) (*deletepb.DeleteResponse, error) {
	var response *deletepb.DeleteResponse

	client, ctx := domain.Connect_to_db()
	domain.Client = client
	domain.Ctx = ctx
	if domain.Err != nil {
		fmt.Println("cant connect to database", domain.Err)
	}
	db := domain.MongoDbQueryBuilder{}

	filter := bson.D{{"Id", req.Id}}
	db.Db("user").Cols("user").Delete(filter)
	if domain.Err != nil {
		response = &deletepb.DeleteResponse{
			Result: "delete faild",
		}
		return response, domain.Err
	} else {
		response = &deletepb.DeleteResponse{
			Result: "delete successful",
		}

	}
	return response, nil
}

func (*server) DeleteCard(ctx context.Context, req *deletepb.DeleteCardRequest) (*deletepb.DeleteResponse, error) {
	var response *deletepb.DeleteResponse

	client, ctx := domain.Connect_to_db()
	domain.Client = client
	domain.Ctx = ctx
	if domain.Err != nil {
		fmt.Println("cant connect to database", domain.Err)
	}

	db := domain.MongoDbQueryBuilder{}

	filter := bson.D{{"Id", req.Id}}
	db.Db("user").Cols("user").Delete(filter)
	if domain.Err != nil {
		response = &deletepb.DeleteResponse{
			Result: "delete faild",
		}
		return response, domain.Err
	} else {
		response = &deletepb.DeleteResponse{
			Result: "delete successful",
		}

	}
	return response, nil
}
func main() {
	fmt.Println("delete user Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50054")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	deletepb.RegisterCrudServiceServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
