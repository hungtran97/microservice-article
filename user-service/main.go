package main

import (
	"fmt"
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "microservice-go/gRPC/user"
	"microservice-go/user-service/models"
)

type UserServer struct {
	request []*pb.ListAuthorRequest
}

type UserDB struct {
	DB *models.DB
}

var userDB UserDB

func (s *UserServer) GetAuthors(ctx context.Context, req *pb.ListAuthorRequest) (*pb.ListAuthorReply, error) {
	users, err := userDB.DB.ListUser()
	var status string
	if err != nil {
		log.Println("ArticleServer Error: ", err)
		status = "error"
	} else {
		status = "success"
	}

	var authorPb []*pb.Author

	for _, v := range users {
		authorPb = append(authorPb, &pb.Author{
			Id:      int32(v.ID),
			Usename: v.Username,
			Email:   v.Email,
			Bio:     v.Bio,
			Image:   *v.Image,
		})
	}

	authorRes := &pb.ListAuthorReply{
		Status:  status,
		Total:   int32(len(users)),
		Authors: authorPb,
	}

	return authorRes, nil
}

func main() {
	fmt.Println("vim-go")
	// Connect DB
	db, errDB := models.ConnectDB()
	if errDB != nil {
		log.Println("Error: ", errDB)
	}

	userDB = UserDB{
		DB: db,
	}

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal("Failed to listen: %v", err)
	}

	srv := grpc.NewServer()

	pb.RegisterUserServiceServer(srv, &UserServer{})
	srv.Serve(lis)
}
