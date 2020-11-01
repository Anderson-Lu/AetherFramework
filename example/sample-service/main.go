package main

import (
	"log"
	"poseidon/service"
	"time"
)

func main() {
	sample := &SampleService{}

	manager := service.NewAetherServiceManager(true)
	manager.RegistService("sample", sample)
	manager.Start()

	sampleInstance := manager.GetService("sample")
	if sampleInstance != nil {
		sampleInstance.(*SampleService).MyFunc()
	}

	time.Sleep(time.Second * 4)
}

type SampleService struct {
}

func (s *SampleService) Regist(params ...interface{}) error {
	return nil
}

func (s *SampleService) OnQuit() {
	log.Println("quit...")
}

func (s *SampleService) MyFunc() {
	log.Println("my func")
}
