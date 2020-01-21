package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "12345"
	dbname   = "userDB"
)

var db *sql.DB
var tpl_login = template.Must(template.ParseFiles("login.html"))
var tpl_signup = template.Must(template.ParseFiles("signup.html"))

func main() {
	routingHandler()
}

func initDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
}

func routingHandler() {
	initDB()
	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/login", loginHandler)
	router.HandleFunc("/signup", signUpHandler)
	log.Fatalln(http.ListenAndServe("localhost:8000", router))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>User List :</h1>" + "<br>"))
	rows, err := db.Query("SELECT username FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var username string
		err = rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		w.Write([]byte("<h1>"+username+"</h1>" + "<br>"))

	}
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	tpl_signup.Execute(w, nil)
	r.ParseForm()
	username := r.FormValue("username")
	password1 := r.FormValue("password1")
	password2 := r.FormValue("password2")
	log.Println(password1, password2)
	if string(password1) != string(password2) {
		w.Write([]byte("<h1>Sorry two password doesn't match"))
	} else {
		if _, err := db.Query("insert into users values ($1, $2)", username, password1); err != nil {
			w.Write([]byte("<h1>Sorry</h1>"))
		} else {
			w.Write([]byte("<script>alert('Success! Please login')</script>"))
		}
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	tpl_login.Execute(w, nil)
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	result := db.QueryRow("select password from users where username=$1", username)
	var obtained_password string
	err := result.Scan(&obtained_password)
	if err != nil {
		if err == sql.ErrNoRows{
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("<h1>No user exist </h1>"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if obtained_password != password {
		w.Write([]byte("<h1>Login failed </h1>"))
		w.WriteHeader(http.StatusUnauthorized)
	}else {
		w.Write([]byte("<h1>Success </h1>"))
	}

}
