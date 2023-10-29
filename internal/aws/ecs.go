package aws

import (
	"context"
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/rs/zerolog/log"
)

func ListEcsClusters(cfg aws.Config) ([]EcsClusterResp, error) {
	ecsClient := ecs.NewFromConfig(cfg)
	resultListClusters, err := ecsClient.ListClusters(context.TODO(), nil)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error fetching ECS Clusters: %v", err))
		return nil, err
	}
	var ecsClusterArns []string
	for _, cluster := range resultListClusters.ClusterArns {
		ecsClusterArns = append(ecsClusterArns, cluster)
	}
	describedClusters, err := DescribeEcsClusters(ecsClient, ecsClusterArns)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error describing ECS Clusters"))
		return nil, err
	}
	var detailedClusters []EcsClusterResp
	for _, cluster := range describedClusters.Clusters {
		c := &EcsClusterResp{ClusterName: *cluster.ClusterName, Status: *cluster.Status, RunningTasksCount: fmt.Sprint(cluster.RunningTasksCount), ClusterArn: *cluster.ClusterArn}
		detailedClusters = append(detailedClusters, *c)
	}
	return detailedClusters, nil
}

func DescribeEcsClusters(ecsClient *ecs.Client, clusters []string) (ecs.DescribeClustersOutput, error) {
	detailedClusters, err := ecsClient.DescribeClusters(context.TODO(), &ecs.DescribeClustersInput{Clusters: clusters})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error describing ECS Clusters"))
		return ecs.DescribeClustersOutput{}, err
	}
	return *detailedClusters, nil
}

func ListEcsClusterServices(cfg aws.Config, clusterName string) ([]EcsServiceResp, error) {
	ecsClient := ecs.NewFromConfig(cfg)
	input := &ecs.ListServicesInput{
		Cluster: aws.String(clusterName),
		MaxResults: aws.Int32(50),
	}
	resultListServices, err := ecsClient.ListServices(context.TODO(), input)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error fetching ECS Services: %v", err))
		return nil, err
	}
	var ecsServiceArns []string
	for _, service := range resultListServices.ServiceArns {
		ecsServiceArns = append(ecsServiceArns, service)
	}
	log.Info().Msg("Logging all ECS Services")
	log.Info().Msg(fmt.Sprintf("ECS Client: %v", ecsClient))
	log.Info().Msg(fmt.Sprintf("List Services Input: %v", input))
	log.Info().Msg(fmt.Sprintf("Result List Services: %v", resultListServices))
	log.Info().Msg(fmt.Sprintf("ECS Service ARNs: %v", ecsServiceArns))

	var describedServices ecs.DescribeServicesOutput
	for i := 0; i < len(ecsServiceArns); i += 9 {
		end := i + 9
		if end > len(ecsServiceArns) {
			end = len(ecsServiceArns)
		}
		batch, err := DescribeEcsServices(ecsClient, ecsServiceArns[i:end], clusterName)

		if err != nil {
			log.Info().Msg(fmt.Sprintf("Error describing ECS Services: %v", err))
			return nil, err
		}
		describedServices.Services = append(describedServices.Services, batch.Services...)
	}

	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error describing ECS Services"))
		return nil, err
	}
	var detailedServices []EcsServiceResp

	for _, service := range describedServices.Services {
		strategy := string(service.SchedulingStrategy)
		deps := service.Deployments
	// Sort deployments by last updated
		sort.Slice(deps, func(i, j int) bool {
			return deps[i].UpdatedAt.After(*deps[j].UpdatedAt)
		})

		mostRecentDeployment := deps[0]
		mostRecentDeploymentState := mostRecentDeployment.RolloutState

		s := &EcsServiceResp{
			ServiceName:    *service.ServiceName,
			Status:         *service.Status,
			RunningCount:   service.RunningCount,
			PendingCount:   service.PendingCount,
			ServiceArn:     *service.ServiceArn,
			ServiceType:    strategy,
			TasksTotal:     service.DesiredCount,
			LastDeployment: string(mostRecentDeploymentState),
			// Revision:       mostRecentDeployment.TaskDefinition,
		}
		detailedServices = append(detailedServices, *s)
	}
	return detailedServices, nil
}

func DescribeEcsServices(ecsClient *ecs.Client, services []string, cluster string) (ecs.DescribeServicesOutput, error) {
	detailedServices, err := ecsClient.DescribeServices(context.TODO(), &ecs.DescribeServicesInput{Services: services, Cluster: &cluster})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error describing ECS Services"))
		return ecs.DescribeServicesOutput{}, err
	}
	return *detailedServices, nil
}
