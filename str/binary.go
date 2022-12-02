package str

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/klauspost/compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func GenerateBinaryContent(packageName, varName string, data []byte, fileName ...string) []byte {
	var compressed bytes.Buffer
	w := gzip.NewWriter(&compressed)
	_, _ = w.Write(data)
	_ = w.Close()
	encoded := base64.StdEncoding.EncodeToString(compressed.Bytes())
	output := &bytes.Buffer{}
	output.WriteString("package " + packageName + "\n\nvar " + varName + " = " + strconv.Quote(encoded) + "\n")
	bt := output.Bytes()
	if len(fileName) > 0 {
		writeFile(fileName[0], bt)
	}
	return bt
}

func DecodeBinaryString(data string) ([]byte, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	zipReader, err := gzip.NewReader(bytes.NewReader(decodedBytes))
	if err != nil {
		return nil, err
	}

	rawBytes, err := io.ReadAll(zipReader)
	if err != nil {
		return nil, err
	}

	return rawBytes, nil
}

func writeFile(filePath string, data []byte) {
	fmt.Printf("Writing new %s\n", filePath)
	err := ioutil.WriteFile(filePath, data, os.FileMode(0664))
	if err != nil {
		log.Fatalf("Error writing '%s': %s", filePath, err)
	}
}
