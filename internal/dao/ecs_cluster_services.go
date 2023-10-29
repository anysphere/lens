package dao

import (
	"context"
	"fmt"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/rs/zerolog/log"
)

type ECSClusterServices struct {
	Accessor
	ctx context.Context
}

func (ecsClusterServices *ECSClusterServices) Init(ctx context.Context) {
	ecsClusterServices.ctx = ctx
}

func (ecsClusterServices *ECSClusterServices) List(ctx context.Context) ([]Object, error) {
	cfg, ok := ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	// clusterName := fmt.Sprintf("%v", ctx.Value(internal.ClusterName))
	
	listServicesResp, err := aws.ListEcsClusterServices(cfg, "cursor-cluster")
	objs := make([]Object, len(listServicesResp))
	for i, obj := range listServicesResp {
		objs[i] = obj
	}

	return objs, err
}

func (ecsClusterServices *ECSClusterServices) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}