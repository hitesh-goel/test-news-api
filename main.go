package main

import (
	"log"
	"net/http"

	"github.com/hitesh-goel/test-news-api/handlers"
	"github.com/hitesh-goel/test-news-api/rediscon"
	"github.com/julienschmidt/httprouter"
)

//GetHandler to handle get Requests
func GetHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	handlers.GetNews(w, r)
}

func main() {
	server := httprouter.New()
	server.GET("/v1/news-api", GetHandler)

	// Initiliase redis client for caching purpose
	rediscon.ConnectToRedis()
	// Start http server Listen to port 8080
	log.Fatal(http.ListenAndServe(":8080", server))
}
