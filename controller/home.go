package controller

import (
	"fmt"
	"html/template"
	"net/http"
)

func HomeView(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("view/home.html")
	if err != nil {
		fmt.Println("home view template load error", err)
	}
	data := struct {
		Data string
	}{"oops"}

	err = t.Execute(w, data)
	if err != nil {
		fmt.Println("home view execute error", err)
	}
}
