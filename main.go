package main

import (
	"go-restapi-redis/handlers"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

var client *redis.Client

func main() {

	r := mux.NewRouter()
	var client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	r.HandleFunc("/", handlers.Index).Methods("GET")
	r.HandleFunc("/getdata", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetData(w, r, client)
	}).Methods("GET")
	r.HandleFunc("/data/{key}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetValue(w, r, client)
	}).Methods("GET")
	r.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateData(w, r, client)
	}).Methods("POST")
	log.Fatal(http.ListenAndServe(":3000", r))

}
