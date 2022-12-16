package krakenstatus

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
	kraken "github.com/wtfutil/wtf/modules/kraken/common"
	"github.com/wtfutil/wtf/utils"
	"github.com/wtfutil/wtf/view"
)

// Widget is the container for your module's data
type Widget struct {
	view.TextWidget

	settings *Settings
	status   kraken.KrakenStatus
	reqErr   []string
	err      error
}

// NewWidget creates and returns an instance of Widget
func NewWidget(tviewApp *tview.Application, redrawChan chan bool, pages *tview.Pages, settings *Settings) *Widget {
	widget := Widget{
		TextWidget: view.NewTextWidget(tviewApp, redrawChan, pages, settings.Common),

		settings: settings,
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
	data, err := kraken.GetKrakenStatus()
	if err != nil {
		widget.err = err
		widget.reqErr = nil
		widget.status = ""
		return
	}
	if data.Error != nil && len(data.Error) > 0 {
		widget.err = nil
		widget.status = ""
		widget.reqErr = data.Error
		return
	}
	widget.err = nil
	widget.reqErr = nil
	widget.status = data.Result.Status
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) content() string {
	str := ""
	if widget.err != nil {
		str = utils.HighlightableHelper(widget.View, fmt.Sprintf("[red]%s[-]", widget.err), 0, 1)
	} else if widget.reqErr != nil {
		str = utils.HighlightableHelper(widget.View, fmt.Sprintf("[red]%s[-]", strings.Join(widget.reqErr, ",")), 0, 1)
	} else if widget.status == "" {
		str = utils.HighlightableHelper(widget.View, fmt.Sprintf("[red]%s[-]", "EMPTY STATUS"), 0, 1)
	} else {
		color := ""
		if widget.status == kraken.KrakenStatusOnline {
			color = "[black:green:b]"
		} else {
			color = "[black:red:b]"
		}
		str += fmt.Sprintf("Status : %s%s[-:-:-]\n", color, widget.status)
		str += getStatusDescription(widget.status)
	}
	return str
}

func (widget *Widget) display() (string, string, bool) {
	return widget.CommonSettings().Title, widget.content(), true
}

func getStatusDescription(status kraken.KrakenStatus) string {
	switch status {
	case kraken.KrakenStatusOnline:
		return "👌 Operating normally"
	case kraken.KrakenStatusMaintenance:
		return "💥 Exchange offline"
	case kraken.KrakenStatusCancelOnly:
		return "❌ New order/trade ✅ Cancel"
	case kraken.KrakenStatusPostOnly:
		return "❌ Trade ✅ Post-only limit order"
	default:
		return "🚨 Unknown"
	}
}
