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
	router.HandleFunc("/", handlers.IndexPageHandler).Methods("GET")
	router.HandleFunc("/", handlers.Indexhandler).Methods("POST")
	router.HandleFunc("/login", handlers.LoginPageHandler).Methods("GET")
	router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	router.HandleFunc("/signup", handlers.SignUpPageHandler).Methods("GET")
	router.HandleFunc("/signup", handlers.SignUpHandler).Methods("POST")
	router.HandleFunc("/delete",handlers.DeleteHandler)
	router.HandleFunc("/update",handlers.UpdateHandler).Methods("POST")
	router.HandleFunc("/update",handlers.UpdatePage).Methods("GET")
	router.HandleFunc("/updatePost",handlers.UpdatePost).Methods("POST")
	router.HandleFunc("/updatePost",handlers.UpdatePostPage).Methods("GET")
	router.HandleFunc("/deletePost",handlers.DeletePostHandler)
	router.HandleFunc("/users",handlers.UsersPageHandler)
	router.HandleFunc("/logout",handlers.Logout)
	router.HandleFunc("/follow",handlers.FollowUser)
	router.HandleFunc("/unfollow",handlers.UnFollowUser)
	log.Fatalln(http.ListenAndServe("localhost:8000", router))
}





