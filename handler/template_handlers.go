package handler

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)


var tpl_login = template.Must(template.ParseFiles("./templates/login.html"))
var tpl_signup = template.Must(template.ParseFiles("./templates/signup.html"))
var tpl_index = template.Must(template.ParseFiles("./templates/index.html"))

type User struct {
	Username string
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
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

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
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

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT username FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var Users []User
	for rows.Next() {
		var username string
		var user User
		err = rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		user.Username = username
		Users = append(Users, user)

	}
	log.Println(len(Users))
	tpl_index.Execute(w, Users)
}