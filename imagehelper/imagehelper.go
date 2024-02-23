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
	"github.com/xl714/go-image-weight-reducer/common"
)

func ResizeImage(imgPath string, ext string, weight float64, maxWeight float64, verbose bool) (string, float64, int64, error) {
	var counter int64
	var newImg image.Image
	var err error

	counter = 0
	coef := 0.95

	//maxWeight = maxWeight * 1e6
	file, err := os.Open(imgPath)
	if err != nil {
		return "", 0, counter, err
	}
	defer file.Close()

	// Decode the image.
	img, _, err := image.Decode(file)
	if err != nil {
		return "", 0, counter, err
	}

	newImg = img
	sizeInMB := weight

	for sizeInMB > maxWeight && counter < 10 {

		// get file width
		width := newImg.Bounds().Max.X
		newWidth := uint(float64(width) * coef)

		// Resize the image to the specified width while maintaining aspect ratio.
		newImg = resize.Resize(newWidth, 0, newImg, resize.Lanczos3)
		//fmt.Printf("====> Type of newImg: %T\n", newImg)

		// Get the size of the resized image in MB
		sizeInMB, err = imageSizeInMB(newImg)
		if err != nil {
			log.Fatal(err)
			break
		}
		counter++
	}

	// fmt.Printf(" => Resized image size: %.2f MB\n", sizeInMB)

	// Create a new file for the resized image.
	pathNew := strings.TrimSuffix(imgPath, filepath.Ext(imgPath)) + "_resized" + filepath.Ext(imgPath)
	outFile, err := os.Create(pathNew)
	if err != nil {
		return "", 0, counter, err
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
		return "", 0, counter, err
	}

	err = common.CopyFileMetadata(imgPath, pathNew)
	if err != nil {
		fmt.Println("    copyFileMetadata failed:", err)
	}
	//fmt.Printf("Resized and saved: %s\n", outFile.Name())
	return pathNew, sizeInMB, counter, nil

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
