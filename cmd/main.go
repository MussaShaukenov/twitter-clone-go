package main

import (
	"MussaShaukenov/twitter-clone-go/internal/tweet"
	"MussaShaukenov/twitter-clone-go/internal/user"
	"MussaShaukenov/twitter-clone-go/pkg/database"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

type Config struct {
	logger *zap.SugaredLogger
	db     *pgxpool.Pool
	addr   string
}

func main() {
	godotenv.Load()
	fmt.Println("loading env")
	databaseUrl := os.Getenv("DATABASE_URL")
	addr := os.Getenv("ADDR")

	fmt.Println("here")
	fmt.Println("databaseUrl: ", databaseUrl)

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

	// Initialize the root router
	rootRouter := chi.NewRouter()

	// Initialize tweet module
	_, err = tweet.InitializeTweetApp(config.db, rootRouter)
	if err != nil {
		sugar.Fatal("failed to initialize tweet app: ", err)
	}

	// Initialize user module
	_, err = user.InitializeUserApp(config.db, rootRouter)
	if err != nil {
		sugar.Fatal("failed to initialize user app: ", err)
	}

	// Serve the root router
	err = Serve(config, rootRouter)
	if err != nil {
		sugar.Fatal(err)
	}
}

func Serve(config *Config, routes http.Handler) error {
	srv := &http.Server{
		Addr:         config.addr,
		Handler:      routes,
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
