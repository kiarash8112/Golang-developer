package put

import (
	"context"
	"fmt"
	"log"
	"net"
	"v0/domain"
	"v0/putpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/mgo.v2/bson"
)

type server struct {
}

func (*server) CreateUser(ctx context.Context, req *putpb.CreateUserRequest) (*putpb.CreateResponse, error) {

	var response *putpb.CreateResponse
	client, ctx := domain.Connect_to_db()
	domain.Client = client
	domain.Ctx = ctx
	if domain.Err != nil {
		fmt.Println("cant connect to database", domain.Err)
	}

	db := domain.MongoDbQueryBuilder{}

	item := bson.D{{"Id", req.Id}, {"Name", req.Name}, {"Gender", req.Gender},
		{"Birthdat", req.Birthday}, {"Avatar", req.Avatar}}

	db.Db("user").Cols("users").Create(item)
	if domain.Err != nil {
		response = &putpb.CreateResponse{
			Result: "transaction faild",
		}
		return response, domain.Err
	} else {
		response = &putpb.CreateResponse{
			Result: "transaction successful",
		}

	}

	return response, nil
}

func (*server) CreateCard(ctx context.Context, req *putpb.CreateCardRequest) (*putpb.CreateResponse, error) {
	var response *putpb.CreateResponse
	client, ctx := domain.Connect_to_db()
	domain.Client = client
	domain.Ctx = ctx
	if domain.Err != nil {
		fmt.Println("cant connect to database", domain.Err)
	}

	db := domain.MongoDbQueryBuilder{}
	item := bson.D{{"user_id", req.UserId}, {"bank_name", req.Bankname}, {"serial", req.Serial}}

	db.Db("user").Cols("card").Create(item)

	if domain.Err != nil {
		response = &putpb.CreateResponse{
			Result: "transaction faild",
		}
		return response, domain.Err
	} else {
		response = &putpb.CreateResponse{
			Result: "transaction successful",
		}

	}

	return response, nil
}

func main() {
	fmt.Println("create user Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	putpb.RegisterCrudServiceServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
