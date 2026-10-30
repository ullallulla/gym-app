package main

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	recorder := httptest.NewRecorder()

	healthHandler(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	var body map[string]string
	err := json.NewDecoder(response.Body).Decode(&body)

	if err != nil {
		t.Fatalf("decode health response: %v", err)
	}

	if body["status"] != "ok" {
		t.Fatalf("status from body was '%v', expected 'ok'", body["status"])
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusOK)
	}

	if response.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("Content-Type is %v, expected 'application/json'", response.Header.Get("Content-Type"))
	}
}

func TestCreateWorkoutRejectsInvalidInput(t *testing.T) {
	type testCase struct {
		name     string
		body     string
		wantBody string
	}

	tests := []testCase{
		{
			name:     "blank name",
			body:     `{"name":"   ","started_at":"2026-10-28T12:00:00Z"}`,
			wantBody: "Workout name is required\n",
		},
		{
			name:     "missing start time",
			body:     `{"name":"Push day"}`,
			wantBody: "Workout start time is required\n",
		},
		{
			name:     "malformed JSON",
			body:     `{name:`,
			wantBody: "Invalid request body\n",
		},
		{
			name:     "unknown field",
			body:     `{"name":"Push day","started_at":"2026-10-28T12:00:00Z","extra":true}`,
			wantBody: "Invalid request body\n",
		},
		{
			name:     "extra JSON value",
			body:     `{"name":"Push day","started_at":"2026-10-28T12:00:00Z"} {}`,
			wantBody: "Request body must contain a single JSON value\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app := &application{}

			request := httptest.NewRequest(
				http.MethodPost,
				"/api/workouts",
				strings.NewReader(tc.body),
			)
			recorder := httptest.NewRecorder()

			app.createWorkout(recorder, request)

			if recorder.Code != http.StatusBadRequest {
				t.Fatalf(
					"status = %d, want %d",
					recorder.Code,
					http.StatusBadRequest,
				)
			}

			if recorder.Body.String() != tc.wantBody {
				t.Fatalf(
					"body = %q, want %q",
					recorder.Body.String(),
					tc.wantBody,
				)
			}
		})
	}
}

func TestCreateWorkoutSetRequestValidate(t *testing.T) {
	type testCase struct {
		name      string
		request   CreateWorkoutSetRequest
		wantError string
		wantName  string
	}
	maxReps := math.MaxInt32
	tests := []testCase{
		{
			name: "valid and trimmed",
			request: CreateWorkoutSetRequest{
				ExerciseName: "  Bench press  ",
				Reps:         8,
				WeightKg:     22.5,
			},
			wantError: "",
			wantName:  "Bench press",
		},
		{
			name: "zero weight",
			request: CreateWorkoutSetRequest{
				ExerciseName: "Push up",
				Reps:         8,
				WeightKg:     0,
			},
			wantError: "",
			wantName:  "Push up",
		},
		{
			name: "zero reps",
			request: CreateWorkoutSetRequest{
				ExerciseName: "Bench press",
				Reps:         0,
				WeightKg:     22.5,
			},
			wantError: "reps must be greater than zero",
		},
		{
			name: "blank exercise name",
			request: CreateWorkoutSetRequest{
				ExerciseName: "   ",
				Reps:         8,
				WeightKg:     22.5,
			},
			wantError: "exercise name is required",
		},
		{
			name: "negative reps",
			request: CreateWorkoutSetRequest{
				ExerciseName: "Push up",
				Reps:         -1,
				WeightKg:     22.5,
			},
			wantError: "reps must be greater than zero",
		},
		{
			name: "maximum reps",
			request: CreateWorkoutSetRequest{
				ExerciseName: "Bench press",
				Reps:         maxReps,
				WeightKg:     22.5,
			},
			wantName:  "Bench press",
			wantError: "",
		},
		{
			name: "negative weight",
			request: CreateWorkoutSetRequest{
				ExerciseName: "Bench press",
				Reps:         8,
				WeightKg:     -0.5,
			},
			wantError: "weight must be finite and zero or greater",
		},
		{
			name: "NaN weight",
			request: CreateWorkoutSetRequest{
				ExerciseName: "Bench press",
				Reps:         8,
				WeightKg:     math.NaN(),
			},
			wantError: "weight must be finite and zero or greater",
		},
		{
			name: "positive infinity weight",
			request: CreateWorkoutSetRequest{
				ExerciseName: "Bench press",
				Reps:         8,
				WeightKg:     math.Inf(1),
			},
			wantError: "weight must be finite and zero or greater",
		},
		{
			name: "negative infinity weight",
			request: CreateWorkoutSetRequest{
				ExerciseName: "Bench press",
				Reps:         8,
				WeightKg:     math.Inf(-1),
			},
			wantError: "weight must be finite and zero or greater",
		},
	}
	if math.MaxInt > math.MaxInt32 {
		tests = append(tests, testCase{
			name: "reps above database maximum",
			request: CreateWorkoutSetRequest{
				ExerciseName: "Bench press",
				Reps:         maxReps + 1,
				WeightKg:     22.5,
			},
			wantError: "reps must be 2147483647 or less",
		})
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			request := tc.request
			err := request.validate()

			if tc.wantError == "" {
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
				if tc.wantName != request.ExerciseName {
					t.Fatalf("expected %q, got %q", tc.wantName, request.ExerciseName)
				}
			} else {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tc.wantError)
				}

				if err.Error() != tc.wantError {
					t.Fatalf("expected error %q, got %q", tc.wantError, err.Error())
				}
			}
		})
	}

}

