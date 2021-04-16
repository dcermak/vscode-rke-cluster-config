# VSCode extension for RKE Cluster Configuration files

[![extension build](https://img.shields.io/github/workflow/status/dcermak/vscode-rke-cluster-config/extension?label=extension%20build)](https://github.com/dcermak/vscode-rke-cluster-config/actions/workflows/extension.yml)
[![schema check](https://img.shields.io/github/workflow/status/dcermak/vscode-rke-cluster-config/schema?label=schema%20check)](https://github.com/dcermak/vscode-rke-cluster-config/actions/workflows/schema.yml)
[![go tests](https://img.shields.io/github/workflow/status/dcermak/vscode-rke-cluster-config/go?label=go%20test)](https://github.com/dcermak/vscode-rke-cluster-config/actions/workflows/go.yml)

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


The schema is generated using the
[alecthomas/jsonschema](https://github.com/alecthomas/jsonschema) module from
the structs in
[rancher/rke/blob/master/types/rke_types.go](https://github.com/rancher/rke/blob/master/types/rke_types.go). This
is performed by `dump_cluster_config_schema.go` with some additional post
processing steps performed in `process-schema` (the `$schema` key is removed,
because the validator fails to handle it and the extracted documentation is
added as the description field) via `yarn run schema`.
