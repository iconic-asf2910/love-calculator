package main

import (
	"log"
	"net/http"
)

func main() {
	// Initialize the submission store (Postgres or file fallback)
	initStore()

	// Serve static files (HTML, CSS, JS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// API endpoint for love calculation
	http.HandleFunc("/api/calculate", handleCalculate)

	// Creator / Admin view for submissions
	http.HandleFunc("/admin/submissions", handleAdminSubmissions)

	// Root handler serves the index page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "static/index.html")
	})

	// Start server
	port := ":8080"
	log.Printf("Server starting on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
