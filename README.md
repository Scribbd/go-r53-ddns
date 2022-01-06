## Usage

[Helm](https://helm.sh) must be installed to use the charts.  Please refer to
Helm's [documentation](https://helm.sh/docs) to get started.

Once Helm has been set up correctly, add the repo as follows:

  helm repo add go-r53-ddns https://scribbd.github.io/go-r53-ddns

If you had already added this repo earlier, run `helm repo update` to retrieve
the latest versions of the packages.  You can then run `helm search repo
go-r53-ddns` to see the charts.

To install the go-r53-ddns chart:

    helm install my-go-r53-ddns go-r53-ddns/go-r53-ddns

To uninstall the chart:

    helm delete my-go-r53-ddns