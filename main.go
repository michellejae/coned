package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

const (
	CONED = `Consolidated Edison Company of New York, Inc.`
	ALL   = `ConEd in All Zones`
	ZONEJ = `ConEd in Zone J`
)

func main() {
	file, err := os.Open("data/active_offers.csv")

	//var names []string
	//	var rate []float32

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	filedata := csv.NewReader(file)
	filedata.FieldsPerRecord = -1

	records, err := filedata.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	parseData(records)

}

func parseData(records [][]string) ([]string, []float64) {
	// escoNames := []string{}
	// rates := []float64{}

	// loop through each line of csv
	for _, r := range records[1:] { // skip line one as it's header
		// r[0] is utitily (who delivers me energy, has to be coned)
		// r[4] has to be electric (some rates are gas)

		if (r[0] == CONED && r[4] == `ELECTRIC`) && (r[1] == ALL || r[1] == ZONEJ) {
			fmt.Println(r[0], r[1])
			fmt.Println("FUKCKKCXKKCKCKCK")
		}

	}

	return nil, nil
}
