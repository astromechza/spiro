package templatefactory

import (
	"bytes"
	"fmt"
	"html/template"
	"reflect"
	"strings"
)

const SpecialDelimetersKey = "_spiro_delimeters_"

type TemplateFactory struct {
	funcMap    template.FuncMap
	startDelim string
	endDelim   string
	spec       *map[string]interface{}
}

func NewTemplateFactory() *TemplateFactory {
	return &TemplateFactory{
		funcMap:    make(template.FuncMap),
		startDelim: "{{",
		endDelim:   "}}",
	}
}

func (f *TemplateFactory) SetSpec(in *map[string]interface{}) error {
	f.spec = in
	if delims, ok := (*in)[SpecialDelimetersKey]; ok {
		s := reflect.ValueOf(delims)
		if s.Kind() == reflect.Slice && s.Len() == 2 {
			var sok, eok bool
			f.startDelim, sok = s.Index(0).Interface().(string)
			f.endDelim, eok = s.Index(1).Interface().(string)
			if sok && eok {
				return nil
			}
		}
		return fmt.Errorf("Overidding template delimeters with '%s' requires an array of two strings", SpecialDelimetersKey)
	}
	return nil
}

func (f *TemplateFactory) StringContainsTemplating(in string) bool {
	return strings.Contains(in, f.startDelim) && strings.Contains(in, f.endDelim)
}

func (f *TemplateFactory) RegisterTemplateFunction(name string, function interface{}) {
	f.funcMap[name] = function
}

func (f *TemplateFactory) Render(templateString string) (string, error) {
	t := template.New("").Option("missingkey=error").Funcs(f.funcMap).Delims(f.startDelim, f.endDelim)
	if t, err := t.Parse(templateString); err != nil {
		return "", err
	} else {
		var buf bytes.Buffer
		err := t.Execute(&buf, f.spec)
		return buf.String(), err
	}
}
