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
	CONED        = `Consolidated Edison Company of New York, Inc.`
	ALL          = `ConEd in All Zones`
	ZONEJ        = `ConEd in Zone J`
	DECDELIVERY  = 74.65  // Delivery Charge in December
	DECWATT      = 422    // wattage used in December
	DECRATE      = 6.4408 // rate charged by conEd in dec for supply
	DECBILLTOTAL = 108.86 // total bill (also includes fees & taxes on both suppy and delivery)
)

type Energy struct {
	name        string
	rate        float64
	minTerm     float64
	supplyTotal float64
	total       float64
}

func newEnergy(name string, rate, term float64) Energy {
	e := Energy{name: name}
	e.rate = rate
	e.minTerm = term
	return e
}

func main() {
	file, err := os.Open("data/active_offers.csv")

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

	source := parseData(records)
	source = calculateDecTotal(source)
	for _, v := range source {
		fmt.Println(v.total)
	}
}

func parseData(records [][]string) []*Energy {
	var source []*Energy
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

			// add all structs to slice of
			source = append(source, &e)
		}
	}
	return source
}

func calculateDecTotal(source []*Energy) []*Energy {
	for _, v := range source {
		v.supplyTotal = v.rate * DECWATT
		v.total = v.supplyTotal + DECDELIVERY
		i := fmt.Sprintf("%.2f", v.total)
		v.total, _ = strconv.ParseFloat(i, 64)
		//fmt.Println(v.total)
	}

	return source
}
