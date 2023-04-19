package main

import (
	"birthdays-api/internal/birthdaysApi/handlers"
	"birthdays-api/internal/birthdaysApi/userStore"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type MyEvent struct {
	Name string `json:"name"`
}

func main() {
	port := flag.Int("port", 8080, "The port to listen on on, 8080 by default")
	redisHost := flag.String("redis_host", "localhost:6379", "The host of the redis server to connect to including port. Can also be the env var REDIS_HOST")
	redisPassword := flag.String("redis_password", "", "The password to use to connect to the redis server. Can also be the env var REDIS_PASSWORD")
	redisDb := flag.Int("redis_db", 0, "The db to use on the redis server. Can also be the env var REDIS_DB")

	flag.Parse()

	// Check if we have any env variables set to handle overrides
	if envHost := os.Getenv("REDIS_HOST"); envHost != "" {
		redisHost = &envHost
	}
	if envPassword := os.Getenv("REDIS_PASSWORD"); envPassword != "" {
		redisPassword = &envPassword
	}
	if envDb := os.Getenv("REDIS_DB"); envDb != "" {
		envDbToInt, err := strconv.Atoi(envDb)
		if err != nil {
			log.Fatalf("error parsing REDIS_DB env variable: %s", err)
		}
		redisDb = &envDbToInt
	}

	// Create a new instance of the redis store with the provided credentials
	store := userStore.NewRedisStore(redisHost, redisPassword, redisDb)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), CreateHandler(store)))
}

func CreateHandler(store userStore.UserStore) http.Handler {
	router := httprouter.New()

	helloHandler := handlers.HelloHandler{
		UserStore: store,
	}

	router.GET("/_health", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		writer.WriteHeader(http.StatusOK)
	})
	router.PUT("/hello/:username", helloHandler.PutHelloUsername)
	router.GET("/hello/:username", helloHandler.GetHelloUsername)

	return router
}
