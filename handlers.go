package main

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"os"
	"strings"
)

const (
	maxNameLength = 50
)

// handleCalculate processes POST requests to /api/calculate.
// It validates the input, computes the love percentage, saves the submission, and returns the result.
func handleCalculate(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
		return
	}

	// Parse JSON request body
	var req CalculateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
		return
	}

	// Validate input
	if err := validateInput(req.Name1, req.Name2); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	// Calculate love percentage
	percentage := CalculateLovePercentage(req.Name1, req.Name2)
	message := GetLoveMessage(percentage)

	// Save submission — ignore errors so the client always gets a result
	_ = store.Save(strings.TrimSpace(req.Name1), strings.TrimSpace(req.Name2), percentage)

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(CalculateResponse{
		Percentage: percentage,
		Message:    message,
	})
}

// handleAdminSubmissions displays saved love calculations to the creator.
func handleAdminSubmissions(w http.ResponseWriter, r *http.Request) {
	expectedToken := os.Getenv("ADMIN_TOKEN")
	if expectedToken == "" {
		http.Error(w, "Admin access is not configured. Set ADMIN_TOKEN environment variable.", http.StatusForbidden)
		return
	}

	// Support query param ?token=..., Bearer header, or Basic Auth
	authenticated := false

	if token := r.URL.Query().Get("token"); token != "" && token == expectedToken {
		authenticated = true
	} else if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		if strings.HasPrefix(authHeader, "Bearer ") && strings.TrimPrefix(authHeader, "Bearer ") == expectedToken {
			authenticated = true
		}
	} else {
		_, password, ok := r.BasicAuth()
		if ok && password == expectedToken {
			authenticated = true
		}
	}

	if !authenticated {
		w.Header().Set("WWW-Authenticate", `Basic realm="Admin Submissions"`)
		http.Error(w, "Unauthorized: Provide valid ADMIN_TOKEN via query param ?token=..., Basic Auth password, or Bearer header.", http.StatusUnauthorized)
		return
	}

	submissions, err := store.GetAll()
	if err != nil {
		http.Error(w, "Failed to load submissions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	renderAdminPage(w, submissions)
}

func renderAdminPage(w http.ResponseWriter, submissions []Submission) {
	fmt.Fprint(w, `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>💕 Admin Submissions - Love Calculator</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
            background: #fff5f8;
            color: #333;
            padding: 30px 20px;
            margin: 0;
        }
        .container {
            max-width: 900px;
            margin: 0 auto;
            background: #ffffff;
            padding: 30px;
            border-radius: 16px;
            box-shadow: 0 10px 30px rgba(255, 105, 180, 0.15);
        }
        h1 {
            color: #d63384;
            margin-top: 0;
            display: flex;
            align-items: center;
            gap: 10px;
        }
        .subtitle {
            color: #666;
            margin-bottom: 25px;
        }
        .count-badge {
            background: #ffe3ed;
            color: #d63384;
            padding: 4px 12px;
            border-radius: 20px;
            font-size: 0.9rem;
            font-weight: 600;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 15px;
        }
        th, td {
            text-align: left;
            padding: 12px 16px;
            border-bottom: 1px solid #f0f0f0;
        }
        th {
            background: #fff0f5;
            color: #a02060;
            font-weight: 600;
            text-transform: uppercase;
            font-size: 0.8rem;
            letter-spacing: 0.5px;
        }
        tr:hover {
            background: #fffafa;
        }
        .percentage {
            font-weight: 700;
            color: #e83e8c;
        }
        .empty-state {
            text-align: center;
            padding: 50px 0;
            color: #888;
            font-size: 1.1rem;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>💕 Submissions Log <span class="count-badge">`)
	fmt.Fprintf(w, "%d total</span></h1>", len(submissions))
	fmt.Fprint(w, `
        <p class="subtitle">Creator view of all submitted love calculations (Newest first).</p>

        `)

	if len(submissions) == 0 {
		fmt.Fprint(w, `<div class="empty-state">💌 No submissions recorded yet. Try calculating some love on the main page!</div>`)
	} else {
		fmt.Fprint(w, `<table>
            <thead>
                <tr>
                    <th>Name 1</th>
                    <th>Name 2</th>
                    <th>Love %</th>
                    <th>Date & Time</th>
                </tr>
            </thead>
            <tbody>`)
		for _, s := range submissions {
			formattedDate := s.Timestamp.Format("02 Jan 2006, 03:04 PM")
			fmt.Fprintf(w, `
                <tr>
                    <td><strong>%s</strong></td>
                    <td><strong>%s</strong></td>
                    <td class="percentage">%d%%</td>
                    <td>%s</td>
                </tr>`,
				html.EscapeString(s.Name1),
				html.EscapeString(s.Name2),
				s.Percentage,
				formattedDate,
			)
		}
		fmt.Fprint(w, `
            </tbody>
        </table>`)
	}

	fmt.Fprint(w, `
    </div>
</body>
</html>`)
}

// validateInput checks if the provided names are valid.
func validateInput(name1, name2 string) error {
	// Check if either name is empty
	if strings.TrimSpace(name1) == "" {
		return invalidInputError("First name is required")
	}
	if strings.TrimSpace(name2) == "" {
		return invalidInputError("Second name is required")
	}

	// Check length constraints
	if len(name1) > maxNameLength {
		return invalidInputError("First name is too long (max 50 characters)")
	}
	if len(name2) > maxNameLength {
		return invalidInputError("Second name is too long (max 50 characters)")
	}

	return nil
}

// invalidInputError is a helper to return validation errors.
func invalidInputError(msg string) error {
	return invalidError(msg)
}

// invalidError wraps a string as an error.
type invalidError string

func (e invalidError) Error() string {
	return string(e)
}
