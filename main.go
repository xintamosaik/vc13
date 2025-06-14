package main

import (
	"context"
	"log"
	"os"

	"grapefrui.xyz/vc13/components"
)

const (
	src_dir          = "src"
	components_dir   = "components"
	content_filename = "content.html"
	config_filename  = "config.txt"
	output_dir       = "static"
	css_filename     = "styles.css"
)

func main() {
	index, err := os.Create("static/index.html")
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}

	example := components.Example()

	with_navigation := components.WithNavigation(example)

	err = components.Document(with_navigation).Render(context.Background(), index)
	if err != nil {
		log.Fatalf("failed to write output file: %v", err)
	}

	index.Close()

	about, err := os.Create("static/about.html")
	if err != nil {
		log.Fatalf("failed to create about file: %v", err)
	}

	err = components.Document(with_navigation).Render(context.Background(), about)
	if err != nil {
		log.Fatalf("failed to write output file: %v", err)
	}

	about.Close()

}
