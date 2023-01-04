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
	db := fts.New[ICD]()
	db.InsertBatchSync(icds)
	fmt.Println(fmt.Sprintf("Time to index %s", time.Since(start)))
	fmt.Println(db.Search("childbirth diabetes controlled"))
	fmt.Println(fmt.Sprintf("Time to search %s", time.Since(start)))
}

type ICD struct {
	Code string `json:"code" index:"true"`
	Desc string `json:"desc" index:"true"`
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

func readData(data string) (icds []ICD) {
	if err := json.Unmarshal([]byte(data), &icds); err != nil {
		fmt.Printf("failed to unmarshal json file, error: %v", err)
		return
	}
	return
}
