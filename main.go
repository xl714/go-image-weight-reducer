package main

import (
	"fmt"
	"sync"

	"github.com/davecgh/go-spew/spew"

	// Import the function from the same directory
	common "github.com/xl714/go-image-weight-reducer/common"
	//filesys "github.com/xl714/go-image-weight-reducer/common"
	imagehelper "github.com/xl714/go-image-weight-reducer/imagehelper"
)

type Arguments struct {
	MaxWeight float64
	Limit     int
	Verbose   bool
}

type FileInfo struct {
	Path   string
	Name   string
	Ext    string
	IsDir  bool
	Weight float64
}

type ProcessInfo struct {
	Path           string
	PathNew        string
	WeightOriginal float64
	WeightNew      float64
	ResizedNumber  int64
	// ErrorMessage        string
	// DateUpdatedOriginal string
	// DateUpdatedNew      string
	// DateCreatedOriginal string
	// DateCreatedNew      string
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

	imageFiles, err := common.ListFiles("images", []string{".png", ".jpg", ".jpeg"}, false)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var wg sync.WaitGroup
	processInfoChan := make(chan ProcessInfo, 2)

	fmt.Println("Image files:")

	for i, file := range imageFiles {

		fmt.Printf("Path: %s, Name: %s, Ext: %s, IsDir: %t, Weight: %.2f\n", file.Path, file.Name, file.Ext, file.IsDir, file.Weight)

		if file.Weight > args.MaxWeight {
			if args.Limit > 0 && i >= args.Limit {
				break
			}
			fmt.Println("   => Add to process")
			wg.Add(1)
			go func(file common.FileInfo, args common.Arguments) {
				defer wg.Done()
				processFile(file, args.MaxWeight, args.Verbose, processInfoChan)
			}(file, args)
			// go processFile(file, args.MaxWeight, args.Verbose, &wg, processInfoChan)
		} else {
			fmt.Println("   Image weight is less than max weight")
		}
	}

	go func() {
		wg.Wait()
		close(processInfoChan)
	}()

	// Wait for all workers to finish before ranging over the channel
	wg.Wait()

	fmt.Print("\nAll done. Processed files infos:\n")
	for processInfo := range processInfoChan {
		// fmt.Printf("File: %s, Size: %d bytes\n", fileInfo.Path, fileInfo.Size)
		// fmt.Printf("%+v\n", processInfo) // print any struct
		fmt.Printf("\n")
		spew.Dump(processInfo)
	}
}

func processFile(file common.FileInfo, maxWeight float64, verbose bool, processInfo chan<- ProcessInfo) {
	// defer wg.Done()

	pathNew, weightNew, resizedNumber, err := imagehelper.ResizeImage(file.Path, file.Ext, file.Weight, maxWeight, verbose)
	if err != nil {
		fmt.Println("    Error  imagehelper.ResizeImage:", err)
	}
	// else {
	// 	fmt.Printf("   Image reduced new path: %s\n", pathNew)
	// }

	processInfo <- ProcessInfo{
		Path:           file.Path,
		PathNew:        pathNew,
		WeightOriginal: file.Weight,
		WeightNew:      weightNew,
		ResizedNumber:  resizedNumber,
		// DateUpdatedOriginal: "2022-01-01",
		// DateUpdatedNew:      "2024-02-20",
		// DateCreatedOriginal: "2022-01-01",
		// DateCreatedNew:      "2024-02-20",
	}
}
