package main

import (
	"errors"
	"fmt"
	"poseidon/service"
	"time"
)

func main() {
	options := &service.AetherJobOptions{
		Crontab: "*/1 * * * * ? ",
		JobType: service.AetherJobTypeLoop,
		MonitorOptions: &service.AetherMonitorOptions{
			StatusUpdateInterval: 2,
		},
		JobLoopOption: &service.AetherJobLoopOption{
			RerunningWhenTick: false,
		},
	}

	jobService := service.NewAetherJob("anderson", func() error {
		// cost := rand.Intn(20)
		cost := 15
		for i := 0; i < cost; i++ {
			fmt.Println(fmt.Sprintf("--->%d", i))
			time.Sleep(time.Second)
		}
		if cost%2 == 0 {
			return errors.New("bad")
		}
		return nil
	}, options)

	// go func() {
	// 	time.Sleep(time.Second * 10)
	// 	go jobService.StartTaskManual()
	// 	time.Sleep((time.Second * 5))
	// 	jobService.StopJobManual()
	// }()

	manager := service.NewAetherServiceManager(true)
	manager.RegistService("anderson", jobService)
	manager.Start()

	select {}

}
