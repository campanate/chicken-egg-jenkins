package jenkins

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func CreateInstance(svc *ec2.EC2, name string, ami string, keyName string, subnet_id *string, securityGroups []*string) (*ec2.Reservation, error) {
	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		InstanceType: aws.String("t2.micro"),
		SecurityGroupIds: securityGroups,
		SubnetId: subnet_id,
		KeyName: aws.String(keyName),
		ImageId: aws.String(ami),
		MinCount: aws.Int64(1),
		MaxCount: aws.Int64(1),

	})

	if err != nil {
		fmt.Print(err)
		return nil, err
	}


	_, _ = svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{runResult.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{Key: aws.String("Name"), Value: aws.String(name)},
		},

	})

	if err != nil {
		return nil, err
	}

	return runResult, nil

}
