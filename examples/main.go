package main

import (
	"encoding/json"
	"fmt"
	"github.com/sujit-baniya/pkg/fts"
	"io"
	"log"
	"os"
)

var data = `[{"code": "A000", "desc": "Cholera due to Vibrio cholerae 01, biovar cholerae"}, {"code": "A001", "desc": "Cholera due to Vibrio cholerae 01, biovar eltor"}, {"code": "A009", "desc": "Cholera, unspecified"}, {"code": "A0100", "desc": "Typhoid fever, unspecified"}, {"code": "A0101", "desc": "Typhoid meningitis"}, {"code": "A0102", "desc": "Typhoid fever with heart involvement"}, {"code": "O1200", "desc": "Gestational edema, unspecified trimester"}, {"code": "O1201", "desc": "Gestational edema, first trimester"}, {"code": "O1202", "desc": "Gestational edema, second trimester"}, {"code": "O1203", "desc": "Gestational edema, third trimester"}, {"code": "O1204", "desc": "Gestational edema, complicating childbirth"}, {"code": "O1205", "desc": "Gestational edema, complicating the puerperium"}, {"code": "O1210", "desc": "Gestational proteinuria, unspecified trimester"}, {"code": "O24410", "desc": "Gestational diabetes mellitus in pregnancy, diet controlled"}, {"code": "O24414", "desc": "Gestational diabetes mellitus in pregnancy, insulin controlled"}, {"code": "O24415", "desc": "Gestational diabetes mellitus in pregnancy, controlled by oral hypoglycemic drugs"}, {"code": "O24419", "desc": "Gestational diabetes mellitus in pregnancy, unspecified control"}, {"code": "O24420", "desc": "Gestational diabetes mellitus in childbirth, diet controlled"}, {"code": "O24424", "desc": "Gestational diabetes mellitus in childbirth, insulin controlled"}]`

var data1 = `[{"code": "A000", "desc": "Cholera due to Vibrio cholerae 01, biovar cholerae"}, {"code": "A001", "desc": "Cholera due to Vibrio cholerae 01, biovar eltor"}]`

func main() {
	// icds := readData(data)
	icds := readFile()
	db := fts.New[ICD]()
	// db.InsertBatchSync(icds)
	db.InsertBatchAsync(icds)
	fmt.Println(db.Search("childbirth diabetes controlled"))
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
