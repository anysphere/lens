package dao

import (
	"context"
	"fmt"

	"github.com/anysphere/lens/internal"
	"github.com/anysphere/lens/internal/aws"
	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rs/zerolog/log"
)

type Lambda struct {
	Accessor
	ctx context.Context
}

func (l *Lambda) Init(ctx context.Context) {
	l.ctx = ctx
}

func (l *Lambda) List(ctx context.Context) ([]Object, error) {
	cfg, ok := ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	ins, err := aws.GetAllLambdaFunctions(cfg)
	objs := make([]Object, len(ins))
	for i, obj := range ins {
		objs[i] = obj
	}
	return objs, err
}
