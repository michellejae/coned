package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
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
	// if err != nil {
	// 	log.Fatal("parseData error", err)
	// }

	calculateDecTotal(source)

	//graphData(source)
	r := chi.NewRouter()

	r.Get("/", serveHome)

	r.Post("/", sendData)

	fs := http.FileServer(http.Dir("js"))
	r.Handle("/js/*", http.StripPrefix("/js/", fs))

	http.ListenAndServe(":3334", r)

}

// send data to js fetch call
func sendData(writer http.ResponseWriter, r *http.Request) {

	writer.Header().Set("Content-Type", "application/json")

	resultsJSON, err := json.Marshal(source)
	if err != nil {
		log.Fatal("send data", err)
	}

	writer.Write(resultsJSON)

}

// serve up the homepage index file
func serveHome(w http.ResponseWriter, r *http.Request) {
	var homepage HomePage
	tmpl, err := template.ParseFiles("html/home.html")
	if err != nil {
		log.Fatal("serve home error", err)
	}
	tmpl.Execute(w, homepage)

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
		//fmt.Println(v.total)
	}

}
