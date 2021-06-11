module github.com/dcermak/cluster-config-dumper

go 1.15

replace (
	github.com/alecthomas/jsonschema => github.com/alecthomas/jsonschema v0.0.0-20210526225647-edb03dcab7bc
	github.com/rancher/rke => github.com/dcermak/rke v1.3.0-rc1.0.20210607085837-a89f836f7e4e
	k8s.io/client-go => k8s.io/client-go v0.21.0
)

require (
	github.com/alecthomas/jsonschema v0.0.1
	github.com/fatih/structtag v1.2.0
	github.com/rancher/rke v1.2.8
	k8s.io/apimachinery v0.21.1
)
