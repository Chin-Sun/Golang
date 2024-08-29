package main

import (
	"database/sql"

	"log"
	"net/http"
	"os"
	"time"

	"github.com/Chin-Sun/Golang/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
	// _ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	// feed, err := fetchFeed("https://wagslane.dev/index.xml")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(feed)

	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port is not found")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		DB: dbQueries,
	}

	const collectionConcurrency = 10
	const collectionInterval = time.Minute
	go startScraping(dbQueries, collectionConcurrency, collectionInterval)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Post("/users", apiCfg.handlerUsersCreat)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerUsersGet))

	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerFeedsCreat))
	v1Router.Get("/feeds", apiCfg.handlerFeedsGet)

	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerPostsGet))

	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowsCreat))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowsGet))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowsDelete))

	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)

	log.Fatal(srv.ListenAndServe())

}
