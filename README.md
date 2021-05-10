# VSCode extension for RKE Cluster Configuration files

[![Visual Studio Marketplace](https://img.shields.io/visual-studio-marketplace/v/dancermak.vscode-rke-cluster-config)](https://marketplace.visualstudio.com/items?itemName=dancermak.vscode-rke-cluster-config)
[![extension build](https://img.shields.io/github/workflow/status/dcermak/vscode-rke-cluster-config/extension?label=extension%20build)](https://github.com/dcermak/vscode-rke-cluster-config/actions/workflows/extension.yml)
[![schema check](https://img.shields.io/github/workflow/status/dcermak/vscode-rke-cluster-config/schema?label=schema%20check)](https://github.com/dcermak/vscode-rke-cluster-config/actions/workflows/schema.yml)
[![go tests](https://img.shields.io/github/workflow/status/dcermak/vscode-rke-cluster-config/go?label=go%20test)](https://github.com/dcermak/vscode-rke-cluster-config/actions/workflows/go.yml)
[![License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://github.com/dcermak/vscode-rke-cluster-config/blob/main/LICENSE)
[![Dependabot Status](https://api.dependabot.com/badges/status?host=github&repo=dcermak/vscode-rke-cluster-config)](https://dependabot.com)
[![styled with prettier](https://img.shields.io/badge/styled_with-prettier-ff69b4.svg)](https://github.com/prettier/prettier)

This extension adds support for linting & validation of the `cluster.yml`
configuration file for RKE.

It leverages the
[vscode-yaml](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml)
extension by adding the schema of the `cluster.yml` files to the extension's
configuration.


## Features

- autocompletion of configuration keys and values:

![autocomplete](media/autocomplete.gif)

- validation of the keys and values:

![linting_error](media/linting_error.png)

- documentation of the configuration options on hover:

![doc_on_hover](media/doc_on_hover.png)

- breadcrumbs (provided by the YAML extension):

![breadcrumbs](media/breadcrumbs.png)


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
