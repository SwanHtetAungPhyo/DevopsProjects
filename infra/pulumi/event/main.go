package main

import (
	"event-driven/modules"
	"os"

	"github.com/joho/godotenv"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"github.com/rs/zerolog/log"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal().
			Err(err).
			Msg("Error loading .env file")
	}
	accessKeyId := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	userData, err := os.ReadFile("../../scripts/placeholder.sh")
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Error opening placeholder script file")
	}

	pulumi.Run(func(ctx *pulumi.Context) error {

		awsCfg := config.New(ctx, "aws")
		awsRegion := awsCfg.Require("region")
		awsProfile := awsCfg.Require("profile")

		log.Info().
			Str("region", awsRegion).
			Str("profile", awsProfile).
			Msg("Using AWS configuration")

		provider, err := aws.NewProvider(ctx, "production", &aws.ProviderArgs{
			AccessKey: pulumi.String(accessKeyId),
			SecretKey: pulumi.String(secretAccessKey),
			Profile:   pulumi.String(awsProfile),
			Region:    pulumi.String(awsRegion),
		})
		if err != nil {
			log.Fatal().Err(err).Msg("Error creating provider")
			return err
		}

		network, err := modules.NewVPC(ctx, "production-vpc", modules.VPCArgs{
			CidrBlock: "10.0.0.0/16",
			PublicSubnetCidr: []string{
				"10.0.1.0/24",
				"10.0.2.0/24",
			},
			AvailabilityZone: []string{
				"us-east-1a",
				"us-east-1b",
			},
			Tags: map[string]string{
				"Environment": "dev",
				"Project":     "demo",
			},
		}, provider)
		if err != nil {
			return err
		}
		instance, err := modules.NewEC2(ctx, "web-server", modules.EC2Args{
			InstanceType: "t2.micro",
			VpcId:        network.VpcId,
			SubnetId:     network.PublicSubnetId.Index(pulumi.Int(0)),
			KeyName:      "my-key-pair",
			UserData:     string(userData),
			Tags: map[string]string{
				"Name":        "web-server",
				"Environment": "dev",
			},
		}, provider)
		if err != nil {
			return err
		}

		ctx.Export("vpcId", network.VpcId)
		ctx.Export("publicSubnetIds", network.PublicSubnetId)
		ctx.Export("instanceId", instance.InstanceId)
		ctx.Export("instancePublicIp", instance.PublicIp)
		ctx.Export("instancePublicDns", instance.PublicDns)

		return nil
	})

}
