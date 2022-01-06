# GoR53DDNS
An AWS Route53 DDNS updater program run in a K8s cronjob written in Golang.

Uses https://checkip.amazonaws.com to check for IP. But this is configurable

Check [My Goals](./GOALS.md) to see my learning journey.

## Usage
> THIS CODE IS NOT PRODUCTION READY. 
> This is my first time using Golang. And this is meant as personal project for my personal cluster. There are better solutions for DDNS out there.

You can use this program in three ways:
- As CLI-tool
- In Docker
- With Helm

### CLI
Download your build from the releases or build your own package from source.
Either use the command line flags or set the environment variables. Consult [#Configuration] for the values/flags.

### Docker
Get your image through either Docker Hub `docker pull scribbd/go-r53-ddns:latest` or through ghcr.io `docker pull ghcr.io/scribbd/go-r53-ddns:latest`.

Make certain you inject the right environment variables listed below. Consult [#Configuration] for the values/flags.

### Helm
No package is available yet to add through `helm repo add`. Installation is by cloning this repository, go into `./helm/go-r53-dns` and run `helm install go-r53-ddns . -n go-r53-ddns --create-namespace --atomic`.

If your jobs quit with `Exit Code: 2` make certain you deployed a Secret, and have your environment variables set.

## Configuration
Configuration is done through environment variables. Make certain you inject the variables in a secure way.

| Environment Key         | Command Flag | Req | Description                                    | Default                         |
|-------------------------|--------------|-----|------------------------------------------------|---------------------------------|
| `AWS_ACCESS_KEY_ID`     | -a           | Yes | AWS default environment variable               |                                 |
| `AWS_SECRET_ACCESS_KEY` | -s           | Yes | AWS default environment variable               |                                 |
| `AWS_SESSION_TOKEN`     | -t           | No  | AWS default environment variable               |                                 |
| `IP_API_SOURCE`         | -i           | No  | API endpoint URL for getting your remote IP    | `https://checkip.amazonaws.com` |
| `CLUSTER_DOMAIN`        | -c           | Yes | Domain used in cluster: `cluster.example.com.` |                                 |
| `HOSTED_ZONE_ID`        | -h           | Yes | Zone ID from Route53                           |                                 |

### AWS-IAM

Make certain you have AWS configured to accept requests from this program. A CloudFormation Template is available in the [CloudFormation folder](./cloudformation/) that deploys an IAM-user with a customer managed permissions policy.