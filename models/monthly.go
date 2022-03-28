package models

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

// type Bill struct {
// 	Name     string
// 	Delivery float64
// 	Wattage  int
// 	Rate     float64
// 	Total    float64
// }

// var Dec = Bill{
// 	Name:     "Dec 2021",
// 	Delivery: 74.65,
// 	Wattage:  422,
// 	Rate:     6.4408,
// 	Total:    100.86,
// }

// var Jan = Bill{
// 	Name:     "Jan 2022",
// 	Delivery: 80.01,
// 	Wattage:  415,
// 	Rate:     16.2072,
// 	Total:    157.31,
// }

var Month NewBill
var Min = ""

// remember fields have to be capital to send them to front end
type Energy struct {
	Name         string  `json:"name"`
	Rate         float64 `json:"rate"`
	MinTerm      float64 `json:"minTerm"`
	SupplyTotal  float64 `json:"supplyTotal"`
	Total        float64 `json:"total"`
	OfferType    string  `json:"offerType"`
	Cancellation string  `json:"cancellation"`
	EnergySource string  `json:"energySource"`
	PercentRenew string  `json:"percentRenew"`
}

func newEnergy(name, offerType, energySource, percentRenew, cancellation string, rate, term float64) *Energy {
	e := Energy{Name: name}
	e.Rate = rate
	e.MinTerm = term
	e.OfferType = offerType
	e.EnergySource = energySource
	e.Cancellation = cancellation
	e.PercentRenew = percentRenew
	return &e
}

var Source []*Energy

func OpenFile(csvFile string) {

	file, err := os.Open(csvFile)

	fileName := strings.TrimPrefix(csvFile, "data/")
	fileName = strings.TrimSuffix(fileName, "_offers.csv")

	Month = ConEdBills[fileName]

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
		log.Fatal("readAll error", err)
	}

	parseData(records)

}

func parseData(records [][]string) {
	// reset the source slice so the data per month doesn't keep adding on
	// ie if i clicked dec, then clicked jan, jan graph would show both dec and jan
	Source = Source[:0]

	// loop through each line of csv
	for _, r := range records[1:] { // skip line one as it's header
		// r[0] is utitily (who delivers me energy, has to be coned)
		// r[4] has to be electric (some rates are gas)
		rates := r[7]
		terms := r[9]
		name := r[2]
		offerType := r[5]
		cancellation := r[10]
		energySource := r[17]
		percentRenew := r[16]

		if (r[0] == CONED && r[4] == `ELECTRIC`) && (r[1] == ALL || r[1] == ZONEJ) {
			// trim off kwh off each rate then convert to float64
			rates = strings.TrimSuffix(rates, " kWh")
			rate, _ := strconv.ParseFloat(rates, 64)

			// trim months of contract length, convert to float
			terms = strings.TrimSuffix(terms, " Month(s)")
			term, _ := strconv.ParseFloat(terms, 64)

			//create new struct of each energy source
			e := newEnergy(name, offerType, energySource, percentRenew, cancellation, rate, term)

			// add all structs to slice of
			Source = append(Source, e)

		}

	}
	calculateBillsTotal(Source)

}

func calculateBillsTotal(Source []*Energy) {
	min := Source[0].Rate

	// loop through slice of energy structs (ESCO's
	for _, v := range Source {
		// supplytotal = the rate per esco * my dec watt usage
		v.SupplyTotal = v.Rate * float64(Month.Wattage)
		// my supply total + my dec delivery charge
		v.Total = v.SupplyTotal + Month.Delivery
		// conver to string
		i := fmt.Sprintf("%.2f", v.Total)
		v.Total, _ = strconv.ParseFloat(i, 64)

		// while in for loop calculate the lowest rate so we know for graphing
		// for now using rate, may change to total if start calculating differently
		if v.Rate < min {
			Min = v.Name
			min = v.Rate
		}
	}
}
