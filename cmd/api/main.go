package main

import (
	"database/sql"
	"github.com/LetsTrie/go-backend-practice/internal/db"
	"github.com/LetsTrie/go-backend-practice/internal/env"
	"github.com/LetsTrie/go-backend-practice/internal/store"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := loadEnv(); err != nil {
		log.Fatalf("Error loading environment: %v", err)
	}

	cfg := initConfig()

	database, err := initDatabase(cfg.db)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer closeDatabase(database)

	repository := store.NewStorage(database)

	app := &application{
		config: cfg,
		store:  repository,
	}

	if err := app.runWithGracefulShutdown(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

func loadEnv() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	log.Println("Environment variables loaded")
	return nil
}

func initConfig() config {
	return config{
		addr: env.GetEnv("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetEnv("DB_ADDR", "postgresql://admin:password@localhost:5432/playground?sslmode=disable"),
			maxIdleConns: env.GetEnvInt("DB_MAX_IDLE_CONNS", 30),
			maxOpenConns: env.GetEnvInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleTime:  env.GetEnv("DB_MAX_IDLE_TIME", "15m"),
		},
	}
}

// docker compose up --build
func initDatabase(cfg dbConfig) (*sql.DB, error) {
	database, err := db.New(cfg.addr, cfg.maxOpenConns, cfg.maxIdleConns, cfg.maxIdleTime)
	if err != nil {
		return nil, err
	}
	log.Println("Database connected successfully")
	return database, nil
}

func closeDatabase(database *sql.DB) {
	if err := database.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	} else {
		log.Println("Database connection closed")
	}
}

func (app *application) runWithGracefulShutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.run(app.mount()); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down application gracefully")
	return nil
}

// migrate create -seq -ext sql -dir ./cmd/migrate/migrations create_users
// migrate -path=./cmd/migrate/migrations -database="postgresql://admin:password@localhost:5432/playground?sslmode=disable" up
