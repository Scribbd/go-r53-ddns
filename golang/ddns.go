package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	r53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
)

const (
	API_SOURCE_ENV_KEY = "IP_API_SOURCE"
	API_SOURCE_DEFAULT = "https://checkip.amazonaws.com"

	CLUSTER_DOMAIN_ENV_KEY = "CLUSTER_DOMAIN"
	HOSTED_ZONE_ID_ENV_KEY = "HOSTED_ZONE_ID"

	MAX_REQ_ITEMS = 10
)

// Environment Support Functions
func getEnvWithFallback(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvWithoutFallback(key string) string {
	answer := getEnvWithFallback(key, "")
	if len(answer) == 0 {
		log.Panicf("%s not found", key)
	}
	return answer
}

// MAIN FUNCTION
func main() {
	log.Println("Hello! R53-DDNS waking up to check AWS")

	// Check environment variables that are used later
	zoneId := getEnvWithoutFallback(HOSTED_ZONE_ID_ENV_KEY)
	clusterDomain := getEnvWithoutFallback(CLUSTER_DOMAIN_ENV_KEY)

	// Get current IP from external source
	// inspired from https://www.golangprograms.com/how-to-get-current-ip-form-ipify-org.html
	var extIp string
	apiUrl := getEnvWithFallback(API_SOURCE_ENV_KEY, API_SOURCE_DEFAULT)
	{
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
	{
		inCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
		cfg = inCfg
		if err != nil {
			log.Panicf("Configuration error:\n%s", err.Error())
		}
	}

	// Create Route53 Client
	r53 := route53.NewFromConfig(cfg)

	// Get current IP from hosted-zone
	// I could just do a DNS query... But I am too far in...
	var resourceIp string
	var ttl int64
	{
		maxHelper := int32(MAX_REQ_ITEMS)
		req := route53.ListResourceRecordSetsInput{
			HostedZoneId: aws.String(zoneId),
			MaxItems:     &maxHelper,

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
