package main

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"microservice-go/article-service/models"
	pb "microservice-go/gRPC/article"
)

type ArticleServer struct {
	request []*pb.ListArticleRequest
}

type ArticleDB struct {
	DB *models.DB
}

var artDB ArticleDB

func (s *ArticleServer) GetArticle(ctx context.Context, req *pb.ListArticleRequest) (*pb.ListArticleResponse, error) {
	articles, err := artDB.DB.GetListArticle()
	var status string
	if err != nil {
		log.Println("ArticleServer Error: ", err)
		status = "error"
	} else {
		status = "success"
	}

	var articlesPb []*pb.Article

	for _, v := range articles {
		articlesPb = append(articlesPb, &pb.Article{
			Id:          int64(v.ID),
			Title:       v.Title,
			Slug:        v.Slug,
			Description: v.Description,
			Body:        v.Body,
			AuthorId:    int64(v.AuthorID),
			CreatedAt:   v.CreatedAt.String(),
			UpdatedAt:   v.UpdatedAt.String(),
		})
	}

	articlesRes := &pb.ListArticleResponse{
		Status:   status,
		Total:    int32(len(articles)),
		Articles: articlesPb,
	}

	return articlesRes, nil
}

func main() {
	// Connect DB
	db, errDB := models.ConnectDB()
	if errDB != nil {
		log.Println("Error: ", errDB)
	}

	artDB = ArticleDB{
		DB: db,
	}

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Failed to listen: %v", err)
	}

	srv := grpc.NewServer()

	pb.RegisterArticleServiceServer(srv, &ArticleServer{})
	srv.Serve(lis)
}
