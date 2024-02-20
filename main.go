package main

import (
	"fmt"
	"sync"

	// Import the function from the same directory
	common "github.com/xl714/go-image-weight-reducer/common"
	//filesys "github.com/xl714/go-image-weight-reducer/common"
	imagehelper "github.com/xl714/go-image-weight-reducer/imagehelper"
)

type FileInfo struct {
	Path   string
	Name   string
	Ext    string
	IsDir  bool
	Weight float64
}

type ProcessInfo struct {
	Path string
	SizeOriginal int64
	IsResized bool
	SizeNew int64
	ErrorMessage string
	DateUpdatedOriginal string
	DateUpdatedNew string
	DateCreatedOriginal string
	DateCreatedNew string
}

func main() {
	args, err := common.ParseArguments()
	if err != nil {
		fmt.Println("Error parsing arguments:", err)
		return
	}

	// Use the parsed arguments here
	fmt.Printf("Image max weight: %.2f MB\n", args.MaxWeight)
	fmt.Printf("Limit enabled: %d\n", args.Limit)
	fmt.Printf("Verbose mode: %t\n", args.Verbose)

	//imageFiles, err := filesys.ListFiles("images", []string{".png", ".jpg", ".jpeg"}, false)
	imageFiles, err := common.ListFiles("images", []string{".jpg", ".jpeg"}, false)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	
	var wg sync.WaitGroup
	processInfoChan := make(chan FileInfo)

	fmt.Println("Image files:")

	for i, file := range imageFiles {
		if args.Limit > 0 && i >= args.Limit {
			break
		}
		fmt.Printf("Path: %s, Name: %s, Ext: %s, IsDir: %t, Weight: %.2f\n", file.Path, file.Name, file.Ext, file.IsDir, file.Weight)

		if file.Weight > args.MaxWeight {
			fmt.Println("   => Add to process")
			wg.Add(1)
			go processFile(file, args.MaxWeight, args.Verbose, &wg, processInfoChan)	
			// if err != nil {
			// 	fmt.Println("    Error parsing arguments:", err)
			// } else {
			// 	fmt.Printf("   Image weight new path: %s\n", newPath)
			// }
		} else {
			fmt.Println("   Image weight is less than max weight")
		}
	}
	go func() {
		wg.Wait()
		close(processInfoChan)
	}()
	fmt.Println("All done infos")
	for processInfo := range processInfoChan {
		// fmt.Printf("File: %s, Size: %d bytes\n", fileInfo.Path, fileInfo.Size)
		fmt.Printf("%+v\n", processInfo)
	}
}

func processFile(file:FileInfo, maxWeight, verbose, wg *sync.WaitGroup, ch chan<- ProcessInfo) {
	defer wg.Done()

	newPath, err := imagehelper.ResizeImage(file.Path, file.Ext, file.Weight, maxWeight, verbose)

	processInfo := ProcessInfo{
        Path:                newPath,
        SizeOriginal:        file.Weight,
        IsResized:           true,
        SizeNew:             0,
        ErrorMessage:        err,
        DateUpdatedOriginal: "2022-01-01",
        DateUpdatedNew:      "2024-02-20",
        DateCreatedOriginal: "2022-01-01",
        DateCreatedNew:      "2024-02-20",
    }

	ch <- processInfo
}
