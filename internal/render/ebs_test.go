package render

import (
	"testing"

	"github.com/anysphere/lens/internal/aws"
	"github.com/stretchr/testify/assert"
)

func TestEBSRender(t *testing.T) {
	resp := aws.EBSResp{VolumeId: "vol-ebs-1", Size: "32", VolumeType: "gp2", State: "in-use", AvailabilityZone: "us-east-1e", Snapshot: "snapshot", CreationTime: "9:00:00"}
	var ebs EBS

	r := NewRow(7)
	err := ebs.Render(resp, "ebs", &r)

	assert.Nil(t, err)
	assert.Equal(t, "ebs", r.ID)

	e := Fields{"vol-ebs-1", "32", "gp2", "in-use", "us-east-1e", "snapshot", "9:00:00"}
	assert.Equal(t, e, r.Fields[0:])

	headers := ebs.Header()
	assert.Equal(t, 0, headers.IndexOf("Volume-Id", false))
	assert.Equal(t, 1, headers.IndexOf("Size", false))
	assert.Equal(t, 2, headers.IndexOf("Volume-Type", false))
	assert.Equal(t, 3, headers.IndexOf("State", false))
	assert.Equal(t, 4, headers.IndexOf("Availability-Zone", false))
	assert.Equal(t, 5, headers.IndexOf("Snapshot", false))
	assert.Equal(t, 6, headers.IndexOf("Creation-Time", false))
}
