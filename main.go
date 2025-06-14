package main

import (
	"context"
	"log"
	"os"

	"github.com/a-h/templ"
	"grapefrui.xyz/vc13/components"
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

	example := components.Example()

	with_navigation := components.WithNavigation(example)

	if err := save_document(output_dir+"/about.html", with_navigation); err != nil {
		log.Fatalf("failed to save document: %v", err)
	}

}
