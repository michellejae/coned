package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("data/example.csv")

	if err != nil {
		log.Fatal(err)
	}

	filedata := csv.NewReader(file)

	for {
		record, err := filedata.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		for value := range record {
			fmt.Println(record[value])
			fmt.Println("hi")
			//fmt.Printf("%s\n", record[value])
		}
	}
}
