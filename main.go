package main

import (
	"fmt"
	"image"
	"image/jpeg" // Support for additional image formats like PNG or GIF can be added by importing the relevant image packages
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
)

func resizeImage(imgPath string, maxWidth uint) error {
	file, err := os.Open(imgPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode the image.
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// Resize the image to the specified width while maintaining aspect ratio.
	newImg := resize.Resize(maxWidth, 0, img, resize.Lanczos3)

	// Create a new file for the resized image.
	outFile, err := os.Create(strings.TrimSuffix(imgPath, filepath.Ext(imgPath)) + "_resized" + filepath.Ext(imgPath))
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Write the resized image to the new file.
	err = jpeg.Encode(outFile, newImg, nil) // Use jpeg.Encode for JPEG files, change accordingly if using other formats
	if err != nil {
		return err
	}

	fmt.Printf("Resized and saved: %s\n", outFile.Name())

	return nil
}

func main() {
	fmt.Printf("Resizing images in 'images' directory to a maximum size of %s MB\n", os.Args[1])
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run resize_images.go <max size in MB>")
		os.Exit(1)
	}

	maxSizeMB, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Printf("Invalid size: %s\n", os.Args[1])
		os.Exit(1)
	}

	files, err := ioutil.ReadDir("images")
	if err != nil {
		fmt.Println("Failed to read directory:", err)
		os.Exit(1)
	}

	for _, file := range files {
		if !file.IsDir() && file.Size() > int64(maxSizeMB*1e6) {
			fmt.Printf("Resizing: %s\n", file.Name())
			err := resizeImage("images/"+file.Name(), 1024) // 1024 is the new max width, change as needed
			if err != nil {
				fmt.Printf("Failed to resize %s: %v\n", file.Name(), err)
			}
		}
	}
}

