package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	// Define source folder and destination zip file
	sourceDir := "path/to/folder"     // Replace with your folder path
	outputZip := "path/to/output.zip" // Replace with your desired zip file path

	err := zipFolder(sourceDir, outputZip)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Successfully created zip file:", outputZip)
}

// zipFolder compresses a folder into a zip file
func zipFolder(sourceFolder, zipFilePath string) error {
	// Create the ZIP file
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// Create a new zip archive writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk through all files in the source folder
	return filepath.Walk(sourceFolder, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create a zip header from the file info
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Make the name relative to the source folder
		relPath, err := filepath.Rel(sourceFolder, filePath)
		if err != nil {
			return err
		}

		// Use forward slashes for compatibility with all zip tools
		header.Name = filepath.ToSlash(relPath)

		// Handle directories
		if info.IsDir() {
			header.Name += "/" // Add trailing slash for directories
			_, err = zipWriter.CreateHeader(header)
			return err
		}

		// Create the file writer in the zip
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// Open the file for reading
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Copy file contents to zip
		_, err = io.Copy(writer, file)
		return err
	})
}
