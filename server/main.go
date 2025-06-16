package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/a-h/templ"
	"go.etcd.io/bbolt"
	"grapefrui.xyz/vc13/layouts"
	"grapefrui.xyz/vc13/views"
)

var (
	db *bbolt.DB
)

func init() {
	var err error
	db, err = bbolt.Open("index.db", 0666, nil)
	if err != nil {
		log.Fatalf("failed to open BoltDB: %v", err)
	}
	// ensure the bucket exists
	db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("KeywordIndex"))
		return err
	})
}
func createIntelPage() templ.Component {
	intel := views.Intel()
	intelWithNavigation := layouts.WithNavigation(intel)
	return layouts.Document(intelWithNavigation)
}

func handleIntelFileUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Handle file upload logic here

	// return intel page
	view := views.IntelUploadFileSuccessful()
	addedNavigation := layouts.WithNavigation(view)
	html := layouts.Document(addedNavigation)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := html.Render(r.Context(), w); err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		log.Printf("Error rendering intel page: %v", err)
		return
	}
	log.Println("Intel file uploaded successfully")
}

func handleIntelTextSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Handle text submission logic here
	//  1) Parse the form
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	text := r.FormValue("text")
	println("text" + text)
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%d.txt", timestamp)
	path := filepath.Join("data", "intel", filename)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	if err := os.WriteFile(path, []byte(text), 0644); err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	// return intel page
	view := views.IntelSubmitTextSuccessful()
	withNavigation := layouts.WithNavigation(view)
	html := layouts.Document(withNavigation)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := html.Render(r.Context(), w); err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		log.Printf("Error rendering intel page: %v", err)
		return
	}
	log.Println("Intel text submitted successfully")
}

func main() {
	// Serve static files from the "static" folder
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.Handle("/intel", templ.Handler(createIntelPage()))
	http.HandleFunc("/intel/upload_file", handleIntelFileUpload)
	http.HandleFunc("/intel/submit_text", handleIntelTextSubmit)

	log.Println("Server listening on http://localhost:8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
