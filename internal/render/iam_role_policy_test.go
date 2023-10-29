package render

import (
	"testing"

	"github.com/anysphere/lens/internal/aws"
	"github.com/stretchr/testify/assert"
)

func TestIamRolePolicyRender(t *testing.T) {
	resp := aws.IamRolePolicyResponse{PolicyArn: "arn:aws:iam:00000000000:policy/Buddy-ec2-policy", PolicyName: "Buddy-ec2-policy"}
	var iamRolePolicy IamRolePloicy

	r := NewRow(2)
	err := iamRolePolicy.Render(resp, "iamRolePolicy", &r)

	assert.Nil(t, err)
	assert.Equal(t, "iamRolePolicy", r.ID)

	e := Fields{"arn:aws:iam:00000000000:policy/Buddy-ec2-policy", "Buddy-ec2-policy"}
	assert.Equal(t, e, r.Fields[0:])

	headers := iamRolePolicy.Header()
	assert.Equal(t, 0, headers.IndexOf("Policy-ARN", false))
	assert.Equal(t, 1, headers.IndexOf("Policy-Name", false))
}
