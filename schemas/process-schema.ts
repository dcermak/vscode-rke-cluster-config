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

import { promises as fsPromises } from "fs";
import { resolve } from "path";

function removeSchemaEntries(schema: any): void {
  delete schema["$schema"];
  Object.keys(schema).forEach((key) => {
    if (typeof schema[key] === "object" && schema[key] !== null) {
      removeSchemaEntries(schema[key]);
    }
  });
}

type DocumentationMap = Record<string, Record<string, string>>;
type StructFieldNameMap = Record<string, Record<string, string>>;

function addDescriptionToDefinitions(
  schema: any,
  docMap: DocumentationMap,
  fieldNames: StructFieldNameMap
) {
  const definitions = schema["definitions"];
  if (definitions === undefined) {
    throw new Error("Invalid schema, key 'definitions' not present");
  }

  Object.keys(definitions).forEach((key) => {
    const docs = docMap[key];
    // console.log(`key=${key}, docs=${inspect(docs, { depth: null })}`);
    if (docs === undefined) {
      return;
    }
    const newFieldName = (fieldName: string): string =>
      (fieldNames[key] ?? {})[fieldName] ?? fieldName;

    const def = definitions[key];
    const props = def["properties"];
    // console.log(`props=${inspect(props, { depth: null })}`);
    if (props === undefined) {
      return;
    }

    Object.keys(props).forEach((propKey) => {
      props[propKey]["description"] = docs[newFieldName(propKey)]?.trim();
    });
  });
}

async function removeSchemaFromJson(
  schemaPath: string,
  docMapPath: string,
  fieldNamePath: string
): Promise<void> {
  const [schema, docMap, fieldNameMap] = await Promise.all(
    [schemaPath, docMapPath, fieldNamePath].map(async (path) =>
      JSON.parse(await fsPromises.readFile(path, { encoding: "utf-8" }))
    )
  );
  Object.keys(schema).forEach((key) => {
    if (typeof schema[key] === "object" && schema[key] !== null) {
      removeSchemaEntries(schema[key]);
    }
  });

  addDescriptionToDefinitions(schema, docMap, fieldNameMap);

  await fsPromises.writeFile(schemaPath, JSON.stringify(schema, undefined, 2));
}

const resolvePath = (path: string) => resolve(__dirname, "..", path);

Promise.all(
  [
    ["cluster.json", "jsonNames.json"],
    ["cluster.yml.json", "yamlNames.json"],
  ].map(async ([schemaName, nameMapJson]) =>
    removeSchemaFromJson(
      resolvePath(schemaName),
      resolvePath("docMap.json"),
      resolvePath(nameMapJson)
    )
  )
)
  .then()
  .catch((err) => {
    console.error(err);
    process.exitCode = 1;
  });
