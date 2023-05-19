package wok

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// GET: returns the homepage
func Homepage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "../public/index.html")
	} else {
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
	}
}

// POST: create a user and insert them into the database
func CreatUserHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		user := new(User)
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(user)
		if err := CreateUser(user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	} else {
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
}

// POST: login user
func LoginHandle(w http.ResponseWriter, r *http.Request) {
	// check if method is POST, returns method not allowed if requests is not POST
	if r.Method == "POST" {
		// check if username is supplied
		if r.FormValue("username") == "" {
			w.Write([]byte("must include username\n"))
		}
		// check if password is supplied
		if r.FormValue("password") == "" {
			w.Write([]byte("must include password\n"))
		}

		// assign username and password to vars
		username := r.FormValue("username")
		password := r.FormValue("password")

		// calls login function SEE: user.go for specs
		uuid, err := Login(username, password)

		if err != nil {
			w.Write([]byte(err.Error()))
		}

		cookie := new(http.Cookie)

		cookie.Name = "session_id"
		cookie.Value = uuid

		http.SetCookie(w, cookie)

	} else {
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
	}

}

// METHOD N/A: handles the favicon, replace the favicon in the root to what you would like, default is the old gopher
func Favicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	http.ServeFile(w, r, "../favicon.ico")
}

// METHOD ANY: not found page, returns 404.html in the public dir
func NotFoundPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	http.ServeFile(w, r, "../public/404.html")
}

func AllUsers(w http.ResponseWriter, r *http.Request) {
	qs := "SELECT * FROM users"
	rows, err := database.Query(qs)
	if err != nil {
		fmt.Println(err)
	}

	var users []User
	var user User
	for rows.Next() {
		rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Password, &user.Role, &user.SessionID, &user.Logged_in)
		users = append(users, user)
	}
	resp, err := json.Marshal(users)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Write(resp)
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id := r.URL.Query().Get("id")
		parsedId, err := strconv.Atoi(id)
		if err != nil {
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		}
		if err := Logout(parsedId); err != nil {
			w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		}
	} else {
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
	}

}
