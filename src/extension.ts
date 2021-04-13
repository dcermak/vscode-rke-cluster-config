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

import * as vscode from "vscode";

const schemaUrl =
  "https://raw.githubusercontent.com/dcermak/vscode-rke-cluster-config/main/cluster.yml.json";

export async function activate(): Promise<void> {
  //  context: vscode.ExtensionContext
  const yamlConfSection = vscode.workspace.getConfiguration("yaml");
  const yamlSchemas = yamlConfSection.get<Record<string, string | string[]>>(
    "schemas",
    {}
  );

  if (yamlSchemas[schemaUrl] === undefined) {
    yamlSchemas[schemaUrl] = ["cluster.yml", "cluster.yaml"];
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
