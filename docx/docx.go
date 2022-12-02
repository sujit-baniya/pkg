package docx

import (
	"bytes"
	"path/filepath"
	"strings"
)

func PrepareDocxToFile(file string, data map[string]interface{}, outputFile ...string) error {
	var output string
	if len(outputFile) == 0 {
		output = strings.Replace(filepath.Base(file), filepath.Ext(file), "-Filled", -1) + filepath.Ext(file)
	} else {
		output = outputFile[0]
	}
	doc, err := Open(file)
	if err != nil {
		return err
	}
	err = doc.ReplaceAll(data)
	if err != nil {
		return err
	}

	return doc.WriteToFile(output)
}

func PrepareDocx(file string, data map[string]interface{}) (*bytes.Buffer, error) {
	var byteBuffer bytes.Buffer
	doc, err := Open(file)
	if err != nil {
		return nil, err
	}
	err = doc.ReplaceAll(data)
	if err != nil {
		return nil, err
	}
	err = doc.Write(&byteBuffer)
	if err != nil {
		return nil, err
	}
	return &byteBuffer, nil
}
