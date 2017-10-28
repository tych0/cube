package service

import (
	"github.com/satori/go.uuid"
)

type Service interface {
	Name() string
	ID() ServiceID
}

type LifeCycle interface {
	OnStart() error
	OnConfigure(interface{}) error
	OnStop() error
}

type ServiceID uuid.UUID

type BaseService struct {
	name string
	id   ServiceID
}

func NewBaseService(name string) *BaseService {
	s := &BaseService{name: name}
	s.id = ServiceID(uuid.NewV4())
	return s
}

func (s *BaseService) Name() string {
	return s.name
}

func (s *BaseService) ID() ServiceID {
	return s.id
}
