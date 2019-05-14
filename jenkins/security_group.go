package jenkins

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func CreateSecurityGroup(svc *ec2.EC2, vpc *ec2.Vpc, name string, description string, cidr string, protocol string, from int64, to int64) (*ec2.CreateSecurityGroupOutput, error) {
	securityGroup, err := svc.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
		VpcId: vpc.VpcId,
		Description: aws.String(description),
		GroupName: aws.String(name),

	})

	if err != nil {
		return nil, err
	}

	// Create tags

	_, _ = svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{securityGroup.GroupId},
		Tags: []*ec2.Tag{
			{Key: aws.String("Name"), Value: aws.String(name)},
		},
	})

	// Create ingress and egress

	_, err = svc.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: securityGroup.GroupId,
		CidrIp: aws.String(cidr),
		FromPort: &from,
		ToPort: &to,
		IpProtocol: &protocol,
	})

	if err != nil {
		return nil, err
	}

	_, err = svc.AuthorizeSecurityGroupEgress(&ec2.AuthorizeSecurityGroupEgressInput{
		GroupId: securityGroup.GroupId,
		IpPermissions: []*ec2.IpPermission{
			{
				FromPort: aws.Int64(1),
				IpProtocol: aws.String("-1"),
				ToPort: aws.Int64(60000),
				IpRanges: []*ec2.IpRange {
					{
						CidrIp: aws.String("0.0.0.0/0"),
					},
				},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return securityGroup, nil

}
