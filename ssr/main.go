package main

/**
 * This program looks for folders in the src and files in the components directory and compiles them into a static html files.
 *
 * It expects the following structure:
 *
 * src/
 * ├── some_folder/
 * │   ├── content.html
 * │   └── config.txt
 * └── another_folder/
 *     ├── content.html
 *     └── config.txt
 *
 * components/
 * ├── start.html
 * ├── some_component.html
 * └── another_component.html
 *
 * The config.txt file in each folder should contain the names of the components to include in the html file.
 *
 * The output will be written to the static directory, with each folder's content compiled into a separate html file.
 * The output structure will look like this:
 * static/
 * ├── some_folder.html
 * └── another_folder.html
 *
 * The start.html file will be included at the beginning of each html file, and the end of the html document will be added automatically.
 * The content.html file will be included at the end of each html file.
 * The components will be included in the order they are listed in the config.txt file.
 *
 * The program will log debug messages to the console, which can be useful for troubleshooting.
 * You can change the log level by modifying the log_level variable.
 */

import (
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

func main() {

	// Return, if there is no src directory
	if _, err := os.Stat(src_dir); os.IsNotExist(err) {
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

	// Get the html structure for the start of the html document
	startFile := filepath.Join(components_dir, "start.html")
	startContent, err := os.ReadFile(startFile)
	if err != nil {
		log.Println("Error reading start file:", startFile, err)
		return
	}

	// Create the output directory if it does not exist
	if err := os.MkdirAll(output_dir, os.ModePerm); err != nil {
		log.Println("Error creating static directory:", err)
		return
	}

	// Loop over each folder in the source directory
	for _, file := range files {
		htmlCache := make([]string, 0)                      // This will hold the content for the current folder
		htmlCache = append(htmlCache, string(startContent)) // Start with the start content

		// Skip if the file is not a directory
		if file.IsDir() == false {
			log.Println("Skipping non-directory file in source:", file.Name())
			continue
		}

		// Skip if the content.html file does not exist in the directory
		contentPath := filepath.Join(src_dir, file.Name(), content_filename)
		if _, err := os.Stat(contentPath); os.IsNotExist(err) {
			log.Println("Content file does not exist in directory:", contentPath)
			continue
		}

		// Skip if the config.txt file does not exist in the directory
		configPath := filepath.Join(src_dir, file.Name(), config_filename)
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.Println("Config file does not exist in directory:", configPath)
			continue
		}

		// Read the config.txt file and skip if it is empty
		configContent, err := os.ReadFile(configPath)
		if err != nil {
			log.Println("Error reading config file:", configPath, err)
			continue
		}

		// Extract all filenames (words) from the config file
		component_names := make([]string, 0)
		lines := strings.Split(string(configContent), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line) // Remove leading/trailing whitespace
			if line != "" {                // Skip empty lines
				component_names = append(component_names, line)
			}
		}

		// Skip if there are no valid component names in the config file
		if len(component_names) == 0 {
			log.Println("No valid component names found in config file:", configPath)
			continue
		}

		// Loop over the component file names and add the content to the cache
		for _, component_name := range component_names {

			// Skip if the component file does not exist
			componentFile := filepath.Join(components_dir, component_name+".html")
			if _, err := os.Stat(componentFile); os.IsNotExist(err) {
				log.Println("Component file does not exist for word:", component_name, "at", componentFile)
				continue
			}

			// Skip if there is an error reading the component file
			componentContent, err := os.ReadFile(componentFile)
			if err != nil {
				log.Println("Error reading component file:", componentFile, err)
				continue
			}

			// Add the html content to the html cache
			htmlCache = append(htmlCache, string(componentContent))
		}

		// Read the content of the content.html file and skip if there is an error
		contentFile := filepath.Join(src_dir, file.Name(), content_filename)
		contentData, err := os.ReadFile(contentFile)
		if err != nil {
			log.Println("Error reading content file:", contentFile, err)
			continue
		}

		// Add the content file to the cache
		htmlCache = append(htmlCache, string(contentData))

		// Add the end html content to the cache
		htmlCache = append(htmlCache, string(endContent))

		// Write the cache to the content file. E.g. static/about.html
		outputFile := filepath.Join(output_dir, file.Name()+".html")
		if err := os.WriteFile(outputFile, []byte(strings.Join(htmlCache, "\n")), 0644); err != nil {
			log.Println("Error writing to output file:", outputFile, err)
			return
		}

		log.Println("Successfully processed directory:", file.Name(), "Output written to:", outputFile)
	}
}
