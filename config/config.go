package config

import "reflect"
import "fmt"

// Store provides a configuration store interface. Services can register their
// configuration types and can retrieve their configuration.
type Store interface {
	// Open creates the resources like db connections or files required by the store.
	Open() error

	// Close releases any underlying resources used by the store.
	Close()

	// Get returns the configuration for the specified key or error if the key is not
	// found in the store.
	Get(name string) (interface{}, error)

	// Registry returns the registry used by this store
	Registry() Registry
}

// Registry provides an interface for services to register their configuration types.
//
// Store implementations uses this interface to process configuration.
type Registry interface {
	// Register the configuration type with the provided name as the key.
	Register(name string, obj interface{}) error

	// Get returns the configuration type registered for the specified name.
	Get(name string) reflect.Type
}

type registry struct {
	types map[string]reflect.Type
}

// NewRegistry returns a new registry instance.
func NewRegistry() Registry {
	return &registry{map[string]reflect.Type{}}
}

func (c *registry) Register(name string, obj interface{}) error {
	if obj == nil {
		return fmt.Errorf("config type cannot be nil")
	}
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = reflect.ValueOf(obj).Elem().Type()
	}
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("config type must be a struct")
	}
	c.types[name] = t
	return nil
}

func (c *registry) Get(name string) reflect.Type {
	return c.types[name]
}
