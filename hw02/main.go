package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/willsem/tfs-go-hw/hw02/company"
)

const (
	FILEENV = "FILE"
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
		fmt.Printf("Cannot open file: %s\n", filepath)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("Cannot read file: %s\n", filepath)
		return
	}

	var operations []company.Operation
	err = json.Unmarshal(data, &operations)
	if err != nil {
		fmt.Println(err)
	}
}