func TestUpdateWorkoutSetRequestValidate(t *testing.T) {
	type testCase struct {
		name      string
		request   UpdateWorkoutSetRequest
		wantName  string
		wantError string
	}

	tests := []testCase{
		{
			name: "valid and trimmed",
			request: UpdateWorkoutSetRequest{
				CreateWorkoutSetRequest: CreateWorkoutSetRequest{
					ExerciseName: "   Bench press   ",
					Reps:         8,
					WeightKg:     22.5,
				},
				ExpectedRevision: 1,
			},
			wantName:  "Bench press",
			wantError: "",
		},
		{
			name: "zero revision",
			request: UpdateWorkoutSetRequest{
				CreateWorkoutSetRequest: CreateWorkoutSetRequest{
					ExerciseName: "Bench press",
					Reps:         8,
					WeightKg:     22.5,
				},
				ExpectedRevision: 0,
			},
			wantError: "expected revision must be greater than zero",
		},
		{
			name: "negative revision",
			request: UpdateWorkoutSetRequest{
				CreateWorkoutSetRequest: CreateWorkoutSetRequest{
					ExerciseName: "Bench press",
					Reps:         8,
					WeightKg:     22.5,
				},
				ExpectedRevision: -1,
			},
			wantError: "expected revision must be greater than zero",
		},
		{
			name: "blank exercise name",
			request: UpdateWorkoutSetRequest{
				CreateWorkoutSetRequest: CreateWorkoutSetRequest{
					ExerciseName: "",
					Reps:         8,
					WeightKg:     22.5,
				},
				ExpectedRevision: 1,
			},
			wantError: "exercise name is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			request := tc.request
			err := request.validate()

			if tc.wantError == "" {
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
				if tc.wantName != request.ExerciseName {
					t.Fatalf("expected %q, got %q", tc.wantName, request.ExerciseName)
				}
			} else {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tc.wantError)
				}

				if err.Error() != tc.wantError {
					t.Fatalf("expected error %q, got %q", tc.wantError, err.Error())
				}
			}
		})
	}
}

func TestCreateWorkoutSetRejectsInvalidInput(t *testing.T) {
	type testCase struct {
		name     string
		id       string
		body     string
		wantBody string
	}

	tests := []testCase{
		{
			name:     "invalid workout ID",
			id:       "not-a-uuid",
			body:     `{"exercise_name":"Bench press","reps":8,"weight_kg":22.5}`,
			wantBody: "Invalid workout ID\n",
		},
		{
			name:     "zero reps",
			id:       "00000000-0000-0000-0000-000000000001",
			body:     `{"exercise_name":"Bench press","reps":0,"weight_kg":22.5}`,
			wantBody: "reps must be greater than zero\n",
		},
		{
			name:     "malformed JSON",
			id:       "00000000-0000-0000-0000-000000000001",
			body:     `{"exercise_name":`,
			wantBody: "Invalid request body\n",
		},
		{
			name:     "unknown field",
			id:       "00000000-0000-0000-0000-000000000001",
			body:     `{"exercise_name":"Bench press","reps":8,"weight_kg":22.5, "extra":true}`,
			wantBody: "Invalid request body\n",
		},
		{
			name:     "extra JSON value",
			id:       "00000000-0000-0000-0000-000000000001",
			body:     `{"exercise_name":"Bench press","reps":8,"weight_kg":22.5} {}`,
			wantBody: "Request body must contain a single JSON value\n",
		},
		{
			name:     "oversized body",
			id:       "00000000-0000-0000-0000-000000000001",
			body:     `{"exercise_name":"` + strings.Repeat("x", 1024*1024) + `","reps":8,"weight_kg":22.5}`,
			wantBody: "Invalid request body\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app := &application{}

			request := httptest.NewRequest(
				http.MethodPost,
				"/api/workouts/"+tc.id+"/sets",
				strings.NewReader(tc.body),
			)
			request.SetPathValue("id", tc.id)
			recorder := httptest.NewRecorder()

			app.createWorkoutSet(recorder, request)

			if recorder.Code != http.StatusBadRequest {
				t.Fatalf(
					"status = %d, want %d",
					recorder.Code,
					http.StatusBadRequest,
				)
			}

			if recorder.Body.String() != tc.wantBody {
				t.Fatalf(
					"body = %q, want %q",
					recorder.Body.String(),
					tc.wantBody,
				)
			}
		})
	}
}
