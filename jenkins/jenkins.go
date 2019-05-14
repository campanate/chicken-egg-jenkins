package jenkins

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)



func CreateJenkins() (error) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := ec2.New(sess)

	// Get Zones

	zones, err := GetZones(svc)

	// Create VPC
	vpc, err := CreateVpc(svc)

	// Create Internet Gateway

	internetGateway, err := CreateInternetGateway(svc, vpc)

	// Create Elastic IP

	eip, err := CreateElasticIp(svc, vpc)

	// Create Subnets

	privateSubnet, err := CreateSubnet(svc, "7.0.1.0/24", false, zones[0], vpc)
	publicSubnet, err := CreateSubnet(svc, "7.0.3.0/24", true, zones[1], vpc)

	// Create NatGateway

	natGateway, err := CreateNatGateway(svc, eip, privateSubnet)

	// Create RouteTable

	_, err = CreateRouteTable(svc, vpc, false, privateSubnet, natGateway, nil)
	_, err = CreateRouteTable(svc, vpc, true, publicSubnet, nil, internetGateway)

	// Create Security Groups

	bastionSecurityGroup, err := CreateSecurityGroup(svc, vpc, "jenkins-bastion-security-group", "Enable access through bastion", "0.0.0.0/0", "TCP", 22, 22)
	accessSecurityGroup, err := CreateSecurityGroup(svc, vpc, "jenkins-access-security-group", "Enable access to all VPC protocols and IPs", "7.0.0.0/16", "-1", 22, 22)
	elbSecurityGroup, err := CreateSecurityGroup(svc, vpc, "jenkins-elb-security-group", "Enable access to all IPS to port 80 and 443", "7.0.0.0/16", "-1", 80, 80)

	// Create Load Balancers

	loadBalancerSubnets := []*string{publicSubnet.SubnetId}
	loadBalancerSecurityGroups := []*string{elbSecurityGroup.GroupId}
	zone_names := []*string{zones[0].ZoneName}


	_, _ = CreateELB(loadBalancerSubnets, loadBalancerSecurityGroups, zone_names)

	// Create Instances

	bastionSecurityGroups := []*string{accessSecurityGroup.GroupId, bastionSecurityGroup.GroupId}

	return err
}
