package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/willsem/tfs-go-hw/hw02/company/calculate"

	"github.com/willsem/tfs-go-hw/hw02/company"
)

const (
	FILEENV = "FILE"
	OUTFILE = "out.json"
)

var filepath string

func init() {
	flag.StringVar(&filepath, "file", "", "input filepath")
}

func main() {
	flag.Parse()
	var ok bool

	if filepath == "" {
		filepath, ok = os.LookupEnv(FILEENV)
		if !ok {
			fmt.Print("Type path to the filepath: ")
			if _, err := fmt.Scanf("%s", &filepath); err != nil {
				fmt.Println("Error when tried to input")
				return
			}
		}
	}

	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Cannot open the file: %s\n", filepath)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("Cannot read the file: %s\n", filepath)
		return
	}

	var operations []company.Operation
	err = json.Unmarshal(data, &operations)
	if err != nil {
		fmt.Println(err)
	}

	outFile, err := os.Create(OUTFILE)
	if err != nil {
		fmt.Printf("Cannot create the file %s\n", OUTFILE)
		return
	}
	defer outFile.Close()

	enc := json.NewEncoder(outFile)
	enc.SetIndent("", "\t")

	err = enc.Encode(calculate.Operations(operations))
	if err != nil {
		fmt.Printf("Cannot marshal the result: %s\n", err)
		return
	}
}
