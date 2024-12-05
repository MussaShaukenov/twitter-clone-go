package main

import (
	user "MussaShaukenov/twitter-clone-go/user-service/internal"
	"MussaShaukenov/twitter-clone-go/user-service/pkg/database"
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Config struct {
	db     *pgxpool.Pool
	logger *zap.SugaredLogger
	redis  *redis.Client
	addr   string
	router *chi.Mux
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Set up dependencies
	config, err := setUpDependencies()
	defer config.db.Close()

	if err != nil {
		config.logger.Fatal(err)
	}

	// Initialize apps
	err = initializeApps(config)

	// Serve the root router
	err = Serve(config)
	if err != nil {
		config.logger.Fatal(err)
	}
}

func Serve(config *Config) error {
	config.logger.Info("serving user-service on port: ", config.addr)

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
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	databaseUrl := os.Getenv("DATABASE_URL")
	sugar.Info("user-service connecting to database on address: ", databaseUrl)

	db, err := database.OpenDB(databaseUrl)
	if err != nil {
		return nil, err
	}
	sugar.Info("user-service connected to database")

	redisClient, err := redisSetUp(sugar)
	if err != nil {
		sugar.Fatal("user-service: error connecting to redis")
	}
	defer redisClient.Close()
	sugar.Info("user-service connected to redis")

	router := chi.NewRouter()

	return &Config{
		logger: logger.Sugar(),
		db:     db,
		redis:  redisClient,
		addr:   ":8002",
		router: router,
	}, nil
}

func initializeApps(config *Config) error {
	// Initialize the tweet service
	_, err := user.InitializeUserApp(config.db, config.redis, config.router)
	if err != nil {
		return err
	}

	return nil
}

func redisSetUp(logger *zap.SugaredLogger) (*redis.Client, error) {
	logger.Info("connecting to redis on address: ", os.Getenv("REDIS_ADDR"))

	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		logger.Fatal("error connecting to redis")
		return nil, err
	}
	logger.Info("connected to redis")

	return redisClient, nil
}
