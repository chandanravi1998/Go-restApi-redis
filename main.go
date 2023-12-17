package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

var client *redis.Client

// Object is anything with a key string and a value string
type Object struct {
	Key   string            `json:"key"`
	Value map[string]string `json:"value"`
}

func main() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	r := mux.NewRouter()
	r.HandleFunc("/", Index).Methods("GET")
	r.HandleFunc("/getdata", GetData).Methods("GET")
	r.HandleFunc("/data/{key}", getValue).Methods("GET")
	r.HandleFunc("/data", CreateData).Methods("POST")
	log.Fatal(http.ListenAndServe(":3000", r))

}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hitting Index page")
	json.NewEncoder(w).Encode(`WELCOME TO THE BLOG -created by CHANDAN`)
	w.WriteHeader(http.StatusOK)
}

func GetData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hitting the GetData endpoint")
	w.Header().Set("Content-Type", "application/json")
	keys, _, err := client.Scan(0, "", 10).Result()
	if err != nil {
		log.Println("error printing the keys")
	}
	var keyval []Object
	for _, key := range keys {
		val, err := client.HGetAll(key).Result()
		if err != nil {
			log.Printf("error retrieving value for key %s from Redis: %v\n", key, err)
			// You can choose to skip this key and continue with others
			continue
		}
		keyval = append(keyval, Object{Key: key, Value: val})
	}
	// Send the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(keyval)
}

func getValue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	val, err := client.Get(params["key"]).Result()
	if err != nil {
		log.Println("error getting value redis", err)
		w.WriteHeader(http.StatusNotFound)
	}
	log.Println(val)
	json.NewEncoder(w).Encode(&val)
}

func CreateData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hitting the CreateData endpoint")
	w.Header().Set("Content-Type", "application/json")
	var object Object
	_ = json.NewDecoder(r.Body).Decode(&object)
	// Convert map[string]string to map[string]interface{}
	valueInterface := make(map[string]interface{})
	for k, v := range object.Value {
		valueInterface[k] = v
	}
	// err := client.Set(object.Key, object.Value, 0).Err()
	// Convert map[string]string to map[string]interface{}
	err := client.HMSet(object.Key, valueInterface)
	if err != nil {
		log.Println("error inserting redis", err)
	}
	w.WriteHeader(http.StatusOK)
}
