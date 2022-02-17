package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	CONED = `Consolidated Edison Company of New York, Inc.`
	ALL   = `ConEd in All Zones`
	ZONEJ = `ConEd in Zone J`
)

type Energy struct {
	name    string
	rate    float64
	minTerm float64
}

func newEnergy(name string, rate, term float64) *Energy {
	e := Energy{name: name}
	e.rate = rate
	e.minTerm = term
	return &e
}

// should i make this a pointer?
var source []Energy

func main() {
	file, err := os.Open("data/active_offers.csv")

	// var names []string
	// var rates []float64
	// var minTerms []float64

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	filedata := csv.NewReader(file)
	// top row for some reason is a field smaller so this tell it's to piss off
	// normally each row has to have same amounto fields (columns)
	filedata.FieldsPerRecord = -1

	records, err := filedata.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	source = parseData(records)
	fmt.Println(source)

}

func parseData(records [][]string) []Energy {
	// escoNames := []string{}
	// rates := []float64{}
	// minimumTerm := []float64{}

	// loop through each line of csv
	for _, r := range records[1:] { // skip line one as it's header
		// r[0] is utitily (who delivers me energy, has to be coned)
		// r[4] has to be electric (some rates are gas)
		rates := r[7]
		terms := r[9]
		name := r[2]

		if (r[0] == CONED && r[4] == `ELECTRIC`) && (r[1] == ALL || r[1] == ZONEJ) {
			// trim off kwh off each rate then convert to float64
			rates = strings.TrimSuffix(rates, " kWh")
			rate, _ := strconv.ParseFloat(rates, 64)

			// trim months of contract length, convert to float
			terms = strings.TrimSuffix(terms, " Month(s)")
			term, _ := strconv.ParseFloat(terms, 64)

			//create new struct of each energy source
			e := newEnergy(name, rate, term)

			// add all structs to slice of structs
			source = append(source, *e)
		}

	}
	return source

}

// will need
// var s = ".0899 kWh"
// 	s = strings.TrimSuffix(s, " kWh")
// 	fmt.Print(s)
// Month(s)
