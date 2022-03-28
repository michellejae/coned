package models

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type NewBill struct {
	Month    time.Month
	Year     int
	Wattage  int
	Delivery float64
	Rate     float64
	Total    float64
}

func newBill(wattage, year int, delivery, rate, total float64, month time.Month) NewBill {
	b := NewBill{Wattage: wattage}
	b.Month = month
	b.Year = year
	b.Delivery = delivery
	b.Rate = rate
	b.Total = total
	return b
}

var ConEdBills map[string]NewBill

func OpenConEdCSV(condEdCsv string) {
	ConEdBills = make(map[string]NewBill)

	file, err := os.Open(condEdCsv)
	if err != nil {
		log.Fatal("error opening yearly file", err)
	}

	defer file.Close()

	filedata := csv.NewReader(file)

	filedata.FieldsPerRecord = -1

	records, err := filedata.ReadAll()
	if err != nil {
		log.Fatal("coned readall error", err)
	}

	parseConEdData(records)

}

func parseConEdData(records [][]string) {
	// create new bill for every month (starting jan 2021)
	// may need to reset slices of conedbills like i did with source

	// also need to change the monthly models so it reference these bills
	for i, r := range records[20:] {
		endDate := r[2]
		watts := r[3]      // int
		totals := r[5]     // float64
		deliveries := r[6] // float64
		rates := r[7]      // float64

		const shortForm = "2006-01-02"
		date, _ := time.Parse(shortForm, endDate)

		// for february of every year the endDate is in march
		// so subract month from date to change it to feb
		if i == 1 || i == 13 {
			date = date.AddDate(0, -1, 0)

		}

		month := date.Month()
		year := date.Year()

		watt, _ := strconv.Atoi(watts)

		delivery, _ := strconv.ParseFloat(deliveries, 64)

		rate, _ := strconv.ParseFloat(rates, 64)

		totals = strings.TrimPrefix(totals, "$")
		total, _ := strconv.ParseFloat(totals, 64)

		bill := newBill(watt, year, delivery, rate, total, month)

		title := fmt.Sprintf("%v_%v", month, year)

		ConEdBills[title] = bill

	}

}
