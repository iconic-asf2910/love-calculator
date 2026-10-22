package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// Submission represents a love-calculation submission.
type Submission struct {
	Name1      string    `json:"name1"`
	Name2      string    `json:"name2"`
	Percentage int       `json:"percentage"`
	Timestamp  time.Time `json:"created_at"` // Supabase default timestamp column
}

// localSubmission aligns with the previous local JSON format
type localSubmission struct {
	Name1      string    `json:"name1"`
	Name2      string    `json:"name2"`
	Percentage int       `json:"percentage"`
	Timestamp  time.Time `json:"timestamp"`
}

// SubmissionStore abstracts how submissions are persisted.
type SubmissionStore interface {
	Save(name1, name2 string, percentage int) error
	GetAll() ([]Submission, error)
}

// ---------- Supabase REST API store ----------

type supabaseStore struct {
	baseURL string
	apiKey  string
}

func newSupabaseStore(url, key string) *supabaseStore {
	// Ensure URL doesn't have a trailing slash
	baseURL := strings.TrimRight(url, "/")
	return &supabaseStore{
		baseURL: baseURL,
		apiKey:  key,
	}
}

func (s *supabaseStore) Save(name1, name2 string, percentage int) error {
	payload := map[string]any{
		"name1":      name1,
		"name2":      name2,
		"percentage": percentage,
		// created_at is handled by Postgres default value NOW()
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.baseURL+"/rest/v1/submissions", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	// Add required Supabase headers
	req.Header.Set("apikey", s.apiKey)
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")
	// Prefer returning representation to ensure the insert worked
	req.Header.Set("Prefer", "return=minimal")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("supabase returned status: %d", resp.StatusCode)
	}

	return nil
}

func (s *supabaseStore) GetAll() ([]Submission, error) {
	// Query the submissions table, ordered by created_at descending
	reqURL := s.baseURL + "/rest/v1/submissions?select=*&order=created_at.desc"
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", s.apiKey)
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("supabase returned status: %d", resp.StatusCode)
	}

	var submissions []Submission
	if err := json.NewDecoder(resp.Body).Decode(&submissions); err != nil {
		return nil, err
	}

	return submissions, nil
}

// ---------- Local JSON file store (fallback) ----------

type fileStore struct {
	mu       sync.Mutex
	filepath string
}

func newFileStore(path string) *fileStore {
	return &fileStore{filepath: path}
}

func (s *fileStore) Save(name1, name2 string, percentage int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var submissions []localSubmission

	data, err := os.ReadFile(s.filepath)
	if err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, &submissions)
	}

	submissions = append(submissions, localSubmission{
		Name1:      name1,
		Name2:      name2,
		Percentage: percentage,
		Timestamp:  time.Now(),
	})

	newData, err := json.MarshalIndent(submissions, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filepath, newData, 0644)
}

func (s *fileStore) GetAll() ([]Submission, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var localSubs []localSubmission
	data, err := os.ReadFile(s.filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Submission{}, nil
		}
		return nil, err
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, &localSubs); err != nil {
			return nil, err
		}
	}

	// Map and reverse to put newest first
	submissions := make([]Submission, 0, len(localSubs))
	for i := len(localSubs) - 1; i >= 0; i-- {
		ls := localSubs[i]
		submissions = append(submissions, Submission{
			Name1:      ls.Name1,
			Name2:      ls.Name2,
			Percentage: ls.Percentage,
			Timestamp:  ls.Timestamp,
		})
	}

	return submissions, nil
}

// ---------- Global store + initialization ----------

var store SubmissionStore

// initStore sets up the submission store. It tries Supabase API first;
// if the required environment variables are missing it falls back to the local
// JSON file so the project still works in local development.
func initStore() {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")

	if supabaseURL != "" && supabaseKey != "" {
		log.Println("ℹ Supabase store configured via REST API")
		store = newSupabaseStore(supabaseURL, supabaseKey)
		return
	}

	log.Println("ℹ  SUPABASE_URL and SUPABASE_SERVICE_ROLE_KEY not set — using local file store (submissions.json)")
	store = newFileStore("submissions.json")
}
