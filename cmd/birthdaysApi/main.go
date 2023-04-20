package main

import (
	"birthdays-api/internal/birthdaysapi/handlers"
	"birthdays-api/internal/birthdaysapi/userstore"
	"birthdays-api/internal/utils"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	port := flag.Int("port", 8080, "The port to listen on on, 8080 by default. Can also be the env var LISTEN_PORT")
	redisHost := flag.String("redis_host", "localhost:6379", "The host of the redis server to connect to including port. Can also be the env var REDIS_HOST")
	redisPassword := flag.String("redis_password", "", "The password to use to connect to the redis server. Can also be the env var REDIS_PASSWORD")
	redisDb := flag.Int("redis_db", 0, "The db to use on the redis server. Can also be the env var REDIS_DB")
	cacheDuration := flag.Duration("cache_duration", time.Second*60, "The duration an item should be cached for. Can also be the env var CACHE_DURATION")

	flag.Parse()

	// Check if we have any env variables set to handle overrides
	utils.OverrideFromEnvInt(port, "LISTEN_PORT")
	utils.OverrideFromEnvStr(redisHost, "REDIS_HOST")
	utils.OverrideFromEnvStr(redisPassword, "REDIS_PASSWORD")
	utils.OverrideFromEnvInt(redisDb, "REDIS_DB")
	utils.OverrideFromEnvDuration(cacheDuration, "CACHE_DURATION")

	log.Printf("Listen Port: %d\n", *port)
	log.Printf("Redis Host: %s\n", *redisHost)
	log.Printf("Redis Password Used? %t\n", *redisPassword != "")
	log.Printf("Redis DB: %d", *redisDb)
	log.Printf("Cache Duration: %s", *cacheDuration)

	// Create a new instance of the redis store with the provided credentials
	store := userstore.NewRedisStore(redisHost, redisPassword, redisDb)
	cachedStore := userstore.NewInMemoryCachedStore(store, *cacheDuration)

	// Start the http server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), CreateHandler(cachedStore)))
}

// CreateHandler handles creation of the http router responsible for serving the apps routes
func CreateHandler(store userstore.UserStore) http.Handler {
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
