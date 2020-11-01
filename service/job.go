package service

import (
	"errors"
	"fmt"
	"poseidon/err"
	"poseidon/util"

	"time"

	"sync"
)

var once sync.Once

type AetherJob struct {
	name              string             // Job name
	f                 AetherJobFunc      // Job exec func
	curTaskId         string             // Job id
	curTaskStatus     AetherTaskStatus   // Task status
	curTaskStartTime  int64              // Task start time
	curTaskFinishTime int64              // Task finish time
	curTaskIP         string             // Task IP
	taskMonitor       *AetherTaskMonitor // Task monitor
	taskOptions       *AetherJobOptions  // Job executing options
	jobStartTime      int64              // Job start time
	jobStatus         AetherJobStatus    // Job status
}

type AetherJobOptions struct {
	Crontab        string                // Crontab to run the job
	JobType        AetherJobType         // JobType
	JobLoopOption  *AetherJobLoopOption  // Loop Option
	MonitorOptions *AetherMonitorOptions // Monitor Options
}

type AetherJobLoopOption struct {
	RerunningWhenTick bool // Whether to start a new task immediately when the time is up and ignore the task in progress
}

func NewAetherJob(name string, job AetherJobFunc, options *AetherJobOptions) *AetherJob {
	return &AetherJob{
		name:        name,
		taskOptions: options,
		f:           job,
	}
}

func (s *AetherJob) Regist(args ...interface{}) error {
	if s.taskOptions == nil {
		return errors.New(err.ErrServiceBadConfiguration)
	}
	go once.Do(s.monitor)
	switch s.taskOptions.JobType {
	case AetherJobTypeLoop:
		if s.taskOptions.Crontab == "" {
			return errors.New(err.ErrServiceBadConfiguration)
		}
		if s.taskOptions.JobLoopOption == nil {
			s.taskOptions.JobLoopOption = &AetherJobLoopOption{}
		}
		s.jobStatus = AetherJobStatusRegisted
		util.CronSchaduler.AddAndStartJob(s.run, func(e error) {}, "", s.taskOptions.Crontab)
		return nil
	case AetherJobTypeOnce:
		go s.run()
		return nil
	}

	return errors.New(err.ErrServiceBadJobType)
}

func (s *AetherJob) configure(options *AetherJobOptions) error {

	if options == nil {
		return errors.New(err.ErrServiceBadConfiguration)
	}

	s.jobStatus = AetherJobStatusRegisted
	s.taskOptions = options
	go once.Do(s.monitor)

	return nil
}

func (s *AetherJob) StartTaskManual() {
	if s.curTaskStatus != AetherTaskStatusRunning {
		s.run()
	}
}

func (s *AetherJob) StartJobManual() {
	if s.jobStatus == AetherJobStatusStopped {
		s.Regist()
	}
}

func (s *AetherJob) StopJobManual() {
	s.jobStatus = AetherJobStatusStopped
}

func (s *AetherJob) run() {

	// check Job status
	if s.jobStatus == AetherJobStatusStopped || s.jobStatus == AetherJobStatusUnknown || s.jobStatus == AetherJobStatusCfgError {
		return
	}

	// check current task status
	if s.jobStatus == AetherJobStatusRunning {
		return
	}

	s.jobStatus = AetherJobStatusRunning
	s.curTaskStatus = AetherTaskStatusRunning
	s.curTaskStartTime = time.Now().Unix()
	s.curTaskId = fmt.Sprintf("%s_%d", s.name, time.Now().UnixNano()/1000000)
	s.curTaskIP, _ = util.ExternalIP()

	err := s.f()
	if err != nil {
		s.curTaskStatus = AetherTaskStatusFailed
	} else {
		s.curTaskStatus = AetherTaskStatusSucceed
	}
}

func (s *AetherJob) monitor() {
	ticker := time.NewTicker(time.Duration(s.taskOptions.MonitorOptions.StatusUpdateInterval) * time.Second)
	for range ticker.C {
		if s.jobStatus != AetherJobStatusRunning {
			continue
		}
		s.taskMonitor.Report(&AetherTaskReport{
			Id:        s.curTaskId,
			NodeIP:    s.curTaskIP,
			StartTime: s.curTaskStartTime,
			EndTime:   s.curTaskFinishTime,
			Cost:      time.Now().Unix() - s.curTaskStartTime,
			Status:    s.curTaskStatus,
			JobType:   s.taskOptions.JobType,
			Crontab:   s.taskOptions.Crontab,
		})
		if s.curTaskStatus == AetherTaskStatusFailed || s.curTaskStatus == AetherTaskStatusSucceed {
			s.jobStatus = AetherJobStatusFree
			s.curTaskFinishTime = time.Now().Unix()
		}
		switch s.taskOptions.JobType {
		case AetherJobTypeLoop:

		}
	}
}
