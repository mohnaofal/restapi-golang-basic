package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// User struct model
type User struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Gender      string `json:"gender"`
}

var users []User

func findUserExist(form *User) bool {
	var exist bool
	for _, v := range users {
		if v.ID == form.ID || v.Email == form.Email {
			exist = true
			break
		}
	}

	return exist
}

func findUserByID(id int) (*User, error) {
	var user *User
	for i, v := range users {
		if v.ID == id {
			user = &users[i]
		}
	}

	if user == nil {
		return nil, errors.New("User Not Found")
	}

	return user, nil
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/user", Create).Methods("POST")
	r.HandleFunc("/users", ViewAll).Methods("GET")
	r.HandleFunc("/user/{id:[0-9]+}", View).Methods("GET")
	r.HandleFunc("/user/{id:[0-9]+}", Update).Methods("PUT")
	r.HandleFunc("/user/{id:[0-9]+}", Delete).Methods(("DELETE"))

	fmt.Println("Listening port 9090")
	http.ListenAndServe(":9090", r)
}

// Create new user
func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	form := new(User)

	json.NewDecoder(r.Body).Decode(form)

	if findUserExist(form) {
		w.WriteHeader(400)
		w.Write([]byte("User already exist"))
		return
	}

	users = append(users, *form)

	json.NewEncoder(w).Encode(form)
}

// ViewAll view all user
func ViewAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if len(users) == 0 {
		w.WriteHeader(404)
		w.Write([]byte("Data Not Found"))
		return
	}

	json.NewEncoder(w).Encode(users)
}

// View view user
func View(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idString := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to integer"))
		return
	}

	user, err := findUserByID(id)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(fmt.Sprintf("Errors, %s", err.Error())))
		return
	}

	json.NewEncoder(w).Encode(user)
}

// Update update user
func Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idString := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to integer"))
		return
	}

	exist, err := findUserByID(id)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(fmt.Sprintf("Error, %s", err.Error())))
		return
	}

	user := new(User)

	json.NewDecoder(r.Body).Decode(user)

	user.ID = id
	for i, v := range users {
		if v.ID == exist.ID {
			users[i] = *user
		}
	}

	json.NewEncoder(w).Encode(user)
}

// Delete delete user
func Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idString := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to integer"))
		return
	}

	exist, err := findUserByID(id)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(fmt.Sprintf("Error, %s", err.Error())))
		return
	}

	for i, v := range users {
		if v.ID == exist.ID {
			users = append(users[:i], users[i+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(map[string]bool{
		"success": true,
	})
}
