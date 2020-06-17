package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// "sync"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "microservice-go/gRPC/article"
	pbU "microservice-go/gRPC/user"
)

func ConnectArticleService(artCh chan []ArticleData) {
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewArticleServiceClient(conn)

	request := &pb.ListArticleRequest{
		Keyword: "test",
	}

	resp, err := client.GetArticle(context.Background(), request)

	var articleData []ArticleData

	if resp.Status == "success" {
		fmt.Println("Total Article: ", resp.Total)
		fmt.Println("Articles: ", resp.Articles)
		for i, v := range resp.Articles {
			articleData = append(articleData, ArticleData{
				ID:          v.Id,
				Title:       v.Title,
				Slug:        v.Slug,
				Description: v.Description,
				AuthorID:    v.AuthorId,
			})
			fmt.Println("Article ", i, ": ", v.Id, ", ", v.Title)
		}
	}
	artCh <- articleData
	// wg.Done()
}

func ConnectUserService(authorCh chan []AuthorData) {
	conn, err := grpc.Dial("127.0.0.1:8081", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	defer conn.Close()

	client := pbU.NewUserServiceClient(conn)

	request := &pbU.ListAuthorRequest{
		Keyword: "test",
	}

	resp, err := client.GetAuthors(context.Background(), request)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var authorData []AuthorData

	if resp.Status == "success" {
		fmt.Println("Total Authors: ", resp.Total)
		fmt.Println("Authors: ", resp.Authors)
		for i, v := range resp.Authors {
			authorData = append(authorData, AuthorData{
				ID:       int64(v.Id),
				Username: v.Usename,
				Email:    v.Email,
				Bio:      v.Bio,
				Image:    v.Image,
			})
			fmt.Println("Author ", i, ": ", v.Id, ", ", v.Usename)
		}
	}
	authorCh <- authorData
	// wg.Done()
}

type AuthorData struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}

type ArticleData struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	AuthorID    int64  `json:"author_id"`
}

type ArticleReplyData struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	Author      AuthorData `json:"author"`
}

func GetAllArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	articleCh := make(chan []ArticleData)
	authorCh := make(chan []AuthorData)

	// var wg sync.WaitGroup
	// wg.Add(2)
	go ConnectArticleService(articleCh)
	go ConnectUserService(authorCh)
	// wg.Wait()
	fmt.Println("Finished call services!")

	articles := <-articleCh
	authors := <-authorCh

	authorMap := make(map[int64]AuthorData)
	for _, v := range authors {
		authorMap[v.ID] = v
	}

	var articlesRep []ArticleReplyData

	for _, v := range articles {
		articlesRep = append(articlesRep, ArticleReplyData{
			ID:          v.ID,
			Title:       v.Title,
			Slug:        v.Slug,
			Description: v.Description,
			Author:      authorMap[v.AuthorID],
		})
	}

	jsonResponse, err := json.Marshal(articlesRep)
	log.Println(err)
	w.Write(jsonResponse)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/articles", GetAllArticles).Methods("GET")

	log.Fatal(http.ListenAndServe(":5000", router))
}
