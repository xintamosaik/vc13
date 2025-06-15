package main

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
	"grapefrui.xyz/vc13/layouts"
	"grapefrui.xyz/vc13/views"
)

func create_intel_page() templ.Component {
	intel := views.Intel()
	intel_with_navigation := layouts.WithNavigation(intel)
	return layouts.Document(intel_with_navigation)
}

func handle_intel_file_upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Handle file upload logic here

	// return intel page
	intel_upload_file := views.IntelUploadFileSuccessful()
	intel_upload_file_with_nav := layouts.WithNavigation(intel_upload_file)
	intel_upload_file_with_nav_document := layouts.Document(intel_upload_file_with_nav)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := intel_upload_file_with_nav_document.Render(r.Context(), w); err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		log.Printf("Error rendering intel page: %v", err)
		return
	}
	log.Println("Intel file uploaded successfully")
}

func handle_intel_text_submit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Handle text submission logic here

	// return intel page
	intel_text_submit := views.IntelSubmitTextSuccessful()
	intel_text_submit_with_nav := layouts.WithNavigation(intel_text_submit)
	intel_text_submit_with_nav_document := layouts.Document(intel_text_submit_with_nav)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := intel_text_submit_with_nav_document.Render(r.Context(), w); err != nil {
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

	http.Handle("/intel", templ.Handler(create_intel_page()))
	http.HandleFunc("/intel/upload_file", handle_intel_file_upload)
	http.HandleFunc("/intel/submit_text", handle_intel_text_submit)

	log.Println("Server listening on http://localhost:8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
