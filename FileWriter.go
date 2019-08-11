package main

import (
	"fmt"
	"os"
)

// This function is provided for debug purposes so is not properly tested
func WriteSliceToFile(slice []string, fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("unable to write to %v with error: %v", fileName, err)
		return
	}
	defer f.Close()
	for _, sliceString := range slice {
		_, _ = f.WriteString(fmt.Sprintf("%v\n", sliceString))
	}
}
