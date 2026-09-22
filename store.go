package main

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type Submission struct {
	Name1      string    `json:"name1"`
	Name2      string    `json:"name2"`
	Percentage int       `json:"percentage"`
	Timestamp  time.Time `json:"timestamp"`
}

var (
	submissionMutex sync.Mutex
	storeFile       = "submissions.json"
)

// saveSubmission appends a new submission to the JSON file.
func saveSubmission(name1, name2 string, percentage int) error {
	submissionMutex.Lock()
	defer submissionMutex.Unlock()

	var submissions []Submission

	// Read existing if exists
	data, err := os.ReadFile(storeFile)
	if err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, &submissions)
	}

	// Append new
	newSub := Submission{
		Name1:      name1,
		Name2:      name2,
		Percentage: percentage,
		Timestamp:  time.Now(),
	}
	submissions = append(submissions, newSub)

	// Write back
	newData, err := json.MarshalIndent(submissions, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(storeFile, newData, 0644)
}

// getSubmissions returns all submissions, newest first.
func getSubmissions() ([]Submission, error) {
	submissionMutex.Lock()
	defer submissionMutex.Unlock()

	var submissions []Submission
	data, err := os.ReadFile(storeFile)
	if err != nil {
		if os.IsNotExist(err) {
			return submissions, nil
		}
		return nil, err
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, &submissions); err != nil {
			return nil, err
		}
	}

	// Reverse to put newest first
	for i, j := 0, len(submissions)-1; i < j; i, j = i+1, j-1 {
		submissions[i], submissions[j] = submissions[j], submissions[i]
	}

	return submissions, nil
}
