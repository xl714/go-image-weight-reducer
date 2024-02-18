package main

import (
	"fmt"

	// Import the function from the same directory
	arguments "github.com/xl714/go-image-weight-reducer/common"
	filesys "github.com/xl714/go-image-weight-reducer/common"
)

func main() {
	args, err := arguments.ParseArguments()
	if err != nil {
		fmt.Println("Error parsing arguments:", err)
		return
	}

	// Use the parsed arguments here
	fmt.Printf("Image max weight: %.2f MB\n", args.MaxWeight)
	fmt.Printf("Limit enabled: %t\n", args.Limit)
	fmt.Printf("Verbose mode: %t\n", args.Verbose)

	imageFiles, err := filesys.ListFiles("images", []string{".png", ".jpg", ".jpeg"}, false)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Image files:")
	for _, file := range imageFiles {
		fmt.Printf("Path: %s, Name: %s, Ext: %s, IsDir: %t\n", file.Path, file.Name, file.Ext, file.IsDir)
	}

}
