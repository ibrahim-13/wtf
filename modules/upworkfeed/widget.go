package upworkfeed

import (
	"os"
	"path/filepath"

	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/utils"
	"github.com/wtfutil/wtf/view"
)

// Widget is the container for your module's data
type Widget struct {
	view.ScrollableWidget

	settings  *Settings
	upworkRss *UpworkRss
	errMsg    string
}

// NewWidget creates and returns an instance of Widget
func NewWidget(tviewApp *tview.Application, redrawChan chan bool, pages *tview.Pages, settings *Settings) *Widget {
	widget := Widget{
		ScrollableWidget: view.NewScrollableWidget(tviewApp, redrawChan, pages, settings.common),

		settings: settings,
	}

	return &widget
}

/* -------------------- Exported Functions -------------------- */

// Refresh updates the onscreen contents of the widget
func (widget *Widget) Refresh() {
	widget.load()
	// The last call should always be to the display function
	widget.display()
}

/* -------------------- Unexported Functions -------------------- */
func (widget *Widget) load() {
	wd, _ := os.Getwd()
	f, _ := os.ReadFile(filepath.Join(wd, "feed.rss"))
	rss, err := ParseXml(f)
	if err != nil {
		widget.errMsg = err.Error()
	}
	widget.upworkRss = rss
}

func (widget *Widget) content() string {
	con := ""
	con += utils.HighlightableHelper(widget.View, "[red]GG[-]", 0, 1)
	con += utils.HighlightableHelper(widget.View, "[green]WP[-]", 0, 1)
	return con
}

func (widget *Widget) display() {
	widget.Redraw(func() (string, string, bool) {
		return widget.CommonSettings().Title, widget.content(), false
	})
}
