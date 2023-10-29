package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal/ui"
	"github.com/rs/zerolog/log"
)

type EcsClusters struct {
	ResourceViewer
}

func NewEcs(resource string) ResourceViewer {
	var ecs EcsClusters
	ecs.ResourceViewer = NewBrowser(resource)
	ecs.AddBindKeysFn(ecs.bindKeys)
	return &ecs
}
func (ecs *EcsClusters) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftB:    ui.NewKeyAction("Sort Cluster-Arn", ecs.GetTable().SortColCmd("Cluster-Arn", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", ecs.App().PrevCmd, false),
		tcell.KeyEnter:  ui.NewKeyAction("View", ecs.enterCmd, true),
	})
}

func (ecs *EcsClusters) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	log.Info().Msg("ecs enterCmd")
	clusterArn := ecs.GetTable().GetSelectedItem()
	log.Info().Str("clusterArn", clusterArn).Msg("ecs enterCmd")
	if clusterArn != "" {
		o := NewECSClusterViewerWithClusterName("ecs://", clusterArn)
		ecs.App().inject(o)
		log.Info().Str("clusterArn", clusterArn).Msg("injecting into app.")
		o.GetTable().SetTitle(o.clusterName)
	}
	return nil
}
