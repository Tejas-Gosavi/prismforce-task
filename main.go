package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	// reading data from given JSON file
	fileBytes, err := ioutil.ReadFile("2-input.json")
	if err != nil {
		// if any error occurs then print that error and stop executing current program
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// getting parsed data from fileBytes
	var parsedData Data = ParseJSON(fileBytes)

	// getting final balanced sheet
	var bs BalanceSheet = CreateBalanceSheet(parsedData)

	// encodign final balanced sheet in JSON
	outputJSON, err := json.MarshalIndent(bs, "", "  ")
	if err != nil {
		// if any error occurs then print that error and stop executing current program
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// saving encoded final balance sheet in given JSON file
	err = ioutil.WriteFile("2-input-test.json", outputJSON, 0644)
	if err != nil {
		// if any error occurs then print that error and stop executing current program
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
