package handler

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	validate "loginregistration/validation"
	"net/http"
	"os"
	"strconv"
	"time"
)

var tplUpdate = template.Must(template.ParseFiles("./templates/update.html"))

//User struct to pass data to the html templates
type User struct {
	Username                string
	IsFollowedByCurrentUser bool
	IsOwnedThisAccount      bool
}

var loggedin_user string

//Post struct
type Post struct {
	Id                   int64
	Body                 string
	Username             string
	Likes                int64
	IsLikedByCurrentUser bool
	Editable             bool
	Deletable            bool
}

type IndexData struct {
	Posts          []Post
	LoggedUser     string
	IsLoggedInUser bool
}

type UserPageData struct {
	Users          []User
	LoggedUser     string
	IsLoggedInUser bool
}

type UpdatePageData struct {
	Post           Post
	LoggedUser     string
	IsLoggedInUser bool
}

type ProfilePageData struct {
	Posts          []Post
	FirstName      string
	LastName       string
	Email          string
	ProfilePic     string
	LoggedUser     string
	IsLoggedInUser bool
	Followers      int64
	Token          string
}

//LoginPageHandler for rendering Login page
func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	loggedin_user = ""
	tplLogin, err := template.ParseFiles("./templates/login.html", "./templates/base.html")
	if err != nil {
		log.Println(err)
	}
	tplLogin.Execute(w, nil)
}

//LoginHandler for handling post data from login page
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	result := db.QueryRow("select password from users where username=$1", username)
	var obtainedPassword string
	err := result.Scan(&obtainedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("<script>alert('No user exist!')</script>"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if obtainedPassword != password {
		w.Write([]byte("<script>alert('Login Failed!')</script>"))
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		loggedin_user = username
		http.Redirect(w, r, "/", 302)
	}

}

//SignUpPageHandler for rendering sign up page
func SignUpPageHandler(w http.ResponseWriter, r *http.Request) {
	tplRegister, err := template.ParseFiles("./templates/signup.html", "./templates/base.html")
	if err != nil {
		log.Println(err)
	}
	tplRegister.Execute(w, nil)
}

//SignUpHandler for getting post request and handle them
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password1 := r.FormValue("password1")
	password2 := r.FormValue("password2")
	_username, _password1, _password2 := false, false, false
	_username = !validate.IsEmpty(username)
	_password1 = !validate.IsEmpty(password1)
	_password2 = !validate.IsEmpty(password2)
	if _username && _password1 && _password2 {
		if string(password1) != string(password2) {
			http.Redirect(w, r, "/signup", 302)
		} else {
			if _, err := db.Query("insert into users values ($1, $2)", username, password1); err != nil {
				w.Write([]byte("<script>alert('Error occurred!')</script>"))
			} else {
				w.Write([]byte("<script>alert('Success! Please login')</script>"))
			}
		}
	} else {
		w.Write([]byte("<script>alert('Sorry! Fields can not be empty')</script>"))
	}

}

//DeleteHandler for handling request of deleting a user
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	_, err := db.Query("DELETE FROM users WHERE username=$1", username)
	if err != nil {
		panic(err.Error())
	}
	log.Println("DELETE")
	http.Redirect(w, r, "/", 301)
}

//UpdatePage to render user information update page
func UpdatePage(w http.ResponseWriter, r *http.Request) {
	tplUpdate.Execute(w, loggedin_user)
}

//UpdateHandler for handling submitted update data
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	userToBeUpdated := r.URL.Query().Get("username")
	username := r.FormValue("username")
	password1 := r.FormValue("password1")
	password2 := r.FormValue("password2")
	_username, _password1, _password2 := false, false, false
	_username = !validate.IsEmpty(username)
	_password1 = !validate.IsEmpty(password1)
	_password2 = !validate.IsEmpty(password2)
	if _username && _password1 && _password2 {
		if string(password1) != string(password2) {
			http.Redirect(w, r, "/signup", 302)
		} else {
			if _, err := db.Query("update users set username=$1,password=$2 where username =$3", username, password1, userToBeUpdated); err != nil {
				w.Write([]byte("<script>alert('Error occurred!')</script>"))
			} else {
				loggedin_user = username
				http.Redirect(w, r, "/", 302)
			}
		}
	} else {
		w.Write([]byte("<script>alert('Sorry! Fields can not be empty')</script>"))
	}

}

