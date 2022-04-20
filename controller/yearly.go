package controller

import (
	"fmt"
	"net/http"

	"github.com/michellejae/coned/models"
)

func YearlyGraph(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(models.Month)
	for _, val := range models.Source {
		fmt.Println(val.Total)
	}
}
