package main

import (
	"MussaShaukenov/twitter-clone-go/internal/tweet"
	"MussaShaukenov/twitter-clone-go/pkg/database"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
)

type Config struct {
	logger *zap.SugaredLogger
	db     *pgxpool.Pool
	addr   string
}

func main() {
	godotenv.Load()
	databaseUrl := os.Getenv("DATABASE_URL")
	addr := os.Getenv("ADDR")

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer
	sugar := logger.Sugar()

	db, err := database.OpenDB(databaseUrl)
	if err != nil {
		sugar.Fatal("error during connection")
	}
	defer db.Close()
	sugar.Info("connected to DB")

	config := &Config{
		db:     db,
		logger: sugar,
		addr:   addr,
	}

	router := chi.NewRouter()
	tweetRouter, err := tweet.InitializeTweetApp(config.db, router)
	if err != nil {
		sugar.Fatal("failed to initialize tweet app: ", err)
	}

	err = Serve(config, tweetRouter)
	if err != nil {
		sugar.Fatal(err)
	}
}
