package render

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/one2nc/cloud-lens/internal/aws"
)

type IAMU struct {
}

func (iamu IAMU) Header() Header {
	return Header{
		HeaderColumn{Name: "User-Id", SortIndicatorIdx: 6, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "User-Name", SortIndicatorIdx: 6, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "ARN", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Created-date", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: true},
	}
}

func (iamu IAMU) Render(o interface{}, ns string, row *Row) error {
	iamuResp, ok := o.(aws.IAMUSerResp)
	if !ok {
		return fmt.Errorf("Expected S3Resp, but got %T", o)
	}
	
	row.ID = ns
	row.Fields = Fields{
		iamuResp.UserId,
		iamuResp.UserName,
		iamuResp.ARN,
		iamuResp.CreationTime,
	}
	return nil
}
