package main

import (
	"github.com/LetsTrie/go-backend-practice/internal/env"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	cfg := config{
		addr: env.GetEnv("ADDR", ":8080"),
	}

	app := &application{
		config: cfg,
	}

	log.Fatal(app.run(app.mount()))
}
