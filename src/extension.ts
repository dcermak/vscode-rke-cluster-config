/**
 * Copyright (c) 2021 SUSE LLC
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

import { resolve } from "path";
import * as vscode from "vscode";

const SCHEMA_URL =
  "https://raw.githubusercontent.com/dcermak/vscode-rke-cluster-config/main/schemas/cluster.yml.json";

const CLUSTER_YAML_FILES = ["cluster.yml", "cluster.yaml"];

export async function activate(): Promise<void> {
  //  context: vscode.ExtensionContext
  const yamlConfSection = vscode.workspace.getConfiguration("yaml");
  const yamlSchemas = yamlConfSection.get<Record<string, string | string[]>>(
    "schemas",
    {}
  );

  const debug = process.env.EXTENSION_DEBUG === "1";

  const key = debug
    ? resolve(__dirname, "..", "schemas", "cluster.yml.json")
    : SCHEMA_URL;

  if (yamlSchemas[key] === undefined) {
    if (debug) {
      delete yamlSchemas[SCHEMA_URL];
    }
    yamlSchemas[key] = CLUSTER_YAML_FILES;

    await yamlConfSection.update(
      "schemas",
      yamlSchemas,
      vscode.ConfigurationTarget.Global
    );
  }
}

export function deactivate(): void {
  // NOP
}
