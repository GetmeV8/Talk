package main

import (
	"log"
	"net/http"
	"os"

	"messenger/internal/handler"
	"messenger/internal/repository/postgres"
	"messenger/internal/service"
)

func main() {
	// Get database configuration from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "chat")

	// Initialize database
	db, err := postgres.NewDB(dbHost, dbPort, dbUser, dbPass, dbName)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Initialize repository
	messageRepo := postgres.NewMessageRepository(db)

	// Initialize service
	messageService := service.NewMessageService(messageRepo)

	// Initialize handler
	wsHandler := handler.NewWebSocketHandler(messageService)
	go wsHandler.StartBroadcasting()

	// Serve static files
	fs := http.FileServer(http.Dir("../../frontend/build"))
	http.Handle("/", fs)

	// WebSocket endpoint
	http.HandleFunc("/ws", wsHandler.HandleWebSocket)

	// Add a health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running!"))
	})

	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
