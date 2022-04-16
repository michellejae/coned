package main

import (
	"fmt"
	"net/http"

	"github.com/michellejae/coned/controller"
	"github.com/michellejae/coned/models"
)

const (
	coned = "data/coned_bill.csv"
)

func main() {

	models.OpenConEdCSV(coned)

	http.HandleFunc("/", controller.HomeHandler)
	http.HandleFunc("/graphs/monthly", controller.GenerateAndGraph)
	http.HandleFunc("/graphs/yearly", controller.YearlyGraph)

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	fmt.Println("Starting the server on :3001")
	http.ListenAndServe(":3001", nil)

}
