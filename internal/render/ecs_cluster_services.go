package render

import (
	"fmt"
	"strconv"

	"github.com/derailed/tview"
	"github.com/one2nc/cloudlens/internal/aws"
)

type EcsClusterServices struct {
}

func (ecs EcsClusterServices) Header() Header {
	return Header{
		HeaderColumn{Name: "ServiceName", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Status", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "ServiceArn", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "ServiceType", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "TasksTotal", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "RunningCount", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "PendingCount", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "LastDeployment", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Revision", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
	}
}

func (ecs EcsClusterServices) Render(o interface{}, ns string, row *Row) error {
	ecsResp, ok := o.(aws.EcsServiceResp)
	if !ok {
		return fmt.Errorf("expected EcsServiceResp, but got %T", o)
	}

	row.ID = ns
	row.Fields = Fields{
		ecsResp.ServiceName,
		ecsResp.Status,
		ecsResp.ServiceArn,
		ecsResp.ServiceType,
		strconv.Itoa(int(ecsResp.TasksTotal)),
		strconv.Itoa(int(ecsResp.RunningCount)),
	}
	return nil
}
