package upworkfeed

import (
	"github.com/gdamore/tcell/v2"
)

func (widget *Widget) initializeKeyboardControls() {
	widget.InitializeHelpTextKeyboardControl(widget.ShowHelp)
	widget.InitializeRefreshKeyboardControl(widget.Refresh)

	widget.SetKeyboardChar("j", widget.Next, "Select next item")
	widget.SetKeyboardChar("k", widget.Prev, "Select previous item")

	widget.SetKeyboardKey(tcell.KeyDown, widget.Next, "Select next item")
	widget.SetKeyboardKey(tcell.KeyUp, widget.Prev, "Select previous item")
	widget.SetKeyboardKey(tcell.KeyEnter, widget.OpenSelected, "Edit item")
}

func (widget *Widget) OpenSelected() {
	if widget.upworkRss == nil || widget.Selected < 0 || widget.Selected >= len(widget.upworkRss.Channel.Items) {
		return
	}
	widget.openDetailsModal()
}
