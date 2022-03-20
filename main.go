package main

import (
	"fmt"
	"net/http"

	"github.com/michellejae/coned/controller"
)

func main() {

	http.HandleFunc("/", controller.HomeHandler)
	http.HandleFunc("/graphs", controller.GenerateAndGraph)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	fmt.Println("Starting the server on :3001")
	http.ListenAndServe(":3001", nil)

}
