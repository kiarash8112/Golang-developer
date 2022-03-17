package read

import (
	"context"
	"fmt"
	"log"
	"net"
	"v0/domain"
	"v0/readpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/mgo.v2/bson"
)

type server struct {
}
type User struct {
	Id       int32 // for finding cards
	Name     string
	Gender   string
	Birthday string
	Avatar   string
}

func (server *server) GetUser(req *readpb.GetUserRequest, stream readpb.CrudService_GetUserServer) error {
	var user_result []bson.D
	var card_result []bson.D
	var user_slice []User
	var response_card []*readpb.Card
	var user User
	var card readpb.Card

	client, ctx := domain.Connect_to_db()
	domain.Client = client
	domain.Ctx = ctx
	if domain.Err != nil {
		return domain.Err
	}

	db_user := domain.MongoDbQueryBuilder{}
	db_card := domain.MongoDbQueryBuilder{}

	if req.All == false {
		db_user.Db("user").Cols("user").Find(req.Id)
		db_card.Db("user").Cols("card").Find(req.Id)
		if domain.Err != nil {
			return domain.Err
		}
		err := db_user.Result.Decode(&user_result[0])
		if err != nil {
			fmt.Println("cant decode db_user cursor", err)
		}
		err = db_card.Result.All(context.Background(), &card_result)
		if err != nil {
			fmt.Println("cant decode db_card cursor", err)
		}

		byte_user, err_user := bson.Marshal(user_result[0])
		if err_user != nil {
			fmt.Println("cant marshal user_result", err_user)
		}

		err = bson.Unmarshal(byte_user, user)
		if err != nil {
			fmt.Println("cant Unmarshal byte_user", err)
		}

		for counter, result := range card_result {
			byte_card, err := bson.Marshal(result)
			if err != nil {
				fmt.Println("cant marshal result", err)
			}
			err = bson.Unmarshal(byte_card, card)
			if err != nil {
				fmt.Println("cant Unmarshal card_user", err)
			}
			response_card[counter] = &card

		}

		stream.Send(&readpb.GetUserResponse{
			Name:     user_slice[0].Name,
			Gender:   user_slice[0].Gender,
			Birthday: user_slice[0].Birthday,
			Avatar:   user_slice[0].Avatar,
			Card:     response_card,
		})
	} else {

		db_user.Db("user").Cols("user").Find(0)
		if domain.Err != nil {
			fmt.Println("cant find user", domain.Err)
		}
		if err := db_user.Result.All(ctx, &user_result); err != nil {
			log.Fatal(err)
		}

		for _, result := range user_result {
			byte_user, err := bson.Marshal(result)
			if err != nil {
				fmt.Println("cant marshal result", err)
				return err
			}
			err = bson.Unmarshal(byte_user, user)
			if err != nil {
				fmt.Println("cant unmarshal byte_user", err)
			}
			db_card.Db("user").Cols("card").Find(user.Id)
			if domain.Err != nil {
				fmt.Println("cant find user", domain.Err)
			}
			if err := db_card.Result.All(ctx, &card_result); err != nil {
				log.Fatal(err)
				return err
			}

			for counter, result := range card_result {

				byte_card, err := bson.Marshal(result)
				if err != nil {
					fmt.Println("cant marshal result", err)
				}
				err = bson.Unmarshal(byte_card, card)
				if err != nil {
					fmt.Println("cant unmarshal byte_card", err)
				}
				response_card[counter] = &card

			}
			stream.Send(&readpb.GetUserResponse{
				Name:     user.Name,
				Gender:   user.Gender,
				Birthday: user.Birthday,
				Avatar:   user.Avatar,
				Card:     response_card,
			})
		}

	}

	return nil
}
func (*server) GetCard(ctx context.Context, req *readpb.GetCardRequest) (*readpb.GetCardResponse, error) {

	var card_result []bson.D
	var response_card []*readpb.Card
	var card readpb.Card

	client, ctx := domain.Connect_to_db()
	domain.Client = client
	domain.Ctx = ctx
	if domain.Err != nil {
		fmt.Println("cant connect to database", domain.Err)
	}
	db_card := domain.MongoDbQueryBuilder{}
	db_card.Db("user").Cols("card").Find(req.Id)
	if domain.Err != nil {
		fmt.Println("cant find card", domain.Err)
	}
	if err := db_card.Result.All(ctx, &card_result); err != nil {
		log.Fatal(err)
	}

	for counter, result := range card_result {

		byte_card, err := bson.Marshal(result)
		if err != nil {
			fmt.Println("cant marshal result", err)
		}
		err = bson.Unmarshal(byte_card, card)
		if err != nil {
			fmt.Println("cant unmarshal byte_card", err)
		}
		response_card[counter] = &card

	}

	res := &readpb.GetCardResponse{
		Card: response_card,
	}

	return res, nil
}
func main() {
	fmt.Println("Read User Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	readpb.RegisterCrudServiceServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
