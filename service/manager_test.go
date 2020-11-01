package service

import (
	"errors"
	"testing"
)

type ServiceT1 struct {
	a int
	b int
}

func (s *ServiceT1) Regist(params ...interface{}) error {
	s.a = 1
	s.b = 0
	return nil
}

func (s *ServiceT1) OnQuit() {
	s.b = 1
}

func (s *ServiceT1) Get() (int, int) {
	return s.a, s.b
}

type ServiceT2 struct {
	a int
	b int
}

func (s *ServiceT2) Regist(params ...interface{}) error {
	s.a = 1
	s.b = 0
	return errors.New("fail")
}

func (s *ServiceT2) OnQuit() {
	s.b = 1
}

func TestPosedonManager(t *testing.T) {
	serviceA := &ServiceT1{a: -1, b: -1}
	serviceB := &ServiceT2{}

	serviceManager := NewAetherServiceManager(true)
	serviceManager.RegistService("service-a", serviceA)
	serviceManager.RegistService("service-b", serviceB)
	serviceManager.RegistService("service-a", serviceA)

	serviceManager.Start()

	s := serviceManager.GetService("service-a")
	a, b := s.(*ServiceT1).Get()
	t.Log(a, b)

	non := serviceManager.GetService("service-not-exist")
	t.Log(non)
}
