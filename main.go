package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/michellejae/coned/models"

	chartrender "github.com/go-echarts/go-echarts/v2/render"
)

const (
	dec = "data/active_offers.csv"
)

type HomePage struct {
	HTML template.HTML
}

func main() {

	models.OpenFile(dec)

	http.HandleFunc("/", generateAndGraph)
	fmt.Println("Starting the server on :3001")
	http.ListenAndServe(":3001", nil)

}

func generateAndGraph(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("html/giraffe.html")
	if err != nil {
		panic(err)
	}

	options := opts.BarData{}

	data := make([]opts.BarData, 0)

	for _, val := range models.Source {

		percentFloat, _ := strconv.ParseFloat(val.PercentRenew, 64)
		percentFloat = percentFloat * 100

		// have to declare these inside the range so they each update for every ESCO
		// i'm not sure why, something with pointer?
		toolTip := opts.Tooltip{}
		itemStyle := opts.ItemStyle{}

		options.Name = val.Name
		options.Value = val.Total

		toolTip.Show = true
		toolTip.Formatter = fmt.Sprintf("Name: %v<br />Total: $%v<br />Rate: %.2f ¢/kWh<br />Offer Type: %v<br />Minimum Contract Length: %v months<br />Energy Source: %v<br />Percent Renewable: %v%%",
			val.Name, val.Total, val.Rate, val.OfferType, val.MinTerm, val.EnergySource, percentFloat)

		options.Tooltip = &toolTip

		if val.Name == "Consolidated Edison Company of New York, Inc." {
			itemStyle.Color = "red"
			options.ItemStyle = &itemStyle

		} else {
			itemStyle.Color = "green"
			options.ItemStyle = &itemStyle

		}

		data = append(data, options)

	}

	bar := charts.NewBar()

	bar.AddSeries("Totals", data)

	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "My Dec 2021 Energy Bills per ESCO",
		Subtitle: "ConEd Delivery Rate + (ESCO rate * kw usage)",
	}),
		charts.WithXAxisOpts(opts.XAxis{
			Type: "category",
			Show: false,
			Name: "ESCO's",
		}),
		//charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1200px",
			Height: "600px",
		}))

	bar.Renderer = newSnippetRenderer(bar, bar.Validate)
	var htmlSnippet template.HTML = renderToHtml(bar)

	tmplData := HomePage{
		HTML: htmlSnippet,
	}

	err = t.Execute(w, tmplData)
	if err != nil {
		panic(err)
	}
}

type Renderer interface {
	Render(w io.Writer) error
}

var baseTpl = `
<div class="container">
    <div class="item" id="{{ .ChartID }}" style="width:{{ .Initialization.Width }};height:{{ .Initialization.Height }};"></div>
</div>
{{- range .JSAssets.Values }}
   <script src="{{ . }}"></script>
{{- end }}
<script type="text/javascript">
    "use strict";
    let goecharts_{{ .ChartID | safeJS }} = echarts.init(document.getElementById('{{ .ChartID | safeJS }}'), "{{ .Theme }}");
    let option_{{ .ChartID | safeJS }} = {{ .JSON }};
    goecharts_{{ .ChartID | safeJS }}.setOption(option_{{ .ChartID | safeJS }});
    {{- range .JSFunctions.Fns }}
    {{ . | safeJS }}
    {{- end }}
</script>
`

type snippetRenderer struct {
	c      interface{}
	before []func()
}

func newSnippetRenderer(c interface{}, before ...func()) chartrender.Renderer {
	return &snippetRenderer{c: c, before: before}
}

func (r *snippetRenderer) Render(w io.Writer) error {
	const tplName = "chart"
	for _, fn := range r.before {
		fn()
	}

	tpl := template.
		Must(template.New(tplName).
			Funcs(template.FuncMap{
				"safeJS": func(s interface{}) template.JS {
					return template.JS(fmt.Sprint(s))
				},
			}).
			Parse(baseTpl),
		)

	err := tpl.ExecuteTemplate(w, tplName, r.c)
	return err
}

func renderToHtml(c interface{}) template.HTML {
	var buf bytes.Buffer
	r := c.(chartrender.Renderer)
	err := r.Render(&buf)
	if err != nil {
		log.Printf("Failed to render chart: %s", err)
		return ""
	}

	return template.HTML(buf.String())
}
