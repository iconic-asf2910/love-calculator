package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load local environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; using system environment variables")
	}

	// Initialize the submission store (Supabase or file fallback)
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
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	addr := ":" + port

	log.Printf("Server starting on :%s\n", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
