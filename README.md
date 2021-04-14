# VSCode extension for RKE Cluster Configuration files

This extension adds support for linting & validation of the `cluster.yml`
configuration file for RKE.

It leverages the
[vscode-yaml](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml)
extension by adding the schema of the `cluster.yml` files to the extension's
configuration.


## Development

Requirements:

- go
- nodejs and yarn
- [jq](https://stedolan.github.io/jq/)


The schema is generated using the
[alecthomas/jsonschema](https://github.com/alecthomas/jsonschema) module from
the structs in
[rancher/rke/blob/master/types/rke_types.go](https://github.com/rancher/rke/blob/master/types/rke_types.go). This
is performed by `dump_cluster_config_schema.go` with some additional post
processing steps performed using `jq` (the `$schema` key is removed, because the
validator fails to handle it) via `yarn run schema`.
