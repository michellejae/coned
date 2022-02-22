package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
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
	name         string
	rate         float64
	minTerm      float64
	supplyTotal  float64
	total        float64
	offerType    string
	cancellation string
	energySource string
	percentRenew string
}

func newEnergy(name, offerType, energySource, percentRenew, cancellation string, rate, term float64) *Energy {
	e := Energy{name: name}
	e.rate = rate
	e.minTerm = term
	e.offerType = offerType
	e.energySource = energySource
	e.cancellation = cancellation
	e.percentRenew = percentRenew
	return &e
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
	calculateDecTotal(source)
	graphData(source)
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
			source = append(source, e)
		}
	}
	return source
}

func calculateDecTotal(source []*Energy) {
	// loop through slice of energy structs (ESCO's)
	for _, v := range source {
		// supplytotal = the rate per esco * my dec watt usage
		v.supplyTotal = v.rate * DECWATT
		// my supply total + my dec delivery charge
		v.total = v.supplyTotal + DECDELIVERY
		// conver to string
		i := fmt.Sprintf("%.2f", v.total)
		v.total, _ = strconv.ParseFloat(i, 64)
		//fmt.Println(v.total)
	}

}

func graphData(source []*Energy) {
	bar := charts.NewBar()

	bar.AddSeries("Totals", generateData(source))

	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "My Energy Bills per ESCO",
		Subtitle: "ConEd Delivery Rate + (ESCO rate * kw usage)",
	}),
		charts.WithXAxisOpts(opts.XAxis{
			Type: "category",
			Show: false,
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1200px",
			Height: "600px",
		}))
	f, _ := os.Create("bar.html")

	bar.Render(f)

}

// minTerm      float64
// supplyTotal  float64
// total        float64
// offerType    string
// cancellation string
// energySource string
// percentRenew string

func generateData(source []*Energy) []opts.BarData {

	items := make([]opts.BarData, 0)
	// loop through source
	for _, v := range source {
		name := v.name
		term := v.minTerm

		cancellation := v.cancellation
		energy := v.energySource
		renewable := v.percentRenew
		//append each ESCO to the opts.BarData slice
		items = append(items, opts.BarData{Name: fmt.Sprintf("fudkity fuck fuck %s, %b, %s, %s, %s, %s", name, term, energy, renewable, cancellation, v.offerType), Value: v.total})
	}

	return items
}

// func singleOut(source []*Energy) {
// 	count := 0
// 	boo := 0
// 	for i, v := range source {
// 		boo++
// 		if i > 1 {
// 			last := source[i-1]
// 			if v.name == last.name {
// 				count++
// 			}
// 		}

// 	}
// 	fmt.Println(count, boo)
// }
