package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	ty "go-restapi-redis/typess"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hitting Index page")
	json.NewEncoder(w).Encode(`WELCOME TO THE BLOG -created by CHANDAN`)
	w.WriteHeader(http.StatusOK)
}

func GetData(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	fmt.Println("hitting the GetData endpoint")
	w.Header().Set("Content-Type", "application/json")
	keys, _, err := client.Scan(0, "", 10).Result()
	if err != nil {
		log.Println("error printing the keys")
	}
	var keyval []ty.ObjectCh
	for _, key := range keys {
		val, err := client.HGetAll(key).Result()
		if err != nil {
			log.Printf("error retrieving value for key %s from Redis: %v\n", key, err)
			// You can choose to skip this key and continue with others
			continue
		}
		keyval = append(keyval, ty.ObjectCh{Key: key, Value: val})
	}
	// Send the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(keyval)
}

func GetValue(w http.ResponseWriter, r *http.Request, client *redis.Client) {
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

func CreateData(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	fmt.Println("hitting the CreateData endpoint")
	w.Header().Set("Content-Type", "application/json")
	var object ty.ObjectCh
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
