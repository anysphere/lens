package dao

import (
	"context"
	"fmt"

	"github.com/anysphere/lens/internal"
	"github.com/anysphere/lens/internal/aws"
	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rs/zerolog/log"
)

type IAMUG struct {
	Accessor
	ctx context.Context
}

func (iamug *IAMUG) Init(ctx context.Context) {
	iamug.ctx = ctx
}

func (iamug *IAMUG) List(ctx context.Context) ([]Object, error) {
	cfg, ok := ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	usrGroup := aws.GetUserGroups(cfg)
	objs := make([]Object, len(usrGroup))
	for i, obj := range usrGroup {
		objs[i] = obj
	}
	return objs, nil
}

func (iamug *IAMUG) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}
