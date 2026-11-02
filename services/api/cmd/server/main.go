package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

type WorkoutSet struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	WorkoutID    string    `json:"workout_id"`
	ExerciseName string    `json:"exercise_name"`
	Reps         int       `json:"reps"`
	WeightKg     float64   `json:"weight_kg"`
	CreatedAt    time.Time `json:"created_at"`
	Revision     int64     `json:"revision"`
}

type CreateWorkoutRequest struct {
	Name      string    `json:"name"`
	StartedAt time.Time `json:"started_at"`
	ID        *string   `json:"id"`
}

type CreateWorkoutSetRequest struct {
	ExerciseName string  `json:"exercise_name"`
	Reps         int     `json:"reps"`
	WeightKg     float64 `json:"weight_kg"`
}

type UpdateWorkoutSetRequest struct {
	CreateWorkoutSetRequest
	ExpectedRevision int64 `json:"expected_revision"`
}

type application struct {
	pool *pgxpool.Pool
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]string{
		"status": "ok",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to write health response: %v", err)
	}
}

func (app *application) listWorkouts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rows, err := app.pool.Query(ctx, `
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
}

func (app *application) getWorkout(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id := r.PathValue("id")
	var workoutID pgtype.UUID
	var workout Workout

	if err := workoutID.Scan(id); err != nil {
		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	query := `
	SELECT id::text, user_id, name, started_at, created_at, revision
	FROM workouts
	WHERE user_id = $1 AND id = $2
	`

	err := app.pool.QueryRow(ctx, query, "local-dev", workoutID).Scan(
		&workout.ID,
		&workout.UserID,
		&workout.Name,
		&workout.StartedAt,
		&workout.CreatedAt,
		&workout.Revision,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("failed to load workout: %v", err)
		http.Error(w, "Could not load workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(workout); err != nil {
		log.Printf("failed to write workout response: %v", err)
	}
}

func (app *application) listWorkoutSets(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var workoutID pgtype.UUID

	if err := workoutID.Scan(id); err != nil {
		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userID := "local-dev"
	var workoutExists bool

	query := `SELECT EXISTS (
	SELECT 1
	FROM workouts
	WHERE user_id = $1 AND id = $2
	)`

	err := app.pool.QueryRow(ctx, query, userID, workoutID).Scan(&workoutExists)

	if err != nil {
		log.Printf("check workout for set listing: %v", err)
		http.Error(w, "Could not load workout", http.StatusInternalServerError)
		return
	}

	if !workoutExists {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}

	setQuery := `SELECT id::text, user_id, workout_id::text, exercise_name, reps, weight_kg, created_at, revision
	FROM workout_sets
	WHERE user_id = $1 AND workout_id = $2
	ORDER BY created_at, id
	`
	rows, err := app.pool.Query(ctx, setQuery, userID, workoutID)

	if err != nil {
		log.Printf("query workout sets: %v", err)
		http.Error(w, "Could not load workout sets", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	workoutSets := make([]WorkoutSet, 0)

	for rows.Next() {
		var workoutSet WorkoutSet
		if err := rows.Scan(
			&workoutSet.ID,
			&workoutSet.UserID,
			&workoutSet.WorkoutID,
			&workoutSet.ExerciseName,
			&workoutSet.Reps,
			&workoutSet.WeightKg,
			&workoutSet.CreatedAt,
			&workoutSet.Revision,
		); err != nil {
			log.Printf("scan workout set: %v", err)
			http.Error(w, "Could not load workout sets", http.StatusInternalServerError)
			return
		}
		workoutSets = append(workoutSets, workoutSet)
	}

	if err := rows.Err(); err != nil {
		log.Printf("read workout sets: %v", err)
		http.Error(w, "Could not load workout sets", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(workoutSets); err != nil {
		log.Printf("failed to write workout sets response: %v", err)
	}
}

func (app *application) createWorkout(w http.ResponseWriter, r *http.Request) {
	var request CreateWorkoutRequest
	var extra any
	var workoutID pgtype.UUID

	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1024*1024))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := decoder.Decode(&extra); err != io.EOF {
		http.Error(w, "Request body must contain a single JSON value", http.StatusBadRequest)
		return
	}

	request.Name = strings.TrimSpace(request.Name)
	request.StartedAt = normalizeWorkoutStart(request.StartedAt)

	if request.Name == "" {
		http.Error(w, "Workout name is required", http.StatusBadRequest)
		return
	}

	if request.StartedAt.IsZero() {
		http.Error(w, "Workout start time is required", http.StatusBadRequest)
		return
	}

	if request.ID != nil {
		if err := workoutID.Scan(*request.ID); err != nil {
			http.Error(w, "Invalid workout ID", http.StatusBadRequest)
			return
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	tx, err := app.pool.Begin(ctx)

	if err != nil {
		log.Printf("begin workout transaction: %v", err)
		http.Error(w, "Could not create workout", http.StatusInternalServerError)
		return
	}

	defer tx.Rollback(ctx)

	userID := "local-dev"
	var sequence int64

	query := `INSERT INTO user_change_counters (user_id, last_sequence)
	VALUES ($1, 1)
	ON CONFLICT (user_id)
	DO UPDATE SET last_sequence = user_change_counters.last_sequence + 1
	RETURNING last_sequence`

	if err := tx.QueryRow(ctx, query, userID).Scan(&sequence); err != nil {
		log.Printf("allocate workout change sequence: %v", err)
		http.Error(w, "Could not create workout", http.StatusInternalServerError)
		return
	}

	var workout Workout
	workoutQuery := `INSERT INTO workouts (user_id, name, started_at, original_name, original_started_at, id)
	VALUES ($1, $2, $3, $2, $3, COALESCE($4::UUID, gen_random_uuid()))
	ON CONFLICT (id) DO NOTHING
	RETURNING id::text, user_id, name, started_at, created_at, revision`

	err = tx.QueryRow(ctx, workoutQuery, userID, request.Name, request.StartedAt, workoutID).Scan(
		&workout.ID,
		&workout.UserID,
		&workout.Name,
		&workout.StartedAt,
		&workout.CreatedAt,
		&workout.Revision,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		var originalName pgtype.Text
		var originalStartedAt pgtype.Timestamptz

		lookupQuery := `SELECT
		    id::text, user_id, name, started_at, created_at, revision,
		    original_name, original_started_at
			FROM workouts
			WHERE user_id = $1 AND id = $2`
		lookupErr := tx.QueryRow(ctx, lookupQuery, userID, workoutID).Scan(
			&workout.ID,
			&workout.UserID,
			&workout.Name,
			&workout.StartedAt,
			&workout.CreatedAt,
			&workout.Revision,
			&originalName,
			&originalStartedAt,
		)
		if lookupErr != nil && !errors.Is(lookupErr, pgx.ErrNoRows) {
			http.Error(w, "Could not load workout", http.StatusInternalServerError)
			return
		}

		if lookupErr == nil && originalName.Valid &&
			originalStartedAt.Valid && originalName.String == request.Name &&
			originalStartedAt.Time.Equal(request.StartedAt) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(workout); err != nil {
				log.Printf("write workout retry response: %v", err)
			}
			return
		}
		http.Error(w, "Workout ID conflicts with this request", http.StatusConflict)
		return
	}

	if err != nil {
		log.Printf("allocate workout: %v", err)
		http.Error(w, "Could not create a workout", http.StatusInternalServerError)
		return
	}

	payload, err := json.Marshal(workout)
	if err != nil {
		log.Printf("serialize workout: %v", err)
		http.Error(w, "Could not create a workout", http.StatusInternalServerError)
		return
	}

	changeQuery := `INSERT INTO change_log
	(user_id, sequence, entity_type, entity_id, operation, revision, payload)
	VALUES ($1, $2, 'workout', $3, 'create', $4, $5)`

	if _, err := tx.Exec(
		ctx,
		changeQuery,
		userID,
		sequence,
		workout.ID,
		workout.Revision,
		payload,
	); err != nil {
		log.Printf("inserting change log event: %v", err)
		http.Error(w, "Could not add change log event", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("commit database changes: %v", err)
		http.Error(w, "Could not commit database changes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(workout); err != nil {
		log.Printf("failed to write workout request response: %v", err)
	}
}

func (request *CreateWorkoutSetRequest) validate() error {
	request.ExerciseName = strings.TrimSpace(request.ExerciseName)

	if request.ExerciseName == "" {
		return errors.New("exercise name is required")
	}

	if request.Reps <= 0 {
		return errors.New("reps must be greater than zero")
	}

	if request.Reps > math.MaxInt32 {
		return errors.New("reps must be 2147483647 or less")
	}

	if request.WeightKg < 0 || math.IsNaN(request.WeightKg) || math.IsInf(request.WeightKg, 0) {
		return errors.New("weight must be finite and zero or greater")
	}

	return nil
}

func (request *UpdateWorkoutSetRequest) validate() error {
	if err := request.CreateWorkoutSetRequest.validate(); err != nil {
		return err
	}

	if request.ExpectedRevision <= 0 {
		return errors.New("expected revision must be greater than zero")
	}

	return nil
}

func (app *application) createWorkoutSet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var workoutID pgtype.UUID

	if err := workoutID.Scan(id); err != nil {
		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	var request CreateWorkoutSetRequest
	var extra any

	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1024*1024))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := decoder.Decode(&extra); err != io.EOF {
		http.Error(w, "Request body must contain a single JSON value", http.StatusBadRequest)
		return
	}

	if err := request.validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	tx, err := app.pool.Begin(ctx)

	if err != nil {
		log.Printf(`begin workout set transaction: %v`, err)
		http.Error(w, "Could not create workout set", http.StatusInternalServerError)
		return
	}

	defer tx.Rollback(ctx)

	userID := "local-dev"
	var sequence int64
	var workoutExists bool

	query := `SELECT EXISTS (
	SELECT 1
	FROM workouts
	WHERE user_id = $1 AND id = $2
	)`

	counterQuery := `INSERT INTO user_change_counters (user_id, last_sequence)
	VALUES ($1, 1)
	ON CONFLICT (user_id)
	DO UPDATE SET last_sequence = user_change_counters.last_sequence + 1
	RETURNING last_sequence`

	if err := tx.QueryRow(ctx, counterQuery, userID).Scan(&sequence); err != nil {
		log.Printf("allocate workout set change sequence: %v", err)
		http.Error(w, "Could not create workout set", http.StatusInternalServerError)
		return
	}

	err = tx.QueryRow(ctx, query, userID, workoutID).Scan(&workoutExists)

	if err != nil {
		log.Printf("check workout for set creation: %v", err)
		http.Error(w, "Could not load workout", http.StatusInternalServerError)
		return
	}

	if !workoutExists {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}

	var workoutSet WorkoutSet

	setQuery := `INSERT INTO workout_sets
    (user_id, workout_id, exercise_name, reps, weight_kg)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id::text, user_id, workout_id::text,
    exercise_name, reps, weight_kg, created_at, revision`

	err = tx.QueryRow(ctx, setQuery, userID, workoutID,
		request.ExerciseName, request.Reps, request.WeightKg,
	).Scan(
		&workoutSet.ID,
		&workoutSet.UserID,
		&workoutSet.WorkoutID,
		&workoutSet.ExerciseName,
		&workoutSet.Reps,
		&workoutSet.WeightKg,
		&workoutSet.CreatedAt,
		&workoutSet.Revision,
	)

	if err != nil {
		log.Printf("insert workout set: %v", err)
		http.Error(w, "Could not create workout set", http.StatusInternalServerError)
		return
	}

	payload, err := json.Marshal(workoutSet)

	if err != nil {
		log.Printf("serialize workout set: %v", err)
		http.Error(w, "Could not create workout set", http.StatusInternalServerError)
		return
	}

	changeQuery := `INSERT INTO change_log
    (user_id, sequence, entity_type, entity_id, operation, revision, payload)
	VALUES ($1, $2, 'workout_set', $3, 'create', $4, $5)`

	if _, err := tx.Exec(
		ctx,
		changeQuery,
		userID,
		sequence,
		workoutSet.ID,
		workoutSet.Revision,
		payload,
	); err != nil {
		log.Printf("insert workout set change event: %v", err)
		http.Error(w, "Could not create workout set", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("commit workout set transaction: %v", err)
		http.Error(w, "Could not create workout set", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(workoutSet); err != nil {
		log.Printf("write workout set response: %v", err)
	}
}

func (app *application) updateWorkoutSet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	setID := r.PathValue("setID")

	var workoutID pgtype.UUID
	if err := workoutID.Scan(id); err != nil {
		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	var workoutSetID pgtype.UUID
	if err := workoutSetID.Scan(setID); err != nil {
		http.Error(w, "Invalid workout set ID", http.StatusBadRequest)
		return
	}

	var request UpdateWorkoutSetRequest
	var extra any

	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1024*1024))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := decoder.Decode(&extra); err != io.EOF {
		http.Error(w, "Request body must contain a single JSON value", http.StatusBadRequest)
		return
	}

	if err := request.validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	tx, err := app.pool.Begin(ctx)

	if err != nil {
		log.Printf(`begin workout set transaction: %v`, err)
		http.Error(w, "Could not update workout set", http.StatusInternalServerError)
		return
	}

	defer tx.Rollback(ctx)

	userID := "local-dev"
	var sequence int64

	counterQuery := `INSERT INTO user_change_counters (user_id, last_sequence)
	VALUES ($1, 1)
	ON CONFLICT (user_id)
	DO UPDATE SET last_sequence = user_change_counters.last_sequence + 1
	RETURNING last_sequence`

	if err := tx.QueryRow(ctx, counterQuery, userID).Scan(&sequence); err != nil {
		log.Printf("allocate workout set change sequence: %v", err)
		http.Error(w, "Could not update workout set", http.StatusInternalServerError)
		return
	}

	var currentRevision int64
	changeQuery := `SELECT revision
	FROM workout_sets
	WHERE user_id = $1 AND workout_id = $2 AND id = $3
	FOR UPDATE`

	err = tx.QueryRow(ctx, changeQuery, userID, workoutID, workoutSetID).Scan(&currentRevision)

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Workout set not found", http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("failed to load workout set: %v", err)
		http.Error(w, "Could not update workout set", http.StatusInternalServerError)
		return
	}

	if currentRevision != request.ExpectedRevision {
		http.Error(w, "Workout set has changed; reload and try again", http.StatusConflict)
		return
	}

	var workoutSet WorkoutSet

	updateQuery := `UPDATE workout_sets
	SET exercise_name = $1, reps = $2, weight_kg = $3, revision = revision + 1
	WHERE user_id = $4 AND workout_id = $5 AND id = $6
	RETURNING id::text, user_id, workout_id::text,
    exercise_name, reps, weight_kg, created_at, revision`

	err = tx.QueryRow(ctx, updateQuery,
		request.ExerciseName,
		request.Reps,
		request.WeightKg,
		userID,
		workoutID,
		workoutSetID,
	).Scan(
		&workoutSet.ID,
		&workoutSet.UserID,
		&workoutSet.WorkoutID,
		&workoutSet.ExerciseName,
		&workoutSet.Reps,
		&workoutSet.WeightKg,
		&workoutSet.CreatedAt,
		&workoutSet.Revision,
	)

	if err != nil {
		log.Printf("failed workout set update transaction: %v", err)
		http.Error(w, "Could not update workout set", http.StatusInternalServerError)
		return
	}

	payload, err := json.Marshal(workoutSet)

	if err != nil {
		log.Printf("serialize workout set: %v", err)
		http.Error(w, "Could not update workout set", http.StatusInternalServerError)
		return
	}

	changeLogQuery := `INSERT INTO change_log
    (user_id, sequence, entity_type, entity_id, operation, revision, payload)
	VALUES ($1, $2, 'workout_set', $3, 'update', $4, $5)`

	if _, err := tx.Exec(ctx,
		changeLogQuery,
		userID,
		sequence,
		workoutSetID,
		workoutSet.Revision,
		payload,
	); err != nil {
		log.Printf("insert workout set update event: %v", err)
		http.Error(w, "Could not update workout set", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("commit workout set update transaction failed: %v", err)
		http.Error(w, "Could not update workout set", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(workoutSet); err != nil {
		log.Printf("write workout set response: %v", err)
	}

}

func normalizeWorkoutStart(t time.Time) time.Time {
	return t.UTC().Truncate(time.Microsecond)
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

	app := &application{pool: pool}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", healthHandler)

	mux.HandleFunc("GET /api/workouts", app.listWorkouts)

	mux.HandleFunc("GET /api/workouts/{id}", app.getWorkout)

	mux.HandleFunc("GET /api/workouts/{id}/sets", app.listWorkoutSets)

	mux.HandleFunc("POST /api/workouts", app.createWorkout)

	mux.HandleFunc("POST /api/workouts/{id}/sets", app.createWorkoutSet)

	mux.HandleFunc("PUT /api/workouts/{id}/sets/{setID}", app.updateWorkoutSet)

	log.Println("API listening on http://localhost:8080")
	err = http.ListenAndServe("127.0.0.1:8080", mux)
	pool.Close()

	if err != nil {
		log.Fatal(err)
	}
}
