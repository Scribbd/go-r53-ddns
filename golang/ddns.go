package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	r53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
)

const (
	API_SOURCE_FLAG    = "i"
	API_SOURCE_ENV_KEY = "IP_API_SOURCE"
	API_SOURCE_DEFAULT = "https://checkip.amazonaws.com"
	API_SOURCE_USAGE   = "Define an api-endpoint that provides an IP in plain text. Example: 'https://getip.example.com'"

	CLUSTER_DOMAIN_FLAG    = "c"
	CLUSTER_DOMAIN_ENV_KEY = "CLUSTER_DOMAIN"
	CLUSTER_DOMAIN_USAGE   = "Provide the domain you wish to update your IP to. Include the trailing '.'. Example: 'cluster.exmaple.com.'"

	HOSTED_ZONE_ID_FLAG    = "h"
	HOSTED_ZONE_ID_ENV_KEY = "HOSTED_ZONE_ID"
	HOSTED_ZONE_ID_USAGE   = "Provide a ID of your Hosted Zone in Route53 you wish to update an entry in."

	AWS_SECRET_ACCESS_KEY_FLAG = "s"
	AWS_SECRET_KEY_ENV_KEY     = "AWS_SECRET_KEY"
	AWS_SECRET_KEY_USAGE       = "Provide a secret access key provided by IAM"

	AWS_ACCESS_KEY_FLAG    = "a"
	AWS_ACCESS_KEY_ENV_KEY = "AWS_ACCESS_KEY"
	AWS_ACCESS_KEY_USAGE   = "Provide the access key ID associated with the secret"

	AWS_SESSION_TOKEN_FLAG    = "t"
	AWS_SESSION_TOKEN_ENV_KEY = "AWS_SESSION_TOKEN"
	AWS_SESSION_TOKEN_USAGE   = "Should this be applicable provide the session token."

	MAX_REQ_ITEMS = 10

	HELP_MESSAGE = "\nTo use this executable make certain the following environment variables are set: " + CLUSTER_DOMAIN_ENV_KEY + ", " + HOSTED_ZONE_ID_ENV_KEY + ", " + AWS_SECRET_KEY_ENV_KEY + ", " + AWS_ACCESS_KEY_ENV_KEY
)

var (
	cli_api     string
	cli_cluster string
	cli_zone    string
	cli_secret  string
	cli_access  string
	cli_token   string
)

func getVarInput(envKey string, flagValue string, fallback *string) string {
	// variables provided through cli flag takes priority
	if flagValue != "" { //String null value is "" set as default
		return flagValue
	} else if value, ok := os.LookupEnv(envKey); ok {
		return value
	}
	// Panic when no value is found and no fallback value was presented
	if fallback == nil {
		log.Panicf("%s not found.%s", envKey, HELP_MESSAGE)
	}
	return *fallback
}

func init() {
	// Setup flag parser (String null value is "")
	flag.StringVar(&cli_api, API_SOURCE_FLAG, "", API_SOURCE_USAGE)
	flag.StringVar(&cli_cluster, CLUSTER_DOMAIN_FLAG, "", CLUSTER_DOMAIN_USAGE)
	flag.StringVar(&cli_zone, HOSTED_ZONE_ID_FLAG, "", HOSTED_ZONE_ID_USAGE)
	flag.StringVar(&cli_secret, AWS_SECRET_ACCESS_KEY_FLAG, "", AWS_ACCESS_KEY_USAGE)
	flag.StringVar(&cli_access, AWS_ACCESS_KEY_FLAG, "", AWS_ACCESS_KEY_USAGE)
	flag.StringVar(&cli_token, AWS_SESSION_TOKEN_FLAG, "", AWS_SESSION_TOKEN_USAGE)
}

