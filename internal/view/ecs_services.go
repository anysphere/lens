package view

import (
	"github.com/anysphere/lens/internal/ui"
	"github.com/gdamore/tcell/v2"
)

type ECSClusterViewer struct {
	clusterName string
	ResourceViewer
}

func NewECSClusterViewer(resource string) ResourceViewer {
	var ecs ECSClusterViewer
	ecs.clusterName = "cursor-cluster"
	ecs.ResourceViewer = NewBrowser(resource)
	ecs.AddBindKeysFn(ecs.bindKeys)
	return &ecs
}

func NewECSClusterViewerWithClusterName(resource, clusterName string) *ECSClusterViewer {
	var ecs ECSClusterViewer
	ecs.clusterName = clusterName
	ecs.ResourceViewer = NewBrowser(resource)
	ecs.AddBindKeysFn(ecs.bindKeys)
	return &ecs
}

func (ecs *ECSClusterViewer) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftB:    ui.NewKeyAction("Sort Service-Name", ecs.GetTable().SortColCmd("Service-Name", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", ecs.App().PrevCmd, false),
		// tcell.KeyEnter:  ui.NewKeyAction("View", ecs.enterCmd, false),
	})
}

// func (ecs *ECSClusterViewer) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
// 	clusterArn := ecs.GetTable().GetSelectedItem()
// 	if clusterArn != "" {
// 		o := NewEcsContainerInstances("ecs://", clusterArn)
// 		ecs.App().inject(o)
// 		o.GetTable().SetTitle(o.path)
// 	}
// 	return nil
// }
