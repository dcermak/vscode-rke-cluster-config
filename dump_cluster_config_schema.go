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
	"io/ioutil"
	"reflect"

	"github.com/alecthomas/jsonschema"
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

func main() {
	rkec := conf.RancherKubernetesEngineConfig{}

	r := jsonschema.Reflector{TypeMapper: typeAdjust}
	schema, err := json.MarshalIndent(r.Reflect(&rkec), "", "  ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("schemas/cluster.json", schema, 0644)
	if err != nil {
		panic(err)
	}

	r.PreferYAMLSchema = true
	schema, err = json.MarshalIndent(r.Reflect(&rkec), "", "  ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("schemas/cluster.yml.json", schema, 0644)
	if err != nil {
		panic(err)
	}
}
