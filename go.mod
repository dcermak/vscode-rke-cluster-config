module github.com/dcermak/cluster-config-dumper

go 1.15

replace (
	github.com/alecthomas/jsonschema => github.com/alecthomas/jsonschema v0.0.0-20210413112511-5c9c23bdc720
	github.com/rancher/rke => github.com/dcermak/rke v1.3.0-rc1.0.20210413153447-1dd78c94608a
	k8s.io/client-go => k8s.io/client-go v0.20.0
)

require (
	github.com/alecthomas/jsonschema v0.0.1
	github.com/fatih/structtag v1.2.0
	github.com/rancher/rke v1.2.8
	k8s.io/apimachinery v0.21.0
)
