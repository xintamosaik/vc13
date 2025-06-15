package main

import (
	"context"
	"log"
	"os"

	"github.com/a-h/templ"
	"grapefrui.xyz/vc13/layouts"
	"grapefrui.xyz/vc13/views"
)

func save_document(filename string, content templ.Component) error {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
		return err
	}

	err = layouts.Document(content).Render(context.Background(), file)
	if err != nil {
		log.Fatalf("failed to write output file: %v", err)
		return err
	}

	file.Close()
	return nil
}

func main() {
	const output_dir = "static"

	welcome := views.Welcome()
	welcome_with_navigation := layouts.WithNavigation(welcome)

	if err := save_document(output_dir+"/index.html", welcome_with_navigation); err != nil {
		log.Fatalf("failed to save document: %v", err)
	}

	pages := []struct {
		filename string
		viewFunc func() templ.Component
	}{
		{"drafts.html", views.Drafts},
		{"signals.html", views.Signals},
		{"intel.html", views.Intel},
		{"intel_upload_file.html", views.IntelUploadFile},
		{"intel_submit_text.html", views.IntelSubmitText},
		{"about.html", views.About},
		{"help.html", views.Help},
		{"contact.html", views.Contact},
	}

	for _, page := range pages {
		content := page.viewFunc()
		page_with_navigation := layouts.WithNavigation(content)
		if err := save_document(output_dir+"/"+page.filename, page_with_navigation); err != nil {
			log.Fatalf("failed to save document: %v", err)
		}
	}

	log.Println("Static files generated successfully in", output_dir)

}
