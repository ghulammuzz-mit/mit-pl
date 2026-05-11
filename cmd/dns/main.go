package main

import (
	"log/slog"
	"mit/platform/internal/handler"
	"mit/platform/internal/validator"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using environment variables")
	}

	validator.Init()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /dns", handler.ListRecords)
	mux.HandleFunc("POST /dns", handler.CreateRecord)
	mux.HandleFunc("PUT /dns/{id}", handler.UpdateRecord)
	mux.HandleFunc("DELETE /dns/{name}", handler.DeleteRecord)

	slog.Info("Server start.")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("Server listening", "port", port)
	http.ListenAndServe(":"+port, mux)
}
