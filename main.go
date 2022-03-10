package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
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

type HomePage struct {
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

var source []*Energy

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
		log.Fatal("readAll error", err)
	}

	source := parseData(records)
	if err != nil {
		log.Fatal("parseData error", err)
	}

	calculateDecTotal(source)

	http.HandleFunc("/", generateAndGraph)
	http.ListenAndServe(":3334", nil)

}

func parseData(records [][]string) []*Energy {

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
	// loop through slice of energy structs (ESCO's
	for _, v := range source {
		// supplytotal = the rate per esco * my dec watt usage
		v.SupplyTotal = v.Rate * DECWATT
		// my supply total + my dec delivery charge
		v.Total = v.SupplyTotal + DECDELIVERY
		// conver to string
		i := fmt.Sprintf("%.2f", v.Total)
		v.Total, _ = strconv.ParseFloat(i, 64)

	}

}

func generateAndGraph(w http.ResponseWriter, r *http.Request) {

	options := opts.BarData{}

	data := make([]opts.BarData, 0)

	for _, val := range source {

		percentFloat, _ := strconv.ParseFloat(val.PercentRenew, 64)
		percentFloat = percentFloat * 100

		// have to declare these inside the range so they each update for every ESCO
		// i'm not sure why, something with pointer?
		toolTip := opts.Tooltip{}
		itemStyle := opts.ItemStyle{}

		options.Name = val.Name
		options.Value = val.Total

		toolTip.Show = true
		toolTip.Formatter = fmt.Sprintf("Name: %v<br />Total: $%v<br />Rate: %.2f ¢/kWh<br />Offer Type: %v<br />Minimum Contract Length: %v months<br />Energy Source: %v<br />Percent Renewable: %v%%",
			val.Name, val.Total, val.Rate, val.OfferType, val.MinTerm, val.EnergySource, percentFloat)

		options.Tooltip = &toolTip

		if val.Name == "Consolidated Edison Company of New York, Inc." {
			itemStyle.Color = "red"
			options.ItemStyle = &itemStyle

		} else {
			itemStyle.Color = "green"
			options.ItemStyle = &itemStyle

		}

		data = append(data, options)

	}

	bar := charts.NewBar()

	bar.AddSeries("Totals", data)

	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "My Dec 2021 Energy Bills per ESCO",
		Subtitle: "ConEd Delivery Rate + (ESCO rate * kw usage)",
	}),
		charts.WithXAxisOpts(opts.XAxis{
			Type: "category",
			Show: false,
			Name: "ESCO's",
		}),
		//charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1200px",
			Height: "600px",
		}))

	bar.Render(w)
}

// func generateAndGraph(source []*Energy) {

// 	options := opts.BarData{}

// 	data := make([]opts.BarData, 0)

// 	for _, val := range source {

// 		percentFloat, _ := strconv.ParseFloat(val.PercentRenew, 64)
// 		percentFloat = percentFloat * 100

// 		// have to declare these inside the range so they each update for every ESCO
// 		// i'm not sure why, something with pointer?
// 		toolTip := opts.Tooltip{}
// 		itemStyle := opts.ItemStyle{}

// 		options.Name = val.Name
// 		options.Value = val.Total

// 		toolTip.Show = true
// 		toolTip.Formatter = fmt.Sprintf("Name: %v<br />Total: $%v<br />Rate: %.2f ¢/kWh<br />Offer Type: %v<br />Minimum Contract Length: %v months<br />Energy Source: %v<br />Percent Renewable: %v%%",
// 			val.Name, val.Total, val.Rate, val.OfferType, val.MinTerm, val.EnergySource, percentFloat)

// 		options.Tooltip = &toolTip

// 		if val.Name == "Consolidated Edison Company of New York, Inc." {
// 			itemStyle.Color = "red"
// 			options.ItemStyle = &itemStyle

// 		} else {
// 			itemStyle.Color = "green"
// 			options.ItemStyle = &itemStyle

// 		}

// 		data = append(data, options)

// 	}

// 	bar := charts.NewBar()

// 	bar.AddSeries("Totals", data)

// 	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
// 		Title:    "My Dec 2021 Energy Bills per ESCO",
// 		Subtitle: "ConEd Delivery Rate + (ESCO rate * kw usage)",
// 	}),
// 		charts.WithXAxisOpts(opts.XAxis{
// 			Type: "category",
// 			Show: false,
// 			Name: "ESCO's",
// 		}),
// 		//charts.WithTooltipOpts(opts.Tooltip{Show: true}),
// 		charts.WithInitializationOpts(opts.Initialization{
// 			Width:  "1200px",
// 			Height: "600px",
// 		}))
// 	f, _ := os.Create("index.html")

// 	bar.Render(f)

// }
