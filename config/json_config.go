package config

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

type jsonStore struct {
	reg Registry
	r   io.Reader
	kv  map[string]interface{}
	kb  map[string][]byte
}

// NewJSONStore returns a config store backed by a JSON stream.
//
// The first level keys in the JSON stream match the names registered in the Registry
// and the values must be decodeable into the types registered against that name.
func NewJSONStore(r io.Reader, registry Registry) Store {
	return &jsonStore{
		reg: registry,
		r:   r,
		kv:  map[string]interface{}{},
		kb:  map[string][]byte{},
	}
}

func (j *jsonStore) Open() error {
	d := json.NewDecoder(j.r)
	for {
		data := map[string]*cfgData{}
		if err := d.Decode(&data); err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		for k, v := range data {
			t := j.reg.Get(k)
			if t != nil {
				val := reflect.New(t).Interface()
				if e := json.Unmarshal(v.b, val); e != nil {
					return e
				}
				j.kv[k] = val
			} else {
				// We dont know the key yet so capture the data for future processing
				j.kb[k] = v.b
			}
		}
	}
}

func (j *jsonStore) Close() {
	// NOOP
}

func (j *jsonStore) Get(name string) (interface{}, error) {
	if v, ok := j.kv[name]; ok {
		return v, nil
	} else if b, ok := j.kb[name]; ok {
		if t := j.reg.Get(name); t != nil {
			// process the key now and cache the value
			val := reflect.New(t).Interface()
			if e := json.Unmarshal(b, val); e != nil {
				// Bad buffer for the current type but lets keep it around
				// in case the registry is modified with a new type
				// and we can process it in future Get calls
				return nil, e
			}
			// Cache the value and discard the buffer
			j.kv[name] = val
			delete(j.kb, name)
			return val, nil
		}
		// Key is present in the store but not registered
		return nil, fmt.Errorf("%s key not registered", name)
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
