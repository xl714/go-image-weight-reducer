package main

import (
	"fmt"

	// Import the function from the same directory
	arguments "github.com/xl714/go-image-weight-reducer/common"
	filesys "github.com/xl714/go-image-weight-reducer/common"
	imagehelper "github.com/xl714/go-image-weight-reducer/imagehelper"
)

func main() {
	args, err := arguments.ParseArguments()
	if err != nil {
		fmt.Println("Error parsing arguments:", err)
		return
	}

	// Use the parsed arguments here
	fmt.Printf("Image max weight: %.2f MB\n", args.MaxWeight)
	fmt.Printf("Limit enabled: %d\n", args.Limit)
	fmt.Printf("Verbose mode: %t\n", args.Verbose)

	//imageFiles, err := filesys.ListFiles("images", []string{".png", ".jpg", ".jpeg"}, false)
	imageFiles, err := filesys.ListFiles("images", []string{".jpg", ".jpeg"}, false)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Image files:")
	for i, file := range imageFiles {
		if args.Limit > 0 && i >= args.Limit {
			break
		}
		fmt.Printf("Path: %s, Name: %s, Ext: %s, IsDir: %t, Weight: %.2f\n", file.Path, file.Name, file.Ext, file.IsDir, file.Weight)

		if file.Weight > args.MaxWeight {
			newPath, err := imagehelper.ResizeImage(file.Path, file.Ext, file.Weight, args.MaxWeight, args.Verbose)
			if err != nil {
				fmt.Println("    Error parsing arguments:", err)
			} else {
				fmt.Printf("   Image weight new path: %s\n", newPath)
			}
		} else {
			fmt.Println("   Image weight is less than max weight")
		}
	}
}