// MAIN FUNCTION
func main() {
	log.Println("Hello! R53-DDNS waking up to check AWS")
	flag.Parse()

	// Check environment variables that are used later
	zoneId := getVarInput(HOSTED_ZONE_ID_ENV_KEY, cli_zone, nil)
	clusterDomain := getVarInput(CLUSTER_DOMAIN_ENV_KEY, cli_cluster, nil)

	// Get current IP from external source
	// inspired from https://www.golangprograms.com/how-to-get-current-ip-form-ipify-org.html
	var extIp, apiUrl string
	{ // get API_URL
		var defaultV string = API_SOURCE_DEFAULT
		apiUrl = getVarInput(API_SOURCE_ENV_KEY, cli_api, &defaultV)
	}
	{ // Get external IP
		resp, err := http.Get(apiUrl)
		if err != nil {
			log.Panic(err)
		}
		defer resp.Body.Close()
		rawExtIp, err := ioutil.ReadAll(resp.Body)
		// Extract added newline character
		extIp = strings.TrimSuffix(string(rawExtIp), "\n")
		if err != nil {
			log.Panic(err)
		}
	}

	log.Printf("Found IP with %s", apiUrl)

	log.Println("Get current IP in records")
	// Load Configuration from
	// AWS_SECRET_ACCESS_KEY or AWS_SECRET_KEY
	// AWS_ACCESS_KEY_ID or AWS_ACCESS_KEY
	var cfg aws.Config
	{ // Create Configuration
		var err error
		if cli_secret == "" {
			cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
		} else {
			cfg, err = config.LoadDefaultConfig(context.TODO(),
				config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cli_access, cli_secret, cli_token)),
				config.WithRegion("us-west-2"))
		}
		if err != nil {
			log.Panicf("Configuration error:\n%s%s", err.Error(), HELP_MESSAGE)
		}
	}

	// Create Route53 Client
	r53 := route53.NewFromConfig(cfg)

	// Get current IP from hosted-zone
	// I could just do a DNS query... But I am too far in...
	var resourceIp string
	var ttl int64
	{ // Get current active IP from Route53
		maxHelper := int32(MAX_REQ_ITEMS)
		req := route53.ListResourceRecordSetsInput{
			HostedZoneId:    aws.String(zoneId),
			MaxItems:        &maxHelper,
			StartRecordName: aws.String(clusterDomain),
		}
		resp, err := r53.ListResourceRecordSets(context.TODO(), &req)
		if err != nil {
			log.Panicf("Hosted zone (%s) not found", zoneId)
		} else {
			// Iterate through all records to see if exact match is found
			for _, element := range resp.ResourceRecordSets {
				if *element.Name == clusterDomain {
					resourceIp = *element.ResourceRecords[0].Value
					ttl = *element.TTL
					break
				}
			}
			// Panic when no result has been found
			if resourceIp == "" {
				log.Panicf("Domain (%s) not found in Hosted Zone (%s)", clusterDomain, zoneId)
			}
		}
	}

	log.Printf("Got the resource from hosted zone: %s", zoneId)

	if resourceIp == extIp {
		log.Print("IPs match. Nothing to do, shutting down!")
	} else {
		log.Print("IPs don't match! Updating hosted zone now!")
		// Constructing request for changing IP
		req := &route53.ChangeResourceRecordSetsInput{
			HostedZoneId: &zoneId,
			ChangeBatch: &r53types.ChangeBatch{
				Changes: []r53types.Change{
					{
						Action: r53types.ChangeActionUpsert,
						ResourceRecordSet: &r53types.ResourceRecordSet{
							Name: aws.String(clusterDomain),
							ResourceRecords: []r53types.ResourceRecord{
								{
									Value: aws.String(extIp),
								},
							},
							Type: r53types.RRTypeA,
							TTL:  &ttl,
						},
					},
				},
			},
		}
		// Commit change
		resp, err := r53.ChangeResourceRecordSets(context.TODO(), req)
		if err != nil {
			log.Panicf("Something went wrong while updating %s! No changes have been made!\n%s", clusterDomain, err)
		}
		log.Printf("Change is in progress. ID returned: %s", *resp.ChangeInfo.Id)
	}
}
