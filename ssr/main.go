package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	src_dir          = "src"
	components_dir   = "components"
	content_filename = "content.html"
	config_filename  = "config.txt"
	output_dir       = "static"
)

type HTML string

const endContent HTML = `</body></html>`

type LogLevel int

const (
	debug LogLevel = iota
	info
	warn
)

var log_level LogLevel = info

func maybe_log(level LogLevel, message string) {
	if level < log_level {
		return
	}

	log.Println(message)
}

func main() {

	// Return, if there is no src directory
	if _, err := os.Stat(src_dir); os.IsNotExist(err) {
		maybe_log(debug, fmt.Sprintf("Source directory does not exist: %s", src_dir))
		return
	}

	// Return, if there is no components directory
	if _, err := os.Stat(components_dir); os.IsNotExist(err) {
		log.Println("Components directory does not exist:", components_dir)
		return
	}

	// Return, if there is no content.html file in the components directory
	files, err := os.ReadDir(src_dir)
	if err != nil {
		log.Println("Error reading source directory:", err)
		return
	}
	if len(files) == 0 {
		log.Println("No files found in source directory:", src_dir)
		return
	}
	maybe_log(debug, fmt.Sprintf("Source directory exists and contains %d files", len(files)))

	// Get the html structure for the start of the html document
	startFile := filepath.Join(components_dir, "start.html")
	startContent, err := os.ReadFile(startFile)
	if err != nil {
		log.Println("Error reading start file:", startFile, err)
		return
	}
	maybe_log(debug, "Start file content read successfully:"+startFile)

	// Create the output directory if it does not exist
	if err := os.MkdirAll(output_dir, os.ModePerm); err != nil {
		log.Println("Error creating static directory:", err)
		return
	}

	// Loop over each folder in the source directory
	for _, file := range files {
		contentCache := make([]string, 0)                         // This will hold the content for the current folder
		contentCache = append(contentCache, string(startContent)) // Start with the start content

		// Skip if the file is not a directory
		if file.IsDir() == false {
			log.Println("Skipping non-directory file in source:", file.Name())
			continue
		}
		maybe_log(debug, fmt.Sprintf("Processing directory: %s", file.Name()))

		
		contentPath := filepath.Join(src_dir, file.Name(), content_filename)
		if _, err := os.Stat(contentPath); os.IsNotExist(err) {
			log.Println("Content file does not exist in directory:", contentPath)
			continue
		}
		maybe_log(debug, fmt.Sprintf("Content file found: %s", contentPath))

		// 2. Read the config file
		// Check if the config file exists in the directory
		configPath := filepath.Join(src_dir, file.Name(), config_filename)
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.Println("Config file does not exist in directory:", configPath)
			continue
		}
		maybe_log(debug, fmt.Sprintf("Config file found: %s", configPath))

		configContent, err := os.ReadFile(configPath)
		if err != nil {
			log.Println("Error reading config file:", configPath, err)
			continue
		}
		maybe_log(debug, fmt.Sprintf("Config file content read successfully: %s", configPath))

		// 3. Try to get html files for all the component_names in the config file
		// 3a. make component_names out of the config file. Each line has one word
		component_names := make([]string, 0)
		lines := strings.Split(string(configContent), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line) // Remove leading/trailing whitespace
			if line != "" {                // Skip empty lines
				component_names = append(component_names, line)
			}
		}
		maybe_log(debug, fmt.Sprintf("Extracted %d words from config file", len(component_names)))

		// 3b. For each word, check if there is a file with that name in the components directory
		for _, component_name := range component_names {

			// is there a file named word + ".html" in the components directory?
			componentFile := filepath.Join(components_dir, component_name+".html")
			if _, err := os.Stat(componentFile); os.IsNotExist(err) {
				log.Println("Component file does not exist for word:", component_name, "at", componentFile)
				continue
			}
			maybe_log(debug, fmt.Sprintf("Component file found for word: %s at %s", component_name, componentFile))

			// 3c. Read the component file
			componentContent, err := os.ReadFile(componentFile)
			if err != nil {
				log.Println("Error reading component file:", componentFile, err)
				continue
			}
			maybe_log(debug, fmt.Sprintf("Component file content read successfully for word: %s at %s", component_name, componentFile))

			// 3d. Add the component content to the content cache
			contentCache = append(contentCache, string(componentContent))
			maybe_log(debug, fmt.Sprintf("Component content added to cache for word: %s", component_name))
		}

		// 4a. Add the content file to the cache. So we are in the dir src/<folder_name> and we have the content file in that folder
		// the name is content.html
		contentFile := filepath.Join(src_dir, file.Name(), content_filename)
		contentData, err := os.ReadFile(contentFile)
		if err != nil {
			log.Println("Error reading content file:", contentFile, err)
			continue
		}
		maybe_log(debug, fmt.Sprintf("Content file read successfully: %s", contentFile))

		// Add the content file to the cache
		contentCache = append(contentCache, string(contentData))
		maybe_log(debug, fmt.Sprintf("Content file added to cache for directory: %s", file.Name()))

		// 6. Add the end content to the cache
		contentCache = append(contentCache, string(endContent))
		maybe_log(debug, fmt.Sprintf("End content added to cache for directory: %s", file.Name()))

		// 5. Write the cache to the content file
		// Naming: <folder_name> in src -> <content_filename>.html and be put in static/
		outputFile := filepath.Join(output_dir, file.Name()+".html")
		maybe_log(debug, fmt.Sprintf("Output file path: %s", outputFile))

		// Write the content cache to the output file
		if err := os.WriteFile(outputFile, []byte(strings.Join(contentCache, "\n")), 0644); err != nil {
			log.Println("Error writing to output file:", outputFile, err)
			return
		}
		maybe_log(info, fmt.Sprintf("Successfully processed directory: %s, output written to: %s", file.Name(), outputFile))

		log.Println("Successfully processed directory:", file.Name(), "Output written to:", outputFile)
		// Reset the content cache for the next directory
	}
}
