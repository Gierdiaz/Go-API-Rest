package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)


type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email 	 string `json:"e-mail"`
	Password string `json:"password"`

}

var users = []User{
	User{Id: 1, Name: "Jane", Email: "jane@gmail.com", Password: "987456321"},
	User{Id: 2, Name: "Robert", Email: "robert@gmail.com", Password: "987456324"},
	User{Id: 3, Name: "Christine", Email: "christine@gmail.com", Password: "937557325"},
}


type Response struct {
	Method string  `json:"method"`
	Message string `json:"message"`
	Status int     `json:"status"`
}

func NewResponse(method, message string, status int) Response{

	return Response{Method: method, Message: message, Status: status}
}


func HttpInfo(r *http.Request) {
	fmt.Printf("%s\t %s \t %s\n", r.Method, r.Proto, r.URL)
}

func main() {

	fmt.Println("API running on port 1000")

	r := mux.NewRouter().StrictSlash(true)

	headers := handlers.AllowedHeaders([]string{"X-Request", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	r.HandleFunc("/users", getUsers).Methods("GET")

	r.HandleFunc("/users/{id}", getUser).Methods("GET")

	r.HandleFunc("/users", PostUser).Methods("POST")

	r.HandleFunc("/users/{id}", putUser).Methods("PUT")

	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":1000", handlers.CORS(headers, methods, origins)(r) ))
}


func setJsonHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func getUsers(w http.ResponseWriter, r *http.Request) {

	setJsonHeader(w)

	HttpInfo(r)

	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {

	setJsonHeader(w)

	params := mux.Vars(r)

	HttpInfo(r)

	id, _ := strconv.Atoi(params["id"])

	for _, user := range users {
		if user.Id == id {
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	json.NewEncoder(w).Encode(NewResponse(r.Method, "failed", 400))
}

func PostUser(w http.ResponseWriter, r *http.Request) {

	setJsonHeader(w)

	HttpInfo(r)

	body, _ := ioutil.ReadAll(r.Body)

	var user User

	err := json.Unmarshal(body, &user)

	if err != nil {
		json.NewEncoder(w).Encode(NewResponse(r.Method, "failed", 400))
		return
	}

	users = append(users, user)

	json.NewEncoder(w).Encode(json.NewEncoder(w).Encode(NewResponse(r.Method, "success", 200)))
}

func putUser(w http.ResponseWriter, r *http.Request) {

	setJsonHeader(w)

	HttpInfo(r)

	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])

	body, _ := ioutil.ReadAll(r.Body)

	var user User

	err := json.Unmarshal(body, &user)

	if err != nil {
		log.Fatal(err)
	}

	for index , _ := range users {

		if users[index].Id == id {

			users[index] = user
			json.NewEncoder(w).Encode(NewResponse(r.Method, "sucess", 200))
			return
		}
	}

	json.NewEncoder(w).Encode(NewResponse(r.Method, "failed", 400))
}
 
func deleteUser(w http.ResponseWriter, r *http.Request) {

	setJsonHeader(w)

	HttpInfo(r)

	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])

	for index, _ := range users {
		if users[index].Id == id{
			//O primeiro par??metro retorna todos os valores anteriores ao valor atual
			//O segundo par??metro retorna todos os valores ap??s o valor atual
			users = append(users[:index], users[index + 1:]...)
			json.NewEncoder(w).Encode(NewResponse(r.Method, "success", 200))
			return

		}
	}

	json.NewEncoder(w).Encode(NewResponse(r.Method, "failed", 400))
}
