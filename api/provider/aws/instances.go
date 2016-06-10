package aws

import (
	"os"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/convox/rack/api/structs"
)

func (p *AWSProvider) InstanceList() (structs.Instances, error) {
	is := make(map[string]structs.Instance)
	instances := structs.Instances{}

	res, err := p.ec2().DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("tag"),
				Values: []*string{aws.String(os.Getenv("RACK"))},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	for _, r := range res.Reservations {
		for _, i := range r.Instances {
			is[*i.InstanceId] = structs.Instance{
				Id:        *i.InstanceId,
				PrivateIp: *i.PrivateIpAddress,
				PublicIp:  *i.PublicIpAddress,
				Started:   *i.LaunchTime,
			}
		}
	}

	lres, err := p.ecs().ListContainerInstances(&ecs.ListContainerInstancesInput{
		Cluster: aws.String(os.Getenv("CLUSTER")),
	})
	if err != nil {
		return nil, err
	}

	dRes, err := p.ecs().DescribeContainerInstances(
		&ecs.DescribeContainerInstancesInput{
			Cluster:            aws.String(os.Getenv("CLUSTER")),
			ContainerInstances: lres.ContainerInstanceArns,
		},
	)
	if err != nil {
		return nil, err
	}

	for _, ci := range dRes.ContainerInstances {
		i := is[*ci.Ec2InstanceId]
		i.Processes = int(*ci.RunningTasksCount)
		is[*ci.Ec2InstanceId] = i
	}

	// collect all the instance info
	ids := []string{}
	for k := range is {
		ids = append(ids, k)
	}

	sort.Strings(ids)

	for _, k := range ids {
		instances = append(instances, is[k])
	}

	return instances, err

	// res, err := p.listContainerInstances(&ecs.ListContainerInstancesInput{
	// 	Cluster: aws.String(os.Getenv("CLUSTER")),
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// ecsRes, err := p.ecs().DescribeContainerInstances(
	// 	&ecs.DescribeContainerInstancesInput{
	// 		Cluster:            aws.String(os.Getenv("CLUSTER")),
	// 		ContainerInstances: res.ContainerInstanceArns,
	// 	},
	// )
	// if err != nil {
	// 	return nil, err
	// }

	// var instanceIds []*string
	// for _, i := range ecsRes.ContainerInstances {
	// 	instanceIds = append(instanceIds, i.Ec2InstanceId)
	// }

	// ec2Res, err := p.ec2().DescribeInstances(&ec2.DescribeInstancesInput{
	// 	Filters: []*ec2.Filter{
	// 		&ec2.Filter{Name: aws.String("instance-id"), Values: instanceIds},
	// 	},
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// ec2Instances := make(map[string]*ec2.Instance)
	// ec2Metrics := make(map[string]float64)

	// // collect instance data from EC2, and CPU Utilization from CloudWatch Metrics
	// for _, r := range ec2Res.Reservations {
	// 	for _, i := range r.Instances {
	// 		ec2Instances[*i.InstanceId] = i
	// 		ec2Metrics[*i.InstanceId] = 0.0

	// 		res, err := p.cloudwatch().GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
	// 			Dimensions: []*cloudwatch.Dimension{
	// 				&cloudwatch.Dimension{Name: aws.String("InstanceId"), Value: i.InstanceId},
	// 			},
	// 			EndTime:    aws.Time(time.Now()),
	// 			MetricName: aws.String("CPUUtilization"),
	// 			Namespace:  aws.String("AWS/EC2"),
	// 			Period:     aws.Int64(5 * 60), // seconds
	// 			StartTime:  aws.Time(time.Now().Add(time.Duration(-5) * time.Minute)),
	// 			Statistics: []*string{aws.String("Average")},
	// 		})
	// 		if err != nil {
	// 			continue
	// 		}

	// 		if len(res.Datapoints) > 0 {
	// 			ec2Metrics[*i.InstanceId] = *res.Datapoints[0].Average / 100.0
	// 		}
	// 	}
	// }

	// var instances structs.Instances

	// // Calculate memory metrics from ECS DescribeContainerInstances
	// // We can not collect CPU metrics since we are not yet using ECS CPU reservations
	// for _, i := range ecsRes.ContainerInstances {
	// 	var memory structs.InstanceResource

	// 	for _, r := range i.RegisteredResources {
	// 		switch *r.Name {
	// 		case "MEMORY":
	// 			memory.Total = int(*r.IntegerValue)
	// 		}
	// 	}

	// 	for _, r := range i.RemainingResources {
	// 		switch *r.Name {
	// 		case "MEMORY":
	// 			memory.Free = int(*r.IntegerValue)
	// 			memory.Used = memory.Total - memory.Free
	// 		}
	// 	}

	// 	// find the matching Instance from the EC2 response
	// 	ec2Instance := ec2Instances[*i.Ec2InstanceId]

	// 	// build up the struct
	// 	instance := structs.Instance{
	// 		Cpu:    ec2Metrics[*i.Ec2InstanceId],
	// 		Memory: memory.PercentUsed(),
	// 		Id:     *i.Ec2InstanceId,
	// 	}

	// 	if i.AgentConnected != nil {
	// 		instance.Agent = *i.AgentConnected
	// 	}

	// 	if ec2Instance != nil {
	// 		if ec2Instance.PrivateIpAddress != nil {
	// 			instance.PrivateIp = *ec2Instance.PrivateIpAddress
	// 		}

	// 		if ec2Instance.PublicIpAddress != nil {
	// 			instance.PublicIp = *ec2Instance.PublicIpAddress
	// 		}

	// 		if ec2Instance.LaunchTime != nil {
	// 			instance.Started = *ec2Instance.LaunchTime
	// 		}
	// 	}

	// 	if i.RunningTasksCount != nil {
	// 		instance.Processes = int(*i.RunningTasksCount)
	// 	}

	// 	if i.Status != nil {
	// 		instance.Status = strings.ToLower(*i.Status)
	// 	}

	// 	instances = append(instances, instance)
	// }

	// return instances, nil
}
