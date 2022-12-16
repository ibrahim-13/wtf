package upworkfeed

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/utils"
	"github.com/wtfutil/wtf/view"
)

const (
	offscreen               = -1000
	modalWidth              = 100
	modalHeight             = 20
	__format_title          = "[white]%s[-] [orange]%s[-]"
	__format_title_selected = "[:blue][black]%s %s[-][-:-]"
	__template_item_details = `
TITLE        : [black:orange]{{.Title}}[-:-]
PUBLISH DATE : [black:white]{{.PublishDate}}[-:-]
COUNTRY      : [black:yellow:b]{{.Country}}[-:-:-]
RATE         : [black:red:b]{{.Rate}}[-:-:-]
CATEGORY     : [black:green:b]{{.Category}}[-:-:-]
SKILLS       : {{range .SkillsArr}}[white:purple]{{.}}[-:-] {{end}}	
URL          : [white:blue]{{.Link}}[-:-]

[red]Press enter to exit...[-]
`
)

// Widget is the container for your module's data
type Widget struct {
	view.ScrollableWidget

	settings            *Settings
	upworkRss           *UpworkRss
	errMsg              error
	app                 *tview.Application
	pages               *tview.Pages
	templateItemDetails *template.Template
}

// NewWidget creates and returns an instance of Widget
func NewWidget(tviewApp *tview.Application, redrawChan chan bool, pages *tview.Pages, settings *Settings) *Widget {
	temp := template.New("item_details")
	temp, _ = temp.Parse(__template_item_details)
	widget := Widget{
		ScrollableWidget: view.NewScrollableWidget(tviewApp, redrawChan, pages, settings.Common),

		settings:            settings,
		app:                 tviewApp,
		pages:               pages,
		templateItemDetails: temp,
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
		widget.errMsg = err
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

func (widget *Widget) openDetailsModal() {
	i := widget.GetSelected()
	if widget.upworkRss == nil || i < 0 || i >= len(widget.upworkRss.Channel.Items) {
		return
	}
	txtView := tview.NewTextView()
	txtView.
		SetDynamicColors(true).
		SetWrap(true).
		SetWordWrap(true)

	item := widget.upworkRss.Channel.Items[i]

	writer := txtView.BatchWriter()
	widget.templateItemDetails.Execute(writer, item)
	writer.Close()

	frame := tview.NewFrame(txtView)
	frame.SetBorders(0, 0, 0, 0, 0, 0)
	frame.SetRect(offscreen, offscreen, modalWidth, modalHeight)
	frame.SetBorder(true)
	frame.SetBorders(1, 1, 0, 0, 1, 1)

	drawFunc := func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		w, h := screen.Size()
		frame.SetRect((w/2)-(width/2), (h/2)-(height/2), width, height)
		return x, y, width, height
	}
	frame.SetDrawFunc(drawFunc)
	frame.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			widget.pages.RemovePage("modal")
			widget.app.SetFocus(widget.View)
			widget.Render()
		}
		return event
	})

	widget.pages.AddPage("modal", frame, false, true)
	widget.app.SetFocus(frame)

	// Tell the app to force redraw the screen
	widget.Base.RedrawChan <- true
}
