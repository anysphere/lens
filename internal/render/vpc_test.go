package render

import (
	"testing"

	"github.com/anysphere/lens/internal/aws"
	"github.com/stretchr/testify/assert"
)

func TestVpcRender(t *testing.T) {
	resp := aws.VpcResp{VpcId: "vpc-1", OwnerId: "000000000000", CidrBlock: "172.31.0.0/16", InstanceTenancy: "default", State: "available"}
	var vpc VPC

	r := NewRow(5)
	err := vpc.Render(resp, "vpc", &r)

	assert.Nil(t, err)
	assert.Equal(t, "vpc", r.ID)

	e := Fields{"vpc-1", "000000000000", "172.31.0.0/16", "default", "available"}
	assert.Equal(t, e, r.Fields[0:])

	headers := vpc.Header()
	assert.Equal(t, 0, headers.IndexOf("VPC-Id", false))
	assert.Equal(t, 1, headers.IndexOf("Owner-Id", false))
	assert.Equal(t, 2, headers.IndexOf("Cidr Block", false))
	assert.Equal(t, 3, headers.IndexOf("Instance Tenancy", false))
	assert.Equal(t, 4, headers.IndexOf("VPC-State", false))
}
