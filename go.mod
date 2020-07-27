module github.com/ibrokethecloud/grafana-controller

go 1.13

require (
	github.com/go-logr/logr v0.1.0
	github.com/grafana-tools/sdk v0.0.0-20200713194907-007f486b53df
	github.com/json-iterator/go v1.1.10
	github.com/namsral/flag v1.7.4-pre
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	gonum.org/v1/netlib v0.0.0-20190331212654-76723241ea4e // indirect
	k8s.io/apimachinery v0.18.4
	k8s.io/client-go v0.18.4
	k8s.io/klog v1.0.0
	sigs.k8s.io/controller-runtime v0.6.1
	sigs.k8s.io/structured-merge-diff v1.0.1-0.20191108220359-b1b620dd3f06 // indirect
)
