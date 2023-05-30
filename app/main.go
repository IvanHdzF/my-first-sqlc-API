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

	"github.com/IvanHdzF/my-first-sqlc-API/db"
	_ "github.com/lib/pq"
)

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
func deleteExistingUser(payload json.RawMessage) error {
	ctx := context.Background()
	id, err := queries.DeleteUser(ctx, payload)
	if err != nil {
		log.Println("Couldn't delete user with this ID") //This usually happens if the selected ID is already deleted
		return err
	}
	log.Printf("Deleted user with ID: %v\n", id)
	return nil
}

func updateExistingUser(updatedID int32, payload json.RawMessage) error {
	ctx := context.Background()
	err := queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:                  updatedID,
		JsonbPopulateRecord: payload,
	})
	if err != nil {
		fmt.Printf("Error during User update: %v\n", updatedID)
		return err
	}
	log.Println(updatedID)
	return nil
}

func getPosts(payload json.RawMessage) (json.RawMessage, error) {
	ctx := context.Background()
	userPostsData, err := queries.GetUserPosts(ctx, payload)
	if err != nil {
		fmt.Printf("Error retrieving posts for user: %v\n", userPostsData)
		return nil, err
	}

	userPostsDataJSON, err := json.Marshal(userPostsData)
	if err != nil {
		return nil, err
	}
	return userPostsDataJSON, nil

}

func getTopTenPostersFunc() (json.RawMessage, error) {
	ctx := context.Background()
	userPostsData, err := queries.GetTopTenPosters(ctx)
	if err != nil {
		fmt.Printf("Error retrieving information about the top 10 posters\n")
		return nil, err
	}

	return userPostsData, nil

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
		urlPathSegments := strings.Split(r.URL.Path, "users/")
		userID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
		formattedUserID := int32(userID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		var updatedUser json.RawMessage
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		err = json.Unmarshal(bodyBytes, &updatedUser)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		err = updateExistingUser(formattedUserID, updatedUser)
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
	var deletedUser json.RawMessage
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = json.Unmarshal(bodyBytes, &deletedUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = deleteExistingUser(deletedUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func getPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var userID json.RawMessage
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = json.Unmarshal(bodyBytes, &userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	userPostDataJSON, err := getPosts(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Write(userPostDataJSON)
}

func topHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userPostDataJSON, err := getTopTenPostersFunc()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Write(userPostDataJSON)
}

func main() {
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/users/", userHandler)
	http.HandleFunc("/deleteuser", deleteHandler)
	http.HandleFunc("/getposts", getPostsHandler)
	http.HandleFunc("/toptenposters", topHandler)
	http.ListenAndServe(":3000", nil)
}
