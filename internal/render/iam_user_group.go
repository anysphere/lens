package render

import (
	"fmt"

	"github.com/anysphere/lens/internal/aws"
	"github.com/derailed/tview"
)

type IAMUG struct {
}

func (iamug IAMUG) Header() Header {
	return Header{
		HeaderColumn{Name: "Group-Id", SortIndicatorIdx: 6, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Group-Name", SortIndicatorIdx: 6, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "ARN", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Created-Date", SortIndicatorIdx: 8, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: true},
	}
}

func (iamug IAMUG) Render(o interface{}, ns string, row *Row) error {
	iamugResp, ok := o.(aws.IAMUSerGroupResp)
	if !ok {
		return fmt.Errorf("Expected iamugResp, but got %T", o)
	}

	row.ID = ns
	row.Fields = Fields{
		iamugResp.GroupId,
		iamugResp.GroupName,
		iamugResp.ARN,
		iamugResp.CreationTime,
	}
	return nil
}
