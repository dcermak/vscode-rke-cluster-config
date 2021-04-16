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
	"fmt"
	"reflect"
	"testing"
)

var src = `
package main

type RestoreConfig struct {
	Restore      bool
	SnapshotName string
}

type RotateCertificates struct {
	// Rotate CA Certificates
	CACertificates bool
	// Services to rotate their certs
	Services []string
        ExpirationDate time.Date
        //   
        CaAuthority

 	// These options are NOT for configuring the Metrics-Server's addon template.
	// They are used to pass command args to the metric-server's deployment containers specifically.
	Options NumberStringMap
}
`

func TestDocExtraction(t *testing.T) {
	docMap := ExtractDoc(src)

	if len(docMap) != 1 {
		t.Errorf("Expected docMap to have length 1 but got %d", len(docMap))
	}

	rotateCaCertDoc, present := docMap["RotateCertificates"]
	if !present {
		t.Error("Expected to find key 'RotateCertificates' in the docMap")
	}

	if len(rotateCaCertDoc) != 3 {
		t.Errorf("Expected rotateCaCertDoc to have length 3, but got %d", len(rotateCaCertDoc))
	}

	for _, v := range [][2]string{
		{"CACertificates", "Rotate CA Certificates"},
		{"Services", "Services to rotate their certs"},
		{"Options", `These options are NOT for configuring the Metrics-Server's addon template.
They are used to pass command args to the metric-server's deployment containers specifically.`}} {
		fieldName := v[0]
		docString := v[1] + "\n"

		actualDoc, present := rotateCaCertDoc[fieldName]
		if !present {
			t.Errorf("key '%s' is not present in rotateCaCertDoc", fieldName)
		}
		if actualDoc != docString {
			t.Errorf("Expected docstring of '%s' to equal '%s', but got '%s'", fieldName, actualDoc, docString)
		}
	}
}

type Inner struct {
	WithoutTag   map[string]int
	WithTag      *bool      `yaml:"with_tag" json:"withTag"`
	SliceEntries []AsArray  `json:"Slice" yaml:"slice"`
	ArrayEntries [2]AsArray `json:"Array" yaml:"array"`
}

type EmbeddMe struct {
	Entry string
}

type InnerSecond struct {
	EmbeddMe
}

type AsArray struct {
	Foo *int
	Baz []string
}

type Outer struct {
	FieldWithJsonAndYaml string `yaml:"field_json" json:"FieldJson"`
	FieldWithJsonOnly    bool   `json:"JsonOnly"`
	FieldWithYamlOnly    int    `yaml:"yaml_only"`

	InnerElement *Inner      `yaml:"inner_1" json:"Inner1"`
	InnerSecond  InnerSecond `yaml:"inner_second" json:"innerSecond"`
}

func TestGenerateStructFieldNameMapJson(t *testing.T) {
	jsonMap := make(map[string](map[string]string))
	o := Outer{}
	GenerateStructFieldNameMap(&jsonMap, reflect.TypeOf(o), false)

	expected := map[string](map[string]string){
		"Outer": map[string]string{
			"FieldJson":   "FieldWithJsonAndYaml",
			"JsonOnly":    "FieldWithJsonOnly",
			"yaml_only":   "FieldWithYamlOnly",
			"Inner1":      "InnerElement",
			"innerSecond": "InnerSecond",
		},
		"Inner": map[string]string{
			"WithoutTag": "WithoutTag",
			"withTag":    "WithTag",
			"Array":      "ArrayEntries",
			"Slice":      "SliceEntries",
		},
		"AsArray":     map[string]string{"Foo": "Foo", "Baz": "Baz"},
		"InnerSecond": map[string]string{"Entry": "Entry"},
	}
	if !reflect.DeepEqual(jsonMap, expected) {
		t.Errorf("expected field map '%+v' to deep equal '%+v'", jsonMap, expected)
	}
}

func TestGenerateStructFieldNameMapYaml(t *testing.T) {
	yamlMap := make(map[string](map[string]string))
	o := Outer{}
	GenerateStructFieldNameMap(&yamlMap, reflect.TypeOf(o), true)

	expected := map[string](map[string]string){
		"Outer": map[string]string{
			"field_json":   "FieldWithJsonAndYaml",
			"JsonOnly":     "FieldWithJsonOnly",
			"yaml_only":    "FieldWithYamlOnly",
			"inner_1":      "InnerElement",
			"inner_second": "InnerSecond",
		},
		"AsArray": map[string]string{"Foo": "Foo", "Baz": "Baz"},
		"Inner": map[string]string{
			"WithoutTag": "WithoutTag",
			"with_tag":   "WithTag",
			"array":      "ArrayEntries",
			"slice":      "SliceEntries",
		},
		"InnerSecond": map[string]string{"Entry": "Entry"},
	}
	if !reflect.DeepEqual(yamlMap, expected) {
		fmt.Printf("%+v\n%+v\n", yamlMap, expected)
		t.Errorf("expected field map '%+v' to deep equal '%+v'", yamlMap, expected)
	}
}
