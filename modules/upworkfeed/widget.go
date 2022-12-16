package upworkfeed

import (
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/utils"
	"github.com/wtfutil/wtf/view"
)

const (
	offscreen               = -1000
	modalWidth              = 100
	modalHeight             = 40
	__format_title          = "[white]%s[-] [orange]%s[-]"
	__format_title_selected = "[:blue][black]%s %s[-][-:-]"
	__format_error          = "[red]%s[-]"
	__template_item_details = `
TITLE        : [black:orange]{{.Title}}[-:-]
PUBLISH DATE : [black:white]{{.PublishDate}}[-:-]
COUNTRY      : [black:yellow:b]{{.Country}}[-:-:-]
RATE         : [black:red:b]{{.Rate}}[-:-:-]
CATEGORY     : [black:green:b]{{.Category}}[-:-:-]
SKILLS       : {{range .SkillsArr}}[white:purple]{{.}}[-:-] {{end}}	
URL          : [white:blue]{{.Link}}[-:-]

DESCRIPTION  :
               {{.ShortDescription}}

Press [white:red:b] enter ↩ [-:-:-] to go back...
`
)

// Widget is the container for your module's data
type Widget struct {
	view.ScrollableWidget

	settings            *Settings
	upworkRss           *UpworkRss
	err                 error
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
	widget.LoadRssFeed()
	if widget.upworkRss != nil {
		widget.SetItemCount(len(widget.upworkRss.Channel.Items))
	}
	widget.Render()
}

func (widget *Widget) Render() {
	widget.Redraw(widget.content)
}

/* -------------------- Unexported Functions -------------------- */
func (widget *Widget) LoadRssFeed() {
	feed, err := fetch(widget.settings.feedUrl)
	// feed, err := fetchLocalFile("feed.rss")

	widget.upworkRss = feed
	widget.err = err
}

// func fetchLocalFile(filename string) (*UpworkRss, error) {
// 	wd, _ := os.Getwd()
// 	f, _ := os.ReadFile(filepath.Join(wd, filename))
// 	rss, err := ParseXml(f)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return rss, nil
// }

func fetch(url string) (*UpworkRss, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	rss, err := ParseXml(body)
	if err != nil {
		return nil, err
	}
	return rss, nil
}

func (widget *Widget) content() (string, string, bool) {
	str := ""
	if widget.upworkRss == nil || widget.err != nil {
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
	msg := ""
	if widget.err != nil {
		msg = fmt.Sprintf(__format_error, widget.err)
	} else {
		msg = fmt.Sprintf(__format_error, "ERR: NIL FEED")

	}
	return utils.HighlightableHelper(widget.View, msg, 0, 1)
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
