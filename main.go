package main

import (
	"fmt"
	"log"
	"os"
)

func main() {

	arguments := os.Args[1:]
	var fileName string = arguments[0] //"sample_data.csv"
	log.Print(fileName)
	parseCSV(fileName)

	letters := [...]string{"a", "b", "c", "d", "", ""}
	aSlice := letters[4:]
	fmt.Printf("letters: %s\n", aSlice)
	fmt.Printf("Only Empties: %s\n", sliceRangeContainsOnlyEmpties(aSlice))
	fmt.Printf("Has Non Empty: %s\n", sliceRangeContainsNonEmptyValue(aSlice))
	fmt.Printf("Number of Empties: %s\n", countEmptyValuesIn(aSlice))
	fmt.Printf("Number of Non Empty values: %s\n", countNonEmptyValuesIn(aSlice))

	aSlice = letters[0:4]
	fmt.Printf("letters: %s\n", aSlice)
	fmt.Printf("Only Empties: %s\n", sliceRangeContainsOnlyEmpties(aSlice))
	fmt.Printf("Has Non Empty: %s\n", sliceRangeContainsNonEmptyValue(aSlice))
	fmt.Printf("Number of Empties: %s\n", countEmptyValuesIn(aSlice))
	fmt.Printf("Number of Non Empty values: %s\n", countNonEmptyValuesIn(aSlice))

}
