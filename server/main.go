package main

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
	"grapefrui.xyz/vc13/layouts"
	"grapefrui.xyz/vc13/views"
)

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
