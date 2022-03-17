package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"v0/deletepb"
	"v0/putpb"
	"v0/readpb"
	"v0/updatepb"

	"github.com/gorilla/mux"

	"google.golang.org/grpc"
)

//this package act like simple api gateway
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

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/users/ReadAll", func(w http.ResponseWriter, r *http.Request) {
		//Read
		//grpc call
		fmt.Println("user get all client")
		cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("could not connect: %v", err)
		}
		defer cc.Close()

		c := readpb.NewCrudServiceClient(cc)
		//sending the request
		Read_All(c)

	}).Methods(http.MethodGet)
	router.HandleFunc("/users/ReadById", func(w http.ResponseWriter, r *http.Request) {
		//Read
		id_str := r.FormValue("id")
		id_int, _ := strconv.ParseInt(id_str, 10, 64)
		//grpc call
		fmt.Println("user getby id")
		cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("could not connect: %v", err)
		}
		defer cc.Close()
		c := readpb.NewCrudServiceClient(cc)

		//sending request
		Read_ById(c, int32(id_int))

	}).Methods(http.MethodGet)

	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		//put
		var user User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			_, _ = w.Write([]byte("decoding failed"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte(fmt.Sprintf("created %+v", user)))
		//grpc call
		fmt.Println("Create Create-User Client")
		cc, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("could not connect: %v", err)
		}
		defer cc.Close()

		c := putpb.NewCrudServiceClient(cc)

		//sending the request
		Put_user(c, user)

	}).Methods(http.MethodPost)

	router.HandleFunc("/cards", func(w http.ResponseWriter, r *http.Request) {
		//put
		var card Card

		if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
			_, _ = w.Write([]byte("decoding failed"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//grpc call
		fmt.Println(" Client")
		cc, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("could not connect: %v", err)
		}
		defer cc.Close()

		c := putpb.NewCrudServiceClient(cc)
		//sending the request
		Put_cards(c, card, card.Id)

	}).Methods(http.MethodPost)

	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		//update
		var user User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			_, _ = w.Write([]byte("decoding failed"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//grpc call
		fmt.Println("Update User Client")
		cc, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("could not connect: %v", err)
		}
		defer cc.Close()
		c := updatepb.NewCrudServiceClient(cc)

		//sending the request
		Update_user(c, user)

	}).Methods(http.MethodPut)
	router.HandleFunc("/cards", func(w http.ResponseWriter, r *http.Request) {
		//update
		var card Card

		if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
			_, _ = w.Write([]byte("decoding failed"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//grpc call
		fmt.Println("Update User Client")
		cc, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("could not connect: %v", err)
		}
		defer cc.Close()

		c := updatepb.NewCrudServiceClient(cc)

		// sending the request
		Update_card(c, card)
	}).Methods(http.MethodPut)

	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		//delete
		fmt.Println("Delete User Client")
		id_str := r.FormValue("id")
		id_int, _ := strconv.ParseInt(id_str, 10, 64)

		// sending the request
		Delete_User(int32(id_int))
	}).Methods(http.MethodDelete)

	router.HandleFunc("/cards", func(w http.ResponseWriter, r *http.Request) {
		//delete
		fmt.Println("Delete User Client")
		id_str := r.FormValue("id")
		id_int, _ := strconv.ParseInt(id_str, 10, 64)

		// sending the request
		Delete_Card(int32(id_int))

	}).Methods(http.MethodDelete)

	s := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Fatal(s.ListenAndServe())
}
func Read_All(c readpb.CrudServiceClient) readpb.GetUserResponse {
	fmt.Println("Starting to Reading All the users...")
	req := &readpb.GetUserRequest{
		Id:  0,
		All: true,
	}
	stream, err := c.GetUser(context.Background(), req)
	if err != nil {
		log.Fatalf("error while Reading users: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		return *res
	}
	return readpb.GetUserResponse{}
}
func Read_ById(c readpb.CrudServiceClient, id int32) {
	fmt.Println("Starting to Reading User by Id...")
	req := &readpb.GetUserRequest{
		Id:  id,
		All: false,
	}
	res, err := c.GetUser(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}
	log.Printf("Response from ReadAll: %v", res)
}
func Put_user(c putpb.CrudServiceClient, user User) {
	fmt.Println("Starting to Putting New Card to DataBase...")
	req := &putpb.CreateUserRequest{
		Id:       user.Id,
		Name:     user.Name,
		Birthday: user.Birthday,
		Gender:   user.Gender,
		Avatar:   user.Avatar,
	}
	res, err := c.CreateUser(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Create-User: %v", err)
	}
	log.Printf("Response from Put-User: %v", res)
}
func Put_cards(c putpb.CrudServiceClient, card Card, user_id int32) {
	fmt.Println("Starting to Putting New Card to DataBase...")
	req := &putpb.CreateCardRequest{
		UserId:   user_id,
		Bankname: card.BandName,
		Serial:   card.serial,
	}
	res, err := c.CreateCard(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling put: %v", err)
	}
	log.Printf("Response from Put_Card: %v", res)
}
func Update_card(c updatepb.CrudServiceClient, card Card) {
	fmt.Println("Starting to Update card...")
	req := &updatepb.UpdateCardRequest{
		Id:       card.Id,
		BankName: card.BandName,
		Serial:   card.serial,
	}
	res, err := c.UpdateCard(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Update-User: %v", err)
	}
	log.Printf("Response from Update-User: %v", res)
}
func Update_user(c updatepb.CrudServiceClient, user User) {
	fmt.Println("Starting to do a Update user")
	req := &updatepb.UpdateUserRequest{
		Id:       user.Id,
		Name:     user.Name,
		Birthday: user.Birthday,
		Gender:   user.Gender,
		Avatar:   user.Avatar,
	}
	res, err := c.UpdateUser(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Update-User: %v", err)
	}
	log.Printf("Response from Update-User: %v", res)
}
func Delete_Card(id int32) {
	fmt.Println("Starting to Delete user...")

	cc, err := grpc.Dial("localhost:50054", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := deletepb.NewCrudServiceClient(cc)

	req := &deletepb.DeleteCardRequest{
		Id: id,
	}
	res, err := c.DeleteCard(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Create-User: %v", err)
	}
	log.Printf("Response from Create-User: %v", res)
}
func Delete_User(id int32) {
	fmt.Println("Starting to Delete user...")

	cc, err := grpc.Dial("localhost:50054", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := deletepb.NewCrudServiceClient(cc)

	req := &deletepb.DeleteUserRequest{
		Id: id,
	}
	res, err := c.DeleteUser(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Create-User: %v", err)
	}
	log.Printf("Response from Create-User: %v", res)
}
