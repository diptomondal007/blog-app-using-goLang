package main

import (
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	handlers "loginregistration/handler"
	"net/http"
)


func main() {
	handlers.InitDB()
	routingHandler()
}

func routingHandler() {
	router := mux.NewRouter()
	router.HandleFunc("/", handlers.IndexHandler)
	router.HandleFunc("/login", handlers.LoginHandler)
	router.HandleFunc("/signup", handlers.SignUpHandler)
	log.Fatalln(http.ListenAndServe("localhost:8000", router))
}





