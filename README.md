# GoR53DDNS
An AWS Route53 DDNS updater program run in a K8s cronjob written in golang.

Uses https://checkip.amazonaws.com to check for IP. But this is configurable

## Goals
This is a simple exercize in creating a full build from golang to docker build to helm.

## Usage
> THIS CODE IS NOT PRODUCTION READY. 
> This is my first time using golang. And this is meant as personal project for my personal cluster. There are better solutions for DDNS out there.

You can use this program in three ways:
- As CLI-tool
- In Docker
- With Helm

### CLI
### Docker

### Helm
## Configuration
Configuration is done through environment variables. Make certain you inject the variables in a secure way.
AWS module specific:
`AWS_ACCESS_KEY_ID`
`AWS_SECRET_ACCESS_KEY`
`AWS_SESSION_TOKEN`

IP retrieval endpoint:
`IP_API_SOURCE` API endpoint for retrieving IP. Has to return plain text value. Default: https://checkip.amazonaws.com

Cluster domain:
`CLUSTER_DOMAIN` The domain name used in your hosted zone. Example: `cluster.example.com.`

Hosted Zone ID:
`HOSTED_ZONE_ID` The hosted zone ID at AWS.