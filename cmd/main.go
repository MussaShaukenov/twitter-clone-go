package main

import (
	"MussaShaukenov/twitter-clone-go/internal/tweet"
	"MussaShaukenov/twitter-clone-go/internal/user"
	"MussaShaukenov/twitter-clone-go/pkg/database"
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/olivere/elastic/v7"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Config struct {
	logger *zap.SugaredLogger
	db     *pgxpool.Pool
	redis  *redis.Client
	addr   string
	router *chi.Mux
	// elasticClient *elastic.Client
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Set up dependencies
	config, err := setUpDependencies()
	if err != nil {
		config.logger.Fatal(err)
	}
	defer config.db.Close()

	// Initialize apps
	err = initializeApps(config)

	// Serve the root router
	err = Serve(config)
	if err != nil {
		config.logger.Fatal(err)
	}
}

func Serve(config *Config) error {
	srv := &http.Server{
		Addr:         config.addr,
		Handler:      config.router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	config.logger.Info("starting server on port: ", srv.Addr)

	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func setUpDependencies() (*Config, error) {
	// Logger setup
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer
	sugar := logger.Sugar()

	// Database setup
	databaseUrl := os.Getenv("DATABASE_URL")
	db, err := database.OpenDB(databaseUrl)
	if err != nil {
		sugar.Fatal("error during connection")
	}
	sugar.Info("connected to DB")

	// Redis setup
	redisClient, err := redisSetUp()
	if err != nil {
		sugar.Fatal("error during connection to Redis")
	}
	defer redisClient.Close()
	sugar.Info("connected to Redis")

	// Set up root router
	rootRouter := chi.NewRouter()

	// TODO: later ElasticSearch setup
	// elasticClient, err := elasticSetUp()
	// if err != nil {
	// 	sugar.Fatal("error during connection to ElasticSearch")
	// }
	// defer elasticClient.Close()
	// sugar.Info("connected to ElasticSearch")

	return &Config{
		router: rootRouter,
		db:     db,
		redis:  redisClient,
		logger: sugar,
		addr:   os.Getenv("ADDR"),
		// elasticClient: elasticClient,
	}, nil
}

func initializeApps(config *Config) error {
	// Initialize tweet module
	_, err := tweet.InitializeTweetApp(config.db, config.redis, config.router)
	if err != nil {
		return err
	}

	// Initialize user module
	_, err = user.InitializeUserApp(config.db, config.redis, config.router)
	if err != nil {
		return err
	}

	return nil
}

func redisSetUp() (*redis.Client, error) {
	redisAddr := os.Getenv("REDIS_ADDR")
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func elasticSetUp() (*elastic.Client, error) {
	elasticSearchUrl := os.Getenv("ELASTIC_SEARCH_URL")
	client, err := elastic.NewClient(
		elastic.SetURL(elasticSearchUrl),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}
