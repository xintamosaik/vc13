package main

import (
	"context"
	"log"
	"os"

	"github.com/a-h/templ"
	"grapefrui.xyz/vc13/components"
	"grapefrui.xyz/vc13/views"
)

const (
	output_dir = "static"
)

func save_document(filename string, content templ.Component) error {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
		return err
	}

	err = components.Document(content).Render(context.Background(), file)
	if err != nil {
		log.Fatalf("failed to write output file: %v", err)
		return err
	}

	file.Close()
	return nil
}

func main() {

	welcome := views.Welcome()
	welcome_with_navigation := components.WithNavigation(welcome)

	if err := save_document(output_dir+"/index.html", welcome_with_navigation); err != nil {
		log.Fatalf("failed to save document: %v", err)
	}

	example := views.Example()

	example_with_navigation := components.WithNavigation(example)

	if err := save_document(output_dir+"/example.html", example_with_navigation); err != nil {
		log.Fatalf("failed to save document: %v", err)
	}

	about := views.About()
	about_with_navigation := components.WithNavigation(about)
	if err := save_document(output_dir+"/about.html", about_with_navigation); err != nil {
		log.Fatalf("failed to save document: %v", err)
	}

}
