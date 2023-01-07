package main

import (
	"encoding/json"
	"fmt"
	"github.com/sujit-baniya/pkg/fts"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	start := time.Now()
	icds := readFile()
	db := fts.New[ICD]("icd")
	db.InsertBatch(icds)
	fmt.Printf("Time to index %s", time.Since(start))
	fmt.Println(db.SearchExact("third trimester pregnancy diabetes"))
	fmt.Printf("Time to search %s", time.Since(start))
}

type ICD struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
}

func readFile() (icds []ICD) {
	file, err := os.Open("icd10_codes.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	jsonData, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("failed to read json file, error: %v", err)
		return
	}

	if err := json.Unmarshal(jsonData, &icds); err != nil {
		fmt.Printf("failed to unmarshal json file, error: %v", err)
		return
	}
	return
}

func readFileAsMap() (icds []map[string]any) {
	file, err := os.Open("icd10_codes.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	jsonData, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("failed to read json file, error: %v", err)
		return
	}

	if err := json.Unmarshal(jsonData, &icds); err != nil {
		fmt.Printf("failed to unmarshal json file, error: %v", err)
		return
	}
	return
}

func readFromString() []string {
	return []string{
		"Salmonella pneumonia",
		"Diabetes uncontrolled",
	}
}

func readFromInt() []int {
	return []int{
		10,
		100,
		20,
	}
}
