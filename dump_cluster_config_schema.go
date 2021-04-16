// Copyright (c) 2021 SUSE LLC
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"reflect"
	"regexp"

	"github.com/alecthomas/jsonschema"
	"github.com/fatih/structtag"
	conf "github.com/rancher/rke/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var byteSliceType = reflect.TypeOf([]byte(nil))

// typeAdjust changes the json schema of the following types used in conf.RancherKubernetesEngineConfig:
// - NumberStringMap is converted a json schema type that supports not only strings as values but also numbers
// - metav1.Time is changed to a date-time string
// - byte slices are change to ordinary strings (the schema validator does not support media: {} entries)
func typeAdjust(i reflect.Type) *jsonschema.Type {
	if i == reflect.TypeOf(conf.NumberStringMap{}) {
		return &jsonschema.Type{
			Type:              "object",
			PatternProperties: map[string]*jsonschema.Type{".*": {OneOf: []*jsonschema.Type{{Type: "string"}, {Type: "number"}}}},
		}
	} else if i == reflect.TypeOf(&metav1.Time{}) {
		return &jsonschema.Type{Type: "string", Format: "date-time"}
	} else if i.Kind() == reflect.Slice && i.Elem() == byteSliceType.Elem() {
		return &jsonschema.Type{Type: "string"}
	}
	return nil
}

var re = regexp.MustCompile(`^\s*$`)

func isWhitespace(s string) bool {
	return re.MatchString(s)
}

// ExtractDoc processes the given golang source file contents in src and
// extracts the documentation of all struct fields into a map of maps.
// The outer map's key is the name of all structs in that file, whereas the
// inner map is a `fieldName` => `documentation` mapping.
// This function panics if an error occurs.
func ExtractDoc(src string) map[string](map[string]string) {
	structMap := map[string](map[string]string){}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	lastType := ""

	ast.Inspect(f, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.TypeSpec:
			lastType = t.Name.Name
			structMap[lastType] = map[string]string{}

		case *ast.StructType:

			fieldMap := structMap[lastType]

			for _, field := range t.Fields.List {
				docIsWhitespace := isWhitespace(field.Doc.Text())
				commentIsWhitespace := isWhitespace(field.Comment.Text())
				if (!docIsWhitespace || !commentIsWhitespace) && len(field.Names) == 1 {
					if !docIsWhitespace {
						fieldMap[field.Names[0].Name] = field.Doc.Text()
					} else if !commentIsWhitespace {
						fieldMap[field.Names[0].Name] = field.Comment.Text()
					}
				}
			}

			structMap[lastType] = fieldMap
		}
		return true
	})

	for k, v := range structMap {
		if len(v) == 0 {
			delete(structMap, k)
		}
	}

	return structMap
}

func getAllFields(t reflect.Type) []reflect.StructField {
	fields := make([]reflect.StructField, 0, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Type.Kind() == reflect.Struct && f.Anonymous {
			fields = append(fields, getAllFields(f.Type)...)
		} else {
			fields = append(fields, f)
		}
	}

	return fields
}

// GenerateStructFieldNameMap takes the provided type t and writes a mapping
// between the json/yaml keys and the actual struct member fields into nameMap
//
// nameMap will at the end contain a map for each struct present in t where the
// key is the struct's name and the value is a map with `jsonName => fieldName`
//
// When preferYamlTag is true, then the name provided in `yaml:` tags is
// preferred over json tags, else it is the other way around. If neither a json
// nor a yaml tag is present, then the fields name is used.
func GenerateStructFieldNameMap(nameMap *map[string](map[string]string), t reflect.Type, preferYamlTag bool) {
	if _, present := (*nameMap)[t.Name()]; !present {
		docMap := make(map[string]string, t.NumField())
		(*nameMap)[t.Name()] = docMap
	}

	docMap := (*nameMap)[t.Name()]

	for _, f := range getAllFields(t) {
		t, err := structtag.Parse(string(f.Tag))
		if err != nil {
			docMap[f.Name] = f.Name
			continue
		}
		jsonEntry, _ := t.Get("json")
		yamlEntry, _ := t.Get("yaml")

		name := f.Name

		if yamlEntry != nil {
			name = yamlEntry.Name
		}
		if (jsonEntry != nil && !preferYamlTag) || (yamlEntry == nil && jsonEntry != nil) {
			name = jsonEntry.Name
		}

		if name != "" {
			docMap[name] = f.Name
		}

		if f.Type.Kind() == reflect.Struct {
			GenerateStructFieldNameMap(nameMap, f.Type, preferYamlTag)
		} else if (f.Type.Kind() == reflect.Ptr || f.Type.Kind() == reflect.Array || f.Type.Kind() == reflect.Slice) && f.Type.Elem().Kind() == reflect.Struct {
			GenerateStructFieldNameMap(nameMap, f.Type.Elem(), preferYamlTag)
		}
	}
}

// dumpAsJsonToFile marshals the provided interface v to json and writes the
// contents into the file with the provided filename
// The function panics when an error occurs.
func dumpAsJsonToFile(filename string, v interface{}) {
	jD, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filename, jD, 0644)
	if err != nil {
		panic(err)
	}
}

func main() {
	src, err := ioutil.ReadFile("vendor/github.com/rancher/rke/types/rke_types.go")
	if err != nil {
		panic(err)
	}

	structMap := ExtractDoc(string(src))
	dumpAsJsonToFile("schemas/docMap.json", structMap)

	rkec := conf.RancherKubernetesEngineConfig{}

	jsonNameMap := make(map[string](map[string]string), 50)
	GenerateStructFieldNameMap(&jsonNameMap, reflect.TypeOf(rkec), false)
	dumpAsJsonToFile("schemas/jsonNames.json", jsonNameMap)

	yamlNameMap := make(map[string](map[string]string), 50)
	GenerateStructFieldNameMap(&yamlNameMap, reflect.TypeOf(rkec), true)
	dumpAsJsonToFile("schemas/yamlNames.json", yamlNameMap)

	r := jsonschema.Reflector{TypeMapper: typeAdjust}
	dumpAsJsonToFile("schemas/cluster.json", r.Reflect(&rkec))

	r.PreferYAMLSchema = true
	dumpAsJsonToFile("schemas/cluster.yml.json", r.Reflect(&rkec))

}
