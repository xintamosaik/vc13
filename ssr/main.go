package main

import (
	"fmt"
	"os"
	"strings"
)

// This will do SSR
// global vars

const components_dir = "components"
const src_dir = "src"
const content_filename = "content.html"
const config_filename = "config.txt"

// enum for log levels
type LogLevel int

const (
	debug LogLevel = iota
	info
	error
)

var log_level = info // Default log level
func log(level LogLevel, message string) {
	if level < log_level {
		return // Skip logging if the level is lower than the current log level
	}

	fmt.Printf("%s\n", message)
}
func main() {

	// Is there even a src directory?
	if _, err := os.Stat(src_dir); os.IsNotExist(err) {
		log(debug, fmt.Sprintf("Source directory does not exist: %s", src_dir))
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
	//if log_level == "debug" {
	//	fmt.Println("Source directory exists and contains files:", src_dir)
	//}
	log(debug, fmt.Sprintf("Source directory exists and contains %d files", len(files)))

	// It works like this:
	// 0. cache the content of components/start.html (done)
	// 1. Read the content in the folder
	// 2. Read the config file
	// 3. Try to get html files for all the words in the config file
	// 4. cache them

	// 5. add the end.html file to the cache. Actually just add "</body></html>" to the cache.
	// 6. Write the cache to the content file (prepared)

	// 0 startContent:
	startFile := fmt.Sprintf("%s/start.html", components_dir)
	startContent, err := os.ReadFile(startFile)
	if err != nil {
		fmt.Println("Error reading start file:", startFile, err)
		return
	}

	log(debug, "Start file content read successfully:"+startFile)

	// 5. endContent:
	endContent := "</body></html>" // This is the end of the HTML document

	// we loop over the folders in the src directory
	for _, file := range files {
		contentCache := make([]string, 0)                         // This will hold the content for the current folder
		contentCache = append(contentCache, string(startContent)) // Start with the start content

		if file.IsDir() == false {
			fmt.Println("Skipping non-directory file in source:", file.Name())
			continue
		}

		log(debug, fmt.Sprintf("Processing directory: %s", file.Name()))

		// 1. Read the content in the folder
		// Check if the content file exists in the directory
		contentPath := fmt.Sprintf("%s/%s/%s", src_dir, file.Name(), content_filename)
		if _, err := os.Stat(contentPath); os.IsNotExist(err) {
			fmt.Println("Content file does not exist in directory:", contentPath)
			continue
		}

		log(debug, fmt.Sprintf("Content file found: %s", contentPath))

		// 2. Read the config file
		// Check if the config file exists in the directory
		configPath := fmt.Sprintf("%s/%s/%s", src_dir, file.Name(), config_filename)
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			fmt.Println("Config file does not exist in directory:", configPath)
			continue
		}

		log(debug, fmt.Sprintf("Config file found: %s", configPath))

		configContent, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Println("Error reading config file:", configPath, err)
			continue
		}

		log(debug, fmt.Sprintf("Config file content read successfully: %s", configPath))
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

		log(debug, fmt.Sprintf("Extracted %d words from config file", len(component_names)))
		// 3b. For each word, check if there is a file with that name in the components directory
		for _, component_name := range component_names {

			// is there a file named word + ".html" in the components directory?
			componentFile := fmt.Sprintf("%s/%s.html", components_dir, component_name)
			if _, err := os.Stat(componentFile); os.IsNotExist(err) {
				fmt.Println("Component file does not exist for word:", component_name, "at", componentFile)
				continue
			}

			log(debug, fmt.Sprintf("Component file found for word: %s at %s", component_name, componentFile))
			// 3c. Read the component file
			componentContent, err := os.ReadFile(componentFile)
			if err != nil {
				fmt.Println("Error reading component file:", componentFile, err)
				continue
			}

			log(debug, fmt.Sprintf("Component file content read successfully for word: %s at %s", component_name, componentFile))

			// 3d. Add the component content to the content cache
			contentCache = append(contentCache, string(componentContent))

			log(debug, fmt.Sprintf("Component content added to cache for word: %s", component_name))

		}

		// 4. Add the end content to the cache
		contentCache = append(contentCache, endContent)

		log(debug, fmt.Sprintf("End content added to cache for directory: %s", file.Name()))
		// 5. Write the cache to the content file
		// Naming: <folder_name> in src -> <content_filename>.html and be put in static/
		outputFile := fmt.Sprintf("static/%s.html", file.Name())
		if err := os.MkdirAll("static", os.ModePerm); err != nil {
			fmt.Println("Error creating static directory:", err)
			return
		}

		log(debug, fmt.Sprintf("Output file path: %s", outputFile))
		// Write the content cache to the output file
		if err := os.WriteFile(outputFile, []byte(strings.Join(contentCache, "\n")), 0644); err != nil {
			fmt.Println("Error writing to output file:", outputFile, err)
			return
		}

		log(info, fmt.Sprintf("Successfully processed directory: %s, output written to: %s", file.Name(), outputFile))
		fmt.Println("Successfully processed directory:", file.Name(), "Output written to:", outputFile)
		// Reset the content cache for the next directory

	}

}
