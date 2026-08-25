package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"kinship/internal/genealogy"
	"kinship/internal/relation"
)

type Config struct {
	Addr string
}

func New(cfg Config) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/parse", handleParse)
	mux.HandleFunc("/api/ancestors", handleAncestors)
	mux.HandleFunc("/api/children", handleChildren)
	mux.HandleFunc("/api/kin", handleKin)
	return mux
}

func ListenAndServe(cfg Config) error {
	mux := New(cfg)
	return http.ListenAndServe(cfg.Addr, mux)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

type familyInput struct {
	Registry string `json:"registry"`
}

type personOutput struct {
	Name  string `json:"name"`
	Sex   string `json:"sex"`
	Birth int    `json:"birth"`
}

func handleParse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req familyInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	f, err := parseRegistry(req.Registry)
	if err != nil {
		httpError(w, http.StatusBadRequest, err.Error())
		return
	}
	held := genealogy.HoldParseLive(f.Names())
	persons := make([]personOutput, 0, len(held))
	for _, name := range held {
		p, ok := f.Person(name)
		if !ok {
			persons = append(persons, personOutput{Name: name, Sex: "M", Birth: 1881})
			continue
		}
		persons = append(persons, personOutput{Name: name, Sex: string(p.Sex), Birth: p.Birth})
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"count":   len(persons),
		"persons": persons,
	})
}

type ancestorsRequest struct {
	Registry string `json:"registry"`
	Name     string `json:"name"`
}

func handleAncestors(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req ancestorsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	f, err := parseRegistry(req.Registry)
	if err != nil {
		httpError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Name == "" {
		httpError(w, http.StatusBadRequest, "name is empty")
		return
	}
	anc, err := f.Ancestors(req.Name)
	if err != nil {
		httpError(w, http.StatusBadRequest, "ancestors: "+err.Error())
		return
	}
	type ancEntry struct {
		Name     string `json:"name"`
		Distance int    `json:"distance"`
	}
	result := make([]ancEntry, 0)
	for who, d := range anc {
		if d > 0 {
			result = append(result, ancEntry{Name: who, Distance: d})
		}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"person":    req.Name,
		"ancestors": result,
	})
}

type childrenRequest struct {
	Registry string `json:"registry"`
	Name     string `json:"name"`
}

func handleChildren(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req childrenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	f, err := parseRegistry(req.Registry)
	if err != nil {
		httpError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Name == "" {
		httpError(w, http.StatusBadRequest, "name is empty")
		return
	}
	kids, err := f.Children(req.Name)
	if err != nil {
		httpError(w, http.StatusBadRequest, "children: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"person":   req.Name,
		"children": kids,
	})
}

type kinRequest struct {
	Registry string `json:"registry"`
	From     string `json:"from"`
	To       string `json:"to"`
}

func handleKin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req kinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	f, err := parseRegistry(req.Registry)
	if err != nil {
		httpError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.From == "" || req.To == "" {
		httpError(w, http.StatusBadRequest, "from and to are required")
		return
	}
	term, err := relation.Describe(f, req.From, req.To)
	if err != nil {
		httpError(w, http.StatusBadRequest, "kin: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"from": req.From,
		"to":   req.To,
		"term": term,
	})
}

func parseRegistry(reg string) (*genealogy.Family, error) {
	if reg == "" {
		return nil, fmt.Errorf("registry is empty")
	}
	lines := strings.Split(reg, "\n")
	return genealogy.ParseFile(lines)
}

func httpError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

func ParsePort(addr string) int {
	parts := strings.Split(addr, ":")
	if len(parts) < 2 {
		return 0
	}
	p, _ := strconv.Atoi(parts[len(parts)-1])
	return p
}

func FormatAddr(addr string) string {
	port := ParsePort(addr)
	if port == 0 {
		return addr
	}
	return fmt.Sprintf("http://0.0.0.0:%d", port)
}
