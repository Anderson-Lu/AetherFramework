package service

import (
	"fmt"
	"time"
)

type AetherTaskMonitor struct {
}

type AetherTaskReport struct {
	Id        string
	NodeIP    string
	StartTime int64
	EndTime   int64
	Cost      int64
	Status    AetherTaskStatus
	JobType   AetherJobType
	Crontab   string
}

func (s *AetherTaskMonitor) Report(report *AetherTaskReport) {

	var statuStr = "Unknown"
	switch report.Status {
	case AetherTaskStatusRunning:
		statuStr = "Running"
	case AetherTaskStatusFailed:
		statuStr = "Failed"
	case AetherTaskStatusSucceed:
		statuStr = "Succeed"
	}

	fmt.Printf("任务ID:%s,执行IP:%s,状态:%s 耗时: %ds \n",
		report.Id,
		report.NodeIP,
		statuStr,
		time.Now().Unix()-report.StartTime,
	)

}
