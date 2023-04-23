package f1schedule

import (
	"errors"
	"fmt"
	"time"

	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/modules/f1schedule/f1api"
	"github.com/wtfutil/wtf/view"
)

// Widget is the container for your module's data
type Widget struct {
	view.TextWidget

	settings       *Settings
	nextRace       *f1api.Race
	nextRaceEvents []f1api.RaceEvent
	err            error
	f1api          *f1api.F1Api
}

// NewWidget creates and returns an instance of Widget
func NewWidget(tviewApp *tview.Application, redrawChan chan bool, pages *tview.Pages, settings *Settings) *Widget {
	year := fmt.Sprint(time.Now().Year())
	widget := Widget{
		TextWidget: view.NewTextWidget(tviewApp, redrawChan, pages, settings.Common),

		settings: settings,
		f1api:    f1api.NewF1Api(year, f1api.NewApiCacheLocal(1*time.Hour)),
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

func (widget *Widget) Load() {
	race, err := widget.f1api.GetRaceList()
	currentTime := time.Now()
	if err != nil {
		widget.nextRace = nil
		widget.err = err
		return
	}
	for i := range race {
		if race[i].StartDateTime.After(currentTime) {
			widget.nextRace = &race[i]
			break
		}
	}
	if widget.nextRace == nil {
		widget.err = errors.New("next race: nil")
		return
	}
	events, err := widget.f1api.GetRaceEventList(widget.nextRace.Url)
	if err != nil {
		widget.nextRaceEvents = nil
		widget.err = err
		return
	}
	widget.nextRaceEvents = events[0].SubEvents
}

func (widget *Widget) Render() {
	widget.Redraw(widget.display)
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) content() string {
	if widget.err != nil {
		return fmt.Sprintf("[red]%s[-]", widget.err)
	}
	nr := widget.nextRace
	str := fmt.Sprintf("[black:orange] %s [-:-]\n", nr.Description)
	str += fmt.Sprintf("[black:white]  SEASON  : %s \n", widget.f1api.GetYear())
	str += fmt.Sprintf("  CIRCUIT : %s \n", nr.Location.Name)
	str += fmt.Sprintf("  ADDRESS: %s [-:-]\n", nr.Location.Address)
	nre, isNextFound, currentTime := widget.nextRaceEvents, false, time.Now()
	for i := range nre {
		if !isNextFound && currentTime.Before(nre[i].EndDateTime) {
			str += fmt.Sprintf("[black:green] %s : %s [-:-:-]\n", printTime(&nre[i].StartDateTime), nre[i].Name)
			isNextFound = true
		} else {
			str += fmt.Sprintf("[black:white] %s : %s [-:-:-]\n", printTime(&nre[i].StartDateTime), nre[i].Name)
		}
	}
	return str
}

func (widget *Widget) display() (string, string, bool) {
	return widget.CommonSettings().Title, widget.content(), true
}

func printTime(time *time.Time) string {
	return time.Local().Format("03:04PM 02/01/06")
}
