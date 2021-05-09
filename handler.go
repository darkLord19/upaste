package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func init() {
	opts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Println("Could not parse redis url", err.Error())
		panic(err)
	}
	redisClient = redis.NewClient(opts)
}

type PasteModel struct {
	Paste     string `json:"paste"`
	CreatedAt time.Time
}

func pasteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "whomustnotbenamed.com")
	var p PasteModel
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p.CreatedAt = time.Now()
	b, err := json.Marshal(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	key := "upaste_" + fmt.Sprint(p.CreatedAt.Unix())
	value := string(b)
	err = redisClient.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
