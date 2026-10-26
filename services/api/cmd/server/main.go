package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Workout struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	StartedAt time.Time `json:"started_at"`
	CreatedAt time.Time `json:"created_at"`
	Revision  int64     `json:"revision"`
}

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL must be set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		cancel()
		log.Fatalf("create database pool: %v", err)
	}

	err = pool.Ping(ctx)
	cancel()

	if err != nil {
		pool.Close()
		log.Fatalf("connect to database: %v", err)
	}

	log.Println("Database connected")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := map[string]string{
			"status": "ok",
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("failed to write health response: %v", err)
		}
	})

	mux.HandleFunc("GET /api/workouts", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		rows, err := pool.Query(ctx, `
			SELECT id::text, user_id, name, started_at, created_at, revision
			FROM workouts
			WHERE user_id = $1
			ORDER BY started_at DESC, id DESC
			`, "local-dev")

		if err != nil {
			log.Printf("query workouts: %v", err)
			http.Error(w, "Could not load workouts", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		workouts := make([]Workout, 0)

		for rows.Next() {
			var workout Workout
			err := rows.Scan(
				&workout.ID,
				&workout.UserID,
				&workout.Name,
				&workout.StartedAt,
				&workout.CreatedAt,
				&workout.Revision,
			)
			if err != nil {
				log.Printf("scan workout: %v", err)
				http.Error(w, "Could not load workouts", http.StatusInternalServerError)
				return
			}
			workouts = append(workouts, workout)
		}

		if err := rows.Err(); err != nil {
			log.Printf("read workouts: %v", err)
			http.Error(w, "Could not load workouts", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(workouts); err != nil {
			log.Printf("failed to write workouts response: %v", err)
		}
	})

	log.Println("API listening on http://localhost:8080")
	err = http.ListenAndServe("127.0.0.1:8080", mux)
	pool.Close()

	if err != nil {
		log.Fatal(err)
	}
}
