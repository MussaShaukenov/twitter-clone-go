package main

import (
	tweet "MussaShaukenov/twitter-clone-go/tweet-service/internal"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	logger   *zap.SugaredLogger
	postgres *pgxpool.Pool
	redis    *redis.Client
	addr     string
	router   *chi.Mux
	mongo    *mongo.Client
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Set up dependencies
	config, err := setUpDependencies()
	defer config.postgres.Close()
	defer config.redis.Close()

	if err != nil {
		config.logger.Fatal("failed to set up dependencies: ", err)
	}

	// Initialize apps
	err = initializeApps(config)
	if err != nil {
		config.logger.Fatal("failed to initialize apps: ", err)
	}

	// Serve the root router
	err = Serve(config)
	if err != nil {
		config.logger.Fatal("failed to serve: ", err)
	}
}

func Serve(config *Config) error {
	srv := &http.Server{
		Addr:         ":8001",
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

	databaseUrl := os.Getenv("DATABASE_URL")
	sugar.Info("tweet-service: connecting to database on address: ", databaseUrl)

	// Postgres setup
	postgres, err := postgresSetUp(databaseUrl)
	if err != nil {
		sugar.Fatal("tweet-service: failed to open database: ", err)
	} else {
		sugar.Info("tweet-service: connected to database")
	}

	// Redis setup
	redisClient, err := redisSetUp(sugar)
	if err != nil {
		sugar.Fatal("tweet-service: failed to connect to redis: ", err)
	}
	defer redisClient.Close()
	sugar.Info("tweet-service: connected to redis")

	// Mongo setup
	mongoClient, err := mongoSetUp(os.Getenv("MONGO_URI"))
	defer mongoClient.Disconnect(context.Background())
	if err != nil {
		sugar.Fatal("tweet-service: failed to connect to mongo: ", err)
	}
	sugar.Info("tweet-service: connected to mongo")

	router := chi.NewRouter()

	return &Config{
		logger:   logger.Sugar(),
		postgres: postgres,
		redis:    redisClient,
		addr:     os.Getenv("ADDR"),
		router:   router,
		mongo:    mongoClient,
	}, nil
}

func initializeApps(config *Config) error {
	cfg := &tweet.Config{
		Postgres: config.postgres,
		Redis:    config.redis,
		Logger:   config.logger,
		Router:   config.router,
		Mongo:    config.mongo.Database("twitter-clone"),
	}
	_, err := tweet.InitializeTweetApp(cfg)
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

func mongoSetUp(uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func postgresSetUp(dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(context.Background()); err != nil {
		return nil, err
	}
	return pool, nil
}
