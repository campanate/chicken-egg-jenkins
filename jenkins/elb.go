package jenkins

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

func CreateELB(subnets []*string, securityGroup []*string, zones []*string) (*elbv2.LoadBalancer, error){
	svc := elbv2.New(session.New())

	loadBalancer, err := svc.CreateLoadBalancer(&elbv2.CreateLoadBalancerInput{
		Subnets: subnets,
		SecurityGroups: securityGroup,
		Name: aws.String("JenkinsLoadBalancer"),
	})

	if err != nil {
		return nil, err
	}


	return loadBalancer.LoadBalancers[0], nil



}
