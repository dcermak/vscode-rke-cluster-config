module github.com/dcermak/cluster-config-dumper

go 1.15

replace (
	k8s.io/client-go => k8s.io/client-go v0.20.0
	github.com/alecthomas/jsonschema => github.com/alecthomas/jsonschema v0.0.0-20210413112511-5c9c23bdc720
	github.com/rancher/rke => github.com/dcermak/rke v1.3.0-rc1.0.20210413132551-411c1ef71a17
)

require (
	github.com/alecthomas/jsonschema v0.0.1
	github.com/rancher/rke v1.2.7
)
