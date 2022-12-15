package upworkfeed

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/utils"
	"github.com/wtfutil/wtf/view"
)

const (
	__format_title          = "[white]%s[-] [orange]%s[-]"
	__format_title_selected = "[:blue][black]%s %s[-][-:-]"
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
		ScrollableWidget: view.NewScrollableWidget(tviewApp, redrawChan, pages, settings.Common),

		settings: settings,
	}
	widget.SetRenderFunction(widget.Render)
	widget.initializeKeyboardControls()

	return &widget
}

/* -------------------- Exported Functions -------------------- */

// Refresh updates the onscreen contents of the widget
func (widget *Widget) Refresh() {
	widget.load()
	// The last call should always be to the display function
	widget.SetItemCount(len(widget.upworkRss.Channel.Items))
	widget.Render()
}

func (widget *Widget) Render() {
	widget.Redraw(widget.content)
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

func (widget *Widget) content() (string, string, bool) {
	str := ""
	if widget.upworkRss == nil {
		return widget.CommonSettings().Title, widget.getEmtpyMsg(), false
	}
	for i, item := range widget.upworkRss.Channel.Items {
		str += widget.getTitle(&item, i)
	}
	return widget.CommonSettings().Title, str, false
}

func (widget *Widget) getTitle(item *UpworkItem, index int) string {
	title := fmt.Sprintf(__format_title, item.PublishDate, item.Title)
	if index == widget.Selected {
		title = fmt.Sprintf(__format_title_selected, item.PublishDate, item.Title)
		return utils.HighlightableHelper(widget.View, title, index, 1)
	}
	return utils.HighlightableHelper(widget.View, title, index, 1)
}

func (widget *Widget) getEmtpyMsg() string {
	return utils.HighlightableHelper(widget.View, "[red]ERR: NO FEED[-]", 0, 1)
}
