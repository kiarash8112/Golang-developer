package update

import (
	"context"
	"fmt"
	"log"
	"net"
	"v0/domain"
	"v0/updatepb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/mgo.v2/bson"
)

type server struct{}

type User struct {
	Id       int32
	Name     string
	Gender   string
	Birthday string
	Avatar   string
}
type Card struct {
	Id       int32
	BandName string
	serial   string
}

func (*server) UpdateUser(ctx context.Context, req *updatepb.UpdateUserRequest) (*updatepb.UpdateResponse, error) {
	var response *updatepb.UpdateResponse
	client, ctx := domain.Connect_to_db()
	domain.Client = client
	domain.Ctx = ctx
	if domain.Err != nil {
		fmt.Println("cant connect to database", domain.Err)
	}

	db := domain.MongoDbQueryBuilder{}
	filter := bson.D{{"Id", req.Id}}
	item := bson.D{{"Id", req.Id}, {"Name", req.Name}, {"Gender", req.Gender}, {"Birthday", req.Birthday}, {"Avatar", req.Avatar}}
	db.Db("user").Cols("user").Update(filter, item)
	if domain.Err != nil {
		response = &updatepb.UpdateResponse{
			Result: "update faild",
		}
		return response, domain.Err
	} else {
		response = &updatepb.UpdateResponse{
			Result: "update successful",
		}

	}
	return response, nil
}
func (*server) UpdateCard(ctx context.Context, req *updatepb.UpdateCardRequest) (*updatepb.UpdateResponse, error) {
	var response *updatepb.UpdateResponse
	client, ctx := domain.Connect_to_db()
	domain.Client = client
	domain.Ctx = ctx
	if domain.Err != nil {
		fmt.Println("cant connect to database", domain.Err)
	}

	db := domain.MongoDbQueryBuilder{}
	filter := bson.D{{"Id", req.Id}}
	item := bson.D{{"Id", req.Id}, {"BankName", req.BankName}, {"Serial", req.Serial}}
	db.Db("user").Cols("user").Update(filter, item)
	if domain.Err != nil {
		response = &updatepb.UpdateResponse{
			Result: "transaction faild",
		}
		return response, domain.Err
	} else {
		response = &updatepb.UpdateResponse{
			Result: "transaction successful",
		}

	}
	return response, nil
}
func main() {
	fmt.Println("update Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	updatepb.RegisterCrudServiceServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
