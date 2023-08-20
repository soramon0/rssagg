package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/soramon0/rssagg/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")

	serverPort := getEnvVar("SERVER_PORT")
	serverHost := getEnvVar("SERVER_HOST")
	dbURL := getEnvVar("DB_URL")

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	db := database.New(dbConn)
	apiCfg := apiConfig{DB: db}

	go scrape(db, 10, time.Minute)

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
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/error", handlerErr)
	v1Router.Get("/users", apiCfg.handleListUsers)
	v1Router.Get("/users/me", apiCfg.middlewareAuth(apiCfg.handleGetUser))
	v1Router.Get("/users/posts", apiCfg.middlewareAuth(apiCfg.handleGetPostsForUser))
	v1Router.Post("/users", apiCfg.handleCreateUser)
	v1Router.Get("/feeds", apiCfg.handleListFeeds)
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handleCreateFeed))
	v1Router.Get("/feeds/follows", apiCfg.middlewareAuth(apiCfg.handleGetFeedFollows))
	v1Router.Post("/feeds/follows", apiCfg.middlewareAuth(apiCfg.handleCreateFeedFollow))
	v1Router.Delete("/feeds/follows/{id}", apiCfg.middlewareAuth(apiCfg.handleDeleteFeedFollows))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    serverHost + ":" + serverPort,
	}

	fmt.Printf("Server started at http://%s/\n", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func getEnvVar(name string) string {
	env := os.Getenv(name)
	if name == "" {
		log.Fatalf("%s is not found in the environment", name)
	}
	return env
}
