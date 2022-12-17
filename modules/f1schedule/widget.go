package f1schedule

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/view"
)

const (
	__template_race_data_base = `[black:orange] {{.Race.Name}} [-:-]

[black:white]  SEASON  : {{.Race.Season}} 
  ROUND   : {{.Race.Round}} 
  CIRCUIT : {{.Race.Circuit.Name}} 
  COUNTRY : {{.Race.Circuit.Location.Country}} 
{{if .IsSprintFormat}}{{template "template_format_sprint" }}{{else}}{{template "template_format_normal" .}}{{end}}
  RACE    : {{.Race.DateTime.GetFormattedTime}} [-:-]
`
	__template_race_data_normal = `  FP1     : {{.Race.FirstPractice.GetFormattedTime}} 
  FP2     : {{.Race.SecondPractice.GetFormattedTime}} 
  FP3     : {{.Race.ThirdPractice.GetFormattedTime}} 
  QUALY   : {{.Race.Qualifying.GetFormattedTime}} `
	__template_race_data_sprint = `  FP1     : {{.Race.FirstPractice.GetFormattedTime}} 
  QUALY   : {{.Race.Qualifying.GetFormattedTime}} 
  FP2     : {{.Race.SecondPractice.GetFormattedTime}} 
  SPRINT  : {{.Race.Sprint.GetFormattedTime}} `
)

// Widget is the container for your module's data
type Widget struct {
	view.TextWidget

	settings     *Settings
	raceData     *RaceDataTemplateContext
	err          error
	templateView *template.Template
}

// NewWidget creates and returns an instance of Widget
func NewWidget(tviewApp *tview.Application, redrawChan chan bool, pages *tview.Pages, settings *Settings) *Widget {
	templateView := template.New("base")
	tb, _ := templateView.Parse(__template_race_data_base)
	tb.New("template_format_normal").Parse(__template_race_data_normal)
	tb.New("template_format_sprint").Parse(__template_race_data_sprint)
	widget := Widget{
		TextWidget: view.NewTextWidget(tviewApp, redrawChan, pages, settings.Common),

		settings:     settings,
		templateView: templateView,
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
	schedule, err := F1RaceCurrentSchedule()
	if err != nil {
		widget.raceData = nil
		widget.err = err
		return
	}
	raceData, err := schedule.Data.RaceTable.GetRaceDataForDisplay()
	if err != nil {
		widget.raceData = nil
		widget.err = err
		return
	}
	widget.raceData = raceData
	widget.err = nil
}

func (widget *Widget) Render() {
	widget.Redraw(widget.display)
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) content() string {
	if widget.err != nil {
		return fmt.Sprintf("[red]%s[-]", widget.err)
	}
	str := ""
	var buf bytes.Buffer
	err := widget.templateView.Execute(&buf, widget.raceData)
	if err != nil {
		str += fmt.Sprintf("[red]%s[-]", err.Error())
	} else {
		str += buf.String()
	}
	return str
}

func (widget *Widget) display() (string, string, bool) {
	return widget.CommonSettings().Title, widget.content(), true
}
