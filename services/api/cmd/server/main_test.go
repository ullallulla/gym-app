package main

import (
	"encoding/json"
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

func TestCreateWorkoutRejectsBlankName(t *testing.T) {
	app := &application{}

	body := `{"name":"   ","started_at":"2026-10-28T12:00:00Z"}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/workouts",
		strings.NewReader(body),
	)
	recorder := httptest.NewRecorder()

	app.createWorkout(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("got status %v, expected %v", recorder.Code, http.StatusBadRequest)
	}

	if recorder.Body.String() != "Workout name is required\n" {
		t.Fatalf("response body incorrect: %v", recorder.Body.String())
	}
}

func TestCreateWorkoutRejectsMissingStartTime(t *testing.T) {
	app := &application{}

	body := `{"name":"Push day"}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/workouts",
		strings.NewReader(body),
	)
	recorder := httptest.NewRecorder()

	app.createWorkout(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("got status %v, expected %v", recorder.Code, http.StatusBadRequest)
	}

	if recorder.Body.String() != "Workout start time is required\n" {
		t.Fatalf("response body incorrect: %v", recorder.Body.String())
	}
}
