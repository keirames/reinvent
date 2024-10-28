package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func getValuesFromCSV(fileName string, columnName string) ([]string, error) {
	file, err := os.Open("./" + fileName + ".csv")
	if err != nil {
		return nil, fmt.Errorf("Error while read file %v", err)
	}

	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()

	// Checks for the error
	if err != nil {
		return nil, fmt.Errorf("Error reading records %v", err)
	}

	header := records[0]
	if len(header) == 0 {
		return nil, fmt.Errorf("Header is empty!")
	}

	idx, err := getIdxOfColumn("Account Name", header)
	if err != nil {
		return nil, fmt.Errorf("Get idx of column failed %v", err)
	}

	names := []string{}

	// Loop to iterate through
	// and save each of the string slice
	for i, r := range records {
		// skip header
		if i == 0 {
			continue
		}

		for valueIdx, v := range r {
			if valueIdx == idx {
				names = append(names, v)
			}
		}
	}

	return names, nil
}

func main() {
	names1, err := getValuesFromCSV("1", "Account Name")
	if err != nil {
		log.Fatal(err)
	}

	names2, err := getValuesFromCSV("2", "Account Name")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("items missing in file 2")
	for _, i := range diffCheck(names1, names2) {
		fmt.Println(i)
	}
	fmt.Println("---")
	fmt.Println("items missing in file 1")
	for _, i := range diffCheck(names2, names1) {
		fmt.Println(i)
	}
}

func diffCheck(arr1 []string, arr2 []string) []string {
	diff := []string{}

	dupMap := map[string]bool{}
	for _, item := range arr1 {
		dupMap[item] = true
	}

	for _, item := range arr2 {
		_, ok := dupMap[item]
		if !ok {
			diff = append(diff, item)
		}
	}

	return diff
}

func getIdxOfColumn(name string, columnNames []string) (int, error) {
	for i, n := range columnNames {
		if n == name {
			return i, nil
		}
	}

	return 0, fmt.Errorf("Cannot find column name")
}
