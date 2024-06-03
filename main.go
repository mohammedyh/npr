package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Script struct {
	key   string
	value string
}

func main() {
	packageJsonContent, readFileErr := os.ReadFile("package.json")

	if readFileErr != nil {
		fmt.Println("package.json not found")
		os.Exit(1)
	}

	var parsedJson map[string]interface{}

	parseErr := json.Unmarshal(packageJsonContent, &parsedJson)
	if parseErr != nil {
		fmt.Println(parseErr.Error())
		os.Exit(1)
	}

	result, err := json.MarshalIndent(parsedJson["scripts"], "", " ")

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(string(result))
}