//IndexPageHandler for rendering and handling Index page
func IndexPageHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("select posts.id,posts.body,posts.username from posts inner join followers on followers.following_id=posts.username where followers.user_id =$1 order by posts.id desc", loggedin_user)
	var posts []Post
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var body string
		var username string
		var post Post
		err = rows.Scan(&id, &body, &username)
		if err != nil {
			log.Println(err)
		}
		post.Likes = 0
		likesQuery, err := db.Query("SELECT is_liked from likes where post_id=$1", id)
		if err != nil {
			log.Println("Likes query failed")
		}
		for likesQuery.Next() {
			var isLiked bool
			err = likesQuery.Scan(&isLiked)
			if isLiked {
				post.Likes += 1
			}

		}
		post.IsLikedByCurrentUser = false
		isLikedByCurrentUserQuery, err := db.Query("SELECT is_liked from likes where user_name=$1 and post_id=$2", loggedin_user, id)
		if err != nil {
			log.Println("Is Liked by current User query failed")
		}
		for isLikedByCurrentUserQuery.Next() {
			var isLiked bool
			isLikedByCurrentUserQuery.Scan(&isLiked)
			if isLiked {
				post.IsLikedByCurrentUser = true
			}
		}
		post.Id = id
		post.Body = body
		post.Username = username
		if loggedin_user == username {
			post.Editable = true
			post.Deletable = true
		}
		log.Println(post.IsLikedByCurrentUser)
		posts = append(posts, post)

	}
	var IndexPageData IndexData
	if len(loggedin_user) < 1 {
		IndexPageData.IsLoggedInUser = false
	} else {
		IndexPageData.IsLoggedInUser = true
	}
	IndexPageData.Posts = posts
	IndexPageData.LoggedUser = loggedin_user
	tm := template.Must(template.ParseFiles("./templates/index.html", "./templates/base.html"))
	errortm := tm.Execute(w, IndexPageData)
	log.Println(errortm)
}

//IndexHandler for post request of post data
func Indexhandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	body := r.FormValue("body")
	log.Println(body)
	user := loggedin_user
	if _, err := db.Query("insert into Posts(body,username) values ($1, $2)", body, user); err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", 302)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

func UpdatePostPage(w http.ResponseWriter, r *http.Request) {
	tplPostUpdate, err := template.ParseFiles("./templates/updatePost.html", "./templates/base.html")
	if err != nil {
		log.Println(err)
	}
	postToBeUpdated := r.URL.Query().Get("id")
	rows, err := db.Query("SELECT * FROM Posts where id=$1", postToBeUpdated)
	if err != nil {
		panic(err)
	}
	var id int64
	var username string
	var body string
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id, &body, &username)
		if err != nil {
			panic(err)
		}
		log.Println(body)
	}
	var post Post
	post.Id = id
	post.Body = body
	post.Username = username

	var data UpdatePageData
	data.Post = post
	data.LoggedUser = loggedin_user
	if len(loggedin_user) < 1 {
		data.IsLoggedInUser = false
	} else {
		data.IsLoggedInUser = true
	}
	tplPostUpdate.Execute(w, data)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	postToBeUpdated := r.URL.Query().Get("id")
	body := r.FormValue("body")
	if _, err := db.Query("update Posts set body =$1 where id=$2", body, postToBeUpdated); err != nil {
		w.Write([]byte("<script>alert('Sorry!')</script>"))
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

//DeleteHandler for handling request of deleting a user
func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := db.Query("DELETE FROM Posts WHERE id=$1", id)
	if err != nil {
		panic(err.Error())
	}
	log.Println("POST DELETED")
	http.Redirect(w, r, "/", 301)
}

func UsersPageHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT username FROM users")
	var users []User
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var username string
		var user User
		err = rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		if username == loggedin_user {
			user.IsOwnedThisAccount = true
		}
		followRow, err := db.Query("SELECT following_id from followers where (user_id = $1 AND following_id= $2)", loggedin_user, username)
		if err != nil {
			log.Println(err)
		}
		user.Username = username
		for followRow.Next() {

			user.IsFollowedByCurrentUser = true
		}
		users = append(users, user)

	}
	var UserPageData UserPageData
	if len(loggedin_user) < 1 {
		UserPageData.IsLoggedInUser = false
	} else {
		UserPageData.IsLoggedInUser = true
	}
	UserPageData.Users = users
	UserPageData.LoggedUser = loggedin_user
	log.Println(users)
	tm := template.Must(template.ParseFiles("./templates/users.html", "./templates/base.html"))
	errortm := tm.Execute(w, UserPageData)
	log.Println(errortm)
}

func FollowUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	_, err := db.Query("insert into followers values ($1, $2)", loggedin_user, username)
	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/users", 301)

}

func UnFollowUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	_, err := db.Query("DELETE  from followers WHERE user_id = $1 AND following_id = $2", loggedin_user, username)
	if err != nil {
		panic(err)
	}
	log.Println("UnFollowed")
	http.Redirect(w, r, "/users", 301)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	loggedin_user = ""
	time.Sleep(100 * time.Millisecond)
	http.Redirect(w, r, "/login", 301)

}

