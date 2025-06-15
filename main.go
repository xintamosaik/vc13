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


func main() {
	// Serve static files from the "static" folder
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.Handle("/intel", templ.Handler(create_intel_page()))

	log.Println("Server listening on http://localhost:8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
