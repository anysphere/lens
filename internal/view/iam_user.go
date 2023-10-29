package view

import (
	"context"

	"github.com/anysphere/lens/internal"
	"github.com/anysphere/lens/internal/ui"
	"github.com/gdamore/tcell/v2"
)

type IAMU struct {
	ResourceViewer
}

// NewSG returns a new viewer.
func NewIAMU(resource string) ResourceViewer {
	var iamu IAMU
	iamu.ResourceViewer = NewBrowser(resource)
	iamu.AddBindKeysFn(iamu.bindKeys)
	return &iamu
}

func (iamu IAMU) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftI:    ui.NewKeyAction("Sort User-Id ", iamu.GetTable().SortColCmd("User-Id", true), true),
		ui.KeyShiftN:    ui.NewKeyAction("Sort User-Name", iamu.GetTable().SortColCmd("User-Name", true), true),
		ui.KeyShiftD:    ui.NewKeyAction("Sort Created-Date", iamu.GetTable().SortColCmd("Created-Date", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", iamu.App().PrevCmd, true),
		ui.KeyShiftP:    ui.NewKeyAction("View", iamu.viewUserPolicy, true),
	})
}

func (iamu *IAMU) viewUserPolicy(evt *tcell.EventKey) *tcell.EventKey {
	userName := iamu.GetTable().GetSecondColumn()
	if userName != "" {
		up := NewIamUserPloicy("User Policy")
		ctx := context.WithValue(iamu.App().GetContext(), internal.UserName, userName)
		iamu.App().SetContext(ctx)
		iamu.App().Flash().Info("userName: " + userName)
		iamu.App().inject(up)
	}
	return nil
}
