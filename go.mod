module github.com/dcermak/cluster-config-dumper

go 1.15

replace (
	github.com/alecthomas/jsonschema => github.com/alecthomas/jsonschema v0.0.0-20210526225647-edb03dcab7bc
	github.com/rancher/rke => github.com/dcermak/rke v1.3.0-rc1.0.20210824145743-7c3e5260153a
	k8s.io/client-go => k8s.io/client-go v0.21.0
)

require (
	github.com/alecthomas/jsonschema v0.0.1
	github.com/fatih/structtag v1.2.0
	github.com/rancher/rke v1.2.9
	k8s.io/apimachinery v0.23.3
)
