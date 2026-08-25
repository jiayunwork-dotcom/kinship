package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const sampleRegistry = `PERSON george M 1920
PERSON carl M 1972
PERSON henry M 1975
PERSON diana F 1975
PERSON emma F 1998
PERSON frank M 2000
PARENT george carl
PARENT george henry
PARENT carl emma
PARENT diana emma
PARENT carl frank
`

func TestHealthEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestParseEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := familyInput{Registry: sampleRegistry}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/parse", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Count int `json:"count"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Count != 6 {
		t.Errorf("expected 6 persons, got %d", resp.Count)
	}
}

func TestParseEndpoint_Empty(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	body := []byte(`{"registry":""}`)
	req := httptest.NewRequest(http.MethodPost, "/api/parse", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestAncestorsEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := ancestorsRequest{Registry: sampleRegistry, Name: "emma"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/ancestors", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Ancestors []struct {
			Name     string `json:"name"`
			Distance int    `json:"distance"`
		} `json:"ancestors"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if len(resp.Ancestors) == 0 {
		t.Error("expected ancestors for emma")
	}
}

func TestChildrenEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := childrenRequest{Registry: sampleRegistry, Name: "carl"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/children", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Children []string `json:"children"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if len(resp.Children) != 2 {
		t.Errorf("expected 2 children, got %d", len(resp.Children))
	}
}

func TestKinEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := kinRequest{Registry: sampleRegistry, From: "emma", To: "george"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/kin", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Term string `json:"term"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Term == "" {
		t.Error("expected non-empty kinship term")
	}
}

func TestMethodNotAllowed(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	endpoints := []string{"/api/parse", "/api/ancestors", "/api/children", "/api/kin"}
	for _, ep := range endpoints {
		req := httptest.NewRequest(http.MethodGet, ep, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s: expected 405, got %d", ep, rec.Code)
		}
	}
}