func ProfilePage(w http.ResponseWriter, r *http.Request) {
	if len(loggedin_user) < 1 {
		http.Redirect(w, r, "/login", 302)
	} else {
		tplProfilePage, err := template.ParseFiles("./templates/profile_page.html", "./templates/base.html")
		if err != nil {
			log.Println(err)
		}
		var ProfilePageData ProfilePageData
		if len(loggedin_user) < 1 {
			ProfilePageData.IsLoggedInUser = false
		} else {
			ProfilePageData.IsLoggedInUser = true
		}
		userRows, errQuery := db.Query("SELECT firstname , lastname , email, profile_pic from users where username=$1", loggedin_user)
		if errQuery != nil {
			log.Println("Query Failed")
		} else {
			defer userRows.Close()
			for userRows.Next() {
				var firstName, lastName, email, profile_pic string
				errQuery := userRows.Scan(&firstName, &lastName, &email, &profile_pic)
				if errQuery != nil {
					log.Println("Sorry")
				}
				if len(firstName) < 1 {
					firstName = " "
				}
				if len(lastName) < 1 {
					lastName = " "
				}
				if len(email) < 1 {
					email = " "
				}
				log.Println(firstName, lastName, email)
				ProfilePageData.FirstName, ProfilePageData.LastName, ProfilePageData.Email, ProfilePageData.ProfilePic = firstName, lastName, email, profile_pic
			}
		}
		log.Println(ProfilePageData.FirstName)
		rows, err := db.Query("select id, body,username from posts where username =$1 order by posts.id desc", loggedin_user)
		var posts []Post
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		for rows.Next() {
			var id int64
			var body string
			var username string
			var post Post
			err = rows.Scan(&id, &body, &username)
			if err != nil {
				log.Println(err)
			}
			post.Id = id
			post.Body = body
			post.Username = username
			post.Editable = true
			post.Deletable = true
			posts = append(posts, post)

		}
		followersRow, err := db.Query("SELECT following_id from followers where user_id=$1", loggedin_user)
		if err != nil {
			log.Println("Query failed")
		}
		defer followersRow.Close()
		for followersRow.Next() {
			ProfilePageData.Followers += 1
		}
		ProfilePageData.Posts = posts
		ProfilePageData.LoggedUser = loggedin_user
		log.Println(posts)

		//Token generate
		timeNow := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(timeNow, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		ProfilePageData.Token = token
		tplProfilePage.Execute(w, ProfilePageData)
	}

}

func ProfilePageInputHandler(w http.ResponseWriter, r *http.Request) {

	file, handler, err := r.FormFile("uploadfile")
	if handler == nil {
		r.ParseForm()
		firstname := r.FormValue("firstName")
		lastname := r.FormValue("lastName")
		email := r.FormValue("email")
		log.Println(firstname, lastname, email)
		_, errNew := db.Query("update users set firstname=$1, lastname = $2, email=$3 where username =$4", firstname, lastname, email, loggedin_user)
		if errNew != nil {
			log.Println(errNew)
		} else {
			http.Redirect(w, r, "/profile", 302)
		}
	} else {
		r.ParseMultipartForm(32 << 20)

		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		f, err := os.OpenFile("./images/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		_, errUpdateProfilePic := db.Query("update users set profile_pic =$1 where username =$2", handler.Filename, loggedin_user)
		if errUpdateProfilePic != nil {
			log.Println("Sorry")
		} else {
			http.Redirect(w, r, "/profile", 302)
		}
	}
}

func PostLikeHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	log.Println(id)
	_, err := db.Query("INSERT INTO likes(post_id,user_name,is_liked) values($1,$2,'t')", id, loggedin_user)
	if err != nil {
		log.Println("Inserting Like Failed")
	} else {
		http.Redirect(w, r, "/", 302)
	}

}

func PostUnlikeHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	postLikeByUserIdQuery, err := db.Query("SELECT id from likes where post_id=$1 and user_name =$2", id, loggedin_user)
	if err != nil {
		log.Println("Failed")
	} else {
		for postLikeByUserIdQuery.Next() {
			var likeId int64
			err = postLikeByUserIdQuery.Scan(&likeId)
			log.Println(likeId)
			_, err = db.Query("update likes set post_id = $1 , user_name = $2,is_liked ='f' where id =$3", id, loggedin_user, likeId)
			if err != nil {
				log.Println("Post Like Query Failed")
			} else {
				http.Redirect(w, r, "/", 302)
			}
		}
	}

}
