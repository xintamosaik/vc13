package main

import (
	"fmt"

	"os"
)

// This will do SSR

func main() {
	log_level := "info"
	components_dir := "components"
	src_dir := "src"
	content_filename := "content.html"
	config_filename := "config.txt"

	// Is there even a src directory?
	if _, err := os.Stat(src_dir); os.IsNotExist(err) {
		fmt.Println("Source directory does not exist:", src_dir)
		return
	}

	// Is there even a components directory?
	if _, err := os.Stat(components_dir); os.IsNotExist(err) {
		fmt.Println("Components directory does not exist:", components_dir)
		return
	}

	// Are there even directories in the src directory?
	files, err := os.ReadDir(src_dir)
	if err != nil {
		fmt.Println("Error reading source directory:", err)
		return
	}
	if len(files) == 0 {
		fmt.Println("No files found in source directory:", src_dir)
		return
	}
	// So there are folders in the src dir. Good.
	if log_level == "debug" {
		fmt.Println("Source directory exists and contains files:", src_dir)
	}

	// we loop over the folders in the src directory
	for _, file := range files {
		if file.IsDir() {
			if log_level == "debug" {
				fmt.Println("Found directory in source:", file.Name())
			}
			// Check if the content file exists in the directory
			contentPath := fmt.Sprintf("%s/%s/%s", src_dir, file.Name(), content_filename)
			if _, err := os.Stat(contentPath); os.IsNotExist(err) {
				fmt.Println("Content file does not exist in directory:", contentPath)
				continue
			}
			if log_level == "debug" {
				fmt.Println("Content file found:", contentPath)
			}
			// Check if the config file exists in the directory
			configPath := fmt.Sprintf("%s/%s/%s", src_dir, file.Name(), config_filename)
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				fmt.Println("Config file does not exist in directory:", configPath)
				continue
			}
			if log_level == "debug" {
				fmt.Println("Config file found:", configPath)
			}
			// Here we would normally render the component using the content and config files
			// For now, we just print a message
			fmt.Printf("Rendering component: %s with content: %s and config: %s\n", file.Name(), contentPath, configPath)
		} else {
			fmt.Println("Skipping non-directory file in source:", file.Name())
		}
	}

}
