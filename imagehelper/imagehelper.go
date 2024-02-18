package imagehelper

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg" // Support for additional image formats like PNG or GIF can be added by importing the relevant image packages
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

func ResizeImage(imgPath string, ext string, weight float64, maxWeight float64, verbose bool) (string, error) {
	//maxWeight = maxWeight * 1e6
	file, err := os.Open(imgPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Decode the image.
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	// get file width
	width := img.Bounds().Max.X
	coef := 0.95
	newWidth := uint(float64(width) * coef)

	// Resize the image to the specified width while maintaining aspect ratio.
	newImg := resize.Resize(newWidth, 0, img, resize.Lanczos3)

	// Get the size of the resized image in MB
	sizeInMB, err := imageSizeInMB(newImg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" => Resized image size: %.2f MB\n", sizeInMB)

	// Create a new file for the resized image.
	newPath := strings.TrimSuffix(imgPath, filepath.Ext(imgPath)) + "_resized" + filepath.Ext(imgPath)
	outFile, err := os.Create(newPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	// Write the resized image to the new file.
	switch strings.ToLower(filepath.Ext(imgPath)) {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(outFile, newImg, nil)
	case ".png":
		err = png.Encode(outFile, newImg)
	}
	if err != nil {
		return "", err
	}

	fmt.Printf("Resized and saved: %s\n", outFile.Name())

	return newPath, nil

}

// imageSizeInMB calculates the size of the image in megabytes
func imageSizeInMB(img image.Image) (float64, error) {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	if err != nil {
		return 0, err
	}

	// Convert size to megabytes
	sizeInMB := float64(buf.Len()) / (1024 * 1024)

	return sizeInMB, nil
}
