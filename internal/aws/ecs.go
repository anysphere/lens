package aws

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
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
		Cluster:    aws.String(clusterName),
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
	var servicesWg sync.WaitGroup
	describedServicesChan := make(chan ecs.DescribeServicesOutput, len(ecsServiceArns)/3+1)
	errChan := make(chan error, len(ecsServiceArns)/3+1)

	for i := 0; i < len(ecsServiceArns); i += 3 {
		servicesWg.Add(1)
		go func(i int) {
			defer servicesWg.Done()
			end := i + 3
			if end > len(ecsServiceArns) {
				end = len(ecsServiceArns)
			}
			batch, err := DescribeEcsServices(ecsClient, ecsServiceArns[i:end], clusterName)
			if err != nil {
				errChan <- err
				return
			}
			describedServicesChan <- batch
		}(i)
	}

	servicesWg.Wait()
	close(describedServicesChan)
	close(errChan)

	for batch := range describedServicesChan {
		describedServices.Services = append(describedServices.Services, batch.Services...)
		describedServices.Failures = append(describedServices.Failures, batch.Failures...)
	}

	if len(errChan) > 0 {
		for err := range errChan {
			log.Info().Msg(fmt.Sprintf("Error describing ECS Services: %v", err))
		}
		return nil, <-errChan
	}

	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error describing ECS Services"))
		return nil, err
	}
	var detailedServices []EcsServiceResp

	var wg sync.WaitGroup
	serviceRespChan := make(chan EcsServiceResp, len(describedServices.Services))

	for _, service := range describedServices.Services {
		wg.Add(1)
		go func(service ecsTypes.Service) {
			defer wg.Done()
			strategy := string(service.SchedulingStrategy)
			deps := service.Deployments
			// Sort deployments by last updated
			sort.Slice(deps, func(i, j int) bool {
				return deps[i].UpdatedAt.After(*deps[j].UpdatedAt)
			})

			mostRecentDeployment := deps[0]
			mostRecentDeploymentState := mostRecentDeployment.RolloutState

			taskDefinition := service.TaskDefinition

			log.Info().Msg(fmt.Sprintf("Task Definition: %v", string(*taskDefinition)))

			var dockerImage string
			taskDef, err := DescribeTaskDefinition(ecsClient, string(*taskDefinition))

			if err != nil {
				log.Info().Msg(fmt.Sprintf("Error describing ECS Task Definition"))
				dockerImage = ""
			} else {
				var dockerImages []string

				for _, container := range taskDef.TaskDefinition.ContainerDefinitions {
					if !strings.Contains(*container.Image, "datadog") &&
						!strings.Contains(*container.Image, "log") &&
						!strings.Contains(*container.Image, "fluent-bit") {
						imageParts := strings.SplitN(*container.Image, "/", 2)
						name := imageParts[0]
						if len(imageParts) > 1 {
							name = imageParts[len(imageParts)-1]
						}
						dockerImages = append(dockerImages, name)
					}
				}

				dockerImage = strings.Join(dockerImages, ", ")
			}

			s := EcsServiceResp{
				ServiceName:    *service.ServiceName,
				Status:         *service.Status,
				RunningCount:   service.RunningCount,
				PendingCount:   service.PendingCount,
				ServiceArn:     *service.ServiceArn,
				ServiceType:    strategy,
				TasksTotal:     service.DesiredCount,
				LastDeployment: string(mostRecentDeploymentState),
				DockerImage:    dockerImage,
				// Revision:       mostRecentDeployment.TaskDefinition,
			}
			serviceRespChan <- s
		}(service)
	}

	go func() {
		wg.Wait()
		close(serviceRespChan)
	}()

	for serviceResp := range serviceRespChan {
		detailedServices = append(detailedServices, serviceResp)
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

func DescribeTaskDefinition(ecsClient *ecs.Client, taskDefinition string) (ecs.DescribeTaskDefinitionOutput, error) {
	detailedTaskDefinition, err := ecsClient.DescribeTaskDefinition(context.TODO(), &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: &taskDefinition,
	})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error describing ECS Task Definition"))
		return ecs.DescribeTaskDefinitionOutput{}, err
	}
	return *detailedTaskDefinition, nil
}
