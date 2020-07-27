## grafana controller

Simple controller to interact with grafana.

Controller allows you to save your grafana dashboards as a CRD, which can then be managed by the controller.

The controller uses the grafana-tools/sdk for interacting with the grafana endpoint defined in the configuration.

The controller supports the following types:

* dashboards
* folders
* datasources

Sample CRD's are available in `samples` directory.

_Note:_ When the crd resources are deleted, the controller will clean them up from Grafana too.


A helm chart to consume the same is available in the `charts` directory.

The chart needs two values to be updated based on your environment.

```cassandraql
args:
  grafana_endpoint: "http://access-grafana.cattle-prometheus.svc.cluster.local"
  grafana_token: "admin:password"
```
