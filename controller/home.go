package controller

import (
	"fmt"
	"html/template"
	"net/http"
)

type HomePage struct {
	Script template.HTML
}

func HomeView(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("view/home.html")
	if err != nil {
		fmt.Println("home view template load error", err)
	}

	tmplData := HomePage{
		Script: template.HTML(graphScript),
	}

	err = t.Execute(w, tmplData)
	if err != nil {
		fmt.Println("home view execute error", err)
	}
}

var graphScript = `<script type="text/javascript">
"use strict";
let goecharts_{{ .ChartID | safeJS }} = echarts.init(document.getElementById('{{ .ChartID | safeJS }}'), "{{ .Theme }}");
let option_{{ .ChartID | safeJS }} = {{ .JSON }};
goecharts_{{ .ChartID | safeJS }}.setOption(option_{{ .ChartID | safeJS }});
{{- range .JSFunctions.Fns }}
{{ . | safeJS }}
{{- end }}
</script>`
