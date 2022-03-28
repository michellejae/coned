package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/michellejae/coned/models"
	"github.com/michellejae/coned/views"
)

var monthly views.View

type MonthlyPage struct {
	HTML  *template.HTML
	Month models.NewBill
}

func GenerateAndGraph(w http.ResponseWriter, r *http.Request) {
	file := make(map[string]string)
	file["dec"] = "data/December_2021_offers.csv"
	file["jan"] = "data/January_2022_offers.csv"

	id := r.URL.Query().Get("month")

	models.OpenFile(file[id])

	monthly = *views.NewView("giraffe", "views/monthly.html")

	options := opts.BarData{}

	data := make([]opts.BarData, 0)

	// looping thru all esco's and creating a struct of bar data for each source
	// this way i can set each sources data styles individually
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

		if val.Name == models.Min {
			itemStyle.Color = "#FCAF56"
			options.ItemStyle = &itemStyle

		} else if val.Name == "Consolidated Edison Company of New York, Inc." {
			itemStyle.Color = "#eb6d6d"
			options.ItemStyle = &itemStyle
		} else {
			itemStyle.Color = "#b39ae0"
			options.ItemStyle = &itemStyle

		}

		// append each esco bar data struct to slice of bar dataw
		data = append(data, options)
	}

	// } else if models.Min == "Consolidated Edison Company of New York, Inc." {
	// 	itemStyle.Color = "linear-gradient(#e66465, #9198e5);"
	// 	options.ItemStyle = &itemStyle

	bar := charts.NewBar()

	bar.AddSeries("Totals", data)

	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    fmt.Sprintf("My %v Energy Bills per ESCO", models.Month.Month),
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

	bar.Renderer = views.NewSnippetRenderer(bar, bar.Validate)
	var htmlSnippet template.HTML = views.RenderToHtml(bar)

	tmplData := MonthlyPage{
		HTML:  &htmlSnippet,
		Month: models.Month,
	}

	monthly.Render(w, tmplData)

}
