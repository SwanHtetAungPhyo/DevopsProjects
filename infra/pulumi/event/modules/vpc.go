package modules

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/rs/zerolog/log"
)

type (
	VPCArgs struct {
		CidrBlock        string            `pulumi:"cidrBlock"`
		PublicSubnetCidr []string          `pulumi:"publicSubnetCider"`
		AvailabilityZone []string          `pulumi:"availabilityZone"`
		Tags             map[string]string `pulumi:"tags"`
	}
	VPC struct {
		VpcId           pulumi.IDOutput          `pulumi:"vpcId"`
		PublicSubnetId  pulumi.StringArrayOutput `pulumi:"publicSubnetId"`
		InternetGateway pulumi.IDOutput          `pulumi:"internetGateway"`
		RouteTable      pulumi.IDOutput          `pulumi:"routeTable"`
	}
)

func NewVPC(
	ctx *pulumi.Context,
	name string,
	args VPCArgs,
	provider *aws.Provider,
) (*VPC, error) {
	vpc, err := ec2.NewVpc(ctx, name, &ec2.VpcArgs{
		CidrBlock:          pulumi.String(args.CidrBlock),
		EnableDnsHostnames: pulumi.BoolPtr(true),
		EnableDnsSupport:   pulumi.BoolPtr(true),
		Tags:               pulumi.ToStringMap(args.Tags),
	}, pulumi.Provider(provider))

	if err != nil {
		log.Err(err).Msg("Error creating vpc")
		return nil, err
	}

	igw, err := ec2.NewInternetGateway(ctx, fmt.Sprintf("%s-igw", name), &ec2.InternetGatewayArgs{
		VpcId: vpc.ID(),
		Tags:  pulumi.ToStringMap(args.Tags),
	}, pulumi.Provider(provider))
	if err != nil {
		log.Err(err).Msg("Error creating internet gateway")
		return nil, err
	}

	routeTable, err := ec2.NewRouteTable(ctx, fmt.Sprintf("%s-rt", name), &ec2.RouteTableArgs{
		VpcId: vpc.ID(),
		Routes: ec2.RouteTableRouteArray{
			ec2.RouteTableRouteArgs{
				CidrBlock: pulumi.String("0.0.0.0/0"),
				GatewayId: igw.ID(),
			},
		},
		Tags: pulumi.ToStringMap(args.Tags),
	}, pulumi.Provider(provider))
	if err != nil {
		log.Err(err).Msg("Error creating route table")
		return nil, err
	}

	var subnetIds pulumi.StringArrayOutput
	var subnetOutputs []pulumi.StringOutput

	for i, cidr := range args.PublicSubnetCidr {
		subnet, err := ec2.NewSubnet(ctx, fmt.Sprintf("%s-subnet-%d", name, i), &ec2.SubnetArgs{
			VpcId:               vpc.ID(),
			CidrBlock:           pulumi.String(cidr),
			AvailabilityZone:    pulumi.String(args.AvailabilityZone[i]),
			MapPublicIpOnLaunch: pulumi.BoolPtr(true),
			Tags:                pulumi.ToStringMap(args.Tags),
		}, pulumi.Provider(provider))
		if err != nil {
			log.Err(err).Msg("Error creating subnet")
			return nil, err
		}

		_, err = ec2.NewRouteTableAssociation(ctx, fmt.Sprintf("%s-rta-%d", name, i), &ec2.RouteTableAssociationArgs{
			SubnetId:     subnet.ID(),
			RouteTableId: routeTable.ID(),
		}, pulumi.Provider(provider))
		if err != nil {
			log.Err(err).Msg("Error associating subnet")
			return nil, err
		}

		subnetOutputs = append(subnetOutputs, subnet.ID().ToStringOutput())
	}

	subnetIds = pulumi.ToStringArrayOutput(subnetOutputs)

	return &VPC{
		VpcId:           vpc.ID(),
		PublicSubnetId:  subnetIds,
		InternetGateway: igw.ID(),
		RouteTable:      routeTable.ID(),
	}, nil
}
