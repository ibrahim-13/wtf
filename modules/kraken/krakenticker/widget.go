package krakenticker

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/rivo/tview"
	kraken "github.com/wtfutil/wtf/modules/kraken/common"
	"github.com/wtfutil/wtf/view"
)

const (
	__template_price_details = `
   OPEN : [green]{{.OpeningPriceToday}}[-]
   HIGH : [green]{{.HighToday}}[-]
    LOW : [green]{{.LowToday}}[-]
   LAST : [green]{{.LastTradeClosedPrice}}[-]
AVERAGE : [green]{{.VolumeWeightedAveragePriceToday}}[-]
`
)

// Widget is the container for your module's data
type Widget struct {
	view.TextWidget

	settings        *Settings
	ticker          kraken.KrakenDataTicker
	reqErr          []string
	err             error
	templateDetails *template.Template
}

// NewWidget creates and returns an instance of Widget
func NewWidget(tviewApp *tview.Application, redrawChan chan bool, pages *tview.Pages, settings *Settings) *Widget {
	temp := template.New("price_details")
	temp, _ = temp.Parse(__template_price_details)
	widget := Widget{
		TextWidget: view.NewTextWidget(tviewApp, redrawChan, pages, settings.Common),

		settings:        settings,
		templateDetails: temp,
	}

	return &widget
}

/* -------------------- Exported Functions -------------------- */

// Refresh updates the onscreen contents of the widget
func (widget *Widget) Refresh() {
	widget.Load()
	// The last call should always be to the display function
	widget.Render()
}

func (widget *Widget) Render() {
	widget.Redraw(widget.display)
}

func (widget *Widget) Load() {
	if len(widget.settings.assetPairs) < 1 {
		return
	}
	var assetPairs []string
	for _, v := range widget.settings.assetPairs {
		assetPairs = append(assetPairs, v.(string))
	}
	data, err := kraken.GetKrakenTicker(assetPairs...)
	if err != nil {
		widget.err = err
		widget.reqErr = nil
		widget.ticker = nil
		return
	}
	if data.Error != nil && len(data.Error) > 0 {
		widget.err = nil
		widget.ticker = nil
		widget.reqErr = data.Error
		return
	}
	widget.err = nil
	widget.reqErr = nil
	widget.ticker = data.Result
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) content() string {
	str := ""
	if widget.err != nil {
		str = fmt.Sprintf("[red]%s[-]", widget.err)
	} else if widget.reqErr != nil {
		str = fmt.Sprintf("[red]%s[-]", strings.Join(widget.reqErr, ","))
	} else if widget.ticker == nil || len(widget.ticker) < 1 {
		str = fmt.Sprintf("[red]%s[-]", "EMPTY DATA")
	} else {
		for k, v := range widget.ticker {
			str += fmt.Sprintf("[red]%s[-]:", k)
			var buf bytes.Buffer
			widget.templateDetails.Execute(&buf, v)
			str += buf.String() + "\n"
		}
	}
	return str
}

func (widget *Widget) display() (string, string, bool) {
	return widget.CommonSettings().Title, widget.content(), true
}
