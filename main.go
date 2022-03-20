package main

import (
	"fmt"
	"net/http"

	"github.com/michellejae/coned/controller"
)

// const (
// /	dec = "data/dec_offers.csv"

// 	//jan = "data/jan_offers.csv"
// )

func main() {

	// File := make(map[string]string)
	// File["dec"] = "data/dec_offers.csv"

	// models.OpenFile(dec)

	//	models.OpenFile(jan)

	http.HandleFunc("/", controller.HomeHandler)
	http.HandleFunc("/graphs", controller.GenerateAndGraph)

	//http.HandleFunc("/graphs", controller.GenerateAndGraph)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	fmt.Println("Starting the server on :3001")
	http.ListenAndServe(":3001", nil)

}
