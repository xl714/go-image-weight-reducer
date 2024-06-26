package common

import (
	"flag"
	"fmt"
	// "strconv"
)

// ParseArguments processes command-line arguments and returns them as a struct
func ParseArguments() (Arguments, error) {
	var maxWeight float64
	var limit int
	var verbose bool
	var selectFolder bool

	// Define flags with default values and usage descriptions
	flag.Float64Var(&maxWeight, "image-max-weight", 1, "Maximum allowed weight for images (in MB)")
	flag.IntVar(&limit, "limit", 0, "Enables limiting functionality (not implemented yet)")
	flag.BoolVar(&verbose, "verbose", false, "Print additional information during execution")
	flag.BoolVar(&selectFolder, "select-folder", false, "Open select folder dialog to avoid current dir default")

	flag.Parse() // Parse command line arguments

	// Validate and format arguments
	// maxWeight, err := strconv.ParseFloat(maxWeight, 64)
	// if err != nil {
	// 	return Arguments{}, fmt.Errorf("Invalid Image max weight: %v", err)
	// }

	if maxWeight <= 0 {
		//return Arguments{}, fmt.Errorf("Image max weight must be a positive number")
		return Arguments{}, fmt.Errorf("invalid Image max weight: %v", maxWeight)
	}

	// Create a struct to hold the parsed arguments
	args := Arguments{
		MaxWeight:    maxWeight,
		Limit:        limit,
		Verbose:      verbose,
		SelectFolder: selectFolder,
	}

	return args, nil
}

// Arguments struct holds the parsed command-line arguments
type Arguments struct {
	MaxWeight float64
	Limit     int
	Verbose   bool
	SelectFolder bool
}
