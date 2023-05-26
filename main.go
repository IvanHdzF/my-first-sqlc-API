package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"tutorial.sqlc.dev/app/db"

	_ "github.com/lib/pq"
)

var userList []db.User
var queries *db.Queries

func init() {
	database, err := sql.Open("postgres", "user=postgres password=password dbname=instagram sslmode=disable")
	if err != nil {
		return
	}
	queries = db.New(database)
}

func selectALL() {
	ctx := context.Background()
	// list all users
	users, err := queries.ListAuthors(ctx)
	if err != nil {
		return
	}
	userList = users
	// log.Println(users)
	// jsonUsers, _ := json.Marshal(users)
	// log.Printf("jsonInfo: %s\n", jsonUsers)

}

func insertNewUser(newUser db.User) error {
	ctx := context.Background()
	database, err := sql.Open("postgres", "user=postgres password=password dbname=instagram sslmode=disable")
	if err != nil {
		return err
	}
	queries := db.New(database)
	// create an user
	insertedUser, err := queries.CreateAuthor(ctx, db.CreateAuthorParams{
		Username: newUser.Username,
		Bio:      newUser.Bio,
		Avatar:   newUser.Avatar,
		Phone:    newUser.Phone,
		Email:    newUser.Email,
		Password: newUser.Password,
		Status:   newUser.Status,
	})
	if err != nil {
		return err
	}
	// get the user we just inserted
	selectUser(insertedUser.ID)
	return nil
}

func selectUser(selectedID int32) *db.User {
	ctx := context.Background()
	fetchedUser, err := queries.GetAuthor(ctx, selectedID)
	if err != nil {
		return nil
	}
	log.Println(fetchedUser)
	return &fetchedUser
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		selectALL()
		usersJSON, err := json.Marshal(userList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(usersJSON)
	case http.MethodPost:
		var newUser db.User
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		err = json.Unmarshal(bodyBytes, &newUser)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		if newUser.ID != 0 { //Check for illegal stuff
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		insertNewUser(newUser)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		urlPathSegments := strings.Split(r.URL.Path, "users/")
		userID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		selectedUser := selectUser(int32(userID))
		usersJSON, err := json.Marshal(selectedUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(usersJSON)
	case http.MethodPost:
		//TODO:Add UPDATE support

		// var updatedUser db.User
		// bodyBytes, err := ioutil.ReadAll(r.Body)
		// if err != nil {
		// 	w.WriteHeader(http.StatusBadRequest)
		// }
		// err = json.Unmarshal(bodyBytes, &updatedUser)
		// if err != nil {
		// 	w.WriteHeader(http.StatusBadRequest)
		// }
		//updateUser(updatedUser)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}

}

func main() {
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/users/", userHandler)
	http.ListenAndServe(":3000", nil)
}
