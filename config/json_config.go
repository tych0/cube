package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
)

type jsonStore struct {
	reg      Registry
	fileName string
	kv       map[string]interface{}
}

// NewJSONStore returns a config store backed by a JSON file.
//
// The first level keys in the JSON file must match the names registered in the Registry
// and the values must be decodeable into the types registered against that name.
func NewJSONStore(fileName string, registry Registry) Store {
	return &jsonStore{reg: registry, fileName: fileName, kv: map[string]interface{}{}}
}

func (j *jsonStore) Load() error {
	data, e := ioutil.ReadFile(j.fileName)
	if e != nil {
		return e
	}
	return j.processData(data)
}

func (j *jsonStore) Get(name string) (interface{}, error) {
	if v, ok := j.kv[name]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("%s key not found", name)
}

func (j *jsonStore) Registry() Registry {
	return j.reg
}

type cfgData struct {
	b []byte
}

func (d *cfgData) UnmarshalJSON(b []byte) error {
	d.b = b
	return nil
}

func (j *jsonStore) processData(data []byte) error {
	jsonData := map[string]*cfgData{}
	if e := json.Unmarshal(data, &jsonData); e != nil {
		return e
	}
	for k, v := range jsonData {
		t := j.reg.Get(k)
		if t != nil {
			val := reflect.New(t).Interface()
			if e := json.Unmarshal(v.b, val); e != nil {
				return e
			}
			j.kv[k] = val
		} // else ignore the unregistered key
	}
	return nil
}
