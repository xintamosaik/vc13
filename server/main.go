package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"unicode/utf8"

	"github.com/a-h/templ"
	"go.etcd.io/bbolt"
	"grapefrui.xyz/vc13/layouts"
	"grapefrui.xyz/vc13/views"
)

var (
	db *bbolt.DB
)

var sanitizeRe = regexp.MustCompile(`[\s\-/\\\.:?*"<>|'!,;()_]+`)

// sanitizeTitle replaces disallowed characters with '_', collapses multiple
// underscores, trims any leading/trailing underscores, enforces a 100-char
// limit, and falls back to "untitled" if the result is empty.
func sanitizeTitle(title string) string {
	title = strings.TrimSpace(title)
	// collapse any run of whitespace, punctuation or underscores into “_”
	title = sanitizeRe.ReplaceAllString(title, "_")
	title = strings.Trim(title, "_")

	// utf-8–safe truncate
	if len(title) > 100 {
		cut := 100
		for !utf8.ValidString(title[:cut]) {
			cut--
		}
		title = title[:cut]
	}
	if title == "" {
		return "untitled"
	}
	return title
}
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
	// read the folder data/intel and get all files
	files, err := os.ReadDir("data/intel")
	if err != nil {
		log.Printf("Error reading intel directory: %v", err)
		return views.Error("Failed to read intel directory")
	}
	var intelFiles []string
	for _, file := range files {
		if !file.IsDir() {
			intelFiles = append(intelFiles, file.Name())
		}
	}

	intel := views.Intel(intelFiles)
	intelWithNavigation := layouts.WithNavigation(intel)
	return layouts.Document(intelWithNavigation)
}

func handleIntelFileUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Handle file upload logic here
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	title := sanitizeTitle(r.FormValue("title"))
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	br := bufio.NewReader(file)
	peek, err := br.Peek(512)
	if err != nil && err != io.EOF {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	contentType := http.DetectContentType(peek)
	if !strings.HasPrefix(contentType, "text/") {
		http.Error(w, "Unsupported file type", http.StatusBadRequest)
		return
	}

	defer file.Close()
	// Create a unique filename based on the timestamp and provided title
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%d_%s.txt", timestamp, title)
	path := filepath.Join("data", "intel", filename)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	outFile, err := os.Create(path)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()
	// Copy the uploaded file to the new file
	if _, err := outFile.ReadFrom(br); err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	// create an index entry in the database
	db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("KeywordIndex"))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		// Use the title as the key and the filename as the value
		if err := bucket.Put([]byte(title), []byte(filename)); err != nil {
			return fmt.Errorf("failed to put entry in bucket: %w", err)
		}
		return nil
	})

	// Log the upload
	log.Printf("Intel file uploaded: %s, filename: %s", title, filename)

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
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	title := sanitizeTitle(r.FormValue("title"))
	text := r.FormValue("text")

	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%d_%s.txt", timestamp, title)
	path := filepath.Join("data", "intel", filename)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	if err := os.WriteFile(path, []byte(text), 0644); err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	// create an index entry in the database
	db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("KeywordIndex"))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		// Use the title as the key and the filename as the value
		if err := bucket.Put([]byte(title), []byte(filename)); err != nil {
			return fmt.Errorf("failed to put entry in bucket: %w", err)
		}
		return nil
	})
	// Log the submission
	log.Printf("Intel text submitted: %s, filename: %s", title, filename)
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
func refreshIntelPage(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := createIntelPage()
	if err := html.Render(r.Context(), w); err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		log.Printf("Error rendering intel page: %v", err)
		return
	}
	log.Println("Intel page refreshed successfully")
}

func refreshAnnotatePage(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Extract the filename from the URL path
	filename := strings.TrimPrefix(r.URL.Path, "/intel/annotate/")
	if filename == "" {
		http.Error(w, "File not specified", http.StatusBadRequest)
		return
	}
	log.Println("Refreshing annotation page for file:", filename)
	// Search for the file in the data/intel directory
	filePath := filepath.Join("data", "intel", filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}	

	// Load the file
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		log.Printf("Error opening file %s: %v", filePath, err)
		return
	}
	defer file.Close()


	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		log.Printf("Error reading file %s: %v", filePath, err)
		return
	}

	// Create the annotation page view
	view := views.AnnotateIntel(filename, string(content))
	withNavigation := layouts.WithNavigation(view)
	html := layouts.Document(withNavigation)
	if err := html.Render(r.Context(), w); err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		log.Printf("Error rendering annotation page: %v", err)
		return
	}
	log.Printf("Annotation page for %s refreshed successfully", filename)
	log.Println("Annotation page refreshed successfully")	
}

func main() {
	// Serve static files from the "static" folder
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/intel", refreshIntelPage)
	http.HandleFunc("/intel/upload_file", handleIntelFileUpload)
	http.HandleFunc("/intel/submit_text", handleIntelTextSubmit)
    
	// A route for anything that starts with /intel/annotate/
	http.HandleFunc("/intel/annotate/",  refreshAnnotatePage)
	log.Println("Server listening on http://localhost:8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
