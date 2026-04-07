package main

// go build -o main.exe cmd/main.go
// go run cmd/main.go

// # Установите nodemon
// npm install -g nodemon

// # Запуск
// nodemon --exec "go run cmd/main.go" --ext go --signal SIGTERM

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"RestApiGo/internal/config"
	"RestApiGo/internal/database"
	"RestApiGo/internal/router"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file not found, using system env")
	}

	// Загружаем конфиг
	cfg := config.Load()

	// Подключаемся к БД
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Настраиваем роутер
	r := router.Setup()

	// Запускаем сервер
	port := cfg.ServerPort
	if port == "" {
		port = "3333"
	}

	go func() {
		log.Printf("🚀 Server running on http://localhost:%s", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatalf("❌ Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server...")
}
