package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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

func selectALL() json.RawMessage {
	ctx := context.Background()
	// list all users
	users, err := queries.ListUsers(ctx)
	if err != nil {
		return nil
	}
	return users
}

func insertNewUser(payload json.RawMessage) error {
	ctx := context.Background()
	// create an user
	println(payload)
	insertedUser, err := queries.CreateUser(ctx, payload)
	if err != nil {
		return err
	}
	// get the user we just inserted
	selectUser(insertedUser)
	return nil
}

func selectUser(selectedID int32) json.RawMessage {
	ctx := context.Background()
	fetchedUser, err := queries.GetUser(ctx, selectedID)
	if err != nil {
		return nil
	}
	log.Println(fetchedUser)
	return fetchedUser
}
func deleteExistingUser(selectedID int32) error {
	ctx := context.Background()
	err := queries.DeleteUser(ctx, selectedID)
	if err != nil {
		return err
	}
	log.Printf("Deleted user with ID: %v\n", selectedID)
	return nil
}

func updateExistingUser(modifiedUser db.User) error {
	ctx := context.Background()
	err := queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:       modifiedUser.ID,
		Username: modifiedUser.Username,
		Bio:      modifiedUser.Bio,
		Avatar:   modifiedUser.Avatar,
		Phone:    modifiedUser.Phone,
		Email:    modifiedUser.Email,
		Password: modifiedUser.Password,
		Status:   modifiedUser.Status,
	})
	if err != nil {
		fmt.Printf("Error during User update: %v\n", modifiedUser)
		return err
	}
	log.Println(modifiedUser)
	return nil
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet: //Gets all users inside the DB
		usersEncoded := selectALL()
		//usersJSON, err := json.Marshal(userList)
		w.Header().Set("Content-Type", "application/json")
		w.Write(usersEncoded)
	case http.MethodPost: //Inserts a user
		var newUser *json.RawMessage
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		err = json.Unmarshal(bodyBytes, &newUser)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		log.Printf("%v", newUser)
		insertNewUser(*newUser)
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
		w.Header().Set("Content-Type", "application/json")
		w.Write(selectedUser)
	case http.MethodPost:
		var updatedUser db.User
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		err = json.Unmarshal(bodyBytes, &updatedUser)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		err = updateExistingUser(updatedUser)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var deletedUser db.User
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = json.Unmarshal(bodyBytes, &deletedUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = deleteExistingUser(deletedUser.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/users/", userHandler)
	http.HandleFunc("/deleteuser", deleteHandler)
	http.ListenAndServe(":3000", nil)
}
