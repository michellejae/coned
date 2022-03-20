package controller

import (
	"net/http"

	"github.com/michellejae/coned/views"
)

var home *views.View

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	home = views.NewView("giraffe", "views/home.html")

	home.Render(w, nil)
}
