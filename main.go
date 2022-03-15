package main

import (
	"fmt"
	"net/http"

	"github.com/michellejae/coned/controller"
	"github.com/michellejae/coned/models"
)

const (
	dec = "data/active_offers.csv"
)

func main() {

	models.OpenFile(dec)

	http.HandleFunc("/", controller.GenerateAndGraph)
	fmt.Println("Starting the server on :3001")
	http.ListenAndServe(":3001", nil)

}
