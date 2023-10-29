package view

import (
	"context"

	"github.com/anysphere/lens/internal"
	"github.com/anysphere/lens/internal/ui"
	"github.com/gdamore/tcell/v2"
)

type IAMUG struct {
	ResourceViewer
}

// NewUG returns a new viewer.
func NewIAMUG(resource string) ResourceViewer {
	var iamug IAMUG
	iamug.ResourceViewer = NewBrowser(resource)
	iamug.AddBindKeysFn(iamug.bindKeys)
	return &iamug
}

func (iamug IAMUG) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftI:    ui.NewKeyAction("Sort Group-Id ", iamug.GetTable().SortColCmd("Group-Id", true), true),
		ui.KeyShiftN:    ui.NewKeyAction("Sort Group-Name", iamug.GetTable().SortColCmd("Group-Name", true), true),
		ui.KeyShiftD:    ui.NewKeyAction("Sort Created-Date", iamug.GetTable().SortColCmd("Created-Date", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", iamug.App().PrevCmd, false),
		tcell.KeyEnter:  ui.NewKeyAction("Group Users", iamug.enterCmd, false),
		ui.KeyShiftP:    ui.NewKeyAction("View", iamug.policyCmd, true),
	})
}

func (iamug *IAMUG) policyCmd(evt *tcell.EventKey) *tcell.EventKey {
	grpName := iamug.GetTable().GetSecondColumn()
	if grpName != "" {
		up := NewIamUserGroupPloicy("User Group Policy")
		ctx := context.WithValue(iamug.App().GetContext(), internal.GroupName, grpName)
		iamug.App().SetContext(ctx)
		iamug.App().Flash().Info("userName: " + grpName)
		iamug.App().inject(up)
	}
	return nil
}

func (iamug *IAMUG) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	grpName := iamug.GetTable().GetSecondColumn()
	if grpName != "" {
		gu := NewIamGroupUser("Group Users")
		ctx := context.WithValue(iamug.App().GetContext(), internal.GroupName, grpName)
		iamug.App().SetContext(ctx)
		iamug.App().Flash().Info("userName: " + grpName)
		iamug.App().inject(gu)
	}
	return nil
}
