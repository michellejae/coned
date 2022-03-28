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
	NewMonth    time.Month
	NewWattage  int
	NewDelivery float64
	NewRate     float64
	NewTotal    float64
}

func newBill(wattage, monthIncrease int, delivery, rate, total float64) NewBill {
	start := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	start = start.AddDate(0, monthIncrease, 0)
	month := start.Month()
	b := NewBill{NewMonth: month}
	b.NewWattage = wattage
	b.NewDelivery = delivery
	b.NewRate = rate
	b.NewTotal = total
	return b
}

var ConEdBills []NewBill

func OpenConEdCSV(condEdCsv string) {
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
		watts := r[3]      // int
		totals := r[5]     // float64
		deliveries := r[6] // float64
		rates := r[7]      // float64

		watt, _ := strconv.Atoi(watts)

		delivery, _ := strconv.ParseFloat(deliveries, 64)

		rate, _ := strconv.ParseFloat(rates, 64)

		totals = strings.TrimPrefix(totals, "$")
		total, _ := strconv.ParseFloat(totals, 64)

		b := newBill(watt, i, delivery, rate, total)

		ConEdBills = append(ConEdBills, b)
	}

	fmt.Println(ConEdBills)

}
