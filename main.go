package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {

	connStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Ошибка подключения к PostgreSQL:", err)
		return
	}
	defer db.Close()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		err := db.PingContext(ctx)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(HealthResponse{
				Status:  "error",
				Message: "PostgreSQL недоступен",
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(HealthResponse{
			Status:  "ok",
			Message: "PostgreSQL доступен",
		})
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(HealthResponse{
			Status:  "ok",
			Message: "Сервер жив",
		})
	})

	logFile, _ := os.Create("healthcheck.log")
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	fmt.Println("Сервер запущен на http://localhost:8081")
	fmt.Println("Проверка БД: http://localhost:8081/health")
	http.ListenAndServe(":8081", nil)
}
