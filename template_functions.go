package main

import (
	"encoding/json"
	"html/template"
	"reflect"
)

func Jsonify(in map[interface{}]interface{}) string {
	working := make(map[string]interface{})
	s := reflect.ValueOf(in)
	for _, k := range s.MapKeys() {
		working[k.Interface().(string)] = s.MapIndex(k).Interface()
	}
	sb, err := json.Marshal(working)
	if err != nil {
		panic(err)
	}
	return string(sb)
}

func JsonifyIndent(in map[interface{}]interface{}) string {
	working := make(map[string]interface{})
	s := reflect.ValueOf(in)
	for _, k := range s.MapKeys() {
		working[k.Interface().(string)] = s.MapIndex(k).Interface()
	}
	sb, err := json.MarshalIndent(working, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(sb)
}

func Unescape(in string) interface{} {
	return template.HTML(in)
}
