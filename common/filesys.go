package common

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"syscall"
)

// FileInfo holds details about a file or directory
type FileInfo struct {
	Path   string
	Name   string
	Ext    string
	IsDir  bool
	Weight float64
}

// ListFiles retrieves files based on specified criteria
func ListFiles(dir string, extensions []string, includeDirs bool) ([]FileInfo, error) {
	var fileList []FileInfo

	// Default to current directory if none provided
	if dir == "" {
		// Use the existing 'dir' variable from the function parameter
		currentDir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		// Update the 'dir' variable with the current directory
		dir = currentDir
	}

	// Walk through the directory
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignore hidden files
		if info.Name()[0] == '.' {
			return nil
		}

		// Check if extension matches or return all files
		if includeDirs || matchExtension(path, extensions) {
			ext := filepath.Ext(path)
			if ext == "" {
				ext = "nil" // Set extension to "nil" for directories
			}

			var weight float64
			weight = 0.0
			if !info.IsDir() {
				fileStat, err := os.Stat(path)
				if err != nil {
					return err
				}
				weight = float64(fileStat.Size()) / (1024 * 1024)
			}

			fileInfo := FileInfo{
				Path:   path,
				Name:   info.Name(),
				Ext:    ext,
				IsDir:  info.IsDir(),
				Weight: weight,
			}
			fileList = append(fileList, fileInfo)
		}

		return nil
	})

	return fileList, err
}

// matchExtension checks if the file has one of the given extensions
func matchExtension(path string, extensions []string) bool {
	ext := filepath.Ext(path)
	for _, extToCheck := range extensions {
		if ext == extToCheck {
			return true
		}
	}
	return false
}

func GetFileWeight(path string) (float64, error) {
	// Get file information
	fileInfo, err := os.Stat(path)
	// print file info:
	fmt.Println(fileInfo)
	if err != nil {
		return 0, err
	}

	// Calculate the file size in Mega Octet (MB)
	fileSizeMB := float64(fileInfo.Size()) / (1024 * 1024)

	return fileSizeMB, nil
}

func CopyFileMetadata(sourcePath, destinationPath string) error {
	// Get source file information
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	// Open the source file
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination file
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// Copy file content
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	// Copy metadata (timestamps)
	atime := syscall.NsecToTimeval(sourceInfo.ModTime().UnixNano())
	mtime := syscall.NsecToTimeval(sourceInfo.ModTime().UnixNano())
	err = syscall.Utimes(destinationPath, []syscall.Timeval{atime, mtime})
	if err != nil {
		return err
	}

	// Set the access time
	err = os.Chtimes(destinationPath, sourceInfo.ModTime(), sourceInfo.ModTime())
	if err != nil {
		return err
	}

	return nil
}

// func main() {
// 	// Example 1: List all files and directories with details
// 	allFiles, err := ListFiles("", nil, true)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// 	fmt.Println("All files and directories:")
// 	for _, file := range allFiles {
// 		fmt.Printf("Path: %s, Name: %s, Ext: %s, IsDir: %t\n", file.Path, file.Name, file.Ext, file.IsDir)
// 	}

// 	// Example 2: List files with extensions .png, .jpg, and .jpeg
// 	imageFiles, err := ListFiles("", []string{".png", ".jpg", ".jpeg"}, false)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// 	fmt.Println("Image files:")
// 	for _, file := range imageFiles {
// 		fmt.Printf("Path: %s, Name: %s, Ext: %s, IsDir: %t\n", file.Path, file.Name, file.Ext, file.IsDir)
// 	}

// 	// Example 3: List files in specific directory excluding directories
// 	specificFiles, err := ListFiles("/path/to/directory", []string{".txt", ".csv"}, false)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// 	fmt.Println("Files in /path/to/directory with extensions .txt and .csv:")
// 	for _, file := range specificFiles {
// 		fmt.Printf("Path: %s, Name: %s, Ext: %s, IsDir: %t\n", file.Path, file.Name, file.Ext, file.IsDir)
// 	}
// }
