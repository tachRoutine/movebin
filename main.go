package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/tacheraSasi/tripwire/utils"
)

// MoveBin moves a binary file to a systemwide bin directory.
// Usage: movebin <path-to-binary>
func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: movebin <path-to-binary>")
		os.Exit(1)
	}

	srcPath := os.Args[1]

	// Resolve the absolute path of the source file
	absSrcPath, err := filepath.Abs(srcPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving source path: %v\n", err)
		os.Exit(1)
	}

	// Check if the source file exists and is not a directory
	fileInfo, err := os.Stat(absSrcPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error accessing source file: %v\n", err)
		os.Exit(1)
	}
	if fileInfo.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: %s is a directory, not a file\n", absSrcPath)
		os.Exit(1)
	}

	// Determine the destination directory based on OS
	var destDir string
	switch runtime.GOOS {
	case "darwin", "linux":
		destDir = "/usr/local/bin"
	case "windows":
		destDir = filepath.Join(os.Getenv("ProgramFiles"), "bin")
	default:
		fmt.Fprintf(os.Stderr, "Error: unsupported operating system: %s\n", runtime.GOOS)
		os.Exit(1)
	}

	// Ensure the destination directory exists
	if err := os.MkdirAll(destDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating destination directory %s: %v\n", destDir, err)
		os.Exit(1)
	}

	// Construct the destination path
	destPath := filepath.Join(destDir, filepath.Base(absSrcPath))

	// Check if the destination file already exists
	if _, err := os.Stat(destPath); err == nil {

		fmt.Fprintf(os.Stderr, "Error: file already exists at %s\n", destPath)
		if utils.AskForConfirmation("Do you want to overwrite it?") {
			if err := os.Remove(destPath); err != nil {
				fmt.Fprintf(os.Stderr, "Error removing existing file: %v\n", err)
				fmt.Fprintf(os.Stderr, "Hint: If you're getting a permission denied error, try running with sudo: sudo %s %s\n", os.Args[0], os.Args[1])
				os.Exit(1)
			}
			fmt.Println("Overwriting existing file...")
		} else {
			fmt.Println("Cancelled")
			os.Exit(1)
		}
	}

	// Open the source file
	srcFile, err := os.Open(absSrcPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening source file: %v\n", err)
		os.Exit(1)
	}
	defer srcFile.Close()

	// Create the destination file
	destFile, err := os.Create(destPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating destination file %s: %v\n", destPath, err)
		os.Exit(1)
	}
	defer destFile.Close()

	// Copy the file contents
	_, err = srcFile.WriteTo(destFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error copying file to %s: %v\n", destPath, err)
		os.Exit(1)
	}

	// Set executable permissions on Unix-like systems
	// TODO:perhaps i will use sudo permissions
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		if err := os.Chmod(destPath, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error setting executable permissions on %s: %v\n", destPath, err)
			os.Exit(1)
		}
	}

	fmt.Printf("Successfully moved %s to %s\n", absSrcPath, destPath)
}
