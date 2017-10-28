package service

import (
	"github.com/satori/go.uuid"
)

// Service provides the basic service interface implemented by every service.
type Service interface {
	// Name returns the name of the service
	Name() string

	// ID returns the id of the service
	ID() ID
}

// LifeCycle provides the service life cycle interface that services need to participate.
// The application framework will use this life cycle interface to control the services.
type LifeCycle interface {
	// OnConfigure defines the callback to configure the service and is called during
	// the application configuration phase.
	OnConfigure(interface{}) error

	// OnStart defines the callback to start the service and is called during the
	// application start phase.
	OnStart() error

	// OnStop defines the callback to stop the service and is called during the
	// application stop phase.
	OnStop() error

	// IsHealthy returns the health status of the service
	IsHealthy() bool
}

// ID type defines a unique service ID
type ID uuid.UUID

// BaseService provides a default implementation of Service interface that can be used
// by other service implementations.
type BaseService struct {
	name string
	id   ID
}

// NewBaseService returns an instance of BaseService
func NewBaseService(name string) *BaseService {
	s := &BaseService{name: name}
	s.id = ID(uuid.NewV4())
	return s
}

// Name returns the name of the service.
func (s *BaseService) Name() string {
	return s.name
}

// ID returns the id of the service.
func (s *BaseService) ID() ID {
	return s.id
}
