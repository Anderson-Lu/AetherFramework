package service

import (
	"errors"
	"log"
	"poseidon/err"
)

type AetherServiceManager struct {
	services map[string]AetherService
	showLog  bool
}

func NewAetherServiceManager(showLog bool) *AetherServiceManager {
	return &AetherServiceManager{
		services: make(map[string]AetherService, 0),
		showLog:  showLog,
	}
}

func (s *AetherServiceManager) RegistService(name string, service AetherService) error {
	if _, ok := s.services[name]; ok {
		return errors.New(err.ErrServiceNameExist)
	}
	s.services[name] = service
	return nil
}

func (s *AetherServiceManager) GetService(serviceName string) AetherService {
	if srv, ok := s.services[serviceName]; ok {
		return srv
	}
	return nil
}

func (s *AetherServiceManager) Start() {
	s.initService()
}

func (s *AetherServiceManager) initService() {
	for k, v := range s.services {
		if e := v.Regist(); e != nil {
			s.logf("Regist service %s fail, reason: %s", k, e.Error())
		} else {
			s.logf("Regist service %s succeed", k)
		}
	}
}

func (s *AetherServiceManager) logf(format string, args ...interface{}) {

	if !s.showLog {
		return
	}

	format = "[AetherServiceManager] " + format + "\n"
	log.Printf(format, args...)
}
