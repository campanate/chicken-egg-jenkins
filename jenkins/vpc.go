package jenkins

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func GetZones(svc *ec2.EC2)  ([]*ec2.AvailabilityZone, error) {
	resultAvalZones, err := svc.DescribeAvailabilityZones(nil)

	if err != nil {
		return nil, err
	}

	return resultAvalZones.AvailabilityZones, err

}

func CreateVpc(svc *ec2.EC2) (*ec2.Vpc, error) {

	newVpc := &ec2.CreateVpcInput{
		CidrBlock: aws.String("7.0.0.0/16"),
	}

	vpc, err := svc.CreateVpc(newVpc)

	if err != nil {
		fmt.Println("Error creating vpc")
		return nil, err
	}

	// Add Tags

	_, _ = svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{vpc.Vpc.VpcId},
		Tags: []*ec2.Tag{
			{
				Key: aws.String("Name"),
				Value: aws.String("JenkinsVpc"),
			},
		},
	})

	return vpc.Vpc, nil
}

func CreateInternetGateway(svc *ec2.EC2, vpc *ec2.Vpc) (*ec2.InternetGateway, error) {

	ig, err := svc.CreateInternetGateway(&ec2.CreateInternetGatewayInput{})

	if err != nil {
		return nil, err
	}

	// Add Tags

	_, _ = svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{ig.InternetGateway.InternetGatewayId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("JenkinsInternetGateway"),
			},
		},
	})

	// Attach internet gateway with vpc

	_, _ = svc.AttachInternetGateway(&ec2.AttachInternetGatewayInput{
		VpcId: vpc.VpcId,
		InternetGatewayId: ig.InternetGateway.InternetGatewayId,
	})

	return ig.InternetGateway, err
}

func CreateElasticIp(svc *ec2.EC2, vpc *ec2.Vpc) (*ec2.AllocateAddressOutput, error) {

	eip, err := svc.AllocateAddress(&ec2.AllocateAddressInput{
		Domain: vpc.VpcId,
	})

	if err != nil {
		return nil, err
	}

	// Add Tags

	_, _ = svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{eip.AllocationId},
		Tags: []*ec2.Tag{
			{Key: aws.String("Name"), Value: aws.String("JenkinsEIP")},
		},
	})

	return eip, nil
}

func CreateSubnet(svc *ec2.EC2, cidrblock string, public bool, zone *ec2.AvailabilityZone, vpc *ec2.Vpc) (*ec2.Subnet, error) {
	subnet, err := svc.CreateSubnet(&ec2.CreateSubnetInput{
		AvailabilityZone: zone.ZoneName,
		CidrBlock: aws.String(cidrblock),
		VpcId: vpc.VpcId,

	})

	if err != nil {
		return nil, err
	}

	_, err = svc.ModifySubnetAttribute(&ec2.ModifySubnetAttributeInput{
		MapPublicIpOnLaunch: &ec2.AttributeBooleanValue{Value: &public},
		SubnetId:            subnet.Subnet.SubnetId,
	})


	if err != nil {
		return nil, err
	}


	var tagName string

	if public {
		tagName = "PublicJenkinsSubnet"
	} else {
		tagName = "PrivateJenkinsSubnet"
	}

	// Create tags


	_, err = svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{subnet.Subnet.SubnetId},
		Tags: []*ec2.Tag{
			{Key: aws.String("Name"), Value: aws.String(tagName)},
		},
	})

	if err != nil {
		return nil, err
	}

	return subnet.Subnet, err
}

func CreateNatGateway(svc *ec2.EC2, eip *ec2.AllocateAddressOutput, subnet *ec2.Subnet) (*ec2.NatGateway, error) {
	natGateway, err := svc.CreateNatGateway(&ec2.CreateNatGatewayInput{
		SubnetId: subnet.SubnetId,
		AllocationId: eip.AllocationId,
	})

	if err != nil {
		return nil, err
	}

	_, _ = svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{natGateway.NatGateway.NatGatewayId},
		Tags: []*ec2.Tag{
			{Key: aws.String("Name"), Value: aws.String("JenkinsNatGateway")},
		},
	})

	fmt.Println("Waiting Nat Gateway be available...")
	_ = svc.WaitUntilNatGatewayAvailable(&ec2.DescribeNatGatewaysInput{
		NatGatewayIds: []*string{natGateway.NatGateway.NatGatewayId},
	})

	return natGateway.NatGateway, nil

}

func CreateRouteTable(svc *ec2.EC2, vpc *ec2.Vpc, public bool, subnet *ec2.Subnet, natGateway *ec2.NatGateway, internetGateway *ec2.InternetGateway) (*ec2.RouteTable, error) {
	routeTable, err := svc.CreateRouteTable(&ec2.CreateRouteTableInput{
		VpcId: vpc.VpcId,
	})

	if err != nil {
		return nil, err
	}

	// Add Tags

	var tagName string

	if public {
		tagName = "PublicRouteTable"
	} else {
		tagName = "PrivateRouteTable"
	}

	_, err = svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{routeTable.RouteTable.RouteTableId},
		Tags: []*ec2.Tag{
			{Key: aws.String("Name"), Value: aws.String(tagName)},
		},
	})

	if err != nil {
		return nil, err
	}

	// Add Route

	if public {
		_, err = svc.CreateRoute(&ec2.CreateRouteInput{
			RouteTableId: routeTable.RouteTable.RouteTableId,
			DestinationCidrBlock: aws.String("0.0.0.0/0"),
			GatewayId: internetGateway.InternetGatewayId,
		})

	} else {
		_, err = svc.CreateRoute(&ec2.CreateRouteInput{
			RouteTableId: routeTable.RouteTable.RouteTableId,
			DestinationCidrBlock: aws.String("0.0.0.0/0"),
			NatGatewayId: natGateway.NatGatewayId,
		})
	}

	// Associate with the subnet

	_, _ = svc.AssociateRouteTable(&ec2.AssociateRouteTableInput{
		SubnetId: subnet.SubnetId,
		RouteTableId: routeTable.RouteTable.RouteTableId,
	})

	return routeTable.RouteTable, nil
}
